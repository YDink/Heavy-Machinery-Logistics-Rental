# API Contracts & Schema Governance — Dink Machinery

This directory contains the canonical API definitions and schema specifications for Dink Machinery. The project enforces a strict **Contract-First** discipline: no endpoint, request payload, or response structure may be implemented in backend or client applications without first being defined, linted, and reviewed in this directory.

---

## 1. Canonical Contract Inventory

- **`openapi.yaml`**: The primary OpenAPI 3.1.0 specification defining all external and internal HTTP REST endpoints across the modular monolith, mobile application, and administrative web portal.

---

## 2. Invariant Design Rules

Every operation declared within `openapi.yaml` must adhere to the platform-wide architectural invariants:

### A. Universal Response Envelope
Endpoints must not return raw domain entities, flat arrays, or unadorned status strings. Every response must resolve to one of two envelope structures:

1. **Success Envelope (`200`, `201`, `202`):**
   ```json
   {
     "success": true,
     "data": {},
     "error": null,
     "meta": {
       "version": "1.0",
       "requestId": "req_01HP8XYZ12345678",
       "timestamp": "2026-03-31T08:30:00Z",
       "pagination": null,
       "warnings": []
     }
   }
   ```
2. **Error Envelope (RFC 9457 Problem Details):**
   ```json
   {
     "success": false,
     "data": null,
     "error": {
       "type": "[https://api.dinkmachinery.com/errors/RESOURCE_LOCKED](https://api.dinkmachinery.com/errors/RESOURCE_LOCKED)",
       "title": "Resource Conflict",
       "status": 409,
       "detail": "The requested machinery fleet is locked under an active hold.",
       "code": "FLEET_HOLD_CONFLICT",
       "invalidParams": []
     },
     "meta": {
       "version": "1.0",
       "requestId": "req_01HP8XYZ12345678",
       "timestamp": "2026-03-31T08:30:00Z"
     }
   }
   ```

### B. Integer Minor-Unit Currency
- Floating-point representations (`float`, `double`, number with decimal places) for monetary amounts are strictly prohibited.
- All monetary fields must use integer types representing the smallest currency unit (e.g., Santim for ETB, Cents for KES/USD).
- Every monetary amount must be accompanied by an ISO 4217 three-letter currency code:
  ```yaml
  MoneyAmount:
    type: object
    required: [amountMinor, currency]
    properties:
      amountMinor:
        type: integer
        format: int64
        description: Amount in minor currency units (e.g., 10050 = 100.50).
        example: 10050
      currency:
        type: string
        pattern: '^[A-Z]{3}$'
        description: ISO 4217 three-letter currency code.
        example: "ETB"
  ```

### C. Standard Headers
- **`X-Country-Code` (Required on localized endpoints):** ISO 3166-1 alpha-2 country code (`ET`, `KE`) resolving multi-country operational rules via `country_configs`.
- **`Idempotency-Key` (Required on mutating financial & booking endpoints):** UUID v4 or ULID preventing duplicate execution on `POST` and state-transition actions.
- **`traceparent` (W3C Trace Context):** Propagated across gateway and internal services for distributed tracing.

### D. Cursor-Based Pagination
- Operational datasets (fleets, bookings, ledger entries, notifications) must paginate via opaque cursors.
- Query parameters must specify `limit` (integer, max 100) and `cursor` (string, base64 encoded).
- Response metadata must include `meta.pagination` containing `nextCursor`, `prevCursor`, and `hasMore`.

---

## 3. Contract Linting & Quality Gates

The contract is validated on every pre-commit hook and in CI via **Spectral CLI** using the custom ruleset defined in `.spectral.yaml`.

### Running Contract Linter Locally
```bash
# Lint openapi.yaml against strict project standards
npx @stoplight/spectral-cli lint contracts/openapi.yaml --ruleset .spectral.yaml
```

### Enforced Rules
- `oas3-schema`: Valid OpenAPI 3.1 structure.
- `operation-operationId`: Every operation must define a unique `operationId`.
- `operation-tags`: Every operation must belong to at least one logical module tag.
- `rfc9457-error-response`: Default or 4xx/5xx responses must match the RFC 9457 Problem Details schema.
- `no-floating-point-money`: Schemas containing "amount", "price", or "balance" must type-check as integer minor units.

---

## 4. Code Generation Targets (Phase 1+)

Upon transition to Phase 1 (Bootstrap), client SDKs and server scaffolding are generated directly from `contracts/openapi.yaml`:

| Target | Toolchain | Output Directory | Notes |
|---|---|---|---|
| **Go 1.27.3 Backend** | `oapi-codegen` | `backend/internal/transport/http/` | Generates strict interfaces, chi/echo routing wrappers, and DTO structs. |
| **Flutter 3.44 Mobile** | `openapi-generator-cli` | `mobile/packages/dink_api_client/` | Generates Dart models and type-safe Dio API client. |
| **Next.js 16 Admin** | `openapi-typescript` + `zod` | `admin/lib/api/generated/` | Generates TypeScript interfaces and runtime validation schemas. |

---

## 5. Contract Change Protocol

1. **Additive Compatibility:** Changes to `contracts/openapi.yaml` must be backward-compatible whenever possible. Adding optional request fields or new response properties is permitted.
2. **Breaking Changes:** Removing endpoints, deleting properties, renaming fields, or modifying validation bounds requires an approved **L1 ADR** (`docs/decisions/ADR-NNN-<slug>.md`) before merging to `main`.
3. **Deprecation Window:** Deprecated endpoints must be annotated with `deprecated: true` and emit RFC 8594 `Deprecation` and `Sunset` headers for a minimum of 90 days prior to removal.