# Architecture Blueprint — Dink Machinery

**Target Architecture (Generated at Bootstrap):** Feature-Based Modular Monolith in Go 1.27.3, Flutter 3.47.9 Mobile, Next.js 16 Admin.  
**Canonical Source:** `docs/source/Architecture_Design_Document.md`

---

## 1. System Containers (C4 Level 2)

```mermaid
graph TD
  subgraph Clients [Client Layer]
    MobileApp["Flutter Mobile App (Renter/Owner/Operator)"]
    AdminWeb["Next.js 16 Admin Panel (Dynamic UI)"]
  end

  subgraph Gateway [Ingress Layer]
    Traefik{"Traefik API Gateway (TLS 1.3, Rate Limiting, W3C Tracing)"}
  end

  MobileApp -->|HTTPS REST / JSON| Traefik
  AdminWeb -->|HTTPS REST / Server Actions| Traefik

  subgraph Backend [Go 1.27.3 Modular Monolith]
    direction TB
    subgraph Middleware [Shared Middleware Chain]
      AuthMW["Auth & RBAC Middleware"]
      CountryMW["Country Context Middleware (LOC-01)"]
      TraceMW["OpenTelemetry Tracing Middleware"]
      ErrorMW["RFC 9457 Error Middleware"]
    end

    subgraph Modules [Logical Business Modules]
      AuthMod["auth module (Identity, OTP, MFA, Sessions)"]
      MachMod["machinery module (Fleets, Specs, EAV)"]
      BookMod["booking module (20-State Engine, Holds)"]
      PayMod["payment module (OwnerWallet, Financial Ledger)"]
      NotifMod["notification module (Prefix Router, SMS)"]
      AdminMod["admin module (Dynamic UI, CountryConfig, Audits)"]
      GeoMod["geo module (PostGIS Radius Search)"]
    end
  end

  Traefik --> Middleware
  Middleware --> Modules

  subgraph Persistence [Storage & Telemetry Layer]
    Postgres[(PostgreSQL 18.6 + PostGIS)]
    Valkey[(Valkey 9.1.3 Cache & Sessions)]
    SeaweedFS[(SeaweedFS S3 Storage)]
    Tempo["Grafana Tempo (Traces)"]
  end

  Modules --> Postgres
  Modules --> Valkey
  Modules --> SeaweedFS
  TraceMW --> Tempo
```

---

## 2. Core Architectural Invariants

### 1. In-Transaction Financial Atomicity (`BK-17-BE`)
State machine transitions that alter financial standing execute status updates, `FinancialLedgerEntry` writes, and `OwnerWalletTransaction` debits within a single `pgx.Tx` transaction protected by optimistic locking (`version` column). Financial mutations must never rely on asynchronous events.

### 2. Multi-Country Parameterization (`LOC-01-BE`, `ADR-015`)
Incoming requests resolve an active `country_code` via the `X-Country-Code` header or user context. Dialing codes, currency symbols, minor units, tax rates, and regulatory retention periods are dynamically resolved from `country_configs`.

### 3. Server-Driven Admin UI (`AD-06-BE`)
The Next.js administrative dashboard renders navigation hierarchies, table columns, filter rules, and action triggers dynamically from metadata fetched from `/api/v1/admin/ui-config`.

---

## 3. Subsystem Architectural Details

Detailed component designs are maintained in:
- Backend Architecture: `docs/shared/architecture/04-backend-architecture-go.md`
- Mobile Architecture: `docs/shared/architecture/05-flutter-mobile-app-architecture.md`
- Admin Web Architecture: `docs/shared/architecture/06-nextjs-admin-panel-architecture.md`
- Database Design: `docs/shared/architecture/07-database-design.md`
- Security Architecture: `docs/shared/architecture/09-security-architecture.md`