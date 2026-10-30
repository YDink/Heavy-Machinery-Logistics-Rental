# Partition Decisions: 05-07 — Rental Completion, Return Inspection & Damage Arbitration Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-07`.

- **Governing System ADRs:** `ADR-004`, `ADR-010`
- **Subsystem Scope:** backend, mobile, admin
- **Requirement Targets:** RET-01-BE to RET-14-BE, DSP-01-BE to DSP-10-BE

---

## D-05-07-001: Dual-Party Photographic Handover vs Return Image Diffing

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [RET-02-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
At rental return, disputes routinely arise over scratches, broken glass, undercarriage damage, or missing bucket teeth. Determining whether damage occurred during the rental requires objective comparison against handover baseline condition.

### 2. Considered Alternatives
1. **Option A:** Standardized 6-angle camera-enforced return inspection matching handover angles with split-screen lightbox comparison in Admin Console
   - *Pros:* Definitive visual evidence; side-by-side inspection photo alignment; rapid dispute resolution.
   - *Cons:* Requires operator to capture 6 mandatory inspection angles at both pickup and return.
2. **Option B:** Single general return photo with optional damage text description
   - *Pros:* Fast return completion.
   - *Cons:* Insufficient evidence to resolve contentious multi-thousand Birr damage claims.
3. **Option C:** In-person physical inspection by third-party marketplace adjusters
   - *Pros:* Professional neutral assessment.
   - *Cons:* Prohibitive operational cost; delays equipment demobilization by days.

### 3. Chosen Solution & Technical Rationale
Option A: Standardized 6-angle camera-enforced return inspection paired with split-screen visual diffing in the Next.js Admin Dispute Workbench (`AD-04-WEB`). Return photos must match the required angles (Front, Rear, Left, Right, Hour Meter, Cab).

### 4. Rejected Alternatives
- Option B: Rejected because ambiguous photos make legal dispute arbitration impossible.
- Option C: Rejected due to human staffing bottlenecks across rural construction corridors.

### 5. Consequences & Invariants
- **Invariants Maintained:** A return inspection cannot transition to COMPLETED without all 6 mandatory photographic angles captured and scanned clean of malware.
- **Operational Trade-offs:** Operators must spend 5-10 minutes photographing the machine at job site handover.

### 6. Verification
- **Test File:** `backend/internal/modules/inspection/application/return_inspection_test.go`
- **Assertion:** TestReturnInspection_RequiresAllSixAngles asserts that submitting 5 angles fails with HTTP 422 MISSING_INSPECTION_ANGLES.

### 7. Revisit Trigger
- Revisit if automated computer-vision edge damage detection is introduced in Phase 3.

---

## D-05-07-002: 72-Hour Statutory Administrative Dispute Resolution SLA Engine

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [DSP-05-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
When damage claims or excess hour disputes are lodged, unresolved conflicts tie up owner equipment and renter security deposits. The platform must enforce strict arbitration timelines conforming to jurisdictional SLAs (72 hours in Ethiopia, 48 hours in Kenya).

### 2. Considered Alternatives
1. **Option A:** Automated SLA state machine with cron escalation, tri-party evidence chat, and binding administrative ruling overrides
   - *Pros:* Guaranteed resolution timeline; clear audit trail; prevents open-ended dispute paralysis.
   - *Cons:* Requires administrative support staff to review escalated claims within SLA windows.
2. **Option B:** Open-ended informal negotiation between renter and owner without platform deadlines
   - *Pros:* Zero administrative overhead.
   - *Cons:* Disputes drag on for months, leading to civil litigation and marketplace abandonment.
3. **Option C:** Automatic default ruling in favor of the owner after 48 hours of silence
   - *Pros:* Zero manual review required.
   - *Cons:* Severely penalizes renters whose communication was interrupted by regional cellular blackouts.

### 3. Chosen Solution & Technical Rationale
Option A: Automated SLA state machine with dynamic jurisdictional deadlines loaded from country_configs.dispute_sla_hours. If parties do not settle mutually within the window, a cron escalates the case to Admin Arbitration, allowing binding state transitions (DISPUTE_RESOLVED_RENTER, DISPUTE_RESOLVED_OWNER, DAMAGE_PARTIAL_RESOLVED).

### 4. Rejected Alternatives
- Option B: Rejected because unresolved disputes cause customer churn and frozen equipment.
- Option C: Rejected because unilateral default judgments violate fair trade principles and procedural justice.

### 5. Consequences & Invariants
- **Invariants Maintained:** Dispute deadlines are strictly calculated using jurisdictional country_configs.dispute_sla_hours excluding national public holidays.
- **Operational Trade-offs:** Requires dedicated dispute arbitrators on duty during business hours.

### 6. Verification
- **Test File:** `backend/internal/modules/dispute/application/sla_escalation_test.go`
- **Assertion:** TestDisputeSLA_AutoEscalatesAfterDeadline asserts that an unresolved dispute transitions to ESCALATED_TO_ADMIN upon 72-hour timer expiration.

### 7. Revisit Trigger
- Revisit if dispute volume warrants automated AI settlement recommendations based on historical rulings.

---

