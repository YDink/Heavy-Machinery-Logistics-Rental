# Task 00a: Phase 0 Kickoff — Architecture, Canonical Contracts & Governance Baseline

**Task Identifier:** `TASK-00a-PHASE-0-KICKOFF`  
**Phase:** 0 (Contract Definition, Documentation & Architectural Governance)  
**Target Completion State:** Readiness to execute `BOOTSTRAP.md` (Phase 1)  
**Governing Protocols:** `AGENTS.md`, `CONTEXT.md`, `TECH_STACK.md`, `CODE_STANDARDS.md`

---

## 1. Executive Summary & Objective

The objective of Task 00a is to establish the complete architectural foundation, contract specifications, decision tracking mechanisms, and automated governance tooling for **Dink Machinery**. 

Per platform governance invariants, no operational application code (Go modules, Flutter widgets, Next.js views) may be generated before all Phase 0 artifacts are authored, strictly validated, and checked into version control. Phase 0 establishes the single source of truth across the modular monolith backend, mobile app, and web admin console.

---

## 2. Invariant Governance Gates

Before completing this task and transitioning to Phase 1, the following platform invariants must be verified:

1. **Contract-First Precedence:** All API endpoints, request schemas, response models, and error structures must be declared in `contracts/openapi.yaml` conforming to OpenAPI 3.1.0 and Spectral linting rules.
2. **Standardized Response Envelopes:** No raw JSON entities or naked arrays. All endpoints must implement either the Success Envelope or the RFC 9457 Problem Details Error Envelope.
3. **Integer Minor-Unit Currency:** All financial fields must be 64-bit signed integers coupled with ISO 4217 currency identifiers. Floating-point numbers are prohibited.
4. **Architectural Decision Tracking:** Any design choice satisfying the 5-point qualification test must be recorded either as a system-wide L1 ADR (`docs/decisions/ADR-NNN-<slug>.md`) or an L2 Partition Decision (`docs/decisions/partitions/<partition-id>.md`).
5. **Bidirectional Decision Linkage:** All decision records must have reciprocal references in corresponding requirement partitions.
6. **Zero-PII Assurance (OWASP ASVS L3):** Telephone numbers, National IDs, TINs, and financial accounts must be masked in logging patterns and omitted from error envelopes.

---

## 3. Scope & Phase 0 Deliverables Matrix

| Artifact Path | Category | Status | Verification Protocol |
|---|---|---|---|
| `.env.dev` | Configuration | Complete | Git tracking verification (zero secrets) |
| `.editorconfig` | Tooling | Complete | Cross-IDE whitespace and indent validation |
| `.gemini/settings.json` | Tooling | Complete | Context file and ignore pattern audit |
| `.github/workflows/security.yml` | CI/CD | Complete | Gitleaks, Semgrep, Govulncheck, Trivy checks |
| `.github/workflows/ci-docs.yml` | CI/CD | Complete | Markdownlint, Spectral, Decision link audit |
| `.github/commitlint.config.js` | Governance | Complete | Conventional commit regex with ticket headers |
| `.github/pull_request_template.md` | Governance | Complete | Checklists for OCC, RFC 9457, ASVS L3, tests |
| `contracts/README.md` | Contracts | Complete | Schema conventions and code-gen mapping guide |
| `contracts/openapi.yaml` | Contracts | Complete | Spectral CLI linting (`.spectral.yaml`) |
| `scripts/scaffold-decisions.sh` | Scripts | Complete | Execution test for L1 ADR & L2 generation |
| `scripts/validate-decision-links.sh` | Scripts | Complete | Zero broken decision or partition references |
| `scripts/validate-decision-template.sh` | Scripts | Complete | 8 mandatory ADR sections & 7 L2 sections |
| `docs/TECH_STACK.md` | Architecture | Complete | Version pinning and library registry audit |
| `OPEN_QUESTIONS.md` | Architecture | Complete | Blocking status and partition impact tracking |
| `ERROR_CATALOG.md` | Specification | Complete | RFC 9457 code mapping, HTTP status, retryability |
| `CODE_STANDARDS.md` | Guidelines | Complete | Invariant rules for Go, Flutter, Next.js, SQL |
| `BOOTSTRAP.md` | Runbook | Complete | Step-by-step Phase 1 execution commands |
| `API_CONVENTIONS.md` | Specification | Complete | REST design guide, cursor pagination, OCC, RPC |
| `docs/quality/BOOTSTRAP_CHECKLIST.md`| Quality | Pending Audit | Full Phase 0 sign-off checklist |
| `CONTEXT.md` | Meta | Active | Atomic update per commit |

---

## 4. Execution & Verification Steps

### Step 1: Contract Validation & Schema Check
Run Spectral CLI across the canonical OpenAPI specification to verify full compliance with RFC 9457 and integer currency rules:
```bash
npx @stoplight/spectral-cli lint contracts/openapi.yaml --ruleset .spectral.yaml
```
- **Expected Outcome:** 0 errors, 0 warnings.
- **Remediation:** If missing `operationId`, unreferenced schemas, or floating-point monetary properties are flagged, correct `contracts/openapi.yaml` immediately.

### Step 2: Governance Script Conformance Drill
Execute the automated decision validation test suite to confirm template compliance and reference integrity:
```bash
# Make scripts executable
chmod +x scripts/*.sh

# Run structural template audit
bash scripts/validate-decision-template.sh

# Run bidirectional reference audit
bash scripts/validate-decision-links.sh
```
- **Expected Outcome:**
  ```text
  SUCCESS: All architectural decision templates strictly adhere to governance specifications.
  SUCCESS: All architectural decision references and bidirectional links are valid.
  ```

### Step 3: Drift & Source-of-Truth Synchronization Audit
Ensure no documentation drift exists between `docs/source/` and partitioned specifications:
```bash
if [ -f "scripts/check-docs-drift.sh" ]; then
  bash scripts/check-docs-drift.sh
fi
```

### Step 4: Verify Secret & Static Governance Linters
Validate that commit lint rules and hook configurations are operational:
```bash
# Verify commit message format
echo "[BOOTSTRAP] feat(platform): test commit message" | npx commitlint --config .github/commitlint.config.js
```

---

## 5. Phase 0 Quality Gate Checklist (Sign-Off Criteria)

Before declaring Phase 0 complete and advancing to Phase 1:

- [ ] **Contract Completeness:** Core bounded contexts (`auth`, `machinery`, `booking`, `payment`, `geo`, `admin`) defined in `contracts/openapi.yaml`.
- [ ] **Error Model Alignment:** All error responses map to entries registered in `ERROR_CATALOG.md`.
- [ ] **Multi-Country Parameterization:** `X-Country-Code` and `country_configs` mechanisms documented in `API_CONVENTIONS.md` and `contracts/openapi.yaml`.
- [ ] **Financial Double-Entry Ledger Invariants:** Explicitly defined in `CODE_STANDARDS.md § 1.4` and `contracts/README.md`.
- [ ] **Automated Tooling Verified:** Validation scripts executed and passing locally.
- [ ] **Zero Open Phase 0 Blockers:** All blocking open questions in `OPEN_QUESTIONS.md` resolved or scheduled for subsequent feature partitions.
- [ ] **Bootstrap Runbook Ready:** `BOOTSTRAP.md` verified against workstation tool versions (Go `1.27.3`, Flutter `3.44`, Next.js `16.x`, PostgreSQL `18.4`, Valkey `8.0`).

---

## 6. Handover to Phase 1 Bootstrap

Upon satisfying all items in Section 5:

1. **Commit Phase 0 Baseline:**
   ```bash
   git add .
   git commit -m "[PHASE-0] docs(governance): complete phase 0 contracts and architectural baseline

   Established canonical OpenAPI 3.1.0 specification, decision validation scripts,
   error catalog, code standards, tech stack specifications, and Phase 1 bootstrap runbook.

   Tests: Spectral passed (0 errors), decision-links passed, commitlint verified.
   Docs: CONTEXT.md, BOOTSTRAP.md, API_CONVENTIONS.md, ERROR_CATALOG.md
   Refs: #0000"
   ```

2. **Trigger Phase 1 Execution:**
   Proceed to execute `BOOTSTRAP.md` following `docs/AI_TASKS/01-bootstrap-subsystems.md`.