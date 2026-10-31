# CONTEXT.md — Dink Machinery Active Session State

**Project Phase:** Phase 2 (Partition Feature Development)  
**Active Requirement Partition:** `05-02-fleet-registration-machinery-logistics` (Fleet Registration, Machinery & Logistics Taxonomy Architecture)  
**Active Branch:** `feature/05-02-fleet-registration-machinery-logistics`  
**Current Stop Point:** Partition 05-02 Fully Implemented & Verified Across Go Backend, Database, and Flutter Mobile

---

## 1. Current Progress Snapshot

### Completed (Partition 05-02: Fleet Registration, Machinery & Logistics Taxonomy)
- **Database Schema & Migrations (`migrations/`):**
  - Created `migrations/000004_create_fleet_and_machinery_tables.up.sql` and `.down.sql`:
    - `categories`: Hierarchical machinery taxonomy with two-level depth, English and Amharic names, icon URLs, and display ordering.
    - Dynamic EAV Schema: `spec_field_groups`, `spec_field_definitions` (DECIMAL, INTEGER, BOOLEAN, ENUM, STRING with unit symbols and boundary constraints), `fleet_spec_values`.
    - `fleets`: PostGIS geography yard point (`location geography(Point, 4326)`), baseline cabin meter reading (`NUMERIC(10, 2)`), daily rates, excess hour rates, cancellation timeouts, fuel types, operator inclusion, self-move capability.
    - `fleet_photos`: Mandatory 6-angle walkaround tracking (`FRONT`, `REAR`, `LEFT`, `RIGHT`, `CABIN`, `METER`).
    - `fleet_spec_indexed`: Logistics performance indexing table with dual-persistence for sub-10ms dimensional query matchmaking (`D-05-02-002`, `FL-14`).
    - `owner_operator_memberships` and `fleet_operator_assignments`.
    - Seeded initial taxonomy from SRS Appendix C (`EXCAVATION`, `EARTHMOVING`, `LIFTING`, `HAULING`, `LOGISTICS`) and spec field groups for mini excavators.
  - Generated sqlc persistence models and typed queries in `backend/internal/platform/database/db/`.
- **Go Backend Core (`backend/`):**
  - `internal/modules/machinery/domain/`:
    - Invariants for non-negative baseline meter (`ValidateBaselineMeter`).
    - Mandatory 6-photo completeness validator (`ValidateWalkaroundPhotos`).
    - Automated trailer class inference matrix according to SRS §7.2 payload thresholds (`InferLogisticsTransportClass`).
  - `internal/modules/machinery/application/`:
    - `UploadPhotoService` with ClamAV stream scanner (:3310 `zINSTREAM`) and SeaweedFS S3 storage (`D-05-02-001`).
    - `MachineryService` with `CreateFleet`, `GetFleetAsset`, `ListCategories`, and `GetCategorySpecDefinitions`, enforcing dual-persistence to `fleet_spec_indexed`.
  - `internal/modules/machinery/infrastructure/`:
    - `PostgresMachineryRepository` implementing `CategoryRepository`, `SpecDefinitionRepository`, `FleetRepository`, `FleetMediaRepository`, and `FleetSpecRepository`.
    - `MemoryMachineryRepository` supporting thread-safe in-memory execution and test seeding.
  - `internal/transport/http/`:
    - Wired `GET /api/v1/categories` with automatic Accept-Language localization (`en` ↔ `am`).
    - Wired `POST /api/v1/fleets` validating Bearer token, mandatory 6 walkaround photos, and baseline engine meter reading.
    - Wired `GET /api/v1/fleets/{id}` and `GET /api/v1/fleets`.
    - 100% backend test pass rate across all packages.
- **Flutter Mobile Client (`mobile/`):**
  - Fully localized with zero hardcoded strings: English (`app_en.arb`) and verified Amharic (`app_am.arb`).
  - Data layer: `FleetRemoteDataSourceImpl`, `FleetLocalDataSourceImpl` (offline registration draft caching), `FleetRepositoryImpl`.
  - Domain layer: Entities `Category`, `SpecFieldGroup`, `SpecFieldDefinition`, `FleetAsset`, `FleetPhoto`, and Use Cases `GetCategoriesUseCase`, `GetSpecDefinitionsUseCase`, `RegisterFleetUseCase`.
  - Presentation layer:
    - `FleetRegistrationBloc` managing category selection, dynamic spec fields loading, walkaround photo uploads, and offline draft storage.
    - `WalkaroundCameraReticle` tracking the 6 mandatory physical angles with visual indicators.
    - `DynamicSpecFormCompiler` compiling category technical specifications dynamically into validated form widgets.
    - `FleetRegistrationPage` providing the complete registration workflow.
  - 100% test pass rate (22/22 unit and widget tests passing); 0 issues on `flutter analyze`.

---

## 2. Next Immediate Tasks (Partition 05-03: Attachment Registration & Inventory)

1. **Database Schema & Migrations (`migrations/`):**
   - Create `migrations/000005_create_attachment_tables.up.sql` and `.down.sql`:
     - `attachment_types`: Attachment taxonomy and mechanical interface standards (pin diameter, ear spread, center-to-center distance, hydraulic pressure/flow limits).
     - `attachments`: Serial numbers, daily rates, compatible machine categories, photos, and current availability status.
     - `fleet_attachment_compatibility_matrix`: Sub-10ms lookup matching machines to attachments.
2. **Backend Domain & Application Core (`backend/internal/modules/attachment/`):**
   - Implement compatibility validation engine matching hydraulic and mechanical coupler envelopes.
3. **Flutter Mobile UI (`mobile/lib/features/attachment/`):**
   - Implement attachment registration and machine-attachment pairing workflows.
