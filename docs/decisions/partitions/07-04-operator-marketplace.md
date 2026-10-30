# Partition Decisions: 07-04 — Phase 3 Independent Certified Operator Marketplace Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `07-04`.

- **Governing System ADRs:** `ADR-006`, `ADR-011`
- **Subsystem Scope:** backend, mobile
- **Requirement Targets:** OPR-MKT-01-BE to OPR-MKT-10-BE

---

## D-07-04-001: Independent Heavy Machinery Operator Registry with Biometric Credential Attestation

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [OPR-MKT-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Contractors who own equipment or rent bare-metal machinery frequently lack certified heavy equipment operators. Inexperienced or uncertified operators operating complex machinery (e.g. 50-tonne mobile cranes) cause fatal accidents and catastrophic equipment destruction.

### 2. Considered Alternatives
1. **Option A:** Verified operator registry with mandatory government competence certificate verification, machinery category endorsements, and biometric Fayda attestation
   - *Pros:* High safety standards; certified skill badges; contractor confidence; liability reduction.
   - *Cons:* Requires compliance review of operator licenses and competency certificates.
2. **Option B:** Open freelance job board with unverified operator profiles
   - *Pros:* Rapid operator onboarding.
   - *Cons:* Dangerous safety liability: fraudulent or untrained individuals hire out on lethal heavy machinery.
3. **Option C:** Only allow equipment owners to provide their own in-house operators
   - *Pros:* Owner assumes operator liability.
   - *Cons:* Prevents bare-metal rentals and excludes thousands of skilled independent operators from earning income.

### 3. Chosen Solution & Technical Rationale
Option A: Verified operator registry with competency certification and biometric attestation. Operators register, select certified machinery categories, upload ministry competency certificates, and undergo compliance verification before appearing in search.

### 4. Rejected Alternatives
- Option B: Rejected due to catastrophic job site safety risks and insurance nullification.
- Option C: Rejected because an independent operator marketplace unlocks high-margin revenue and market liquidity.

### 5. Consequences & Invariants
- **Invariants Maintained:** An operator cannot be dispatched for a machinery category for which they lack an active, verified competence certificate.
- **Operational Trade-offs:** Operators must maintain up-to-date medical and licensing documentation.

### 6. Verification
- **Test File:** `backend/internal/modules/operator/application/certification_test.go`
- **Assertion:** TestOperator_RequiresCategoryEndorsement asserts uncertified operator cannot accept crane booking.

### 7. Revisit Trigger
- Revisit if national operator competency licensing is automated via Ministry of Labor digital portals.

---

## D-07-04-002: Operator On-Demand Booking Dispatch with Dynamic Hourly Wage Escrow Hold

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [OPR-MKT-05-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Operators hired on-demand require guaranteed wage payment upon shift completion, while contractors require protection against operators abandoning job sites or operating equipment negligently.

### 2. Considered Alternatives
1. **Option A:** Escrow wage hold at dispatch: contractor funds operator wage into platform escrow; funds release upon verified shift completion inspection and mutual hour meter sign-off
   - *Pros:* Guaranteed payment for operators; protects contractors against no-shows; eliminates wage payment disputes.
   - *Cons:* Requires escrow holding infrastructure.
2. **Option B:** Direct cash payment from contractor to operator at the end of the day
   - *Pros:* Zero escrow code.
   - *Cons:* Frequent contractor wage theft where operators work 10 hours and are refused payment; platform disintermediation.
3. **Option C:** Monthly post-paid salary billing
   - *Pros:* Standard corporate payroll.
   - *Cons:* Incompatible with gig-economy independent daily operator hiring.

### 3. Chosen Solution & Technical Rationale
Option A: Escrow wage hold at dispatch. When an operator is booked, estimated daily wages are held in escrow. Upon shift completion and dual-party confirmation of operating hours, net wages release to operator wallet while platform fee is debited.

### 4. Rejected Alternatives
- Option B: Rejected because cash payment leads to chronic wage theft and eliminates platform monetization.
- Option C: Rejected because independent operators require daily or weekly wage liquidity.

### 5. Consequences & Invariants
- **Invariants Maintained:** Operator wages held in escrow cannot be released without verified shift completion or dispute ruling.
- **Operational Trade-offs:** Contractors must fund operator wages in advance of shift mobilization.

### 6. Verification
- **Test File:** `backend/internal/modules/operator/application/wage_escrow_test.go`
- **Assertion:** TestWageEscrow_ReleasesOnDualConfirmation asserts escrow releases funds to operator upon verified shift completion.

### 7. Revisit Trigger
- Revisit if operator union collective bargaining agreements mandate standardized wage tables.

---

