# REQUIREMENTS_TRACKER.md — Dink Machinery Master Implementation Tracker

**Classification:** Canonical Project Requirement & Partition Implementation Status Matrix  
**System:** Heavy Machinery & Logistics Rental Marketplace  
**Last Updated:** 2026-11-20T15:45:00+03:00  
**Latest Baseline Commit:** `41a098a` (`[PLATFORM-01] fix(db,mobile): cast queries params and configure android gradle`)  
**Active Phase:** Phase 2 (Partition Feature Development)  
**Active Partition:** `05-08-logistics-transport-integration` (Logistics Transport Integration, Sub-Bookings & Reverse Cascades)

---

## 1. Executive Progress Dashboard

### Status Metrics
- **Total MVP Partitions (Phase 2):** 24
- **Completed Partitions:** 7 (`05-01`, `05-02`, `05-03`, `05-04`, `05-05`, `05-06`, `05-07`) + Architecture Refactor (`Pure Go Stdlib`)
- **In Progress Partitions:** 1 (`05-08`)
- **Pending MVP Partitions:** 17 (`05-08` through `05-24`)
- **Phase 2 MVP Completion Rate:** **29.2%** (7 / 24 partitions complete)
- **Corporate Partitions (Phase 2 Extension):** 5 pending (`06-01` through `06-05`)
- **Advanced Scale Partitions (Phase 3):** 9 pending (`07-01` through `07-09`)

```text
Phase 2 MVP Progress: [███████░░░░░░░░░░░░░░░] 29.2% (7 / 24 Partitions Done)
Quality Gates:        100% Passing (Backend Suite, Mobile 84/84, Spectral OpenAPI 0 errors)
```

---

## 2. Requirement Lifecycle & Tracking Rule

### The Requirement Tracking Invariant
1. **Three Allowed Status States:**
   - `NOT_STARTED` (⏳): Requirement specification exists in `docs/shared/requirements/`, pending development branch.
   - `IN_PROGRESS` (🔄): Feature branch `feature/<partition-id>` created, active implementation or test work underway.
   - `COMPLETED` (✅): 100% test pass rate across all targeted subsystems, OpenAPI bundled, quality gates passed (`G0`-`G2`), and rebase-merged into `main`.
2. **Atomic Tracker Synchronization:**
   - Any commit introducing code, tests, or state transitions for a requirement partition **MUST co-stage both `CONTEXT.md` and `REQUIREMENTS_TRACKER.md` and document verified manual test scenarios under the `Manual-Tests:` footer per `docs/quality/MANUAL_TESTING_PROTOCOL.md`**.
   - Commit messages must record the updated partition status and reference the ticket in the header (e.g., `[RENT-01-BE]`).
3. **Subsystem Verification Standard:**
   A partition cannot be marked `COMPLETED` unless all applicable subsystems are implemented and verified:
   - **`DB` (Database):** Up/down SQL migrations, generated `sqlc` queries, and idempotent test fixtures in `scripts/fixtures/` ensuring multi-machine parity.
   - **`OAS` (Contract):** OpenAPI 3.1.0 schemas and endpoints validated via Spectral (0 errors).
   - **`BE` (Backend):** Pure Go standard library (`net/http.ServeMux` + `log/slog`), strict type & layer isolation (Entities vs DTOs vs DB models), 1-concept-per-file, with 100% unit/integration pass rate.
   - **`MOB` (Mobile):** UI designs generated just-in-time via Stitch MCP per requirement at development start (maintaining flow continuity from preceding screens); Flutter Clean Architecture + BLoC with strict layer segregation (zero entities in data sources), bilingual localization (EN/AM), 56dp touch targets, granular 1-concept-per-file separation, and passing widget tests.
   - **`ADM` (Admin Console):** Next.js App Router dynamic metadata tables / Server Actions (where in scope).

---

## 3. Phase 2 (MVP Core Marketplace) Partitions (`05-01` to `05-24`)

| Partition ID | Partition Name & Functional Scope | Requirement IDs | Subsystems Delivered | Status | Verified Tests | Completed Commit |
|:---:|---|---|:---:|:---:|:---:|:---:|
| **`05-01`** | **Registration & Authentication**<br>OTP auth, JWT session refresh, phone carrier detection, dynamic RBAC. | `AUTH-01` to `AUTH-12`<br>`REG-01-MOB` | `DB`, `OAS`, `BE`, `MOB` | ✅ `COMPLETED` | 100% unit + integration | `aa04956`<br>`8e71af2` |
| **`05-02`** | **Fleet Registration (Machinery/Logistics)**<br>Equipment onboarding, multi-angle photo reticles, category specs, VIN validation. | `FL-01` to `FL-25`<br>`REQ-05-02` | `DB`, `OAS`, `BE`, `MOB` | ✅ `COMPLETED` | 100% unit + integration | `23627f2`<br>`98f8a7b` |
| **`05-03`** | **Attachment Registration & Inventory**<br>Breaker/bucket inventory, GiST exclusion anti-collision constraints, bundled tariffs. | `FL-31` to `FL-36`<br>`REQ-05-03` | `DB`, `OAS`, `BE`, `MOB` | ✅ `COMPLETED` | 100% unit + concurrency | `e26c1e1` |
| **`05-04`** | **Fleet Discovery, Search & Comparison**<br>PostGIS spatial radius search, trigram filtering, dynamic 3-unit comparison matrix. | `SR-01` to `SR-14`<br>`REQ-05-04` | `DB`, `OAS`, `BE`, `MOB` | ✅ `COMPLETED` | 100% unit + spatial | `c1466e2` |
| **`REFACTOR`**| **Pure Go Standard Library Architecture**<br>Zero web frameworks; complete removal of Chi and Zap; native `net/http.ServeMux` + `log/slog`. | `ARCH-01-BE`<br>`REF-01-BE` | `BE` | ✅ `COMPLETED` | 100% backend suite | `de4c378` |
| **`05-05`** | **Booking Reservation Flow**<br>16-state lifecycle FSM, 60-min Valkey hold lock, 50% prepay quoting, bank slip evidence. | `BK-01` to `BK-18`<br>`BK-01-BE` | `DB`, `OAS`, `BE`, `MOB` | ✅ `COMPLETED` | 100% unit, FSM, widgets | `3a23e83` |
| **`05-06`** | **Active Rental Management & Telemetry**<br>5-angle inspection, 6-digit challenge PIN, 8.0 hr shift quota HUD, breakdown timer pause, extensions. | `AR-01` to `AR-05`<br>`RENT-01-BE` | `DB`, `OAS`, `BE`, `MOB` | ✅ `COMPLETED` | 100% unit, HUD, BLoC (72 tests) | `3010010` |
| **`05-07`** | **Rental Completion, Return & Damage**<br>Return inspection, damage triage, excess shift hours billing, deposit release/forfeiture. | `RET-01` to `RET-12` | `DB`, `OAS`, `BE`, `MOB` | ✅ `COMPLETED` | 100% unit, BLoC, widgets, pure Go handlers | `c0e45ab` |
| **`05-08`** | **Logistics & Transport Integration**<br>Lowboy carrier dispatch, freight estimation, transporter coordination. | `LOG-01` to `LOG-10` | `DB`, `OAS`, `BE`, `MOB` | 🔄 `IN_PROGRESS` | 100% backend unit/integration + live HTTP probes | `feature/05-08-logistics-transport-integration` |
| **`05-09`** | **Commission Revenue & Owner Payouts**<br>Double-entry financial ledger, platform commission deduction, owner escrow wallets. | `PAY-01` to `PAY-15` | `DB`, `OAS`, `BE`, `MOB` | ⏳ `NOT_STARTED` | — | — |
| **`05-10`** | **In-App Notification Center**<br>Multi-channel push, SMS, and in-app event notification inbox. | `NOTIF-01` to `NOTIF-08` | `DB`, `OAS`, `BE`, `MOB` | ⏳ `NOT_STARTED` | — | — |
| **`05-11`** | **Admin Dynamic UI Engine**<br>Server-driven dynamic admin panel, approval workflows, dispute resolution console. | `ADM-01` to `ADM-14` | `DB`, `OAS`, `BE`, `ADM` | ⏳ `NOT_STARTED` | — | — |
| **`05-12`** | **Operator GPS Tracking**<br>Periodic telematics GPS pings, breadcrumb trail, and geofence deviation alerts. | `GPS-01` to `GPS-08` | `DB`, `OAS`, `BE`, `MOB` | ⏳ `NOT_STARTED` | — | — |
| **`05-13`** | **Insurance Document Management**<br>Coverage policy verification, document vault, expiry alarm workflows. | `INS-01` to `INS-06` | `DB`, `OAS`, `BE`, `MOB` | ⏳ `NOT_STARTED` | — | — |
| **`05-14`** | **User Onboarding & Help System**<br>Interactive field tutorials, multilingual FAQs, and inline contextual help. | `HELP-01` to `HELP-06` | `MOB`, `ADM` | ⏳ `NOT_STARTED` | — | — |
| **`05-15`** | **Search & Discovery Enhancements**<br>Advanced machinery facet filtering, popular searches, saved searches. | `SR-15` to `SR-20` | `DB`, `OAS`, `BE`, `MOB` | ⏳ `NOT_STARTED` | — | — |
| **`05-16`** | **Accessibility & Internationalization**<br>WCAG 2.2 AAA audit, Ge'ez typography compliance, dynamic locale OTA bundles. | `A11Y-01` to `A11Y-10` | `MOB`, `ADM` | ⏳ `NOT_STARTED` | — | — |
| **`05-17`** | **Offline Low-Connectivity Support**<br>Isar NoSQL offline mutation queue, outbox synchronization worker. | `OFF-01` to `OFF-08` | `MOB` | ⏳ `NOT_STARTED` | — | — |
| **`05-18`** | **Error Handling & User Feedback**<br>Unified RFC 9457 error presentation, offline banners, self-service diagnostics. | `ERR-01` to `ERR-08` | `OAS`, `BE`, `MOB` | ⏳ `NOT_STARTED` | — | — |
| **`05-19`** | **Account Deletion & Data Retention**<br>GDPR/statutory right-to-erasure workflows with 10-year immutable ledger hold. | `DEL-01` to `DEL-06` | `DB`, `BE`, `MOB`, `ADM` | ⏳ `NOT_STARTED` | — | — |
| **`05-20`** | **Policy & Terms Management**<br>Click-wrap versioned contract acceptance, legal audit logging. | `TOS-01` to `TOS-05` | `DB`, `OAS`, `BE`, `MOB` | ⏳ `NOT_STARTED` | — | — |
| **`05-21`** | **Feature Flagging & Remote Configuration**<br>Targeted runtime toggles, dark launches, and kill switches. | `FLAG-01` to `FLAG-06` | `DB`, `OAS`, `BE`, `MOB` | ⏳ `NOT_STARTED` | — | — |
| **`05-22`** | **Global Soft-Delete Policy**<br>Universal temporal soft-delete cascading, audit preserving, purge schedules. | `DEL-07` to `DEL-10` | `DB`, `BE` | ⏳ `NOT_STARTED` | — | — |
| **`05-23`** | **Ethiopian Public Holiday Calendar**<br>Ge'ez Computus algorithmic calendar, religious holiday rental rate adjustments. | `CAL-01` to `CAL-06` | `BE`, `MOB` | ⏳ `NOT_STARTED` | — | — |
| **`05-24`** | **Marketplace Integrity & Abuse Controls**<br>Anti-circumvention phone number masking, duplicate booking mitigation, fraud score. | `INT-01` to `INT-08` | `BE`, `OAS`, `ADM` | ⏳ `NOT_STARTED` | — | — |

---

## 4. Phase 2 (Corporate B2B Extension) Partitions (`06-01` to `06-05`)

| Partition ID | Partition Name & Functional Scope | Requirement IDs | Subsystems Target | Status |
|:---:|---|---|:---:|:---:|
| **`06-01`** | **Corporate Account Management**<br>Enterprise organizational structure, corporate hierarchy, sub-contractors. | `CORP-01` to `CORP-10` | `DB`, `OAS`, `BE`, `ADM` | ⏳ `NOT_STARTED` |
| **`06-02`** | **Corporate Approval Workflows**<br>Multi-tier dual-authorization purchase orders, spend approval gates. | `CORP-11` to `CORP-20` | `DB`, `OAS`, `BE`, `ADM` | ⏳ `NOT_STARTED` |
| **`06-03`** | **Corporate Pricing & Payment (NET-30)**<br>Negotiated volume rate cards, credit limits, post-invoiced settlement. | `CORP-21` to `CORP-30` | `DB`, `OAS`, `BE`, `ADM` | ⏳ `NOT_STARTED` |
| **`06-04`** | **Corporate Fleet Management (Owner Side)**<br>Enterprise yard dispatching, multi-depot fleet allocation. | `CORP-31` to `CORP-40` | `DB`, `OAS`, `BE`, `ADM` | ⏳ `NOT_STARTED` |
| **`06-05`** | **Corporate Equipment & Utilization Reporting**<br>Telematics-driven asset efficiency, fuel burn, carbon/ESG reporting. | `CORP-41` to `CORP-50` | `BE`, `ADM` | ⏳ `NOT_STARTED` |

---

## 5. Phase 3 (Advanced Scale & Automation) Partitions (`07-01` to `07-09`)

| Partition ID | Partition Name & Functional Scope | Requirement IDs | Subsystems Target | Status |
|:---:|---|---|:---:|:---:|
| **`07-01`** | **In-App Payment Gateway Integration**<br>Direct Telebirr, M-Pesa, CBE Birr automated gateway settlement. | `ADV-PAY-01` to `ADV-PAY-10` | `DB`, `BE`, `MOB` | ⏳ `NOT_STARTED` |
| **`07-02`** | **Automated Commission Deduction Engine**<br>Direct automated split-payments and instantaneous tax escrow. | `ADV-COM-01` to `ADV-COM-08` | `DB`, `BE` | ⏳ `NOT_STARTED` |
| **`07-03`** | **Full Logistics Marketplace & Route Optimization**<br>Automated lowboy load balancing and multi-stop dispatch routing. | `ADV-LOG-01` to `ADV-LOG-12` | `DB`, `BE`, `MOB` | ⏳ `NOT_STARTED` |
| **`07-04`** | **Certified Equipment Operator Marketplace**<br>Heavy equipment operator bidding, skill verification, ratings. | `ADV-OP-01` to `ADV-OP-10` | `DB`, `BE`, `MOB` | ⏳ `NOT_STARTED` |
| **`07-05`** | **AI-Based Demand Prediction & Dynamic Surge Pricing**<br>Regional construction seasonal pricing models and demand prediction. | `ADV-AI-01` to `ADV-AI-08` | `BE`, `ADM` | ⏳ `NOT_STARTED` |
| **`07-06`** | **IoT Hardware Telematics Integration**<br>Direct CAN-Bus IoT gateway ingestion, automated engine hours, fault codes. | `ADV-IOT-01` to `ADV-IOT-12` | `DB`, `BE` | ⏳ `NOT_STARTED` |
| **`07-07`** | **White-Label Partner API**<br>External developer OpenAPI portal, OAuth2 client credentials, webhook dispatch. | `ADV-API-01` to `ADV-API-10` | `OAS`, `BE` | ⏳ `NOT_STARTED` |
| **`07-08`** | **Multi-Region & Multi-Currency Expansion**<br>Cross-border East Africa expansion (ETB, KES, USD, TZS, UGX). | `ADV-REG-01` to `ADV-REG-10` | `DB`, `BE`, `MOB`, `ADM` | ⏳ `NOT_STARTED` |
| **`07-09`** | **Advanced Fleet Analytics & Predictive Maintenance**<br>Predictive component failure alarms and residual value depreciation. | `ADV-ANL-01` to `ADV-ANL-10` | `BE`, `ADM` | ⏳ `NOT_STARTED` |

---

## 6. Active Execution Log

| Date | Partition | Transition | Summary of Completed Milestones | Author / Reviewer |
|:---:|:---:|:---:|---|:---:|
| 2026-10-26 | `00` | `COMPLETED` | Phase 0 Documentation & Governance Baseline approved. | Engineering Lead |
| 2026-10-30 | `01` | `COMPLETED` | Phase 1 Machine Bootstrapping & toolchain pinning verified. | Engineering Lead |
| 2026-10-31 | `05-01` | `COMPLETED` | Registration & OTP auth pipeline implemented across backend and mobile. | Agentic Pair |
| 2026-10-31 | `05-02` | `COMPLETED` | Fleet machinery onboarding and multi-photo camera reticles completed. | Agentic Pair |
| 2026-11-01 | `05-03` | `COMPLETED` | Attachment inventory with GiST exclusion constraints completed. | Agentic Pair |
| 2026-11-02 | `05-04` | `COMPLETED` | PostGIS spatial search, filtering, and 3-way comparison matrix completed. | Agentic Pair |
| 2026-11-02 | `REFACTOR` | `COMPLETED` | Pure Go standard library (`net/http.ServeMux` + `log/slog`) completed. | Agentic Pair |
| 2026-11-03 | `05-05` | `COMPLETED` | 16-state booking FSM, 60-min hold lock, quoting, and bank slips completed. | Agentic Pair |
| 2026-11-03 | `05-06` | `COMPLETED` | Active rental HUD, 8.0 hr shift telemetry, breakdown pausing completed. | Agentic Pair |
| 2026-11-05 | `05-07` | `COMPLETED` | Return inspection, damage claims, universal reconciliation statement, and double-blind reviews completed. | Agentic Pair |
