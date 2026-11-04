# AGENTS.md — Dink Machinery Engineering Context & Operating Rules

System: Heavy machinery & logistics rental marketplace.
Operating Model: Docs-first (Phase 0) → Bootstrap (Phase 1) → Partition feature development (Phase 2).
Tool Portability: Tool-agnostic specification. Root and component wrappers (`CLAUDE.md`, `GEMINI.md`) import this file via `@AGENTS.md`.

---

## 1. QA & Test Gates

### Test Strategy & Determinism
- **Test-Driven Design:** Tests are written alongside or prior to implementation. Commits must never introduce production code without corresponding test coverage.
- **Determinism:** Tests must be strictly isolated. Flaky tests (dependent on wall-clock time, unseeded pseudo-randomness, network calls, or shared mutable state) are quarantined within 24 hours.
- **Time Freezing:** Any time-sensitive domain logic (booking timeouts, grace periods, payment verification windows) must inject an abstracted time provider.

### Mandatory Pre-Commit Verification Triad
Every code contribution requires completion of the three-pillar verification before staging or committing:
1. **Automated Unit Testing:** Pure business invariants, state machine transitions, financial math, and BLoC state streams ($\ge 80\%$ line coverage overall; $\ge 95\%$ on financial/auth paths).
2. **Automated Integration Testing:** Repository persistence, SQL execution, database migrations, and HTTP transport handlers.
3. **Mandatory Manual Verification (Interactive Testing):** Every newly added or modified capability must be manually executed and validated live before commit:
   - **Backend APIs:** Live execution of endpoints via real HTTP requests (e.g., `curl` or interactive probe harnesses) verifying 2xx success envelopes, 4xx/5xx RFC 9457 error details, and database mutation side effects.
   - **Database Migrations:** Manual execution and verification of both `up` and `down` migration scripts (rollback drill) to guarantee zero SQL syntax errors or relational constraint collisions.
   - **Mobile UI & BLoC:** Manual execution and visual inspection of widgets, forms, 56dp touch targets, input validation error feedback, and bilingual language toggling (English ↔ Amharic).
   - **Manual Evidence Requirement:** Commits must include an itemized `Manual-Tests:` section documenting the exact manual test scenarios executed and verified.

### Coverage Ratchet
- **Target:** Minimum 80% line coverage on changed code; 95% branch coverage required for financial calculations, ledger operations, auth paths, and payout workflows.
- **Ratchet Rule:** Code coverage metrics are non-decreasing. PRs reducing overall test coverage will fail CI.

### Quality Gate Hierarchy (G0 to G4)
- **G0 (Pre-Commit / Local):** Linter execution, static analysis, unit test suite, integration test suite, mandatory manual verification drill, and atomic doc-commit verification (`CONTEXT.md` and `REQUIREMENTS_TRACKER.md` staged via `lefthook.yml`).
- **G1 (Pull Request):** Full automated test suite, coverage ratchet checks, contract validation tests, security static analysis, and documentation drift checks (`ci-docs.yml` in Phase 0; component CI in Phase 2).
- **G2 (Pre-Merge):** Rebase verification against `main`, semantic commit checks, bidirectional decision link audits, and clean security scans.
- **G3 (Pre-Release):** Full End-to-End (E2E) regression, synthetic load/soak performance profiles, accessibility validation (axe-core zero-critical rule), and database migration rollback drills.
- **G4 (Production):** Live synthetic health monitoring, RED metrics verification, and automated rollback triggers on error rate threshold breach.

### Security Testing
- **Zero PII Leakage:** Automated assertions must ensure names, phone numbers, bank details, and government IDs never enter standard logs, analytics payloads, or RFC 9457 error details.
- **Ledger Invariants:** Financial mutations must run automated double-entry verification tests confirming zero-sum ledger invariants across debit and credit legs.
- **Idempotency Assertions:** Replay tests are mandatory for all payment, booking confirmation, and refund endpoints using identical `Idempotency-Key` headers.

---

## 2. API Standards

### Unified Response Envelope
Every JSON endpoint must return the standardized response envelope. Unwrapped entities, bare arrays, or raw status messages are prohibited.

#### Success Envelope (HTTP 200, 201, 202)
```json
{
  "success": true,
  "data": {},
  "error": null,
  "meta": {
    "version": "1.0",
    "requestId": "req_01HP8XYZ12345678",
    "timestamp": "2026-03-31T08:30:00Z",
    "pagination": null,
    "warnings": []
  }
}
```

#### Error Envelope (RFC 9457 Problem Details)
All client-facing errors (4xx and 5xx) must serve `application/problem+json` semantics wrapped inside the standard structure.
```json
{
  "success": false,
  "data": null,
  "error": {
    "type": "[https://api.dinkmachinery.com/errors/invalid-booking-state](https://api.dinkmachinery.com/errors/invalid-booking-state)",
    "title": "Invalid Booking State Transition",
    "status": 409,
    "detail": "Cannot transition booking from CANCELLED to CONFIRMED.",
    "code": "BOOKING_INVALID_STATE_TRANSITION",
    "invalidParams": [
      {
        "name": "status",
        "reason": "Current state is terminal."
      }
    ]
  },
  "meta": {
    "version": "1.0",
    "requestId": "req_01HP8XYZ12345678",
    "timestamp": "2026-03-31T08:30:00Z"
  }
}
```

### Money & Currency Representation
- **Strictly String or Integer Minor Units:** Never transmit or calculate money using floating-point types (`float`, `double`, or untyped JSON numbers).
- **Format:** Either an integer representing minor units (e.g., `10050` with an explicit currency code representing 100.50) or fixed-point string formatting (`"100.50"`).
- **ISO 4217 Mandatory:** Every monetary amount must be accompanied by its ISO 4217 three-letter currency code (e.g., `ETB`, `USD`).

### Idempotency Enforcement
- Mutating endpoints (`POST`, `PUT`, `PATCH`) involving financial operations, state machine transitions, or side-effects require the `Idempotency-Key` header (UUID v4 or ULID).
- Keys are persisted and verified inside the primary database transaction alongside the mutation. Duplicate submissions within the retention window must return the cached original response.

### Pagination
- **Cursor-Based Only:** Real-time and operational datasets (bookings, ledger entries, fleet browsing) must strictly utilize cursor-based pagination. Offset/limit pagination is prohibited on financial and audit tables.
- Pagination metadata must be passed via `meta.pagination` containing `nextCursor`, `prevCursor`, and `hasMore`.

### Naming & Routing Conventions
- **Resource URIs:** Plural kebab-case (e.g., `/api/v1/fleet-equipments`).
- **Nesting Cap:** Maximum two levels of URI hierarchy (e.g., `/api/v1/bookings/{id}/inspections`). Operations outside direct hierarchy must use query filters or root paths.
- **RPC/Action Style:** Non-CRUD state modifications must append an explicit action tag (e.g., `/api/v1/bookings/{id}:confirm-payment`).
- **Casing Isolation:**
  - URIs: `kebab-case`
  - JSON payload fields: `lowerCamelCase`
  - Database schema/tables/columns: `snake_case`
  - Constants/Enums: `UPPER_SNAKE_CASE`

---

## 3. Code Standards

### Cross-Stack Invariants
- **Consult `CODE_STANDARDS.md`:** All implementations must conform to the project standards defined in `docs/CODE_STANDARDS.md`.
- **Strict Single-Concept File Separation ("No Shared/Lumped Files"):**
  - Shared, catch-all, or monolithic files are strictly forbidden across every tier (`entity.go`, `models.go`, `types.go`, `helpers.go`, `utils.go`, `invariants.go`, `repository.go`, `memory_repository.go`, `postgres_repository.go`, monolithic `*_handler.go`).
  - Every single entity, value object, invariant rule, model, repository interface, repository implementation, use case/service command, HTTP route handler, UI widget, and test MUST reside in its own dedicated, specifically named file in the appropriate directory.
  - "When we say everything, we mean everything." Zero exceptions across Backend, Mobile, Admin, Contracts, and Tests.
- **UI Design-First with Stitch MCP:**
  - **Just-In-Time Generation per Requirement:** Designs must be created individually via Stitch MCP at the exact moment when the development of that specific requirement/feature is initiated, rather than generating all partition screens at once in bulk.
  - **Design System & Flow Continuity:** Every generated screen must adhere to the core design requirements (high-contrast outdoor palette `#111827`, $\ge 7:1$ contrast ratio, $\ge 56\text{dp}$ touch targets) and maintain seamless visual and ergonomic flow from the preceding screen in the user journey (referencing canonical user flows in `docs/shared/flows/` and preceding Stitch screens).
  - **Implementation Target:** Frontend UI code (Flutter widgets/screens and Next.js admin components) must be built directly against and faithfully mirror the generated Stitch designs.
- **No Unapproved Paradigms:** Do not introduce architectural patterns, global state solutions, or external dependencies without an explicit L1 (ADR) or L2 (Partition) decision doc.
- **DTO Isolation:** Database models and persistence entities must never be exposed directly to the controller layer or client responses. DTOs and view models are mandatory.
- **Layered Validation:** Validate format, presence, and boundary constraints at the edge/controller layer. Validate business invariants and relational constraints in the domain/service layer.

### Go Backend (Modular Monolith)
- Package dependencies point inward: `domain` → `application` → `infrastructure`/`presentation`.
- The domain layer contains pure business entities, value objects, and repository interfaces. No SQL queries or external driver imports within domain packages.
- Mutations use command handlers returning `Result[T, E]` or error types. Queries use specialized query handlers returning read-optimized DTOs.
- `sqlc` is the standard for SQL generation. Raw SQL strings within application logic are prohibited.
- HTTP routing relies strictly on pure Go standard library `net/http.ServeMux` and `log/slog` (zero third-party web frameworks).
- **Strict Granular File Layout:**
  - `domain/entities/<entity_name>.go`: Strictly one entity per file.
  - `domain/invariants/<invariant_name>.go`: Strictly one business invariant rule per file.
  - `domain/repositories/<repository_name>.go`: Strictly one repository interface per file.
  - `infrastructure/persistence/postgres/<entity_name>_postgres_repository.go`: Strictly one postgres repository adapter per file.
  - `infrastructure/persistence/memory/<entity_name>_memory_repository.go`: Strictly one in-memory repository adapter per file.
  - `application/usecases/<command_or_query_name>.go`: Strictly one command/query use case per file.
  - `transport/http/handlers/<endpoint_action>_handler.go`: Strictly one route/action handler per file.
  - Tests mirror implementation files 1:1 (`<name>_test.go`). No monolithic test files.

### Flutter Mobile (FVM)
- Subsystem SDK and toolchain are managed strictly via FVM (.fvmrc). All Flutter commands execute via `fvm flutter ...`.
- State management relies strictly on the BLoC pattern (`flutter_bloc`). UI widgets emit events; state changes are received via typed states.
- **UI Design via Stitch MCP:** Each screen and interactive sheet must be designed just-in-time via Stitch MCP when active requirement work begins, maintaining visual and flow continuity from the preceding screen.
- Clean separation & granular file layout:
  - `domain/entities/<entity_name>.dart`: One entity per file.
  - `domain/repositories/<repository_name>.dart`: One repository interface per file.
  - `domain/usecases/<usecase_name>.dart`: One use case per file.
  - `data/models/<model_name>.dart`: One model/DTO per file.
  - `data/datasources/<datasource_name>.dart`: One remote/local data source per file.
  - `data/repositories/<repository_impl_name>.dart`: One repository implementation per file.
  - `presentation/bloc/<feature>_bloc.dart`, `<feature>_event.dart`, `<feature>_state.dart`.
  - `presentation/views/<screen_name>_view.dart`: One view per file.
  - `presentation/views/widgets/<widget_name>.dart`: One component/widget per file.
  - Tests mirror files 1:1 (`<name>_test.dart`).
- All interactive touch targets must be $\ge 56\text{dp}$.
- Network calls are routed through a configured `Dio` client with logging, auth refresh, and error-to-problem-details interceptors.

### Next.js Admin Panel (pnpm)
- Package management relies strictly on `pnpm` (pnpm-lock.yaml). Use of `npm` or `yarn` is prohibited.
- **UI Design via Stitch MCP:** Administrative dashboards, tables, and modal workflows must be designed via Stitch MCP before frontend implementation.
- Granular file separation applies: strictly one component/server-action/schema/hook/test per file. No multi-component or catch-all files.
- Next.js App Router architecture utilizing React Server Components (RSC) for reads and Server Actions for mutations.
- Dynamic administrative panels are server-rendered from metadata schemas; no client-side hardcoded tables or dynamic form structures.
- Forms utilize `react-hook-form` paired with strict `zod` schemas matching OpenAPI contracts.

---

## 4. Git Standards

### Branch-per-Partition Lifecycle
- Development maps directly to SRS partitions. Feature branches follow the pattern:
  `feature/<partition-id>` (e.g., `feature/05-05-booking-reservation-flow`).
- Commits on the partition branch represent atomic micro-commits targeting specific requirement IDs.
- Merging to `main` is executed via rebase-merge to maintain clean linear history and preserve bisectability. Direct pushes to `main` are prevented.

### Commit Format & Structure
Commits must conform to the extended Conventional Commits specification. Every commit must explicitly track corresponding requirement and document updates.

```text
[<REQ-ID>] <type>(<scope>): <short summary>

<detailed description of change, rationale, and consequences>

Tests: <test files or suites executed/added>
Manual-Tests: <itemized manual test steps and verification evidence>
Docs: <list of documentation files updated, e.g., CONTEXT.md, decisions/partitions/05-05.md>
Refs: #<issue-or-task-number>
```

### The Atomic Doc-Commit Rule
- Code and tests must never be committed without staging their operational documentation simultaneously.
- **Mandatory in every commit:** Code changes + Automated Tests + Manual Test Verification Evidence + `CONTEXT.md` + `REQUIREMENTS_TRACKER.md`.
- **Conditional:** If a design decision meets any qualification test, the relevant partition decision doc or ADR must be updated and staged within the exact same commit.

### Pre-Commit Gate Enforcement (`lefthook.yml`)
Local git hooks block commits failing the following:
1. `G0` linter, code style, unit, and integration test runs.
2. Verification that `CONTEXT.md` and `REQUIREMENTS_TRACKER.md` are present in the staged index.
3. Verification that commit messages include the `Docs:` and `Manual-Tests:` footers.
4. Validation that any decision mentioned in commit bodies exists in `docs/decisions/`.
5. Bypassing checks via `--no-verify` is strictly prohibited on all development branches.

---

## 5. Decision Documentation

### Three-Tier Decision Model
All technical and architectural choices must be categorized and recorded according to the appropriate governance level.

| Tier | Scope | Location | ID Format |
|---|---|---|---|
| **L1 — ADR** | Cross-cutting, multi-partition, system-wide, high cost to reverse | `docs/decisions/ADR-NNN-<slug>.md` | `ADR-NNN` |
| **L2 — Partition** | Single partition, non-trivial, localized trade-offs | `docs/decisions/partitions/<partition-id>.md` | `D-<partition-id>-<NNN>` |
| **L3 — Inline** | Obvious, fully reversible, tactical code implementation | Code comments or PR review threads | N/A |

### The 5-Point Qualification Test
If **ANY** of the following tests evaluate to `TRUE`, the decision **MUST** be recorded in the partition decision doc (or escalated to an ADR if the Scope Test is true):
1. **Alternative Test:** Were two or more viable implementation paths considered?
2. **Surprise Test:** Would a future engineer or AI agent ask "Why was it done this way?"
3. **Constraint Test:** Was this design selected due to a non-obvious constraint (performance, memory, compliance, third-party limitation)?
4. **Irreversibility Test:** Would reversing this decision require more than one engineer-day of work or impact external consumers?
5. **Scope Test:** Does this decision affect more than one requirement partition? *(If TRUE, escalate to L1 ADR).*

### Decision Artifact Structure
Every L1 and L2 decision entry must strictly document:
- Context and problem statement.
- Minimum of 2 considered options.
- Chosen solution and technical rationale.
- Explicit list of rejected alternatives with detailed justification.
- Consequences (positive, negative, operational).
- Verification strategy (exact test cases confirming the decision).
- Revisit conditions (under what metrics or scale thresholds this must be reconsidered).

---

## 6. Bootstrap Policy

### Phased Development Progression
- **Phase 0 — Documentation & Contract Definition:** The current project phase. Only documentation, contracts, test strategies, prompt templates, and infrastructure declarations exist. No implementation directories (`backend/`, `mobile/`, `admin/`) may be initialized.
- **Phase 1 — Machine Bootstrapping:** Triggered only when `docs/quality/BOOTSTRAP_CHECKLIST.md` is completely verified. Frameworks, package manifests, and initial project structures are established using latest stable tools.
- **Phase 2 — Partition Feature Work:** Iterative development executed partition-by-partition on dedicated feature branches on standard Git discipline.

### Version Management & Pinning
- Documentation files must never hardcode ephemeral patch or framework versions.
- At bootstrap execution, tooling will fetch the latest stable LTS releases of runtimes, libraries, and frameworks.
- Once created, runtime versions and package trees are pinned via lockfiles (`go.sum`, `pubspec.lock`, `pnpm-lock.yaml`, and `.tool-versions`). Any subsequent major version bump requires an L1 ADR.

---

## 7. Documentation Sync Rule

### Source-of-Truth Hierarchy
1. `docs/source/`: The canonical, complete reference specifications (`Heavy_Machinery_Rental_SRS_v2.md`, `Architecture_Design_Document.md`, `Heavy_Machinery_App_Flow_Diagrams.md`, `Dink_Machinery_UI_UX_Design_Brief.md`).
2. `docs/shared/`: Derived modular partitions mapped directly to implementation partitions.
3. Contracts and Code: Downstream artifacts derived from requirements and architecture.

### Maintenance Protocol
- Never modify files inside `docs/shared/requirements/`, `docs/shared/architecture/`, or `docs/shared/flows/` directly without updating the corresponding parent section in `docs/source/`.
- CI pipelines execute automated drift detection (`ci-docs.yml`). Any discrepancy between `docs/source/` and derived documentation will fail build verification.

---

## 8. Context Loading Rule

### Context Window Optimization
- Avoid bulk-loading the entire documentation tree into active AI context windows.
- Load only the baseline operating files: `AGENTS.md`, `CONTEXT.md`, `REQUIREMENTS_TRACKER.md`, and `contracts/openapi.yaml`.
- Pull specific partition files on demand based on the active task:
  - Requirements: `docs/shared/requirements/<partition-id>.md`
  - Design/Flows: `docs/shared/flows/<flow-id>.md`
  - Decisions: `docs/decisions/partitions/<partition-id>.md`
- Always verify open blockers in `docs/OPEN_QUESTIONS.md` before executing new architectural implementations.

---

## 9. Requirement Traceability & Tracking Protocol

### Status Lifecycle Management
All partitions and requirements follow a strict three-state lifecycle tracked in `REQUIREMENTS_TRACKER.md`:
1. `NOT_STARTED` (⏳): Partition is specified in `docs/shared/requirements/`, but no active development branch exists.
2. `IN_PROGRESS` (🔄): Feature branch `feature/<partition-id>` has been checked out, and active implementation or test work is underway.
3. `COMPLETED` (✅): All required subsystems (`DB`, `OAS`, `BE`, `MOB`, `ADM`) are implemented, 100% tests pass, OpenAPI bundle passes linting with 0 errors, and the branch is rebase-merged into `main`.

### Tracker Synchronization Invariants
- **Branch Checkout / Task Start:** Transition partition from `NOT_STARTED` to `IN_PROGRESS` in `REQUIREMENTS_TRACKER.md` and `CONTEXT.md`.
- **Subsystem Milestone:** Check off completed subsystems (`DB`, `BE`, `MOB`, `OAS`, `ADM`) in the tracker matrix.
- **Merge to Main:** Record the completion date, commit hash, verified test count, and transition status to `COMPLETED`.
- **Single Source of Truth:** `REQUIREMENTS_TRACKER.md` is the authoritative executive status matrix. Discrepancies between git history and `REQUIREMENTS_TRACKER.md` are treated as quality gate failures.