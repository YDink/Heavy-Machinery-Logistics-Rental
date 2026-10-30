# Partition Decisions: 05-06 — Active Rental Management & Operational Telemetry Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-06`.

- **Governing System ADRs:** `ADR-006`, `ADR-012`
- **Subsystem Scope:** backend, mobile
- **Requirement Targets:** RENT-01-BE to RENT-12-BE

---

## D-05-06-001: Real-Time Incident Reporting with Offline Isar Outbox Queuing

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [RENT-04-BE]
- **Subsystem:** [mobile]

### 1. Problem Statement & Context
Mechanical breakdowns, hydraulic line ruptures, and job site work stoppages occur in remote construction zones without cellular reception. Operators must record incidents immediately with photographic evidence, timestamps, and GPS coordinates without app freezes.

### 2. Considered Alternatives
1. **Option A:** Local Isar NoSQL incident queue with immediate offline persistence and background retry worker
   - *Pros:* Zero data loss in dead zones; operator can record incident immediately and resume work; automatic retry on connectivity.
   - *Cons:* Requires local disk storage management on Android handsets.
2. **Option B:** Block incident submission until active network connection is detected
   - *Pros:* Simple synchronous client logic.
   - *Cons:* Operators cannot record incident at moment of failure; delayed reporting leads to disputed downtime hours.
3. **Option C:** Plain SQLite database with custom synchronization scripts
   - *Pros:* Standard SQL syntax.
   - *Cons:* High maintenance overhead and slower Dart serialization compared to native Isar.

### 3. Chosen Solution & Technical Rationale
Option A: Local Isar NoSQL incident queue with background outbox synchronization. Incident photos are compressed via Dart worker isolates, stored in the app sandbox, and queued in Isar with a client-generated UUIDv4 idempotency key.

### 4. Rejected Alternatives
- Option B: Rejected because remote job sites frequently have zero cellular coverage; blocking incident capture causes inaccurate downtime tracking.
- Option C: Rejected because Isar provides superior Dart object serialization and sub-millisecond query performance.

### 5. Consequences & Invariants
- **Invariants Maintained:** All incident reports capture device timestamp, GPS coordinates, operator user ID, and at least 1 photo.
- **Operational Trade-offs:** Incident notifications to owners and support agents are delayed until mobile connectivity is restored.

### 6. Verification
- **Test File:** `mobile/test/features/rental/data/incident_outbox_test.dart`
- **Assertion:** TestIncidentOutbox_SavesOfflineAndFlushesOnConnect verifies incident persists in Isar during airplane mode and uploads on network restore.

### 7. Revisit Trigger
- Revisit if direct satellite messaging (e.g. Android satellite SMS) becomes standard in East African devices.

---

## D-05-06-002: Daily Operating Log Tracking with Monotonic Engine Hour Validation

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [RENT-08-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Rentals billed on daily shift bases (8 operating hours/day) incur overtime charges for excess hours. Operators or owners may attempt to manipulate hour meter readings to falsify billable operational time.

### 2. Considered Alternatives
1. **Option A:** Strict monotonic hour meter validation: each logged reading must be >= previous verified reading, backed by camera-only photo proof
   - *Pros:* Mathematically eliminates negative hour rollbacks; tamper-evident photographic evidence.
   - *Cons:* Operators must take daily photos of the mechanical gauge.
2. **Option B:** Self-reported manual number input without photographic verification
   - *Pros:* Fast operator input.
   - *Cons:* Rampant falsification of operational hours leading to bitter disputes at rental close.
3. **Option C:** Automated CAN-bus IoT telematics streaming alone
   - *Pros:* Zero manual operator input.
   - *Cons:* Over 85% of machinery in Ethiopia and Kenya lack electronic telematics gateways in MVP phase.

### 3. Chosen Solution & Technical Rationale
Option A: Strict monotonic hour meter validation enforced at both mobile viewfinder (camera-only capture with bounding reticle) and backend validation layer (new_hours >= previous_hours). Readings are flagged for review if daily variance exceeds 24 hours.

### 4. Rejected Alternatives
- Option B: Rejected due to frequent fraudulent reporting of engine hours.
- Option C: Rejected because telematics penetration is too low in MVP emerging markets.

### 5. Consequences & Invariants
- **Invariants Maintained:** Hour meter reading must be monotonically non-decreasing: hours(T2) >= hours(T1).
- **Operational Trade-offs:** Requires operators to submit a daily gauge photo before ending the operational shift.

### 6. Verification
- **Test File:** `backend/internal/modules/booking/domain/hour_meter_validation_test.go`
- **Assertion:** TestHourMeter_RejectsRegression asserts that submitting 1,240 hours after a verified 1,250 reading returns HTTP 422 VALIDATION_ERROR.

### 7. Revisit Trigger
- Revisit when Phase 3 telematics IoT hardware (`ADR-007`, `07-06`) is installed across majority fleets.

---

