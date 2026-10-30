# Partition Decisions: 10-01 — Relational Schema Architecture & PostGIS 3.4 Spatial Taxonomy Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `10-01`.

- **Governing System ADRs:** `ADR-003`, `ADR-005`
- **Subsystem Scope:** database
- **Requirement Targets:** SRS-SEC-10, ARCH-DB-DINK-CORE-V7.0

---

## D-10-01-001: Unified 111-Table Relational Schema with PostGIS 3.4 Spatial Extensions

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [DB-01-INF]
- **Subsystem:** [database]

### 1. Problem Statement & Context
The Dink Machinery domain spans 111 relational tables covering authentication, fleet taxonomies, attachments, bookings, inspections, disputes, dynamic configs, and double-entry accounting. Splitting into fragmented polyglot databases creates immense operational overhead and destroys ACID transactional guarantees.

### 2. Considered Alternatives
1. **Option A:** Single unified PostgreSQL 18.4 database with PostGIS 3.4, uuid-ossp, and pgcrypto extensions managing all 111 platform tables
   - *Pros:* Single database engine to backup, tune, and replicate; ACID transactions across all bounded contexts; spatial GiST indexing.
   - *Cons:* Requires disciplined index tuning and connection pool configuration.
2. **Option B:** Polyglot persistence: MongoDB for fleet specs, Redis for sessions, PostgreSQL for money, Neo4j for logistics routes
   - *Pros:* Specialized database per data structure.
   - *Cons:* Nightmarish operational complexity in local datacenter; partial failure modes; distributed transaction synchronization lag.
3. **Option C:** MySQL 8.0 without PostGIS
   - *Pros:* Common relational database.
   - *Cons:* Lacks spheroidal spatial indexing; inferior JSONB path indexing; poor window function optimization compared to PostgreSQL.

### 3. Chosen Solution & Technical Rationale
Option A: Unified PostgreSQL 18.4 database with PostGIS 3.4. All 111 tables are deployed via sequential golang-migrate migrations, leveraging native GiST spatial indexes, JSONB GIN indexing, and pgcrypto column encryption.

### 4. Rejected Alternatives
- Option B: Rejected due to catastrophic operational and maintenance overhead in emerging-market infrastructure.
- Option C: Rejected because PostgreSQL + PostGIS provides superior spatial capabilities and JSONB performance.

### 5. Consequences & Invariants
- **Invariants Maintained:** All geographic coordinates are persisted as PostGIS geography(POINT, 4326) in WGS 84 coordinate system.
- **Operational Trade-offs:** Requires dedicated DBA monitoring of vacuuming, WAL archiving, and shared_buffers memory tuning.

### 6. Verification
- **Test File:** `backend/internal/platform/database/schema_migration_test.go`
- **Assertion:** TestSchema_MigrationUpAndDownExecutesCleanly asserts migration 000001 applies and rolls back without errors.

### 7. Revisit Trigger
- Revisit if single database writes exceed 25,000 IOPS, requiring read-replica offloading.

---

## D-10-01-002: Strict Foreign Key Integrity and Immutable Double-Entry Ledger Partitioning

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [DB-05-INF]
- **Subsystem:** [database]

### 1. Problem Statement & Context
High-volume financial ledger tables (financial_ledger_entries) accumulate millions of immutable audit records, while referential integrity must be strictly guaranteed across parent booking and user tables.

### 2. Considered Alternatives
1. **Option A:** Strict foreign keys across all relational entities with declarative PostgreSQL range partitioning by fiscal year on financial_ledger_entries
   - *Pros:* Guaranteed referential integrity; fast partition-pruned ledger queries; easy archiving of historical fiscal years.
   - *Cons:* Requires partitioning configuration in initial migration DDL.
2. **Option B:** Unpartitioned monolithic ledger table without foreign keys
   - *Pros:* Simple table definition.
   - *Cons:* Orphaned financial records; sequential scans slow down dramatically as table grows past 10 million rows.
3. **Option C:** Application-layer foreign keys without database constraints
   - *Pros:* Fastest insertion speed.
   - *Cons:* Guaranteed data corruption and orphaned rows when application bugs or race conditions occur.

### 3. Chosen Solution & Technical Rationale
Option A: Strict foreign key constraints with declarative range partitioning by fiscal year (created_at) on financial_ledger_entries and audit_logs. Foreign key constraints maintain relational invariants across all transactional tables.

### 4. Rejected Alternatives
- Option B: Rejected because unpartitioned massive ledger tables cause severe performance degradation on annual financial audits.
- Option C: Rejected because database-enforced foreign keys are non-negotiable for financial ledger integrity.

### 5. Consequences & Invariants
- **Invariants Maintained:** All operational entity tables enforce relational foreign key constraints to prevent orphaned records.
- **Operational Trade-offs:** Queries on partitioned tables must include the partitioning column (created_at) to enable partition pruning.

### 6. Verification
- **Test File:** `backend/internal/platform/database/partitioning_test.go`
- **Assertion:** TestLedger_PartitionPruning asserts queries filtering by year scan only the corresponding fiscal partition table.

### 7. Revisit Trigger
- Revisit if audit log volume exceeds 100 million rows, requiring sub-partitioning by month.

---

