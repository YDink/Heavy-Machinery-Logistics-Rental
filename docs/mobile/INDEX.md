# Heavy Machinery & Logistics Rental Aggregator — Dink Machinery
## Mobile Client Subsystem Specifications — Master Index & Traceability Matrix (`docs/mobile/INDEX.md`)

**Document Identifier:** `INDEX-MOB-DINK-CORE-V7.0`  
**Classification:** Canonical Subsystem UI/UX, State Topology & Technical Behavior Index  
**Governing Documents:** [`Dink_Machinery_UI_UX_Design_Brief.md`](../source/Dink_Machinery_UI_UX_Design_Brief.md), [`00-design-system.md`](../shared/design/00-design-system.md), [`design-tokens.json`](../shared/design/design-tokens.json), [`mobile/DESIGN.md`](../components/mobile/DESIGN.md), [`contracts/openapi.yaml`](../../contracts/openapi.yaml)  
**Target Runtime:** Flutter 3.47.9 / Dart 3.9 Cross-Platform Mobile Client (Android 8.0+ API 26+, iOS 15.0+)  
**Supported Personas:** Civil Engineering Contractor (*Renter*), Fleet Owner & Yard Manager (*Owner*), Certified Heavy Equipment Operator (*Operator*), Lowboy Hauler (*Transporter*), Corporate Procurement & Approver (*Corporate*)  
**Core Architecture:** Feature-First Clean Architecture + `flutter_bloc` + `freezed` State Unions + `Isar` NoSQL Encrypted Offline Outbox + `Dio` 4-Tier Interceptor Stack  
**Total Screen & Subsystem Specifications:** 65 Specifications across 10 Bounded Modules

---

## 1. Architectural Invariants Enforced Across All Mobile Specifications

1. **Non-Custodial Financial Settlement (`ADR-010`, `SRS 05-05`, `05-09`):**
   - **Phase 1 MVP:** Strictly non-custodial under National Bank of Ethiopia (NBE) directives. Renters transfer rental funds directly to the Fleet Owner's verified bank coordinates (`CBE`, `Awash`, `Abyssinia`, `Dashen`, `Equity`, `KCB`) and upload the payment slip (`BK-07`). Security deposits are **informational-only** and settled off-platform (`BK-04`). Fleet Owners maintain a pre-funded `OwnerWallet` (`CW-01`–`CW-12`) with a `-5,000.00 ETB` (`-500,000 Santims`) credit buffer (`CW-06`) from which the 10% platform commission is atomically debited (`BK-17-BE`).
   - **Phase 2 Corporate (`SRS 06-01`–`06-05`):** Adds multi-tier organizational approval chains (`AWAITING_CORP_APPROVAL`), revolving corporate credit lines, and `NET-30` consolidated monthly invoicing.
   - **Phase 3 Advanced (`SRS 07-01`–`07-09`):** Adds in-app payment gateways (Telebirr, CBEBirr, M-Pesa STK Push) with automated commission deduction and third-party logistics/operator bidding marketplaces.
2. **Industrial Ergonomics & Outdoor High-Contrast Mode (`SRS 05-16`, `DS-DINK-CORE-V7.0`):**
   - Minimum **$56 \times 56\text{ dp}$** glove-friendly touch targets (`{dimension.space.56}`) and **$64\text{ dp}$** primary field actuators (`{dimension.space.64}`) with $\ge 8\text{ dp}$ separation.
   - Ambient light sensor auto-switching ($> 80,000\text{ lux}$ direct sunlight) to **Outdoor High-Contrast Mode** (`outdoorHighContrast` pure `#000000` AMOLED canvas, `#FFFFFF` text, `#FBBF24` glare yellow, $2\text{ dp}$ solid borders, WCAG 2.2 Level AAA $\ge 7:1$ / $21:1$ contrast).
3. **Monotonic Hour Meter & ClamAV Media Pipeline (`SRS 05-06`, `05-07`):**
   - Field inspections enforce live camera capture via a $320 \times 120\text{ dp}$ viewfinder reticle (gallery uploads disabled, accelerometer drift $< 0.05g$), background Dart Isolate compression ($1600\text{px}$ WebP at $82\%$ quality), and ClamAV `zINSTREAM` scanning.
   - Strict non-regression inequality: $\text{Baseline Meter} \le \text{Pre-Rental Check-In Meter} \le \text{Post-Rental Check-Out Meter}$ (`INSPECT_METER_REGRESSION` HTTP 422).
4. **Offline-First Local Execution Perimeter (`SRS 05-17`, `ADR-010`):**
   - Field inspections, walkaround photos, daily shift logs, and HMAC-SHA256 dynamic QR handshakes execute 100% offline in encrypted `Isar` NoSQL collections and replay idempotently (`Idempotency-Key == localDraftId`) upon reconnection.

---

## 2. Master Specification Catalog & Requirement Traceability Matrix

### 2.1 Foundations & Industrial Ergonomics (`docs/mobile/foundations/`)
| Spec ID | Specification Document | Covered SRS Requirements | Covered Flows & ADRs | Primary Focus |
|---|---|---|---|---|
| `SPEC-MOB-FND-01` | [`01-industrial-ergonomics.md`](./foundations/01-industrial-ergonomics.md) | `REQ-05-16` (`AC-01`–`AC-04`) | `FLOW-ENTRY`, `ADR-006` | $56\text{dp}$/$64\text{dp}$ glove touch targets, $>80,000\text{ lux}$ auto-switching, one-handed thumb zones. |
| `SPEC-MOB-FND-02` | [`02-color-typography-tokens.md`](./foundations/02-color-typography-tokens.md) | `REQ-05-16` (`AC-05`–`AC-07`) | `ADR-013` | W3C DTCG token bindings, `outdoorHighContrast`, OTA `Noto Sans Ethiopic` ($1.38\times$ line height) & `Space Mono`. |
| `SPEC-MOB-FND-03` | [`03-haptics-motion-feedback.md`](./foundations/03-haptics-motion-feedback.md) | `REQ-05-16`, `REQ-05-18` | `FLOW-07`, `FLOW-09` | Tactile haptic pulses, $85\text{dB}$ acoustic jobsite chimes, $<0.05g$ accelerometer shutter stability gate. |
| `SPEC-MOB-FND-04` | [`04-offline-outbox-architecture.md`](./foundations/04-offline-outbox-architecture.md) | `REQ-05-17` (`OF-01`–`OF-07`) | `FLOW-APPENDIX-A`, `ADR-010` | Encrypted Isar NoSQL outbox, Isolate WebP compression, Tus chunked uploads, `409 Conflict` resolution. |

### 2.2 Dynamic App Shell & Navigation (`docs/mobile/navigation/`)
| Spec ID | Specification Document | Covered SRS Requirements | Covered Flows & ADRs | Primary Focus |
|---|---|---|---|---|
| `SPEC-MOB-NAV-01` | [`01-role-switcher-app-shell.md`](./navigation/01-role-switcher-app-shell.md) | `REQ-04`, `REQ-05-01` | `FLOW-01`, `FLOW-02` | Seamless persona switching (Renter $\leftrightarrow$ Owner $\leftrightarrow$ Operator $\leftrightarrow$ Transporter) with isolated BLoC trees. |
| `SPEC-MOB-NAV-02` | [`02-bottom-navigation-drawer.md`](./navigation/02-bottom-navigation-drawer.md) | `REQ-04`, `REQ-05-10` | `FLOW-01`, `FLOW-17` | Role-parameterized bottom navigation bars, unread badge sync, offline outbox status drawer. |

### 2.3 Authentication, Onboarding & KYC (`docs/mobile/auth-kyc/`)
| Spec ID | Specification Document | Covered SRS Requirements | Covered Flows & ADRs | Primary Focus |
|---|---|---|---|---|
| `SPEC-MOB-AUTH-01` | [`01-splash-language-selection.md`](./auth-kyc/01-splash-language-selection.md) | `REQ-05-01`, `REQ-05-16` | `FLOW-01`, `ADR-013`, `ADR-015` | Country SIM detection (`ET`/`KE`), 6-locale selector (`en`, `am`, `om`, `ti`, `so`, `sw`), OTA font pre-fetch. |
| `SPEC-MOB-AUTH-02` | [`02-passwordless-phone-otp.md`](./auth-kyc/02-passwordless-phone-otp.md) | `REQ-05-01` (`REG-01`–`REG-04`) | `FLOW-02`, `FLOW-APPENDIX-B` | E.164 phone input, 6-digit SMS autofill, 60s resend cooldown, ES256 JWT + RTR family rotation. |
| `SPEC-MOB-AUTH-03` | [`03-persona-role-selection.md`](./auth-kyc/03-persona-role-selection.md) | `REQ-04`, `REQ-05-01`, `REQ-05-14` | `FLOW-01`, `FLOW-02` | Interactive role selection cards with progressive KYC requirement checklists per persona. |
| `SPEC-MOB-AUTH-04` | [`04-fayda-national-id-kyc.md`](./auth-kyc/04-fayda-national-id-kyc.md) | `REQ-05-01` (`REG-05`–`REG-08`) | `FLOW-03`, `ADR-011` | Pluggable Fayda FIN / Kebele / Kenyan Huduma ID camera capture, glare detection, OCR verification. |
| `SPEC-MOB-AUTH-05` | [`05-trade-license-business-onboarding.md`](./auth-kyc/05-trade-license-business-onboarding.md) | `REQ-05-01` (`REG-09`), `REQ-06-01` | `FLOW-03`, `FLOW-14` | Commercial Trade License, TIN certificate & MoTRI registration upload for Owners and Corporate Renters. |
| `SPEC-MOB-AUTH-06` | [`06-phone-migration-recovery.md`](./auth-kyc/06-phone-migration-recovery.md) | `REQ-05-01` (`REG-14`), `REQ-05-24` | `FLOW-15`, `FLOW-APPENDIX-B` | Secure phone number migration flow with dual OTP, National ID match, and liveness selfie challenge. |

### 2.4 Equipment Renter Experience (`docs/mobile/renter-experience/`)
| Spec ID | Specification Document | Covered SRS Requirements | Covered Flows & ADRs | Primary Focus |
|---|---|---|---|---|
| `SPEC-MOB-RNT-01` | [`01-renter-home-discovery.md`](./renter-experience/01-renter-home-discovery.md) | `REQ-05-04` (`FD-01`–`FD-05`) | `FLOW-06` | 3-level category taxonomy grid, featured nearby machines, active hold banner, quick search header. |
| `SPEC-MOB-RNT-02` | [`02-postgis-radius-search-map.md`](./renter-experience/02-postgis-radius-search-map.md) | `REQ-05-04` (`FD-06`–`FD-09`) | `FLOW-06` | PostGIS `ST_DWithin` radius slider ($5\text{–}200\text{ km}$), spatial cluster markers, preview carousel. |
| `SPEC-MOB-RNT-03` | [`03-corridor-route-filtering.md`](./renter-experience/03-corridor-route-filtering.md) | `REQ-05-04` (`FD-10`–`FD-13`) | `FLOW-06`, `FLOW-11` | Origin-to-destination corridor filter, dynamic EAV spec range facets, wet/dry lease toggles. |
| `SPEC-MOB-RNT-04` | [`04-fleet-detail-spec-sheet.md`](./renter-experience/04-fleet-detail-spec-sheet.md) | `REQ-05-02`, `REQ-05-03`, `FD-14` | `FLOW-05`, `FLOW-06` | 5–20 photo gallery, Dynamic EAV Spec Accordion, Baseline Hour Meter badge, Compatible Attachments. |
| `SPEC-MOB-RNT-05` | [`05-side-by-side-comparison.md`](./renter-experience/05-side-by-side-comparison.md) | `REQ-05-04` (`FD-15`–`FD-17`) | `FLOW-06` | Synchronized 2-to-4 machine comparison matrix highlighting EAV spec differences and minor-unit rates. |
| `SPEC-MOB-RNT-06` | [`06-booking-reservation-flow.md`](./renter-experience/06-booking-reservation-flow.md) | `REQ-05-05` (`BK-01`–`BK-03`) | `FLOW-07` | Dual Gregorian/Ge'ez date range selector, bundled demountable attachment picker (`BK-03`), job site pin. |
| `SPEC-MOB-RNT-07` | [`07-checkout-price-breakdown.md`](./renter-experience/07-checkout-price-breakdown.md) | `REQ-05-05` (`BK-04`, `BK-06B`) | `FLOW-07`, `ADR-002`, `ADR-010` | Minor-unit tariff breakdown (Base + Attachments + 15%/16% VAT) + **Informational-Only Security Deposit notice**. |
| `SPEC-MOB-RNT-08` | [`08-60-minute-hold-timer.md`](./renter-experience/08-60-minute-hold-timer.md) | `REQ-05-05` (`BK-05`, `BK-06`) | `FLOW-07`, `FLOW-16` | 60-minute Valkey hold countdown ring (`Emerald -> Amber -> Crimson`), copyable Owner Bank Coordinates card. |
| `SPEC-MOB-RNT-09` | [`09-bank-slip-upload-verification.md`](./renter-experience/09-bank-slip-upload-verification.md) | `REQ-05-05` (`BK-07`, `BK-08`, `BK-12`) | `FLOW-08` | Non-custodial bank wire slip upload, ClamAV scan, SHA-256 duplicate check, Renter Decline Remediation (`BK-12`). |
| `SPEC-MOB-RNT-10` | [`10-logistics-sub-booking-modal.md`](./renter-experience/10-logistics-sub-booking-modal.md) | `REQ-05-08` (`LG-01`–`LG-08`) | `FLOW-11`, `FLOW-18` | Linked Transport Sub-Booking (`parent_booking_id`), payload-matched lowboy suggestions, Reverse Cascade alerts. |
| `SPEC-MOB-RNT-11` | [`11-corporate-approval-pending-view.md`](./renter-experience/11-corporate-approval-pending-view.md) | `REQ-06-02` (`CA-01`–`CA-08`) | `FLOW-14` | `AWAITING_CORP_APPROVAL` multi-tier approval chain tracker, SLA timers, cost-center budget tag. |
| `SPEC-MOB-RNT-12` | [`12-saved-searches-waitlist-alerts.md`](./renter-experience/12-saved-searches-waitlist-alerts.md) | `REQ-05-15` (`SD-01`–`SD-07`) | `FLOW-06` | Saved search subscriptions, "Notify Me When Available" waitlist alerts, smart alternative suggestions. |
| `SPEC-MOB-RNT-13` | [`13-active-rental-hud-extensions-breakdowns.md`](./renter-experience/13-active-rental-hud-extensions-breakdowns.md) | `REQ-05-06` (`AR-03`–`AR-05`) | `FLOW-09` | Active Rental Shift HUD, live hour meter gauge, Lease Extension request (`AR-04`), Breakdown Timer Pause (`AR-05`). |
| `SPEC-MOB-RNT-14` | [`14-post-rental-reconciliation-mutual-rating.md`](./renter-experience/14-post-rental-reconciliation-mutual-rating.md) | `REQ-05-07` (`RC-03`, `RC-05`–`RC-07`) | `FLOW-10` | Final Balance Reconciliation Statement (excess hours + damage), final wire upload, Blind 7-Day Mutual Rating. |

### 2.5 Corporate Renter & Enterprise Experience (`docs/mobile/corporate-experience/`)
| Spec ID | Specification Document | Covered SRS Requirements | Covered Flows & ADRs | Primary Focus |
|---|---|---|---|---|
| `SPEC-MOB-CRP-01` | [`01-corporate-org-hierarchy-cost-centers.md`](./corporate-experience/01-corporate-org-hierarchy-cost-centers.md) | `REQ-06-01` | `FLOW-14` | Multi-tier corporate division & project site switcher, cost-center budget bars, team RBAC management. |
| `SPEC-MOB-CRP-02` | [`02-corporate-requisition-approval-inbox.md`](./corporate-experience/02-corporate-requisition-approval-inbox.md) | `REQ-06-02` | `FLOW-14` | Mobile approver queue for Procurement/Finance Directors, one-tap approve/reject, authority delegation. |
| `SPEC-MOB-CRP-03` | [`03-corporate-credit-line-net30-invoices.md`](./corporate-experience/03-corporate-credit-line-net30-invoices.md) | `REQ-06-03` | `FLOW-14`, `ADR-010` | Revolving corporate credit line gauge, negotiated rate cards, `NET-30` monthly statement & WHT upload. |
| `SPEC-MOB-CRP-04` | [`04-enterprise-multi-yard-fleet-console.md`](./corporate-experience/04-enterprise-multi-yard-fleet-console.md) | `REQ-06-04` | `FLOW-05` | Multi-depot yard switcher for enterprise fleet owners, yard dispatcher assignment, maintenance schedules. |
| `SPEC-MOB-CRP-05` | [`05-corporate-cost-center-utilization-reports.md`](./corporate-experience/05-corporate-cost-center-utilization-reports.md) | `REQ-06-05` | `FLOW-14` | Project cost-center burn rate, active vs idle engine hour telemetry, exportable PDF/CSV executive reports. |

### 2.6 Fleet Owner Experience (`docs/mobile/owner-experience/`)
| Spec ID | Specification Document | Covered SRS Requirements | Covered Flows & ADRs | Primary Focus |
|---|---|---|---|---|
| `SPEC-MOB-OWN-01` | [`01-owner-dashboard-fleet-summary.md`](./owner-experience/01-owner-dashboard-fleet-summary.md) | `REQ-05-02`, `REQ-05-09` | `FLOW-05`, `FLOW-12` | Fleet utilization KPIs, active leases, maintenance alerts, pre-funded `OwnerWallet` credit limit bar. |
| `SPEC-MOB-OWN-02` | [`02-fleet-registration-stepper.md`](./owner-experience/02-fleet-registration-stepper.md) | `REQ-05-02` (`FL-01`–`FL-12`) | `FLOW-05` | Unified Machinery & Logistics onboarding wizard, 3-level category picker, VIN/Plate, `Libre` & Insurance upload. |
| `SPEC-MOB-OWN-03` | [`03-dynamic-eav-spec-form.md`](./owner-experience/03-dynamic-eav-spec-form.md) | `REQ-05-02` (`FL-13`–`FL-22`) | `FLOW-05` | Server-driven `SpecFieldGroup` accordion engine, missing-field counter pills, auto-scroll to error. |
| `SPEC-MOB-OWN-04` | [`04-machinery-360-photo-uploader.md`](./owner-experience/04-machinery-360-photo-uploader.md) | `REQ-05-02` (`FL-23`–`FL-27`) | `FLOW-05` | 5–20 photo capture with guided 4-angle + baseline hour meter reticles, Isolate WebP compression, ClamAV stream. |
| `SPEC-MOB-OWN-05` | [`05-attachment-inventory-management.md`](./owner-experience/05-attachment-inventory-management.md) | `REQ-05-03` (`AT-01`–`AT-11`) | `FLOW-05` | Demountable vs Integrated attachment registry, `compatible_fleet_ids` linking, bundled vs standalone rates. |
| `SPEC-MOB-OWN-06` | [`06-availability-calendar-pricing.md`](./owner-experience/06-availability-calendar-pricing.md) | `REQ-05-02` (`FL-28`–`FL-36`), `05-23` | `FLOW-05`, `FLOW-07` | Interactive monthly availability board (`FL-35`), maintenance blocks, Wet/Dry/Per-Km rates, holiday rules. |
| `SPEC-MOB-OWN-07` | [`07-booking-request-approval-sheet.md`](./owner-experience/07-booking-request-approval-sheet.md) | `REQ-05-05`, `FL-38` | `FLOW-07` | Incoming reservation review, contractor KYC badge, Certified Operator assignment (`FL-38`), readiness status. |
| `SPEC-MOB-OWN-08` | [`08-direct-payment-confirmation.md`](./owner-experience/08-direct-payment-confirmation.md) | `REQ-05-05` (`BK-09`, `BK-10`, `BK-17-BE`) | `FLOW-08`, `ADR-010` | Zoomable bank slip audit drawer, SHA-256 duplicate warning (`BK-07B`), atomic 10% `OwnerWallet` commission debit. |
| `SPEC-MOB-OWN-09` | [`09-owner-wallet-payout-request.md`](./owner-experience/09-owner-wallet-payout-request.md) | `REQ-05-09` (`CW-01`–`CW-12`), `07-02` | `FLOW-12`, `ADR-010` | Pre-funded `OwnerWallet` dashboard, `-5,000 ETB` credit limit bar, bank slip top-up (`CW-03`), balance refund (`CW-11`). |
| `SPEC-MOB-OWN-10` | [`10-payment-decline-remediation-dispute.md`](./owner-experience/10-payment-decline-remediation-dispute.md) | `REQ-05-05` (`BK-11`–`BK-13`) | `FLOW-08` | Structured payment decline modal (`AMOUNT_MISMATCH`, `ILLEGIBLE_SLIP`), Renter rebuttal & Admin dispute view. |
| `SPEC-MOB-OWN-11` | [`11-return-damage-claim-24h-docket.md`](./owner-experience/11-return-damage-claim-24h-docket.md) | `REQ-05-07` (`RC-04`–`RC-06`) | `FLOW-10` | 24-hour post-return itemized damage claim filing wizard, pre/post photo pairing, repair quote upload. |

### 2.7 Certified Field Operator Experience (`docs/mobile/operator-experience/`)
| Spec ID | Specification Document | Covered SRS Requirements | Covered Flows & ADRs | Primary Focus |
|---|---|---|---|---|
| `SPEC-MOB-OPR-01` | [`01-operator-job-dispatch-list.md`](./operator-experience/01-operator-job-dispatch-list.md) | `REQ-04`, `FL-38`, `GPS-04` | `FLOW-04`, `FLOW-09` | Owner-to-Operator linking invitation, assigned dispatch queue, explicit GPS tracking consent (`GPS-04`). |
| `SPEC-MOB-OPR-02` | [`02-check-in-walkaround-inspection.md`](./operator-experience/02-check-in-walkaround-inspection.md) | `REQ-05-06` (`AR-01`, `AR-02`) | `FLOW-09` | Pre-rental check-in stepper: 4 mandatory angle photos, fluids/safety/attachment pin checklist, damage pins. |
| `SPEC-MOB-OPR-03` | [`03-monotonic-meter-reticle-capture.md`](./operator-experience/03-monotonic-meter-reticle-capture.md) | `REQ-05-06` (`AR-01`), `RC-02` | `FLOW-09`, `FLOW-10` | Camera-only $320 \times 120\text{ dp}$ meter reticle, $<0.05g$ tilt check, monotonic validation, 5%-step fuel slider. |
| `SPEC-MOB-OPR-04` | [`04-peer-qr-handshake-screen.md`](./operator-experience/04-peer-qr-handshake-screen.md) | `REQ-05-06` (`AR-02`), `OF-04` | `FLOW-09`, `FLOW-APPENDIX-A` | Dual Challenge 6-Digit PIN & Offline HMAC-SHA256 Peer-to-Peer Dynamic QR Handshake. |
| `SPEC-MOB-OPR-05` | [`05-active-job-daily-telematics-log.md`](./operator-experience/05-active-job-daily-telematics-log.md) | `REQ-05-06` (`AR-03`, `AR-05`), `05-12` | `FLOW-09` | 30s WebSocket GPS streamer, $1,000\text{m}$ geofence breach alert, shift hour logger, breakdown pause trigger. |
| `SPEC-MOB-OPR-06` | [`06-check-out-return-inspection.md`](./operator-experience/06-check-out-return-inspection.md) | `REQ-05-07` (`RC-01`–`RC-03`) | `FLOW-10` | Post-rental check-out inspection stepper, baseline diff view, automated excess operating hour calculator. |
| `SPEC-MOB-OPR-07` | [`07-glass-digital-signature-canvas.md`](./operator-experience/07-glass-digital-signature-canvas.md) | `REQ-05-06` (`AR-08`) | `FLOW-09`, `FLOW-10` | Orientation-locked vector signature canvas ($\ge 20$ nodes), signee role tag, fallback Phone OTP attestation. |

### 2.8 Lowboy Transporter & Logistics Experience (`docs/mobile/transporter-experience/`)
| Spec ID | Specification Document | Covered SRS Requirements | Covered Flows & ADRs | Primary Focus |
|---|---|---|---|---|
| `SPEC-MOB-TRN-01` | [`01-transporter-job-radar.md`](./transporter-experience/01-transporter-job-radar.md) | `REQ-05-08` (`LG-01`–`LG-04`), `07-03` | `FLOW-11` | Phase 1 Linked Transport Sub-Booking requests & Phase 3 Corridor Haulage Job Radar by payload tonnage. |
| `SPEC-MOB-TRN-02` | [`02-transport-job-detail.md`](./transporter-experience/02-transport-job-detail.md) | `REQ-05-08` (`LG-03`–`LG-05`) | `FLOW-11` | Transport manifest, depot-to-site route clearance, machine weight/dimensions, mobilization delay notice. |
| `SPEC-MOB-TRN-03` | [`03-carrier-fleet-registration.md`](./transporter-experience/03-carrier-fleet-registration.md) | `REQ-05-02`, `REQ-05-08` (`LG-02`) | `FLOW-05`, `FLOW-11` | `Category = LOGISTICS` fleet onboarding: Lowboy, Flatbed, Crane Truck payload capacity, dimensions, axle count. |
| `SPEC-MOB-TRN-04` | [`04-pickup-handshake-verification.md`](./transporter-experience/04-pickup-handshake-verification.md) | `REQ-05-08` (`LG-06`) | `FLOW-11` | Depot machinery loading & chain tie-down safety checklist, 4-angle loaded trailer capture, release QR scan. |
| `SPEC-MOB-TRN-05` | [`05-site-delivery-confirmation.md`](./transporter-experience/05-site-delivery-confirmation.md) | `REQ-05-08` (`LG-07`, `LG-08`) | `FLOW-11`, `FLOW-18` | Job site offloading confirmation, GPS arrival check, contractor sign-off, Bidirectional Cancellation Cascade handling. |

### 2.9 Cross-Persona Shared Features (`docs/mobile/shared-features/`)
| Spec ID | Specification Document | Covered SRS Requirements | Covered Flows & ADRs | Primary Focus |
|---|---|---|---|---|
| `SPEC-MOB-SHR-01` | [`01-notification-inbox-preferences.md`](./shared-features/01-notification-inbox-preferences.md) | `REQ-05-10` (`NT-01`–`NT-09`) | `FLOW-17`, `FLOW-APPENDIX-A` | Priority-tiered notification center (`CRITICAL`, `HIGH`, `NORMAL`), deep-link router, quiet-hours preferences. |
| `SPEC-MOB-SHR-02` | [`02-insurance-policy-claim-filing.md`](./shared-features/02-insurance-policy-claim-filing.md) | `REQ-05-13` (`IN-01`–`IN-08`) | `FLOW-19` | Insurance & `Libre` vault, 30/14/7-day expiry banners, field incident dossier & claim filingstepper. |
| `SPEC-MOB-SHR-03` | [`03-dispute-arbitration-tracker.md`](./shared-features/03-dispute-arbitration-tracker.md) | `REQ-05-05` (`BK-13`), `REQ-05-07` | `FLOW-08`, `FLOW-10` | 72-hour SLA dispute timeline, rebuttal evidence uploader, scoped arbitrator messaging (`renter-only`/`owner-only`). |
| `SPEC-MOB-SHR-04` | [`04-pdf-contract-invoice-viewer.md`](./shared-features/04-pdf-contract-invoice-viewer.md) | `REQ-05-05`, `REQ-05-09` (`CW-12`) | `FLOW-07`, `FLOW-12` | Offline-cached PDF viewer for signed lease agreements, ERCA/KRA VAT invoices, and 2% WHT certificates. |
| `SPEC-MOB-SHR-05` | [`05-user-profile-security-settings.md`](./shared-features/05-user-profile-security-settings.md) | `REQ-05-01`, `REQ-05-19` (`DL-01`–`DL-06`) | `FLOW-APPENDIX-B` | WebAuthn Passkey/Biometric binding, active session revocation, 30-day cooling-off Account Deletion wizard. |
| `SPEC-MOB-SHR-06` | [`06-onboarding-walkthrough-contextual-help.md`](./shared-features/06-onboarding-walkthrough-contextual-help.md) | `REQ-05-14` (`ON-01`–`ON-07`) | `FLOW-01` | Interactive persona coach marks, contextual field tooltips, offline-cached safety manual & support escalation. |
| `SPEC-MOB-SHR-07` | [`07-policy-consent-ethiopian-calendar-picker.md`](./shared-features/07-policy-consent-ethiopian-calendar-picker.md) | `REQ-05-20` (`PT-01`–`PT-05`), `05-23` | `FLOW-01`, `FLOW-07` | Versioned legal policy re-consent gate (`SHA-256` receipt) + Dual Gregorian & Ethiopian Ge'ez 13-month date picker. |
| `SPEC-MOB-SHR-08` | [`08-anti-circumvention-qa-integrity-shields.md`](./shared-features/08-anti-circumvention-qa-integrity-shields.md) | `REQ-05-24` (`MI-01`–`MI-07`) | `FLOW-06`, `FLOW-08` | Pre-confirmation Q&A regex PII masking (blocking phone/email sharing), velocity fraud lockouts, strike banners. |

### 2.10 Mobile Quality Assurance & Chaos Verification (`docs/mobile/quality/`)
| Spec ID | Specification Document | Covered SRS Requirements | Covered Flows & ADRs | Primary Focus |
|---|---|---|---|---|
| `SPEC-MOB-QAS-01` | [`01-mobile-qa-test-matrix.md`](./quality/01-mobile-qa-test-matrix.md) | `REQ-09` | All Mobile Flows | Low-end Android ($2\text{GB}$ RAM) to iOS hardware matrix, $<1.8\text{s}$ cold start, `bloc_test` & Golden suite. |
| `SPEC-MOB-QAS-02` | [`02-accessibility-wcag-audit.md`](./quality/02-accessibility-wcag-audit.md) | `REQ-05-16` (`AC-01`–`AC-07`) | All Mobile Flows | WCAG 2.2 Level AAA Outdoor contrast ($21:1$), $56\text{dp}$ touch target linter, TalkBack/VoiceOver semantics. |
| `SPEC-MOB-QAS-03` | [`03-network-chaos-test-plan.md`](./quality/03-network-chaos-test-plan.md) | `REQ-05-17` (`OF-01`–`OF-07`) | `FLOW-APPENDIX-A`, `FLOW-APPENDIX-B` | 1,000-mutation Isar offline outbox chaos replay under 100% packet drop, RTR token replay revocation drill. |
