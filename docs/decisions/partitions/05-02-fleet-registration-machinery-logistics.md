# Partition Decisions: 05-02 — Fleet Registration, Machinery & Logistics Taxonomy Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-02`.

- **Governing System ADRs:** `ADR-003`, `ADR-005`, `ADR-012`
- **Subsystem Scope:** backend, database
- **Requirement Targets:** FL-01-BE to FL-14-BE

---

## D-05-02-001: Multi-Angle Photo Ingestion Pipeline with Synchronous ClamAV Stream Scanning

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [FL-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Fleet listings require 5 to 20 high-resolution photographs covering all equipment angles, serial number plates, and mechanical condition. Uploaded media must be protected against embedded malware, and stored efficiently in local sovereign storage.

### 2. Considered Alternatives
1. **Option A:** Synchronous ClamAV stream scanning over TCP socket (:3310) prior to SeaweedFS S3 persistence with isolate image compression
   - *Pros:* Zero infected files ever written to disk; immediate client feedback; full S3 API compatibility.
   - *Cons:* Adds ~150ms processing latency per uploaded photo during stream scanning.
2. **Option B:** Asynchronous post-upload quarantine bucket scanning
   - *Pros:* Fastest upload latency from client perspective.
   - *Cons:* Infected files sit on storage volumes before background worker scan; complex rollback choreography.
3. **Option C:** Public cloud S3 with AWS Rekognition and third-party SaaS antivirus
   - *Pros:* Managed AI tagging and scanning.
   - *Cons:* Violates Ethiopian Data Protection Proclamation No. 1321/2024 and incurs foreign exchange USD charges.

### 3. Chosen Solution & Technical Rationale
Option A: Synchronous ClamAV stream scanning (:3310) with local SeaweedFS S3 persistence. Photos are compressed on the mobile client via background Dart worker isolates (WebP, 1600px max, 82% quality) and streamed through ClamAV into SeaweedFS.

### 4. Rejected Alternatives
- Option B: Rejected because persisting unscanned binaries introduces an unacceptable attack window in multi-tenant environments.
- Option C: Rejected due to national data residency legal prohibitions and foreign currency payment restrictions.

### 5. Consequences & Invariants
- **Invariants Maintained:** Every persisted image is guaranteed clean of malware signatures. Listings enforce minimum 5 and maximum 20 verified photos.
- **Operational Trade-offs:** Requires running clamd daemon containers with persistent signature databases.

### 6. Verification
- **Test File:** `backend/internal/modules/machinery/application/upload_photo_test.go`
- **Assertion:** TestUploadPhoto_MalwareRejection asserts that an EICAR test string payload returns HTTP 422 INSPECT_MEDIA_SCAN_FAILED.

### 7. Revisit Trigger
- Revisit if hardware acceleration for antivirus streaming becomes necessary under massive daily upload spikes.

---

## D-05-02-002: Hybrid EAV Dynamic Spec Modeling with PostgreSQL JSONB GIN Indexing

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [FL-07-BE]
- **Subsystem:** [database]

### 1. Problem Statement & Context
Heavy machinery categories (excavators, bulldozers, tower cranes, lowboys) possess fundamentally different mechanical parameters (operating weight, boom length, lift capacity, engine power). Storing all dynamic fields in a monolithic table creates hundreds of sparse columns.

### 2. Considered Alternatives
1. **Option A:** Core relational table (fleets) paired with flattened fleet_spec_indexed JSONB column and GIN jsonb_path_ops index
   - *Pros:* ACID transactional integrity for status and pricing; flexible schemaless attributes; fast sub-15ms filtered search.
   - *Cons:* Requires application-layer validation against spec_field_definitions templates.
2. **Option B:** Pure Entity-Attribute-Value (EAV) relational table with one row per specification attribute
   - *Pros:* Pure relational model with foreign key validation per cell.
   - *Cons:* Horrific query performance requiring 20+ joins for single machine comparison queries.
3. **Option C:** Dedicated relational table per machinery category (excavators, cranes, etc.)
   - *Pros:* Strict typed columns per category.
   - *Cons:* Requires database schema migrations every time a new category or attribute is introduced.

### 3. Chosen Solution & Technical Rationale
Option A: Core relational table fleets combined with fleet_spec_indexed storing dynamic attributes in JSONB, indexed via GIN jsonb_path_ops. Category templates define required and comparable keys validated before insertion.

### 4. Rejected Alternatives
- Option B: Rejected because multi-join EAV tables degrade faceted search performance exponentially.
- Option C: Rejected because frequent category expansions would require disruptive database DDL migrations.

### 5. Consequences & Invariants
- **Invariants Maintained:** All comparable attributes conform strictly to unit definitions (kg, mm, hp, kPa) specified in category templates.
- **Operational Trade-offs:** GIN indexes increase database disk footprint and write amplification on fleet updates.

### 6. Verification
- **Test File:** `backend/internal/modules/machinery/infrastructure/spec_validation_test.go`
- **Assertion:** TestSpecValidation_RejectsMissingRequiredFields asserts that submitting an excavator without operating_weight_kg returns validation failure.

### 7. Revisit Trigger
- Revisit if custom machine attributes exceed 500 distinct keys per single category.

---

