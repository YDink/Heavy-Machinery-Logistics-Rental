# CONTEXT.md — Dink Machinery Active Session State

**Project Phase:** Phase 2 (Partition Feature Development)  
**Active Requirement Partition:** `05-05-booking-reservation-flow` (20-State Booking Machine & Hold Concurrency Engine)  
**Active Branch:** `main` (Pending checkout of `feature/05-05-booking-reservation-flow`)  
**Current Stop Point:** Partition 05-04 Verified & Merged into Main. Ready to begin Partition 05-05.

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
  - Added endpoint contracts: `POST /attachments`, `GET /attachments/inventory`, `POST /fleets/{id}/attachments/link`, `POST /bookings/check-bundled-availability`.
  - Redocly bundled, generated TypeScript DTOs and Go chi server types.
- **Go Backend Core Subsystem (`backend/`):**
  - Domain entities, invariant validators (`ValidateAttachmentPhotos`, `ValidateCarrierCompatibility`, `ValidateBundledRate`), `AttachmentService`, `PostgresMachineryRepository` and `MemoryMachineryRepository`.
  - Concurrency tests in `attachment_booking_test.go` verifying dual-booking temporal GiST protection.
  - 100% backend test pass rate across all packages (`go test ./...`).
- **Flutter Mobile Client (`mobile/`):**
  - Full-stack localization: English (`app_en.arb`) and Amharic (`app_am.arb`).
  - Data & Domain layers: `AttachmentModel`, `AttachmentRepositoryImpl`, use cases.
  - Presentation: `AttachmentBloc`, `AttachmentPhotoReticle`, `AttachmentRegistrationPage`, `AttachmentInventoryPage`.
  - 100% mobile test pass rate; 0 issues on `fvm flutter analyze`.

### Completed (Partition 05-04: Fleet Discovery, Spatial Search, Multi-Variable Filtering & Dynamic Comparison Matrix)
- **Database Schema & Migrations (`migrations/`):**
  - Created `migrations/000006_create_discovery_search_tables.up.sql` and `.down.sql`:
    - Enabled `pg_trgm` extension for full-text search across makes, models, and depot locations.
    - Added `country_code` (VARCHAR(2)) and `search_vector` generated trigram column to `fleets`.
    - Added `specs_json` (JSONB) with GIN `jsonb_path_ops` index on `fleet_spec_indexed` (`D-05-04-002`).
    - Created `owner_aggregates` table for ratings and fulfillment tracking.
    - Created `fleet_insurances` table with status checking.
    - Created `fleet_availability` table with composite GiST exclusion constraint (`EXCLUDE USING gist (fleet_id WITH =, hire_period WITH &&)`).
    - Seeded dynamic RBAC permissions `fleets:search` and `fleets:compare` mapped to standard roles.
  - Added `sqlc` queries in `backend/internal/platform/database/queries.sql` (`SearchFleetsSpatial`, `GetComparableSpecsForFleets`, `GetFleetsForComparison`, `ListFleetAvailabilityRanges`, `ListFleetBlackoutsInRange`, `CreateFleetAvailability`, `UpsertOwnerAggregate`, `GetOwnerAggregate`, `UpsertFleetInsurance`, `GetFleetInsuranceByFleetID`).
  - Generated Go queries and models cleanly in `backend/internal/platform/database/db/`.
- **OpenAPI 3.1.0 Canonical Contracts & Type Generation:**
  - Added schema models in `contracts/components/schemas/discovery.yaml`: `FleetSearchSuccessResponse`, `FleetCompareRequest`, `FleetCompareSuccessResponse`, `FleetAvailabilitySuccessResponse`.
  - Added RFC 9457 error response definition `contracts/components/responses/unprocessable-entity.yaml`.
  - Added endpoint contracts:
    - `GET /fleets/search` (`contracts/endpoints/machinery/fleets-search.yaml`): Keyset cursor spatial search, distance radius, multi-facet filtering, deterministic scoring.
    - `POST /fleets/compare` (`contracts/endpoints/machinery/fleets-compare.yaml`): 2-to-3 machine side-by-side dynamic EAV spec matrix with difference highlighting (`SR-08`, `FL-15`).
    - `GET /fleets/{id}/availability` (`contracts/endpoints/machinery/fleet-availability.yaml`): Month/day availability heatmap and pre-flight hold diagnostics.
  - Updated `contracts/openapi.yaml`, bundled via Redocly (`contracts/openapi.bundled.yaml`), generated TypeScript DTOs (`admin/src/types/api.generated.ts`), and Go chi server types with `oapi-codegen` (`backend/internal/transport/http/generated/api.gen.go`).
- **Go Backend Core Subsystem (`backend/`):**
  - `internal/modules/machinery/domain/`:
    - Domain entities & value objects: `FleetSearchFilter`, `FleetSearchResult`, `FleetCompareResult`, `FleetAvailabilityInfo`.
    - Invariants & math algorithms: `ValidateSearchFilter`, `ValidateComparisonSet` (homogeneity enforcement), `CalculateRecommendationScore` ($0.40R + 0.30P + 0.20U + 0.10C$), `BuildComparisonMatrix`, `SortSearchResults`.
    - Domain unit tests in `discovery_test.go` (100% pass rate).
  - `internal/modules/machinery/application/`:
    - `DiscoveryService` with spatial search orchestration, Valkey hold lease pre-flight interception (`HoldLeaseChecker`, `MemoryHoldLeaseChecker`), dynamic comparison matrix generation, and availability range aggregation.
    - Application unit tests in `discovery_service_test.go` (100% pass rate).
  - `internal/modules/machinery/infrastructure/`:
    - `MemoryMachineryRepository` implementing `DiscoveryRepository` with Haversine distance, trigram text matching, blackout/reservation checking, and comparable spec pivoting.
    - `PostgresMachineryRepository` implementing `DiscoveryRepository`.
  - `internal/transport/http/`:
    - Handlers: `SearchFleets`, `CompareFleets`, `GetFleetAvailability`, `SearchFleetsByRadius`.
    - Integrated with `cmd/api/main.go` and full HTTP integration test suite (`discovery_test.go`).
    - 100% backend test pass rate across all packages (`go test ./...`).
- **Flutter Mobile Client (`mobile/`):**
  - Full-stack localization: English (`app_en.arb`) and verified Amharic (`app_am.arb`) with zero hardcoded user-visible strings.
  - Data layer: `FleetSearchResultModel`, `FleetCompareModel`, `FleetAvailabilityModel`, `DiscoveryRemoteDataSourceImpl`, `DiscoveryRepositoryImpl`.
  - Domain layer: Entities `FleetSearchFilter`, `FleetSearchResult`, `FleetCompareMatrix`, `FleetAvailability`, and use cases `SearchFleetsUseCase`, `CompareFleetsUseCase`, `GetFleetAvailabilityUseCase`.
  - Presentation layer:
    - `DiscoveryBloc` handling search filtering, clustered map/list toggling, 3-machine comparison selection, and real-time availability checks.
    - `FleetFeedCard` with 56x56dp tap targets and WCAG 2.2 AAA contrast readouts.
    - `MultiFacetFilterSheet` with radius slider, daily rate range, crew toggle, self-move toggle, and fuel chips.
    - `FleetComparisonMatrixView` side-by-side spec comparison table with difference highlighting.
    - `DiscoveryPage` and `FleetComparisonPage`.
  - Registered all dependencies in `lib/core/di/injection.dart` (`GetIt`).
  - 100% mobile test pass rate (40/40 tests passing); 0 issues on `fvm flutter analyze`.
- **Governance & Quality Gates:**
  - `scripts/validate-decision-template.sh`: Verified all decisions conform strictly to governance specifications.
  - `scripts/validate-decision-links.sh`: Verified all 19 ADRs and 80 L2 partition decisions have valid bidirectional links.

### Completed (Zero Web Framework Migration — Pure Go Standard Library Architecture)
- **Elimination of Web Frameworks & Third-Party Logging:**
  - Removed `github.com/go-chi/chi/v5` and `go.uber.org/zap` completely from `backend/go.mod` and all backend packages.
  - Migrated HTTP routing to native Go standard library `net/http.ServeMux` with Go 1.22+ method-and-path pattern matching (`mux.HandleFunc("POST /api/v1/auth/otp/request", ...)`).
  - Path parameters extracted strictly using native `r.PathValue()`.
  - Migrated structured application logging to Go standard library `log/slog` with contextual handlers (`slog.NewTextHandler` for Dev, `slog.NewJSONHandler` for Prod) and zero-PII masking.
  - Re-configured `oapi-codegen` to generate pure type DTOs (`-generate types`) without framework-coupled server scaffolding.
  - Auth, Machinery, Attachment, and Discovery transport handlers refactored to standard `func (s *Server) Handler(w http.ResponseWriter, r *http.Request)`.
  - Created standard library HTTP middlewares in `backend/internal/transport/http/middleware.go` (`RequestID`, `Recoverer`, `RequestLogger`).
  - 100% backend test pass rate across all packages (`go test ./...` and `go test -count=1 ./...`).
  - 100% mobile test pass rate (40/40 tests passing; `fvm flutter test`).
  - Updated `docs/TECH_STACK.md` and `docs/components/backend/AGENTS.md`.

---

## 2. Next Immediate Tasks (Partition 05-05: Booking, Reservation & Lifecycle State Machine Architecture)

1. **Database Schema & Migrations (`migrations/`):**
   - Create `migrations/000007_create_booking_state_machine_tables.up.sql` and `.down.sql`:
     - `bookings` table: 16-state/20-state lifecycle topology (`DRAFT`, `HELD_PENDING_WIRE`, `PAYMENT_EVIDENCE_SUBMITTED`, `CONFIRMED`, `DISPATCHED_TO_SITE`, `IN_TRANSIT`, `ON_SITE_ACTIVE`, `BREAKDOWN_PAUSED`, `EXTENSION_REQUESTED`, `RETURN_INITIATED`, `PENDING_FINAL_PAYMENT`, `COMPLETED`, `DISPUTE_RAISED`, `CANCELLED_NO_FEE`, `CANCELLED_WITH_PENALTY`, `REFUNDED`).
     - Composite keyset pagination indexing (`created_at`, `id`, `renter_id`, `owner_id`, `status`).
     - Optimistic Concurrency Control (OCC) `version INT NOT NULL DEFAULT 1`.
     - `booking_audit_logs` append-only audit trail with previous state, new state, actor, reason, timestamp.
     - `booking_bank_slips` table for non-custodial bank wire settlement and SHA-256 duplicate proof.
     - Dynamic RBAC permissions for booking transitions and administrative interventions.
   - Generate `sqlc` queries in `backend/internal/platform/database/db/`.
2. **OpenAPI 3.1.0 Contracts:**
   - Define canonical wire contracts in `contracts/`:
     - `POST /bookings` (`contracts/endpoints/bookings/create-booking.yaml`): Requisition creation & 60-min Valkey hold lock.
     - `POST /bookings/{id}:submit-slip` (`contracts/endpoints/bookings/submit-payment-slip.yaml`): Slip ingestion & evidence submission.
     - `POST /bookings/{id}:confirm-payment` (`contracts/endpoints/bookings/confirm-payment.yaml`): Owner receipt verification.
     - `POST /bookings/{id}:cancel` (`contracts/endpoints/bookings/cancel-booking.yaml`): Bilateral cancellation with penalty calculation.
     - `GET /bookings/{id}` & `GET /bookings`: Single booking inspector and keyset paginated history.
   - Bundle OpenAPI specs and generate TypeScript DTOs & Go types.
3. **Backend Domain & Application Core (`backend/internal/modules/booking/`):**
   - Implement finite state machine (FSM) with strict transition matrix validation (`ValidateTransition`).
   - Implement deterministic cost quoting & tariff breakdown (`BK-04`).
   - Implement 60-minute Valkey distributed hold lease (`BK-05`) with automatic expiration cleanup.
   - Implement ClamAV slip scanning hook and SHA-256 idempotency check.
   - Comprehensive unit and integration concurrency tests.
4. **Flutter Mobile Client (`mobile/lib/features/booking/`):**
   - Domain, Data, and Presentation layers:
     - 60-minute countdown hold timer with 15-min warning push notification.
     - Non-custodial bank slip upload and preview.
     - Booking checkout price breakdown and status HUD.
   - 100% test coverage and zero analyzer issues.
