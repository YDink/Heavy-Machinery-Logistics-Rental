# Session Startup & Context Orientation Protocol — Dink Machinery

**Protocol ID:** `PROTO-01-SESSION-START`  
**Target:** AI Agents & Engineering Workstations  
**Scope:** Deterministic initialization at the beginning of any engineering or AI development turn.  
**Governing Documents:** `AGENTS.md`, `CONTEXT.md`, `OPEN_QUESTIONS.md`, `CODE_STANDARDS.md`

---

## 1. Objective & Invariant Mindset

Every AI agent or engineer beginning a session on **Dink Machinery** must run through this startup sequence before generating code, executing migrations, or modifying requirements. 

This protocol guarantees:
1. **Zero Context Drift:** The working agent anchors on `CONTEXT.md` as the single source of operational truth.
2. **Branch & Partition Alignment:** Operations occur strictly within the designated feature branch (`feature/<partition-id>`).
3. **Pre-Flight Invariant Assurance:** Contract-first discipline, RFC 9457 error models, ASVS Level 3 Zero-PII rules, and integer minor-unit money invariants are actively loaded into working memory.

---

## 2. Sequential Startup Procedure

```text
[Step 1: Read State] ---> [Step 2: Inspect Branch] ---> [Step 3: Audit Blockers] ---> [Step 4: Verify Tools] ---> [Step 5: Align & Execute]
```

### Step 1: Ingest Context & Progress Snapshot
Read and parse the root tracking document `CONTEXT.md`:
- Locate the **Current Progress Snapshot**.
- Identify the **Active Partition** and **Subsystem Focus** (`backend`, `mobile`, `admin`, `contracts`, `infra`).
- Note the **Next Immediate Tasks** listed in `CONTEXT.md`.
- Review recently approved ADRs or Partition Decisions in `docs/decisions/`.

```bash
# Display immediate status and focus areas
head -n 45 CONTEXT.md
```

### Step 2: Validate Git Branch & Working Tree
Confirm that the working environment matches the active partition recorded in `CONTEXT.md`:
- Working on `main` directly for feature development is strictly prohibited.
- Working tree must be clean or have only known staged edits.

```bash
git status --short
git branch --show-current
```
*If on `main` and initiating a new partition, create the partition branch:*
```bash
# Format: feature/<partition-id>-<slug>
git checkout -b feature/05-01-core-identity-auth
```

### Step 3: Check Open Questions & Blockers
Cross-reference the active requirement partition against `OPEN_QUESTIONS.md`:
- If an open question blocks the current partition, do not implement speculative or throwaway logic.
- Confirm whether required ADRs (`docs/decisions/ADR-NNN-*.md`) or Partition Decisions (`docs/decisions/partitions/<partition-id>.md`) have been finalized.

### Step 4: Toolchain & Local Infrastructure Sanity Check
Confirm that required local engines and compilers are online and pinned:

```bash
# Check compiler versions
go version         # Must be Go 1.27.3
flutter --version  # Must be Flutter 3.47.9
node -v            # Must be Node v26.13.0

# Check container infrastructure (if developing backend or integration tests)
docker compose ps
```

### Step 5: Fast Invariant & Pre-Commit Lint Drill
Verify that baseline contracts and linters currently pass:

```bash
# 1. Lint OpenAPI contract
pnpm dlx @stoplight/spectral-cli lint contracts/openapi.yaml --ruleset .spectral.yaml

# 2. Run decision template & reference validation
bash scripts/validate-decision-template.sh
bash scripts/validate-decision-links.sh
```

---

## 3. Working Memory Invariant Checklist

Before beginning code edits in this session, affirm adherence to platform invariants:

- [ ] **Contract-First:** No API endpoint or JSON structure may be implemented unless defined in `contracts/openapi.yaml`.
- [ ] **Response Envelopes:** All responses wrapped in `success: true` or RFC 9457 Problem Details `success: false`.
- [ ] **Integer Minor-Unit Money:** Currency fields stored as `int64` minor units coupled with ISO 4217 code. Zero floating-point money.
- [ ] **Multi-Country Parameterization (`LOC-01-BE`):** Country rules resolved dynamically via `country_configs`; no hardcoded constants.
- [ ] **Zero-PII Exposure (ASVS L3):** Telephone numbers, National IDs, TINs, and bank accounts masked in structured logs and omitted from errors.
- [ ] **Optimistic Locking:** Stateful entity mutations assert `WHERE id = $1 AND version = $2` and increment `version`.
- [ ] **Soft Deletes (`SD-01-BE`):** Queries filter `deleted_at IS NULL` or target `v_active_*` views.
- [ ] **Atomic Doc-Commit Rule:** `CONTEXT.md` must be updated and staged within every commit that advances task state.

---

## 4. Session Kickoff Status Template

When starting an agent turn, output the following structured brief before taking action:

```markdown
### 🚀 Session Initialized

- **Active Partition:** `[Partition ID & Name]`
- **Current Branch:** `[Branch Name]`
- **Target Subsystem:** `[Go Backend | Flutter Mobile | Next.js Admin | Contracts]`
- **Immediate Task ID:** `[TASK-ID / REQ-ID]`
- **Relevant Decisions / ADRs:** `[ADR-NNN / D-XX-YY-ZZZ]`
- **Blocking Questions:** `[None / OQ-NNN]`
- **Action Plan:** `[1-3 sentence summary of exact steps about to be executed]`
```