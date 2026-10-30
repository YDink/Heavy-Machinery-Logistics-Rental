# Partition Decisions: 08-01 — Canonical 20-State Booking Machine & Concurrency Engine Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `08-01`.

- **Governing System ADRs:** `ADR-001`, `ADR-004`
- **Subsystem Scope:** backend
- **Requirement Targets:** BK-01-BE to BK-24-BE, SRS-SEC-08

---

## D-08-01-001: Canonical 20-State Booking Transition Table with Strict Non-Bypassable Invariants

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [BK-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
The rental lifecycle spans 20 discrete operational states (PENDING_PAYMENT, PAYMENT_EVIDENCE_UPLOADED, CONFIRMED, ACTIVE, RETURN_INITIATED, DISPUTED_HOUR_METER, etc.). Ad-hoc status updates in arbitrary handlers cause illegal state bypasses, uncollected commissions, and corrupted rental lifecycles.

### 2. Considered Alternatives
1. **Option A:** Centralized state transition table in domain core rejecting any status mutation not explicitly registered in the allowed transition graph
   - *Pros:* Impossible to bypass business invariants; single source of truth; compile-time and runtime transition verification.
   - *Cons:* Adding a new valid transition requires updating the core state machine definition.
2. **Option B:** Allowing arbitrary status updates in any handler via UPDATE bookings SET status = $1
   - *Pros:* Fast and flexible for developers.
   - *Cons:* Catastrophic state corruption; handlers transition directly from CANCELLED to CONFIRMED or skip mandatory handover inspections.
3. **Option C:** Database CHECK constraint listing all valid transitions
   - *Pros:* Relational enforcement.
   - *Cons:* Cannot execute rich domain side-effects (wallet debits, ledger logging, outbox events) within database constraints.

### 3. Chosen Solution & Technical Rationale
Option A: Centralized state transition table in internal/modules/booking/domain/state_machine.go. All status mutations must call booking.TransitionTo(newState). If the transition is illegal, it returns an RFC 9457 HTTP 409 Conflict (BOOKING_INVALID_STATE_TRANSITION) without executing database updates.

### 4. Rejected Alternatives
- Option B: Rejected because unrestricted state mutations create severe business logic bugs and financial fraud vulnerabilities.
- Option C: Rejected because database check constraints cannot orchestrate cross-module transactional side-effects.

### 5. Consequences & Invariants
- **Invariants Maintained:** Zero booking records may transition between states without passing through the canonical state machine validation graph.
- **Operational Trade-offs:** Developers must define and test all state transitions in the domain package.

### 6. Verification
- **Test File:** `backend/internal/modules/booking/domain/state_machine_test.go`
- **Assertion:** TestStateMachine_RejectsIllegalTransition asserts transitioning from CANCELLED to CONFIRMED returns ErrInvalidStateTransition.

### 7. Revisit Trigger
- Revisit if dynamic workflow state machine graphs configurable by enterprise tenants are required in Phase 3.

---

## D-08-01-002: Event-Driven State Transition Audit Journal with Full OCC Version Validation

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [BK-03-BE]
- **Subsystem:** [database, backend]

### 1. Problem Statement & Context
State transitions must be completely auditable for legal dispute resolution and post-incident forensics. Concurrent state mutations from multiple actors (e.g. renter cancels while owner confirms) must be detected and resolved without data corruption.

### 2. Considered Alternatives
1. **Option A:** Versioned Optimistic Concurrency Control (version column) paired with immutable booking_status_histories audit logging on every transition
   - *Pros:* Prevents race-condition overwrites; immutable forensic timeline of who transitioned what state, when, and from what IP.
   - *Cons:* Requires inserting an audit row on every status change.
2. **Option B:** Updating updated_at timestamp on bookings table alone without history rows
   - *Pros:* Zero historical storage overhead.
   - *Cons:* Zero auditability; cannot prove who approved a payment or when an inspection took place.
3. **Option C:** Pessimistic row locking across all state transitions
   - *Pros:* Strict sequential updates.
   - *Cons:* High lock contention and deadlock risks during high-frequency status queries.

### 3. Chosen Solution & Technical Rationale
Option A: Versioned OCC with booking_status_histories audit logging. Every transition executes: UPDATE bookings SET status = $1, version = version + 1 WHERE id = $2 AND version = $3. Inside the same transaction, a row is inserted into booking_status_histories capturing previous_status, new_status, actor_id, and reason.

### 4. Rejected Alternatives
- Option B: Rejected because lack of historical state logs makes resolving customer disputes legally impossible.
- Option C: Rejected due to database lock contention and connection pool exhaustion.

### 5. Consequences & Invariants
- **Invariants Maintained:** Every status transition inserts an immutable row into booking_status_histories within the same transaction.
- **Operational Trade-offs:** Status history table grows linearly with platform rental transaction volume.

### 6. Verification
- **Test File:** `backend/internal/modules/booking/infrastructure/status_history_test.go`
- **Assertion:** TestStatusHistory_RecordsAuditRow asserts valid transition writes matching row in booking_status_histories.

### 7. Revisit Trigger
- Revisit if regulatory mandates require cryptographic block-chain hashing of state history records.

---

