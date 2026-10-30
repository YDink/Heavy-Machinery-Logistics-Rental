# Dink Machinery — Architectural Decision Records (L1 ADRs) Master Index

**Classification:** Platform Architecture Governance & Engineering Standards Baseline  
**Governing Documents:** `AGENTS.md § 5`, `docs/AI_TASKS/04-new-adr.md`, `contracts/openapi.yaml`, `docs/quality/BOOTSTRAP_CHECKLIST.md`  
**Target Runtimes:** Go 1.27.3 Modular Monolith, Flutter 3.44 Mobile Client, Next.js 16 Admin Console, PostgreSQL 18.4 + PostGIS 3.4, Valkey 8.0, SeaweedFS S3  
**Status:** Approved Engineering Baseline

---

## 1. Decision Governance & Three-Tier Hierarchy

Dink Machinery enforces a strict three-tier technical governance hierarchy:

| Tier | Governance Level | Scope | Storage Location | Identifier Format |
|---|---|---|---|---|
| **L1** | **System ADR** | System-wide, multi-partition, foundational runtimes, high reversibility cost. | `docs/decisions/ADR-NNN-<slug>.md` | `ADR-NNN` |
| **L2** | **Partition Decision** | Localized to a single requirement partition or bounded context. | `docs/decisions/partitions/<id>.md` | `D-<partition-id>-<NNN>` |
| **L3** | **Inline Decision** | Tactical, localized, easily reversible code implementation details. | Code comments / PR review discussions | N/A |

### The 5-Point Qualification Evaluation Gate
Before authoring an L1 System ADR, the decision must evaluate against the five canonical criteria:
1. **Alternative Test:** Evaluates two or more technically viable implementation approaches with concrete trade-offs.
2. **Surprise Test:** Prevents counterintuitive choices or non-obvious architecture from perplexing future engineers.
3. **Constraint Test:** Driven by hard constraints (OWASP ASVS L3, PostGIS spatial latency, zero-custody NBE regulatory directives).
4. **Irreversibility Test:** Reversing the decision would require $> 1$ engineer-day of refactoring or a public schema breaking change.
5. **Scope Test:** Cross-cuts two or more requirement partitions or subsystems (`[backend, mobile, admin, database, contracts]`).

---

## 2. Master L1 ADR Register

| ADR ID | Decision Title | Status | Affected Partitions | Primary Rationale & Core Invariant |
|---|---|:---:|---|---|
| [`ADR-001`](./ADR-001-modular-monolith-mvp.md) | Modular Monolith in Go 1.27.3 over Distributed Microservices | Accepted | `01`, `03`, `05-05`, `05-08`, `08`, `12` | Atomic `pgx.Tx` transactions across bookings & ledger; single deployable binary. |
| [`ADR-002`](./ADR-002-rest-external-communication.md) | OpenAPI 3.1 RESTful External Communication with RFC 9457 & Minor-Unit Integers | Accepted | `02`, `05-05`, `05-09`, `05-11`, `05-18`, `14` | Provable financial precision; strict contract-first synchronization; universal error handling. |
| [`ADR-003`](./ADR-003-clean-architecture-feature-modules.md) | Strict Clean Architecture Feature Modules (`domain` $\rightarrow$ `application` $\rightarrow$ `infrastructure`) | Accepted | `05-02`, `05-03`, `05-22`, `10` | Inward dependency rule; pure domain models without database driver leakage. |
| [`ADR-004`](./ADR-004-cqrs-inside-modules.md) | Tactical In-Process CQRS and Versioned Optimistic Concurrency Control (OCC) | Accepted | `05-04`, `05-05`, `05-07`, `05-21`, `08` | High-throughput reads; OCC `version` column preventing lost updates and deadlocks. |
| [`ADR-005`](./ADR-005-pgx-sqlc-over-orm.md) | Type-Safe SQL Compilation (`sqlc` + `pgx/v5`) on PostgreSQL 18.4 + PostGIS 3.4 | Accepted | `05-02`, `05-03`, `05-04`, `05-12`, `10` | Zero runtime reflection; native binary protocol; PostGIS GiST spatial indexing. |
| [`ADR-006`](./ADR-006-bloc-flutter-state-management.md) | Flutter 3.44 Feature-First Clean Architecture with BLoC and Isar 3.1 Offline Outbox | Accepted | `05-06`, `05-12`, `05-14`, `05-16`, `05-17` | Event-driven deterministic state machines; zero-connectivity field resilience. |
| [`ADR-007`](./ADR-007-in-process-event-bus-mvp.md) | In-Process Typed Domain Event Bus with Transactional Outbox for MVP | Accepted | `05-05`, `05-08`, `05-10`, `08` | Eliminates external message broker operational overhead while preserving loose coupling. |
| [`ADR-008`](./ADR-008-dynamic-admin-ui-day-one.md) | Next.js 16 App Router (React 19 RSC) with Server-Driven Dynamic Admin UI Engine | Accepted | `04`, `05-11`, `05-18`, `05-21`, `11` | Back-office UI synthesis from `/ui-config` metadata; zero client-side token exposure. |
| [`ADR-009`](./ADR-009-opentelemetry-tracing-mvp.md) | OpenTelemetry Distributed Tracing, RED Metrics & In-Flight Zero-PII Log Scrubbing | Accepted | `05-10`, `05-24`, `09`, `12` | ASVS Level 3 privacy compliance; structured logging with in-memory Zap core masking. |
| [`ADR-010`](./ADR-010-dual-settlement-wallet-commission.md) | Non-Custodial Direct Bank Transfer Settlement & Pre-Funded `OwnerWallet` Commission | Accepted | `03`, `05-05`, `05-07`, `05-09`, `05-24` | Compliance with NBE non-banking enterprise rules; zero escrow deposit custody. |
| [`ADR-011`](./ADR-011-pluggable-identity-providers.md) | Pluggable National Identity Verification Framework (Fayda ET & IPRS KE) | Accepted | `05-01`, `05-13`, `07-04` | Extensible eKYC interface with automated score fallback to manual compliance docket. |
| [`ADR-012`](./ADR-012-phase-0-foundations-contract-first.md) | Phase 0 Contract-First Governance, Sovereign Tier III Colocation & SeaweedFS S3 | Accepted | `01`, `05-06`, `05-13`, `09`, `12` | Ethiopian Data Protection Proclamation No. 1321/2024 compliance; synchronous ClamAV. |
| [`ADR-013`](./ADR-013-server-driven-ota-localization.md) | Server-Driven OTA Translation Bundles & Dynamic Ethiopic Font Delivery Pipeline | Accepted | `05-14`, `05-16`, `05-20`, `05-23` | App bundle size optimization ($< 25\text{MB}$); dynamic Ge'ez typography injection. |
| [`ADR-014`](./ADR-014-auth-registration-partition.md) | Dual-Token ES256 JWT Session Architecture & Valkey Refresh Token Family Rotation | Accepted | `04`, `05-01`, `06-01`, `09` | Asymmetric token verification; immediate family revocation on refresh token replay. |
| [`ADR-015`](./ADR-015-multi-country-parameterization.md) | Dynamic Multi-Country Regulatory & Fiscal Parameterization via `country_configs` | Accepted | `05-09`, `05-20`, `05-23`, `11` | Zero country `if/else` branching; dynamic VAT/WHT basis points and holiday calendars. |

---

## 3. ADR Lifecycle & Traceability Governance

Every merged ADR is subject to continuous CI enforcement via `scripts/validate-decision-template.sh` and `scripts/validate-decision-links.sh`. Modifying or superseding an accepted ADR requires an approved RFC and a superseding ADR record referencing the prior identifier.
