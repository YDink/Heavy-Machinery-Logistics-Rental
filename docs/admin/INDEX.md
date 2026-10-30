# Heavy Machinery & Logistics Rental Aggregator — Dink Machinery
## Admin Web Console Subsystem Specifications — Master Index & Traceability Matrix (`docs/admin/INDEX.md`)

**Document Identifier:** `INDEX-ADM-DINK-CORE-V7.0`  
**Classification:** Canonical Administrative Web Console UI/UX, Server-Driven Metadata & Technical Behavior Index  
**Governing Documents:** [`Dink_Machinery_UI_UX_Design_Brief.md`](../source/Dink_Machinery_UI_UX_Design_Brief.md), [`00-design-system.md`](../shared/design/00-design-system.md), [`design-tokens.json`](../shared/design/design-tokens.json), [`admin/DESIGN.md`](../components/admin/DESIGN.md), [`contracts/openapi.yaml`](../../contracts/openapi.yaml)  
**Target Runtime:** Next.js 16.x (React 19, App Router), TypeScript 5.7+, Tailwind CSS v4, `@tanstack/react-table`, `@tanstack/react-virtual`, `react-hook-form` + `zod`  
**Supported Administrative Personas:** Platform Super Admin (`SUPER_ADMIN`), Compliance & KYC Officer (`COMPLIANCE_OFFICER`), Financial Controller (`FINANCE_CONTROLLER`), Dispute Arbitrator (`DISPUTE_ARBITRATOR`), Operations & Dispatch Specialist (`OPS_SUPPORT`)  
**Core Architecture:** Server-Driven Dynamic Metadata Engine (`/api/v1/admin/ui-config`, `AD-06-BE`) + React Server Components (RSC) + Mutating Server Actions (`'use server'`) + ASVS Level 3 `HttpOnly` Cookie Session Enclave  
**Total Operational Console Specifications:** 38 Specifications across 9 Bounded Modules

---

## 1. Architectural Invariants Enforced Across All Admin Specifications

1. **Server-Driven Dynamic UI Synthesis (`AD-06-BE`, `SRS 05-11`, `ADR-008`):**
   - Navigation sidebars (`admin_sidebar_items`), data grid column schemas (`admin_table_configs`), faceted filter bars, and RBAC-gated row actions (`admin_actions`) are synthesized dynamically from `/api/v1/admin/ui-config`. Hardcoded client-side permission matrices or static column definitions are prohibited.
2. **Non-Custodial Financial Ledger & Pre-Funded `OwnerWallet` Accounting (`ADR-010`, `BK-17-BE`, `SRS 05-09`):**
   - All monetary values are represented strictly as 64-bit signed integer minor units (`Santims` for `ETB`, `Cents` for `KES`/`USD`) paired with an ISO 4217 currency code (`MoneyAmount`).
   - Every financial state mutation writes balanced `DEBIT` and `CREDIT` legs to the immutable `financial_ledger_entries` table ($\sum \text{amount\_minor} = 0$) inside a single `pgx.Tx` transaction alongside `owner_wallets` balance updates.
3. **OWASP ASVS Level 3 Security & Zero-PII Masking (`SRS 05-01`, `05-24`):**
   - Admin authentication enforces WebAuthn FIDO2 Passkeys and hardware TOTP with `HttpOnly; Secure; SameSite=Strict` ES256 session cookies (zero tokens in `localStorage`).
   - Sensitive PII (National IDs, TINs, bank account numbers, phone numbers) is masked by default (`****5678`) and requires an audited "Click to Reveal" justification modal logged to `admin_audit_logs`.

---

## 2. Master Specification Catalog & Requirement Traceability Matrix

### 2.1 Application Shell, Authentication & RBAC (`docs/admin/shell/`)
| Spec ID | Specification Document | Covered SRS Requirements | Covered Flows & ADRs | Primary Focus |
|---|---|---|---|---|
| `SPEC-ADM-SHL-01` | [`01-admin-layout-navigation.md`](./shell/01-admin-layout-navigation.md) | `REQ-05-11` (`AD-01`, `AD-06`) | `FLOW-13`, `ADR-008` | Next.js 16 App Router shell, dynamic sidebar from `/api/v1/admin/ui-config`, `ET`/`KE` jurisdiction switcher, `Cmd+K` palette. |
| `SPEC-ADM-SHL-02` | [`02-admin-auth-mfa-passkey.md`](./shell/02-admin-auth-mfa-passkey.md) | `REQ-05-01` (`REG-10`–`REG-15`) | `FLOW-APPENDIX-B`, `ADR-011` | ASVS L3 WebAuthn FIDO2 Passkey login, RFC 6238 TOTP fallback, `HttpOnly` cookie session, 15-min idle lock. |
| `SPEC-ADM-SHL-03` | [`03-role-based-access-control.md`](./shell/03-role-based-access-control.md) | `REQ-04`, `REQ-05-11` (`AD-02`) | `FLOW-13` | Granular RBAC matrix, Maker-Checker Separation of Duties (SoD), audited PII "Click to Reveal" masking gate. |

### 2.2 Operations, Compliance & Integrity Queues (`docs/admin/operations-queues/`)
| Spec ID | Specification Document | Covered SRS Requirements | Covered Flows & ADRs | Primary Focus |
|---|---|---|---|---|
| `SPEC-ADM-OPS-01` | [`01-executive-kpi-dashboard.md`](./operations-queues/01-executive-kpi-dashboard.md) | `REQ-05-11` (`AD-01`), `REQ-07-09` | `FLOW-13` | Live operational cockpit: bookings by state, minor-unit GMV & commission volume, SLA breach alerts, PostGIS heatmap. |
| `SPEC-ADM-OPS-02` | [`02-kyc-user-verification-queue.md`](./operations-queues/02-kyc-user-verification-queue.md) | `REQ-05-01` (`REG-05`–`REG-08`) | `FLOW-03`, `FLOW-15` | Fayda Digital ID / Kebele ID inspection workbench, phone migration liveness audit (`FLOW-15`), structured rejection codes. |
| `SPEC-ADM-OPS-03` | [`03-business-license-queue.md`](./operations-queues/03-business-license-queue.md) | `REQ-05-01` (`REG-09`), `REQ-06-01` | `FLOW-03`, `FLOW-14` | Commercial Trade License, TIN certificate & MoTRI verification workbench for Fleet Owners and Corporate Renters. |
| `SPEC-ADM-OPS-04` | [`04-fleet-compliance-approval-queue.md`](./operations-queues/04-fleet-compliance-approval-queue.md) | `REQ-05-02` (`FL-01`–`FL-27`), `AD-03` | `FLOW-05`, `FLOW-13` | Machinery & Logistics listing audit: Title Deed (`Libre`), EAV spec completeness, 5-angle + baseline meter photos. |
| `SPEC-ADM-OPS-05` | [`05-operator-certification-queue.md`](./operations-queues/05-operator-certification-queue.md) | `REQ-04`, `FL-38`, `IN-03` | `FLOW-04`, `FLOW-13` | Certified Operator heavy machinery driving license & category endorsement verification workbench. |
| `SPEC-ADM-OPS-06` | [`06-marketplace-integrity-abuse-triage.md`](./operations-queues/06-marketplace-integrity-abuse-triage.md) | `REQ-05-24` (`MI-01`–`MI-07`) | `FLOW-08`, `FLOW-13` | Anti-circumvention chat redaction queue, duplicate SHA-256 bank slip velocity fraud, Sybil clusters, 3-strike suspensions. |
| `SPEC-ADM-OPS-07` | [`07-insurance-document-expiry-monitor.md`](./operations-queues/07-insurance-document-expiry-monitor.md) | `REQ-05-13` (`IN-01`–`IN-08`) | `FLOW-13`, `FLOW-19` | Automated 30/14/7-day insurance & `Libre` expiry compliance tracker, auto-suspended fleet quarantine & reactivation. |

### 2.3 Booking Lifecycle & Logistics Operations (`docs/admin/booking-operations/`)
| Spec ID | Specification Document | Covered SRS Requirements | Covered Flows & ADRs | Primary Focus |
|---|---|---|---|---|
| `SPEC-ADM-BKG-01` | [`01-booking-operations-console.md`](./booking-operations/01-booking-operations-console.md) | `REQ-05-05`, `REQ-08`, `AD-06` | `FLOW-07`, `FLOW-13` | Keyset-paginated (`nextCursor`/`prevCursor`) `@tanstack/react-table` booking register across all 16/20 states. |
| `SPEC-ADM-BKG-02` | [`02-booking-lifecycle-timeline.md`](./booking-operations/02-booking-lifecycle-timeline.md) | `REQ-05-05`, `REQ-08`, `AD-07` | `FLOW-07`, `FLOW-09`, `FLOW-10` | Forensic booking inspector: immutable `booking_audit_logs` timeline, Valkey hold status, Pre/Post inspections, ledger legs. |
| `SPEC-ADM-BKG-03` | [`03-logistics-dispatch-management.md`](./booking-operations/03-logistics-dispatch-management.md) | `REQ-05-08` (`LG-01`–`LG-09`) | `FLOW-11`, `FLOW-18` | Parent-Child transport sub-booking monitor (`parent_booking_id`), `AWAITING_FLEET_READINESS` & Reverse Cascade console. |
| `SPEC-ADM-BKG-04` | [`04-manual-intervention-modal.md`](./booking-operations/04-manual-intervention-modal.md) | `REQ-05-05`, `REQ-08`, `AD-08` | `FLOW-07`, `FLOW-13` | Dual-authorized administrative state override modal with OCC `version` lock, mandatory reason, and atomic ledger reversal. |

### 2.4 Dispute Arbitration & Damage Lightbox (`docs/admin/dispute-arbitration/`)
| Spec ID | Specification Document | Covered SRS Requirements | Covered Flows & ADRs | Primary Focus |
|---|---|---|---|---|
| `SPEC-ADM-DSP-01` | [`01-dispute-resolution-workbench.md`](./dispute-arbitration/01-dispute-resolution-workbench.md) | `REQ-05-05` (`BK-13`), `REQ-05-07` | `FLOW-08`, `FLOW-10`, `FLOW-13` | Payment slip & return damage dispute triage queue with 72h SLA timers and scoped (`owner-only`/`renter-only`) chat. |
| `SPEC-ADM-DSP-02` | [`02-split-screen-damage-lightbox.md`](./dispute-arbitration/02-split-screen-damage-lightbox.md) | `REQ-05-06`, `REQ-05-07`, `AD-04` | `FLOW-10`, `FLOW-13` | Dual-pane synchronized pan/zoom ($100\%\text{–}400\%$) comparator: Pre-Rental Check-In vs Post-Rental Check-Out photos & meters. |
| `SPEC-ADM-DSP-03` | [`03-escrow-arbitration-decision.md`](./dispute-arbitration/03-escrow-arbitration-decision.md) | `REQ-05-05` (`BK-13`), `RC-05`, `CW-05` | `FLOW-08`, `FLOW-10`, `ADR-010` | Binding arbitration ruling execution panel (`DISPUTE_RESOLVED_RENTER`/`OWNER`, `RULE_PARTIAL_SPLIT`) inside atomic `pgx.Tx`. |

### 2.5 Finance, Double-Entry Ledger & `OwnerWallet` (`docs/admin/finance-ledger/`)
| Spec ID | Specification Document | Covered SRS Requirements | Covered Flows & ADRs | Primary Focus |
|---|---|---|---|---|
| `SPEC-ADM-FIN-01` | [`01-bank-slip-reconciliation-queue.md`](./finance-ledger/01-bank-slip-reconciliation-queue.md) | `REQ-05-05`, `REQ-05-09` (`CW-03`, `CW-04`) | `FLOW-08`, `FLOW-12` | Bank wire receipt audit queue for `OwnerWallet` top-ups (`CW-04`) & disputed booking slips (`BK-13`), SHA-256 dup check. |
| `SPEC-ADM-FIN-02` | [`02-double-entry-ledger-browser.md`](./finance-ledger/02-double-entry-ledger-browser.md) | `REQ-05-05` (`BK-17-BE`), `REQ-05-09` | `FLOW-12`, `ADR-002`, `ADR-010` | Immutable `financial_ledger_entries` journal browser (`1010`, `2010`, `3010`, `2020`) with $\sum \text{amount\_minor} = 0$ verifier. |
| `SPEC-ADM-FIN-03` | [`03-owner-payout-batch-processing.md`](./finance-ledger/03-owner-payout-batch-processing.md) | `REQ-05-09` (`CW-11`), `07-01`, `07-02` | `FLOW-12`, `ADR-010` | Phase 1 `OwnerWallet` unused balance refund approvals (`CW-11`) & Phase 3 Gateway Escrow B2C/RTGS batch disbursements. |
| `SPEC-ADM-FIN-04` | [`04-withholding-tax-management.md`](./finance-ledger/04-withholding-tax-management.md) | `REQ-05-09` (`CW-12`), `LOC-01-BE` | `FLOW-12`, `ADR-015` | Ethiopia ERCA (1500 bps VAT, 200 bps WHT) & Kenya KRA (1600 bps VAT) tax invoice & WHT certificate compliance export. |
| `SPEC-ADM-FIN-05` | [`05-platform-commission-analytics.md`](./finance-ledger/05-platform-commission-analytics.md) | `REQ-05-09`, `REQ-07-09` | `FLOW-12`, `FLOW-13` | Realized commission revenue vs negative-wallet receivables, category take-rate margins, `ET` vs `KE` financial cohorts. |
| `SPEC-ADM-FIN-06` | [`06-owner-wallet-topup-credit-lockout.md`](./finance-ledger/06-owner-wallet-topup-credit-lockout.md) | `REQ-05-09` (`CW-03`–`CW-09`) | `FLOW-12` | `OwnerWallet` credit limit monitor (`-5,000.00 ETB` floor `CW-06`), automated fleet delisting quarantine (`CW-07`), top-up unlock. |

### 2.6 Corporate Governance & Enterprise B2B (`docs/admin/corporate-governance/`)
| Spec ID | Specification Document | Covered SRS Requirements | Covered Flows & ADRs | Primary Focus |
|---|---|---|---|---|
| `SPEC-ADM-CRP-01` | [`01-corporate-tenant-provisioning-hierarchy.md`](./corporate-governance/01-corporate-tenant-provisioning-hierarchy.md) | `REQ-06-01`, `REQ-06-02` | `FLOW-14` | Enterprise tenant onboarding, multi-tier org hierarchy builder, custom approval workflow matrix & threshold configurator. |
| `SPEC-ADM-CRP-02` | [`02-corporate-credit-underwriting-net30-ledger.md`](./corporate-governance/02-corporate-credit-underwriting-net30-ledger.md) | `REQ-06-03` | `FLOW-14`, `ADR-010` | Corporate revolving credit line underwriting, custom rate-card overrides, `NET-30` monthly invoice generation & aging dunning. |
| `SPEC-ADM-CRP-03` | [`03-enterprise-fleet-bulk-import-analytics.md`](./corporate-governance/03-enterprise-fleet-bulk-import-analytics.md) | `REQ-06-04`, `REQ-06-05` | `FLOW-05`, `FLOW-14` | Bulk CSV/Excel fleet onboarding validator for enterprise owners, multi-yard depot mapping, cross-project SLA reports. |

### 2.7 Dynamic Metadata & Schema Engine (`docs/admin/schema-engine/`)
| Spec ID | Specification Document | Covered SRS Requirements | Covered Flows & ADRs | Primary Focus |
|---|---|---|---|---|
| `SPEC-ADM-SCH-01` | [`01-dynamic-ui-engine-overview.md`](./schema-engine/01-dynamic-ui-engine-overview.md) | `REQ-05-11` (`AD-06-BE`), `REQ-11` | `FLOW-13`, `ADR-008` | End-to-end architecture of `/api/v1/admin/ui-config`, metadata schema tables, Valkey caching & instant invalidation. |
| `SPEC-ADM-SCH-02` | [`02-eav-category-field-builder.md`](./schema-engine/02-eav-category-field-builder.md) | `REQ-05-02` (`FL-13`–`FL-22`), `AD-05` | `FLOW-05`, `FLOW-13` | Visual 3-level Category & EAV `SpecFieldGroup`/`SpecFieldDefinition` schema builder with live mobile form preview. |
| `SPEC-ADM-SCH-03` | [`03-admin-table-column-customizer.md`](./schema-engine/03-admin-table-column-customizer.md) | `REQ-05-11` (`AD-06-WEB`), `REQ-11` | `FLOW-13` | Server-driven table column, formatter (`MONEY_MINOR`, `STATUS_BADGE`), sort key, and RBAC row action configurator. |
| `SPEC-ADM-SCH-04` | [`04-notification-template-editor.md`](./schema-engine/04-notification-template-editor.md) | `REQ-05-10`, `REQ-05-16`, `REQ-11` | `FLOW-17`, `ADR-013` | 6-locale Mustache notification template IDE (`SMS`, `FCM_PUSH`, `EMAIL`, `IN_APP`) with Ge'ez SMS segment counter. |

### 2.8 Multi-Country Parameters & Platform Configuration (`docs/admin/configuration/`)
| Spec ID | Specification Document | Covered SRS Requirements | Covered Flows & ADRs | Primary Focus |
|---|---|---|---|---|
| `SPEC-ADM-CFG-01` | [`01-multi-country-parameters.md`](./configuration/01-multi-country-parameters.md) | `REQ-05-23`, `REQ-07-08`, `LOC-01-BE` | `FLOW-13`, `ADR-015` | `country_configs` console (`ET`/`KE` VAT & WHT basis points, currency, dialing prefix, hold duration, credit limits). |
| `SPEC-ADM-CFG-02` | [`02-ethiopian-calendar-holidays.md`](./configuration/02-ethiopian-calendar-holidays.md) | `REQ-05-23` (`HC-01`–`HC-05`) | `FLOW-07`, `FLOW-13` | Ethiopian Ge'ez & Kenyan public holiday registry, moveable feast calculator, non-billable day & overtime multiplier rules. |
| `SPEC-ADM-CFG-03` | [`03-feature-flags-remote-config.md`](./configuration/03-feature-flags-remote-config.md) | `REQ-05-21` (`FF-01`–`FF-05`), `REQ-11` | `FLOW-13` | Canary percentage rollouts, jurisdiction/persona targeting rules, and instant operational kill-switches. |
| `SPEC-ADM-CFG-04` | [`04-audit-log-security-monitoring.md`](./configuration/04-audit-log-security-monitoring.md) | `REQ-05-11` (`AD-09`), `REQ-05-24` | `FLOW-13`, `FLOW-APPENDIX-B` | Immutable `admin_audit_logs` explorer with side-by-side JSON Diff Viewer (`previous_state` vs `new_state`) and `traceId`. |
| `SPEC-ADM-CFG-05` | [`05-policy-versioning-soft-delete-trash.md`](./configuration/05-policy-versioning-soft-delete-trash.md) | `REQ-05-19`, `REQ-05-20`, `REQ-05-22` | `FLOW-13` | Versioned legal policy publisher (`PT-01`), 30-day account deletion grace queue (`DL-01`), soft-delete (`deleted_at`) restore console. |

### 2.9 Admin Quality Assurance & Stress Benchmarks (`docs/admin/quality/`)
| Spec ID | Specification Document | Covered SRS Requirements | Covered Flows & ADRs | Primary Focus |
|---|---|---|---|---|
| `SPEC-ADM-QAS-01` | [`01-admin-qa-test-matrix.md`](./quality/01-admin-qa-test-matrix.md) | `REQ-09` | All Admin Flows | Desktop/Tablet browser matrix, Playwright E2E regression suite, Server Action contract & idempotency tests. |
| `SPEC-ADM-QAS-02` | [`02-accessibility-wcag-aa-audit.md`](./quality/02-accessibility-wcag-aa-audit.md) | `REQ-05-16` (`AC-01`–`AC-07`) | All Admin Flows | `@axe-core/playwright` zero-critical gate, WCAG 2.2 Level AA contrast, full keyboard grid navigation & focus traps. |
| `SPEC-ADM-QAS-03` | [`03-data-grid-stress-benchmark.md`](./quality/03-data-grid-stress-benchmark.md) | `REQ-09` | `FLOW-13` | 50,000-row `@tanstack/react-virtual` 60fps scroll benchmark, $<100\text{ms}$ cursor keyset pagination & facet filter SLA. |
