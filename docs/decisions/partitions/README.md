# Level 2 (L2) Partition Decision Records — Dink Machinery

**Classification:** Bounded Context Architectural Governance & Local Technical Decisions  
**Governing Standard:** `AGENTS.md § 5`, `docs/AI_TASKS/04-new-adr.md`, `docs/quality/BOOTSTRAP_CHECKLIST.md`  
**Scope:** Localized engineering decisions within specific functional partitions. Cross-cutting decisions are escalated to Level 1 System ADRs (`docs/decisions/ADR-NNN-*.md`).

---

## 1. Governance & 5-Point Qualification Routing

Level 2 Partition Decisions document non-trivial technical trade-offs localized entirely within a single requirement partition or subsystem. They must satisfy at least two criteria of the **5-Point Qualification Test** while answering **NO** to the Scope Test (confined to a single partition).

Every L2 decision conforms strictly to the standardized 7-section layout validated by `scripts/validate-decision-template.sh` and `scripts/validate-decision-links.sh`.

---

## 2. Partition Decisions Directory Index (40 Partitions)

| Partition ID | File Path | Focus Area & Bounded Context | Target Decisions |
|---|---|---|---|
| **`05-01`** | [`05-01-registration-authentication.md`](./05-01-registration-authentication.md) | User Registration & Auth | `D-05-01-001`, `D-05-01-002` |
| **`05-02`** | [`05-02-fleet-registration-machinery-logistics.md`](./05-02-fleet-registration-machinery-logistics.md) | Fleet Registry & Spec Engine | `D-05-02-001`, `D-05-02-002` |
| **`05-03`** | [`05-03-attachment-registration-inventory.md`](./05-03-attachment-registration-inventory.md) | Attachment Compatibility | `D-05-03-001`, `D-05-03-002` |
| **`05-04`** | [`05-04-fleet-discovery-search-comparison.md`](./05-04-fleet-discovery-search-comparison.md) | Spatial Radius Discovery | `D-05-04-001`, `D-05-04-002` |
| **`05-05`** | [`05-05-booking-reservation-flow.md`](./05-05-booking-reservation-flow.md) | Booking Holds & Atomicity | `D-05-05-001`, `D-05-05-002` |
| **`05-06`** | [`05-06-active-rental-management.md`](./05-06-active-rental-management.md) | Active Rental Monitoring | `D-05-06-001`, `D-05-06-002` |
| **`05-07`** | [`05-07-rental-completion-return-damage.md`](./05-07-rental-completion-return-damage.md) | Return Triage & Disputes | `D-05-07-001`, `D-05-07-002` |
| **`05-08`** | [`05-08-logistics-transport-integration.md`](./05-08-logistics-transport-integration.md) | Lowboy Carrier Booking | `D-05-08-001`, `D-05-08-002` |
| **`05-09`** | [`05-09-commission-revenue-owner-payouts.md`](./05-09-commission-revenue-owner-payouts.md) | Double-Entry Commission | `D-05-09-001`, `D-05-09-002` |
| **`05-10`** | [`05-10-notifications-in-app-center.md`](./05-10-notifications-in-app-center.md) | Multi-Channel Outbox | `D-05-10-001`, `D-05-10-002` |
| **`05-11`** | [`05-11-admin-dashboard-dynamic-ui-engine.md`](./05-11-admin-dashboard-dynamic-ui-engine.md) | Dynamic UI Schema Engine | `D-05-11-001`, `D-05-11-002` |
| **`05-12`** | [`05-12-operator-gps-tracking.md`](./05-12-operator-gps-tracking.md) | Operator Telematics & GPS | `D-05-12-001`, `D-05-12-002` |
| **`05-13`** | [`05-13-insurance-document-management.md`](./05-13-insurance-document-management.md) | Insurance Expiry Alarms | `D-05-13-001`, `D-05-13-002` |
| **`05-14`** | [`05-14-user-onboarding-help.md`](./05-14-user-onboarding-help.md) | In-App Step Steppers | `D-05-14-001`, `D-05-14-002` |
| **`05-16`** | [`05-16-accessibility-internationalisation.md`](./05-16-accessibility-internationalisation.md) | WCAG AAA Outdoor Mode | `D-05-16-001`, `D-05-16-002` |
| **`05-17`** | [`05-17-offline-low-connectivity-support.md`](./05-17-offline-low-connectivity-support.md) | Isar NoSQL Outbox Queue | `D-05-17-001`, `D-05-17-002` |
| **`05-18`** | [`05-18-error-handling-user-feedback.md`](./05-18-error-handling-user-feedback.md) | RFC 9457 Problem Details | `D-05-18-001`, `D-05-18-002` |
| **`05-19`** | [`05-19-account-deletion.md`](./05-19-account-deletion.md) | Right-to-Erasure & 10Y Hold | `D-05-19-001`, `D-05-19-002` |
| **`05-20`** | [`05-20-policy-terms-management.md`](./05-20-policy-terms-management.md) | Click-Wrap Version Gate | `D-05-20-001`, `D-05-20-002` |
| **`05-21`** | [`05-21-feature-flagging-remote-configuration.md`](./05-21-feature-flagging-remote-configuration.md) | Canary Feature Toggles | `D-05-21-001`, `D-05-21-002` |
| **`05-22`** | [`05-22-global-soft-delete-policy.md`](./05-22-global-soft-delete-policy.md) | deleted_at Standard | `D-05-22-001`, `D-05-22-002` |
| **`05-23`** | [`05-23-ethiopian-public-holiday-calendar.md`](./05-23-ethiopian-public-holiday-calendar.md) | Ge'ez Calendar Computus | `D-05-23-001`, `D-05-23-002` |
| **`05-24`** | [`05-24-marketplace-integrity-abuse-controls.md`](./05-24-marketplace-integrity-abuse-controls.md) | Anti-Circumvention Redaction | `D-05-24-001`, `D-05-24-002` |
| **`06-01`** | [`06-01-corporate-account-management.md`](./06-01-corporate-account-management.md) | Enterprise Workspaces | `D-06-01-001`, `D-06-01-002` |
| **`06-02`** | [`06-02-corporate-approval-workflows.md`](./06-02-corporate-approval-workflows.md) | Multi-Tier Approvals | `D-06-02-001`, `D-06-02-002` |
| **`06-03`** | [`06-03-corporate-pricing-payment.md`](./06-03-corporate-pricing-payment.md) | NET-30 Underwriting | `D-06-03-001`, `D-06-03-002` |
| **`06-04`** | [`06-04-corporate-fleet-management-owner-side.md`](./06-04-corporate-fleet-management-owner-side.md) | Depot Yard Dispatch | `D-06-04-001`, `D-06-04-002` |
| **`06-05`** | [`06-05-corporate-reporting.md`](./06-05-corporate-reporting.md) | ESG Fuel & Utilization | `D-06-05-001`, `D-06-05-002` |
| **`07-01`** | [`07-01-in-app-payment-processing.md`](./07-01-in-app-payment-processing.md) | Mobile Money Gateway | `D-07-01-001`, `D-07-01-002` |
| **`07-02`** | [`07-02-automatic-commission-deduction.md`](./07-02-automatic-commission-deduction.md) | Gateway Fee Splits | `D-07-02-001`, `D-07-02-002` |
| **`07-03`** | [`07-03-full-logistics-marketplace.md`](./07-03-full-logistics-marketplace.md) | Automated Freight Matching | `D-07-03-001`, `D-07-03-002` |
| **`07-04`** | [`07-04-operator-marketplace.md`](./07-04-operator-marketplace.md) | Independent Operator Pool | `D-07-04-001`, `D-07-04-002` |
| **`07-05`** | [`07-05-ai-based-demand-prediction-surge-pricing.md`](./07-05-ai-based-demand-prediction-surge-pricing.md) | AI Surge Pricing Caps | `D-07-05-001`, `D-07-05-002` |
| **`07-06`** | [`07-06-iot-hardware-telematics-integration.md`](./07-06-iot-hardware-telematics-integration.md) | CAN-Bus MQTT Ingestion | `D-07-06-001`, `D-07-06-002` |
| **`07-07`** | [`07-07-white-label-api.md`](./07-07-white-label-api.md) | Partner mTLS Gateway | `D-07-07-001`, `D-07-07-002` |
| **`07-08`** | [`07-08-multi-region-multi-currency.md`](./07-08-multi-region-multi-currency.md) | Multi-Currency & Carnet | `D-07-08-001`, `D-07-08-002` |
| **`07-09`** | [`07-09-advanced-analytics-reporting.md`](./07-09-advanced-analytics-reporting.md) | Predictive Maintenance | `D-07-09-001`, `D-07-09-002` |
| **`08`** | [`08-booking-state-machine.md`](./08-booking-state-machine.md) | 20-State Machine Graph | `D-08-01-001`, `D-08-01-002` |
| **`10`** | [`10-data-model.md`](./10-data-model.md) | 111-Table Relational DDL | `D-10-01-001`, `D-10-01-002` |
| **`11`** | [`11-dynamic-config.md`](./11-dynamic-config.md) | Cascading Config Engine | `D-11-01-001`, `D-11-01-002` |
