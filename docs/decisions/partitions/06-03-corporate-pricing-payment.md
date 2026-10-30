# Partition Decisions: 06-03 — Corporate Negotiated Pricing, NET-30 Credit & Master Agreements

This document records localized Level 2 (L2) architectural decisions governing partition `06-03`.

- **Governing System ADRs:** `ADR-002`, `ADR-010`, `ADR-015`
- **Subsystem Scope:** backend, database
- **Requirement Targets:** CORP-PRC-01-BE to CORP-PRC-10-BE

---

## D-06-03-001: Negotiated Volume Rate Cards and Custom Contract Tariff Overrides

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [CORP-PRC-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Large enterprise contractors negotiate bespoke volume discount tariffs with equipment owners (e.g. 15% discount for long-term rentals > 6 months). The platform must calculate custom pricing rules for specific corporate agreements without polluting public catalog prices.

### 2. Considered Alternatives
1. **Option A:** Corporate pricing agreements entity (corporate_pricing_agreements) evaluated during pricing calculation with precedence over base fleet tariffs
   - *Pros:* Automated negotiated rate application; private pricing isolated to authorized corporate members; full auditability.
   - *Cons:* Pricing calculation pipeline must check for active corporate agreements.
2. **Option B:** Manual price adjustment coupons issued by customer support
   - *Pros:* Zero custom pricing database tables.
   - *Cons:* High administrative friction; prone to coupon leakage and human calculation errors.
3. **Option C:** Duplicate private fleet listings created specifically for the corporate client
   - *Pros:* Uses existing listing models.
   - *Cons:* Clutters database with duplicate machine assets and creates availability synchronization nightmares.

### 3. Chosen Solution & Technical Rationale
Option A: Corporate pricing agreements entity. When an authenticated corporate user views a listing, the pricing engine checks for active agreements linking organization_id and owner_id/fleet_id, applying agreed rate overrides transparently.

### 4. Rejected Alternatives
- Option B: Rejected because manual coupons do not scale to multi-machine enterprise contracts.
- Option C: Rejected because duplicating fleet records destroys single-source-of-truth machine telemetry and availability.

### 5. Consequences & Invariants
- **Invariants Maintained:** Negotiated rates only apply to bookings created within the valid date range of an active corporate pricing agreement.
- **Operational Trade-offs:** Pricing engine must execute agreement lookups in Valkey cache to maintain sub-5ms quote performance.

### 6. Verification
- **Test File:** `backend/internal/modules/pricing/application/corporate_rate_test.go`
- **Assertion:** TestCorporateRate_AppliesNegotiatedDiscount asserts corporate user receives contracted rate rather than public base rate.

### 7. Revisit Trigger
- Revisit if dynamic commodity index fuel price adjustments are bundled into long-term machine rate cards.

---

## D-06-03-002: Consolidated NET-30 Corporate Invoicing and Credit Line Underwriting Engine

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [CORP-PRC-06-BE]
- **Subsystem:** [backend, database]

### 1. Problem Statement & Context
Enterprise contractors refuse to execute individual manual bank transfers for hundreds of daily machine shifts, demanding consolidated monthly invoicing with 30-day payment terms (NET-30). Offering credit lines without strict underwriting exposes the marketplace to severe credit default risk.

### 2. Considered Alternatives
1. **Option A:** NET-30 credit line underwriting: corporate organizations receive an approved credit_limit_minor based on audited bank guarantees and financials; monthly consolidated invoices are generated with auto-suspension on overdue aging (>45 days)
   - *Pros:* Meets enterprise procurement requirements; strictly bounded credit exposure; automated billing consolidation; automated account suspension on default.
   - *Cons:* Requires platform finance controllers to manually underwrite and set credit limits.
2. **Option B:** Unrestricted post-paid billing for all corporate accounts
   - *Pros:* Maximum customer onboarding velocity.
   - *Cons:* Guaranteed bankruptcy from massive uncollectible debt defaults.
3. **Option C:** Reject NET-30 entirely and force enterprise clients to prepay every rental
   - *Pros:* Zero credit risk.
   - *Cons:* Complete failure to capture enterprise contractor market; tier-1 construction firms will not use the platform.

### 3. Chosen Solution & Technical Rationale
Option A: Underwritten NET-30 credit lines. Organizations are underwritten by finance controllers with an approved credit limit. Bookings draw down available credit without upfront bank slips. On the 1st of each month, a consolidated tax invoice is generated with 30-day settlement terms.

### 4. Rejected Alternatives
- Option B: Rejected due to catastrophic commercial financial default risk.
- Option C: Rejected because enterprise commercial contractors operate exclusively on credit invoicing cycles.

### 5. Consequences & Invariants
- **Invariants Maintained:** Total outstanding unbilled usage + unpaid invoices cannot exceed organization approved credit_limit_minor.
- **Operational Trade-offs:** Requires legal collection processes and bank guarantee verification by Dink finance staff.

### 6. Verification
- **Test File:** `backend/internal/modules/corporate/application/credit_limit_test.go`
- **Assertion:** TestCreditLimit_BlocksBookingWhenCeilingExceeded asserts booking fails when drawdown breaches approved credit line.

### 7. Revisit Trigger
- Revisit if trade credit insurance providers underwrite platform receivables in Phase 3.

---

