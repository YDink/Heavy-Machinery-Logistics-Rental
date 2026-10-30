# ADR-007: In-Process Typed Domain Event Bus with Transactional Outbox for MVP

- **Status:** Accepted
- **Date:** 2026-10-26
- **Decision Owner:** Dink Machinery Architecture Review Board
- **Affected Partitions:** 05-05, 05-08, 05-10, 08-booking-state-machine
- **Related Requirements:** SYS-02-BE, NOTIF-01-BE, BK-17-BE

---

## 1. Context & Problem Statement

When booking status transitions occur (e.g. `PAYMENT_EVIDENCE_UPLOADED`, `CONFIRMED`), downstream side-effects must trigger: dispatching push notifications, updating search cache availability, sending SMS alerts, and queuing bank transfer reconciliations. Introducing an external distributed message broker (Kafka, RabbitMQ) for the MVP adds significant deployment, monitoring, and operational complexity to our Tier III on-premise infrastructure.

---

## 2. The 5-Point Qualification Evaluation

- [x] **Alternative Test:** Evaluated in-process typed event bus with transactional outbox vs. external Apache Kafka cluster vs. Valkey Pub/Sub.
- [x] **Surprise Test:** Side-effects run asynchronously within the Go process or via transactional outbox rather than an external message broker.
- [x] **Constraint Test:** Zero message loss during process restarts; critical notification events cannot be dropped.
- [x] **Irreversibility Test:** Standardizes event publication and consumption interfaces across all backend modules.
- [x] **Scope Test:** Governs cross-module notification dispatch, analytics ingestion, and background job queuing.

---

## 3. Considered Options

### Option 1: In-Process Typed Go Event Bus with PostgreSQL Transactional Outbox
- **Summary:** Domain modules publish typed Go event structs to an in-memory event bus (`eventbus.Bus`). For events requiring guaranteed delivery (e.g. notifications, ledger webhooks), events are written to a `notification_events` outbox table within the primary database transaction (`pgx.Tx`) and processed by a local background worker.
- **Pros:**
  - Zero external infrastructure dependency (no Kafka/RabbitMQ clusters to manage).
  - Guaranteed delivery: outbox entries commit atomically with the domain state mutation.
  - Strongly typed event subscribers with zero network serialization overhead.
- **Cons:**
  - Event processing is bound to the single monolithic server or requires database polling workers.

### Option 2: Apache Kafka Cluster
- **Summary:** Deploy a 3-node Apache Kafka cluster for distributed event streaming.
- **Pros:**
  - High event throughput and independent horizontal scaling of event consumer groups.
- **Cons:**
  - Extreme operational complexity, JVM memory footprint, Zookeeper/KRaft maintenance, and high bare-metal infrastructure costs in Addis Ababa.

### Option 3: Valkey Pub/Sub Alone
- **Summary:** Publish domain events directly to Valkey Pub/Sub channels.
- **Pros:**
  - Low latency and leverages our existing Valkey container.
- **Cons:**
  - Valkey Pub/Sub is fire-and-forget; disconnected consumers or server restarts result in permanent message loss.

---

## 4. Chosen Solution & Architectural Rationale

We adopt **Option 1: In-Process Typed Go Event Bus with PostgreSQL Transactional Outbox**.

1. **Lightweight In-Process Event Bus:** Module interactors publish domain events to an in-memory event bus (`eventbus.Bus`). Subscribers register typed handler callbacks:
   ```go
   type EventBus interface {
       Subscribe(eventType string, handler EventHandler)
       Publish(ctx context.Context, event Event) error
   }
   ```
2. **Guaranteed Transactional Outbox:** Critical asynchronous events (such as `BookingConfirmedEvent`) write a row to `notification_events` inside the same `pgx.Tx` that updates the booking. A background worker picks up unprocessed events using `SELECT ... FOR UPDATE SKIP LOCKED` and dispatches FCM/SMS alerts with exponential backoff retries.
3. **Future Portability:** Because handlers implement a clean `EventHandler` interface, the underlying bus can be swapped to an external broker in Phase 3 without altering domain application code.

---

## 5. Rejected Alternatives & Technical Justification

- **Option 2 (Apache Kafka):** Rejected due to massive operational overhead, high RAM requirements, and unnecessary complexity for MVP transaction volumes ($< 500\text{ events/sec}$).
- **Option 3 (Valkey Pub/Sub):** Rejected because fire-and-forget messaging does not satisfy the zero-message-loss invariant required for financial confirmation alerts and legal dispute notices.

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - Minimal operational footprint: no extra broker services to monitor or restart.
  - Atomic transaction guarantees: an event is only persisted if the business state transition succeeds.
  - Sub-millisecond local event processing.
- **Negative / Technical Debt:**
  - High outbox polling frequency can create light database load.
- **Operational / Observability Impacts:**
  - Outbox processing lag is tracked via Prometheus gauge `dink_outbox_queue_depth`.

---

## 7. Verification Strategy & Invariant Tests

- **Outbox Atomicity Test:** `backend/internal/modules/booking/infrastructure/outbox_test.go` verifies that rolling back a booking transaction guarantees zero outbox rows are committed.
- **Dead-Letter Recovery Test:** Asserts that after 5 failed FCM dispatch attempts, the event is marked `FAILED` with error details captured in `notification_events`.

---

## 8. Revisit Conditions

- Platform event volume exceeds 5,000 events/second or background worker processing begins to saturate primary PostgreSQL I/O.
