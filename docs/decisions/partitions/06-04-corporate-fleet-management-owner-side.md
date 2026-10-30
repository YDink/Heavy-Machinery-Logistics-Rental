# Partition Decisions: 06-04 — Corporate Fleet Management & Depot Yard Dispatch Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `06-04`.

- **Governing System ADRs:** `ADR-003`, `ADR-006`
- **Subsystem Scope:** backend, mobile
- **Requirement Targets:** CORP-FLT-01-BE to CORP-FLT-10-BE

---

## D-06-04-001: Enterprise Depot Yard Hierarchy with Delegated Fleet Dispatcher Roles

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [CORP-FLT-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Large equipment rental companies operate multiple equipment yards across regional industrial centers (Addis Ababa Kality, Hawassa, Dire Dawa). Yard managers need autonomy to dispatch machines and conduct inspections in their specific yard without accessing financial payout accounts.

### 2. Considered Alternatives
1. **Option A:** Depot yard entity model with organization_fleet_managers linking dispatchers to specific geographic yards with restricted dispatch/inspection permissions
   - *Pros:* Operational delegation; yard managers cannot access bank accounts or executive financials; geographic asset scoping.
   - *Cons:* Requires yard-level permission filtering in API middleware.
2. **Option B:** Single shared fleet owner account credentials given to all yard managers
   - *Pros:* Zero new permission code.
   - *Cons:* Severe security risk; yard managers can view financial payouts, alter bank accounts, or delete assets.
3. **Option C:** Independent owner accounts for each yard
   - *Pros:* Isolated accounts.
   - *Cons:* Enterprise executive leadership loses consolidated fleet analytics and unified tax invoicing.

### 3. Chosen Solution & Technical Rationale
Option A: Depot yard entity model with delegated dispatcher roles. Fleet assets belong to an organization and are assigned to a primary yard. Dispatchers are granted yard-scoped access, allowing them to manage machine readiness and handovers without seeing company financial ledgers.

### 4. Rejected Alternatives
- Option B: Rejected because credential sharing violates OWASP ASVS Level 3 security and invites internal theft.
- Option C: Rejected because enterprise owners require unified financial consolidation across all operational yards.

### 5. Consequences & Invariants
- **Invariants Maintained:** A yard dispatcher can only initiate handover inspections for assets assigned to their authorized depot yard.
- **Operational Trade-offs:** Machine relocation between yards requires an administrative transfer transaction.

### 6. Verification
- **Test File:** `backend/internal/modules/machinery/application/yard_dispatch_test.go`
- **Assertion:** TestYardDispatch_RestrictsAccessToAssignedYard asserts Kality dispatcher cannot conduct inspections for Dire Dawa machines.

### 7. Revisit Trigger
- Revisit if automated RFID gate transponders automate yard check-in/check-out.

---

## D-06-04-002: Multi-Asset Batch Availability Scheduling and Maintenance Blackout Calendars

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [CORP-FLT-05-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Owners managing 50+ machines need to update availability calendars, schedule periodic 250-hour hydraulic oil servicing, and set maintenance blackout dates in bulk without manually editing each listing one by one.

### 2. Considered Alternatives
1. **Option A:** Batch availability endpoints accepting arrays of fleet IDs and date ranges with atomic calendar conflict validation
   - *Pros:* Fast operational updates; prevents maintenance scheduling during active rentals; bulk efficiency.
   - *Cons:* Requires batched locking across multiple fleet availability records.
2. **Option B:** Editing availability individually on a per-machine screen
   - *Pros:* Simple existing API.
   - *Cons:* Extremely tedious for corporate fleet managers; takes hours to schedule routine maintenance across a fleet.
3. **Option C:** Allowing maintenance blackouts to overwrite and cancel active customer bookings automatically
   - *Pros:* Guarantees machine servicing.
   - *Cons:* Disastrous customer experience: cancels customer rentals on short notice, incurring severe platform penalties.

### 3. Chosen Solution & Technical Rationale
Option A: Batch availability scheduling with conflict validation. Fleet managers select multiple machines and apply maintenance blackout ranges. If any selected machine has an active or confirmed booking during that window, the transaction reports the specific conflict without cancelling rentals.

### 4. Rejected Alternatives
- Option B: Rejected because corporate fleet operations require batch tooling to operate efficiently.
- Option C: Rejected because active customer reservations must never be unilaterally cancelled for routine maintenance.

### 5. Consequences & Invariants
- **Invariants Maintained:** A maintenance blackout cannot be scheduled over dates overlapping an existing CONFIRMED or ACTIVE rental.
- **Operational Trade-offs:** Batch operations must acquire table locks carefully to prevent deadlocks under concurrent bookings.

### 6. Verification
- **Test File:** `backend/internal/modules/machinery/application/batch_availability_test.go`
- **Assertion:** TestBatchAvailability_RejectsWhenBookingExists asserts batch blackout fails for machines with active rentals.

### 7. Revisit Trigger
- Revisit if telematics-driven predictive maintenance automatically reserves service blackouts based on engine hours.

---

