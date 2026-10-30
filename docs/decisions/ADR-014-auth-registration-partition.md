# ADR-014: Dual-Token ES256 JWT Session Architecture and Refresh Token Family Rotation

- **Status:** Accepted
- **Date:** 2026-10-26
- **Decision Owner:** Dink Machinery Architecture Review Board
- **Affected Partitions:** 04-user-roles, 05-01, 06-01, 09-non-functional-requirements
- **Related Requirements:** AUTH-01-BE, AUTH-06-BE, ASVS-V2.1, ASVS-V3.1

---

## 1. Context & Problem Statement

Dink Machinery users operate across varied form factors: field operators on mobile devices across intermittent cellular signals, corporate procurement managers on desktop browsers, and administrators on back-office consoles. Centralized server-side session lookups across all API calls create database bottlenecks, while long-lived static tokens pose severe security risks if intercepted. Furthermore, if a refresh token is stolen, the system must detect the replay attack and immediately invalidate all sessions in that family.

---

## 2. The 5-Point Qualification Evaluation

- [x] **Alternative Test:** Evaluated dual-token asymmetric ES256 JWT with Refresh Token Rotation (RTR) vs. stateful database sessions vs. symmetric HMAC-SHA256 JWTs.
- [x] **Surprise Test:** Uses asymmetric ECDSA P-256 (`ES256`) where edge proxies and services verify tokens using public keys without database calls, while refresh tokens use family rotation.
- [x] **Constraint Test:** OWASP ASVS Level 3 session management and token replay detection compliance.
- [x] **Irreversibility Test:** Central security and authentication framework spanning mobile, admin web, and backend runtimes.
- [x] **Scope Test:** Directly governs user authentication, token issuance, logout, role evaluation, and API gateway authorization.

---

## 3. Considered Options

### Option 1: Dual-Token Asymmetric ES256 JWT + Valkey Refresh Token Family Rotation (RTR)
- **Summary:** Short-lived (15-minute) Access Tokens signed using ECDSA P-256 (`ES256`). Edge gateways and backend modules verify signatures locally using the public key with zero database queries. Long-lived (30-day) Refresh Tokens are stored as SHA-256 hashes in Valkey organized by token families. Every refresh issues a new refresh token and invalidates the previous one. If an invalidated refresh token is reused, the entire family is revoked immediately.
- **Pros:**
  - Stateless high-throughput request authentication ($< 1\text{ms}$ latency; zero database hits).
  - Immediate detection and neutralization of stolen refresh tokens via family rotation.
  - Asymmetric keys: private signing key resides in OpenBao KMS; public keys distributed widely.
- **Cons:**
  - Revoking an access token prior to its 15-minute expiration requires checking an in-memory Valkey blocklist for critical events (password changes, account suspensions).

### Option 2: Pure Stateful Database Sessions
- **Summary:** Traditional server-side sessions stored in PostgreSQL or Valkey, passing a random session ID cookie or bearer token.
- **Pros:**
  - Instant revocation by deleting the session row.
- **Cons:**
  - Every single API request requires a database/cache lookup, creating a central bottleneck under high traffic.

### Option 3: Symmetric HMAC-SHA256 (`HS256`) Tokens
- **Summary:** Standard JWT tokens using a shared secret key.
- **Pros:**
  - Simple setup with a single secret string.
- **Cons:**
  - The shared secret must be distributed to all microservices and edge proxies; compromising any single service compromises the entire platform's token signing authority.

---

## 4. Chosen Solution & Architectural Rationale

We adopt **Option 1: Dual-Token Asymmetric ES256 JWT + Valkey Refresh Token Family Rotation (RTR)**.

1. **Access Token Specifications:**
   - Standard: RFC 7519 JWT.
   - Algorithm: ECDSA P-256 with SHA-256 (`ES256`).
   - Lifetime: 15 minutes.
   - Claims: `sub` (User UUID), `role` (super_admin, owner, renter, etc.), `country` (ET, KE), `org_id` (if corporate), `exp`, `jti` (Token UUID).
2. **Refresh Token Family Rotation (RTR):**
   - Refresh tokens are 256-bit cryptographically secure random strings.
   - Stored in Valkey: `auth:refresh:{token_family_id}`.
   - Replay Detection: If Token $N$ is submitted after Token $N+1$ has already been issued, the backend immediately purges the entire token family, killing all active sessions for that user on that device.
3. **Web Back-Office Security:**
   - Admin access and refresh tokens are delivered strictly in `HttpOnly`, `SameSite=Strict`, `Secure` cookies with Path scoping.

---

## 5. Rejected Alternatives & Technical Justification

- **Option 2 (Pure Stateful Sessions):** Rejected due to performance degradation and cache saturation from checking session tables on every single HTTP read and write.
- **Option 3 (Symmetric HS256):** Rejected because sharing a single symmetric signing key across edge gateways, admin services, and worker jobs violates the principle of least privilege and ASVS L3 cryptographic requirements.

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - Micro-second token validation with zero database load during normal browsing.
  - Automatic isolation and lockout if a user's mobile device is cloned or token intercepted.
  - Strong cryptographic non-repudiation with ECDSA P-256.
- **Negative / Technical Debt:**
  - Requires maintaining a Valkey blocklist for emergency immediate session revocations.
- **Operational / Observability Impacts:**
  - Token refresh family replay anomalies trigger high-priority security alarms to the SIEM.

---

## 7. Verification Strategy & Invariant Tests

- **Token Replay Attack Simulation:** `backend/internal/modules/auth/application/token_rotation_test.go` simulates an attacker submitting a previously used refresh token; asserts that all tokens in the family are revoked and HTTP 401 is returned.
- **Signature Algorithm Enforcement:** Asserts that the JWT parser rejects `none` and `HS256` algorithms, strictly requiring `ES256`.

---

## 8. Revisit Conditions

- Emerging quantum computing risks necessitating migration from ECDSA P-256 to post-quantum signature algorithms (e.g. ML-DSA / Dilithium).
