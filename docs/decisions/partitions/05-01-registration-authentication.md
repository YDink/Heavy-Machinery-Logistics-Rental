# Partition Decisions: 05-01 — User Registration & Authentication Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-01`.

- **Governing System ADRs:** `ADR-011`, `ADR-014`
- **Subsystem Scope:** backend, mobile
- **Requirement Targets:** REG-01-BE to REG-09-BE, AUT-01-BE to AUT-11-BE

---

## D-05-01-001: Phone OTP Authentication with Leaky Bucket Rate Limiting and Carrier Failover

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [REG-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
User onboarding relies on mobile phone numbers (MSISDN) in Ethiopia (+251) and Kenya (+254). SMS delivery failure rates on regional telcos can spike during network congestion. Furthermore, SMS toll fraud and brute-force credential stuffing must be prevented without frustrating legitimate users.

### 2. Considered Alternatives
1. **Option A:** Leaky bucket distributed rate limiting in Valkey with multi-carrier SMS gateway fallback (Ethio Telecom -> Safaricom ET -> WhatsApp)
   - *Pros:* Eliminates toll fraud; prevents brute-force; maintains high OTP delivery rates during carrier degradation.
   - *Cons:* Requires maintaining integration contracts with multiple SMS gateways.
2. **Option B:** Single carrier SMS gateway with fixed 3-attempt limit per phone number
   - *Pros:* Simple implementation with a single provider.
   - *Cons:* Carrier outages completely block all user signups in that region; rigid lockout frustrates legitimate contractors on poor reception.
3. **Option C:** Email-first registration with optional phone verification
   - *Pros:* Zero SMS costs; familiar for corporate users.
   - *Cons:* Fails in East African construction reality where field operators and machine owners operate via mobile phones and lack frequent email access.

### 3. Chosen Solution & Technical Rationale
Option A: Leaky bucket distributed rate limiting in Valkey with multi-carrier SMS gateway fallback. OTP verification sessions allow a maximum of 5 attempts within a 120-second TTL. If primary SMS delivery fails within 45 seconds, the system automatically cascades to secondary carrier or WhatsApp Business API.

### 4. Rejected Alternatives
- Option B: Rejected because single-carrier reliance causes unacceptable registration drop-offs during frequent regional telecom dips.
- Option C: Rejected because mobile phone authentication is the universal identification anchor across target demographics.

### 5. Consequences & Invariants
- **Invariants Maintained:** OTP verification attempts are strictly bounded to 5 tries per 120s window. Plaintext OTPs are never logged or persisted to disk.
- **Operational Trade-offs:** Requires monitoring multi-gateway SMS delivery latencies and maintaining balance across telco accounts.

### 6. Verification
- **Test File:** `backend/internal/modules/auth/application/otp_service_test.go`
- **Assertion:** TestOTPService_RateLimiting asserts that the 6th verification attempt triggers HTTP 429 Too Many Requests and locks out the challenge.

### 7. Revisit Trigger
- Revisit if national biometric passkey authentication (WebAuthn/FIDO2) reaches ubiquitous smartphone penetration in East Africa.

---

## D-05-01-002: Dual-Token ES256 Minting with Valkey Refresh Token Family Rotation

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [AUT-06-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Field operators and contractors require persistent mobile sessions across intermittent connectivity, while back-office web admins require strict session isolation. Token theft or credential leakage must be instantly detectable without database bottlenecks on every request.

### 2. Considered Alternatives
1. **Option A:** Asymmetric ECDSA P-256 (ES256) short-lived JWT access tokens paired with Valkey refresh token family rotation (RTR)
   - *Pros:* Edge verification without database queries; immediate invalidation of entire token family if refresh token replay is detected.
   - *Cons:* Requires KMS management for ECDSA private key and in-memory blacklist for emergency revocation.
2. **Option B:** Stateful database sessions stored in PostgreSQL
   - *Pros:* Instant revocation by deleting session rows.
   - *Cons:* Database bottleneck under high concurrent request volumes; degrades mobile offline performance.
3. **Option C:** Symmetric HS256 long-lived JWTs
   - *Pros:* Simple single-secret configuration.
   - *Cons:* High vulnerability to token interception; shared secret must be distributed across multiple services.

### 3. Chosen Solution & Technical Rationale
Option A: Asymmetric ECDSA P-256 (ES256) dual-token architecture. Access tokens expire in 15 minutes. Refresh tokens are tracked in Valkey with family IDs. Replay of an invalidated refresh token immediately purges the entire token family.

### 4. Rejected Alternatives
- Option B: Rejected due to severe database I/O overhead on high-frequency mobile GPS and catalog requests.
- Option C: Rejected because shared symmetric secret keys violate OWASP ASVS Level 3 cryptographic separation.

### 5. Consequences & Invariants
- **Invariants Maintained:** Zero plaintext tokens in logs. Any reused refresh token invalidates all associated user sessions immediately.
- **Operational Trade-offs:** Requires OpenBao KMS for private key signing and Valkey memory monitoring.

### 6. Verification
- **Test File:** `backend/internal/modules/auth/application/token_rotation_test.go`
- **Assertion:** TestTokenRotation_ReplayAttack asserts that presenting a used refresh token triggers immediate family revocation and returns HTTP 401.

### 7. Revisit Trigger
- Revisit if post-quantum cryptographic standards mandate migration to lattice-based signature algorithms.

---

