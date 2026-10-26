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

### Coverage Ratchet
- **Target:** Minimum 80% line coverage on changed code; 95% branch coverage required for financial calculations, ledger operations, auth paths, and payout workflows.
- **Ratchet Rule:** Code coverage metrics are non-decreasing. PRs reducing overall test coverage will fail CI.

### Quality Gate Hierarchy (G0 to G4)
- **G0 (Pre-Commit / Local):** Linter execution, static analysis, unit test suite, and atomic doc-commit verification (`lefthook.yml`).
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
- **No Unapproved Paradigms:** Do not introduce architectural patterns, global state solutions, or external dependencies without an explicit L1 (ADR) or L2 (Partition) decision doc.
- **DTO Isolation:** Database models and persistence entities must never be exposed directly to the controller layer or client responses. DTOs and view models are mandatory.
- **Layered Validation:** Validate format, presence, and boundary constraints at the edge/controller layer. Validate business invariants and relational constraints in the domain/service layer.

### Go Backend (Modular Monolith)
- Package dependencies point inward: `domain` → `application` → `infrastructure`/`presentation`.
- The domain layer contains pure business entities, value objects, and repository interfaces. No SQL queries or external driver imports within domain packages.
- Mutations use command handlers returning `Result[T, E]` or error types. Queries use specialized query handlers returning read-optimized DTOs.
- `sqlc` is the standard for SQL generation. Raw SQL strings within application logic are prohibited.

### Flutter Mobile
- State management relies strictly on the BLoC pattern (`flutter_bloc`). UI widgets emit events; state changes are received via typed states.
- Clean separation: Presentation (Widgets/BLoC) → Domain (Use Cases/Entities) → Data (Data Sources/Repositories).
- Network calls are routed through a configured `Dio` client with logging, auth refresh, and error-to-problem-details interceptors.

### Next.js Admin Panel
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
Docs: <list of documentation files updated, e.g., CONTEXT.md, decisions/partitions/05-05.md>
Refs: #<issue-or-task-number>
```

### The Atomic Doc-Commit Rule
- Code and tests must never be committed without staging their operational documentation simultaneously.
- **Mandatory in every commit:** Code changes + Tests + `CONTEXT.md`.
- **Conditional:** If a design decision meets any qualification test, the relevant partition decision doc or ADR must be updated and staged within the exact same commit.

### Pre-Commit Gate Enforcement (`lefthook.yml`)
Local git hooks block commits failing the following:
1. `G0` linter, code style, and test runs.
2. Verification that `CONTEXT.md` is present in the staged index.
3. Verification that commit messages include the `Docs:` footer.
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
- Once created, runtime versions and package trees are pinned via lockfiles (`go.sum`, `pubspec.lock`, `package-lock.json`, and `.tool-versions`). Any subsequent major version bump requires an L1 ADR.

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
- Load only the baseline operating files: `AGENTS.md`, `CONTEXT.md`, and `contracts/openapi.yaml`.
- Pull specific partition files on demand based on the active task:
  - Requirements: `docs/shared/requirements/<partition-id>.md`
  - Design/Flows: `docs/shared/flows/<flow-id>.md`
  - Decisions: `docs/decisions/partitions/<partition-id>.md`
- Always verify open blockers in `docs/OPEN_QUESTIONS.md` before executing new architectural implementations.