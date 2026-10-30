# Partition Decisions: 05-08 — Logistics Transport & Lowboy Carrier Integration Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-08`.

- **Governing System ADRs:** `ADR-001`, `ADR-005`, `ADR-007`
- **Subsystem Scope:** backend, mobile
- **Requirement Targets:** LOG-01-BE to LOG-12-BE

---

## D-05-08-001: Lowboy Carrier Linked Booking with Dual-Leg Transition Coordination

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [LOG-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Heavy machinery (excavators, bulldozers) cannot drive on paved highways and requires specialized lowboy trailer transport to and from project sites. Transport sub-bookings must coordinate seamlessly with the primary machinery rental schedule without desynchronization.

### 2. Considered Alternatives
1. **Option A:** Linked sub-booking entity (parent_booking_id) with atomic state machine synchronization across delivery and return legs
   - *Pros:* Machine and transport schedules stay synchronized; cancellation of machine auto-triggers transport review; single checkout.
   - *Cons:* Increases state machine complexity with linked parent-child transitions.
2. **Option B:** Completely unlinked standalone transport bookings made separately by the renter
   - *Pros:* Simple independent booking models.
   - *Cons:* High risk of machine arriving without transport or transport arriving when machine is not ready; contractor frustration.
3. **Option C:** Machinery owners bundle transport into their rental rate and manage lowboys offline
   - *Pros:* Zero logistics software required.
   - *Cons:* Restricts rental radius to owners who own their own lowboy carriers, cutting out 70% of available fleet inventory.

### 3. Chosen Solution & Technical Rationale
Option A: Linked sub-booking entity with parent_booking_id. When a renter selects transport, a linked logistics reservation is created. The primary machine cannot enter ACTIVE until the delivery leg inspection is marked DELIVERED_ON_SITE.

### 4. Rejected Alternatives
- Option B: Rejected because asynchronous uncoordinated transport bookings cause massive site mobilization delays.
- Option C: Rejected because the majority of equipment owners do not operate lowboy transport trailers.

### 5. Consequences & Invariants
- **Invariants Maintained:** A machinery booking with required transport cannot transition to ACTIVE until the transport sub-booking completes delivery inspection.
- **Operational Trade-offs:** Transport carrier delays directly impact machinery rental start times.

### 6. Verification
- **Test File:** `backend/internal/modules/booking/application/logistics_coordination_test.go`
- **Assertion:** TestLogistics_ParentBookingBlockedUntilTransportDelivery asserts machinery cannot be marked ACTIVE while transport is IN_TRANSIT.

### 7. Revisit Trigger
- Revisit when Phase 3 autonomous freight matching (`07-03`) introduces automated spot-market carrier bidding.

---

## D-05-08-002: Axle Load and Gross Vehicle Weight (GVM) Safety Enforcement Against Regional Limits

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [LOG-04-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Overloading transport carriers violates East African Community (EAC) vehicle load limits (56-tonne gross vehicle mass cap; axle limits) and results in severe highway authority fines, vehicle impoundment, and structural bridge collapse risks.

### 2. Considered Alternatives
1. **Option A:** Automated mechanical weight matching: machine operating weight + trailer tare weight must be <= carrier payload rating and jurisdictional road ceilings
   - *Pros:* Prevents illegal and dangerous carrier overloading; automatic rejection of undersized trailers.
   - *Cons:* Requires accurate tare and payload specifications during fleet and carrier registration.
2. **Option B:** Disclaimer leaving road compliance entirely to the transporter driver
   - *Pros:* Zero backend validation rules.
   - *Cons:* Impounded machinery causes massive project delays, insurance cancellations, and legal liability.
3. **Option C:** Manual review of every transport pairing by operations dispatchers
   - *Pros:* Human oversight of heavy haul permits.
   - *Cons:* Slow booking confirmation that fails during weekend or after-hours emergency mobilizations.

### 3. Chosen Solution & Technical Rationale
Option A: Automated mechanical weight matching. The logistics dispatch engine computes: total_mass = machine.operating_weight_kg + carrier.tare_weight_kg. If total_mass exceeds carrier.max_payload_kg or national highway limits, pairing is blocked with an RFC 9457 error.

### 4. Rejected Alternatives
- Option B: Rejected due to catastrophic legal, insurance, and road safety liability.
- Option C: Rejected because manual checks create unacceptable booking confirmation bottlenecks.

### 5. Consequences & Invariants
- **Invariants Maintained:** Carrier assignment enforces: machine.operating_weight_kg <= carrier.max_payload_kg AND total_gvm <= country_configs.max_gvm_kg.
- **Operational Trade-offs:** Owners must provide certified weighbridge or manufacturer specification documentation for heavy haul trailers.

### 6. Verification
- **Test File:** `backend/internal/modules/machinery/domain/gvm_validation_test.go`
- **Assertion:** TestGVM_RejectsOverloadedTrailer asserts pairing a 45-tonne excavator with a 30-tonne capacity lowboy fails validation.

### 7. Revisit Trigger
- Revisit if national transport authorities issue special multi-axle heavy transport permit integrations in Phase 3.

---

