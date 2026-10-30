# Definition of Done (DoD) & Release Acceptance Criteria — Dink Machinery

**System:** Heavy Machinery & Logistics Rental Marketplace  
**Specification:** Canonical Engineering Completion & Acceptance Standards  
**Subsystems:** Go 1.27.3 Backend, Flutter 3.47.9 Mobile, Next.js 16 Admin, PostgreSQL 18.6 + PostGIS 3.6.4, Valkey 9.1.3  
**Governing Documents:** `AGENTS.md`, `TECH_STACK.md`, `CODE_STANDARDS.md`, `API_CONVENTIONS.md`, `ERROR_CATALOG.md`  
**Compliance Standards:** OWASP ASVS 4.0.3 Level 3, RFC 9457 (Problem Details), ISO/IEC/IEEE 29119, W3C WCAG 2.1/2.2 AA

---

## 1. Executive Summary & Philosophy of "Done"

In **Dink Machinery**, a task, requirement, or feature partition is never considered "Done" merely because code compiles or works on a local workstation. Because the platform orchestrates mission-critical physical equipment operations and multi-jurisdictional financial transactions, "Done" represents an auditable state of verified compliance with system-wide architectural invariants.

### The Five Non-Negotiable Invariants of "Done"
1. **Contract-First Synchronization:** No network interaction, schema attribute, or HTTP status code is permitted to ship unless defined in `contracts/openapi.yaml`, validated by Spectral CLI, and generated into typed models across Go, Flutter, and Next.js.
2. **Strict Financial & Minor-Unit Math:** Floating-point monetary fields are forbidden. All transactions, rates, deposits, and fee calculations must use 64-bit signed integers (`int64` / `BIGINT`) representing currency minor units coupled with ISO 4217 codes. Ledger writes must balance to zero minor units within single database transactions (`BK-17-BE`).
3. **Optimistic Concurrency & State Integrity:** Stateful entities (`bookings`, `fleets`, `wallets`, `inspections`) must enforce optimistic locking (`version` increment) and validate state machine rules before committing updates.
4. **OWASP ASVS Level 3 & Zero-PII:** Telephone numbers, National IDs, TINs, and financial accounts must be encrypted at rest (AES-256-GCM), masked in structured logs, and omitted from client-facing error envelopes.
5. **Atomic Context Tracking:** Every commit advancing system capability must atomically update `CONTEXT.md` to guarantee continuous, zero-drift operational tracking.

---

## 2. Multi-Level Definition of Done Gates

```
                            The Dink Machinery DoD Funnel
                                          │
       ┌──────────────────────────────────┼──────────────────────────────────┐
       ▼                                  ▼                                  ▼
[Level 1: Local Task &             [Level 2: Feature Partition         [Level 3: PR & Main Merge
 Subsystem Gate]                    Integration Gate]                   Verification Gate]
 - Format & Static Analysis         - Contract Code-Gen Clean           - Spectral 0 Errors/Warnings
 - Unit Tests Passing               - Testcontainers Integration        - oasdiff Zero Breaking
 - Invariant In-Line Checks         - Dynamic UI & Forms Synced         - Coverage Ratchet Pass
 - Zero-PII Regex Clean             - RFC 9457 Errors Registered        - Decision Links Validated
```

---

## 3. Level 1: Subsystem & Local Task DoD

Before any unit of work is considered complete within an individual subsystem, it must satisfy all requirements in this section.

### 3.1. Go 1.27.3 Backend Subsystem
- [ ] **Clean Architecture Boundaries:** Domain entities contain zero dependencies on transport, database drivers, or external frameworks. Interactors accept ports; adapters implement ports.
- [ ] **Type-Safe Database Interactions:** All SQL queries are compiled via `sqlc`. Raw SQL strings and runtime query concatenations are absent.
- [ ] **Parameterized Execution:** Queries use binary parameter bindings exclusively (`$1, $2`).
- [ ] **Universal Response Envelopes:** Handlers return either the platform `SuccessEnvelope` or RFC 9457 `ErrorEnvelope`.
- [ ] **Deterministic Time Handling:** System time is injected via an abstract `Clock` interface; direct calls to `time.Now()` inside domain business logic are prohibited.
- [ ] **Static Analysis & Formatting:**
  - `gofmt -s` and `goimports` applied cleanly.
  - `golangci-lint run ./...` returns 0 warnings and 0 errors.
  - `govulncheck ./...` reports 0 reachable vulnerabilities.
- [ ] **Concurrency & Race Conditions:** Tests pass with the race detector enabled: `go test -v -race ./...`.

### 3.2. Flutter 3.47.9 Mobile Subsystem
- [ ] **BLoC Unidirectional Data Flow:** Presentation components emit BLoC events; UI state rebuilds exclusively from immutable `freezed` state unions.
- [ ] **Network & Error Mapping:** All remote I/O flows through `Dio` configured with `AuthInterceptor`, `CountryCodeInterceptor`, `IdempotencyInterceptor`, and `ProblemDetailsInterceptor`.
- [ ] **Offline Schema Alignment:** Local caching tables in `isar` match corresponding backend DTO versions.
- [ ] **Secure Storage Separation:** JWT refresh tokens, WebAuthn keys, and biometrics flags reside exclusively in `flutter_secure_storage`.
- [ ] **A11y & Touch Boundaries:** All actionable components satisfy the $48 \times 48\text{ dp}$ minimum touch target (`androidTapTargetGuideline` and `iOSTapTargetGuideline`).
- [ ] **Dynamic Text Scalability:** Screen layouts accommodate $200\%$ system text scaling without overflow bars or text truncation.
- [ ] **Static Verification:** `fvm flutter analyze` reports 0 issues; `fvm flutter test` executes with 100% pass rate.

### 3.3. Next.js 16 Admin Web Console
- [ ] **React Server Components (RSC):** Pages utilize RSC by default; `'use client'` is confined strictly to leaf interactive components.
- [ ] **Server-Driven Data Grids (`AD-06-BE`):** Table columns and filter controls are loaded dynamically from `/api/v1/admin/ui-config` and rendered via `@tanstack/react-table`.
- [ ] **OpenAPI Schema Synchronization:** TypeScript interfaces are derived directly from `src/types/api.generated.ts`.
- [ ] **Validation Binding:** Forms validate via `react-hook-form` paired with `@hookform/resolvers/zod`. RFC 9457 `invalidParams` bind directly to field error messages.
- [ ] **Keyboard & Focus Integrity:** Tab navigation order is logical; focus visible rings are present; modal dialogs implement focus trapping and `Escape` key close handlers.
- [ ] **Static Verification:** `pnpm run lint` and `pnpm run build` succeed with 0 TypeScript or linting errors.

---

## 4. Level 2: Feature Partition & Integration DoD

When completing a feature partition (e.g., `05-01` Auth, `05-04` Machinery, `05-05` Booking, `05-06` Payments), the following integration standards must be met:

### 4.1. Contract Verification & Code Generation
- [ ] **OpenAPI Specification Updated:** All new or altered endpoints, schemas, parameters, and response envelopes are committed to `contracts/openapi.yaml`.
- [ ] **Spectral Lint Clean:** `pnpm dlx @stoplight/spectral-cli lint contracts/openapi.yaml --ruleset .spectral.yaml` passes with 0 errors and 0 warnings.
- [ ] **Zero Code-Gen Drift:** Code generation scripts (`generate-contracts.sh`) have been executed, and `git diff --exit-code` confirms generated files in `backend/`, `mobile/`, and `admin/` match the contract.
- [ ] **Error Registration:** Any new error codes used in the partition are registered in `ERROR_CATALOG.md` with HTTP status, RFC 9457 type URI, and retryability semantics.

### 4.2. Database & Migration Invariants
- [ ] **Reversible Migrations:** Migration pairs (`.up.sql` and `.down.sql`) are committed to `migrations/` with sequential timestamp naming. The rollback migration has been tested and executes cleanly.
- [ ] **Optimistic Concurrency Control:** Stateful tables include `version INT NOT NULL DEFAULT 1`.
- [ ] **Soft-Delete Conformance (`SD-01-BE`):** Mutable entity tables define `deleted_at TIMESTAMPTZ NULL`. Hard deletes are absent from application queries. Active views (`v_active_*`) or `WHERE deleted_at IS NULL` filters are applied.
- [ ] **Spatial Indexing:** Geographic point columns (`GEOMETRY(Point, 4326)`) are indexed with Generalized Search Trees (`GIST`).
- [ ] **Multi-Country Parameterization (`LOC-01-BE`):** Country-specific rules (tax rates, currencies, dialing codes, timeouts) are resolved dynamically from `country_configs`.

### 4.3. Financial Ledger Atomicity (`BK-17-BE`)
- [ ] **Double-Entry Balance:** All financial transactions (holds, releases, commissions, disputes, payouts) write debit and credit entries to `financial_ledger_entries` summing exactly to zero minor units.
- [ ] **Single Transaction Boundary:** State transitions and ledger writes execute within a single database transaction (`pgx.Tx`).
- [ ] **Integer Math Verification:** All intermediate fee and commission computations perform integer basis-point math with floor rounding:
  $$\text{feeMinor} = \frac{\text{amountMinor} \times \text{feeBps}}{10000}$$

### 4.4. Integration Test Suites
- [ ] **Real Engine Integration:** Integration tests run against containerized PostgreSQL 18.6 + PostGIS and Valkey instances using `testcontainers-go`. Mocking database engines in integration suites is forbidden.
- [ ] **Idempotency Verification:** Replaying mutating requests with the same `Idempotency-Key` returns the cached response envelope without duplicate side effects or duplicate ledger writes.
- [ ] **Wire Schema Conformance:** Provider tests validate HTTP responses against `contracts/openapi.yaml` schemas using `kin-openapi`.

---

## 5. Level 3: Pull Request, CI & Release Candidate DoD

Before a feature branch can be merged into `main` or promoted to an operational release candidate, it must clear all automated CI gates and architectural reviews.

### 5.1. Pull Request Hygiene & Traceability
- [ ] **Branch Naming Standard:** Branch follows the pattern `feature/<partition-id>-<slug>` (e.g., `feature/05-05-booking-engine`).
- [ ] **Conventional Commit Format:** Commits conform to `.github/commitlint.config.js`:
  ```text
  [REQ-XX-YY-ZZ] <type>(<scope>): <subject>
  ```
- [ ] **Atomic Context Update:** `CONTEXT.md` is updated and staged in the exact same commit, reflecting the updated progress snapshot, active partition status, and next tasks.
- [ ] **PR Description Template:** All sections of `.github/pull_request_template.md` are completed, including the invariant verification checklist.

### 5.2. Automated CI Gate Checklist

| CI Gate | Command / Tool | Success Threshold |
|---|---|---|
| **Contract Linting** | `pnpm dlx @stoplight/spectral-cli lint ...` | 0 errors, 0 warnings |
| **Breaking API Changes** | `oasdiff breaking origin/main:...` | 0 unapproved breaking changes |
| **Decision Templates** | `bash scripts/validate-decision-template.sh` | 0 template structure errors |
| **Decision Links** | `bash scripts/validate-decision-links.sh` | 0 broken bidirectional references |
| **Secret Scanning** | `gitleaks detect --verbose` | 0 leaks detected |
| **Static Security (SAST)** | `semgrep --config p/owasp-top-ten --error` | 0 security violations |
| **Zero-PII Logging** | Static regex audit across codebases | 0 unmasked phone/TIN/bank entries |
| **Backend Code Coverage** | `go test -coverprofile=...` | $\ge 80\%$ overall, $\ge 95\%$ on financial/auth paths |
| **Mobile Test Suite** | `fvm flutter test --coverage` | 100% pass rate, zero analyzer warnings |
| **Admin Production Build** | `pnpm run build` | 0 TypeScript errors, 0 lint warnings |
| **Accessibility Audit** | `@axe-core/playwright` | 0 WCAG 2.1 AA violations |

---

## 6. Architectural Decision Governance DoD

Whenever an architectural choice or implementation trade-off is made:

- [ ] **5-Point Qualification Evaluated:** The decision has been assessed against the Alternative, Surprise, Constraint, Irreversibility, and Scope tests.
- [ ] **Record Scaffolded:**
  - If cross-cutting: Level 1 System ADR scaffolded in `docs/decisions/ADR-NNN-<slug>.md`.
  - If partition-local: Level 2 Partition Decision scaffolded in `docs/decisions/partitions/<partition-id>.md`.
- [ ] **Required Sections Completed:**
  - L1 ADRs include all 8 mandatory sections; minimum 2 considered options; completed 5-point evaluation; explicit revisit conditions.
  - L2 Decisions include all 7 mandatory sections; minimum 2 considered alternatives.
- [ ] **Zero Boilerplate Placeholders:** All HTML comments (`<!-- ... -->`) and template placeholders are removed.
- [ ] **Bidirectional Links Verified:** The decision references its affected requirement partitions, and the requirement partitions reciprocally cite the decision ID. Verified via `scripts/validate-decision-links.sh`.

---

## 7. Definition of Done Sign-Off Checklist Template

Attach this checklist to every pull request to verify completion before requesting engineering review:

```markdown
### Definition of Done Sign-Off

#### Architectural Invariants
- [ ] Contract-First: All schemas match `contracts/openapi.yaml` (Spectral passed).
- [ ] Currency Math: All monetary fields are integer minor units (zero floats).
- [ ] In-Transaction Ledger: Double-entry balanced inside single `pgx.Tx` (`BK-17-BE`).
- [ ] Concurrency: Stateful mutations assert and increment `version` column (OCC).
- [ ] Soft Deletes: Queries filter `deleted_at IS NULL` (`SD-01-BE`).
- [ ] Multi-Country: No hardcoded country logic; parameterized via `country_configs` (`LOC-01-BE`).
- [ ] Zero-PII: Masked logging verified; zero PII in error envelopes (ASVS L3).

#### Quality & Testing
- [ ] Unit & Integration tests passing locally with `-race`.
- [ ] Code coverage satisfies threshold ($\ge 80\%$ general, $\ge 95\%$ financial/auth).
- [ ] End-to-end wire conformance validated against OpenAPI schemas.
- [ ] Mobile touch targets ($\ge 48\times 48\text{ dp}$) and $200\%$ text scaling verified.
- [ ] Web accessibility passes WCAG 2.1 AA keyboard navigation and Axe scans.

#### Documentation & Tracking
- [ ] Architectural decisions recorded (L1 ADR or L2 Decision) if 5-point test met.
- [ ] `scripts/validate-decision-links.sh` and `scripts/validate-decision-template.sh` pass.
- [ ] `ERROR_CATALOG.md` updated with any new domain error codes.
- [ ] `CONTEXT.md` updated and staged in this commit.
```