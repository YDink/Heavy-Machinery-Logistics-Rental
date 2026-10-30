# Partition Decisions: 07-01 — Phase 3 In-App Mobile Money Gateway & Escrow Settlement Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `07-01`.

- **Governing System ADRs:** `ADR-002`, `ADR-010`, `ADR-015`
- **Subsystem Scope:** backend
- **Requirement Targets:** PAY-01-BE to PAY-12-BE

---

## D-07-01-001: Direct Mobile Money Gateway Orchestration (Telebirr, M-Pesa, CBE Birr) with Webhook Verification

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [PAY-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
In Phase 3, migrating from manual bank transfer slips to automated in-app digital checkout requires integrating regional mobile money APIs (Telebirr and CBE Birr in Ethiopia, M-Pesa Daraja in Kenya). Gateway webhooks can be replayed, delayed, or forged by attackers.

### 2. Considered Alternatives
1. **Option A:** Pluggable PaymentGateway interface with cryptographic webhook signature verification (HMAC-SHA256 / RSA), idempotency keys, and asynchronous status polling fallback
   - *Pros:* Automated instant checkout; protection against forged webhooks; automatic recovery if gateway webhook drops.
   - *Cons:* Requires maintaining bespoke gateway adapters and security certificates.
2. **Option B:** Relying exclusively on unverified incoming webhook parameters
   - *Pros:* Fast prototype.
   - *Cons:* Catastrophic vulnerability: attackers can forge fake payment confirmation webhooks and steal equipment rentals.
3. **Option C:** Third-party international payment aggregator (e.g. Stripe)
   - *Pros:* Single unified API.
   - *Cons:* Stripe does not operate domestically in Ethiopia and cannot process Telebirr or CBE Birr in local currency.

### 3. Chosen Solution & Technical Rationale
Option A: Pluggable PaymentGateway interface with cryptographic webhook validation. Webhook payloads are verified against provider public keys. In the event of a dropped webhook, a reconciliation worker polls the gateway status API after 3 minutes.

### 4. Rejected Alternatives
- Option B: Rejected because unauthenticated webhooks violate basic financial security invariants.
- Option C: Rejected because international aggregators do not support domestic Ethiopian banking rails.

### 5. Consequences & Invariants
- **Invariants Maintained:** No payment transaction transitions to SUCCESS without cryptographic signature verification or direct provider API polling confirmation.
- **Operational Trade-offs:** Requires secure storage of mobile money gateway API credentials in OpenBao KMS.

### 6. Verification
- **Test File:** `backend/internal/modules/payment/infrastructure/telebirr_adapter_test.go`
- **Assertion:** TestTelebirr_RejectsInvalidSignature asserts webhook with mismatched HMAC signature returns HTTP 401 Unauthorized.

### 7. Revisit Trigger
- Revisit if national unified QR payment standards (e.g. EthSwitch National Switch) standardize mobile money APIs.

---

## D-07-01-002: Idempotent Payment Intent Life-Cycle with Two-Phase Commit Reservation Finalization

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [PAY-05-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Under unstable mobile connections, users may double-tap payment buttons or network timeouts may leave payment state ambiguous between the mobile app, backend, and external telecom gateway.

### 2. Considered Alternatives
1. **Option A:** Payment intent state machine (CREATED -> PROCESSING -> SUCCEEDED / FAILED) with client-generated UUIDv4 Idempotency-Key and atomic reservation finalization
   - *Pros:* Eliminates double-billing; deterministic recovery from network dropouts; clear auditability.
   - *Cons:* Requires intent state transition modeling.
2. **Option B:** Direct charge execution on checkout button press without intent state tracking
   - *Pros:* Fewer database tables.
   - *Cons:* Frequent duplicate charges when users retry during slow cellular network responses.
3. **Option C:** Refund duplicate charges manually via customer support
   - *Pros:* Zero technical prevention.
   - *Cons:* Disastrous customer experience and high payment gateway refund dispute fees.

### 3. Chosen Solution & Technical Rationale
Option A: Payment intent lifecycle with strict idempotency keys. The client creates a payment intent, receives a payment token, and submits payment. Duplicate submissions with the same idempotency key return the existing intent status without re-charging.

### 4. Rejected Alternatives
- Option B: Rejected because duplicate charges for high-value equipment rentals destroy user trust.
- Option C: Rejected because manual refunds are slow and incur unnecessary transaction fees.

### 5. Consequences & Invariants
- **Invariants Maintained:** Submitting a payment request with an identical Idempotency-Key within 24 hours returns the cached original result.
- **Operational Trade-offs:** Idempotency keys must be indexed in PostgreSQL and cached in Valkey.

### 6. Verification
- **Test File:** `backend/internal/modules/payment/application/payment_intent_test.go`
- **Assertion:** TestPaymentIntent_IdempotentReplayReturnsSameResponse asserts duplicate requests execute exactly one financial charge.

### 7. Revisit Trigger
- Revisit if central bank mandates real-time ISO 20022 instant payment rail architecture.

---

