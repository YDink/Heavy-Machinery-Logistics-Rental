# ADR-011: Pluggable Identity Verification Framework with Manual Fallback

- **Status:** Accepted
- **Date:** 2026-10-26
- **Decision Owner:** Dink Machinery Architecture Review Board
- **Affected Partitions:** 05-01, 05-13, 07-04
- **Related Requirements:** AUTH-04-BE, KYC-01-BE, AD-13-BE, ECA-KYC-01

---

## 1. Context & Problem Statement

Dink Machinery operates across multi-national jurisdictions where government digital identity platforms are in varying stages of national rollout. In Ethiopia, the Fayda National ID (eKYC) system is rolling out nationally but experiences periodic gateway maintenance and connectivity dips; in Kenya, the Integrated Population Registration System (IPRS) and Huduma Namba govern citizen attestation. Rigidly hardcoding a single identity API causes nationwide onboarding paralysis during third-party gateway outages.

---

## 2. The 5-Point Qualification Evaluation

- [x] **Alternative Test:** Evaluated pluggable identity adapter framework vs. hardcoded single-provider integration vs. pure manual document review.
- [x] **Surprise Test:** Implements an automated circuit-breaker that seamlessly falls back to an administrative manual review queue without blocking user registration.
- [x] **Constraint Test:** Regulatory mandate requiring verified identity before heavy machinery asset custody or rental operations.
- [x] **Irreversibility Test:** Core architectural contract governing user identity records and compliance workflows across all jurisdictions.
- [x] **Scope Test:** Directly affects user registration, owner fleet verification, operator licensing, and back-office admin queues.

---

## 3. Considered Options

### Option 1: Pluggable `IdentityProvider` Adapter Interface with Automated Queue Fallback
- **Summary:** Define an internal Go interface `IdentityProvider` with jurisdictional implementations (e.g. `FaydaAdapter` for ET, `IprsAdapter` for KE). If the automated provider returns a high confidence match, verification is immediate; if the provider is unreachable, times out (>5s), or returns an ambiguous confidence score, the user document is routed automatically to the Admin Manual Verification Docket (`AD-13-BE`).
- **Pros:**
  - High resilience: third-party API downtime never halts user onboarding.
  - Zero code modifications to swap or add new national identity providers.
  - Strict audit trail tracking automated vs. manual verification decisions.
- **Cons:**
  - Requires maintaining an operational back-office manual review queue.

### Option 2: Hardcoded Direct Integration with National eKYC API
- **Summary:** Directly call the Ethiopian Fayda eKYC API endpoint inside the auth registration handler.
- **Pros:**
  - Simple initial implementation.
- **Cons:**
  - Any Fayda outage completely halts all registration and rental activity in Ethiopia; impossible to onboard Kenyan users without substantial refactoring.

### Option 3: Pure Manual Document Inspection
- **Summary:** All users upload identity photos, and support staff manually inspect every document.
- **Pros:**
  - Zero third-party API integration dependency.
- **Cons:**
  - Extreme operational human bottleneck; onboarding takes 24–48 hours instead of seconds.

---

## 4. Chosen Solution & Architectural Rationale

We adopt **Option 1: Pluggable `IdentityProvider` Adapter Interface with Automated Queue Fallback**.

1. **Pluggable Architecture:**
   ```go
   type VerificationResult struct {
       Status           VerificationStatus // VERIFIED, REJECTED, REQUIRES_MANUAL_REVIEW
       ConfidenceScore  float64
       NationalIDNumber string
       FullName         string
       ProviderRawRef   string
   }

   type IdentityProvider interface {
       CountryCode() string
       VerifyCitizen(ctx context.Context, req IdentityVerificationRequest) (*VerificationResult, error)
   }
   ```
2. **Dynamic Provider Resolution:**
   The `IdentityService` resolves the provider using the active country context (`X-Country-Code`).
3. **Resilient Circuit Breaking & Manual Docket (`AD-13-BE`):**
   - Automated calls enforce a 5-second timeout.
   - If the provider is offline or confidence is below threshold ($< 0.85$), the system marks `identity_verifications.status = 'REQUIRES_MANUAL_REVIEW'` and enqueues the record in the Next.js Admin Verification Queue.
   - The user is placed in `KYC_PENDING` status, allowing them to browse listings while awaiting compliance sign-off.

---

## 5. Rejected Alternatives & Technical Justification

- **Option 2 (Hardcoded Fayda Integration):** Rejected due to single-point-of-failure vulnerability and complete lack of geographic portability to Kenya or other regional markets.
- **Option 3 (Pure Manual Review):** Rejected because human review costs are prohibitive and turn-around latency severely depresses conversion rates.

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - Onboarding continues seamlessly even during national telecom or government server maintenance windows.
  - Multi-country compliance rules are cleanly separated into dedicated adapter packages.
  - Transparent compliance reporting for government regulators.
- **Negative / Technical Debt:**
  - Administrative staff must staff the manual review docket to maintain the 24-hour verification SLA.
- **Operational / Observability Impacts:**
  - Metric `dink_identity_verification_attempts_total` tracks automated vs. manual routing ratios.

---

## 7. Verification Strategy & Invariant Tests

- **Provider Mock Failover Test:** `backend/internal/modules/auth/application/identity_verification_test.go` simulates a simulated Fayda API 504 Gateway Timeout; asserts that the verification status transitions to `REQUIRES_MANUAL_REVIEW` and enqueues an admin task.
- **Pluggable Interface Conformance:** Asserts all adapters satisfy the `domain.IdentityProvider` interface at compile time.

---

## 8. Revisit Conditions

- National biometric eKYC APIs achieve 99.99% verified uptime SLAs, allowing reduction of the manual review team.
