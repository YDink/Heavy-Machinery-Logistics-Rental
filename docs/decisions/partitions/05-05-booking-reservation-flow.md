# Partition Decisions: 05-05 — Booking Reservation Lifecycle & Payment Hold Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-05`.

- **Governing System ADRs:** `ADR-001`, `ADR-002`, `ADR-004`, `ADR-007`, `ADR-010`
- **Subsystem Scope:** backend, database
- **Requirement Targets:** BK-01-BE to BK-24-BE

---

## D-05-05-001: In-Transaction Synchronous Financial Atomicity (BK-17-BE) for Ledger Mutations

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [BK-17-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
When a booking transitions to CONFIRMED following owner verification of direct bank transfer evidence, marketplace commission must be debited from the owner pre-funded wallet and recorded in the double-entry ledger. If ledger recording fails or is deferred asynchronously, uncollected platform commissions and untracked financial drift occur.

### 2. Considered Alternatives
1. **Option A:** Synchronous in-transaction execution: booking state update, owner wallet debit, and double-entry ledger entry committed atomically within a single pgx.Tx
   - *Pros:* Provable zero financial drift; absolute ledger invariants; instant rollback on insufficient wallet funds.
   - *Cons:* Transaction holds lock on owner_wallets row for ~5-10ms.
2. **Option B:** Asynchronous event bus: booking marks CONFIRMED immediately, and an event consumer debits owner wallet in background
   - *Pros:* Slightly lower HTTP handler latency (~5ms faster).
   - *Cons:* High commercial default risk; if owner wallet has insufficient balance, booking is already confirmed and platform commission is lost.
3. **Option C:** Periodic batch reconciliation cron job debiting commissions nightly
   - *Pros:* Decoupled billing operations.
   - *Cons:* Severe credit risk; owners can withdraw equipment before daily batch runs.

### 3. Chosen Solution & Technical Rationale
Option A: In-Transaction Synchronous Financial Atomicity (BK-17-BE). When owner confirms receipt of bank wire, the handler executes within pgx.Tx: 1) Verifies owner wallet balance/credit ceiling; 2) Decrements wallet balance; 3) Writes double-entry ledger rows; 4) Transitions booking to CONFIRMED. If any step fails, the entire transaction aborts.

### 4. Rejected Alternatives
- Option B: Rejected because asynchronous commission debits create uncollectible debts when owner wallets lack funds.
- Option C: Rejected due to catastrophic commercial credit risk.

### 5. Consequences & Invariants
- **Invariants Maintained:** Zero booking transitions to CONFIRMED without an atomic, balanced double-entry ledger debit from OwnerWallet.
- **Operational Trade-offs:** Requires fleet owners to maintain pre-funded wallet balances or authorized credit limits.

### 6. Verification
- **Test File:** `backend/internal/modules/booking/application/confirm_booking_test.go`
- **Assertion:** TestConfirmBooking_RollbackOnWalletFailure asserts that insufficient owner wallet balance rolls back booking confirmation to PAYMENT_EVIDENCE_UPLOADED.

### 7. Revisit Trigger
- Revisit if migration to automated payment escrow gateways in Phase 3 changes the fund flow topology.

---

## D-05-05-002: 60-Minute Ephemeral Valkey Mutex Hold Locks and OCC Versioning

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [BK-04-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
When a contractor taps Book Now, they require a guaranteed window (default 60 minutes) to execute the direct interbank wire transfer and upload the receipt slip without competing contractors sniping the reservation.

### 2. Considered Alternatives
1. **Option A:** Ephemeral Valkey mutual exclusion hold lock (SET NX EX 3600) paired with database versioned OCC (version column)
   - *Pros:* Zero database locks held while user completes manual bank app transfer; automatic expiry cleanup; robust concurrency.
   - *Cons:* Requires Valkey memory monitoring to avoid key eviction.
2. **Option B:** Database SELECT ... FOR UPDATE row lock held across the 60-minute window
   - *Pros:* Native relational locking.
   - *Cons:* Catastrophic connection pool exhaustion; platform outage within minutes.
3. **Option C:** Optimistic UI without locks; first slip confirmed by owner wins the booking
   - *Pros:* Zero backend locking infrastructure.
   - *Cons:* Multiple contractors wire hundreds of thousands of Birr for the same machine, creating severe legal disputes.

### 3. Chosen Solution & Technical Rationale
Option A: Ephemeral Valkey hold lock (SET lock:fleet:hold:<id> <booking_id> NX EX 3600) combined with database OCC versioning. The fleet availability calendar displays PENDING_PAYMENT during this window. If slip is uploaded, lock extends; if timer expires, lock is auto-deleted.

### 4. Rejected Alternatives
- Option B: Rejected because open database transactions across human bank transfers destroy server scalability.
- Option C: Rejected because concurrent payments for high-value machinery create unacceptable customer financial disputes.

### 5. Consequences & Invariants
- **Invariants Maintained:** A fleet asset cannot be placed on payment hold by more than one user for overlapping dates.
- **Operational Trade-offs:** If Valkey crashes without replication, hold locks must recover from database PENDING_PAYMENT records.

### 6. Verification
- **Test File:** `backend/internal/modules/booking/application/booking_hold_test.go`
- **Assertion:** TestBookingHold_ConcurrentAttemptBlocked asserts that a second user attempting to hold the same fleet within 60 minutes receives HTTP 409 Conflict.

### 7. Revisit Trigger
- Revisit if instant mobile money payment gateways reduce required hold windows from 60 minutes to 300 seconds.

---

