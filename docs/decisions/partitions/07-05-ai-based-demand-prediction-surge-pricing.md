# Partition Decisions: 07-05 — Phase 3 AI Demand Forecasting & Dynamic Surge Pricing Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `07-05`.

- **Governing System ADRs:** `ADR-002`, `ADR-005`
- **Subsystem Scope:** backend
- **Requirement Targets:** AI-PRC-01-BE to AI-PRC-08-BE

---

## D-07-05-001: Spatial Machine Learning Demand Forecasting on Hexagonal Hierarchical Spatial Index (H3)

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [AI-PRC-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Heavy machinery demand surges unpredictably during seasonal construction seasons (e.g. dry season post-Meskerem in Ethiopia) and major infrastructure tenders. Owners price equipment sub-optimally, leading to regional shortages in high-demand corridors.

### 2. Considered Alternatives
1. **Option A:** Uber H3 hexagonal spatial indexing (resolution 7) aggregating historical booking density, seasonal weather patterns, and public tender awards to forecast regional supply deficits
   - *Pros:* Predicts regional equipment shortages 14 days in advance; guides fleet repositioning; spatial indexing efficiency.
   - *Cons:* Requires training and maintaining gradient-boosted forecasting models.
2. **Option B:** Static seasonal multiplier tables edited manually by administrators
   - *Pros:* Simple static lookup.
   - *Cons:* Cannot capture localized demand spikes (e.g. sudden highway bridge reconstruction tender in Awash).
3. **Option C:** Zero forecasting; purely reactive spot pricing
   - *Pros:* Zero AI complexity.
   - *Cons:* Severe equipment bottlenecks and missed revenue optimization for fleet owners.

### 3. Chosen Solution & Technical Rationale
Option A: H3 spatial indexing demand forecasting. The platform aggregates search frequency, booking velocity, and seasonal weather by H3 hex cell, training an offline LightGBM model to predict 14-day forward demand indices per equipment category.

### 4. Rejected Alternatives
- Option B: Rejected because manual tables cannot react to dynamic regional construction tenders.
- Option C: Rejected because predictive intelligence provides fleet owners with compelling monetization value.

### 5. Consequences & Invariants
- **Invariants Maintained:** Demand prediction scores are advisory and never alter owner base listing prices without explicit owner opt-in.
- **Operational Trade-offs:** Model training requires periodic batch retraining on weekly historical logs.

### 6. Verification
- **Test File:** `backend/internal/modules/analytics/application/demand_forecast_test.go`
- **Assertion:** TestDemandForecast_ComputesH3Index asserts valid H3 index generation from geographic coordinates.

### 7. Revisit Trigger
- Revisit if real-time satellite imagery of construction earthworks becomes commercially viable.

---

## D-07-05-002: Algorithmic Surge Pricing Multiplier with Regulatory Ceiling Guardrails

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [AI-PRC-05-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
During extreme equipment shortages, unconstrained surge pricing can lead to contractor price gouging and regulatory antitrust scrutiny. Dynamic pricing must balance owner revenue optimization with fair commercial price ceilings.

### 2. Considered Alternatives
1. **Option A:** Dynamic surge pricing multiplier bounded strictly between 0.85x (discount floor) and 1.35x (ceiling cap) with explicit contractor confirmation
   - *Pros:* Prevents price gouging; complies with fair trade practices; incentivizes owners to mobilize idle machines to high-demand corridors.
   - *Cons:* Caps maximum surge revenue during extreme emergency shortages.
2. **Option B:** Uncapped algorithmic surge pricing (e.g. 3.0x to 5.0x during peak demand)
   - *Pros:* Maximizes short-term revenue.
   - *Cons:* Triggers severe contractor backlash, public relations scandals, and government price control intervention.
3. **Option C:** Completely flat fixed rates with zero dynamic adjustment
   - *Pros:* Total price predictability.
   - *Cons:* Owners refuse to list machines during peak periods, preferring offline cash black markets.

### 3. Chosen Solution & Technical Rationale
Option A: Bounded surge pricing engine. The dynamic pricing multiplier is clamped between 0.85x and 1.35x based on local H3 cell supply-demand ratio. The surge multiplier is displayed transparently to the renter before booking confirmation.

### 4. Rejected Alternatives
- Option B: Rejected because predatory price gouging destroys contractor relationships and invites regulatory sanctions.
- Option C: Rejected because rigid rates lead to artificial supply shortages and platform circumvention.

### 5. Consequences & Invariants
- **Invariants Maintained:** Surge pricing multiplier cannot exceed the regulatory ceiling of 1.35x under any operational condition.
- **Operational Trade-offs:** Owners in high-demand zones cannot charge arbitrary surge premiums above the platform ceiling.

### 6. Verification
- **Test File:** `backend/internal/modules/pricing/domain/surge_multiplier_test.go`
- **Assertion:** TestSurgeMultiplier_ClampedAtCeiling asserts extreme demand ratio never produces a multiplier > 1.35x.

### 7. Revisit Trigger
- Revisit if national competition authorities mandate narrower pricing bandwidths.

---

