# Partition Decisions: 05-23 — Ethiopian Ge'ez Calendar & Public Holiday Booking Engine Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-23`.

- **Governing System ADRs:** `ADR-013`, `ADR-015`
- **Subsystem Scope:** backend, mobile
- **Requirement Targets:** CAL-01-BE to CAL-10-BE

---

## D-05-23-001: Multi-Country Operational Parameters and Ge'ez/Gregorian Calendar Resolution via country_configs and X-Country-Code Middleware

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [CAL-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Ethiopia operates on the 13-month Ge'ez calendar (7-8 years behind Gregorian calendar; Meskerem 1 starts in September; 12 months of 30 days plus Pagume of 5 or 6 days). Kenya operates on the Gregorian calendar. Heavy machinery rental contracts, delivery schedules, and dispute arbitration deadlines must calculate and display dates accurately in both calendar systems without temporal corruption.

### 2. Considered Alternatives
1. **Option A:** Universal UTC timestamp persistence in database combined with dynamic calendar resolution via country_configs and X-Country-Code middleware with client-side dual-calendar conversion engine
   - *Pros:* Zero database date math bugs; pure ISO 8601 UTC wire serialization; localized calendar rendering matching user national jurisdiction.
   - *Cons:* Requires mathematical conversion algorithms between Julian Day Numbers, Ge'ez, and Gregorian calendars.
2. **Option B:** Storing local calendar dates as strings in database (e.g. "2018-01-01" in Ge'ez)
   - *Pros:* Direct representation of local dates.
   - *Cons:* Disastrous: database sorting, duration arithmetic, and cross-border comparisons become impossible.
3. **Option C:** Force all users in Ethiopia and Kenya to use standard Gregorian calendar only
   - *Pros:* Zero dual-calendar code needed.
   - *Cons:* Contractors and government inspectors in Ethiopia cannot verify legal lease contracts; severe market rejection.

### 3. Chosen Solution & Technical Rationale
Option A: Universal UTC timestamp persistence with dynamic calendar resolution. PostgreSQL persists all dates as TIMESTAMPTZ in UTC. The mobile and web frontends convert UTC timestamps to Ge'ez dates (e.g. 11 Meskerem 2018) when active country is ET, using an algorithmic Julian Day Number (JDN) converter, while wire contracts remain strict ISO 8601.

### 4. Rejected Alternatives
- Option B: Rejected because string date storage destroys date math, range indexing, and cross-border accounting.
- Option C: Rejected because Ethiopian business, banking, and legal operations run fundamentally on the Ge'ez calendar.

### 5. Consequences & Invariants
- **Invariants Maintained:** All timestamps in database and wire APIs are strictly ISO 8601 UTC. Ge'ez conversion occurs exclusively at the presentation boundary.
- **Operational Trade-offs:** Date pickers in mobile app must support toggling between Ge'ez and Gregorian views in Ethiopia.

### 6. Verification
- **Test File:** `backend/internal/platform/calendar/geez_converter_test.go`
- **Assertion:** TestGeezConverter_AccuratelyConvertsEnkutatash asserts 2026-09-11 UTC maps precisely to 1 Meskerem 2019 in Ge'ez.

### 7. Revisit Trigger
- Revisit if national legal statutes standardize digital government contracts on a single international calendar.

---

## D-05-23-002: Algorithmic Fasika Computus Calculation vs Admin Seeded Public Holiday Table Overrides

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [CAL-05-BE]
- **Subsystem:** [backend, database]

### 1. Problem Statement & Context
Public holidays affect rental rate calculations (holiday surcharges or rental grace extensions) and pause statutory dispute resolution SLA timers. In Ethiopia, Christian Orthodox movable feasts (Fasika, Siklet) follow the Alexandrian computus, Muslim movable feasts (Eid al-Fitr, Eid al-Adha) follow lunar observations, and national holidays (Adwa Victory, Patriots Day) are fixed.

### 2. Considered Alternatives
1. **Option A:** Hybrid calendar engine: algorithmic Alexandrian computus for Christian feasts combined with admin-managed public_holidays relational table for lunar and government-decreed holiday overrides
   - *Pros:* Automated calculation of movable Easter feasts; admin flexibility to declare emergency or lunar public holidays instantly; SLA timers automatically pause.
   - *Cons:* Requires maintaining public_holidays seed data annually.
2. **Option B:** Pure hardcoded holiday dates in Go application code
   - *Pros:* Zero database lookups.
   - *Cons:* Lunar holidays shift annually; government-decreed unscheduled holidays cannot be added without a backend deployment.
3. **Option C:** Ignoring public holidays in dispute and rental rate calculations
   - *Pros:* Zero holiday logic.
   - *Cons:* Dispute SLA timers expire over national holiday weekends, penalizing contractors when banks and offices are closed.

### 3. Chosen Solution & Technical Rationale
Option A: Hybrid calendar engine. Fixed and computus-based holidays are calculated algorithmically, while the public_holidays database table allows administrators to seed and adjust national holidays dynamically. The dispute engine skips public_holidays when computing 72-hour arbitration deadlines.

### 4. Rejected Alternatives
- Option B: Rejected because hardcoded dates cannot accommodate lunar holiday shifts or sudden government public holidays.
- Option C: Rejected because expiring dispute deadlines during national holidays violates procedural fairness.

### 5. Consequences & Invariants
- **Invariants Maintained:** Dispute arbitration SLA timers automatically exclude verified public holidays for the active jurisdiction.
- **Operational Trade-offs:** Administrators must confirm lunar holiday dates each year following official government proclamations.

### 6. Verification
- **Test File:** `backend/internal/platform/calendar/holiday_sla_test.go`
- **Assertion:** TestHolidaySLA_PausesDuringAdwaVictory asserts that a dispute opened before Adwa Victory does not decrement SLA hours during the holiday.

### 7. Revisit Trigger
- Revisit if regional holiday variations between regional states (e.g. Oromia Irreecha) require regional state calendar partitioning.

---

