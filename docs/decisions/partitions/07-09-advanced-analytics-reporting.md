# Partition Decisions: 07-09 — Phase 3 Machine Learning Asset Valuation & Predictive Maintenance Analytics

This document records localized Level 2 (L2) architectural decisions governing partition `07-09`.

- **Governing System ADRs:** `ADR-004`, `ADR-008`
- **Subsystem Scope:** backend
- **Requirement Targets:** ANL-01-BE to ANL-10-BE

---

## D-07-09-001: Machine Learning Asset Residual Value and Depreciation Forecasting Model

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [ANL-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Fleet owners require accurate asset valuation and depreciation forecasts to optimize equipment lifecycle retirement, fleet expansion, and bank collateral financing based on empirical marketplace rental demand and operating hours.

### 2. Considered Alternatives
1. **Option A:** Machine learning regression model (XGBoost) trained on historical equipment transactions, brand reliability indices, operating hours, and age to output dynamic asset valuation curves
   - *Pros:* Empirical asset valuation; empowers owners to obtain asset-backed bank financing; automated depreciation tracking.
   - *Cons:* Requires sufficient historical transaction volume to maintain statistical accuracy.
2. **Option B:** Linear straight-line accounting depreciation (10% per year) regardless of hours or condition
   - *Pros:* Simple arithmetic calculation.
   - *Cons:* Completely inaccurate: a machine working 3,000 hours/year in a rock quarry depreciates far faster than one parked in a yard.
3. **Option C:** Zero asset valuation analytics on platform
   - *Pros:* Zero ML overhead.
   - *Cons:* Fails to differentiate platform from basic classified listings.

### 3. Chosen Solution & Technical Rationale
Option A: Machine learning asset valuation model. The analytics engine combines equipment metadata (brand, model, age, cumulative engine hours, inspection damage history) with secondary market price indices to compute estimated market replacement value and projected 36-month depreciation curves.

### 4. Rejected Alternatives
- Option B: Rejected because straight-line depreciation ignores real mechanical wear and tear.
- Option C: Rejected because predictive asset intelligence is a core competitive moat for enterprise fleet owners.

### 5. Consequences & Invariants
- **Invariants Maintained:** Valuation estimates cite confidence intervals and training data cohort sample sizes.
- **Operational Trade-offs:** Valuation curves update on a monthly offline batch training schedule.

### 6. Verification
- **Test File:** `backend/internal/modules/analytics/domain/asset_valuation_test.go`
- **Assertion:** TestAssetValuation_HigherHoursReducesResidualValue asserts identical machine with 5,000h values lower than 1,000h unit.

### 7. Revisit Trigger
- Revisit if commercial banks integrate Dink asset valuation APIs directly for machinery collateral underwriting.

---

## D-07-09-002: Predictive Maintenance Telemetry Analysis with Automated Component Wear Scorecards

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [ANL-06-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Catastrophic mechanical component failures (hydraulic pump blowouts, engine piston seizures) during active rentals cause dangerous job site downtime, costly emergency repairs, and contentious rental dispute claims.

### 2. Considered Alternatives
1. **Option A:** Predictive maintenance engine analyzing operating hour intervals, hydraulic duty cycles, and historic incident logs to compute component wear scorecards with preventative service triggers
   - *Pros:* Prevents catastrophic breakdowns before they happen; lowers total ownership cost; guarantees reliable machinery for renters.
   - *Cons:* Requires tracking component service logs in maintenance module.
2. **Option B:** Purely reactive maintenance: repair only when machine breaks down on site
   - *Pros:* Zero predictive analytics.
   - *Cons:* High job site failure rates; stranded machines; angry contractors; lost rental revenue.
3. **Option C:** Calendar-based servicing recommendations (service every 6 months)
   - *Pros:* Simple time tracking.
   - *Cons:* Fails because machines operate at wildly different duty cycles regardless of elapsed calendar time.

### 3. Chosen Solution & Technical Rationale
Option A: Predictive component wear engine. The system calculates cumulative operating hours against manufacturer component service intervals (e.g. hydraulic filter change every 500h, engine oil every 250h). As components approach 90% wear, automated service alerts are dispatched to the fleet owner.

### 4. Rejected Alternatives
- Option B: Rejected because on-site breakdowns destroy marketplace reliability and contractor trust.
- Option C: Rejected because engine operating hours, not calendar time, drive mechanical wear.

### 5. Consequences & Invariants
- **Invariants Maintained:** Machines with overdue critical safety component servicing are automatically flagged with a maintenance warning in catalog.
- **Operational Trade-offs:** Owners must log completed maintenance service records to reset component wear scorecards.

### 6. Verification
- **Test File:** `backend/internal/modules/machinery/domain/predictive_maintenance_test.go`
- **Assertion:** TestPredictiveMaintenance_AlertsAtThreshold asserts component exceeding 90% operational threshold triggers service alert.

### 7. Revisit Trigger
- Revisit when Phase 3 real-time CAN-bus diagnostic trouble codes (DTCs) automate wear detection.

---

