# Partition Decisions: 05-04 — Fleet Discovery, Spatial Search & Comparison Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-04`.

- **Governing System ADRs:** `ADR-004`, `ADR-005`
- **Subsystem Scope:** backend, database
- **Requirement Targets:** SRCH-01-BE to SRCH-12-BE

---

## D-05-04-001: Spatial Radius Discovery via PostGIS ST_DWithin with Keyset Keyset Pagination

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [SRCH-01-BE]
- **Subsystem:** [database]

### 1. Problem Statement & Context
Contractors search for available heavy machinery within specific radii (25km, 50km, 100km) of project sites along regional infrastructure corridors. Deep offset pagination degrades query performance exponentially and produces unstable result feeds under concurrent insertions.

### 2. Considered Alternatives
1. **Option A:** PostGIS ST_DWithin on geography(POINT, 4326) with GiST spatial indexing and keyset cursor pagination (distance_meters, id)
   - *Pros:* Deterministic sub-15ms spatial queries; zero offset pagination degradation; stable scrolling feeds.
   - *Cons:* Cursor pagination requires composite keyset tokens.
2. **Option B:** Bounding box (BBox) rectangle coordinates with standard B-Tree latitude/longitude indexes
   - *Pros:* Simple coordinate range math.
   - *Cons:* Inaccurate spheroidal distance calculations at equatorial latitudes; distortion increases with distance.
3. **Option C:** External Elasticsearch / OpenSearch cluster with geo_point queries
   - *Pros:* Specialized spatial search cluster.
   - *Cons:* Data synchronization drift with primary PostgreSQL; operational complexity in local datacenter.

### 3. Chosen Solution & Technical Rationale
Option A: PostGIS ST_DWithin on geography(POINT, 4326) with GiST spatial indexing paired with keyset cursor pagination. Query results order by spheroidal distance and primary key, returning nextCursor and prevCursor tokens.

### 4. Rejected Alternatives
- Option B: Rejected because flat coordinate math produces inaccurate distance radii across mountainous East African topography.
- Option C: Rejected because managing a separate search cluster introduces synchronization lag on equipment availability.

### 5. Consequences & Invariants
- **Invariants Maintained:** Search queries return only equipment with status = ACTIVE and deleted_at IS NULL within the user-specified radius.
- **Operational Trade-offs:** GiST spatial indexes require tuning vacuum settings to maintain index efficiency.

### 6. Verification
- **Test File:** `backend/internal/modules/machinery/infrastructure/spatial_search_test.go`
- **Assertion:** TestSpatialSearch_RadiusAccuracy asserts that a machine 52km away is excluded from a 50km radius query.

### 7. Revisit Trigger
- Revisit if spatial search volume exceeds 5,000 queries/second across international borders.

---

## D-05-04-002: Dynamic Category-Specific Spec Faceting via Pre-Aggregated JSONB GIN Indexing

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [SRCH-05-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Faceted search requires dynamic filter aggregations based on category (e.g. excavator boom reach, crane tonnage, dump truck volume). Running dynamic aggregate queries across raw unindexed JSON documents causes high CPU spikes.

### 2. Considered Alternatives
1. **Option A:** Pre-indexed fleet_spec_indexed JSONB table with GIN jsonb_path_ops index and Valkey facet cache (15-min TTL)
   - *Pros:* Sub-millisecond facet retrieval; zero CPU overload on repeated filter searches.
   - *Cons:* Requires cache invalidation on new fleet listing approval.
2. **Option B:** Real-time JSON aggregation queries on every search request
   - *Pros:* Always shows real-time facet counts.
   - *Cons:* Database CPU saturation under concurrent user browsing spikes.
3. **Option C:** Static hardcoded filter facets per category in client code
   - *Pros:* Zero backend aggregation overhead.
   - *Cons:* Filters cannot adapt when owners list novel equipment capacities or brands.

### 3. Chosen Solution & Technical Rationale
Option A: Pre-indexed JSONB with GIN jsonb_path_ops index combined with a 15-minute Valkey cache for category facet counts. Cache invalidation is triggered on fleet approval or deactivation.

### 4. Rejected Alternatives
- Option B: Rejected due to heavy database CPU degradation during peak morning project requisition hours.
- Option C: Rejected because hardcoded filters fail to represent actual yard inventory.

### 5. Consequences & Invariants
- **Invariants Maintained:** Facet counts only aggregate active equipment available for booking.
- **Operational Trade-offs:** Facet counts may lag behind listing activations by up to 15 minutes during cache TTL.

### 6. Verification
- **Test File:** `backend/internal/modules/machinery/infrastructure/facet_search_test.go`
- **Assertion:** TestFacetSearch_CachedFacetsMatchDatabase asserts that facet counts match active database listing tallies within tolerance.

### 7. Revisit Trigger
- Revisit if listing volume reaches 100,000+ units requiring dedicated OLAP aggregations.

---

