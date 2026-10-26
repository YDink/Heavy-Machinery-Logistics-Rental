# Open Questions & Architectural Blockers Ledger — Dink Machinery

This document tracks unresolved architectural questions, third-party integration ambiguities, and regulatory constraints across Dink Machinery requirement partitions. Per `AGENTS.md § 8`, active questions must be audited before initiating implementation on affected partitions.

---

## 1. Registry & Status Overview

| ID | Title | Impacted Partitions | Subsystems | Status | Blocking Phase | Target Resolution |
|---|---|---|---|---|---|---|
| `OQ-001` | Telebirr & CBEBirr Webhook Reliability vs. Polling Fallback | `05-06` (Payments) | Backend | **Resolved** | Phase 1 | `ADR-008` (Signed Webhook + Dual Reconcile) |
| `OQ-002` | Offline Site Inspection Handover under Zero-Connectivity | `05-07` (Inspections) | Mobile, Backend | **Open** | Phase 2 (P05-07) | 2026-11-15 |
| `OQ-003` | Operator Ministry Certification Verification Automation | `05-03` (KYC & Auth) | Admin, Backend | **Deferred** | Phase 2 (P05-03) | 2026-12-01 |
| `OQ-004` | Telematics/CAN-Bus IoT Streaming vs. Manual Meter Entry | `05-04` (Machinery) | Backend, Infra | **Resolved** | Phase 1 | `ADR-011` (Dual-Mode EAV Meter Model) |
| `OQ-005` | Cross-Border Equipment Transit & Customs Escrow Retention | `05-02`, `05-06` | Backend, Ledger | **Open** | Phase 2 (P05-02) | 2026-11-30 |
| `OQ-006` | Fuel Tank Discrepancy Margin & Automatic Deduction Caps | `05-07`, `05-08` | Backend, Mobile | **Open** | Phase 2 (P05-08) | 2026-11-20 |
| `OQ-007` | VAT vs. Withholding Tax Application Order on Escrow Payouts | `05-06` (Ledger) | Backend, Ledger | **Resolved** | Phase 1 | `ADR-015` (Multi-Country Config Tax Basis) |

---

## 2. Detailed Problem Statements & Evaluations

### OQ-001: Telebirr & CBEBirr Webhook Reliability vs. Polling Fallback
- **Impacted Requirements:** `[REQ-05-06-01]`, `[BK-17-BE]`
- **Status:** **Resolved** (Integrated into `ADR-008`)
- **Context:** Local Ethiopian payment rails (Telebirr B2C/C2B and CBEBirr) experience periodic network dropouts and delayed HTTP callback delivery, leading to state inconsistencies where payment has settled externally but the platform hold remains in `PENDING_PAYMENT`.
- **Proposed Paths:**
  1. Rely exclusively on external webhooks with exponential backoff on client retry.
  2. Implement an out-of-band asynchronous reconciliation worker that polls provider transaction query endpoints when webhooks are not delivered within 90 seconds.
- **Resolution:** Path 2 adopted. Payment state machine transitions from `PENDING_PAYMENT` to `PAYMENT_PROCESSING`, deploying a scheduled Valkey-backed task that queries gateway transaction APIs if no signed webhook arrives within 90 seconds. Enforced in single-transaction ledger entry.

---

### OQ-002: Offline Site Inspection Handover under Zero-Connectivity
- **Impacted Requirements:** `[REQ-05-07-04]`, `[REQ-05-07-08]`
- **Status:** **Open**
- **Blocking Target:** Feature Partition `feature/05-07-machinery-inspections`
- **Context:** Heavy equipment handovers frequently occur at remote construction corridors (e.g., dams, highways, agricultural expansions) with zero GSM/LTE network coverage. Both Renter and Owner/Operator need to mutually attest condition, engine hour meter reading, and fuel levels without live server connectivity.
- **Proposed Alternatives:**
  - **Alternative A (Cryptographic Peer-to-Peer Handshake):** Use local Bluetooth Low Energy (BLE) or encrypted QR code handshakes where both devices sign an inspection snapshot locally using their private keys stored in `flutter_secure_storage`. The signed bundle is queued locally in Isar DB and synchronized once network connectivity resumes.
  - **Alternative B (Unilateral Draft Sync):** The party with the device captures inspection media and countersigns offline. The inspection remains in `PENDING_PEER_ACK` until synced to backend and verified by the counterparty within a 24-hour server-side grace window.
- **Pending Decision:** Legal enforceability of offline asymmetric cryptographic signatures under Ethiopian Electronic Signature Proclamation No. 1072/2018.

---

### OQ-003: Operator Ministry Certification Verification Automation
- **Impacted Requirements:** `[REQ-05-03-05]`
- **Status:** **Deferred**
- **Blocking Target:** Partition `feature/05-03-operator-directory`
- **Context:** Specialized heavy equipment (excavators, tower cranes, bulldozers) requires certified operator credentials issued by the Ministry of Transport and Logistics. There is currently no publicly accessible government REST API to synchronously validate operator badge numbers.
- **Proposed Alternatives:**
  - **Alternative A (Manual Back-Office OCR & Audit):** Upload scanned certificate via mobile app, run Tesseract/AWS Textract OCR to extract badge ID and expiry, and place profile in `PENDING_ADMIN_REVIEW` until an admin verifies physical ledger copies.
  - **Alternative B (Third-Party Aggregator Integration):** Partner with licensed local verification aggregators once private B2B APIs go live.
- **Current Mitigation:** Manual verification queue supported by Next.js Admin Dynamic Table (`AD-06-BE`). Automated verification deferred to Phase 3.

---

### OQ-004: Telematics/CAN-Bus IoT Streaming vs. Manual Meter Entry
- **Impacted Requirements:** `[REQ-05-04-03]`, `[REQ-05-07-02]`
- **Status:** **Resolved** (Integrated into `ADR-011`)
- **Context:** Modern fleet units (CAT, Komatsu) feature OEM telematics transmitting continuous GPS and engine hour readings over CAN-Bus, whereas older local machinery relies purely on mechanical dashboard meters.
- **Resolution:** Adopted a dual-mode schema in `machinery_specs`: Assets with active telematics hardware bind to an IoT gateway stream via MQTT/HTTP; assets without telematics require photographically attested manual inspection meter entries (`SubmitInspectionRequest.operatingHoursMeter`).

---

### OQ-005: Cross-Border Equipment Transit & Customs Escrow Retention
- **Impacted Requirements:** `[REQ-05-02-04]`, `[REQ-05-06-09]`, `[LOC-01-BE]`
- **Status:** **Open**
- **Blocking Target:** Multi-Country Launch (Kenya/Ethiopia Border Operations)
- **Context:** When machinery cross-leases between jurisdictions (e.g., Moyale border zone between Ethiopia and Kenya), equipment is subject to customs bonded security deposits and differing cross-border withholding tax rates.
- **Proposed Alternatives:**
  - **Alternative A (Multi-Currency Partitioned Escrow):** Escrow engine segregates local rental earnings from customs guarantee deposits, holding customs guarantees in escrow until border clearance documentation is validated.
  - **Alternative B (Disallow Cross-Border Leases):** Confine fleet availability strictly within country borders by enforcing a spatial geofence filter inside PostGIS queries (`country_code` boundary bounding polygon).
- **Owner:** Legal Counsel & Lead Platform Architect.

---

### OQ-006: Fuel Tank Discrepancy Margin & Automatic Deduction Caps
- **Impacted Requirements:** `[REQ-05-07-06]`, `[REQ-05-08-03]`
- **Status:** **Open**
- **Blocking Target:** Partition `feature/05-08-dispute-resolution`
- **Context:** Fuel levels reported between check-in and check-out often differ slightly due to sloping terrain or imprecise analog fuel gauges. We must establish a platform tolerance threshold before triggering an automated deposit deduction.
- **Proposed Alternatives:**
  - **Alternative A:** Strict 5% tolerance threshold. Discrepancies >5% calculate automatic refueling charge based on prevailing fuel price dynamically retrieved from `country_configs`.
  - **Alternative B:** Automatic charge capped at 10% of total deposit hold; larger discrepancies automatically transition the booking into state `DISPUTED` for manual mediation.
- **Target Resolution Date:** Prior to finalizing inspection acceptance tests.

---

### OQ-007: VAT vs. Withholding Tax Application Order on Escrow Payouts
- **Impacted Requirements:** `[REQ-05-06-04]`, `[LOC-01-BE]`
- **Status:** **Resolved** (Integrated into `ADR-015`)
- **Context:** Different jurisdictions require differing computation sequences for standard Value Added Tax (VAT) and Withholding Tax (WHT) on marketplace facilitation fees versus gross owner rental payouts.
- **Resolution:** Parameterized inside `country_configs` using integer basis points (`bps`). The computation pipeline executes:
  1. `gross_amount` computed from base daily rate.
  2. Platform commission deducted from gross.
  3. `vat_amount = (commission * vat_bps) / 10000`.
  4. Withholding tax applied to owner net balance in accordance with jurisdictional tax registration certificate (`withholding_tax_bps`).
  All calculations stored as immutable integer minor units in `financial_ledger_entries`.

---

## 3. Governance for Resolving Blockers

1. **Escalation Path:** Open questions are raised via PR or issue referencing the target partition.
2. **Resolution Mechanism:**
   - Platform-wide architectural impacts (>1 partition) must resolve into an **L1 ADR** (`docs/decisions/ADR-NNN-<slug>.md`).
   - Partition-local technical trade-offs resolve into an **L2 Partition Decision** (`docs/decisions/partitions/<partition-id>.md`).
3. **Closing Protocol:** When an open question resolves, its status is updated in this ledger, linking directly to the corresponding ADR or partition decision record.