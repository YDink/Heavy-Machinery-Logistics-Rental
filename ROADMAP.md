# Dink Machinery — Implementation Roadmap & Status Tracker[cite: 38]

**Single source of truth for cross-machine implementation progress.**[cite: 37, 38]  
Status updates are committed atomically with code and tests on the active partition branch[cite: 37, 38].

## Status Legend[cite: 38]
| Status | Definition |
|---|---|
| `todo` | Not started[cite: 38] |
| `wip` | In active implementation (stop point noted in `CONTEXT.md`)[cite: 38] |
| `done` | Implemented, tests passing, committed and pushed on partition branch[cite: 38] |
| `✔ merged` | Partition merged to `main` following contract and quality verification[cite: 38] |
| `deferred` | Scheduled for post-MVP phases (Phase 2 or Phase 3)[cite: 38] |

---

## Phase 0 — Documentation & Governance Gate (Current Phase)

| Task | Target Artifact | Status | Notes |
|---|---|---|---|
| AI Context Engine | `AGENTS.md`, wrappers, `.gemini/settings.json` | done | All 8 disciplines codified |
| Foundation Docs | `PRD.md`, `ARCHITECTURE.md`, `DESIGN_SYSTEM.md`, `DATABASE.md`, `SECURITY.md` | done | Complete specifications |
| Engineering Standards | `CODE_STYLE.md`, `TESTING.md`, `GIT_STANDARDS.md`, `CONTRIBUTING.md` | done | G0–G4 gate integration |
| QA Architecture | `docs/quality/*` (9 documents) | done | ASVS L3, ISO 25010, axe-core |
| Decision System | `docs/decisions/` (ADR-001..015, templates, scripts)[cite: 38] | done | L1/L2 three-tier model |
| Bootstrap Gate | `docs/quality/BOOTSTRAP_CHECKLIST.md` | wip | Final verification prior to Phase 1 |

---

## Phase 1 — Component Bootstrap (Gated by Phase 0)

| Component | Target Scaffold | Status | Toolchain Rule |
|---|---|---|---|
| Backend | `backend/`[cite: 37] | todo | Go modular monolith, `sqlc`, `pgx/v5`[cite: 32, 38] |
| Mobile | `mobile/`[cite: 37] | todo | Flutter 3.47.8, BLoC, Dio, `isar_community`[cite: 32, 38] |
| Admin Web | `admin/`[cite: 37] | todo | Next.js 16 (App Router), Tailwind, Server Actions[cite: 32, 38] |
| Lockfile Registry | `docs/TECH_STACK.md` | todo | Record installed framework versions |

---

## Phase 2 — MVP Feature Partitions (Branch-per-Partition)[cite: 38]

| Partition ID | Partition Description | Branch[cite: 38] | Status[cite: 38] | Notes[cite: 38] |
|---|---|---|---|---|
| `05-01`[cite: 38] | Registration & Authentication[cite: 38] | `feature/05-01-registration-authentication`[cite: 38] | todo[cite: 38] | REG-R, REG-O, REG-OP, AUTH[cite: 38] |
| `05-22`[cite: 38] | Global Soft-Delete Policy[cite: 38] | `feature/05-22-global-soft-delete-policy`[cite: 38] | todo[cite: 38] | Base SQL views & CI linter[cite: 38] |
| `05-23`[cite: 38] | Country Config & Calendar[cite: 38] | `feature/05-23-ethiopian-public-holiday-calendar`[cite: 38] | todo[cite: 38] | `country_configs` & holiday SLAs[cite: 38] |
| `05-02`[cite: 38] | Fleet Registration[cite: 38] | `feature/05-02-fleet-registration-machinery-logistics`[cite: 38] | todo[cite: 38] | EAV specs, dynamic validation[cite: 38] |
| `05-03`[cite: 38] | Attachment Management[cite: 38] | `feature/05-03-attachment-registration-inventory`[cite: 38] | todo[cite: 38] | Inventory & compatibility[cite: 38] |
| `05-04`[cite: 38] | Discovery & Search[cite: 38] | `feature/05-04-fleet-discovery-search-comparison`[cite: 38] | todo[cite: 38] | PostGIS radius, ranking[cite: 38] |
| `05-05`[cite: 38] | Booking & Reservation[cite: 38] | `feature/05-05-booking-reservation-flow`[cite: 38] | todo[cite: 38] | 20-state machine, holds, disputes[cite: 38] |
| `05-09`[cite: 38] | Commission & Settlement[cite: 38] | `feature/05-09-commission-revenue-owner-payouts`[cite: 38] | todo[cite: 38] | Direct transfer, OwnerWallet[cite: 38] |
| `05-06`[cite: 38] | Active Rental Management[cite: 38] | `feature/05-06-active-rental-management`[cite: 38] | todo[cite: 38] | Handover attestation, incidents[cite: 38] |
| `05-07`[cite: 38] | Completion & Return[cite: 38] | `feature/05-07-rental-completion-return-damage`[cite: 38] | todo[cite: 38] | Hour-meter check, damage math[cite: 38] |
| `05-08`[cite: 38] | Logistics Transport[cite: 38] | `feature/05-08-logistics-transport-integration`[cite: 38] | todo[cite: 38] | Matching, sub-bookings, cascade[cite: 38] |
| `05-10`[cite: 38] | Notification Center[cite: 38] | `feature/05-10-notifications-in-app-center`[cite: 38] | todo[cite: 38] | Prefix router, SMS fallback[cite: 38] |
| `05-11`[cite: 38] | Admin Dynamic Engine[cite: 38] | `feature/05-11-admin-dashboard-dynamic-ui-engine`[cite: 38] | todo[cite: 38] | DB-driven tables & forms[cite: 38] |
| `05-13`[cite: 38] | Insurance Management[cite: 38] | `feature/05-13-insurance-document-management`[cite: 38] | todo[cite: 38] | Document tracking, expiry[cite: 38] |
| `05-16`[cite: 38] | Accessibility & i18n[cite: 38] | `feature/05-16-accessibility-internationalisation`[cite: 38] | todo[cite: 38] | OTA bundles, WCAG AA[cite: 38] |
| `05-17`[cite: 38] | Offline Connectivity[cite: 38] | `feature/05-17-offline-low-connectivity-support`[cite: 38] | todo[cite: 38] | Isar queue, sync banner[cite: 38] |
| `05-18`[cite: 38] | Error Handling & Feedback[cite: 38] | `feature/05-18-error-handling-user-feedback`[cite: 38] | todo[cite: 38] | RFC 9457 error taxonomy[cite: 38] |
| `05-19`[cite: 38] | Account Deletion & Privacy[cite: 38] | `feature/05-19-account-deletion`[cite: 38] | todo[cite: 38] | Grace window, anonymization[cite: 38] |
| `05-20`[cite: 38] | Policy & Terms[cite: 38] | `feature/05-20-policy-terms-management`[cite: 38] | todo[cite: 38] | Blocking version acceptance[cite: 38] |
| `05-21`[cite: 38] | Feature Flagging[cite: 38] | `feature/05-21-feature-flagging-remote-configuration`[cite: 38] | todo[cite: 38] | Valkey dynamic rules[cite: 38] |
| `05-24`[cite: 38] | Fraud & Abuse Controls[cite: 38] | `feature/05-24-marketplace-integrity-abuse-controls`[cite: 38] | todo[cite: 38] | Geovelocity, duplicate checks[cite: 38] |