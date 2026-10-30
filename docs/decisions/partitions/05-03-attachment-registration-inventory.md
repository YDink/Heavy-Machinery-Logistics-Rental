# Partition Decisions: 05-03 — Attachment Registration, Compatibility & Inventory Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-03`.

- **Governing System ADRs:** `ADR-003`, `ADR-005`
- **Subsystem Scope:** backend, database
- **Requirement Targets:** ATT-01-BE to ATT-07-BE

---

## D-05-03-001: Standalone vs Host-Bundled Attachment Availability Modeling

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [ATT-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Attachments (hydraulic breakers, earth augers, trenching buckets) can be rented independently or bundled as optional accessories with a host excavator. If a host machine is rented with an attachment, that attachment must be atomically marked unavailable for standalone rental.

### 2. Considered Alternatives
1. **Option A:** Unified attachment inventory with composite status and relational join linking to active booking leases
   - *Pros:* Prevents double-booking between standalone and bundled rental channels; single atomic update.
   - *Cons:* Requires cross-table lock during booking reservation hold.
2. **Option B:** Separate inventories for standalone vs machine-attached equipment
   - *Pros:* Simple isolated scheduling calendars.
   - *Cons:* Forces fleet owners to artificially partition their physical equipment, drastically reducing asset utilization.
3. **Option C:** Manual owner confirmation of attachment availability on every booking request
   - *Pros:* Zero system concurrency complexity.
   - *Cons:* High friction and slow reservation confirmations leading to customer churn.

### 3. Chosen Solution & Technical Rationale
Option A: Unified attachment inventory with composite availability states. When a host machine booking is placed with bundled attachments, the booking transaction atomically locks and reserves both the fleet record and the attached item records within pgx.Tx.

### 4. Rejected Alternatives
- Option B: Rejected because physical attachments are shared assets that can be swapped between machines or rented standalone.
- Option C: Rejected because manual checks destroy the automated 60-minute reservation workflow.

### 5. Consequences & Invariants
- **Invariants Maintained:** An attachment cannot be in an ACTIVE or CONFIRMED lease in more than one concurrent booking interval.
- **Operational Trade-offs:** Booking confirmation transactions must acquire locks on multiple entity rows.

### 6. Verification
- **Test File:** `backend/internal/modules/machinery/application/attachment_booking_test.go`
- **Assertion:** TestAttachment_BundledLockPreventsStandaloneRental asserts that renting an excavator with a breaker blocks standalone rental of that breaker.

### 7. Revisit Trigger
- Revisit if quick-coupler telematics sensors automate mechanical attachment detection in Phase 3.

---

## D-05-03-002: Mechanical Compatibility Validation via Pin-Diameter and Weight Class Verification

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [ATT-04-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Mounting incompatible attachments (e.g. 5-ton excavator bucket on a 20-ton machine, or mismatched hydraulic flow) causes severe job site damage, hydraulic pump failure, and safety hazards.

### 2. Considered Alternatives
1. **Option A:** Automated compatibility validation engine matching host machine weight class, pin diameter, and hydraulic flow ranges
   - *Pros:* Eliminates human error; prevents hazardous mismatch bookings; builds contractor confidence.
   - *Cons:* Requires rigorous specification data capture during equipment registration.
2. **Option B:** Textual disclaimer shifting compatibility responsibility entirely to the contractor
   - *Pros:* Zero backend validation logic.
   - *Cons:* High dispute rates and mechanical breakdown claims when mismatched equipment arrives on site.
3. **Option C:** Owner manual approval of equipment pairing
   - *Pros:* Relies on owner mechanical knowledge.
   - *Cons:* Slow booking confirmation turn-around and inconsistent standards.

### 3. Chosen Solution & Technical Rationale
Option A: Automated compatibility validation engine. The API evaluates compatibility rules (weight class tolerance +/- 10%, pin diameter matching in mm, minimum hydraulic flow L/min) before permitting an attachment to be added to a host machine booking.

### 4. Rejected Alternatives
- Option B: Rejected because preventable mechanical failures generate costly disputes and damage the marketplace reputation.
- Option C: Rejected due to human error and delay in reservation confirmation.

### 5. Consequences & Invariants
- **Invariants Maintained:** An attachment can only be bundled with a host machine if its weight_class_min_tons <= host.weight_class_tons <= weight_class_max_tons.
- **Operational Trade-offs:** Owners must accurately measure and input pin diameter and hydraulic flow specs during onboarding.

### 6. Verification
- **Test File:** `backend/internal/modules/machinery/domain/compatibility_test.go`
- **Assertion:** TestCompatibility_RejectsUndersizedHost asserts that pairing a 20-ton hydraulic breaker with an 8-ton excavator fails compatibility checks.

### 7. Revisit Trigger
- Revisit if universal quick-hitch adapter standards eliminate pin diameter constraints.

---

