# Heavy Machinery & Logistics Rental Aggregator — Dink Machinery
## Software Requirements Specification — Master Directory & Partition Traceability Matrix (`docs/shared/requirements/INDEX.md`)

**Document Identifier:** `SRS-DINK-CORE-V7.0-INDEX`  
**Classification:** Canonical System Requirements Specification & Traceability Matrix (Single Source of Truth)  
**Governing Standard:** IEEE Std 830-1998 / IEEE/ISO/IEC 29148:2018 (Systems and Software Engineering — Requirements Engineering)  
**Associated Artifacts:** `docs/source/Heavy_Machinery_Rental_SRS_v2.md`, `docs/shared/architecture/INDEX.md`, `docs/shared/flows/INDEX.md`, `docs/decisions/INDEX.md`, `contracts/openapi.yaml`  
**Target Runtimes:** Go 1.27.3 Backend Monolith, Flutter 3.47.9 Mobile Client, Next.js 16.x Admin Console (latest CVE security patch), PostgreSQL 18.6 + PostGIS 3.6.4, Valkey 9.1.3  
**Target Operating Markets:** Ethiopia (`ET`) & Kenya (`KE`) — Dynamic Multi-Country Parameterization (`LOC-01-BE`)  
**Compliance Standards:** OWASP ASVS 4.0.3 Level 3, RFC 9457 (Problem Details), ISO 4217, ISO 3166-1 alpha-2  
**Status:** Approved Engineering Baseline

---

## 1. Executive Summary & Specification Architecture

The Software Requirements Specification (SRS) for Dink Machinery establishes the functional, architectural, regulatory, and quality requirements governing the heavy machinery and logistics rental marketplace. To support rapid agentic pair-programming, deterministic testing, and partition-by-partition feature branch development (`feature/<partition-id>`), the canonical specification is decomposed into modular requirement partition slices.

---

## 2. Requirements Partition Matrix & Traceability Index

| Partition ID | Specification Document | Title & Functional Perimeter | Target Personas | Governing ADRs |
|---|---|---|---|---|
| **`01`** | [`01-introduction.md`](./01-introduction.md) | Introduction, Document Scope & IEEE Governance | All Stakeholders | `ADR-001`, `ADR-012` |
| **`02`** | [`02-glossary.md`](./02-glossary.md) | Universal Taxonomy, Domain Invariants & Conventions | All Stakeholders | `ADR-002`, `ADR-008` |
| **`03`** | [`03-product-vision.md`](./03-product-vision.md) | Strategic Vision, Market Problem & Solution Architecture | Executive, Product | `ADR-001`, `ADR-010` |
| **`04`** | [`04-user-roles.md`](./04-user-roles.md) | User Roles, Granular Permissions & RBAC Matrix | All Personas | `ADR-008`, `ADR-014` |
| **`05-01`** | [`05-01-registration-authentication.md`](./05-01-registration-authentication.md) | User Onboarding, Identity Attestation & Authentication | Renters, Owners, Ops | `ADR-011`, `ADR-014` |
| **`05-02`** | [`05-02-fleet-registration-machinery-logistics.md`](./05-02-fleet-registration-machinery-logistics.md) | Machinery & Carrier Fleet Registration | Fleet Owners | `ADR-003`, `ADR-005` |
| **`05-03`** | [`05-03-attachment-registration-inventory.md`](./05-03-attachment-registration-inventory.md) | Standalone & Bundled Attachment Inventory | Fleet Owners | `ADR-003`, `ADR-005` |
| **`05-04`** | [`05-04-fleet-discovery-search-comparison.md`](./05-04-fleet-discovery-search-comparison.md) | Spatial Radius Discovery, Filtering & Comparison | Renters | `ADR-004`, `ADR-005` |
| **`05-05`** | [`05-05-booking-reservation-flow.md`](./05-05-booking-reservation-flow.md) | 20-State Booking Machine & Hold Concurrency Engine | Renters, Owners | `ADR-001`, `ADR-002`, `ADR-004`, `ADR-010` |
| **`05-06`** | [`05-06-active-rental-management.md`](./05-06-active-rental-management.md) | Active Rental Monitoring & Incident Logging | Renters, Operators | `ADR-006`, `ADR-012` |
| **`05-07`** | [`05-07-rental-completion-return-damage.md`](./05-07-rental-completion-return-damage.md) | Return Triage, Meter Check & Damage Claims | Renters, Owners, Ops | `ADR-004`, `ADR-010` |
| **`05-08`** | [`05-08-logistics-transport-integration.md`](./05-08-logistics-transport-integration.md) | Lowboy Carrier Booking & Dispatch Coordination | Transporters, Renters | `ADR-001`, `ADR-005` |
| **`05-09`** | [`05-09-commission-revenue-owner-payouts.md`](./05-09-commission-revenue-owner-payouts.md) | Double-Entry Commission Ledger & Owner Wallets | Finance, Owners | `ADR-002`, `ADR-010` |
| **`05-10`** | [`05-10-notifications-in-app-center.md`](./05-10-notifications-in-app-center.md) | Multi-Channel Notifications (Push, SMS, In-App) | All Personas | `ADR-007`, `ADR-009` |
| **`05-11`** | [`05-11-admin-dashboard-dynamic-ui-engine.md`](./05-11-admin-dashboard-dynamic-ui-engine.md) | Server-Driven Dynamic Admin UI Engine | Administrators | `ADR-002`, `ADR-008` |
| **`05-12`** | [`05-12-operator-gps-tracking.md`](./05-12-operator-gps-tracking.md) | Operator Telematics & Periodic GPS Pings | Operators, Renters | `ADR-005`, `ADR-006` |
| **`05-13`** | [`05-13-insurance-document-management.md`](./05-13-insurance-document-management.md) | Machinery Insurance Verification & Expiry Alarms | Compliance, Owners | `ADR-005`, `ADR-012` |
| **`05-14`** | [`05-14-user-onboarding-help.md`](./05-14-user-onboarding-help.md) | Interactive In-App Tutorials & Multilingual FAQs | New Users | `ADR-006`, `ADR-013` |
| **`05-16`** | [`05-16-accessibility-internationalisation.md`](./05-16-accessibility-internationalisation.md) | WCAG 2.2 AAA Accessibility & Multilingual OTA | Field Users | `ADR-006`, `ADR-013` |
| **`05-17`** | [`05-17-offline-low-connectivity-support.md`](./05-17-offline-low-connectivity-support.md) | Isar NoSQL Offline Outbox & Retry Sync Queue | Field Operators | `ADR-002`, `ADR-006` |
| **`05-18`** | [`05-18-error-handling-user-feedback.md`](./05-18-error-handling-user-feedback.md) | RFC 9457 Problem Details & Client Resiliency | All Users | `ADR-002`, `ADR-008` |
| **`05-19`** | [`05-19-account-deletion.md`](./05-19-account-deletion.md) | Right-to-Erasure & 10-Year Statutory Ledger Hold | Compliance, Users | `ADR-005`, `ADR-010` |
| **`05-20`** | [`05-20-policy-terms-management.md`](./05-20-policy-terms-management.md) | Versioned Terms of Service & Click-Wrap Binding | Legal, All Users | `ADR-005`, `ADR-015` |
| **`05-21`** | [`05-21-feature-flagging-remote-configuration.md`](./05-21-feature-flagging-remote-configuration.md) | Runtime Feature Toggles & Canary Rollouts | Operations, SRE | `ADR-004`, `ADR-008` |
| **`05-22`** | [`05-22-global-soft-delete-policy.md`](./05-22-global-soft-delete-policy.md) | Temporal Soft Deletes (`deleted_at`) & Auditing | Developers, SRE | `ADR-003`, `ADR-008` |
| **`05-23`** | [`05-23-ethiopian-public-holiday-calendar.md`](./05-23-ethiopian-public-holiday-calendar.md) | Ge'ez Computus & Public Holiday Booking Engine | Renters, Owners | `ADR-013`, `ADR-015` |
| **`05-24`** | [`05-24-marketplace-integrity-abuse-controls.md`](./05-24-marketplace-integrity-abuse-controls.md) | Anti-Circumvention, Phone Redaction & Abuse Defense | Compliance, Ops | `ADR-009`, `ADR-010` |
| **`06-01`** | [`06-01-corporate-account-management.md`](./06-01-corporate-account-management.md) | Enterprise Organizational Hierarchy & Workspaces | Corporate Renters | `ADR-001`, `ADR-014` |
| **`06-02`** | [`06-02-corporate-approval-workflows.md`](./06-02-corporate-approval-workflows.md) | Multi-Tiered Dual-Authorization Workflow Engine | Procurement Leads | `ADR-004`, `ADR-007` |
| **`06-03`** | [`06-03-corporate-pricing-payment.md`](./06-03-corporate-pricing-payment.md) | Negotiated Volume Rate Cards & `NET-30` Credit | Corporate Finance | `ADR-002`, `ADR-010` |
| **`06-04`** | [`06-04-corporate-fleet-management-owner-side.md`](./06-04-corporate-fleet-management-owner-side.md) | Enterprise Fleet Depots & Yard Dispatching | Corporate Owners | `ADR-003`, `ADR-006` |
| **`06-05`** | [`06-05-corporate-reporting.md`](./06-05-corporate-reporting.md) | Equipment Utilization, ESG & Financial Reports | Corporate Directors | `ADR-004`, `ADR-008` |
| **`07-01`** | [`07-01-in-app-payment-processing.md`](./07-01-in-app-payment-processing.md) | Phase 3 In-App Telebirr/M-Pesa Gateway Escrow | Renters, Owners | `ADR-002`, `ADR-010` |
| **`07-02`** | [`07-02-automatic-commission-deduction.md`](./07-02-automatic-commission-deduction.md) | Phase 3 Split Settlement & Automated Deductions | Finance | `ADR-004`, `ADR-010` |
| **`07-03`** | [`07-03-full-logistics-marketplace.md`](./07-03-full-logistics-marketplace.md) | Phase 3 Automated Freight Route Optimization | Transporters | `ADR-001`, `ADR-006` |
| **`07-04`** | [`07-04-operator-marketplace.md`](./07-04-operator-marketplace.md) | Phase 3 Independent Heavy Machinery Operators | Field Operators | `ADR-006`, `ADR-011` |
| **`07-05`** | [`07-05-ai-based-demand-prediction-surge-pricing.md`](./07-05-ai-based-demand-prediction-surge-pricing.md) | Phase 3 AI Demand Forecasting & Dynamic Pricing | Operations | `ADR-002`, `ADR-005` |
| **`07-06`** | [`07-06-iot-hardware-telematics-integration.md`](./07-06-iot-hardware-telematics-integration.md) | Phase 3 CAN-Bus Telematics & Geofence Alarms | Fleet Owners | `ADR-007`, `ADR-009` |
| **`07-07`** | [`07-07-white-label-api.md`](./07-07-white-label-api.md) | Phase 3 B2B Partner White-Label Integration API | Enterprise Partners | `ADR-002`, `ADR-009` |
| **`07-08`** | [`07-08-multi-region-multi-currency.md`](./07-08-multi-region-multi-currency.md) | Phase 3 Regional East Africa Expansion (UG, TZ, RW) | Operations | `ADR-012`, `ADR-015` |
| **`07-09`** | [`07-09-advanced-analytics-reporting.md`](./07-09-advanced-analytics-reporting.md) | Phase 3 Machine Learning Asset Depreciation Models | Executive | `ADR-004`, `ADR-008` |
| **`08`** | [`08-booking-state-machine.md`](./08-booking-state-machine.md) | Canonical 20-State Rental State Machine | Core Backend | `ADR-001`, `ADR-004` |
| **`09`** | [`09-non-functional-requirements.md`](./09-non-functional-requirements.md) | Non-Functional Requirements & ASVS L3 Baselines | SRE, Security | `ADR-009`, `ADR-012` |
| **`10`** | [`10-data-model.md`](./10-data-model.md) | 111-Table Relational Schema & DDL Architecture | Data Architects | `ADR-003`, `ADR-005` |
| **`11`** | [`11-dynamic-config.md`](./11-dynamic-config.md) | Dynamic Multi-Country Configuration Engine | SRE, Ops | `ADR-008`, `ADR-015` |
| **`12`** | [`12-tech-stack.md`](./12-tech-stack.md) | Pinned Technology Stack & Deployment View | Platform Engineers | `ADR-001`, `ADR-012` |
| **`13`** | [`13-roadmap.md`](./13-roadmap.md) | 12-Month Implementation Roadmap & Sprint Plan | Product, Engineering | `ADR-001`, `ADR-012` |
| **`14`** | [`14-appendices.md`](./14-appendices.md) | Canonical Reference Appendices & Seed Catalogs | All Engineers | `ADR-002`, `ADR-008` |
