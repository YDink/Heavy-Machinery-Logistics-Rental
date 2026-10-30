# CONTEXT.md — Dink Machinery Active Session State

**Project Phase:** Phase 1 (Subsystem Bootstrap & Wiring)  
**Active Requirement Partition:** `00-foundations-governance`  
**Active Branch:** `main`  
**Current Stop Point:** Gate G0 Signed Off, Gate G1-A (Docker Infrastructure) & Gate G1-B (Migrations & PostGIS) Verified

---

## 1. Current Progress Snapshot

### Completed (Phase 0 Documentation, Governance & Baseline Infrastructure)
- **Spectral Contract Linting Verified:** Created root `.spectral.yaml` conforming to RFC 9457 and integer minor-unit money rules; `pnpm dlx @stoplight/spectral-cli lint contracts/openapi.yaml --ruleset .spectral.yaml` passes with 0 errors.
- **Architectural Decisions Verified:** `validate-decision-template.sh` and `validate-decision-links.sh` pass with 0 errors across 15 L1 ADRs and 80 L2 Partition Decisions.
- **Local Git Hooks Activated:** `lefthook` configured with commitlint, decision checks, spectral linting, and atomic doc-commit rules; hooks installed in `.git/hooks`.
- **Docker Compose Infrastructure Operational (Gate G1-A):**
  - PostgreSQL 18.6 + PostGIS 3.6.4 container (`dink-postgres`) healthy on port 5432.
  - Valkey 9.1.3 container (`dink-valkey`) healthy on port 6379.
  - SeaweedFS S3-compatible container (`dink-seaweedfs`) healthy on port 9000.
  - ClamAV container (`dink-clamav`) healthy on port 3310.
- **Persistence Baseline Applied (Gate G1-B):**
  - Migration `000001_init_platform_schema.up.sql` executed cleanly via `golang-migrate`.
  - Cryptographic extensions (`uuid-ossp`, `pgcrypto`) and `postgis` (3.6.4) active.
  - `country_configs` table created with updated-at trigger and seeded with Ethiopia (`ET`) and Kenya (`KE`).
- **Canonical Architecture & Design SSOT:** All specifications, mobile specs (65 files), admin specs (38 files), and modular requirements fully populated.
- **Canonical Toolchain & Version Standardization:**
  - Go: `1.27.3` (windows/amd64) configured.
  - Flutter: `3.47.9` (Dart `3.13.7`) with FVM (`.fvmrc`).
  - Node.js: `26.13.0` & pnpm `12.16.1`.
  - Java: OpenJDK 21 LTS configured in Flutter.
  - `flutter doctor -v`: 100% clean (`• No issues found!`).

### In Progress
- Scaffolding Go modular monolith (`backend/`), Flutter mobile workspace (`mobile/`), and Next.js admin portal (`admin/`).

---

## 2. Next Immediate Tasks (In Order)

1. Scaffold Go modular monolith (`backend/`): initialize module, run `oapi-codegen` and `sqlc generate`, build Chi `/healthz`.
2. Scaffold Flutter mobile workspace (`mobile/`): run `fvm flutter create`, pin BLoC, Isar, Dio, generate feature structure.
3. Scaffold Next.js admin app (`admin/`): run `pnpm create next-app`, install UI dependencies, sync OpenAPI types.
4. Verify non-decreasing coverage and lockfile generation (`go.sum`, `pubspec.lock`, `pnpm-lock.yaml`).
5. Open initial Phase 2 partition branch: `feature/05-01-registration-authentication`.

---

## 3. Active Decisions & Open Blockers

### Decisions Formulated
- `ADR-001` through `ADR-015` recorded and verified in `docs/decisions/`.
- `D-05-01-001` through `D-11-01-002` (80 decisions across 40 partitions) recorded and verified in `docs/decisions/partitions/`.
- Canonical Toolchain Resolution: Pinned Go `1.27.3`, Flutter `3.47.9` (via FVM), Node.js `26.13.0`, `pnpm 11.x`, Next.js `16.x` (latest CVE patch), PostgreSQL `18.6` + PostGIS `3.6.4`, and Valkey `9.1.3`.

### Blockers
- None. Phase 0 documentation, contracts, decisions, specifications, toolchains, FVM, and pnpm configurations are 100% complete and validated.

---

## 4. Files Touched in Current Session

- `.fvmrc` (created)
- `.fvm/fvm_config.json` (created)
- `.tool-versions` (updated with pnpm)
- `.gitignore` (updated for FVM)
- `AGENTS.md` (updated for FVM and pnpm)
- `lefthook.yml` (updated to pnpm dlx)
- `docs/TECH_STACK.md` (updated SSOT)
- `docs/BOOTSTRAP.md`
- `docs/quality/BOOTSTRAP_CHECKLIST.md`
- `docs/quality/QUALITY_GATES.md`
- `docs/quality/API_CONTRACT_TESTING.md`
- `docs/quality/DEFINITION_OF_DONE.md`
- `docs/quality/SECURITY_QA.md`
- `docs/quality/TEST_REPORT_TEMPLATE.md`
- `docs/CODE_STANDARDS.md`
- `docs/AI_TASKS/**`
- `docs/components/**`
- `.github/workflows/**`
- `TESTING.md`
- `CONTEXT.md`
