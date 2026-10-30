# Technology Stack & Toolchain Specification — Dink Machinery

**System:** Heavy Machinery & Logistics Rental Marketplace  
**Architecture:** Feature-Based Modular Monolith (Go Backend), Multi-Tenant Mobile (Flutter), Dynamic Web Console (Next.js)  
**Governance:** Strict Version Pinning & Lockfile Discipline (`ADR-002`, `AGENTS.md § 6`)  
**Single Source of Truth (SSOT):** This document and root `.tool-versions` define canonical versions for all subsystems.

---

## 1. Core Runtimes & Infrastructure Engines

| Layer | Technology | Version | Pinning Strategy | Purpose & Architectural Role |
|---|---|---|---|---|
| **Backend Runtime** | Go | `1.27.3` | `.tool-versions`, `go.mod` | High-throughput, low-latency modular monolith backend service. |
| **Mobile Runtime** | Flutter / Dart | `3.47.8` / `3.9.x` | `.tool-versions`, `pubspec.lock` | Cross-platform client for Renters, Asset Owners, and Machine Operators. |
| **Admin Web Runtime** | Node.js / Next.js | `26.13.0` / `16.x` (latest CVE patch) | `.tool-versions`, `package-lock.json` | Metadata-driven administrative and operations console. Requires latest patch of Next.js 16 to mitigate frequent CVEs. |
| **Primary Persistence** | PostgreSQL + PostGIS | `18.6` (PostGIS `3.6.4`) | `docker-compose.yml` (`postgis/postgis:18.6-3.6.4-alpine`) | Relational data, spatial radius indexing, immutable financial ledger. |
| **Cache & Ephemeral** | Valkey | `9.1.3` (`9.1.3-alpine`) | `docker-compose.yml` (`valkey/valkey:9.1.3-alpine`) | Distributed locks, session storage, rate-limiting counters, hot cache. |
| **Object Storage** | SeaweedFS / MinIO | S3-Compatible API | Container image tag | Contract PDFs, inspection media, dynamic font OTA assets. |
| **Malware Inspection** | ClamAV Daemon | Latest stable | Container image tag | Synchronous scanning of uploaded inspection photos and legal docs. |
| **Edge Ingress Gateway** | Traefik | `3.x` | Container image tag | TLS 1.3 termination, rate limiting, and W3C trace propagation. |

---

## 2. Backend Subsystem (Go 1.27.3 Modular Monolith)

The backend follows Clean Architecture principles with Domain-Driven Design (DDD) module boundaries. Domain logic is decoupled from persistence, network protocols, and third-party SDKs.

```
backend/
├── cmd/api/                  # Application entry points and dependency wiring
├── internal/
│   ├── platform/             # Cross-cutting concerns (DB, telemetry, crypto, country)
│   ├── transport/http/       # HTTP handlers, Chi router, RFC 9457 middleware
│   └── modules/              # Bounded context modules (auth, machinery, booking, payment, admin, geo)
│       └── <module>/
│           ├── domain/       # Pure domain entities, invariants, value objects (no external imports)
│           ├── application/  # Use case command/query handlers
│           └── infrastructure/# sqlc repositories, external gateway adapters
```

### Core Libraries & Tooling
- **HTTP Routing & Ingress:** `go-chi/chi/v5` — Lightweight, composable HTTP router adhering strictly to standard `net/http` primitives.
- **Database Driver & Connection Pool:** `jackc/pgx/v5` — Native PostgreSQL wire protocol implementation with connection pooling, structured logging, and binary parameter binding.
- **Query Generation & Type Safety:** `sqlc` (`v1.27+`) — Compiles validated SQL queries into type-safe, idiomatic Go structs. Raw runtime SQL strings and traditional ORMs (e.g., GORM, Ent) are prohibited.
- **Migrations:** `golang-migrate/migrate/v4` — Sequential, directional (`up`/`down`) SQL migration runner integrated into local bootstrapping and CI release pipelines.
- **Contract Code Generation:** `oapi-codegen` — Compiles `contracts/openapi.yaml` into Go request binding schemas, router interfaces, and RFC 9457 error contracts.
- **Validation:** `go-playground/validator/v10` — Edge transport schema validation enforcing payload boundaries prior to application dispatch.
- **Cryptography:** Go standard library `crypto/cipher`, `crypto/aes` — AES-256-GCM authenticated encryption for PII and sensitive payment keys at rest.
- **Linter & Static Analysis:** `golangci-lint` (`v1.64+`) — Comprehensive suite enforcing `errcheck`, `govet`, `staticcheck`, `gosec`, and custom AST-based zero-PII logging rules.

---

## 3. Mobile Subsystem (Flutter 3.47.8)

The mobile client delivers a unified experience across Renter, Owner, and Operator personas using dynamic role gating and responsive layouts.

### Architecture & Key Libraries
- **Runtime:** Flutter `3.47.8` (Dart `3.9.x`) pinned via `.tool-versions`.
- **State Management:** `flutter_bloc` (`v9.x`) — Predictable, unidirectional data flow separating presentation components from business logic via discrete Events and States.
- **Networking & API Client:** `dio` (`v5.x`) — HTTP client configured with custom interceptors for Bearer token refresh, `Idempotency-Key` generation, `X-Country-Code` injection, and RFC 9457 problem detail parsing.
- **Offline Storage & Keyset Caching:** `isar` (`v4.x`) — High-performance, embedded NoSQL database for cached fleet listings, inspection drafts, and localized dictionary records.
- **Secure Token Storage:** `flutter_secure_storage` — Encrypted storage backed by iOS Keychain and Android KeyStore for refresh tokens and biometrics flags.
- **Immutable Models:** `freezed` + `json_serializable` — Type-safe, compile-time code generation for union states, copy mutations, and JSON mapping.
- **Location & Geofencing:** `geolocator` — Accurate device location acquisition for machine discovery and check-in radius verification.
- **Typography & Localization:** OTA dynamic font loader fetching non-Latin script files (`Noto Sans Ethiopic`, `Noto Sans Arabic`) from S3 endpoints dynamically (`ADR-013`).

---

## 4. Admin Web Console (Next.js 16.x - Latest CVE Security Patch, Node.js 26.13.0)

The administrative dashboard provides operations teams, dispute managers, and platform engineers with control over fleet onboarding, financial ledgers, and country configurations.

### Architecture & Key Libraries
- **Core Framework:** Next.js `16.x` (React 19+, App Router) on Node.js `26.13.0`. Always pinned to the latest patch release of Next.js 16 to protect against active CVEs. Full utilization of React Server Components (RSC) for streaming reads and Server Actions for audit-logged mutations.
- **Design System & Styling:** Tailwind CSS (`v4.x`) + `shadcn/ui` — Utility-first styling paired with accessible, copy-pasteable Radix UI primitives.
- **Component Data Grids:** `@tanstack/react-table` (`v8.x`) — Headless table engine powering server-driven dynamic column mappings, filtering, sorting, and keyset pagination.
- **Form Management & Validation:** `react-hook-form` paired with `@hookform/resolvers/zod` — Runtime schema evaluation matching OpenAPI request definitions.
- **State & Server Synchronization:** `@tanstack/react-query` (`v5.x`) — Cache invalidation and client-side background polling for real-time operational feeds.
- **Icons & Visuals:** `lucide-react` — Standard icon set strictly mapped to design token specifications.

---

## 5. Storage, Persistence & Query Models

### PostgreSQL 18.6 + PostGIS 3.6.4
- **Spatial Geometry:** Machinery locations and job sites are indexed using PostGIS spatial points (`GEOMETRY(Point, 4326)`).
- **Spatial Indexing:** 2D R-Tree spatial indices (`GIST`) enable sub-millisecond radius search operations:
  ```sql
  CREATE INDEX spat_idx_fleets_location ON fleets USING GIST (location_point);
  ```
- **Concurrency & Invariants:**
  - Optimistic locking using sequential `version INT NOT NULL DEFAULT 1` columns on all mutable entities.
  - Double-entry immutable financial ledger writes executing inside single database transactions (`BK-17-BE`).
  - Standard soft-delete filtering via `deleted_at IS NULL` indices and active relational views (`v_active_*`).

### Valkey 9.1.3 Engine
- **Session Tokens:** Active session metadata, WebAuthn challenge salts, and device fingerprint validations stored with strict TTL expiration.
- **Ingress Rate-Limiting:** Sliding-window counters enforcing tier-based limits per client IP and authenticated principal.
- **Distributed Hold Locks:** Short-term distributed locks preventing concurrent double-hold requests on identical fleet assets.

---

## 6. Observability, Security & Quality Infrastructure

### Telemetry & Monitoring
- **Tracing:** W3C Trace Context propagated via OpenTelemetry (`go.opentelemetry.io/otel`). Traces stream over OTLP gRPC to Grafana Tempo.
- **Metrics:** Native Prometheus metrics exported over `:9090/metrics` tracking RED metrics (Rate, Errors, Duration), database pool saturation, and state machine transaction latency.
- **Structured Logging:** `uber-go/zap` emitting newline-delimited JSON with automatic masking for telephone numbers, TIN, and banking identifiers.

### Security Tooling (OWASP ASVS Level 3)
- **Static Analysis (SAST):** `semgrep` scanning for OWASP Top 10, Go memory boundaries, and Next.js injection surfaces.
- **Secret Detection:** `gitleaks` executed locally via `lefthook` and enforced in CI to block credential leaks.
- **Vulnerability Auditing:** `govulncheck` (Go toolchain) and `npm audit` (Next.js dependencies).
- **Container Hardening:** `aquasecurity/trivy` inspecting base container images and Docker Compose service configurations.

---

## 7. Version Pinning & Dependency Lockfile Protocol

1. **Deterministic Builds:** Every environment (local, CI, staging, production) must build exclusively against pinned lockfiles:
   - Go: `backend/go.mod` + `backend/go.sum` (Go `1.27.3`)
   - Flutter: `mobile/pubspec.yaml` + `mobile/pubspec.lock` (Flutter `3.47.8`)
   - Next.js: `admin/package.json` + `admin/package-lock.json` (Next.js `16.x` latest CVE patch, Node.js `26.13.0`)
2. **Runtime Consistency:** Developer machines and CI runners align runtimes using root `.tool-versions` (managed via `asdf` or `mise`).
3. **Major Version Governance:** Upgrading major versions of Go, Flutter, Next.js, PostgreSQL, or Valkey requires an approved **L1 ADR** evaluating breaking changes, database migration compatibility, and performance profiles.
