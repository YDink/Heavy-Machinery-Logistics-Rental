# Manual Testing & Pre-Commit Verification Protocol — Dink Machinery

**Document Identifier:** `QA-MANUAL-VERIFICATION-PROTOCOL-V1.0`  
**Classification:** Canonical Quality Assurance Standard & Inviolable Pre-Commit Gate  
**Governing Documents:** `AGENTS.md § 1 & § 4`, `docs/quality/QUALITY_GATES.md (Gate G0)`, `docs/quality/QA_STRATEGY.md`  
**Target Runtimes:** Go 1.27.3 Backend Monolith, Flutter 3.47.9 Mobile Client, Next.js 16.x Admin Console, PostgreSQL 18.6 + PostGIS 3.6.4  
**Status:** Approved Platform Invariant

---

## 1. Executive Mandate & Pre-Commit Condition

In **Dink Machinery**, automated testing (unit and integration tests) is a necessary baseline but is **not sufficient on its own** to authorize a code commit. Automated tests verify code against predefined assumptions; **manual testing verifies runtime behavior, usability, operational safety, and system resilience under realistic field conditions**.

### The Rule
> **No code or database migration may be committed to any branch without prior, documented manual testing across every affected subsystem.**
> 
> Every commit must explicitly record verified manual test scenarios in the commit body/footer under `Manual-Tests:`, and detailed execution evidence must be logged in `CONTEXT.md`.

---

## 2. Subsystem Manual Testing Runbooks

### 2.1. Database Subsystem Manual Verification
Every migration script (`migrations/NNNNNN_*.up.sql` and `.down.sql`) must undergo a live two-way execution drill before commit:

1. **Migration Up Verification:**
   - Execute the `.up.sql` migration against the target PostgreSQL 18 database.
   - Inspect the database schema: confirm tables, columns, indexes, GiST exclusion constraints, and initial triggers exist.
   - Verify default values: confirm OCC `version` defaults to `1`, `created_at` defaults to `NOW()`, and monetary columns are typed `BIGINT` (minor units).
2. **Rollback Drill (Migration Down):**
   - Execute the `.down.sql` migration against the database.
   - Confirm all newly created tables, enum types, and constraints are dropped cleanly with zero dangling locks or dependencies.
3. **Re-Application Verification:**
   - Re-run `.up.sql` to verify clean re-applicability (idempotence of deployment).
4. **Generated Query Validation:**
   - Run `sqlc generate` and verify that generated Go models and query functions compile without warnings.

---

### 2.2. Go Backend HTTP & Service Manual Verification
Every newly added or modified Go standard library handler (`net/http.ServeMux`) must be probed live via real HTTP requests (using `curl` or automated HTTP probe harnesses) before commit:

1. **Happy Path Probe (HTTP 200/201/202):**
   - Dispatch valid JSON payload with required headers (`Authorization: Bearer <token>`, `X-Country-Code: ET`, `Idempotency-Key: <uuid>`).
   - Inspect the response: verify standard success envelope:
     ```json
     {
       "success": true,
       "data": { ... },
       "error": null,
       "meta": {
         "version": "1.0",
         "requestId": "req_...",
         "timestamp": "..."
       }
     }
     ```
2. **Edge Case & Invariant Violation Probes (HTTP 400/409/422):**
   - **Meter Reading Regression:** Send a check-out meter reading below baseline check-in meter. Confirm immediate rejection with HTTP 422 `ERR_METER_READING_REGRESSION`.
   - **Temporal Conflict / Collision:** Send an extension request overlapping an existing booking or blackout. Confirm HTTP 409 `ERR_FLEET_UNAVAILABLE`.
   - **Invalid Status Machine State:** Dispatch action on terminal booking. Confirm HTTP 409 `BOOKING_INVALID_STATE_TRANSITION`.
   - **RFC 9457 Structure:** Confirm error response uses `application/problem+json` semantics wrapped in standard envelope:
     ```json
     {
       "success": false,
       "data": null,
       "error": {
         "type": "https://api.dinkmachinery.com/errors/...",
         "title": "...",
         "status": 422,
         "detail": "...",
         "code": "..."
       }
     }
     ```
3. **Authentication & RBAC Boundary Probes (HTTP 401/403):**
   - Probe endpoint without `Authorization` header: verify HTTP 401 `UNAUTHORIZED`.
   - Probe endpoint with a role lacking the requisite permission (e.g. contractor attempting owner action): verify HTTP 403 `FORBIDDEN`.
4. **Database Mutation Verification:**
   - Query the database following the HTTP call to verify that state transitioned, `version` incremented by 1, and audit log records were appended.

---

### 2.3. Flutter Mobile Client Manual Verification
Every newly added or modified mobile screen, modal sheet, or BLoC state flow must be manually verified:

1. **Visual Contrast & High-Visibility Theme:**
   - Verify screen in the dark, high-contrast operational palette (`#111827` background, `#F9FAFB` high-contrast typography, high-visibility state badges).
   - Verify outdoor readability under simulated sunlight glare (WCAG 2.2 AAA text contrast ratio $\ge 7:1$).
2. **Ergonomic Touch Targets:**
   - Confirm all primary action buttons, inputs, dropdowns, and date pickers have a minimum touch target size of **56dp** (`minHeight: 56.0`).
3. **Bilingual Localization Drill (English ↔ Amharic):**
   - Toggle language from English to Amharic (`locales: [Locale('en'), Locale('am')]`).
   - Verify 100% of user-facing strings (labels, hints, dialog titles, button texts, error messages, units) render properly in Amharic with appropriate Ge'ez typography line-height (`height: 1.45`).
   - Confirm **zero** untranslated fallback English strings or raw code identifiers.
4. **Interactive Form Validation & States:**
   - **Empty Form Submission:** Tap submit without inputs; verify inline field validation errors appear immediately.
   - **Invalid Range Inputs:** Enter invalid engine hours or past dates; verify descriptive error feedback.
   - **Loading State:** Verify submit button transitions to loading indicator (`CircularProgressIndicator`) and disables duplicate taps.
   - **Error / Offline State:** Simulate network error; verify descriptive banner or snackbar renders with retry CTA.

---

### 2.4. Next.js Admin Web Console Manual Verification
Every newly added or modified administrative view or server action must be manually verified:

1. **Server Component & Schema Rendering:**
   - Verify dynamic table loads data without client-side hydration mismatch.
   - Confirm Zod form validation halts invalid inputs before server action submission.
2. **Dispute / Action Overrides:**
   - Execute administrative override (e.g., dispute ruling); verify optimistic UI update and notification toast (`sonner`).

---

### 2.5. Multi-Machine Database Fixtures & State Parity Protocol (Windows ↔ macOS)
Because development and manual testing occur across multiple workstations (e.g. Windows at the office, macOS at home), **uncommitted, local-only database mutations are strictly prohibited**.

1. **Version-Controlled Fixtures:**
   - Every database record inserted or required for manual verification (test users, fleet machinery, geofences, active bookings, telemetry shifts, damage claims) must be scripted in `scripts/fixtures/<partition-id>_fixtures.sql`.
   - Ad-hoc database inserts that exist only in a single workstation's PostgreSQL instance violate Gate G0.
2. **Strict Idempotency:**
   - All fixture SQL statements must use `INSERT INTO ... ON CONFLICT (id) DO NOTHING` (or `DO UPDATE SET ...`).
   - All primary keys must use static, deterministic UUIDs (e.g. `11111111-1111-4111-8111-111111111111`). Dynamic `gen_random_uuid()` is prohibited in test fixtures to ensure exact ID parity across machines.
3. **Session Start Synchronization:**
   - At session start on any machine, running `scripts/fixtures/apply-fixtures.sh` (macOS/Linux) or `.\scripts\fixtures\apply-fixtures.ps1` (Windows) guarantees all test records are present in the local PostgreSQL instance without duplicate key collisions.
4. **Pre-Commit Verification:**
   - The corresponding `scripts/fixtures/<partition-id>_fixtures.sql` file must be staged and committed alongside code and documentation.

---

## 3. Commit Format & Evidence Requirements

### 3.1. Commit Message Standard
Commits must include the `Manual-Tests:` footer enumerating the verified manual scenarios:

```text
[<REQ-ID>] <type>(<scope>): <short summary>

<detailed description of change, rationale, and consequences>

Tests: <automated test suites executed and pass count>
Manual-Tests:
- <subsystem>: <scenario 1 description> -> <verified outcome>
- <subsystem>: <scenario 2 description> -> <verified outcome>
Docs: CONTEXT.md, REQUIREMENTS_TRACKER.md, <other docs>
Refs: #<issue-or-task-number>
```

### 3.2. Logging in `CONTEXT.md`
For every partition in progress, `CONTEXT.md` must maintain an itemized `Manual Verification Log` recording:
- Date and time of manual test execution.
- Executed curl commands / probe scripts with actual HTTP status codes and responses.
- Observed UI behaviors across English and Amharic locales.
- Confirmation of database state transitions and rollback safety.

---

## 4. Pre-Commit Verification Checklist (Gate G0)

Before executing `git commit`, verify every checkbox:

- [ ] **1. Linters & Static Analysis:** `golangci-lint`, `fvm flutter analyze`, ESLint pass with 0 errors and 0 warnings.
- [ ] **2. Automated Unit Tests:** 100% of unit tests pass with $\ge 80\%$ line coverage overall ($\ge 95\%$ on financial/auth).
- [ ] **3. Automated Integration Tests:** 100% of repository and transport integration tests pass.
- [ ] **4. Database Manual Rollback Drill:** Migration `.up.sql` and `.down.sql` executed and verified cleanly.
- [ ] **5. Backend HTTP Live Probes:** Endpoints probed with valid and invalid payloads; success envelopes and RFC 9457 errors verified.
- [ ] **6. Mobile Ergonomics & Localization:** 56dp touch targets confirmed; English and Amharic screens inspected without untranslated strings.
- [ ] **7. Documentation Synchronization:** `CONTEXT.md` and `REQUIREMENTS_TRACKER.md` staged in the same commit.
- [ ] **8. Commit Message Compliance:** Commit message includes `Tests:`, `Manual-Tests:`, and `Docs:` footers.
- [ ] **9. Multi-Machine Database Fixtures:** Any test seed data or database inserts used during manual testing are committed in idempotent SQL scripts under `scripts/fixtures/` with static UUIDs.
