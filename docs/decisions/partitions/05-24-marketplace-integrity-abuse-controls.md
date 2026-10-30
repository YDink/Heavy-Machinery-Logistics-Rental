# Partition Decisions: 05-24 — Marketplace Integrity, Anti-Circumvention & Abuse Controls Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-24`.

- **Governing System ADRs:** `ADR-009`, `ADR-010`
- **Subsystem Scope:** backend
- **Requirement Targets:** ABUSE-01-BE to ABUSE-10-BE

---

## D-05-24-001: In-Flight Phone Number and Bank Detail Redaction in Pre-Booking Chat Messages

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [ABUSE-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Renters and owners routinely attempt to bypass platform commission fees by sharing phone numbers, bank accounts, or telegram handles in pre-booking inquiry chat messages, intending to transact offline and cut out the marketplace.

### 2. Considered Alternatives
1. **Option A:** Synchronous in-flight regex pattern masking in chat transport layer, replacing phone numbers, bank account numbers, and telegram links with redaction warnings before saving to database
   - *Pros:* Immediate interception; contact details never reach the recipient screen; warns both parties of policy violation; preserves platform commission.
   - *Cons:* Requires sophisticated regex matching to prevent obfuscation tactics (e.g. spelled-out digits, spaces, symbols).
2. **Option B:** Post-chat administrative review by support agents
   - *Pros:* Human understanding of context.
   - *Cons:* Contact details already exchanged; off-platform deal already sealed before agent reviews log.
3. **Option C:** Disabling chat messaging entirely until booking is confirmed
   - *Pros:* Zero circumvention via chat.
   - *Cons:* Prevents legitimate pre-booking technical inquiries (e.g. job site soil conditions, attachment requirements), causing booking drop-off.

### 3. Chosen Solution & Technical Rationale
Option A: Synchronous in-flight regex redaction in chat message handlers. Pre-booking messages pass through an anti-circumvention regex filter detecting E.164 phone patterns, domestic 10-13 digit bank account formats (CBE, Awash, Dashen), and social media handles, replacing matches with [REDACTED_CONTACT_INFO] and logging an abuse event.

### 4. Rejected Alternatives
- Option B: Rejected because asynchronous review is too late to stop off-platform disintermediation.
- Option C: Rejected because contractors need to verify mechanical site compatibility before committing large deposits.

### 5. Consequences & Invariants
- **Invariants Maintained:** Direct personal contact details are redacted from all user chat messages prior to booking confirmation.
- **Operational Trade-offs:** Edge-case false positives on model part numbers or serial numbers must be whitelisted.

### 6. Verification
- **Test File:** `backend/internal/modules/chat/application/redaction_filter_test.go`
- **Assertion:** TestChatRedaction_MasksEthiopianPhoneNumber asserts 0911223344 and +251911223344 are replaced with redaction notices.

### 7. Revisit Trigger
- Revisit if natural language processing (NLP) models are needed to detect semantic circumvention intent in Phase 3.

---

## D-05-24-002: Off-Platform Circumvention Detection Engine with Algorithmic Disintermediation Scoring

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [ABUSE-05-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Sophisticated bad actors attempt circumvention by initiating multiple booking holds, extracting owner yard locations, and then cancelling to conduct offline deals. Chronic repeat offenders degrade marketplace liquidity and commission revenue.

### 2. Considered Alternatives
1. **Option A:** Algorithmic disintermediation scoring: calculates risk score based on cancel-after-inquiry rates, phone redaction hits, and radius proximity, triggering automated account suspension reviews
   - *Pros:* Automated pattern detection; flags chronic offenders; protects honest fleet owners and platform revenue.
   - *Cons:* Scoring weights must be tuned to prevent false accusations against legitimate contractors with high project variability.
2. **Option B:** Strict automatic account ban after 2 consecutive cancellations
   - *Pros:* Aggressive enforcement.
   - *Cons:* Bans legitimate contractors whose projects were postponed due to weather or permit delays.
3. **Option C:** Zero automated scoring; rely on competitor tips
   - *Pros:* Zero algorithmic complexity.
   - *Cons:* Complete failure to protect platform monetization.

### 3. Chosen Solution & Technical Rationale
Option A: Algorithmic disintermediation scoring. The system computes a cumulative risk score: +20 points per redacted contact message, +15 points per cancellation within 30 minutes of owner contact disclosure, +30 points per duplicate bank slip upload. Users exceeding 70 points are flagged for Compliance Review (`AD-13-WEB`).

### 4. Rejected Alternatives
- Option B: Rejected because rigid cancellation thresholds penalize honest construction contractors facing normal job site delays.
- Option C: Rejected because voluntary competitor reporting is too rare to deter platform leakage.

### 5. Consequences & Invariants
- **Invariants Maintained:** Users with an abuse risk score > 85 are automatically blocked from initiating new booking holds pending compliance audit.
- **Operational Trade-offs:** Compliance officers must review flagged abuse queues weekly to evaluate suspensions.

### 6. Verification
- **Test File:** `backend/internal/modules/admin/application/abuse_scoring_test.go`
- **Assertion:** TestAbuseScoring_FlagsHighRiskUser asserts 3 contact redactions + 2 cancellations triggers high-risk review status.

### 7. Revisit Trigger
- Revisit if machine learning graph neural networks are integrated to detect multi-account collusion rings.

---

