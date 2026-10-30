# Heavy Machinery & Logistics Rental Aggregator — Dink Machinery
## Software Requirements Specification — 12. Technology Stack & Deployment View (`docs/shared/requirements/12-tech-stack.md`)

**Document Identifier:** `SRS-DINK-CORE-V7.0-SEC-12`  
**Classification:** Canonical System Requirements Specification & Infrastructure Baseline (Single Source of Truth)  
**Governing Standard:** IEEE Std 830-1998 / IEEE/ISO/IEC 29148:2018 (Systems and Software Engineering — Requirements Engineering)  
**Associated Artifacts:** `SRS-DINK-CORE-V7.0` (§12), `ARCH-DEPLOY-DINK-CORE-V7.0`, `ARCH-BACKEND-GO-DINK-CORE-V7.0`, `ARCH-FLUTTER-DINK-CORE-V7.0`, `ARCH-ADMIN-NEXT-DINK-CORE-V7.0`, `ARCH-DB-DINK-CORE-V7.0`, `contracts/openapi.yaml`, `ADR-001`, `ADR-003`, `ADR-005`, `ADR-006`, `ADR-007`, `ADR-009`, `ADR-010`, `ADR-012`  
**Target Runtimes:** Go 1.27.3 Backend Monolith, Flutter 3.47.8 / Dart 3.9 (Cross-Platform Mobile), Next.js 16.x (latest CVE security patch) / React 19 (Web Admin & Corporate Portal), PostgreSQL 18.6 + PostGIS 3.6.4, Valkey 9.1.3, ClamAV Daemon (:3310), SeaweedFS S3, Traefik 3.x  
**Target Operating Markets:** Ethiopia (`ET`) & Kenya (`KE`) — Dynamic Multi-Country Parameterization (`LOC-01-BE`)  
**Compliance Standards:** OWASP ASVS 4.0.3 Level 3, RFC 9457 (Problem Details), ISO 27001, Ethiopian Personal Data Protection Proclamation No. 1321/2024  
**Status:** Approved Engineering Baseline

---

## 1. System Deployment Topology

```
                              System Deployment Topology
                                           │
  [Mobile App / Web Users]                 ▼
                             ┌───────────────────────────┐
                             │    Traefik 3.x Ingress    │  TLS 1.3, Rate Limiting,
                             │     API Edge Gateway      │  HSTS Preloading
                             └─────────────┬─────────────┘
                                           │
                                           ├──────────────────────────────────┐
                                           ▼                                  ▼
                             ┌───────────────────────────┐      ┌───────────────────────────┐
                             │     Next.js 16 Admin      │      │     Go 1.27.3 Backend     │
                             │     Web Back-Office       │      │     Modular Monolith      │
                             │     (React Server Comp)   │      │     (Chi Router, sqlc)    │
                             └───────────────────────────┘      └─────────────┬─────────────┘
                                                                              │
                                           ┌──────────────────────────────────┼──────────────────────────────────┐
                                           ▼                                  ▼                                  ▼
                             ┌───────────────────────────┐      ┌───────────────────────────┐      ┌───────────────────────────┐
                             │      PostgreSQL 18.6      │      │        Valkey 9.1.3         │      │       SeaweedFS S3        │
                             │     + PostGIS 3.6.4 GiST    │      │     Distributed Cache     │      │     Encrypted Storage     │
                             │     (Relational DB)       │      │     & Mutex Hold Locks    │      │     & ClamAV Daemon       │
                             └───────────────────────────┘      └───────────────────────────┘      └───────────────────────────┘
```

---

## 2. Pinned Technology Matrix

| Subsystem Layer | Technology Pinned | Version Pin | Implementation Rules & Architectural Rationale |
|---|---|---|---|
| **Backend Core** | Go (`net/http`, `go-chi/v5`) | `1.27.3` | Modular monolith with Clean Architecture. High concurrency, low memory footprint ($\le 256\text{MB}$ RSS baseline). Zero runtime reflection overhead. |
| **API Gateway** | Traefik / Envoy | `3.x` | Lightweight edge routing, TLS 1.3 termination, rate-limiting, and Gzip/Brotli compression. |
| **Mobile Client** | Flutter (Dart) | `3.44.x` (Dart `3.9.x`) | Feature-First Clean Architecture, BLoC state management (`flutter_bloc`), Isar 3.1 NoSQL offline transactional store. |
| **Admin Web Console** | Next.js (React 19 App Router) | `16.x` (latest CVE patch) | React Server Components by default, Tailwind CSS v4, dynamic data grids via `@tanstack/react-table`, dynamic UI schema engine. |
| **Relational Database** | PostgreSQL + PostGIS | `18.6` + PostGIS `3.6.4` | Spatial WGS 84 (`SRID 4326`) distance queries via GiST indexing, `TIMESTAMPTZ` UTC storage, double-entry ledgers. |
| **Database Driver** | `jackc/pgx/v5` | `v5.7.1` | Native binary protocol encoding with connection pooling (`pgxpool`). Heavy ORMs strictly prohibited. |
| **SQL Compiler** | `sqlc` | `v1.27.0` | Compiles type-safe Go structs and queries directly from plain SQL files with zero runtime overhead. |
| **Cache & Distributed Locks** | Valkey | `9.1.3-alpine` | High-throughput open-source Redis replacement for mutex hold locks (`SET NX EX 3600`) and sliding-window rate limiting. |
| **Object Persistence** | SeaweedFS S3 | `3.79` | High-performance Go-based object storage for inspection photos, PDF statements, and trade licenses. |
| **Antivirus Scanning** | ClamAV `clamd` | `:3310` TCP socket | Synchronous stream scanning (`zINSTREAM`) on all uploaded media prior to S3 persistence. |
| **Telemetry & Observability** | OpenTelemetry, Zap, Loki, Tempo | Latest Stable | Structured JSON logging with Zap PII-masking core, distributed tracing, and RED metrics. |
| **Secrets Management** | OpenBao | Latest Stable | Open-source Vault fork managing column-encryption master keys and JWT signing credentials. |

---

## 3. Subsystem Constraints & Boundaries

### 3.1 Backend Monolith Architecture
- The backend compiles into a single static Linux binary (`bin/api`).
- Internal feature modules (`internal/modules/{auth, machinery, booking, payment, inspection, geo, admin}`) communicate via typed domain events or direct Go interfaces.
- Inter-package dependencies point inward toward pure domain logic; direct cross-module SQL joins or direct database mutations across module boundaries are prohibited.

### 3.2 Database Connection & Transaction Topology
- Database connection pools are managed via `pgxpool.Pool` configured with maximum 50 concurrent connections per pod.
- Mutating operations spanning bookings, inspection verifications, and financial ledger decrements execute within an explicit `pgx.Tx` transaction boundary.
- Keyset cursor pagination (`ORDER BY created_at DESC, id DESC LIMIT $1`) is enforced on all high-volume tables (`bookings`, `financial_ledger_entries`, `audit_logs`).

### 3.3 Security & Encryption Perimeters
- All network ingress terminates TLS 1.3 at the Traefik edge.
- Symmetric AES-256-GCM column encryption is applied to citizen identity numbers (National ID, TIN, phone numbers) before persistence.
- JWT access tokens are signed using asymmetric ECDSA P-256 (`ES256`) with a 15-minute expiration, paired with 30-day Valkey-tracked refresh token rotation families.
