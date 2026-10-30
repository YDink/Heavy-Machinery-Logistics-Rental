# Partition Decisions: 05-19 — Account Deletion, Right-to-Erasure & Statutory Retention Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-19`.

- **Governing System ADRs:** `ADR-005`, `ADR-010`
- **Subsystem Scope:** backend, database
- **Requirement Targets:** DEL-01-BE to DEL-10-BE

---

## D-05-19-001: Two-Stage Soft Deletion Lifecycle with 30-Day Recovery Grace Period

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [DEL-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Users requesting account deletion may change their mind or may have initiated deletion by mistake. Immediate hard deletion is irreversible and can inadvertently delete in-flight rental leases, active damage claims, or uncollected commissions.

### 2. Considered Alternatives
1. **Option A:** Two-stage deletion lifecycle: Stage 1 marks account PENDING_DELETION with a 30-day grace period; Stage 2 permanently anonymizes PII after grace expires
   - *Pros:* Prevents accidental irreversible account loss; allows users to reactivate during 30 days; blocks deletion if active rentals exist.
   - *Cons:* Requires automated cron to process Stage 2 permanently.
2. **Option B:** Immediate hard deletion (DELETE FROM users WHERE id = $1)
   - *Pros:* Instant data removal.
   - *Cons:* Irreversible data loss; breaks foreign key constraints on past bookings, inspections, and financial ledgers.
3. **Option C:** Permanent soft-deletion flag without ever purging or anonymizing PII
   - *Pros:* Preserves all records indefinitely.
   - *Cons:* Violates statutory Right to Erasure under Ethiopian Data Protection Proclamation No. 1321/2024 and Kenya Data Protection Act.

### 3. Chosen Solution & Technical Rationale
Option A: Two-stage deletion lifecycle. Stage 1 checks that user has zero ACTIVE or CONFIRMED rentals and zero negative wallet balance, then sets status = PENDING_DELETION and revocation of all JWT session families. If not cancelled within 30 days, Stage 2 anonymization cron executes.

### 4. Rejected Alternatives
- Option B: Rejected because hard deletes destroy historical relational integrity and prevent dispute arbitration.
- Option C: Rejected because retaining plaintext personal data indefinitely violates sovereign privacy legislation.

### 5. Consequences & Invariants
- **Invariants Maintained:** An account cannot enter PENDING_DELETION if it is associated with any active rental, open dispute, or negative wallet balance.
- **Operational Trade-offs:** Users in PENDING_DELETION cannot book or list equipment unless they explicitly cancel the deletion request.

### 6. Verification
- **Test File:** `backend/internal/modules/auth/application/account_deletion_test.go`
- **Assertion:** TestAccountDeletion_BlockedWithActiveBooking asserts deletion request returns HTTP 409 if user has an ongoing rental.

### 7. Revisit Trigger
- Revisit if national data protection authorities mandate shorter (e.g. 14-day) statutory deletion grace periods.

---

## D-05-19-002: 10-Year Statutory Financial Ledger Retention and Selective PII Cryptographic Shredding

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [DEL-06-BE]
- **Subsystem:** [database]

### 1. Problem Statement & Context
National tax and commercial laws (Commercial Code of Ethiopia, Kenya Revenue Authority) mandate retaining financial transaction records for 10 years. Conversely, data protection laws grant citizens the right to erasure of personal data. The system must reconcile these opposing legal requirements.

### 2. Considered Alternatives
1. **Option A:** Selective cryptographic shredding: personal attributes (full name, phone, national ID, avatar) are overwritten with salted hashes/tombstones, while financial_ledger_entries retain immutable numeric amounts and anonymized user IDs
   - *Pros:* 100% compliant with both 10-year financial tax retention and data privacy erasure laws; preserves audit balance integrity.
   - *Cons:* Requires irreversible column anonymization routines.
2. **Option B:** Hard-deleting financial ledgers along with personal records
   - *Pros:* Complete data removal.
   - *Cons:* Direct violation of national commercial tax codes; exposes Dink Machinery to criminal tax fraud penalties.
3. **Option C:** Exempting all users who ever transacted from account deletion entirely
   - *Pros:* Keeps all data intact.
   - *Cons:* Direct violation of citizen data privacy rights under Proclamation No. 1321/2024.

### 3. Chosen Solution & Technical Rationale
Option A: Selective cryptographic shredding and financial ledger preservation. At Stage 2 execution, user PII columns in users, owner_profiles, and user_documents are overwritten with anonymous placeholders (e.g. "DELETED_USER_1234"), while financial_ledger_entries and booking_status_histories retain immutable currency amounts and audit timestamps.

### 4. Rejected Alternatives
- Option B: Rejected because deleting tax and audit ledger records violates Ethiopian and Kenyan commercial law.
- Option C: Rejected because denying privacy rights exposes the platform to regulatory telecommunications sanctions.

### 5. Consequences & Invariants
- **Invariants Maintained:** Financial ledger records and double-entry balancing rows are never hard-deleted during account erasure.
- **Operational Trade-offs:** Historical booking reporting reflects anonymized user placeholders for deleted accounts.

### 6. Verification
- **Test File:** `backend/internal/modules/auth/infrastructure/anonymization_cron_test.go`
- **Assertion:** TestAnonymization_ShredsPIIWhilePreservingLedger asserts personal name is erased while ledger debit/credit rows remain intact.

### 7. Revisit Trigger
- Revisit if commercial tax retention statutes are amended to allow shorter financial document storage windows.

---

