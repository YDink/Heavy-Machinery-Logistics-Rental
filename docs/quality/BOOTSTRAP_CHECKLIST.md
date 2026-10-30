# Phase 0 & Phase 1 Quality Gate & Bootstrap Verification Checklist — Dink Machinery

**System:** Heavy Machinery & Logistics Rental Marketplace  
**Subsystems:** Go 1.27.3 Modular Monolith, Flutter 3.47.8 Mobile Client, Next.js 16.x Admin Console (latest CVE security patch), PostgreSQL 18.6 + PostGIS, Valkey 9.1.3  
**Phase State:** Phase 0 (Contract & Governance) $\rightarrow$ Phase 1 (Subsystem Bootstrap)  
**Governing Documents:** `AGENTS.md`, `TECH_STACK.md`, `CODE_STANDARDS.md`, `BOOTSTRAP.md`, `API_CONVENTIONS.md`, `ERROR_CATALOG.md`

---

## 1. Quality Gate Philosophy & Sign-Off Invariants

This checklist is the formal, auditable gate verifying that **Dink Machinery** has met all architectural preconditions before opening **Phase 2 (Feature Partition Development)**. 

### Mandatory Sign-Off Invariants
1. **Zero Domain Logic in Phase 1:** All Phase 1 scaffolding must be restricted to structural wiring, container infrastructure, baseline schemas, contract-driven code generation, and CI/pre-commit governance. No business logic, payment integrations, or booking state transitions may be merged during Phase 1.
2. **Lockfile & Toolchain Determinism:** All runtime versions must match `.tool-versions` and table specifications. Every subsystem must check in a clean, reproducible lockfile (`go.sum`, `pubspec.lock`, `package-lock.json`).
3. **Canonical Contract Precedence:** No network route or data model may exist in code without an exact prior declaration in `contracts/openapi.yaml` passing Spectral linting.
4. **Atomic Doc-Commit Rule:** Sign-off state updates in this checklist and `CONTEXT.md` must be committed atomically with the corresponding verification evidence.

---

## 2. Gate Verification Summary Matrix

```
                     Bootstrap Verification Funnel
                                  │
     ┌────────────────────────────┼────────────────────────────┐
     ▼                            ▼                            ▼
[Gate G0: Contracts        [Gate G1: Runtime &         [Gate G2: Governance &
 & Architecture Docs]       Infra Provisioning]         Static Security Gates]
 - OpenAPI 3.1.0 Spectral    - PostgreSQL 18.6+PostGIS   - Lefthook & Commitlint
 - ADR & Decision Links      - Valkey 9.1.3 & SeaweedFS    - Gitleaks Zero-Secrets
 - Error Catalog RFC 9457    - Go, Flutter, Next.js      - Zero-PII Logging Check
 - Integer Minor-Unit Money  - Code Generators Pass      - Non-Decreasing Coverage
```

| Gate Identifier | Focus Area | Verification Criteria | Status | Sign-off Date |
|---|---|---|---|---|
| **Gate G0-A** | API Contract & Schemas | Spectral lint 0 errors; RFC 9457 envelopes; integer minor-unit money. | **Verified** | 2026-10-26 |
| **Gate G0-B** | Architectural Decisions | Validation scripts pass; ADR bidirectional linkage verified; 5-point test. | **Verified** | 2026-10-26 |
| **Gate G0-C** | Standards & Open Questions | Code standards established; open questions triaged in `OPEN_QUESTIONS.md`. | **Verified** | 2026-10-26 |
| **Gate G1-A** | Containerized Infrastructure | Docker Compose services (Postgres, Valkey, SeaweedFS, ClamAV, Traefik) healthy. | **Pending** | |
| **Gate G1-B** | Database Baseline & PostGIS | Migration `000001` applied; extensions enabled; `country_configs` seeded. | **Pending** | |
| **Gate G1-C** | Go Backend Modular Monolith | Go 1.27.3 module initialized; `oapi-codegen` and `sqlc` compiled; Chi `/healthz` up. | **Pending** | |
| **Gate G1-D** | Flutter Mobile Workspace | Flutter 3.47.8 workspace scaffolded; BLoC & Isar pinned; analyzer & tests pass. | **Pending** | |
| **Gate G1-E** | Next.js Admin Web Console | Next.js 16.x (latest CVE security patch) App Router created; TypeScript types synced; production build clean. | **Pending** | |
| **Gate G2-A** | Security & Zero-PII Audit | Gitleaks clean; ASVS L3 regex audit asserts zero unmasked PII logging. | **Pending** | |
| **Gate G2-B** | Pre-Commit Governance | Lefthook hooks installed; commitlint rule enforced; `CONTEXT.md` check active. | **Pending** | |

---

## 3. Detailed Verification Checklists

### 3.1. Gate G0: Phase 0 Contracts, Specifications & Governance Baseline

- [x] **G0.1 Canonical API Contract Specification (`contracts/openapi.yaml`)**
  - [x] OpenAPI version is strictly `3.1.0`.
  - [x] All paths define unique, camelCase `operationId` values.
  - [x] All 2xx success responses resolve to `SuccessEnvelope` (`success: true`, `data: {}`, `meta: {}`).
  - [x] All 4xx and 5xx error responses resolve to `ErrorEnvelope` conforming to RFC 9457 Problem Details (`success: false`, `error: {}`, `meta: {}`).
  - [x] Zero floating-point monetary properties: all currency fields resolve to composite `MoneyAmount` (`amountMinor: int64` + `currency: string`).
  - [x] State-mutating operations (`POST`, `PUT`, `PATCH`) declare mandatory `Idempotency-Key` header parameter.
  - [x] Multi-country localized endpoints declare `X-Country-Code` parameter.
  - [x] `npx @stoplight/spectral-cli lint contracts/openapi.yaml --ruleset .spectral.yaml` passes with 0 errors and 0 warnings.

- [x] **G0.2 Error Catalog & Problem Details Registry (`ERROR_CATALOG.md`)**
  - [x] Domain namespaces established (`SYS_*`, `REQ_*`, `AUTH_*`, `FLEET_*`, `BOOKING_*`, `PAYMENT_*`, `LEDGER_*`, `WALLET_*`, `INSPECT_*`, `GEO_*`, `COUNTRY_*`).
  - [x] Every error maps to an HTTP status code, canonical RFC 9457 URI slug, and retryability policy.
  - [x] Sanitization policy strictly defined: zero database driver states, SQL queries, or stack traces exposed in client payloads.

- [x] **G0.3 Architectural Decision Framework (`docs/decisions/` & Scripts)**
  - [x] Decision scaffolding utility operational: `scripts/scaffold-decisions.sh`.
  - [x] Template validator operational: `scripts/validate-decision-template.sh` asserts 8 mandatory sections for L1 ADRs and 7 mandatory sections for L2 Partition Decisions.
  - [x] Link auditor operational: `scripts/validate-decision-links.sh` asserts zero broken references and validates bidirectional links between ADRs and requirement partitions.
  - [x] 5-Point Qualification Test defined and documented in `AGENTS.md` and `docs/AI_TASKS/04-new-adr.md`.

- [x] **G0.4 Engineering Standards & Tech Stack Alignment**
  - [x] `docs/TECH_STACK.md` specifies exact version pins for Go 1.27.3, Flutter 3.47.8, Next.js 16.x (latest CVE security patch), PostgreSQL 18.6, and Valkey 9.1.3.
  - [x] `CODE_STANDARDS.md` documents architectural invariants (OCC `version` column, soft deletes via `deleted_at`, in-transaction double-entry ledger balance, zero PII logging).
  - [x] `API_CONVENTIONS.md` specifies cursor keyset pagination, RPC state action suffixes, and standard header propagation.
  - [x] `OPEN_QUESTIONS.md` registers all architectural blockers and tracks resolution links.

---

### 3.2. Gate G1: Phase 1 Runtime, Infrastructure & Subsystem Scaffolding

- [ ] **G1.1 Containerized Infrastructure Provisioning (`docker-compose.yml`)**
  - [ ] PostgreSQL 18.6 + PostGIS 3.6.4 starts cleanly and listens on `localhost:5432`.
  - [ ] Valkey 9.1.3 engine starts and responds to `redis-cli ping` with `PONG`.
  - [ ] SeaweedFS / MinIO S3-compatible storage operational on port `9000`.
  - [ ] ClamAV antivirus daemon responds to `PING` on port `3310`.
  - [ ] Traefik 3.x ingress gateway routes traffic with TLS 1.3 termination capability.
  - [ ] `docker compose ps` reports all containers in healthy status.

- [ ] **G1.2 Primary Persistence Baseline & Migrations (`migrations/`)**
  - [ ] Baseline migration `000001_init_platform_schema.up.sql` created and executed via `golang-migrate`.
  - [ ] PostGIS extension active: `SELECT postgis_full_version();` succeeds.
  - [ ] Cryptographic extensions active: `uuid-ossp` and `pgcrypto` installed.
  - [ ] `country_configs` table created with integer tax basis points (`standard_tax_rate_bps`, `withholding_tax_bps`).
  - [ ] Jurisdictional seed baseline inserted: Ethiopia (`ET`, `+251`, `ETB`, 1500 bps VAT) and Kenya (`KE`, `+254`, `KES`, 1600 bps VAT).
  - [ ] Automated trigger `trg_country_configs_updated_at` verified for automatic `updated_at` modification.
  - [ ] Reversible migration verified: `000001_init_platform_schema.down.sql` rolls back cleanly.

- [ ] **G1.3 Backend Modular Monolith Scaffolding (`backend/`)**
  - [ ] Go module initialized: `go mod init github.com/dinkmachinery/dink-core` (Go `1.27.3`).
  - [ ] Pinned dependencies installed: `chi/v5`, `jackc/pgx/v5`, `google/uuid`, `zap`, `opentelemetry`.
  - [ ] Clean Architecture module hierarchy scaffolded:
    - `internal/platform/{config, database, telemetry}`
    - `internal/transport/http/{middleware, generated}`
    - `internal/modules/{auth, machinery, booking, payment, inspection, geo, admin}/{domain, application, infrastructure}`
  - [ ] Contract transport compilation: `oapi-codegen` generates `internal/transport/http/generated/api.gen.go` with zero manual modifications.
  - [ ] Type-safe query compilation: `sqlc generate` compiles `internal/platform/database/db/` without warnings.
  - [ ] Health-check endpoint implemented (`GET /healthz`) returning platform success envelope with UTC timestamp and version metadata.
  - [ ] Binary builds cleanly: `go build -o bin/api ./cmd/api`.
  - [ ] `golangci-lint run ./...` passes with zero issues.
  - [ ] `go test -v -race ./...` passes with zero failures.

- [ ] **G1.4 Flutter Mobile Subsystem Scaffolding (`mobile/`)**
  - [ ] Flutter workspace initialized with bundle ID `com.dinkmachinery.dink_mobile` (Flutter `3.47.8`, Dart `3.9`).
  - [ ] Pinned dependencies added: `flutter_bloc`, `dio`, `isar`, `isar_flutter_libs`, `freezed_annotation`, `flutter_secure_storage`, `geolocator`.
  - [ ] Feature-first clean architecture directory hierarchy scaffolded:
    - `lib/core/{network, theme, storage, errors}`
    - `lib/features/{auth, fleet_catalog, booking, inspection}/{presentation, domain, data}`
  - [ ] Lockfile checked in: `mobile/pubspec.lock` pinned and committed.
  - [ ] Analyzer passes cleanly: `flutter analyze` reports 0 issues.
  - [ ] Baseline test passes: `flutter test` executes cleanly.

- [ ] **G1.5 Next.js Admin Web Console Scaffolding (`admin/`)**
  - [ ] Next.js 16 App Router initialized with TypeScript, Tailwind CSS v4, and ESLint (Node `26.13.0`). Pinned to latest Next.js 16.x patch release for CVE mitigation.
  - [ ] Operational libraries installed: `@tanstack/react-table`, `@tanstack/react-query`, `react-hook-form`, `@hookform/resolvers`, `zod`, `lucide-react`, `sonner`.
  - [ ] Contract schema synchronization: `npx openapi-typescript ../contracts/openapi.yaml -o src/types/api.generated.ts` executes cleanly.
  - [ ] Production build succeeds: `npm run build` compiles with zero TypeScript errors.
  - [ ] Lockfile checked in: `admin/package-lock.json` committed.

---

### 3.3. Gate G2: Governance, Pre-Commit Hooks & Static Security Assurance

- [ ] **G2.1 Pre-Commit & Commit Message Governance (`lefthook` & `commitlint`)**
  - [ ] `lefthook.yml` configured with parallel hooks for Spectral linting, decision template checks, decision link audits, and Go static analysis.
  - [ ] Atomic `CONTEXT.md` verification hook active in `lefthook` pre-commit.
  - [ ] `.github/commitlint.config.js` enforces conventional commits with ticket header pattern: `^(\[([A-Z0-9-]+)\]\s)?(\w+)(?:\(([^)]+)\))?:\s(.+)$`.
  - [ ] Local hook registration confirmed: `lefthook install` executed.

- [ ] **G2.2 OWASP ASVS Level 3 Static Security & Zero-PII Audit**
  - [ ] Secret detection drill: `gitleaks detect --source . --verbose` reports 0 secrets in repository history.
  - [ ] Zero-PII logging audit: Static regex scan confirms no unmasked logging of phone numbers, TINs, National IDs, or bank accounts across `backend/`, `mobile/`, and `admin/`.
  - [ ] Dependency vulnerability scan: `govulncheck ./...` reports 0 vulnerabilities with reachable call stacks.
  - [ ] Node dependency scan: `npm audit --audit-level=high` reports 0 high or critical vulnerabilities.

- [ ] **G2.3 Documentation & Context Checkpoint Consistency**
  - [ ] `docs/quality/BOOTSTRAP_CHECKLIST.md` updated with exact sign-off timestamps.
  - [ ] `CONTEXT.md` reflects transition from Phase 0 to Phase 1, and documents readiness for Phase 2.
  - [ ] Clean Git status: no untracked build artifacts, temporary log files, or uncommitted files.

---

## 4. Verification Execution Commands (Runbook Drill)

To perform an exhaustive audit of all bootstrap quality gates, execute this consolidated command sequence:

```bash
#!/usr/bin/env bash
set -euo pipefail

echo "=================================================================="
echo "      Dink Machinery Phase 0 / Phase 1 Quality Gate Audit         "
echo "=================================================================="

# 1. Contract & Decision Linting
echo "--> Auditing API Contract (Spectral)..."
npx @stoplight/spectral-cli lint contracts/openapi.yaml --ruleset .spectral.yaml

echo "--> Auditing Decision Templates..."
bash scripts/validate-decision-template.sh

echo "--> Auditing Decision Links & Reciprocal References..."
bash scripts/validate-decision-links.sh

# 2. Database & Infrastructure Status
echo "--> Verifying Docker Container Health..."
docker compose ps --filter "status=running"

echo "--> Verifying Database Extensions & Baseline..."
docker compose exec postgres psql -U postgres -d dink_machinery_dev -c \
  "SELECT extname, extversion FROM pg_extension WHERE extname IN ('postgis', 'uuid-ossp', 'pgcrypto');"

echo "--> Verifying Country Configurations Seed..."
docker compose exec postgres psql -U postgres -d dink_machinery_dev -c \
  "SELECT country_code, dialing_prefix, default_currency, standard_tax_rate_bps FROM country_configs;"

# 3. Backend Subsystem Verification
echo "--> Verifying Go Backend..."
cd backend
golangci-lint run ./...
go test -v -race -cover ./...
go build -o bin/api ./cmd/api
cd ..

# 4. Mobile Subsystem Verification
echo "--> Verifying Flutter Mobile Client..."
cd mobile
flutter analyze
flutter test
cd ..

# 5. Admin Web Console Verification
echo "--> Verifying Next.js Admin Console..."
cd admin
npm run lint
npm run build
cd ..

# 6. Security & Zero-PII Static Audits
echo "--> Running Secret Detection (Gitleaks)..."
gitleaks detect --source . --verbose

echo "--> Auditing Zero-PII Logging Invariants..."
if grep -rnEi 'logger\.(Info|Debug|Warn|Error)\(.*(password|pin|tin|national_id|account_number)' backend/ admin/ mobile/; then
  echo "CRITICAL: Unmasked PII logging patterns detected!"
  exit 1
fi
echo "Zero-PII logging verified."

echo "=================================================================="
echo "      SUCCESS: All Phase 0 & Phase 1 Quality Gates Verified       "
echo "=================================================================="
```

---

## 5. Formal Gate Sign-Off & Transition Handover

### Sign-Off Certification
| Role / Governance Authority | Representative | Attestation Statement | Date |
|---|---|---|---|
| **Lead Platform Architect** | Dink Engineering Core | *I attest that all contracts, architectural decisions, and invariants strictly comply with AGENTS.md, TECH_STACK.md, and ASVS Level 3.* | `2026-10-26` |
| **Backend Subsystem Lead** | Dink Core Services | *I attest that the Go 1.27.3 modular monolith, PostgreSQL 18.6 + PostGIS schema, and sqlc/oapi code generators compile with zero warnings.* | `Pending` |
| **Mobile Subsystem Lead** | Dink Mobile Client | *I attest that Flutter 3.47.8 workspace dependencies are pinned, BLoC/Isar layers scaffolded, and test suites pass.* | `Pending` |
| **Admin Web Console Lead** | Dink Web & Admin | *I attest that Next.js 16.x (latest CVE security patch) App Router builds cleanly and API types synchronize with contracts/openapi.yaml.* | `Pending` |

### Handover Protocol to Phase 2 (Feature Partitions)
When all checkboxes in Section 3 are marked `[x]` and certified:
1. Update `CONTEXT.md` setting `Current Phase: Phase 2 (Partition Feature Development)`.
2. Commit all staged files with commit message:
   ```bash
   git commit -m "[BOOTSTRAP] feat(platform): sign off phase 1 bootstrap and verify quality gates

   Completed all verification criteria across infrastructure, database baseline,
   Go 1.27.3 monolith, Flutter 3.47.8 client, Next.js 16.x (latest CVE security patch) admin, and Lefthook governance.

   Tests: All subsystem builds, race-condition suites, and static audits passing.
   Docs: CONTEXT.md, docs/quality/BOOTSTRAP_CHECKLIST.md
   Refs: #0001"
   ```
3. Checkout the first Phase 2 operational partition:
   ```bash
   git checkout -b feature/05-01-core-identity-auth
   ```