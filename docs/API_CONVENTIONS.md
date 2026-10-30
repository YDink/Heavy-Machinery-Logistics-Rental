# API Conventions & REST Design Standards — Dink Machinery

**System:** Heavy Machinery & Logistics Rental Marketplace  
**Specification:** REST API Style Guide & Communication Protocol  
**Compliance:** RFC 9457 (Problem Details), RFC 3339 (Date/Time), ISO 4217 (Currencies), ISO 3166-1 alpha-2 (Countries), OWASP ASVS Level 3

---

## 1. Overview & Core Architectural Invariants

This document establishes the mandatory API conventions for all internal and external HTTP services across Dink Machinery. Every endpoint implemented in the Go 1.27.3 modular monolith, consumed by the Flutter 3.47.9 mobile client, or called by the Next.js 16 administrative console must adhere strictly to these conventions.

### Summary of Platform Invariants
1. **Universal Response Envelope:** Raw unadorned JSON models, flat arrays, or plain text responses are forbidden. Every response must resolve to either a Success or RFC 9457 Error envelope.
2. **Integer Minor-Unit Money:** Floating-point numbers (`float`, `double`) are strictly prohibited in financial schemas. Amounts are formatted as 64-bit signed integers representing the lowest non-divisible minor currency unit coupled with an ISO 4217 currency code.
3. **Cursor-Based Keyset Pagination:** Operational datasets (fleets, bookings, ledger logs) must paginate using opaque, URL-safe cursors. Offset/limit queries are prohibited on transaction tables.
4. **RPC Action Suffixes:** State machine transitions and non-CRUD operational triggers use standard HTTP `POST` appending an action suffix (e.g., `POST /api/v1/bookings/{id}:confirm-payment`).
5. **Strict Casing Isolation:** Casing styles are isolated by architectural tier:
   - Resource URIs: `kebab-case`
   - JSON keys: `lowerCamelCase`
   - Database tables & columns: `snake_case`
   - Enums & Constants: `UPPER_SNAKE_CASE`
6. **Zero PII Exposure (ASVS L3):** User phone numbers, National IDs, Kebele registration numbers, TINs, and bank accounts must never appear in URI segments, query parameters, or client error envelopes.

---

## 2. URI Design & Routing Structure

### 2.1. Base URI & API Versioning
All public and internal endpoints must be prefixed with the active major version identifier:
```text
[https://api.dinkmachinery.com/api/v1](https://api.dinkmachinery.com/api/v1)
```
Major version bumps (e.g., `/api/v2`) require an approved **L1 ADR** and must support side-by-side execution with the prior version for a minimum 90-day deprecation window with RFC 8594 `Deprecation` and `Sunset` headers.

### 2.2. Resource Path Semantics
- **Plural Nouns:** Collection endpoints must use plural nouns (e.g., `/fleets`, `/bookings`, `/categories`).
- **Kebab-Case Paths:** Compound resource names must use lowercase kebab-case (e.g., `/fleet-equipments`, `/country-configs`).
- **Hierarchy Cap (Max 2 Levels):** Sub-resource hierarchies must not exceed two nesting levels:
  - *Compliant:* `/api/v1/bookings/{id}/inspections`
  - *Non-Compliant:* `/api/v1/fleets/{id}/bookings/{bookingId}/inspections/{inspectionId}/media`
  Operations exceeding two levels of relationship must resolve to root-level collections with query filters (e.g., `/api/v1/inspections?bookingId={id}`).

### 2.3. RPC-Style State Actions
When an operation represents a business workflow, state transition, or non-CRUD event rather than a property mutation, format the endpoint as:
```text
POST /api/v1/<resource>/{id}:<action>
```

#### Canonical Platform State Transitions
| Operation | HTTP Verb & URI Path | Success Status | Description |
|---|---|---|---|
| Verify Phone OTP | `POST /api/v1/auth/otp/verify` | `200 OK` | Exchanges challenge code for JWT pair. |
| Confirm Payment | `POST /api/v1/bookings/{id}:confirm-payment` | `200 OK` | Advances booking to `CONFIRMED` and locks escrow. |
| Cancel Booking | `POST /api/v1/bookings/{id}:cancel` | `200 OK` | Cancels hold and triggers refund assessment. |
| Submit Inspection | `POST /api/v1/bookings/{id}/inspections` | `201 Created` | Logs equipment check-in or check-out condition. |
| Initiate Dispute | `POST /api/v1/bookings/{id}:open-dispute` | `200 OK` | Freezes deposit and escalates to admin review. |

---

## 3. HTTP Verbs & Status Code Specifications

### 3.1. HTTP Method Semantics
| Method | Idempotent | Safe | Request Body | Purpose & Constraints |
|---|---|---|---|---|
| `GET` | Yes | Yes | Prohibited | Reads collections or individual resource records. Must never mutate state. |
| `POST` | No (requires header) | No | Required | Creates resources or executes RPC state machine transitions. |
| `PUT` | Yes | No | Required | Replaces an entire mutable resource representation or idempotently sets config. |
| `PATCH` | Yes | No | Required | Applies partial field updates to an existing resource. |
| `DELETE` | Yes | No | Optional | Soft-deletes a resource (`deleted_at = NOW()`). Hard database deletion is prohibited. |

### 3.2. HTTP Status Code Registry
Endpoints must return the most specific standard HTTP status code:

```
200 OK                  # Read success, update success, or RPC state transition complete
201 Created             # New resource created (returns Location header when applicable)
202 Accepted            # Asynchronous background job accepted (e.g., bulk report export)
204 No Content          # Successful deletion or state change returning no body

400 Bad Request         # Request failed schema, type, or edge validation (invalidParams populated)
401 Unauthorized        # Bearer token is missing, expired, or failed signature check
403 Forbidden           # Authenticated principal lacks necessary RBAC permissions
404 Not Found           # Resource does not exist or has deleted_at IS NOT NULL
409 Conflict            # Invariant violation, OCC version mismatch, or state machine failure
422 Unprocessable Entity# Domain validation failed (e.g., negative fuel reading, bad date range)
429 Too Many Requests   # Rate limit quota exceeded (Retry-After header mandatory)

500 Internal Error      # Unhandled server fault (returns anonymized error envelope)
503 Unavailable         # Database pool saturated or maintenance in progress
504 Gateway Timeout     # Upstream payment rail or external service failed within budget
```

---

## 4. Standard Request & Response Headers

### 4.1. Mandatory Inbound Request Headers
| Header Name | Required | Data Type | Example | Purpose & Architectural Invariant |
|---|---|---|---|---|
| `Authorization` | Conditional | String | `Bearer eyJhbGci...` | Required on all protected endpoints. Encrypted JWT session token. |
| `X-Country-Code` | Optional (defaults) | ISO 3166-1 alpha-2 | `ET` | Multi-country routing context (`LOC-01-BE`). Resolves tax rates, currency, timeouts. |
| `Idempotency-Key` | Mandatory on mutations | UUID v4 / ULID | `9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d` | Required on mutating financial and booking requests (`POST`, `PUT`, `PATCH`). |
| `traceparent` | Optional (Gateway sets) | W3C Trace Context | `00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01` | Distributed telemetry trace propagation across services and daemons. |
| `Content-Type` | Mandatory with body | String | `application/json` | Enforces JSON encoding for all request payloads. |

### 4.2. Mandatory Outbound Response Headers
| Header Name | Data Type | Example | Purpose |
|---|---|---|---|
| `Content-Type` | String | `application/json` | Success responses. |
| `Content-Type` (Errors) | String | `application/problem+json` | RFC 9457 Problem Details error responses. |
| `X-Request-Id` | String | `req_01HP8XYZ12345678` | Unique request tracking token generated at ingress; matches `meta.requestId`. |
| `Retry-After` | Integer | `30` | Mandatory on HTTP 429 and 503 responses indicating cooldown duration in seconds. |
| `Deprecation` | HTTP Date / Boolean | `true` | Emitted when an endpoint is marked for future retirement. |
| `Sunset` | HTTP Date | `Wed, 11 Nov 2026 00:00:00 GMT` | Definitive cutoff date for deprecated APIs. |

---

## 5. Universal Response Envelopes

Every JSON response must adhere to one of the two canonical response envelope specifications. Naked entities or disparate root-level structures are prohibited.

### 5.1. Success Envelope (HTTP 200, 201, 202)
```json
{
  "success": true,
  "data": {
    "id": "7f000001-8b22-1d54-818b-22b0f4910001",
    "fleetId": "7f000001-8b22-1d54-818b-22b0f4910002",
    "renterId": "7f000001-8b22-1d54-818b-22b0f4910003",
    "status": "CONFIRMED",
    "startDate": "2026-11-01",
    "endDate": "2026-11-07",
    "totalCost": {
      "amountMinor": 3500000,
      "currency": "ETB"
    },
    "depositHold": {
      "amountMinor": 500000,
      "currency": "ETB"
    },
    "version": 2,
    "createdAt": "2026-10-26T07:12:00Z"
  },
  "error": null,
  "meta": {
    "version": "1.0",
    "requestId": "req_01HP8XYZ12345678",
    "timestamp": "2026-10-26T07:12:00Z",
    "pagination": null,
    "warnings": []
  }
}
```

### 5.2. Error Envelope (RFC 9457 Problem Details)
All client-facing 4xx and 5xx errors must return `Content-Type: application/problem+json` wrapped in the standard error envelope.

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
        "reason": "Current booking status is terminal."
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

#### RFC 9457 Attributes
- **`type` (URI):** Absolute URI resolving to documentation for the problem type.
- **`title` (String):** Short, human-readable summary of the problem type (does not change across instances).
- **`status` (Integer):** The HTTP status code produced by the origin server.
- **`detail` (String):** Human-readable explanation specific to this instance of the error.
- **`code` (String):** Unique machine-readable domain code from `ERROR_CATALOG.md` (e.g., `BOOKING_HOLD_EXPIRED`).
- **`invalidParams` (Array):** List of field-level schema or boundary validation failures. Empty when error is contextual.

---

## 6. Monetary Amounts & Currency Representation

Monetary values must never be modeled, calculated, or transmitted as floating-point numbers.

### 6.1. Money Amount Object Schema
Every monetary attribute (daily rates, totals, platform fees, taxes, escrow amounts) must use the `MoneyAmount` composite object:

```json
{
  "amountMinor": 150000,
  "currency": "ETB"
}
```

### 6.2. Minor-Unit Exponent Rules
Calculations and storage in database columns must use 64-bit integers (`int64` / `BIGINT`):
$$\text{amountMinor} = \text{majorUnits} \times 10^{\text{exponent}}$$

| Currency | ISO 4217 Code | Minor Unit Name | Exponent | Example: 1,500.00 Units |
|---|---|---|---|---|
| **Ethiopian Birr** | `ETB` | Santim | 2 | `150000` |
| **Kenyan Shilling** | `KES` | Cent | 2 | `150000` |
| **US Dollar** | `USD` | Cent | 2 | `150000` |

---

## 7. Cursor-Based Keyset Pagination

To guarantee deterministic reads over high-concurrency transactional datasets, endpoints returning collections must use keyset pagination with opaque cursors. Offset/limit pagination (`?offset=100&limit=20`) is prohibited on operational tables.

### 7.1. Query Parameters
| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `limit` | Integer | No | `20` | Number of items to return. Maximum allowable limit is `100`. |
| `cursor` | String (Base64) | No | `null` | Opaque cursor token encoding the sorting keys of the last seen record. |

### 7.2. Paginated Success Response
Paginated list queries return the items in `data` (as an array) and include keyset markers in `meta.pagination`:

```json
{
  "success": true,
  "data": [
    {
      "id": "7f000001-8b22-1d54-818b-22b0f4910010",
      "model": "Caterpillar 320D3",
      "dailyRate": {
        "amountMinor": 2500000,
        "currency": "ETB"
      }
    }
  ],
  "error": null,
  "meta": {
    "version": "1.0",
    "requestId": "req_01HP8XYZ12345678",
    "timestamp": "2026-10-26T07:12:00Z",
    "pagination": {
      "nextCursor": "ZXlKa1pYWjBaV05vSWpvaU1qQXlOaTB4TVMwd05WUXdPVFE2TVRBNk1EQXpJanAwY0hWbElqb2lkbVZ5YzJsdmJpSTZNVGswTWpZeU1EQXNJbWxrSWpvaU4yWXdNR0F3TVM0NFlqSXlMVEZrTlRRdE9ERTROeTAyTW1Jd1pqUTVNVEEwTmpBaUxDSnBZM1FnSWpvaU1qQXlOaTB4TVMwd05WUXdPVFE2TVRBNk1EQXpJanA5",
      "prevCursor": null,
      "hasMore": true
    },
    "warnings": []
  }
}
```

### 7.3. Cursor Encoding Standard
Cursors must be opaque to the client. The backend constructs cursors by serializing a JSON tuple of the order-by columns and base64-encoding the result:
```json
// Internal Unencoded Cursor Structure
{
  "createdAt": "2026-11-05T09:14:03Z",
  "id": "7f000001-8b22-1d54-818b-22b0f4910010"
}
```
SQL generation applies keyset filters matching the active sort direction:
```sql
WHERE (created_at, id) < ($1, $2)
ORDER BY created_at DESC, id DESC
LIMIT $3;
```

---

## 8. Idempotency & Optimistic Concurrency Protocol

### 8.1. Idempotency Key Lifecycle (`POST`, `PUT`, `PATCH`)
Mutating requests affecting bookings, payment ledger entries, or asset status must include the `Idempotency-Key` header.

```
Client                        Gateway & Middleware                    Backend & DB
  |                                     |                                   |
  |--- POST /bookings (Idem-Key: K) --->|                                   |
  |                                     |--- Check Valkey for Key K ------->|
  |                                     |    [State: IN_PROGRESS]           |
  |                                     |<-- Key exists & in-flight --------|
  |<-- 409 Conflict (Pending) ----------|                                   |
  |                                     |                                   |
  |                                     |--- [No Key Found] --------------->|
  |                                     |    Lock Key K in Valkey           |
  |                                     |    Execute Mutation inside DB Tx  |
  |                                     |    Write Record + Idempotency Log |
  |                                     |    Cache Result Payload in Valkey |
  |                                     |<-- Tx Committed ------------------|
  |<-- 201 Created (Success Envelope) --|                                   |
  |                                     |                                   |
  |--- Replay: POST (Idem-Key: K) ----->|                                   |
  |                                     |--- Check Valkey for Key K ------->|
  |                                     |<-- Return Cached Response --------|
  |<-- 201 Created (Exact Replay) ------|                                   |
```

1. **Key Generation:** Client generates a unique UUID v4 or ULID prior to submitting the request.
2. **In-Flight Lock:** The server acquires a short-lived distributed lock in Valkey (`SET lock:idem:<key> true NX EX 30`). If lock acquisition fails, return `409 Conflict` (`SYS_CONCURRENT_REQUEST_IN_FLIGHT`).
3. **Transactional Storage:** The original request parameters and serialized response envelope are stored in the database within the primary transaction (`idempotency_keys` table).
4. **Replay Invariance:** If a subsequent request arrives with an identical key within the 24-hour retention window, the server returns the cached response with the original HTTP status code without re-executing business logic.

### 8.2. Optimistic Concurrency Control (OCC)
To prevent lost updates in multi-tenant workflows (e.g., simultaneous booking confirmations or owner edits), stateful entities maintain an integer `version` field.

- Mutating payloads must include the expected entity `version`:
  ```json
  {
    "version": 3,
    "status": "CHECKED_IN"
  }
  ```
- Persistence updates assert equality:
  ```sql
  UPDATE bookings
  SET status = $1, version = version + 1, updated_at = NOW()
  WHERE id = $2 AND version = $3;
  ```
- If 0 rows are updated, the service aborts and returns HTTP 409 (`BOOKING_CONCURRENCY_CONFLICT`).

---

## 9. Geospatial Query Conventions

Geospatial discovery endpoints leverage the PostGIS spatial engine (`PostgreSQL 18.6 + PostGIS 3.6.4`).

### 9.1. Coordinates & Projections
- Spatial coordinates must use the **WGS 84** coordinate reference system (`SRID 4326`).
- Floating-point precision: Longitude and latitude parameters must support up to 6 decimal places ($\approx 0.1\text{ meter}$ precision).
- Parameter names:
  - Latitude: `latitude` (range: `-90.0` to `90.0`)
  - Longitude: `longitude` (range: `-180.0` to `180.0`)

### 9.2. Radius Search Endpoint Format
```text
GET /api/v1/geo/fleets/radius-search?latitude=9.0300&longitude=38.7400&radiusKm=50&limit=20
```

#### Query Parameters
| Parameter | Type | Required | Bounds | Description |
|---|---|---|---|---|
| `latitude` | Double | Yes | `-90.0` to `90.0` | Center latitude. |
| `longitude` | Double | Yes | `-180.0` to `180.0` | Center longitude. |
| `radiusKm` | Integer | Yes | `1` to `500` | Search radius in kilometers. |
| `limit` | Integer | No | `1` to `100` (Default: 20) | Result page size limit. |
| `cursor` | String | No | Opaque | Keyset cursor for paginated results. |

---

## 10. Data Types, Formats & Casing Reference

| Data Category | Representation | Example | Validation & Boundary Rule |
|---|---|---|---|
| **Resource Identifier** | UUID v4 (string) | `"7f000001-8b22-1d54-818b-22b0f4910001"` | RFC 4122 compliant; lowercase hexadecimal with hyphens. |
| **Date & Time** | ISO 8601 / RFC 3339 | `"2026-10-26T07:12:00Z"` | Always UTC (`Z`). Trailing time-zone offsets are normalized to UTC upon ingest. |
| **Calendar Date** | ISO 8601 Date | `"2026-11-01"` | `YYYY-MM-DD` string representation. |
| **Country Code** | ISO 3166-1 alpha-2 | `"ET"` | 2-letter uppercase string. Must match registered `country_configs`. |
| **Currency Code** | ISO 4217 | `"ETB"` | 3-letter uppercase alphabetic string. |
| **Telephone Number** | E.164 International | `"+251911223344"` | Leading `+` followed by country dialing code and subscriber number (no spaces). |
| **Tax / Basis Points** | Integer Basis Points | `1500` | Integer minor units representing hundredths of a percent ($1500\text{ bps} = 15.00\%$). |
| **Coordinates** | Object (Lat / Long) | `{"latitude": 9.03, "longitude": 38.74}` | WGS 84 (`EPSG:4326`). |
| **Media Attachments** | Secure URI (HTTPS) | `"https://media.dinkmachinery.com/..."` | Pre-signed S3/SeaweedFS object URL passing ClamAV antivirus clearance. |