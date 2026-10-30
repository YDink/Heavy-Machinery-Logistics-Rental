# Partition Decisions: 07-02 — Phase 3 Automated Gateway Commission Split & Multi-Party Settlement

This document records localized Level 2 (L2) architectural decisions governing partition `07-02`.

- **Governing System ADRs:** `ADR-004`, `ADR-010`
- **Subsystem Scope:** backend
- **Requirement Targets:** SPLIT-01-BE to SPLIT-08-BE

---

## D-07-02-001: Automated Real-Time Gateway Fee Split Settlement at Point of Renter Checkout

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [SPLIT-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
In Phase 3 with full payment gateway integration, customer funds can be split automatically at payment time between platform commission and owner rental payouts, eliminating the need for owner wallet pre-funding.

### 2. Considered Alternatives
1. **Option A:** Automated gateway split settlement: payment gateway divides incoming funds into platform commission account and owner escrow sub-account at checkout
   - *Pros:* Zero manual commission collection; eliminates owner wallet pre-funding requirement; instant revenue capture.
   - *Cons:* Requires payment gateway support for marketplace split payments.
2. **Option B:** Collect 100% into platform account and disburse owner payout manually at end of month
   - *Pros:* Simple single payment inflow.
   - *Cons:* Exposes platform to significant balance-sheet deposit custody liabilities under central bank rules.
3. **Option C:** Continue manual bank transfer model even after payment gateways are integrated
   - *Pros:* Zero change to settlement logic.
   - *Cons:* Leaves high friction in place when modern automated checkout is available.

### 3. Chosen Solution & Technical Rationale
Option A: Automated gateway split settlement. When checkout completes via Telebirr or M-Pesa, the gateway executes a multi-party transfer: platform commission basis points route to Dink corporate account; net rental funds route to owner settlement account.

### 4. Rejected Alternatives
- Option B: Rejected because holding customer deposits violates non-banking regulatory licenses.
- Option C: Rejected because automated split settlement provides superior customer and owner experience.

### 5. Consequences & Invariants
- **Invariants Maintained:** Sum of platform commission split + owner net settlement must equal 100% of the total charged rental amount.
- **Operational Trade-offs:** Requires custom split payment APIs with domestic East African telecom mobile money providers.

### 6. Verification
- **Test File:** `backend/internal/modules/payment/application/split_settlement_test.go`
- **Assertion:** TestSplitSettlement_AccurateBasisPointDivision asserts 1000 bps platform fee splits exactly 10% platform / 90% owner.

### 7. Revisit Trigger
- Revisit if national payment switch provides unified programmatic escrow settlement APIs.

---

## D-07-02-002: Dynamic Platform Fee Basis Point Apportionment with Multi-Party Escrow Release

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [SPLIT-05-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Different equipment categories, promotional campaigns, or long-term leases may apply tiered platform commission rates (e.g. 800 bps for cranes vs 1200 bps for small compaction tools). Fee calculation must be dynamic and tamper-evident.

### 2. Considered Alternatives
1. **Option A:** Commission rule engine evaluating category, lease duration, and owner tier to determine basis points with cryptographic ledger audit stamp
   - *Pros:* Flexible marketplace monetization; supports volume promotions; fully auditable basis point tracking.
   - *Cons:* Fee calculation rules must be evaluated within the booking price calculation pipeline.
2. **Option B:** Hardcoded global 10% commission across all categories and durations
   - *Pros:* Trivial implementation.
   - *Cons:* Inflexible; prevents volume discounts for high-value enterprise machines or promotional owner onboarding.
3. **Option C:** Ad-hoc commission negotiation entered manually by sales reps per booking
   - *Pros:* Custom pricing.
   - *Cons:* Prone to human corruption, kickbacks, and accounting reconciliation errors.

### 3. Chosen Solution & Technical Rationale
Option A: Dynamic commission rule engine. The system resolves commission_rules based on equipment category and duration, computing commission_basis_points and recording the active rule ID in the financial ledger entry.

### 4. Rejected Alternatives
- Option B: Rejected because modern heavy equipment marketplaces require flexible category-specific monetization.
- Option C: Rejected due to severe operational fraud and kickback risks.

### 5. Consequences & Invariants
- **Invariants Maintained:** Commission basis points can never be lower than the jurisdictional floor (e.g. 500 bps) or higher than the ceiling (e.g. 2500 bps).
- **Operational Trade-offs:** Rule engine updates must be versioned to prevent changing rates on existing in-flight bookings.

### 6. Verification
- **Test File:** `backend/internal/modules/payment/domain/commission_rule_test.go`
- **Assertion:** TestCommissionRule_AppliesTieredRate asserts long-term rentals over 30 days receive 800 bps fee instead of standard 1000 bps.

### 7. Revisit Trigger
- Revisit if algorithmic dynamic take-rates based on supply-demand elasticity are tested in Phase 3.

---

