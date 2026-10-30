# Partition Decisions: 06-05 — Corporate Utilization, ESG & Financial Analytics Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `06-05`.

- **Governing System ADRs:** `ADR-004`, `ADR-008`
- **Subsystem Scope:** backend, admin
- **Requirement Targets:** CORP-REP-01-BE to CORP-REP-08-BE

---

## D-06-05-001: Asynchronous Data Warehouse Projection Pipeline for Heavy Machine Utilization Analytics

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [CORP-REP-01-BE]
- **Subsystem:** [backend, database]

### 1. Problem Statement & Context
Corporate enterprise clients require deep monthly analytics: fleet utilization rates (operating hours vs standby hours), cost per cubic meter excavated, downtime incident frequencies, and project budget variances. Running analytical aggregations directly on the operational OLTP database degrades transactional booking performance.

### 2. Considered Alternatives
1. **Option A:** Asynchronous projection pipeline populating read-optimized corporate_reports reporting tables via background workers using transactional outbox events
   - *Pros:* Zero query impact on primary OLTP transactional tables; sub-second report generation; pre-aggregated metrics.
   - *Cons:* Reporting metrics have a minor asynchronous lag (~5-10 minutes).
2. **Option B:** Heavy dynamic SQL aggregate queries joining 10+ operational tables in real-time
   - *Pros:* Immediate real-time data.
   - *Cons:* Severe database CPU spikes and table lock contention, causing API timeouts for mobile users.
3. **Option C:** Exporting raw CSV dumps and requiring corporate accountants to build Excel pivot tables
   - *Pros:* Zero backend reporting architecture.
   - *Cons:* Fails enterprise client expectations; lack of automated executive dashboards.

### 3. Chosen Solution & Technical Rationale
Option A: Asynchronous projection pipeline. Background workers consume operational completion events and update aggregated reporting summary tables. The Next.js corporate portal reads directly from pre-computed reporting tables without locking operational booking records.

### 4. Rejected Alternatives
- Option B: Rejected because running analytical OLAP queries on operational OLTP databases kills transactional booking performance.
- Option C: Rejected because enterprise clients demand interactive web analytics dashboards.

### 5. Consequences & Invariants
- **Invariants Maintained:** Analytical projection workers must not block or hold locks on operational booking or payment tables.
- **Operational Trade-offs:** Utilization reports reflect data up to the last 15-minute aggregation checkpoint.

### 6. Verification
- **Test File:** `backend/internal/modules/corporate/infrastructure/utilization_projection_test.go`
- **Assertion:** TestUtilization_PreComputesAggregates asserts background worker computes monthly operating hours correctly.

### 7. Revisit Trigger
- Revisit if transaction scale warrants streaming CDC (Change Data Capture) into ClickHouse or DuckDB.

---

## D-06-05-002: Automated Monthly Environmental, Social and Governance (ESG) Fuel Consumption Exports

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [CORP-REP-05-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Multinational infrastructure contractors funded by international development banks (World Bank, AfDB, IFC) must submit mandatory monthly ESG emissions and fuel consumption reports. Manually calculating carbon footprints from paper operator logs is error-prone and labor-intensive.

### 2. Considered Alternatives
1. **Option A:** Automated ESG calculation engine computing estimated diesel consumption and CO2 emissions based on engine tier, operating hours, and equipment category
   - *Pros:* Automated compliance reporting; standard carbon conversion factors (kg CO2 / hour); instant PDF/CSV export for compliance auditors.
   - *Cons:* Uses standardized EPA/Tier-3 emission factor estimates when direct fuel sensor telematics are unavailable.
2. **Option B:** Requiring contractors to manually enter fuel receipt litres on every return
   - *Pros:* Actual fuel receipt capture.
   - *Cons:* Field operators lose fuel receipts or submit fraudulent receipts, producing corrupted audit reports.
3. **Option C:** Zero ESG reporting support
   - *Pros:* Zero engineering effort.
   - *Cons:* Disqualifies the platform from major donor-funded infrastructure development projects.

### 3. Chosen Solution & Technical Rationale
Option A: Automated ESG calculation engine. The system multiplies verified machine operating hours by manufacturer specific fuel consumption factors and standard diesel carbon emission coefficients (2.68 kg CO2 per liter), generating audit-ready PDF/CSV ESG reports.

### 4. Rejected Alternatives
- Option B: Rejected due to human error and incomplete fuel receipt tracking.
- Option C: Rejected because ESG reporting is mandatory for tier-1 civil engineering projects.

### 5. Consequences & Invariants
- **Invariants Maintained:** ESG calculations cite the specific emission factor standards applied in report footnotes.
- **Operational Trade-offs:** Estimates must be clearly disclaimed as calculated approximations when direct telematics sensors are absent.

### 6. Verification
- **Test File:** `backend/internal/modules/corporate/domain/esg_calculation_test.go`
- **Assertion:** TestESG_CalculatesAccurateEmissions asserts 100 hours of 20-ton excavator produces expected CO2 metric tonnes within tolerance.

### 7. Revisit Trigger
- Revisit when Phase 3 telematics integrate direct electronic fuel level sensor telemetry.

---

