# ADR-004: Tactical In-Process CQRS and Versioned Optimistic Concurrency Control

- **Status:** Accepted
- **Date:** 2026-10-26
- **Decision Owner:** Dink Machinery Architecture Review Board
- **Affected Partitions:** 05-04, 05-05, 05-07, 05-21, 08-booking-state-machine
- **Related Requirements:** BK-04-BE, BK-06-BE, SRCH-01-BE

---

## 1. Context & Problem Statement

High-demand heavy machinery listings in major regional corridors (Addis Ababa, Awash, Mombasa) experience concurrent reservation attempts from competing contractors. Furthermore, complex spatial catalog search queries demand read-optimized performance that differs sharply from transactional write paths. Long-lived pessimistic database locks exhaust connection pools, while naive updates risk lost-update anomalies.

---

## 2. The 5-Point Qualification Evaluation

- [x] **Alternative Test:** Evaluated in-process CQRS with versioned OCC vs. distributed pessimistic locking vs. full event-sourced CQRS.
- [x] **Surprise Test:** Rejects distributed locks during user workflows; combines ephemeral Valkey hold locks with database OCC.
- [x] **Constraint Test:** Prevents database connection pool exhaustion while upholding zero double-booking invariants.
- [x] **Irreversibility Test:** Fundamental data schema modification requiring `version` columns on all mutable entity tables.
- [x] **Scope Test:** Governs booking state transitions, fleet availability updates, wallet transactions, and catalog read models.

---

## 3. Considered Options

### Option 1: In-Process CQRS with Versioned OCC and Ephemeral Valkey Hold Locks
- **Summary:** Separate command and query paths within the monolith. All mutable entity tables include a `version INT NOT NULL DEFAULT 1` column. Writes execute atomic OCC updates (`WHERE id = $1 AND version = $2`). Booking availability holds use an ephemeral 60-minute Valkey mutex (`SET NX EX 3600`).
- **Pros:**
  - Zero long-lived database transactions during manual user actions (e.g. bank slip upload).
  - High-throughput read queries execute without locking contention.
  - Provable prevention of lost updates and double-booking.
- **Cons:**
  - Requires clients to handle HTTP 409 Conflict retries.

### Option 2: Distributed Pessimistic Locking (SELECT ... FOR UPDATE across Booking Flow)
- **Summary:** Database rows locked from reservation initiation until bank slip upload and owner confirmation.
- **Pros:**
  - Simple conceptual model; transactions fail immediately if row is locked.
- **Cons:**
  - Holding database locks for up to 60 minutes exhausts connection pools and causes platform-wide outage.

### Option 3: Full Event Sourcing with Separate Read Databases
- **Summary:** Immutable event stream with projection workers populating a dedicated read database (e.g. Elasticsearch).
- **Pros:**
  - Complete historical replayability and independent read scaling.
- **Cons:**
  - Eventual consistency lag on availability checks causes false-negative double bookings; heavy operational overhead.

---

## 4. Chosen Solution & Architectural Rationale

We adopt **Option 1: In-Process CQRS with Versioned OCC and Ephemeral Valkey Hold Locks**.

1. **In-Process CQRS:** Commands execute business logic and state machine transitions through command handlers returning `Result[T, error]`. Queries bypass domain aggregates and execute read-optimized SQL queries returning tailored DTOs.
2. **Versioned OCC (`version` column):**
   ```sql
   UPDATE bookings 
   SET status = $1, version = version + 1, updated_at = NOW() 
   WHERE id = $2 AND version = $3;
   ```
   If zero rows are updated, the handler aborts and returns an RFC 9457 HTTP 409 Conflict (`BOOKING_CONCURRENT_MODIFICATION`).
3. **Ephemeral Reservation Holds:** During payment slip upload, availability is held via Valkey (`SET lock:fleet:hold:<id> <booking_id> NX EX 3600`). If payment is not submitted within 60 minutes, the hold expires automatically without database lock contention.

---

## 5. Rejected Alternatives & Technical Justification

- **Option 2 (Pessimistic Locking):** Rejected because human bank slip transfer verification takes up to 60 minutes; holding open database transactions across network drops rapidly exhausts the connection pool.
- **Option 3 (Full Event Sourcing):** Rejected due to the cognitive complexity, projection lag, and operational maintenance overhead of dual databases in emerging-market infrastructure.

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - Connection pool utilization remains low ($< 15\%$) even under heavy traffic spikes.
  - Immediate, deterministic conflict detection.
  - Queries execute with maximum read performance via specialized PostGIS/JSONB indexes.
- **Negative / Technical Debt:**
  - Presentation layers (Flutter and Next.js) must present retry prompts when receiving 409 Conflict responses.
- **Operational / Observability Impacts:**
  - OCC conflict rates are instrumented and monitored via Prometheus counter `dink_occ_conflicts_total`.

---

## 7. Verification Strategy & Invariant Tests

- **Concurrency Double-Booking Test:** `backend/internal/modules/booking/application/booking_concurrency_test.go` executes 100 concurrent goroutines attempting to reserve the same fleet item; asserts exactly 1 succeeds and 99 fail with OCC conflicts.
- **OCC Schema Invariant Test:** CI test inspects database migrations to verify all mutable tables declare `version INT NOT NULL DEFAULT 1`.

---

## 8. Revisit Conditions

- Write conflict rate exceeds 5% of all booking attempts, indicating extreme listing contention requiring virtual queuing.
