# Heavy Machinery & Logistics Rental Aggregator — Dink Machinery
## Software Requirements Specification — 13. MVP Scope & Implementation Roadmap (`docs/shared/requirements/13-roadmap.md`)

**Document Identifier:** `SRS-DINK-CORE-V7.0-SEC-13`  
**Classification:** Canonical System Requirements Specification & Engineering Roadmap Baseline (Single Source of Truth)  
**Governing Standard:** IEEE Std 830-1998 / IEEE/ISO/IEC 29148:2018 (Systems and Software Engineering — Requirements Engineering)  
**Associated Artifacts:** `SRS-DINK-CORE-V7.0` (§13), `ARCH-INTRO-DINK-CORE-V7.0`, `ARCH-PRIN-DINK-CORE-V7.0`, `ARCH-ROADMAP-DINK-CORE-V7.0`, `contracts/openapi.yaml`, `ADR-001`, `ADR-002`, `ADR-006`, `ADR-007`, `ADR-008`, `ADR-010`, `ADR-012`, `ADR-013`  
**Target Runtimes:** Go 1.27.3 Backend Monolith, Flutter 3.47.8 Mobile Client, Next.js 16.x Admin Console (latest CVE security patch), PostgreSQL 18.6 + PostGIS 3.6.4, Valkey 9.1.3, SeaweedFS S3, ClamAV  
**Target Operating Markets:** Ethiopia (`ET`) & Kenya (`KE`) — Dynamic Multi-Country Parameterization (`LOC-01-BE`)  
**Compliance Standards:** OWASP ASVS 4.0.3 Level 3, RFC 9457 (Problem Details), National Bank of Ethiopia (NBE) Directives  
**Status:** Approved Engineering Baseline

---

## 1. Executive Summary & Phased Delivery Model

The Dink Machinery implementation executes across three structured delivery phases:
- **Phase 1 (MVP Foundation — Sprints 1–30 / 12 Months):** Core bilateral marketplace enabling heavy machinery discovery, booking reservations, direct bank transfer settlement, field inspection verification, return triage, dynamic administrative operations, and multi-country parameterization across Ethiopia and Kenya.
- **Phase 2 (Corporate & Multi-Tenant Accounts — Months 13–18):** Enterprise organizational hierarchy, multi-tiered approval workflows, master rental service agreements (MSAs), consolidated monthly invoicing (`NET-30`), and fleet management tools for corporate contractors.
- **Phase 3 (Advanced Telematics & Autonomous Marketplace — Months 19–24):** In-app gateway payment splits (`PLATFORM_HELD`), direct CAN-bus IoT telematics ingestion, automated AI surge pricing, independent operator marketplace, and white-label enterprise APIs.

---

## 2. MVP Scope Boundaries

### 2.1 In Scope (Phase 1 MVP Foundation)
1. **User Onboarding & Authentication:** Individual renter, owner, and operator registration with MSISDN phone OTP, ES256 JWT lifecycle, and pluggable identity attestation (Fayda ET reference adapter and manual verification docket).
2. **Fleet Registry & Dynamic Taxonomy:** Comprehensive fleet listing with 5–20 multi-angle photos, EAV category specification templates, and bundled attachments.
3. **Spatial Discovery & Search:** PostGIS spheroidal radius search (`ST_DWithin`), category filter facets, sort options, side-by-side machinery comparison, and favorites/wishlist.
4. **20-State Booking Machine & Concurrency Engine:** Canonical 20-state reservation lifecycle with 60-minute Valkey mutual exclusion holds and versioned Optimistic Concurrency Control (OCC).
5. **Non-Custodial Settlement (`MANUAL_TRANSFER`):** Direct bank transfer settlement rail with bank slip receipt upload, perceptual hash duplicate detection, owner confirmation, and automated `OwnerWallet` commission debits (`BK-17-BE`).
6. **Field Inspection & Verification Stepper:** Handover inspection workflow with camera-enforced engine hour meter capture, fuel gauge level, ClamAV antivirus stream scanning, and mutual digital signatures.
7. **Return Triage, Damage & Rating:** Return inspection workflow, monotonic hour meter validation, excess hours billing, damage claim filing with split-screen photo diffs, and mutual star ratings.
8. **Server-Driven Admin Console:** Next.js 16 App Router back-office driven by the dynamic metadata schema endpoint (`/api/v1/admin/ui-config`), featuring granular RBAC, dispute arbitration workbench, and immutable audit logs.
9. **Vernacular Localization & Typography:** Server-driven over-the-air (OTA) localization supporting English and Amharic with dynamic Ethiopic font delivery (`ADR-013`).
10. **Offline Low-Connectivity Perimeters:** Flutter local Isar NoSQL draft storage and background outbox queue for network-resilient field operations.

### 2.2 Out of Scope (Deferred Post-MVP)
- Corporate multi-tiered approval hierarchies and purchase order matching (Phase 2).
- Consolidated monthly invoicing and corporate credit lines (Phase 2).
- Automated direct payment gateway commission splits (Phase 3).
- Automated CAN-bus telematics hardware streaming (Phase 3).
- Standalone independent operator marketplace (Phase 3).
- AI demand forecasting and real-time algorithmic surge pricing (Phase 3).

---

## 3. Implementation Roadmap (30 Sprints / 12 Months)

| Sprint Cluster | Focus Area | Deliverables & Milestones | Duration | Target Partitions |
|---|---|---|---|---|
| **Sprints 1–2** | Platform Scaffold & Environments | Base Go modular monolith scaffold, CI/CD pipelines, Docker Compose, Traefik gateway, PostgreSQL 18.6 + PostGIS setup, baseline database migrations. | 1.0 Month | `00-governance`, `12-tech-stack` |
| **Sprints 3–5** | Identity, Auth & KYC Engine | Phone OTP challenges, ES256 JWT lifecycle, Valkey sessions, Fayda identity adapter, owner KYC queues, Argon2id password derivation. | 1.5 Months | `05-01`, `05-19` |
| **Sprints 6–9** | Fleet Registry & Dynamic Specs | Category taxonomies, dynamic EAV spec engine, photo upload pipeline with ClamAV stream scanning, availability calendar. | 2.0 Months | `05-02`, `05-03` |
| **Sprints 10–13**| Discovery & Spatial Search | PostGIS radius search, dynamic category filter facets, sorting queries, compare tables, wishlist, keyset pagination. | 2.0 Months | `05-04` |
| **Sprints 14–16**| Booking Engine & Payment Proof | 20-state machine, 60-min hold mutex locks, bank slip upload, owner verification, cancellation logic, OCC versioning. | 1.5 Months | `05-05`, `08-state-machine` |
| **Sprints 17–18**| Dispute Center & SLA Escalation | Tri-party dispute chat, photographic damage diff viewer, 3-day SLA cron escalation, binding administrative rulings. | 1.0 Month | `05-07`, `05-18` |
| **Sprints 19–20**| Returns, Damage & Rating | Return inspection wizard, monotonic meter check, excess hours billing, review moderation. | 1.0 Month | `05-06`, `05-07` |
| **Sprints 21–22**| Ledger, Wallets & Notifications | Immutable double-entry ledger, owner wallet top-ups and debits, FCM push alerts, SMS gateway fallback. | 1.0 Month | `05-09`, `05-10` |
| **Sprints 23–24**| Admin Dynamic UI & Audit Logs | Server-driven UI config engine (`/ui-config`), RBAC permission matrix, immutable audit log diff viewer. | 1.0 Month | `05-11`, `11-dynamic-config` |
| **Sprints 25–26**| Onboarding, Offline & i18n | Multilingual OTA bundle sync, Ethiopic font delivery, Isar offline draft queue, onboarding tutorials. | 1.0 Month | `05-14`, `05-16`, `05-17`, `05-23` |
| **Sprints 27–28**| Security Hardening & ASVS L3 | Zero-PII logging audit, Semgrep SAST, penetration testing drills, secret scanning, Crashlytics. | 1.0 Month | `05-20`, `05-22`, `05-24` |
| **Sprints 29–30**| Verification, Load Testing & Launch | k6 geospatial load testing (P95 $< 80\text{ms}$), E2E user journeys, staging acceptance drills, production launch. | 1.0 Month | `All MVP Partitions` |
