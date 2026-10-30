# Partition Decisions: 05-13 — Machinery Insurance Verification & Document Governance Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-13`.

- **Governing System ADRs:** `ADR-005`, `ADR-011`, `ADR-012`
- **Subsystem Scope:** backend, admin
- **Requirement Targets:** INS-01-BE to INS-10-BE

---

## D-05-13-001: Comprehensive Machine Policy Ingestion with Mandatory S3 Document Archival

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [INS-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Heavy equipment carries severe operational and third-party liability risks. Uninsured machinery operating on civil construction projects poses catastrophic financial liabilities. The platform must capture, verify, and persist valid Comprehensive Machine Insurance (CMI) policies before equipment can be listed.

### 2. Considered Alternatives
1. **Option A:** Mandatory policy upload (insurer, policy number, coverage amount minor, expiry date) with PDF/image scan stored in encrypted SeaweedFS S3 and verified by compliance officers
   - *Pros:* Complete legal evidentiary chain; immutable document archival; compliance with insurance regulatory standards.
   - *Cons:* Requires compliance officer review before listing activation.
2. **Option B:** Self-certification checkbox stating the machine is insured without document proof
   - *Pros:* Instant listing activation.
   - *Cons:* Catastrophic legal liability when uninsured machines suffer major structural damage or fatal accidents.
3. **Option C:** Direct automated API verification with Ethiopian and Kenyan insurance carriers
   - *Pros:* Instant digital policy verification.
   - *Cons:* No unified digital insurance API exists in Ethiopia or Kenya; insurance records are largely paper-based.

### 3. Chosen Solution & Technical Rationale
Option A: Mandatory policy metadata ingestion paired with scanned document storage in encrypted SeaweedFS S3. Listings remain in PENDING_INSURANCE_VERIFICATION until a compliance officer verifies the document in the Admin Operations Queue (`AD-02-WEB`).

### 4. Rejected Alternatives
- Option B: Rejected because platform reputation and contractor trust require guaranteed insurance coverage.
- Option C: Rejected because regional insurance carriers do not offer public verification APIs.

### 5. Consequences & Invariants
- **Invariants Maintained:** A fleet asset cannot be published to the active search catalog without an approved, non-expired insurance record.
- **Operational Trade-offs:** New fleet onboarding requires up to 24 hours for manual compliance document verification.

### 6. Verification
- **Test File:** `backend/internal/modules/machinery/application/insurance_verification_test.go`
- **Assertion:** TestInsurance_BlocksListingWithoutValidPolicy asserts fleet activation fails if insurance status is not APPROVED.

### 7. Revisit Trigger
- Revisit if the Ethiopian Insurers Association (EIA) launches an automated policy verification gateway.

---

## D-05-13-002: Automated 30-Day and 7-Day Expiry Escalations with Auto-Suspension Guards

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [INS-06-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Insurance policies expire during active listing lifecycles. If an owner fails to renew their policy, the machine becomes illegal to rent. The system must notify owners in advance and automatically suspend expired listings.

### 2. Considered Alternatives
1. **Option A:** Daily automated cron job checking insurance expiry dates, triggering multi-tier alerts (30-day, 7-day, 1-day) and auto-suspending listings upon expiry
   - *Pros:* Zero manual monitoring required; guaranteed prevention of renting expired machines; clear owner renewal reminders.
   - *Cons:* Requires automated calendar scheduling and notifications.
2. **Option B:** Manual periodic audits by compliance officers
   - *Pros:* Zero automated cron infrastructure.
   - *Cons:* High risk of expired policies slipping through between audits, creating uninsurable accidents.
3. **Option C:** Warning banner on listing without suspending availability
   - *Pros:* Maintains owner listing presence.
   - *Cons:* Contractors can still book an uninsured machine, transferring liability to the marketplace.

### 3. Chosen Solution & Technical Rationale
Option A: Daily automated expiry cron. The system scans insurance_documents: at T-30 days, sends renewal reminder SMS/Push; at T-7 days, escalates to high-priority alert; at T-0 (expiry date), automatically updates fleet status to SUSPENDED_INSURANCE_EXPIRED.

### 4. Rejected Alternatives
- Option B: Rejected due to human oversight and liability exposure.
- Option C: Rejected because allowing bookings on expired insurance violates marketplace safety invariants.

### 5. Consequences & Invariants
- **Invariants Maintained:** Any fleet asset whose insurance expiry date is < CURRENT_DATE is automatically blocked from accepting new reservations.
- **Operational Trade-offs:** Owners whose policies expire during an active rental must complete the rental but cannot accept future bookings until renewed.

### 6. Verification
- **Test File:** `backend/internal/modules/machinery/application/insurance_expiry_cron_test.go`
- **Assertion:** TestInsuranceExpiry_AutoSuspendsExpiredFleets asserts expired fleet status transitions to SUSPENDED_INSURANCE_EXPIRED.

### 7. Revisit Trigger
- Revisit if automated multi-year fleet insurance policies eliminate annual renewal cycles.

---

