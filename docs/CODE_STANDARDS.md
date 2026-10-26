# Code Standards & Implementation Invariants — Dink Machinery

**System:** Heavy Machinery & Logistics Rental Marketplace  
**Scope:** Cross-stack engineering guidelines for Go 1.27.3, Flutter 3.44, Next.js 16, and PostgreSQL 18.4 + PostGIS.  
**Compliance Target:** OWASP ASVS Level 3, RFC 9457 Problem Details, ISO/IEC/IEEE 29119 Testing Standards.

---

## 1. Cross-Cutting Architectural Invariants

Every subsystem must strictly enforce these foundational rules. Pull requests violating any invariant will be rejected by automated CI checks.

### 1.1. Universal Response Envelope & RFC 9457 Errors
- **Envelopes Mandatory:** Naked JSON bodies, bare arrays, and raw strings are banned. Every response must match either the Success Envelope (`success: true`) or Error Envelope (`success: false`).
- **Standard Error Format:** Error responses must implement **RFC 9457 Problem Details** with attributes: `type`, `title`, `status`, `detail`, `code`, and optional `invalidParams`.
- **Zero Stack Traces:** Server stack traces, raw SQL queries, or internal database engine errors must never appear in client error payloads. Internal errors emit `SYS_INTERNAL_ERROR` with a unique `requestId` matching backend structured logs.

### 1.2. Monetary Representation & Currency Invariant
- **Integer Minor Units:** All monetary values must be represented, stored, and calculated as 64-bit signed integers (`int64` / PostgreSQL `BIGINT`) representing the lowest non-divisible currency denomination (e.g., Santims for `ETB`, Cents for `KES` and `USD`).
- **No Floating-Point Money:** `float32`, `float64`, `double`, or untyped JSON numbers are strictly forbidden in financial calculations, structs, database columns, or API interfaces.
- **ISO 4217 Coupling:** An amount must always be paired with its three-letter uppercase ISO 4217 currency code (e.g., `{"amountMinor": 150000, "currency": "ETB"}`).

### 1.3. Concurrency Control & State Integrity
- **Optimistic Concurrency Control (OCC):** Every mutable stateful record (`bookings`, `fleets`, `owner_wallets`, `inspections`) must include a `version INT NOT NULL DEFAULT 1` column.
- **Version Increment:** Updates must execute conditionally:
  ```sql
  UPDATE bookings
  SET status = $1, version = version + 1, updated_at = NOW()
  WHERE id = $2 AND version = $3;
  ```
  If zero rows are affected, the application must immediately abort and return `BOOKING_CONCURRENCY_CONFLICT` (HTTP 409).

### 1.4. In-Transaction Financial Atomicity (`BK-17-BE`)
- **Single Transaction:** State transitions involving funds (escrow hold, booking confirmation, dispute chargeback, payout release) must execute state changes, `financial_ledger_entries` writes, and `owner_wallet_transactions` inside a single atomic database transaction (`pgx.Tx`).
- **Double-Entry Balance:** Every financial transaction must be recorded with balanced credit and debit entries summing exactly to zero minor units. Asynchronous eventual consistency is prohibited for core ledger operations.

### 1.5. Multi-Country Parameterization (`LOC-01-BE`)
- **No Hardcoded Country Constants:** Country-specific tax rates, dialing prefixes, currency codes, hold timeouts, or legal age limits must never be hardcoded into code.
- **Context Injection:** The active jurisdiction is resolved per request via the `X-Country-Code` header or user context and loaded dynamically from the `country_configs` table.

### 1.6. Soft-Delete Standard (`SD-01-BE`)
- **Retention via Timestamp:** Soft-deletable tables include `deleted_at TIMESTAMPTZ NULL`. Hard `DELETE` operations are forbidden on domain entities.
- **Query Filter:** Domain queries must explicitly filter `WHERE deleted_at IS NULL` or query the corresponding active relational view (`v_active_<table_name>`).
- **Partial Indexing:** All uniqueness constraints on soft-deletable tables must be scoped:
  ```sql
  CREATE UNIQUE INDEX uq_fleets_reg_active ON fleets (registration_number) WHERE deleted_at IS NULL;
  ```

### 1.7. Zero PII Exposure (ASVS L3)
- **Log Masking:** Phone numbers, National ID / Kebele card numbers, Tax Identification Numbers (TIN), bank account numbers, and cryptographic tokens must be masked in all log pipelines (e.g., `+251911****44`, `TIN: *****1234`).
- **Log Static Checks:** CI checks fail on any unmasked logging statements matching sensitive field names.

---

## 2. Go Backend Standards (Go 1.27.3)

The backend is built as a Clean Architecture Modular Monolith. Packages are structured with clear inbound dependencies: `domain` $\leftarrow$ `application` $\leftarrow$ `infrastructure` / `transport`.

### 2.1. Project Directory Layout
```
backend/
├── cmd/api/                  # Dependency injection and server bootstrap (main.go)
└── internal/
    ├── platform/             # Shared foundational components (telemetry, crypto, db, env)
    ├── transport/http/       # Chi router, HTTP handlers, middleware, RFC 9457 encoders
    └── modules/              # Bounded context modules
        ├── auth/
        ├── booking/
        ├── machinery/
        ├── payment/
        ├── inspection/
        ├── geo/
        └── admin/
            ├── domain/       # Pure entities, domain events, errors, value objects
            ├── application/  # Command/Query handlers, use-case interactors, repository ports
            └── infrastructure/# sqlc persistence adapters, external payment gateway clients
```

### 2.2. Package & Dependency Rules
- **Domain Independence:** The `domain` package must have zero external imports (no database drivers, no HTTP packages, no third-party libraries).
- **Interface Segregation:** Repositories and external service interfaces are declared in `application` (ports) and implemented in `infrastructure` (adapters).
- **No Circular Dependencies:** Cross-module communication must happen via published domain interfaces or structured domain events, never direct circular package imports.

### 2.3. Error Handling & Sentinel Definitions
- **Wrap Errors:** Always wrap lower-level errors with operational context using `%w`:
  ```go
  if err := r.queries.UpdateBookingStatus(ctx, params); err != nil {
      return fmt.Errorf("failed to update booking status for id %s: %w", bookingID, err)
  }
  ```
- **Domain Errors:** Define typed sentinel errors in `domain/errors.go`:
  ```go
  var (
      ErrBookingNotFound       = errors.New("booking asset does not exist")
      ErrInvalidStateTransition = errors.New("invalid state transition")
      ErrConcurrencyConflict    = errors.New("optimistic lock concurrency conflict")
  )
  ```
- **RFC 9457 Mapping:** The `transport/http` layer maps domain errors to RFC 9457 error models via a central error mapper.

### 2.4. Context, Concurrency & Goroutines
- **Context First:** `ctx context.Context` must always be the first parameter in functions performing I/O, database access, or external network requests.
- **Explicit Cancellation:** Long-running operations must monitor `ctx.Done()`. Never ignore context cancellation signals.
- **No Unmanaged Goroutines:** Raw `go func() { ... }()` calls in HTTP handlers are prohibited. Background jobs must run through a managed worker pool with panic recovery, OpenTelemetry trace propagation, and wait-group synchronization.

### 2.5. Persistence & Database Interactions
- **SQLC Exclusive:** All persistence queries must be compiled via `sqlc`. Hand-written runtime string concatenations and heavy ORMs (e.g., GORM, Ent) are strictly banned.
- **Explicit Transactions:** Transactions must be managed via explicit closures:
  ```go
  err := u.txManager.Execute(ctx, func(txQueries *sqlc.Queries) error {
      // All in-transaction calls must pass txQueries
      return nil
  })
  ```
- **Binary Parameters:** Use parameterized binding (`$1, $2`) exclusively. String interpolation in SQL is a critical security violation.

---

## 3. Flutter Mobile Standards (Flutter 3.44 / Dart)

The mobile client serves Renters, Equipment Owners, and Field Operators via a multi-persona responsive interface.

### 3.1. Layering & Clean Architecture
```
mobile/lib/
├── app/                      # Application entry, router, global theme, DI (get_it)
├── core/                     # Common widgets, network interceptors, crypto, secure storage
└── features/                 # Feature-partitioned modules
    └── <feature>/            # e.g., booking, machinery_catalog, inspections
        ├── data/             # DTOs, data sources, repository implementations
        ├── domain/           # Entities, use cases, repository interfaces
        └── presentation/     # BLoCs, screens, widgets, state models
```

### 3.2. State Management via BLoC (`flutter_bloc`)
- **Unidirectional Data Flow:** Presentation widgets emit Events $\rightarrow$ BLoCs evaluate logic $\rightarrow$ BLoCs yield immutable States $\rightarrow$ Widgets rebuild via `BlocBuilder` / `BlocConsumer`.
- **Freezed Models:** States and Events must be immutable and generated using `freezed`:
  ```dart
  @freezed
  class BookingState with _$BookingState {
    const factory BookingState.initial() = _Initial;
    const factory BookingState.loading() = _Loading;
    const factory BookingState.loaded({required Booking booking}) = _Loaded;
    const factory BookingState.failure({required ProblemDetails problem}) = _Failure;
  }
  ```
- **No Mutable State in Widgets:** `StatefulWidget` is permitted only for transient local UI state (e.g., animation controllers, scroll listeners). Domain state must live in BLoCs.

### 3.3. Networking & HTTP Interceptors (`Dio`)
- **Configured Client:** All network calls flow through a shared `Dio` singleton configured with:
  - `AuthInterceptor`: Attaches Bearer JWT; handles 401 refresh token flow.
  - `CountryCodeInterceptor`: Injects `X-Country-Code` header dynamically.
  - `IdempotencyInterceptor`: Generates and caches `Idempotency-Key` for mutation requests.
  - `ProblemDetailsInterceptor`: Deserializes non-2xx responses into typed `ProblemDetails` models.
- **Timeout Defaults:** Connect timeout: 10 seconds; Send timeout: 15 seconds; Receive timeout: 15 seconds.

### 3.4. Offline Caching & Local Storage
- **Embedded Database:** Use `isar` for local persistence (fleet listings, inspection drafts, localized metadata).
- **Sensitive Storage:** Authentication tokens and cryptographic keys must be stored in `FlutterSecureStorage` (iOS Keychain / Android KeyStore EncryptedSharedPreferences).

### 3.5. Mobile UI & Accessibility Standards
- **Touch Targets:** Interactive elements must have a minimum touch target size of $48 \times 48\text{ dp}$.
- **Text Scaling:** Interfaces must gracefully adapt to system font scaling up to $200\%$ without clipping text.
- **Dynamic Fonts:** Non-Latin script typography (`Noto Sans Ethiopic` for Amharic, `Noto Sans Arabic`) must be loaded via the dynamic OTA font pipeline (`ADR-013`).

---

## 4. Next.js Admin Panel Standards (Next.js 16 / React 19)

The administrative web console provides back-office operations, dispute settlement, and fleet verification tools using server-driven metadata schemas.

### 4.1. App Router & Server Component Discipline
- **RSC by Default:** Every page and layout in `app/` is a React Server Component (RSC) unless interactivity is required.
- **Client Boundary Minimization:** Mark components with `'use client'` only when they need browser event listeners, React hooks (`useState`, `useEffect`), or dynamic browser APIs.
- **Mutations via Server Actions:** Form submissions and operational commands must use Server Actions (`'use server'`) with audit logging.

### 4.2. Server-Driven UI Rendering (`AD-06-BE`)
- **Metadata Tables:** Operational tables (fleets, bookings, ledger entries) must not hardcode static columns. Column configurations, filter predicates, and action triggers must be rendered from `/api/v1/admin/ui-config`.
- **Headless Data Grids:** Use `@tanstack/react-table` driven by dynamic schema configurations.

### 4.3. Type Safety & Form Validation
- **OpenAPI Schema Sync:** TypeScript interfaces for API models must be generated directly from `contracts/openapi.yaml`. Hand-written API DTOs are forbidden.
- **Zod Schema Validation:** All forms must use `react-hook-form` paired with `@hookform/resolvers/zod`. Zod schemas must match API payload constraints.
- **Error Mapping:** Server action failures returning RFC 9457 `invalidParams` must automatically bind to corresponding `react-hook-form` field errors.

### 4.4. Security & Authentication
- **HTTP-Only Cookies:** Session tokens must be stored in secure, `SameSite=Strict`, `HttpOnly` cookies. Never write access tokens to `localStorage` or `sessionStorage`.
- **RBAC Gate Enforcement:** Route handlers