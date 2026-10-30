# Partition Decisions: 05-12 — Certified Operator Telematics & GPS Ping Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-12`.

- **Governing System ADRs:** `ADR-005`, `ADR-006`
- **Subsystem Scope:** backend, mobile
- **Requirement Targets:** OPS-01-BE to OPS-12-BE

---

## D-05-12-001: Periodic GPS Ping Aggregation with PostGIS Temporal Breadcrumbs

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [OPS-01-BE]
- **Subsystem:** [backend, database]

### 1. Problem Statement & Context
Equipment owners and contractors require visibility into machine location on large construction sites. Constant real-time WebSocket streaming drains operator mobile phone batteries and overwhelms database write bandwidth in remote areas.

### 2. Considered Alternatives
1. **Option A:** Batched periodic GPS pings (every 5-15 minutes during active shift) aggregated on mobile and inserted into fleet_location_histories with PostGIS geography points
   - *Pros:* Conserves mobile battery; minimizes cellular bandwidth; maintains historical breadcrumb trails with spatial indexing.
   - *Cons:* Location is not sub-second real-time.
2. **Option B:** High-frequency real-time WebSocket GPS streaming (every 1-5 seconds)
   - *Pros:* Live moving vehicle marker on map.
   - *Cons:* Drains operator phone battery in 2 hours; cellular connection drops cause constant reconnect storms.
3. **Option C:** Manual check-in by operator twice per day
   - *Pros:* Zero automated background tracking.
   - *Cons:* Provides zero security protection against equipment theft or unauthorized off-site mobilization.

### 3. Chosen Solution & Technical Rationale
Option A: Batched periodic GPS pings (5-minute interval during active rentals). Mobile app batches location updates and flushes them to /api/v1/operator/gps-ping. Backend inserts into fleet_location_histories with PostGIS geometry and timestamp.

### 4. Rejected Alternatives
- Option B: Rejected because aggressive continuous streaming is unfeasible on low-end operator phones on dusty construction sites.
- Option C: Rejected because manual check-ins fail to prevent equipment theft or perimeter violations.

### 5. Consequences & Invariants
- **Invariants Maintained:** GPS breadcrumb points must be associated with an active booking lease or authorized yard storage.
- **Operational Trade-offs:** Fleet map marker reflects position within a 5-minute freshness window.

### 6. Verification
- **Test File:** `backend/internal/modules/machinery/application/gps_ping_test.go`
- **Assertion:** TestGPSPing_InsertsTemporalBreadcrumb asserts valid coordinates create a new point in fleet_location_histories.

### 7. Revisit Trigger
- Revisit when Phase 3 dedicated hardwired IoT CAN-bus telematics (`07-06`) provide continuous vehicle-powered telemetry.

---

## D-05-12-002: Circular Geofence Violation Detection with Edge-Triggered Alarms

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [OPS-06-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Heavy machinery rented for a designated project site (e.g. Bole Lemi Industrial Park) is strictly prohibited from leaving the site perimeter without owner authorization. Unauthorized mobilization indicates machine poaching or equipment theft.

### 2. Considered Alternatives
1. **Option A:** PostGIS spatial containment assertion (ST_DWithin on site centroid + buffer radius) evaluated on each GPS ping, triggering immediate push/SMS alarms on exit
   - *Pros:* Automated perimeter breach detection; instant alerts to owner and platform ops; spatial database optimization.
   - *Cons:* Requires defining site centroid and radius at booking confirmation.
2. **Option B:** Manual inspection of daily breadcrumb tracks by fleet owners
   - *Pros:* Zero automated geofence computation.
   - *Cons:* Breaches are discovered hours or days after the machine has already left the area.
3. **Option C:** Complex multi-polygon geofences drawn by hand on mobile map
   - *Pros:* Custom boundary shapes.
   - *Cons:* Complex mobile UI drawing; high spatial join computational overhead.

### 3. Chosen Solution & Technical Rationale
Option A: Circular geofence boundary defined by project site location and radius_km. On each GPS ping ingestion, PostGIS evaluates ST_DWithin(ping.location, booking.site_location, booking.geofence_radius_meters). If false, an immediate SECURITY_GEOFENCE_BREACH event is published.

### 4. Rejected Alternatives
- Option B: Rejected because delayed breach awareness permits machine theft and unauthorized third-party subcontracting.
- Option C: Rejected because circular boundaries cover 95% of job site use cases with a fraction of spatial calculation overhead.

### 5. Consequences & Invariants
- **Invariants Maintained:** A geofence breach event automatically records an immutable security incident log and alerts the owner.
- **Operational Trade-offs:** GPS drift on low-end handsets can trigger occasional false perimeter boundary warnings.

### 6. Verification
- **Test File:** `backend/internal/modules/machinery/domain/geofence_test.go`
- **Assertion:** TestGeofence_BreachDetected asserts a ping 2.5km away from a 2.0km geofence generates a GEOFENCE_BREACH event.

### 7. Revisit Trigger
- Revisit if complex polygon cadastral parcel boundaries are integrated from municipal GIS departments.

---

