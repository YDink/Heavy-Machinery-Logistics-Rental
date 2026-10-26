# Product Requirements Document (PRD) — Dink Machinery

**Product:** Dink Machinery Rental Marketplace  
**Market:** Emerging markets, launching in Ethiopia (`country_code = 'ET'`)  
**Architecture Source:** `docs/source/Heavy_Machinery_Rental_SRS_v2.md`

---

## 1. Executive Summary & Vision

Dink Machinery is a digital aggregator for heavy machinery and logistics equipment. The platform connects contractors/renters with verified equipment owners and certified operators, enabling transparent discovery, verified bank-transfer bookings, documented handovers, and structured dispute resolution.

The system is designed from day one with a **Multi-Country Parameterization Engine** (`LOC-01-BE`, `ADR-015`) to enable expansion into new geographical markets without backend refactoring.

---

## 2. User Personas

| Persona | Role | Primary Goal | Critical Pain Point |
|---|---|---|---|
| **Abraham** | Renter (Contractor) | Locate reliable, nearby machinery with transparent rates. | Hidden fees, unverified owners, manual equipment failure disputes. |
| **Tesfaye** | Owner (Fleet Provider) | Maximize machinery utilization and guarantee payment receipt. | Idle equipment, delayed manual payments, default risk. |
| **Henok** | Operator (Machinist) | Receive verified job assignments and track operational hours. | Unpaid overtime, arbitrary dispute fault attribution. |
| **Sarah** | Admin (Operations) | Verify KYC, review listings, and adjudicate disputes fairly. | Fragmented evidence, lack of audit trails, compliance overhead. |

---

## 3. Product Principles

1. **Non-Custodial Direct Transfer Model (MVP):** Under `MANUAL_TRANSFER`, renters pay owners directly via local bank transfers. The platform never holds or escrows renter funds. Commission is debited from an owner-prefunded `OwnerWallet`.
2. **Verifiable Handover & Return Proof:** Machinery custody transitions (`READY_FOR_PICKUP` → `ACTIVE` and `ACTIVE` → `COMPLETED`) require pre/post-rental condition photos, dual hour-meter readings, and verifiable attestation (signatures/OTP).
3. **Pluggable Integration Architecture:** National identity registries (`identity_providers`) and payment processors (`payment_gateways`) are pluggable, country-scoped adapters.
4. **Offline-First Resilience:** Mobile browsing relies on local caching (`isar_community`) and resilient background upload queues for low-connectivity construction environments.

---

## 4. Requirement Partition Index

Detailed specifications are partitioned into `docs/shared/requirements/`:
- **Registration & Auth:** `05-01-registration-authentication.md`
- **Fleet Registration:** `05-02-fleet-registration-machinery-logistics.md`
- **Discovery & Search:** `05-04-fleet-discovery-search-comparison.md`
- **Booking Lifecycle:** `05-05-booking-reservation-flow.md`
- **Active Rental & Handover:** `05-06-active-rental-management.md`
- **Completion & Damage:** `05-07-rental-completion-return-damage.md`
- **Logistics Transport:** `05-08-logistics-transport-integration.md`
- **Commission & Payments:** `05-09-commission-revenue-owner-payouts.md`
- **Admin Dynamic Engine:** `05-11-admin-dashboard-dynamic-ui-engine.md`