# Standard Error Catalog & RFC 9457 Problem Details Registry — Dink Machinery

This document is the canonical registry for all application error codes, HTTP status mappings, RFC 9457 URI type identifiers, retryability semantics, and masking policies across the Dink Machinery modular monolith, mobile application, and admin console.

---

## 1. Architectural Standards & Wire Format

All client-facing errors (HTTP 4xx and 5xx) must adhere strictly to **RFC 9457 (Problem Details for HTTP APIs)** encapsulated within the platform's standardized JSON response envelope. Direct plain-text, unformatted HTML, or raw database error dumps are prohibited.

### A. Error Response Envelope Structure
```json
{
  "success": false,
  "data": null,
  "error": {
    "type": "[https://api.dinkmachinery.com/errors/booking-invalid-state-transition](https://api.dinkmachinery.com/errors/booking-invalid-state-transition)",
    "title": "Invalid Booking State Transition",
    "status": 409,
    "detail": "Cannot transition booking from CANCELLED to CONFIRMED.",
    "code": "BOOKING_INVALID_STATE_TRANSITION",
    "invalidParams": [
      {
        "name": "status",
        "reason": "Current state is terminal."
      }
    ]
  },
  "meta": {
    "version": "1.0",
    "requestId": "req_01HP8XYZ12345678",
    "timestamp": "2026-10-26T07:12:00Z"
  }
}
```

### B. Core Error Invariants
1. **Zero PII Leakage (`ASVS L3`):** Error envelopes must never include unmasked phone numbers, government IDs, National IDs, bank account numbers, credit card tokens, or raw system stack traces.
2. **Deterministic Code Format:** Error codes must be formatted in uppercase snake case prefixed with their domain namespace: `<DOMAIN>_<SPECIFIC_REASON>` (e.g., `AUTH_INVALID_CREDENTIALS`, `BOOKING_HOLD_CONFLICT`).
3. **Internal Error Obfuscation:** Uncaught 500 errors must emit `INTERNAL_SERVER_ERROR` with a generic, safe client message. Detailed debugging metadata, stack traces, and SQL states must only be recorded in backend structured logs tagged with the correlating `requestId`.
4. **Idempotency Replay Invariance:** Replaying an identical `Idempotency-Key` on a failed mutation must return the exact same cached RFC 9457 error response originally generated, rather than generating duplicate side effects or alternative errors.

---

## 2. Global & Transport Error Codes (`SYS_*`, `REQ_*`)

| Error Code | HTTP Status | RFC 9457 Type URI Slug | Retryable | Description / Trigger Conditions |
|---|---|---|---|---|
| `REQ_VALIDATION_FAILED` | 400 | `/errors/request-validation-failed` | No | Request body or query parameters failed JSON schema or boundary checks. Populates `invalidParams`. |
| `REQ_MALFORMED_JSON` | 400 | `/errors/malformed-json` | No | Body could not be parsed as valid JSON syntax. |
| `REQ_MISSING_IDEMPOTENCY_KEY` | 400 | `/errors/missing-idempotency-key` | No | Mutating financial or booking operation omitted mandatory `Idempotency-Key` header. |
| `REQ_INVALID_COUNTRY_CODE` | 400 | `/errors/invalid-country-code` | No | Supplied `X-Country-Code` does not match an enabled ISO 3166-1 alpha-2 jurisdiction. |
| `SYS_RATE_LIMIT_EXCEEDED` | 429 | `/errors/rate-limit-exceeded` | Yes (Exp) | Ingress rate-limiting threshold exceeded. Includes `Retry-After` header. |
| `SYS_SERVICE_UNAVAILABLE` | 503 | `/errors/service-unavailable` | Yes (Exp) | Upstream dependency or persistence pool saturated. |
| `SYS_INTERNAL_ERROR` | 500 | `/errors/internal-server-error` | No | Unhandled runtime panic or unexpected persistence failure. |

---

## 3. Auth & Identity Module (`AUTH_*`)

| Error Code | HTTP Status | RFC 9457 Type URI Slug | Retryable | Description / Trigger Conditions |
|---|---|---|---|---|
| `AUTH_INVALID_TOKEN` | 401 | `/errors/invalid-token` | No | Bearer token is expired, malformed, or cryptographic signature failed verification. |
| `AUTH_TOKEN_REVOKED` | 401 | `/errors/token-revoked` | No | Refresh or access token has been invalidated via remote logout or password reset. |
| `AUTH_OTP_MAX_ATTEMPTS` | 429 | `/errors/otp-max-attempts` | Yes (15m) | User exceeded allowed code verification attempts (max 5 per challenge session). |
| `AUTH_OTP_EXPIRED` | 400 | `/errors/otp-expired` | No | OTP challenge validity window (typically 120 seconds) has elapsed. |
| `AUTH_OTP_INVALID_CODE` | 400 | `/errors/otp-invalid-code` | Yes | Supplied 6-digit numeric verification code does not match pending challenge. |
| `AUTH_INSUFFICIENT_PERMISSIONS` | 403 | `/errors/insufficient-permissions` | No | Authenticated principal role lacks RBAC privileges to access requested endpoint. |
| `AUTH_ACCOUNT_SUSPENDED` | 403 | `/errors/account-suspended` | No | Account is flagged as suspended or deactivated due to regulatory/fraud enforcement. |
| `AUTH_USER_NOT_FOUND` | 404 | `/errors/user-not-found` | No | Targeted user aggregate does not exist or has been soft-deleted. |

---

## 4. Machinery & Fleet Module (`FLEET_*`, `CAT_*`)

| Error Code | HTTP Status | RFC 9457 Type URI Slug | Retryable | Description / Trigger Conditions |
|---|---|---|---|---|
| `CAT_NOT_FOUND` | 404 | `/errors/category-not-found` | No | Machinery taxonomy category identifier does not exist. |
| `FLEET_NOT_FOUND` | 404 | `/errors/fleet-not-found` | No | Machinery fleet asset does not exist or has `deleted_at IS NOT NULL`. |
| `FLEET_REGISTRATION_EXISTS` | 409 | `/errors/fleet-registration-exists` | No | Machinery registration license/chassis number is already registered under active fleet. |
| `FLEET_STATUS_CONFLICT` | 409 | `/errors/fleet-status-conflict` | No | Fleet asset cannot accept bookings because its status is `UNDER_MAINTENANCE` or `DECOMMISSIONED`. |
| `FLEET_ATTACHMENT_UNSUPPORTED` | 422 | `/errors/fleet-attachment-unsupported`| No | Requested machinery attachment is mechanically incompatible with target base unit specs. |
| `FLEET_DAILY_RATE_INVALID` | 422 | `/errors/fleet-daily-rate-invalid` | No | Daily rate minor units fall below minimum threshold or exceed allowed jurisdictional limits. |

---

## 5. Booking & State Machine Engine (`BOOKING_*`)

| Error Code | HTTP Status | RFC 9457 Type URI Slug | Retryable | Description / Trigger Conditions |
|---|---|---|---|---|
| `BOOKING_NOT_FOUND` | 404 | `/errors/booking-not-found` | No | Specified booking identifier does not resolve to an active entity. |
| `BOOKING_DATES_UNAVAILABLE` | 409 | `/errors/booking-dates-unavailable` | No | Target machinery asset is already reserved or held for overlapping calendar dates. |
| `BOOKING_HOLD_EXPIRED` | 409 | `/errors/booking-hold-expired` | No | Reservation payment hold timer expired; asset returned to public search pool. |
| `BOOKING_INVALID_STATE_TRANSITION`| 409 | `/errors/booking-invalid-state` | No | Attempted state transition violates 20-state lifecycle rules (e.g., `CONFIRMED` to `DRAFT`). |
| `BOOKING_CONCURRENCY_CONFLICT` | 409 | `/errors/booking-concurrency-conflict` | Yes | Optimistic locking violation (`version` column mismatch). Client should reload and retry. |
| `BOOKING_INSPECTION_REQUIRED` | 422 | `/errors/booking-inspection-required` | No | Cannot advance booking to `ACTIVE_OPERATION` or `COMPLETED` without signed site inspection. |
| `BOOKING_CANCELLATION_DISALLOWED`| 422 | `/errors/booking-cancellation-disallowed`| No | Booking has advanced past the permissible self-service cancellation window. |

---

## 6. Payment, Ledger & Owner Wallet (`PAYMENT_*`, `LEDGER_*`, `WALLET_*`)

| Error Code | HTTP Status | RFC 9457 Type URI Slug | Retryable | Description / Trigger Conditions |
|---|---|---|---|---|
| `PAYMENT_GATEWAY_TIMEOUT` | 504 | `/errors/payment-gateway-timeout` | Yes (Exp) | Upstream provider (Telebirr/CBEBirr/MPESA) failed to respond within timeout budget. |
| `PAYMENT_GATEWAY_REJECTED` | 402 | `/errors/payment-gateway-rejected` | No | Provider explicitly rejected the charge (e.g., insufficient customer wallet balance). |
| `PAYMENT_TRANSACTION_EXISTS` | 409 | `/errors/payment-transaction-exists`| No | Transaction identifier has already been ingested and settled in ledger. |
| `WALLET_INSUFFICIENT_FUNDS` | 422 | `/errors/wallet-insufficient-funds` | No | Withdrawal request exceeds verified available balance in `OwnerWallet`. |
| `WALLET_LOCKED_PENDING_AUDIT` | 423 | `/errors/wallet-locked-pending-audit`| No | Owner wallet is under administrative compliance lock; disbursements suspended. |
| `LEDGER_DOUBLE_ENTRY_IMBALANCE` | 500 | `/errors/ledger-imbalance` | No | CRITICAL: Transaction debit and credit legs do not balance to zero minor units. Tx rolled back. |
| `LEDGER_IMMUTABLE_MUTATION` | 500 | `/errors/ledger-immutable-mutation` | No | Disallowed SQL execution attempted to update or delete an existing `financial_ledger_entries` record. |

---

## 7. Inspection & Geospatial Modules (`INSPECT_*`, `GEO_*`)

| Error Code | HTTP Status | RFC 9457 Type URI Slug | Retryable | Description / Trigger Conditions |
|---|---|---|---|---|
| `INSPECT_METER_REGRESSION` | 422 | `/errors/inspect-meter-regression` | No | Operating hours meter on check-out is lower than initial check-in reading. |
| `INSPECT_MEDIA_SCAN_FAILED` | 422 | `/errors/inspect-media-scan-failed` | No | Attached inspection photo/document failed ClamAV antivirus inspection. |
| `GEO_INVALID_COORDINATES` | 400 | `/errors/geo-invalid-coordinates` | No | Latitude or longitude falls outside valid geographic boundaries (-90..90, -180..180). |
| `GEO_RADIUS_OUT_OF_BOUNDS` | 400 | `/errors/geo-radius-out-of-bounds` | No | Radius parameter exceeds allowable maximum search envelope (1 km to 500 km). |

---

## 8. Multi-Country & Administration (`COUNTRY_*`, `ADMIN_*`)

| Error Code | HTTP Status | RFC 9457 Type URI Slug | Retryable | Description / Trigger Conditions |
|---|---|---|---|---|
| `COUNTRY_CONFIG_NOT_FOUND` | 404 | `/errors/country-config-not-found` | No | Specified ISO country code is not registered in `country_configs`. |
| `COUNTRY_TAX_RATE_INVALID` | 422 | `/errors/country-tax-rate-invalid` | No | Basis points (`bps`) for VAT or Withholding Tax outside permissible bounds (0 to 10000 bps). |
| `ADMIN_UNAUTHORIZED_OVERRIDE` | 403 | `/errors/admin-unauthorized-override`| No | Super-admin clearance required to bypass manual hold or unlock escrow. |

---

## 9. Client Error Handling & Localization Guidelines

### A. Mobile Client (Flutter BLoC)
- The Dio networking interceptor intercepts non-2xx responses and transforms them into typed `ProblemDetails` models.
- Domain BLoCs switch on `error.code` rather than parsing error message strings.
- User-facing error notifications display localized UI strings looked up via `error.code`. If an unknown code is returned, fallback to a friendly default mapped to `error.status`.

### B. Admin Console (Next.js 16)
- Server Actions intercept RFC 9457 responses and feed validation arrays (`invalidParams`) directly into `react-hook-form` field errors via `@hookform/resolvers/zod`.
- Global errors render as non-blocking `sonner` toast notifications with the `requestId` displayed for operations debugging.