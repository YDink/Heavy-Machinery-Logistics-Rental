# Partition Decisions: 05-22 — Universal Soft-Delete Policy & Tombstone Auditing Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-22`.

- **Governing System ADRs:** `ADR-003`, `ADR-008`
- **Subsystem Scope:** backend, database
- **Requirement Targets:** SD-01-BE to SD-06-BE

---

## D-05-22-001: Universal deleted_at Timestamp Invariant (SD-01-BE) Across All Mutable Relational Tables

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [SD-01-BE]
- **Subsystem:** [database, backend]

### 1. Problem Statement & Context
Accidental or malicious permanent deletion (DELETE FROM ...) of equipment listings, bookings, or user accounts destroys historical transactional records, breaks foreign key constraints, and prevents retrospective forensic audits during legal disputes.

### 2. Considered Alternatives
1. **Option A:** Universal soft-delete standard: all mutable operational entity tables declare deleted_at TIMESTAMPTZ NULL, and domain queries strictly enforce deleted_at IS NULL filters
   - *Pros:* Zero accidental data destruction; complete historical auditable trail; soft-deleted records can be restored by super admins.
   - *Cons:* Developers must remember to include deleted_at IS NULL on queries or use sqlc query templates.
2. **Option B:** Physical hard deletions (DELETE FROM ...) with ON DELETE CASCADE foreign keys
   - *Pros:* Automatic cleanup of child rows; smaller disk tables.
   - *Cons:* Catastrophic accidental data loss; permanently destroys legal evidence and historical financial ledgers.
3. **Option C:** Separate archive tables (e.g. fleets_archive) where deleted rows are copied before removal
   - *Pros:* Active tables remain small.
   - *Cons:* Requires maintaining dual schemas for every table; complex transactional triggers to move rows.

### 3. Chosen Solution & Technical Rationale
Option A: Universal deleted_at timestamp invariant (SD-01-BE). Physical DELETE statements are prohibited on mutable operational tables. Deletion executes via UPDATE ... SET deleted_at = NOW(), version = version + 1. All domain queries filter deleted_at IS NULL.

### 4. Rejected Alternatives
- Option B: Rejected because hard deletes permanently destroy audit trails and violate financial non-repudiation.
- Option C: Rejected because dual archive tables add immense schema migration maintenance overhead.

### 5. Consequences & Invariants
- **Invariants Maintained:** Physical SQL DELETE operations are strictly forbidden on production operational entity tables.
- **Operational Trade-offs:** Partial unique indexes must account for soft deletes: CREATE UNIQUE INDEX ... WHERE deleted_at IS NULL.

### 6. Verification
- **Test File:** `backend/internal/platform/database/soft_delete_audit_test.go`
- **Assertion:** TestSoftDelete_SchemaInvariant asserts that all 111 operational tables declare deleted_at column.

### 7. Revisit Trigger
- Revisit if historical soft-deleted records older than 7 years require automated cold-storage partition purging.

---

## D-05-22-002: Cascading Soft-Delete Policy with Tombstone Auditing in Historical Scopes

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [SD-04-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
When a fleet listing is soft-deleted by its owner, related entities (unbooked calendar availability blocks, draft maintenance logs) must be soft-deleted, but completed historical bookings and ledger records must remain completely intact for accounting auditability.

### 2. Considered Alternatives
1. **Option A:** Selective cascading soft-delete managed via application service layer: soft-deletes dependent future availability, leaves historical bookings untouched, and records an audit log tombstone
   - *Pros:* Preserves complete accounting history; clears future calendar availability cleanly; auditable reason captured.
   - *Cons:* Requires explicit application service logic rather than automatic database triggers.
2. **Option B:** Database triggers automatically cascading deleted_at to all child tables including past bookings
   - *Pros:* Automated in database.
   - *Cons:* Corrupts past rental records, making it appear that completed rentals were deleted.
3. **Option C:** Block fleet soft-deletion if it has ever been booked in the past
   - *Pros:* Preserves everything.
   - *Cons:* Owners can never remove retired or sold equipment from their yard inventory.

### 3. Chosen Solution & Technical Rationale
Option A: Selective application-layer cascading soft-delete. When a fleet is soft-deleted: future unbooked availability slots are soft-deleted; active or pending bookings block the deletion with HTTP 409; historical completed bookings remain untouched; an audit log entry records the deletion event.

### 4. Rejected Alternatives
- Option B: Rejected because cascading deletion to historical completed bookings destroys financial ledger integrity.
- Option C: Rejected because owners must be able to retire sold or decommissioned machinery from active view.

### 5. Consequences & Invariants
- **Invariants Maintained:** Historical COMPLETED, RATED, or CLOSED booking records are never modified when parent fleet assets are soft-deleted.
- **Operational Trade-offs:** Application services must coordinate soft-deletion across multiple repositories within a transaction.

### 6. Verification
- **Test File:** `backend/internal/modules/machinery/application/delete_fleet_test.go`
- **Assertion:** TestDeleteFleet_PreservesHistoricalBookings asserts completed past bookings remain queryable after fleet soft-deletion.

### 7. Revisit Trigger
- Revisit if automated asset lifecycle de-registration workflows are standardized with ministry equipment registries.

---

