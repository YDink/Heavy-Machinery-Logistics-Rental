# ADR-010: Direct Bank Settlement and Pre-Funded Commission Wallet

- **Status:** Accepted
- **Date:** 2026-10-26
- **Decision Owner:** Dink Machinery Architecture Review Board
- **Affected Partitions:** 03-product-vision, 05-05, 05-07, 05-09, 05-24
- **Related Requirements:** BK-07-BE, BK-17-BE, CW-01-BE, NBE-ONPS-02-2020

---

## 1. Context & Problem Statement

Heavy machinery rentals in East Africa involve high monetary values ($50,000\text{ to }1,000,000+\text{ ETB}$). Commercial contractors settle predominantly via interbank wire transfers (Commercial Bank of Ethiopia, Awash, Dashen Bank). Holding customer rental funds in platform escrow requires specialized national payment system operator licenses and exposes Dink Machinery to severe balance-sheet regulatory scrutiny under National Bank of Ethiopia (NBE) directives. Conversely, invoicing fleet owners for commission after jobs complete leads to chronic commercial defaults.

---

## 2. The 5-Point Qualification Evaluation

- [x] **Alternative Test:** Evaluated direct bank transfer + pre-funded wallet vs. platform-held escrow custody vs. post-paid commission invoicing.
- [x] **Surprise Test:** Non-custodial fund flow: rental money bypasses the platform entirely; marketplace revenue is captured via automated wallet commission debits.
- [x] **Constraint Test:** Strict compliance with NBE Directive ONPS/02/2020 (prohibition of unlicensed customer deposit custody).
- [x] **Irreversibility Test:** Fundamental financial settlement flow governing the entire marketplace business model.
- [x] **Scope Test:** Directly affects booking confirmation, financial ledgers, owner onboarding, cancellation penalties, and dispute arbitration.

---

## 3. Considered Options

### Option 1: Direct Bank Transfer Rail (`MANUAL_TRANSFER`) + Pre-Funded `OwnerWallet`
- **Summary:** Renter wires rental funds directly to the owner's bank account and uploads payment evidence. Owner confirms receipt. Upon confirmation, the platform automatically debits marketplace commission basis points from the owner's pre-funded `OwnerWallet` within an atomic transaction (`BK-17-BE`).
- **Pros:**
  - Zero platform custody of customer deposits; 100% compliant with NBE banking directives.
  - Zero platform credit risk on rental payments; contractors utilize existing trusted banking rails.
  - Commission collection is guaranteed upfront; owners cannot bypass platform fees.
- **Cons:**
  - Requires manual bank slip verification and perceptual hash duplicate detection.
  - Owners must maintain pre-funded wallet balances or credit limits.

### Option 2: Platform-Held Escrow Custody Model
- **Summary:** Renters pay into a platform-controlled escrow bank account; platform releases funds to owners post-rental minus commission.
- **Pros:**
  - Centralized dispute hold and automated refund disbursement.
- **Cons:**
  - Violates NBE non-banking regulations without a Payment System Operator (PSO) license; contractor reluctance to route multi-million Birr wires through an aggregator.

### Option 3: Post-Paid Commission Invoicing
- **Summary:** Platform issues monthly invoices to fleet owners for accumulated rental commissions.
- **Pros:**
  - Low friction for owners during booking confirmation.
- **Cons:**
  - Catastrophic default rates; owners routinely ignore platform invoices once machines are mobilized on site.

---

## 4. Chosen Solution & Architectural Rationale

We adopt **Option 1: Direct Bank Transfer Rail (`MANUAL_TRANSFER`) + Pre-Funded `OwnerWallet`**.

1. **Non-Custodial Fund Flow:**
   ```
   [Renter] ──────── (Direct Bank Wire Transfer) ────────► [Owner Bank Account]
      │                                                              │
      │ Uploads Bank Receipt Slip                                   │ Confirms Receipt
      ▼                                                              ▼
   [Dink Machinery Backend] ── (Atomic pgx.Tx) ──► Debit Commission from [OwnerWallet]
   ```
2. **Atomic Commission Deduction (`BK-17-BE`):**
   When the owner confirms receipt of the bank wire:
   - Verification checks owner wallet balance: `balance_minor + credit_limit_minor >= commission_minor`.
   - If funds are insufficient, confirmation is blocked with an RFC 9457 error (`COMMISSION_WALLET_INSUFFICIENT_FUNDS`), instructing the owner to top up.
   - An immutable double-entry ledger record is written (`financial_ledger_entries`), linking debit and credit legs.
3. **Anti-Fraud Receipt Checks (`BK-07B`):**
   Uploaded bank slip images are processed with perceptual hashing (pHash) and SHA-256 digests to detect duplicate slips submitted across multiple bookings.

---

## 5. Rejected Alternatives & Technical Justification

- **Option 2 (Platform Escrow):** Rejected due to legal prohibition under NBE directives and commercial contractor refusal to deposit multi-million Birr sums with an unlicensed technology intermediary.
- **Option 3 (Post-Paid Invoicing):** Rejected because real-world empirical data in the East African machinery sector shows $> 60\%$ non-collection rates on post-paid invoices.

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - Full regulatory compliance from Day One without requiring a banking charter.
  - Zero platform liability for holding third-party financial balances.
  - Guaranteed commission monetization.
- **Negative / Technical Debt:**
  - Fleet owners must be educated on maintaining a pre-funded commission wallet balance.
- **Operational / Observability Impacts:**
  - Automated SMS/Push alerts notify owners when wallet balance drops below 20% of required thresholds.

---

## 7. Verification Strategy & Invariant Tests

- **Ledger Invariant Test:** `backend/internal/modules/payment/application/commission_deduction_test.go` verifies that every commission deduction creates balanced double-entry ledger legs summing to zero.
- **Credit Limit Boundary Test:** Verifies that an owner can utilize pre-authorized credit limits (`credit_limit_minor`) but is blocked when the overdraft ceiling is breached.

---

## 8. Revisit Conditions

- Acquisition of a formal Payment System Operator (PSO) license and central bank clearance for platform payment holding in Phase 3 (`ADR-010` evolution).
