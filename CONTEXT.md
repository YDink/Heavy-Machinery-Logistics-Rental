# CONTEXT.md — Dink Machinery Active Session State

**Project Phase:** Phase 1 Complete $\rightarrow$ Transitioning to Phase 2 (Partition Feature Development)  
**Active Requirement Partition:** `05-01-core-identity-auth`  
**Active Branch:** `main` (Ready for partition branching: `feature/05-01-core-identity-auth`)  
**Current Stop Point:** Phase 1 Machine Bootstrap Complete & Certified Across All Three Subsystems

---

## 1. Current Progress Snapshot

### Completed (Phase 1 Subsystem Bootstrap & Infrastructure)
- **Gate G0 Contracts & Governance:** OpenAPI 3.1.0 contract validated via Spectral (0 errors), 15 L1 ADRs and 80 L2 Decisions verified, zero empty files.
- **Gate G1-A Containerized Infrastructure:** PostgreSQL 18.6 + PostGIS 3.6.4, Valkey 9.1.3, SeaweedFS S3-compatible storage, and ClamAV running and healthy via Docker Compose.
- **Gate G1-B Database Baseline & PostGIS:** Migration `000001_init_platform_schema.up.sql` applied; `uuid-ossp`, `pgcrypto`, and `postgis` active; `country_configs` table created and seeded (`ET`, `KE`).
- **Gate G1-C Backend Modular Monolith (`backend/`):**
  - Go `1.27.3` module initialized with Chi v5, pgx v5, zap, and OpenTelemetry.
  - HTTP transport router bindings generated via `oapi-codegen` (`internal/transport/http/generated/api.gen.go`).
  - Repository data layer generated via `sqlc` (`internal/platform/database/db/`).
  - Health endpoint implemented (`GET /healthz`) conforming to platform universal envelope.
  - `go build` and `go test ./...` compiled and passed cleanly.
- **Gate G1-D Flutter Mobile Subsystem (`mobile/`):**
  - Flutter `3.47.9` (Dart `3.13.7`) workspace scaffolded via FVM (`com.dinkmachinery.dink_mobile`).
  - Architecture dependencies pinned: `flutter_bloc`, `bloc`, `dio`, `isar_community`, `isar_community_generator`, `freezed_annotation`, `flutter_secure_storage`, `geolocator`.
  - Feature-first Clean Architecture folders established.
  - `fvm flutter test` and `fvm flutter analyze` passed with 0 issues.
- **Gate G1-E Next.js Admin Web Console (`admin/`):**
  - Next.js `16.4.3` App Router scaffolded with TypeScript, Tailwind CSS v4, and ESLint via pnpm.
  - Administrative UI libraries installed: `@tanstack/react-table`, `@tanstack/react-query`, `react-hook-form`, `zod`, `lucide-react`, `sonner`.
  - Contract TypeScript types generated from `contracts/openapi.yaml` (`src/types/api.generated.ts`).
  - Turbopack production build (`pnpm run build`) and ESLint passed cleanly.
- **Gate G2-A Security & Zero-PII Audit:** Verified zero unmasked PII logging in application code; zero floating-point currency representations.
- **Gate G2-B Git Governance:** `lefthook` v2.2.2 pre-commit and commit-msg hooks active, intercepting commits, and validating atomic doc updates, decision links, and Conventional Commits.
- **Lockfile Determinism:** All 3 subsystems committed with reproducible lockfiles (`backend/go.sum`, `mobile/pubspec.lock`, `admin/pnpm-lock.yaml`).

### In Progress
- Opening initial Phase 2 partition feature branch: `feature/05-01-core-identity-auth`.

---

## 2. Next Immediate Tasks (In Order)

1. Create feature branch: `feature/05-01-core-identity-auth`.
2. Load partition requirements (`docs/shared/requirements/05-01-registration-authentication.md`), interaction flow (`docs/shared/flows/02-renter-registration-authentication.md`), and decisions (`docs/decisions/partitions/05-01.md`).
3. Implement Domain layer: `internal/modules/auth/domain/` (User, VerificationSession, Role, CountryCode).
4. Implement Application layer: OTP request & verification use cases with time provider injection.
5. Implement Infrastructure layer: Database repository using `sqlc`, SMS gateway adapter, JWT token generator.
6. Implement Transport layer: Register generated Chi router handlers connecting to application handlers.
7. Implement mobile auth presentation BLoC & screens in `mobile/lib/features/auth/`.
8. Implement admin identity inspection table in `admin/src/app/users/`.
9. Verify all quality gates (G0–G3) prior to rebase-merge into `main`.

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
