# Partition Decisions: 05-10 — Notifications, In-App Message Center & Alerting Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-10`.

- **Governing System ADRs:** `ADR-007`, `ADR-009`
- **Subsystem Scope:** backend, mobile
- **Requirement Targets:** NOTIF-01-BE to NOTIF-14-BE

---

## D-05-10-001: Multi-Channel Dispatch (FCM Push, SMS Gateway, In-App Center) via Transactional Outbox

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [NOTIF-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Critical operational alerts (booking confirmation, inspection overdue, payment timeout) must reach contractors and operators regardless of device state. Direct inline API calls to external push and SMS providers slow down HTTP requests and fail permanently if external gateways time out.

### 2. Considered Alternatives
1. **Option A:** PostgreSQL transactional outbox table (notification_events) committed atomically with domain states and dispatched by background workers
   - *Pros:* Zero message loss; HTTP handlers respond immediately; automated retry with exponential backoff on gateway failure.
   - *Cons:* Introduces small sub-second delivery latency for background processing.
2. **Option B:** Synchronous inline HTTP calls to FCM and SMS gateways inside request handlers
   - *Pros:* Immediate notification attempt.
   - *Cons:* Slows down API responses; network timeouts cause failed HTTP transactions and dropped alerts.
3. **Option C:** External Redis Pub/Sub without persistent database outbox
   - *Pros:* Fast in-memory messaging.
   - *Cons:* Messages are dropped if the server restarts or worker crashes before delivery.

### 3. Chosen Solution & Technical Rationale
Option A: PostgreSQL transactional outbox table with background worker dispatch. Domain interactors insert notification events into notification_events inside pgx.Tx. A dedicated worker flushes events via FCM (data + notification payload), fallback SMS, and in-app message feeds.

### 4. Rejected Alternatives
- Option B: Rejected because third-party telco API latencies must never block primary database transactions.
- Option C: Rejected because ephemeral in-memory queues fail the zero-alert-loss invariant for legal dispute notices.

### 5. Consequences & Invariants
- **Invariants Maintained:** Every critical domain notification is persisted in notification_events within the primary transaction.
- **Operational Trade-offs:** Outbox table requires periodic pruning of delivered messages older than 90 days.

### 6. Verification
- **Test File:** `backend/internal/modules/notification/infrastructure/outbox_worker_test.go`
- **Assertion:** TestNotificationOutbox_DispatchesWithRetry asserts that a failed SMS attempt is retried with backoff and logged.

### 7. Revisit Trigger
- Revisit if notification throughput exceeds 20,000 alerts/minute, necessitating a dedicated distributed queue.

---

## D-05-10-002: In-App Notification Center with Cursor Keyset Pagination and Read State Sync

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [NOTIF-06-BE]
- **Subsystem:** [backend, mobile]

### 1. Problem Statement & Context
Users frequently miss push notifications due to battery saver restrictions or device sleep modes. The in-app notification inbox must provide an unread counter, fast cursor pagination, and reliable read receipt synchronization across multiple user devices.

### 2. Considered Alternatives
1. **Option A:** Relational notification archive with keyset cursor pagination (created_at, id) and optimistic read receipt batch updates
   - *Pros:* Fast sub-10ms inbox loading; consistent infinite scroll; optimistic UI updates in Flutter client.
   - *Cons:* Requires unread count cache invalidation in Valkey.
2. **Option B:** Offset/limit pagination on notification table
   - *Pros:* Simple page-number queries.
   - *Cons:* Slow queries as inbox depth grows; duplicate or skipped notifications when new alerts arrive during browsing.
3. **Option C:** Local mobile-only notification storage without server sync
   - *Pros:* Zero server storage overhead.
   - *Cons:* Unread state is desynchronized across web admin and mobile devices.

### 3. Chosen Solution & Technical Rationale
Option A: Keyset cursor pagination on notification records paired with an optimistic read-state sync endpoint. Unread notification counters are cached in Valkey (user:unread_notifs:<id>) for instantaneous badge rendering.

### 4. Rejected Alternatives
- Option B: Rejected due to page-drift anomalies on live notification feeds.
- Option C: Rejected because contractors switch between office web browsers and field mobile phones.

### 5. Consequences & Invariants
- **Invariants Maintained:** Marking a notification as read decrements the user unread counter atomically.
- **Operational Trade-offs:** Requires maintaining Valkey counter sync with actual unread database records.

### 6. Verification
- **Test File:** `backend/internal/modules/notification/application/inbox_test.go`
- **Assertion:** TestNotificationInbox_CursorPagination asserts stable pagination when new notifications are inserted concurrently.

### 7. Revisit Trigger
- Revisit if user inbox archives exceed 10,000 notifications per user, requiring cold-storage archival.

---

