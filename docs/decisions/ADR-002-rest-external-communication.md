# ADR-002: REST External Communication with RFC 9457 and Minor-Unit Money

- **Status:** Accepted
- **Date:** 2026-10-26
- **Decision Owner:** Dink Machinery Architecture Review Board
- **Affected Partitions:** 05-05, 05-09, 05-18, 14-appendices
- **Related Requirements:** API-01-BE, ERR-01-BE, BK-07-BE, ISO-4217

---

## 1. Context & Problem Statement

Dink Machinery coordinates multi-client interactions across Flutter mobile apps, Next.js admin consoles, and third-party bank webhook endpoints. Wire contract drift, unstructured error payloads, and floating-point monetary calculations introduce grave operational and legal risks, including ledger reconciliation mismatches and customer disputes over Birr/Cent rounding errors.

---

## 2. The 5-Point Qualification Evaluation

- [x] **Alternative Test:** Evaluated OpenAPI 3.1 REST, GraphQL, and gRPC with Protobuf.
- [x] **Surprise Test:** Enforces strict integer minor currency units (`amountMinor` in Santims/Cents) rather than standard JSON decimal numbers.
- [x] **Constraint Test:** Driven by financial ledger zero-loss invariants and RFC 9457 Problem Details compliance.
- [x] **Irreversibility Test:** Public and client-facing API contract schema alterations require coordinated multi-platform releases.
- [x] **Scope Test:** Directly governs all HTTP endpoints, mobile network clients, admin data fetchers, and error-handling pipelines.

---

## 3. Considered Options

### Option 1: OpenAPI 3.1 RESTful Contract Supremacy with RFC 9457 and Integer Minor-Unit Math
- **Summary:** Canonical YAML specification (`contracts/openapi.yaml`) acting as Single Source of Truth. All endpoints return standardized `SuccessEnvelope` or RFC 9457 `ErrorEnvelope`. All monetary values represent integer minor units (e.g. $100.50\text{ ETB} = 10050\text{ Santims}$).
- **Pros:**
  - Automated client and server code generation (`oapi-codegen`, `openapi-typescript`).
  - Total elimination of floating-point rounding errors across languages and architectures.
  - Universal, predictable error semantics conforming to international standards.
- **Cons:**
  - Strict contract-first workflow requires updating YAML before authoring application code.

### Option 2: GraphQL Schema with Custom Float Scalar Types
- **Summary:** Single GraphQL endpoint with nested queries and custom financial float scalars.
- **Pros:**
  - Flexible client queries with zero over-fetching.
- **Cons:**
  - Complex caching at the edge/gateway; lack of HTTP status code semantics; third-party bank webhooks cannot parse GraphQL.

### Option 3: gRPC with Protocol Buffers
- **Summary:** Binary Protobuf wire serialization with HTTP/2 transport.
- **Pros:**
  - High binary throughput and low wire serialization overhead.
- **Cons:**
  - Requires gRPC-Web proxy translation for admin browsers; cannot directly ingest webhooks from commercial banks.

---

## 4. Chosen Solution & Architectural Rationale

We adopt **Option 1: OpenAPI 3.1 RESTful Contract Supremacy with RFC 9457 and Integer Minor-Unit Math**.

1. **Contract Supremacy:** `contracts/openapi.yaml` is the canonical reference. Spectral linter rules enforce contract compliance in CI.
2. **Standardized Response Envelopes:**
   - 2xx Responses:
     ```json
     {
       "success": true,
       "data": { ... },
       "error": null,
       "meta": { "requestId": "req_...", "timestamp": "2026-10-26T12:00:00Z" }
     }
     ```
   - 4xx/5xx Responses (RFC 9457 Problem Details):
     ```json
     {
       "success": false,
       "data": null,
       "error": {
         "type": "https://api.dinkmachinery.com/errors/invalid-booking-state",
         "title": "Invalid Booking State Transition",
         "status": 409,
         "detail": "Cannot transition booking from CANCELLED to CONFIRMED.",
         "code": "BOOKING_INVALID_STATE_TRANSITION"
       },
       "meta": { "requestId": "req_...", "timestamp": "2026-10-26T12:00:00Z" }
     }
     ```
3. **Monetary Representation:**
   All financial fields use the composite `MoneyAmount` object or 64-bit integer minor units (`amountMinor` int64 + `currency` string).

---

## 5. Rejected Alternatives & Technical Justification

- **Option 2 (GraphQL):** Rejected due to the operational difficulty of edge rate limiting, lack of standard HTTP status caching, and inability of local banking partners to submit GraphQL payloads.
- **Option 3 (gRPC):** Rejected due to administrative browser compatibility constraints and lack of JSON webhook support from regional payment gateways.

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - Provable financial precision with zero rounding discrepancy between Go, Dart, and TypeScript.
  - Consistent error parsing in mobile and web interceptors.
  - Automated code generation reduces boilerplate code drift.
- **Negative / Technical Debt:**
  - Developers must edit `contracts/openapi.yaml` prior to backend or frontend implementation.
- **Operational / Observability Impacts:**
  - Request IDs (`X-Request-ID`) propagate across all envelopes, logs, and OpenTelemetry trace spans.

---

## 7. Verification Strategy & Invariant Tests

- **Spectral Linter Rule:** `.spectral.yaml` rule `no-floating-point-money` asserts that no schema defines `type: number` for monetary fields.
- **Envelope Compliance Test:** `backend/internal/transport/http/middleware/envelope_test.go` asserts that all handler responses wrap data in `SuccessEnvelope` or `ErrorEnvelope`.
- **Contract Code-Gen Check:** CI runs `oapi-codegen` and verifies git status is clean.

---

## 8. Revisit Conditions

- Introduction of micro-telematics streaming requiring high-frequency binary serialization (>10,000 pings/sec).
