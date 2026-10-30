# Partition Decisions: 05-18 — Error Handling, User Feedback & Problem Details Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-18`.

- **Governing System ADRs:** `ADR-002`, `ADR-008`
- **Subsystem Scope:** backend, mobile, admin
- **Requirement Targets:** ERR-01-BE to ERR-08-BE

---

## D-05-18-001: Universal RFC 9457 Problem Details Envelope Transformation in Transport Middleware

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [ERR-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
API errors originating from database constraints, domain validation rules, payment timeouts, and third-party gateways produce disparate error formats. Leaking raw database driver errors exposes internal table structures, while vague generic errors prevent mobile clients from presenting actionable guidance to users.

### 2. Considered Alternatives
1. **Option A:** Global HTTP middleware intercepting all error types and mapping domain errors to RFC 9457 application/problem+json Problem Details envelopes with sanitization
   - *Pros:* Unified error schema across all endpoints; zero sensitive internal leakage (passwords, SQL queries, stack traces); typed error codes for client branching.
   - *Cons:* Requires mapping every domain error to an HTTP status and RFC 9457 slug.
2. **Option B:** Raw error responses serialized directly from Go error strings (e.g. {"error": err.Error()})
   - *Pros:* Minimal backend mapping code.
   - *Cons:* Leaks database table names, constraint violations, and internal server IPs in client-facing JSON; violates ASVS L3.
3. **Option C:** HTTP status codes alone without response body details
   - *Pros:* Standard HTTP semantics.
   - *Cons:* Clients cannot differentiate between multiple 409 or 422 error causes; poor user feedback.

### 3. Chosen Solution & Technical Rationale
Option A: Universal RFC 9457 Problem Details envelope transformation in Chi transport middleware. All errors return Content-Type: application/problem+json wrapping status, code, title, detail, and invalidParams with requestId metadata. Internal system errors return opaque SYS_INTERNAL_SERVER_ERROR codes while logging stack traces securely.

### 4. Rejected Alternatives
- Option B: Rejected because raw error string leakage violates OWASP ASVS Level 3 security and exposes system vulnerabilities.
- Option C: Rejected because mobile apps require granular machine-readable error codes to guide user recovery.

### 5. Consequences & Invariants
- **Invariants Maintained:** Zero internal database errors, SQL syntax fragments, or stack traces may ever be transmitted in HTTP response bodies.
- **Operational Trade-offs:** All new domain errors must be registered in ERROR_CATALOG.md with assigned status codes and URI slugs.

### 6. Verification
- **Test File:** `backend/internal/transport/http/middleware/error_handler_test.go`
- **Assertion:** TestErrorHandler_SanitizesSQLError asserts a PostgreSQL driver error is masked to SYS_DATABASE_ERROR with UUID correlation ID.

### 7. Revisit Trigger
- Revisit if international standards supersede RFC 9457 for REST API problem representation.

---

## D-05-18-002: Structured Client Error Presentation with Localized Actionable Recovery Hints

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [ERR-05-MOB]
- **Subsystem:** [mobile]

### 1. Problem Statement & Context
When API requests fail (e.g. payment hold expired, insufficient wallet balance, duplicate slip image), generic error banners ("An error occurred") confuse field contractors, leading to redundant retry spam or abandoned bookings.

### 2. Considered Alternatives
1. **Option A:** Dio network interceptor translating RFC 9457 error codes into localized vernacular messages with specific call-to-action buttons (e.g. "Top Up Wallet", "Retake Photo")
   - *Pros:* Clear actionable user recovery; localized text in active language (Amharic, English); prevents panic retries.
   - *Cons:* Requires mapping all error codes to translation keys in mobile app.
2. **Option B:** Displaying the backend detail string directly in snackbars
   - *Pros:* Zero client mapping effort.
   - *Cons:* Backend detail strings are English-only and may be too technical for field equipment operators.
3. **Option C:** Modal dialog popups on every error forcing user acknowledgment
   - *Pros:* Ensures user notices the error.
   - *Cons:* Highly disruptive to mobile workflow navigation; annoys contractors.

### 3. Chosen Solution & Technical Rationale
Option A: Structured client error presentation via Dio interceptors. Errors map RFC 9457 code to localized strings with recovery action widgets (e.g. COMMISSION_WALLET_INSUFFICIENT_FUNDS renders a Top Up Wallet shortcut pill). Non-fatal warnings use toast sonner alerts.

### 4. Rejected Alternatives
- Option B: Rejected because English technical details alienate non-English-speaking equipment operators.
- Option C: Rejected due to severe UI fatigue caused by intrusive modal interruptions.

### 5. Consequences & Invariants
- **Invariants Maintained:** Every client-facing error message must include a localized human-readable explanation and an actionable recovery step.
- **Operational Trade-offs:** Translation keys for all error codes must be synchronized across OTA translation bundles.

### 6. Verification
- **Test File:** `mobile/test/core/errors/error_presenter_test.dart`
- **Assertion:** TestErrorPresenter_MapsWalletErrorToActionableUI asserts COMMISSION_WALLET_INSUFFICIENT_FUNDS renders localized top-up action.

### 7. Revisit Trigger
- Revisit if voice-based audio error announcements are introduced for low-literacy users.

---

