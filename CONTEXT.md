# CONTEXT.md — Dink Machinery Active Session State

**Project Phase:** Phase 2 (Partition Feature Development)  
**Active Requirement Partition:** `05-03-attachment-registration-inventory` (Attachment Registration & Inventory Architecture)  
**Active Branch:** `feature/05-03-attachment-registration-inventory`  
**Current Stop Point:** Partition 05-03 Fully Implemented & Verified Across PostgreSQL GiST, Go Backend Core, OpenAPI Contracts, and Flutter Mobile

---

## 1. Current Progress Snapshot

### Completed (Partition 05-03: Attachment Registration & Inventory Architecture)
- **Database Schema & Migrations (`migrations/`):**
  - Created `migrations/000005_create_attachment_tables.up.sql` and `.down.sql`:
    - Enabled `btree_gist` extension for composite GiST exclusion constraints.
    - `attachments`: Inventory implements (`HYDRAULIC_BREAKER`, `STANDARD_BUCKET`, `TRENCHING_BUCKET`, `AUGER`, `RIPPER_TOOTH`, `SORTING_GRAPPLE`, `CRUSHER`), operating weight, hydraulic boundaries (flow min/max L/min, pressure bar), mechanical coupler pins (diameter mm, spread mm), daily rate minor units, security deposit minor units, status lifecycle.
    - `attachment_compatible_categories`: Mappings to machinery taxonomy subcategories.
    - `attachment_photos`: Implement conditioning photos enforcing 4 mandatory angles (`OVERALL`, `BRACKET`, `PINS`, `TOOL_TIP`).
    - `fleet_bundled_attachments`: Overrides linking attachments to host fleet prime movers with bundled tariffs (`FL-32`, `FL-33`).
    - `attachment_reservations`: Storage-layer temporal GiST exclusion constraint (`EXCLUDE USING gist (attachment_id WITH =, hire_period WITH &&)`) atomically preventing concurrent overlapping bookings between standalone and host-bundled hires (`FL-34`, `D-05-03-001`).
    - Dynamic RBAC permissions: Seeded `attachments:create`, `attachments:read`, `attachments:update`, `attachments:delete`, `attachments:link` mapped to roles (`owner`, `admin`, `super_admin`).
  - Generated `sqlc` queries in `backend/internal/platform/database/db/` (`models.go`, `queries.sql.go`).
- **OpenAPI 3.1.0 Contracts & Type Generation:**
  - Added schemas in `contracts/components/schemas/attachment.yaml`.
  - Added endpoint contracts:
    - `POST /attachments` (`contracts/endpoints/machinery/attachments.yaml`)
    - `GET /attachments/inventory` (`contracts/endpoints/machinery/attachments-inventory.yaml`)
    - `POST /fleets/{id}/attachments/link` (`contracts/endpoints/machinery/fleet-attachments-link.yaml`)
    - `POST /bookings/check-bundled-availability` (`contracts/endpoints/bookings/check-bundled-availability.yaml`)
  - Updated `contracts/openapi.yaml`, bundled via Redocly (`admin/package.json`), generated TypeScript DTOs (`admin/src/types/api.generated.ts`), and Go chi server types with `oapi-codegen` (`backend/internal/transport/http/generated/api.gen.go`).
- **Go Backend Core Subsystem (`backend/`):**
  - `internal/modules/machinery/domain/`:
    - Domain entities: `Attachment`, `AttachmentPhoto`, `FleetBundledAttachment`, `AttachmentReservation`.
    - Invariant validators: `ValidateAttachmentPhotos` (mandating 4 conditioning angles), `ValidateCarrierCompatibility` (carrier weight $\pm 10\%$, pin diameter $\pm 2\text{mm}$, hydraulic pressure and flow envelopes), and `ValidateBundledRate` (bundled rate $\le$ standalone inventory rate).
    - Unit tests in `compatibility_test.go` verifying compatibility rejections, photo constraints, and bundled pricing rules.
  - `internal/modules/machinery/application/`:
    - `AttachmentService`: `RegisterAttachment`, `ListOwnerInventory`, `LinkAttachmentToFleet`, `CheckBundledAvailability`, `ReserveAttachment`.
    - Integration concurrency tests in `attachment_booking_test.go` verifying that reserving a bundled host excavator atomically blocks standalone rental of the attached breaker over overlapping dates (`D-05-03-001`, `FL-34`).
  - `internal/modules/machinery/infrastructure/`:
    - `MemoryMachineryRepository` with thread-safe atomic range overlap checks.
    - `PostgresMachineryRepository` implementing `AttachmentRepository` with `pgtype.Range[pgtype.Date]` storage mapping.
  - `internal/transport/http/`:
    - Handlers: `CreateAttachment`, `ListAttachmentInventory`, `LinkFleetAttachment`, `CheckBundledAvailability`.
    - Injected `AttachmentService` into `cmd/api/main.go` and test suites.
    - 100% backend test pass rate across all packages (`go test ./...`).
- **Flutter Mobile Client (`mobile/`):**
  - Full-stack localization: English (`app_en.arb`) and verified Amharic (`app_am.arb`) with zero hardcoded user-visible strings.
  - Data layer: `AttachmentModel`, `AttachmentInventoryItemModel`, `AttachmentRemoteDataSourceImpl`, `AttachmentLocalDataSourceImpl` (draft caching via `FlutterSecureStorage`), `AttachmentRepositoryImpl`.
  - Domain layer: Entities `Attachment`, `AttachmentPhoto`, `AttachmentType`, `AttachmentPhotoAngle`, `AttachmentInventoryItem`, `BundledAttachmentLink`, `AttachmentRepository`, and use cases `RegisterAttachmentUseCase`, `GetAttachmentInventoryUseCase`, `LinkAttachmentUseCase`.
  - Presentation layer:
    - `AttachmentBloc` handling draft caching, photo capture, registration submission, inventory loading, and fleet bundling.
    - `AttachmentPhotoReticle` displaying the 4 mandatory conditioning camera angles with completion badges.
    - `AttachmentRegistrationPage` with full dimensional, hydraulic, and pricing form inputs.
    - `AttachmentInventoryPage` displaying active attachments, bundled badges, and status pills.
  - Registered all dependencies in `lib/core/di/injection.dart` (`GetIt`).
  - 100% mobile test pass rate (33/33 tests passing); 0 issues on `fvm flutter analyze`.
- **Governance & Quality Gates:**
  - `scripts/validate-decision-template.sh`: Verified all decisions conform strictly to governance specifications.
  - `scripts/validate-decision-links.sh`: Verified all 19 ADRs and 80 L2 partition decisions have valid bidirectional links.

---

## 2. Next Immediate Tasks (Partition 05-04: Operator Verification & Fleet Assignment Architecture)

1. **Database Schema & Migrations (`migrations/`):**
   - Create `migrations/000006_create_operator_tables.up.sql` and `.down.sql`:
     - `operators`: Operator profile, government certificate verification (`OP-01`, `OP-02`), competency grading, operator badge status.
     - `operator_certifications`: Heavy machinery operating license grades, issuing authorities, expiration dates, OCR verified documents.
     - `fleet_operator_assignments`: Assignment of vetted operators to fleet prime movers with assignment history and active pairing.
     - Dynamic RBAC permissions for operators.
2. **OpenAPI 3.1.0 Contracts:**
   - Define endpoints for operator onboarding, certificate document upload, admin certification verification, and fleet assignment.
3. **Backend Domain & Application Core (`backend/internal/modules/`):**
   - Implement certificate validity verification, license class to equipment class mapping, and operator-to-fleet assignment service.
4. **Flutter Mobile Client (`mobile/lib/features/operator/`):**
   - Implement operator profile registration, license scanning/capture, and owner fleet operator assignment views.
