# Quality Gates & Automated CI/CD Enforcement Framework — Dink Machinery

**System:** Heavy Machinery & Logistics Rental Marketplace  
**Specification:** Formal Gate Definitions, Automated Verification Gates & Release Controls  
**Compliance Standards:** ISO/IEC/IEEE 29119, OWASP ASVS 4.0.3 Level 3, W3C WCAG 2.1/2.2 AA, ISO/IEC 25010  
**Subsystems Covered:** Go 1.27.3 Backend Modular Monolith, Flutter 3.44 Mobile Client, Next.js 16 Admin Web Console, PostgreSQL 18.4 + PostGIS 3.4, Valkey 8.0-alpine  
**Governing Documents:** `AGENTS.md`, `TECH_STACK.md`, `CODE_STANDARDS.md`, `API_CONVENTIONS.md`, `ERROR_CATALOG.md`, `DEFINITION_OF_DONE.md`, `QA_STRATEGY.md`

---

## 1. Executive Summary & Gate Philosophy

In **Dink Machinery**, a Quality Gate is an immutable, automated binary control point that halts software delivery immediately upon the detection of an invariant violation, security regression, contract divergence, or quality defect.

Because the platform orchestrates physical heavy equipment handovers and processes irreversible financial disbursements (Telebirr, CBEBirr, MPESA), manual discretion cannot override automated verification gates. Code that fails any gate cannot be merged to `main`, built into deployment containers, or promoted across environments.

### The Six Inviolable System Invariants
1. **Contract Supremacy:** `contracts/openapi.yaml` is the sole source of truth. Any schema mismatch, unlinted operation, or unapproved breaking change immediately trips the contract gate.
2. **Deterministic Minor-Unit Money:** All monetary values must be 64-bit signed integers (`int64` / `BIGINT`) in currency minor units (e.g., Santims, Cents). Any floating-point representation fails static and integration gates.
3. **In-Transaction Financial Atomicity (`BK-17-BE`):** State transitions and double-entry ledger balancing must execute inside a single atomic database transaction (`pgx.Tx`). Asynchronous or unbalanced ledger writes fail the integration gate.
4. **Optimistic Concurrency Control (OCC):** Mutable stateful records (`bookings`, `fleets`, `wallets`, `inspections`) must assert and increment `version INT NOT NULL DEFAULT 1`.
5. **Zero-PII & Secret Scrubbing (ASVS Level 3):** Phone numbers, National IDs, TINs, and bank accounts must be masked in logs, encrypted at rest (AES-256-GCM), and omitted from error envelopes. Gitleaks detects 0 secrets.
6. **The Continuous Ratchet:** Code coverage, latency budgets, accessibility compliance, and linting rules are strictly monotonic: quality metrics must never regress.

---

## 2. Multi-Stage Quality Gate Lifecycle

The quality lifecycle comprises seven sequential gates executed between the developer workstation and production release:

```
                              Dink Machinery Quality Gates Funnel
                                               │
  [Developer Workstation]                      ▼
                                     ┌───────────────────┐
                                     │  Gate G0: Local   │  Lefthook, Commitlint,
                                     │   Pre-Commit      │  CONTEXT.md Staged
                                     └─────────┬─────────┘
                                               │
  [Pull Request Opened]                        ▼
                                     ┌───────────────────┐
                                     │   Gate G1: Schema │  Spectral CLI, oasdiff,
                                     │   & Contract Drift│  Code-Gen Drift = 0
                                     └─────────┬─────────┘
                                               │
  [Automated PR Build]                         ▼
                                     ┌───────────────────┐
                                     │ Gate G2: Static & │  Go Unit (-race), BLoC,
                                     │    Unit Suites    │  Coverage Ratchet (>=80%/95%)
                                     └─────────┬─────────┘
                                               │
  [Testcontainer Integration]                  ▼
                                     ┌───────────────────┐
                                     │ Gate G3: Provider │  Testcontainers PG+Valkey,
                                     │   Wire & Ledger   │  kin-openapi Conformance
                                     └─────────┬─────────┘
                                               │
  [Automated Security Audit]                   ▼
                                     ┌───────────────────┐
                                     │  Gate G4: ASVS L3 │  Gitleaks, Semgrep,
                                     │   Security & SAST │  Zero-PII Regex, Trivy
                                     └─────────┬─────────┘
                                               │
  [PR Approval & Merge to Main]                ▼
                                     ┌───────────────────┐
                                     │ Gate G5: Non-Func │  k6 Performance (P95<80ms),
                                     │   A11y & Load QA  │  Axe-core WCAG 2.1 AA
                                     └─────────┬─────────┘
                                               │
  [Staging Release Candidate]                  ▼
                                     ┌───────────────────┐
                                     │ Gate G6: Staging  │  Patrol Mobile E2E,
                                     │  Acceptance & CD  │  Playwright Admin E2E
                                     └───────────────────┘
```

---

## 3. Detailed Gate Specifications

---

### Gate G0: Local Pre-Commit & Working Tree Governance
- **Execution Environment:** Local workstation (enforced via Git hook runner `lefthook`).
- **Trigger:** Invocation of `git commit`.
- **Blocking Authority:** Local Git process; aborts commit creation on failure.

#### Automated Checkpoints & Invariant Assertions
1. **Contract Schema Linting:**
   ```bash
   npx @stoplight/spectral-cli lint contracts/openapi.yaml --ruleset .spectral.yaml
   ```
   *Threshold:* 0 errors, 0 warnings.
2. **Architectural Decision Template Audit:**
   ```bash
   bash scripts/validate-decision-template.sh
   ```
   *Threshold:* All L1 ADRs contain 8 mandatory sections; all L2 Partition Decisions contain 7 sections; zero raw HTML placeholders.
3. **Bidirectional Decision Link Audit:**
   ```bash
   bash scripts/validate-decision-links.sh
   ```
   *Threshold:* Zero dangling ADR references; reciprocal linkage between requirement partitions and decision records verified.
4. **Conventional Commit Message Linting:**
   ```bash
   npx commitlint --config .github/commitlint.config.js --edit {1}
   ```
   *Threshold:* Must match `^(\[([A-Z0-9-]+)\]\s)?(\w+)(?:\(([^)]+)\))?:\s(.+)$`.
5. **The Atomic Context Checkpoint Invariant:**
   ```bash
   git diff --cached --name-only | grep -q "^CONTEXT.md$" || (echo "ERROR: CONTEXT.md must be updated and staged in every commit." && exit 1)
   ```
   *Threshold:* `CONTEXT.md` must be modified and staged in the commit.

---

### Gate G1: Contract Conformance & Zero-Drift Gate
- **Execution Environment:** GitHub Actions runner (`.github/workflows/ci-contract.yml`).
- **Trigger:** Pull Request opened, synchronized, or targeted against `main`.
- **Blocking Authority:** GitHub Branch Protection Rule; blocks PR merge.

#### Automated Checkpoints & Invariant Assertions
1. **OpenAPI 3.1.0 Spectral Validation:** Full scan across all endpoints, schemas, parameters, and RFC 9457 error models.
2. **Breaking API Change Detection (`oasdiff`):**
   ```bash
   oasdiff breaking origin/main:contracts/openapi.yaml contracts/openapi.yaml --fail-on-diff --format text
   ```
   *Threshold:* Zero breaking changes. Additions must be backward compatible. Breaking changes require an approved L1 ADR and a major version increment (`/api/v2`).
3. **Multi-Target Code-Generation Drift Audit:**
   ```bash
   bash scripts/generate-contracts.sh
   git diff --exit-code backend/internal/transport/http/generated/ admin/src/types/ mobile/packages/
   ```
   *Threshold:* Zero uncommitted drift between `contracts/openapi.yaml` and compiled client/server stubs.

---

### Gate G2: Static Analysis, Unit Testing & Coverage Ratchet
- **Execution Environment:** GitHub Actions runner (`.github/workflows/ci-backend.yml`, `ci-mobile.yml`, `ci-admin.yml`).
- **Trigger:** Pull Request updates; runs in parallel across subsystems.
- **Blocking Authority:** GitHub Branch Protection Rule; blocks PR merge.

#### Automated Checkpoints & Invariant Assertions
1. **Go 1.27.3 Backend Suite:**
   - Static analysis: `golangci-lint run ./...` (0 warnings).
   - Domain unit suite: `go test -v -race -coverprofile=coverage.out ./internal/modules/...`.
   - Thread safety: Zero race conditions detected by `-race`.
2. **Flutter 3.44 Mobile Suite:**
   - Static analyzer: `flutter analyze` (0 issues).
   - BLoC unit tests: `flutter test --coverage` (100% pass rate).
   - Touch target and semantic widget assertions (`androidTapTargetGuideline` $\ge 48\text{ dp}$).
3. **Next.js 16 Admin Web Console Suite:**
   - Type-checking & lint: `npm run lint` and `npm run build` (0 errors).
   - Unit & component tests: `vitest run --coverage` (100% pass rate).
4. **The Non-Decreasing Coverage Ratchet:**
   - Minimum total line coverage: $\ge 80.0\%$.
   - Financial ledger (`modules/payment`), booking state (`modules/booking`), and authentication (`modules/auth`) coverage: $\ge 95.0\%$.
   - $\Delta\text{Coverage} \ge 0.00\%$ against `origin/main`. PRs decreasing total coverage are automatically rejected.

---

### Gate G3: Integration Persistence, Wire Conformance & Ledger Atomicity
- **Execution Environment:** Containerized GitHub Actions runner with Docker Engine.
- **Trigger:** Successful completion of Gate G2.
- **Blocking Authority:** GitHub Branch Protection Rule; blocks PR merge.

#### Automated Checkpoints & Invariant Assertions
1. **Real Infrastructure Verification (`testcontainers-go`):**
   - Tests execute against containerized PostgreSQL 18.4 + PostGIS 3.4 and Valkey 8.0-alpine.
   - Migrations from `migrations/` applied sequentially (`up` and `down`).
2. **In-Transaction Ledger Atomicity Invariant (`BK-17-BE`):**
   - Every booking state change that mutates funds asserts balanced debit and credit entries in `financial_ledger_entries` summing to 0 minor units.
   - Assert database rollback occurs cleanly if any step inside `pgx.Tx` fails.
3. **Optimistic Concurrency Control Invariant:**
   - Competing transactions updating the same entity version must reject the second update with HTTP 409 (`BOOKING_CONCURRENCY_CONFLICT`).
4. **Provider Wire Schema Conformance (`kin-openapi`):**
   - Live HTTP responses from Go Chi handlers are validated against `contracts/openapi.yaml`.
   - Assert all 2xx responses match `SuccessEnvelope` and all 4xx/5xx responses match `ErrorEnvelope` with RFC 9457 attributes.
5. **Soft-Delete Invariant (`SD-01-BE`):**
   - Domain queries assert `WHERE deleted_at IS NULL` or query active views (`v_active_*`).
   - Direct hard `DELETE` statements trigger integration test failure.

---

### Gate G4: OWASP ASVS Level 3 Security, Vulnerability & Zero-PII Audit
- **Execution Environment:** Security pipeline runner (`.github/workflows/security.yml`).
- **Trigger:** Pull Request updates; nightly scheduled runs.
- **Blocking Authority:** Security team block; halts deployment pipeline.

#### Automated Checkpoints & Invariant Assertions
1. **Secret Scanning (Gitleaks):**
   ```bash
   gitleaks detect --source . --verbose
   ```
   *Threshold:* 0 secrets, tokens, private keys, or credentials detected in history or diffs.
2. **Zero-PII Logging Invariant:**
   ```bash
   grep -rnEi 'logger\.(Info|Debug|Warn|Error)\(.*(password|pin|tin|national_id|account_number)' backend/ admin/ mobile/
   ```
   *Threshold:* 0 occurrences. Sensitive fields must pass through masking formatters (`+251911****44`, `TIN: *****1234`).
3. **Static Application Security Testing (Semgrep SAST):**
   ```bash
   semgrep --config p/owasp-top-ten --config p/golang --error backend/ admin/
   ```
   *Threshold:* 0 high or critical security findings.
4. **Dependency Vulnerability Scanning:**
   - Go: `govulncheck ./...` (0 reachable vulnerabilities).
   - Node: `npm audit --audit-level=high` (0 high/critical CVEs).
5. **Container Base Image Hardening (Trivy):**
   ```bash
   trivy image --severity HIGH,CRITICAL dink-backend:latest
   ```
   *Threshold:* 0 unfixed high or critical OS-level CVEs.

---

### Gate G5: Performance Budgets, Accessibility & Load Verification
- **Execution Environment:** Performance & QA testbed (Staging pre-deployment).
- **Trigger:** Merge to `main` branch or release tag creation.
- **Blocking Authority:** Automated pipeline gate; blocks artifact release.

#### Automated Checkpoints & Invariant Assertions
1. **Automated Geospatial Load Test (k6):**
   - Scenario: 500 concurrent VUs executing radius searches centered on Addis Ababa corridor.
   - P95 latency: $\le 80\text{ ms}$.
   - P99 latency: $\le 150\text{ ms}$.
   - HTTP failure rate: $< 0.1\%$.
2. **Booking Hold Concurrency Spike (k6):**
   - Scenario: 100 concurrent requests competing for an identical machinery asset.
   - Assert: Exactly 1 hold allocated; 99 requests cleanly rejected with HTTP 409 (`BOOKING_HOLD_CONFLICT`). Zero deadlocks.
3. **Web Accessibility Scan (`@axe-core/playwright`):**
   - Full automated scan of all Next.js admin routes.
   - *Threshold:* 0 WCAG 2.1 AA violations.
4. **Mobile Rendering Budget:**
   - Flutter integration driver asserts 60 FPS ($< 16.6\text{ ms}$ per frame) during catalog fling scrolling on emulated devices.

---

### Gate G6: Staging Acceptance, End-to-End Journeys & CD Promotion
- **Execution Environment:** Staging environment (mirror of production infrastructure).
- **Trigger:** Artifact deployment to staging.
- **Blocking Authority:** Release management gate; blocks production promotion.

#### Automated Checkpoints & Invariant Assertions
1. **Mobile End-to-End User Journeys (Patrol):**
   - Journey 1: Renter searches equipment by GPS radius $\rightarrow$ acquires 60-minute payment hold $\rightarrow$ confirms payment via mock Telebirr webhook $\rightarrow$ booking transitions to `CONFIRMED`.
   - Journey 2: Inspector captures equipment check-in photos $\rightarrow$ submits hour meter reading $\rightarrow$ asserts ClamAV clean scan $\rightarrow$ transitions to `ACTIVE_OPERATION`.
2. **Admin Console End-to-End Journeys (Playwright):**
   - Journey 1: Admin logs in with MFA $\rightarrow$ navigates dynamic table $\rightarrow$ reviews pending heavy equipment operator certificate $\rightarrow$ approves verification status.
   - Journey 2: Admin inspects dispute fuel discrepancy $\rightarrow$ executes partial deposit forfeiture $\rightarrow$ verifies immutable ledger debit/credit balance.
3. **Database Migration Verification:**
   - Production readiness drill: Run migration `up`, assert schema integrity, run migration `down`, re-run `up`.
4. **Smoke Test Health Check:**
   - Ingress gateway responds with HTTP 200 on `/healthz` adhering to platform `SuccessEnvelope`.

---

## 4. Master Quality Gate Summary Matrix

| Gate | Stage | Scope | Key Tools | Invariant / Failure Condition |
|---|---|---|---|---|
| **G0** | Pre-Commit | Workstation | Lefthook, Spectral, Scripts | Broken links, unformatted commit, missing `CONTEXT.md`. |
| **G1** | Contract | CI Pull Request | Spectral CLI, `oasdiff` | Unlinted schema, breaking change, uncommitted code-gen drift. |
| **G2** | Static & Unit | CI Pull Request | `golangci-lint`, `flutter test`, Vitest | Linter warning, test failure, coverage drops ($< 80\%/95\%$). |
| **G3** | Integration | CI Pull Request | Testcontainers, `kin-openapi` | DB failure, OCC failure, unbalanced ledger, schema drift. |
| **G4** | Security | CI Pull Request | Gitleaks, Semgrep, Trivy | Leaked secret, unmasked PII log, reachable CVE, ASVS L3 fail. |
| **G5** | Performance | CI Staging Pre-Release | k6, `@axe-core/playwright` | P95 $> 80\text{ ms}$, booking deadlock, WCAG AA violation. |
| **G6** | Staging E2E | Staging Acceptance | Patrol, Playwright, Traefik | Critical user journey failure, broken rollback migration. |

---

## 5. The Emergency Break-Glass Override Protocol

Under strictly audited conditions, an emergency bypass of non-critical gates may occur to deploy a hotfix for an active production severity P0 incident.

### 5.1. Non-Bypassable Invariants (Absolute Zero-Exemption)
The following checks can **NEVER** be bypassed under any circumstances, even during emergency production remediation:
- **Gitleaks Secret Detection:** Secrets in code are permanently non-releasable.
- **Ledger In-Transaction Atomicity (`BK-17-BE`):** Unbalanced ledger code will not be deployed.
- **Integer Minor-Unit Currency Math:** Float-based money logic will not be deployed.
- **Parameterized SQL Execution:** Dynamic query concatenations will not be deployed.

### 5.2. Break-Glass Procedure
1. **Invocation:** The incident commander opens an emergency hotfix branch (`hotfix/<incident-id>`).
2. **Dual Sign-Off:** Bypass of transient gates (e.g., Gate G5 full load test or coverage ratchet delta) requires simultaneous cryptographic approval from:
   - Lead Platform Architect
   - Backend Subsystem Lead or Security Officer
3. **Commit Annotation:** The emergency commit must include the break-glass trailer:
   ```text
   [HOTFIX-P0] fix(payment): resolve telebirr webhook timeout parsing

   Emergency patch bypassing G5 load suite under INCIDENT-2026-1026-01.
   
   Break-Glass-Approved-By: @architect-lead, @security-lead
   Incident-Ref: #INC-0104
   Docs: CONTEXT.md
   ```
4. **Retroactive Remediation SLA:** Within 24 hours of incident resolution, the hotfix branch must be normalized with complete test coverage, passing all skipped gates retroactively.

---

## 6. Continuous Integration Workflow Blueprint

The continuous integration architecture executes gates concurrently where possible, failing fast to minimize feedback latency.

```yaml
# .github/workflows/quality-gates.yml
name: Dink Machinery Master Quality Gates

on:
  pull_request:
    branches: [main]
  push:
    branches: [main]

concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

jobs:
  gate-g1-contracts:
    name: "Gate G1: API Contracts & Drift"
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: actions/setup-node@v4
        with:
          node-version: 22
      - run: npm install -g @stoplight/spectral-cli
      - run: npx @stoplight/spectral-cli lint contracts/openapi.yaml --ruleset .spectral.yaml
      - name: Audit Breaking Changes
        run: |
          curl -fsSL [https://raw.githubusercontent.com/tufin/oasdiff/main/install.sh](https://raw.githubusercontent.com/tufin/oasdiff/main/install.sh) | sh
          ./bin/oasdiff breaking origin/main:contracts/openapi.yaml contracts/openapi.yaml --fail-on-diff
      - name: Verify Zero Code-Gen Drift
        run: |
          bash scripts/generate-contracts.sh
          git diff --exit-code

  gate-g4-security:
    name: "Gate G4: ASVS L3 & Secret Audit"
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - name: Run Gitleaks
        uses: gitleaks/gitleaks-action@v2
      - name: Zero-PII Regex Audit
        run: |
          if grep -rnEi 'logger\.(Info|Debug|Warn|Error)\(.*(password|pin|tin|national_id|account_number)' backend/ admin/ mobile/; then
            echo "CRITICAL: Unmasked PII logging patterns detected!"
            exit 1
          fi
      - name: Run Semgrep SAST
        run: npx semgrep --config p/owasp-top-ten --error backend/ admin/

  gate-g2-unit-backend:
    name: "Gate G2: Go Backend Unit & Race"
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.27.3'
      - name: Golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: v1.64
          working-directory: backend
      - name: Run Unit Tests with Race Detection
        run: |
          cd backend
          go test -v -race -coverprofile=coverage.out ./internal/modules/...
      - name: Verify Financial Coverage Ratchet (>= 95%)
        run: |
          cd backend
          FIN_COV=$(go tool cover -func=coverage.out | grep "internal/modules/payment" | awk '{print $3}' | tr -d '%')
          echo "Payment/Ledger Coverage: ${FIN_COV}%"
          awk -v cov="$FIN_COV" 'BEGIN { if (cov < 95.0) exit 1 }'

  gate-g2-unit-mobile:
    name: "Gate G2: Flutter Mobile Suite"
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: subosito/flutter-action@v2
        with:
          flutter-version: '3.44.x'
          cache: true
      - run: cd mobile && flutter analyze
      - run: cd mobile && flutter test --coverage

  gate-g2-unit-admin:
    name: "Gate G2: Next.js Admin Suite"
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 22
      - run: cd admin && npm ci
      - run: cd admin && npm run lint
      - run: cd admin && npm run build

  gate-g3-integration:
    name: "Gate G3: Integration, Wire & Ledger"
    needs: [gate-g1-contracts, gate-g2-unit-backend]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.27.3'
      - name: Run Provider Wire & Testcontainer Integration Tests
        run: |
          cd backend
          go test -v -race ./internal/transport/http/... ./internal/platform/database/...

  gate-clearance:
    name: "Master Quality Gate Sign-Off"
    needs: [gate-g1-contracts, gate-g4-security, gate-g2-unit-backend, gate-g2-unit-mobile, gate-g2-unit-admin, gate-g3-integration]
    runs-on: ubuntu-latest
    steps:
      - run: echo "SUCCESS: All PR Quality Gates cleared for merge."
```

---

## 7. Quality Gate Audit Checklist

Before submitting or reviewing a pull request, affirm compliance with every gate checkpoint:

- [ ] **Gate G0 Cleared:** Local `lefthook` hooks executed; Spectral lint clean; `CONTEXT.md` updated and staged in this commit.
- [ ] **Gate G1 Cleared:** `openapi.yaml` passes Spectral; `oasdiff` reports 0 breaking changes; code-gen drift is 0.
- [ ] **Gate G2 Cleared:** All unit tests pass with `-race`; coverage satisfies thresholds ($\ge 80\%$ general, $\ge 95\%$ financial/auth); coverage did not decrease.
- [ ] **Gate G3 Cleared:** Integration tests against PostgreSQL 18.4 + PostGIS and Valkey pass; ledger balance invariant asserted inside single `pgx.Tx` transactions; wire responses conform to OpenAPI schemas.
- [ ] **Gate G4 Cleared:** Gitleaks reports 0 secrets; static regex confirms 0 unmasked PII logs; Semgrep reports 0 high/critical vulnerabilities.
- [ ] **Gate G5 Prepared:** Load tests and accessibility checks verified against latency ($< 80\text{ ms}$ P95) and WCAG 2.1 AA benchmarks.
- [ ] **Gate G6 Prepared:** Staging acceptance test plans documented for Patrol and Playwright automated journeys.