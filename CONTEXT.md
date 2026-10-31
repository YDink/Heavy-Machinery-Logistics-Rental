# CONTEXT.md — Dink Machinery Active Session State

**Project Phase:** Phase 2 (Partition Feature Development)  
**Active Requirement Partition:** `05-02-fleet-registration-machinery-logistics` (Fleet Registration, Machinery & Logistics Taxonomy Architecture)  
**Active Branch:** `feature/05-02-fleet-registration-machinery-logistics`  
**Current Stop Point:** Day-One Full-Stack Localization & Dynamic RBAC Certified Across Go, Flutter, and Next.js

---

## 1. Current Progress Snapshot

### Completed (Foundational Dynamic RBAC & Full-Stack Localization)
- **Architectural Decision Records (`docs/decisions/`):**
  - Authored and verified `ADR-018: Dynamic Role-Based Access Control (RBAC) and Granular Permission Architecture`.
  - Authored and verified `ADR-019: Full-Stack Vernacular Localization (i18n/l10n) and Ge'ez Script Architecture`.
  - Audited with `scripts/validate-decision-template.sh` and `scripts/validate-decision-links.sh` (19 L1 ADRs, 80 L2 Decisions, 0 errors).
- **Database Schema & Migrations (`migrations/`):**
  - Created `migrations/000003_create_rbac_and_localization_tables.up.sql` and `.down.sql`:
    - `users.preferred_language`: Multi-language preference column (`en`, `am`, `om`, `sw`).
    - Dynamic RBAC: `roles`, `permissions`, `role_permissions`, `user_roles`.
    - Dynamic Localization: `localization_namespaces`, `localization_keys`, `localization_translations`.
    - Seeded canonical 11 roles and 20+ permissions from SRS Appendix B.
    - Seeded baseline English and Amharic dictionaries for auth and errors.
  - Generated sqlc persistence models and queries in `backend/internal/platform/database/db/`.
- **Go Backend Subsystems (`backend/`):**
  - `internal/platform/i18n/`: `LanguageMiddleware` parsing `Accept-Language`, fallback dictionaries for `en` and `am`, context injection.
  - `internal/platform/auth/rbac/`: `RequirePermission` default-deny HTTP middleware returning localized RFC 9457 `403 Forbidden` (`ERR_FORBIDDEN_INSUFFICIENT_PERMISSIONS`).
  - `internal/platform/response/`: `LocalizedError` rendering localized Problem Details.
  - 100% test pass rate across application, domain, rbac, and i18n packages.
- **Flutter Mobile Client (`mobile/`):**
  - Integrated `flutter_localizations` in `pubspec.yaml` with `l10n.yaml`.
  - Created `mobile/lib/l10n/app_en.arb` and `mobile/lib/l10n/app_am.arb`.
  - Compiled typed `AppLocalizations` (`en`, `am`).
  - Implemented `LocaleController` for instant `< 16ms` dynamic language switching.
  - Configured `TypographyConfig` with $1.45\times$ Ge'ez line-height multiplier (`ADR-013`, `SRS §05.16`).
  - Refactored `PhoneInputPage` and `OtpVerificationPage` with zero hardcoded strings.
  - 100% mobile test pass rate (15/15 unit and widget tests); 0 warnings on `fvm flutter analyze`.
- **Next.js Admin Console (`admin/`):**
  - Created `src/lib/i18n/dictionaries.ts` supporting `en` and `am`.
  - Created `src/lib/rbac/permission-guard.tsx` (`<PermissionGuard>` component).
  - ESLint passed with 0 errors; Turbopack production build verified cleanly.

---

## 2. Next Immediate Tasks (Partition 05-02: Fleet Registration & Taxonomy)

1. **Database Schema & Migrations (`migrations/`):**
   - Create `migrations/000004_create_fleet_and_machinery_tables.up.sql` and `.down.sql`:
     - `categories`: Two-level taxonomy (`parent_id`, `code`, `name_en`, `name_am`, `icon_url`, `display_order`).
     - Dynamic EAV: `spec_field_groups`, `spec_field_definitions`, `fleet_spec_values`.
     - `fleets`: PostGIS geography yard point (`location`), baseline meter reading (`NUMERIC(10, 2)`), pricing tariffs, cancellation penalties, operator & logistics flags.
     - `fleet_photos`: Mandatory 6-angle walkaround (`FRONT`, `REAR`, `LEFT`, `RIGHT`, `CABIN`, `METER`).
     - `fleet_spec_indexed`: Logistics-critical denormalized performance table.
     - `owner_operator_memberships` & `fleet_operator_assignments`.
   - Update `backend/internal/platform/database/queries.sql` and generate with `sqlc`.
2. **Backend Domain & Application Core (`backend/internal/modules/machinery/`):**
   - Domain models, monotonic meter non-regression, 6-photo completeness validator, automated transport inference matrix.
   - Antivirus media stream scanning service via ClamAV (:3310 `zINSTREAM`) and SeaweedFS S3 storage (`D-05-02-001`).
   - Use cases for category taxonomy, dynamic spec definition retrieval, and fleet asset creation guarded by RBAC `fleets:create`.
3. **OpenAPI Contracts & HTTP Transport (`contracts/` & `backend/`):**
   - Wire `/categories`, `/categories/{id}/spec-definitions`, and `/fleets` into HTTP server.
4. **Flutter Mobile Dynamic Forms & Media Capture (`mobile/lib/`):**
   - Implement metadata-driven dynamic EAV form generator and walkaround photo capture reticle.
   - Wire `FleetBloc` state machine and `isar_community` local draft storage.
5. **Quality Gates & Certification:**
   - 100% unit, integration, and widget test pass rate.

---

## 3. Active Decisions & Open Blockers

### Decisions Formulated & Enforced
- `ADR-001` through `ADR-019` in `docs/decisions/`.
- `ADR-013`: Server-Driven OTA Translation Bundles & Dynamic Ethiopic Font Delivery Pipeline.
- `ADR-018`: Dynamic Role-Based Access Control (RBAC) and Granular Permission Architecture.
- `ADR-019`: Full-Stack Vernacular Localization (i18n/l10n) and Ge'ez Script Architecture.
- `D-05-01-001`, `D-05-01-002`, `D-05-02-001`, `D-05-02-002`.

### Blockers
- None. Dynamic RBAC and Full-Stack Localization foundations are active and certified.
