# Task Protocol 05: Security Review & OWASP ASVS Level 3 Audit — Dink Machinery

**Protocol ID:** `PROTO-05-SECURITY-REVIEW`  
**Target:** AI Agents, Security Engineers & Platform Architects  
**Scope:** Verification of application security controls, cryptographic standards, zero-PII enforcement, automated SAST/DAST audits, and container hardening.  
**Governing Documents:** `CODE_STANDARDS.md § 1`, `ERROR_CATALOG.md`, `TECH_STACK.md § 6`, `.github/workflows/security.yml`  
**Compliance Standard:** OWASP Application Security Verification Standard (ASVS) 4.0.3 Level 3

---

## 1. Executive Summary & Security Objectives

In **Dink Machinery**, security is not a post-hoc QA check; it is an active architectural invariant. As a high-value marketplace orchestrating physical heavy equipment leases and multi-rail financial transactions (Telebirr, CBEBirr, MPESA), the platform adheres to **OWASP ASVS Level 3** (High-Assurance / Critical Infrastructure).

### Core Security Mandates
1. **Zero-Trust Identity & Session Lifecycle:** Every boundary interaction requires cryptographically verified identity, short-lived JWTs, and strictly rotated refresh credentials.
2. **Zero PII Exposure:** Personal Identifiable Information (phone numbers, National IDs, TINs, bank details) must be masked across all logging systems, omitted from client error envelopes, and encrypted at rest with AES-256-GCM.
3. **In-Transaction Financial Defense:** Double-entry ledger writes must balance to zero minor units, execute within isolated database transactions, and reject tampering or floating-point conversions.
4. **Idempotency & Replay Immunity:** Mutating endpoints require unique, cached `Idempotency-Key` headers to protect against race conditions and network replay attacks.
5. **Malware Defense on Storage Ingress:** Media attachments (machinery inspection photos, operator certifications, proof-of-insurance documents) must be synchronously inspected by ClamAV before permanent object storage.

---

## 2. ASVS Level 3 Domain Verification Matrix

Every feature partition must verify against these specific verification domains prior to code freeze.

```
                              OWASP ASVS Level 3 Audit
                                         │
     ┌───────────────────┬───────────────┴───────────────┬───────────────────┐
     ▼                   ▼                               ▼                   ▼
[V2: Identity &      [V4: Access        [V5: Validation & Sanitization] [V8: Data Protection
 Session Control]     Control / IDOR]     - Parameterized SQL (sqlc)      & Zero-PII]
 - TOTP / WebAuthn   - Row-Level Scopes   - ClamAV File Scanning          - AES-256-GCM
 - Bearer Auth       - Tenant Isolation   - RFC 9457 Problem Details      - Log Masking
 - Replay Lock       - Admin Clearances   - Integer Minor-Unit Money      - Secret Scrubbing
```

### Domain V2: Authentication & Identity Management
- **V2.1 Password & OTP Governance:** OTP verification sessions are rate-limited to 5 attempts maximum with a 120-second TTL. Exhaustion locks the session for 15 minutes (`AUTH_OTP_MAX_ATTEMPTS`).
- **V2.2 Bearer Token Lifecycles:**
  - Access Token (JWT): Signed via asymmetric ECDSA (ES256) or HMAC-SHA256 with cryptographically generated 256-bit keys; TTL strictly capped at 900 seconds (15 minutes).
  - Refresh Token: Stored in Valkey with single-use rotation semantics. Replay of an invalidated refresh token immediately revokes the entire active token family.
- **V2.3 Client Storage Isolation:** Mobile devices store refresh tokens exclusively in `flutter_secure_storage` (iOS Keychain / Android EncryptedSharedPreferences). Admin web browsers store session tokens in `HttpOnly`, `Secure`, `SameSite=Strict` cookies.

### Domain V4: Access Control & Authorization (RBAC & IDOR)
- **V4.1 Insecure Direct Object Reference (IDOR) Defense:** Handlers must verify resource ownership within the database query, rather than relying on route parameters alone:
  ```sql
  -- Prohibited: SELECT * FROM bookings WHERE id = $1;
  -- Compliant:
  SELECT * FROM bookings 
  WHERE id = $1 AND (renter_id = $2 OR fleet_owner_id = $2);
  ```
- **V4.2 Role-Based Access Control (RBAC):** Personas (`RENTER`, `OWNER`, `OPERATOR`, `ADMIN`) are strictly validated at router middleware. Admin endpoints (`/api/v1/admin/*`) require explicit administrative permission claims in the JWT payload.

### Domain V5: Validation, Sanitization & Input Boundaries
- **V5.1 Strict Boundary Schemas:** Request payloads must validate against OpenAPI constraints via `go-playground/validator/v10` (backend) and `zod` (admin). Any schema violation returns HTTP 400 with an RFC 9457 `invalidParams` array.
- **V5.2 SQL Injection Elimination:** All database interactions are executed via `sqlc` with native binary parameter bindings (`$1, $2`). Dynamic string interpolation, unescaped concatenation, or raw runtime SQL formatting is banned.
- **V5.3 File Ingress & Antivirus Pipeline:**
  - Uploaded files are restricted to allowed MIME types (`image/jpeg`, `image/png`, `image/webp`, `application/pdf`).
  - Maximum upload size capped at 15 MB per file.
  - Files are streamed through the ClamAV daemon (`tcp://clamav:3310`). Uploads triggering malware flags are deleted immediately and return HTTP 422 (`INSPECT_MEDIA_SCAN_FAILED`).

### Domain V8: Data Protection & Zero-PII Invariant
- **V8.1 Cryptography at Rest:** Sensitive identification documents (TIN certificates, National ID numbers, bank routing tokens) are encrypted before column insertion using `crypto/aes` with 256-bit GCM mode and random 12-byte initialization vectors (IV).
- **V8.2 Telemetry & Log Sanitization:**
  - Structured logs emitted by `uber-go/zap` run through masking interceptors.
  - Mandatory masking patterns:
    - Telephone numbers: `+251911****44`
    - Tax Identification Numbers: `TIN: *****1234`
    - National ID / Kebele Cards: `ID: *****890`
    - Bank Account Numbers: `ACC: *******5678`
- **V8.3 Client Error Sanitization:** Under zero circumstances may database driver error messages, SQL constraints, or stack traces leak into client error payloads. Unhandled panics emit generic `SYS_INTERNAL_ERROR` referencing only the correlation `requestId`.

---

## 3. Automated Security Tooling Matrix

The CI pipeline runs automated security scanning on every commit. Pull requests failing any check are automatically blocked.

| Tool | Category | Target | Execution Command | Blocking Criteria |
|---|---|---|---|---|
| **Gitleaks** | Secret Detection | Git History & Working Tree | `gitleaks detect --verbose` | $\ge 1$ leaked secret, API key, private key, or credential. |
| **Semgrep** | SAST | Go Backend & Next.js Admin | `semgrep --config p/owasp-top-ten --error` | $\ge 1$ High/Critical vulnerability or unescaped query. |
| **Govulncheck** | Vulnerability Scan | Go Dependencies & Toolchain | `govulncheck ./...` | Known CVEs with available call-path reaches. |
| **pnpm audit** | Vulnerability Scan | Node.js / Next.js Dependencies | `pnpm audit --audit-level=high` | High or Critical vulnerabilities in direct/transitive trees. |
| **Trivy** | Container SAST | Dockerfiles & Base Images | `trivy image dink-backend:latest` | Critical vulnerabilities or unpatched OS packages. |
| **Spectral CLI** | Contract Linting | `contracts/openapi.yaml` | `pnpm dlx @stoplight/spectral-cli lint ...` | Schema validation errors or floating-point currency fields. |

---

## 4. Step-by-Step Security Review Workflow

Run through this procedure during local development before submitting a pull request:

```bash
# ==============================================================================
# Step 1: Scan for Accidental Secret & Credential Commits
# ==============================================================================
gitleaks detect --source . --verbose

# ==============================================================================
# Step 2: Audit Zero-PII Masking in Code & Logging
# ==============================================================================
# Verify no naked sensitive fields are logged directly
grep -rnEi 'logger\.(Info|Debug|Warn|Error)\(.*(password|pin|tin|national_id|account_number)' backend/ admin/ mobile/ && {
  echo "CRITICAL: Potential unmasked PII logging detected!"
  exit 1
} || echo "Zero-PII logging audit passed."

# ==============================================================================
# Step 3: Run Go Static Analysis & Vulnerability Checks
# ==============================================================================
cd backend
govulncheck ./...
golangci-lint run ./...
cd ..

# ==============================================================================
# Step 4: Run Node Dependency Audits
# ==============================================================================
cd admin
pnpm audit --audit-level=high
cd ..

# ==============================================================================
# Step 5: Run Static Security Rules (Semgrep)
# ==============================================================================
semgrep --config p/golang --config p/owasp-top-ten --error backend/ admin/
```

---

## 5. Vulnerability Triage, CVSS Scoring & Remediation SLAs

Identified security issues must be classified under CVSS v3.1 and remediated within the following strict timeframes:

| Severity Tier | CVSS v3.1 Range | Example Scenarios | Remediation SLA | Deployment Procedure |
|---|---|---|---|---|
| **Critical** | `9.0 - 10.0` | Remote Code Execution, Auth Bypass, Unauthenticated Escrow Drain, Double-Spend Exploit. | **Immediate (< 24 Hours)** | Emergency hotfix branch; immediate production rollout bypassing normal sprint windows. |
| **High** | `7.0 - 8.9` | IDOR accessing competitor fleets, Unmasked PII in server logs, SQL injection vector in internal admin tools. | **$\le$ 7 Calendar Days** | Fast-tracked PR; mandatory verification drill by Lead Architect. |
| **Medium** | `4.0 - 6.9` | Ingress rate-limit bypass, Verbose error message disclosing package path, Missing HTTP security header. | **$\le$ 30 Calendar Days** | Scheduled into active feature sprint. |
| **Low** | `0.1 - 3.9` | Minor information disclosure (server banner version), TLS 1.2 legacy cipher suite presence. | **$\le$ 90 Calendar Days** | Backlog hygiene and dependency upgrade cycles. |

---

## 6. Pre-Merge Security Review Sign-Off Checklist

Before any feature partition or architectural change is approved for merge into `main`, the PR must document compliance with this checklist:

### Identity & Authentication
- [ ] No hardcoded passwords, API keys, or JWT secret keys exist in code or Git history.
- [ ] OTP verification flows enforce maximum 5 attempts and a 120-second timeout window.
- [ ] Session tokens use secure, partitioned storage (`flutter_secure_storage` or `HttpOnly` cookies).

### Authorization & State Integrity
- [ ] All database queries assert tenant/user ownership to eliminate IDOR vulnerabilities.
- [ ] State transitions assert optimistic concurrency control (`WHERE id = $1 AND version = $2`).
- [ ] Financial operations execute within atomic database transactions with double-entry balance checks (`BK-17-BE`).
- [ ] Mutating endpoints validate and cache `Idempotency-Key` headers.

### Data Validation & Protection
- [ ] Inputs are sanitized and bound using parameterization; raw string concatenation in SQL is absent.
- [ ] File uploads stream through ClamAV antivirus inspection prior to permanent S3 storage.
- [ ] PII attributes (phone, National ID, TIN, bank accounts) are encrypted at rest with AES-256-GCM.
- [ ] Structured logging asserts masking on all customer identification and financial properties.
- [ ] Errors conform strictly to RFC 9457 Problem Details; zero stack traces or internal engine errors leak to clients.

### Automated Tooling Gates
- [ ] `gitleaks detect` passes with 0 leaks detected.
- [ ] `govulncheck` and `pnpm audit` pass with 0 high/critical vulnerabilities.
- [ ] Semgrep SAST scans pass with 0 rule violations.
- [ ] Spectral contract validation passes with 0 errors.