# Partition Decisions: 07-07 — Phase 3 White-Label Partner API & Enterprise Integration Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `07-07`.

- **Governing System ADRs:** `ADR-002`, `ADR-009`, `ADR-014`
- **Subsystem Scope:** backend
- **Requirement Targets:** API-WHT-01-BE to API-WHT-08-BE

---

## D-07-07-001: Scoped Enterprise Partner API Gateway with Mutual TLS (mTLS) and Rate-Limiting Tiers

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [API-WHT-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Large civil engineering firms, mining conglomerates, and international construction ERP systems require programmatic API access to search fleets, initiate bookings, and sync invoices directly from their enterprise software without using web or mobile apps.

### 2. Considered Alternatives
1. **Option A:** Dedicated partner gateway (/api/v1/partner) with Mutual TLS (mTLS), scoped API keys, granular IP whitelisting, and strict tier-based token bucket rate limits
   - *Pros:* Enterprise-grade security; zero shared credentials; granular API scope isolation (e.g. read:fleets, write:bookings).
   - *Cons:* Requires certificate authority management for partner mTLS client certs.
2. **Option B:** Allowing partners to use standard user bearer tokens over public endpoints
   - *Pros:* Zero new gateway code.
   - *Cons:* User tokens expire frequently; lacks machine-to-machine security and IP whitelisting; vulnerable to rate limit abuse.
3. **Option C:** Manual daily database SQL exports sent via SFTP
   - *Pros:* No real-time API needed.
   - *Cons:* Completely unsuited for real-time equipment reservation and programmatic dispatch.

### 3. Chosen Solution & Technical Rationale
Option A: Scoped enterprise partner API gateway. Partners authenticate using mTLS paired with scoped API keys (dink_live_...) verified by edge Traefik middleware. Token bucket rate limits are enforced in Valkey per partner tier.

### 4. Rejected Alternatives
- Option B: Rejected because user session tokens are designed for human logins, not robust machine-to-machine integrations.
- Option C: Rejected because enterprise clients demand real-time automated equipment requisitioning.

### 5. Consequences & Invariants
- **Invariants Maintained:** Partner API keys must be cryptographically hashed (Argon2id/SHA-256) before storage in database.
- **Operational Trade-offs:** Requires an onboarding portal in Next.js Admin for generating and rotating partner API keys.

### 6. Verification
- **Test File:** `backend/internal/transport/http/middleware/partner_auth_test.go`
- **Assertion:** TestPartnerAuth_ValidatesScope asserts partner key with read:fleets cannot call write:bookings endpoint.

### 7. Revisit Trigger
- Revisit if enterprise clients standardize on open procurement protocols (e.g. OBI / cXML).

---

## D-07-07-002: Asynchronous Partner Webhook Event Delivery Engine with Cryptographic HMAC Signatures

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [API-WHT-05-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Partners integrating their ERPs require real-time notifications of booking state transitions (e.g. BOOKING_CONFIRMED, INSPECTION_COMPLETED). If partner servers are down, webhooks must retry reliably without overwhelming partner endpoints.

### 2. Considered Alternatives
1. **Option A:** Persistent partner webhook engine with HMAC-SHA256 signature headers (X-Dink-Signature), exponential backoff retries (up to 72 hours), and dead-letter monitoring
   - *Pros:* Guaranteed event delivery; partner can verify payload authenticity; resilient against partner server maintenance.
   - *Cons:* Requires webhook delivery logging and retry queues.
2. **Option B:** Synchronous HTTP POST calls directly to partner URLs during state transitions
   - *Pros:* Zero queue infrastructure.
   - *Cons:* Partner server slowness blocks primary booking transitions; network failure drops event permanently.
3. **Option C:** Requiring partners to poll GET /api/v1/partner/events every 10 seconds
   - *Pros:* Zero outbound webhook complexity.
   - *Cons:* Massive wasteful polling load on backend; high latency for event detection.

### 3. Chosen Solution & Technical Rationale
Option A: Asynchronous webhook delivery engine. Outbox events matching partner subscriptions are enqueued in partner_webhook_deliveries. Workers post payloads signed with partner HMAC secret and retry with backoff over 72 hours.

### 4. Rejected Alternatives
- Option B: Rejected because external partner server latency must never block internal transaction commits.
- Option C: Rejected because polling creates immense database read overhead and delays event synchronization.

### 5. Consequences & Invariants
- **Invariants Maintained:** Every partner webhook payload is signed with HMAC-SHA256 using the partner shared secret.
- **Operational Trade-offs:** Failed webhook endpoints must be automatically disabled after 100 consecutive failures to prevent queue bloat.

### 6. Verification
- **Test File:** `backend/internal/modules/integration/application/webhook_dispatch_test.go`
- **Assertion:** TestWebhook_GeneratesValidHMACSignature asserts X-Dink-Signature matches HMAC-SHA256 of payload.

### 7. Revisit Trigger
- Revisit if event streaming protocols (e.g. WebSockets or gRPC streams) replace webhooks for enterprise partners.

---

