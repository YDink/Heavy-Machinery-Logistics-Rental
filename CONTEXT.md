# CONTEXT.md — Dink Machinery Active Session State

**Project Phase:** Phase 0 (Docs-First & Pre-Development Governance)  
**Active Requirement Partition:** `00-foundations-governance`  
**Active Branch:** `main` (Pre-bootstrap baseline)  
**Current Stop Point:** Technology Stack Version Standardization, FVM Configuration & pnpm Package Manager Enforcement Complete

---

## 1. Current Progress Snapshot

### Completed (Phase 0 Documentation & Governance)
- **Canonical Architecture & Design SSOT:** Populated `docs/source/Architecture_Design_Document.md`, `docs/source/Dink_Machinery_UI_UX_Design_Brief.md`, and complete application flow diagrams (`Heavy_Machinery_App_Flow_Diagrams.md` and `App_Flow_Diagrams_Print.html`).
- **Shared Architecture Specifications Partitioned:** Populated modular architecture chapters 01 through 14 and master `INDEX.md` under `docs/shared/architecture/` derived from canonical architecture specifications.
- **Shared Interaction Flow Specifications:** Populated complete interaction flow specifications 01 through 19, master navigation index (`INDEX.md`), and comprehensive technical appendices (Appendix A: Offline Sync & Notification Handling, Appendix B: Authentication, Token Lifecycle & Security) under `docs/shared/flows/`.
- **Shared Design Specifications & DTCG Tokens:** Populated master design system (`00-design-system.md`), user personas and ergonomics (`01-personas.md`), tokens and accessibility architecture (`02-tokens-and-accessibility.md`), and canonical W3C DTCG token definitions (`design-tokens.json`, `design-tokens.yaml`) under `docs/shared/design/`.
- **Tool Portability Established:** `AGENTS.md` finalized across all 8 disciplines; wrappers `CLAUDE.md` and `GEMINI.md` operational.
- **Contract-First Baseline:** `contracts/openapi.yaml` specified with standard envelopes and RFC 9457 error models.
- **Architectural Decision Records (L1 ADRs):** All 15 L1 system-wide Architecture Decision Records (`ADR-001` through `ADR-015`), master `INDEX.md`, and `ADR-TEMPLATE.md` authored, populated, and validated with zero template or linking errors.
- **Partition Decision Records (L2 Decisions):** All 40 partition decision files (`05-01` through `11-dynamic-config`), master `README.md`, and `PARTITION_DECISIONS_TEMPLATE.md` populated with 80 comprehensive decisions strictly meeting the 7-section layout and 5-point qualification test.
- **Shared Requirements Specifications Partitioned:** All modular requirements specifications (chapters 01 through 14, including `12-tech-stack.md`, `13-roadmap.md`, `14-appendices.md`, and master `INDEX.md`) fully populated and synchronized.
- **Mobile Subsystem Specifications Expanded:** Populated 65 bespoke screen and behavior specifications across all 5 persona contexts, industrial foundations, shared features, and QA test matrices alongside master `INDEX.md` under `docs/mobile/`.
- **Admin Subsystem Specifications Expanded:** Populated 38 dedicated operational console specifications covering the admin shell, operations queues, booking operations, dispute arbitration lightbox, double-entry financial ledger, dynamic metadata schema engine, and QA benchmarks alongside master `INDEX.md` under `docs/admin/`.
- **Zero Empty Files Verified:** Automated recursive audit confirms 0 empty (0-byte) files remain across the entire codebase.
- **Quality Assurance OS:** QA strategy, quality gates (G0–G4), security QA, and accessibility protocols articulated under `docs/quality/`.
- **Canonical Toolchain & Version Standardization:**
  - Root `.tool-versions` established (`golang 1.27.3`, `flutter 3.47.9`, `nodejs 26.13.0`, `pnpm 11.10.0`).
  - `docs/TECH_STACK.md` updated as authoritative Single Source of Truth (SSOT).
  - Pinned exact infrastructure runtimes across `docker-compose.yml`, `.env.example`, `.env.dev`:
    - Go: `1.27.3`
    - Flutter: `3.47.9` (configured and managed via FVM `.fvmrc`)
    - Node.js: `26.13.0`
    - Package Manager: `pnpm` (`v11.x`, `pnpm-lock.yaml`) mandatory across repository
    - Next.js: `16.x` (strict requirement to track latest patch release to mitigate frequent CVEs)
    - PostgreSQL: `18.6`
    - PostGIS: `3.6.4` (`postgis/postgis:18.6-3.6.4-alpine`)
    - Valkey: `9.1.3` (`valkey/valkey:9.1.3-alpine`)
- **FVM & pnpm Standardization:**
  - Configured root `.fvmrc` and `.fvm/fvm_config.json` pinning Flutter `3.47.9`.
  - Added `.fvm/flutter_sdk` and `.fvm/cache` to `.gitignore`.
  - Migrated all Flutter commands across documentation, bootstrap runbooks, CI workflows, and quality matrices to `fvm flutter ...`.
  - Enforced `pnpm` as the exclusive package manager, replacing all `npm` and `npx` references with `pnpm`, `pnpm dlx`, and `pnpm-lock.yaml`.
  - Validated architectural decision templates (`validate-decision-template.sh`) and links (`validate-decision-links.sh`) with 0 errors.

### In Progress
- Formal verification and sign-off of `docs/quality/BOOTSTRAP_CHECKLIST.md` Gate G0.

---

## 2. Next Immediate Tasks (In Order)

1. Verify and sign off Gate G0 in `docs/quality/BOOTSTRAP_CHECKLIST.md`.
2. Initialize local Git pre-commit hooks via `lefthook install`.
3. Execute Phase 1 bootstrap per `docs/BOOTSTRAP.md` using `fvm flutter ...` and `pnpm ...`.
4. Generate `docs/TECH_STACK.md` lockfile record from the bootstrap output.
5. Create initial partition feature branch: `feature/05-01-registration-authentication`.

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
