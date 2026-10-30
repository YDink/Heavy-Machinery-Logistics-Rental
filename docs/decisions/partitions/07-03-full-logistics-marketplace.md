# Partition Decisions: 07-03 — Phase 3 Full Logistics Freight Marketplace & Haul Matching Engine

This document records localized Level 2 (L2) architectural decisions governing partition `07-03`.

- **Governing System ADRs:** `ADR-001`, `ADR-006`
- **Subsystem Scope:** backend, mobile
- **Requirement Targets:** FRT-01-BE to FRT-12-BE

---

## D-07-03-001: Automated Heavy Haul Freight Matching Engine with Dynamic Route Feasibility Checks

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [FRT-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Hauling an oversized 45-tonne excavator requires matching with certified multi-axle lowboy trailers and checking route clearances (overhead bridges, toll stations, weight-restricted bridges). Manual carrier sourcing causes days of project mobilization delays.

### 2. Considered Alternatives
1. **Option A:** Algorithmic haul matching engine pairing machinery dimensions/weights with certified lowboy carriers within a 150km radius and checking route corridor restrictions
   - *Pros:* Automated carrier dispatch; verified weight and bridge clearances; reduces empty backhaul miles for truckers.
   - *Cons:* Requires route network graph data and accurate trailer spec profiles.
2. **Option B:** Open bulletin board where truckers manually browse and call contractors
   - *Pros:* Low software complexity.
   - *Cons:* Inefficient; unlicensed or undersized trailers bid on heavy equipment, causing highway accidents and roadside impoundment.
3. **Option C:** Dink Machinery purchases and operates its own fleet of lowboy trucks
   - *Pros:* Total control over logistics.
   - *Cons:* Massive capital expenditure; transforms technology marketplace into an asset-heavy trucking company.

### 3. Chosen Solution & Technical Rationale
Option A: Algorithmic freight matching engine. When machinery rental requires transport, the engine searches lowboy carrier profiles, filters by payload capacity >= machine weight, and presents optimized carrier options with transparent haul tariffs.

### 4. Rejected Alternatives
- Option B: Rejected because open boards lead to severe safety violations and overloaded trucks.
- Option C: Rejected because Dink is an asset-light technology aggregator, not a trucking fleet operator.

### 5. Consequences & Invariants
- **Invariants Maintained:** A carrier cannot be matched with an equipment load exceeding the carrier certified payload capacity.
- **Operational Trade-offs:** Transporter drivers must keep their availability status updated on mobile.

### 6. Verification
- **Test File:** `backend/internal/modules/logistics/application/freight_matching_test.go`
- **Assertion:** TestFreightMatching_FiltersByCapacity asserts lowboy with 30t capacity is excluded from 40t excavator haul requests.

### 7. Revisit Trigger
- Revisit if national highway authority opens real-time bridge weight restriction APIs.

---

## D-07-03-002: Reverse-Auction Lowboy Bidding Architecture with Time-Bounded Tender Windows

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [FRT-05-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Standard fixed-price transport rates fail on remote or difficult access routes. Contractors need competitive haul pricing, while transporters need opportunities to bid on profitable backhauls.

### 2. Considered Alternatives
1. **Option A:** Time-bounded reverse auction (2-hour to 24-hour tender window) where certified transporters submit competitive bids below a maximum ceiling price
   - *Pros:* Lowers mobilization costs for contractors; fills empty carrier backhauls; dynamic market pricing.
   - *Cons:* Requires contractors to wait for bidding window to close.
2. **Option B:** Fixed distance-based tariff (e.g. 150 ETB per km) across all routes
   - *Pros:* Instant booking price.
   - *Cons:* Fails in reality: truckers refuse to drive difficult unpaved or mountainous routes at flat standard rates.
3. **Option C:** Unrestricted price negotiation via phone calls
   - *Pros:* Traditional method.
   - *Cons:* Off-platform disintermediation and total lack of marketplace transparency.

### 3. Chosen Solution & Technical Rationale
Option A: Time-bounded reverse auction. Contractors can post a haul request with a maximum ceiling price. Transporters receive push notifications and submit downward bids. The contractor selects the winning bid or the lowest bid wins automatically upon timer expiry.

### 4. Rejected Alternatives
- Option B: Rejected because flat tariffs do not reflect real-world road terrain and fuel variations.
- Option C: Rejected because offline negotiations defeat the purpose of an automated digital marketplace.

### 5. Consequences & Invariants
- **Invariants Maintained:** All bids in a reverse auction must be <= the contractor maximum ceiling price.
- **Operational Trade-offs:** Emergency transport mobilizations require a "Book Now" fixed-price override option.

### 6. Verification
- **Test File:** `backend/internal/modules/logistics/application/reverse_auction_test.go`
- **Assertion:** TestReverseAuction_BidsMustBeBelowCeiling asserts bid exceeding ceiling price is rejected with HTTP 422.

### 7. Revisit Trigger
- Revisit if algorithmic pricing models become accurate enough to eliminate bidding.

---

