# Partition Decisions: 05-09 — Commission Revenue, Owner Wallets & Settlement Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-09`.

- **Governing System ADRs:** `ADR-002`, `ADR-010`, `ADR-015`
- **Subsystem Scope:** backend, database
- **Requirement Targets:** COM-01-BE to COM-14-BE

---

## D-05-09-001: Automated Commission Deduction from Pre-Funded OwnerWallet with Overdraft Controls

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [COM-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Marketplace platform revenue depends on capturing pre-agreed commission basis points (default 1000 bps / 10%) on confirmed bookings. Because customer funds flow directly to owner bank accounts via interbank wire transfer, commission cannot be held at checkout. The system must collect fees without post-paid invoicing defaults.

### 2. Considered Alternatives
1. **Option A:** Mandatory pre-funded OwnerWallet with atomic debit during booking confirmation and configurable credit overdraft buffer
   - *Pros:* Guaranteed commission capture; zero default risk; trusted owners receive temporary credit buffers.
   - *Cons:* Owners must top up wallets before accepting high-volume reservations.
2. **Option B:** Post-paid monthly commission invoicing with 14-day payment terms
   - *Pros:* Zero upfront financial friction for new fleet owners.
   - *Cons:* Severe default rates (>60% uncollected in regional markets) once equipment is already mobilized.
3. **Option C:** Holding customer payment in platform escrow and disbursing net payout post-rental
   - *Pros:* Direct deduction at payout source.
   - *Cons:* Requires banking payment system operator charter under NBE regulations; high regulatory and compliance overhead.

### 3. Chosen Solution & Technical Rationale
Option A: Mandatory pre-funded OwnerWallet with atomic debit at booking confirmation. When owner verifies bank slip, commission is debited immediately. Trusted owners with verified track records can be granted credit_limit_minor by admin finance controllers.

### 4. Rejected Alternatives
- Option B: Rejected due to catastrophic commercial invoice default rates in the East African construction industry.
- Option C: Rejected because holding customer deposits violates National Bank of Ethiopia non-banking directives.

### 5. Consequences & Invariants
- **Invariants Maintained:** Owner wallet balance + credit_limit_minor must be >= commission_minor before booking confirmation is permitted.
- **Operational Trade-offs:** Fleet owners must be actively notified via SMS/Push when wallet balance is low.

### 6. Verification
- **Test File:** `backend/internal/modules/payment/application/wallet_debit_test.go`
- **Assertion:** TestWalletDebit_BlocksConfirmationOnInsufficientFunds asserts confirmation fails if wallet funds + credit limit are below commission.

### 7. Revisit Trigger
- Revisit if central bank licenses allow direct split settlement in Phase 3 (`07-02`).

---

## D-05-09-002: Immutable Double-Entry Ledger Posting with Reconciled Batch Payout Sequences

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [LED-01-BE]
- **Subsystem:** [database]

### 1. Problem Statement & Context
Financial transactions (wallet top-ups, commission debits, dispute refunds, compensation credits) require absolute auditability. Mutable balance columns without double-entry tracking make external auditing impossible and conceal arithmetic errors.

### 2. Considered Alternatives
1. **Option A:** Immutable append-only double-entry ledger (financial_ledger_entries) where every transaction balances debits and credits to zero sum
   - *Pros:* Provable mathematical correctness; immutable audit trail; conforms to GAAP and statutory financial auditing standards.
   - *Cons:* Requires disciplined bookkeeping accounting models for all platform mutations.
2. **Option B:** Single balance column on owner_wallets table updated via direct increment/decrement
   - *Pros:* Simple CRUD operations.
   - *Cons:* Zero audit trail; balance corruption or arithmetic errors cannot be reconstructed or audited.
3. **Option C:** Periodic reconciliation snapshot table without double-entry legs
   - *Pros:* Low write overhead.
   - *Cons:* Fails financial auditing; cannot prove that money was not created or destroyed.

### 3. Chosen Solution & Technical Rationale
Option A: Immutable double-entry ledger. Every balance mutation inserts paired debit and credit legs into financial_ledger_entries with transaction_type, currency, and reference entity ID within the same transaction.

### 4. Rejected Alternatives
- Option B: Rejected because mutable balance columns without ledger journals violate financial non-repudiation invariants.
- Option C: Rejected due to failure to meet National Bank of Ethiopia audit standards.

### 5. Consequences & Invariants
- **Invariants Maintained:** Sum of debits must equal sum of credits for every transaction ID: SUM(amount_minor) == 0 across debit and credit legs.
- **Operational Trade-offs:** High row count in ledger tables requires partition maintenance by fiscal year.

### 6. Verification
- **Test File:** `backend/internal/modules/payment/infrastructure/ledger_integrity_test.go`
- **Assertion:** TestLedgerIntegrity_ZeroSumRule asserts that attempting to commit an unbalanced ledger transaction triggers database check error.

### 7. Revisit Trigger
- Revisit if financial transaction volume exceeds 50 million rows/year, requiring PostgreSQL table partitioning.

---

