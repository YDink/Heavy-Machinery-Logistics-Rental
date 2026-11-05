# CONTEXT.md — Dink Machinery Active Session State

**Project Phase:** Phase 2 (Partition Feature Development)  
**Active Requirement Partition:** `05-07-rental-completion-return-damage` (Rental Completion, Off-Hire Return Inspection, Damage & Penalty Adjudication Architecture)  
**Active Branch:** `feature/05-07-rental-completion-return-damage` (Active development underway)  
**Current Stop Point:** Enabled Swagger UI (`/swagger/`, `/swagger/index.html`, `/swagger/openapi.yaml`) in the Go backend modular monolith with pure standard library HTTP handlers, runtime disk hot-reload, embedded fallback bundle, and environment toggle (`SWAGGER_ENABLED`). Full test suite passing. Next immediate action: implement Partition 05-07 database schema and migration (`migrations/000009_create_rental_completion_tables.up.sql`).

---

## 1. Current Progress Snapshot

### Completed (Swagger UI & Interactive OpenAPI Documentation Delivery)
- **Configuration & Toggling (`backend/internal/platform/config/`):**
  - Added `SwaggerEnabled bool` to `Config` with `SWAGGER_ENABLED` environment parsing (defaults to `true` in development/staging, `false` in production).
  - 100% unit test coverage in `config_test.go` covering defaults, explicit toggles, and production enforcement.
- **Go Backend Transport Handlers (`backend/internal/transport/http/`):**
  - Strictly adhering to single-concept file separation:
    - `swagger_ui_handler.go`: Serves interactive Swagger UI HTML dashboard (`GET /swagger/`, `GET /swagger/index.html`).
    - `swagger_redirect_handler.go`: Redirects bare `/swagger` to `/swagger/`.
    - `swagger_spec_handler.go`: Serves canonical OpenAPI 3.1.0 contract (`GET /swagger/openapi.yaml`) with runtime filesystem discovery for dev hot-reload and `//go:embed openapi.bundled.yaml` fallback.
    - `server.go`: Added `swaggerEnabled` state with `SetSwaggerEnabled(bool)` and dynamic route mounting on native `net/http.ServeMux`.
    - `cmd/api/main.go`: Wired `server.SetSwaggerEnabled(cfg.SwaggerEnabled)`.
  - Comprehensive unit and integration test suite in `swagger_handler_test.go` verifying 200 OK HTML payload, YAML spec headers, 301 redirect, and 404 when disabled (100% pass rate).

### Completed (System Manual Live Verification & SQL Fixture Alignment Drill)
- **Database Schema & Migrations (`migrations/`):**
  - Aligned RBAC column references in `000005`, `000006`, `000007`, and `000008`: normalized to canonical `resource` column on `permissions` table.
  - Aligned rollback script `000008_create_active_rental_telemetry_tables.down.sql` to reference `code` column.
  - Executed migration rollback drill on migration 000008 (`migrate down 1` -> `migrate up 1`) passing cleanly with zero relational collisions.
- **Multi-Machine Test Data Parity (`scripts/fixtures/`):**
  - Repaired column mapping in `scripts/fixtures/05_03_attachment_fixtures.sql` (`is_active`).
  - Repaired PostgreSQL date range type casting in `scripts/fixtures/05_05_booking_fixtures.sql` (`(CURRENT_DATE - INTERVAL '...')::date`).
  - Repaired telemetry columns in `scripts/fixtures/05_06_active_rental_fixtures.sql` (`start_time`, `end_time`, `accumulated_hours`).
  - Configured default database name to `dink_machinery_dev` across `apply-fixtures.ps1` and `apply-fixtures.sh`.
  - Applied all 7 SQL fixture scripts with 100% idempotency into `dink_machinery_dev`.
- **Mandatory Manual Live API Probing (Interactive Verification):**
  - Executed direct one-by-one live HTTP requests against Go API server connected to PostgreSQL `dink_machinery_dev`:
    - 2xx Success probes: Healthz, OTP Request, OTP Verify, JWT Authenticated Profile, Bilingual Category Localization (English & Amharic), Fleet Lookup (CAT 320D), Attachment Inventory (Soosan SB81N), 50km Spatial Radius Search, Bookings List, and Active Rental HUD Telemetry (1250.50 check-in meter, 3.75h accumulated hours).
    - RFC 9457 Problem Details probes: 400 Bad Request (malformed JSON), 404 Not Found (non-existent fleet UUID), 401 Unauthorized (missing Bearer token).
    - Result: 100% of tested endpoints passed with standardized response envelopes.
- **Client & Subsystem Verification:**
  - Mobile: Full Flutter test suite passing (72 / 72 tests passed).
  - Backend: Full Go unit and integration suite passing (100% across all modules).
  - Admin: ESLint passing with zero warnings/errors.
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
  - Updated `contracts/openapi.yaml`, bundled via Redocly (`contracts/openapi.bundled.yaml`), generated TypeScript DTOs (`admin/src/types/api.generated.ts`), and Go types with `oapi-codegen` (`backend/internal/transport/http/generated/api.gen.go`).
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
  - 100% mobile test pass rate; 0 issues on `fvm flutter analyze`.

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

### Completed (Partition 05-05: Booking Reservation Flow, 16-State Lifecycle FSM & Hold Concurrency Engine)
- **Database Schema & Migrations (`migrations/`):**
  - Created `migrations/000007_create_booking_state_machine_tables.up.sql` and `.down.sql`:
    - `bookings`: 16-state canonical lifecycle (`DRAFT`, `PENDING_PAYMENT`, `PAYMENT_EVIDENCE_UPLOADED`, `PAYMENT_DECLINED`, `DISPUTE_OPENED`, `DISPUTE_RESOLVED_RENTER`, `DISPUTE_RESOLVED_OWNER`, `CONFIRMED`, `AWAITING_FLEET_READINESS`, `PENDING_ADDITIONAL_PAYMENT`, `READY_FOR_PICKUP`, `ACTIVE`, `RETURN_INITIATED`, `PENDING_FINAL_PAYMENT`, `COMPLETED`, `RATED`, `CANCELLED`).
    - Optimistic Concurrency Control (OCC) column `version INT NOT NULL DEFAULT 1`.
    - Strict monetary minor units (`daily_rate_minor`, `base_hire_fee_minor`, `attachments_fee_minor`, `gross_rental_amount_minor`, `prepayment_amount_minor`, `penalty_amount_minor`, `refund_amount_minor`).
    - 4-digit site handover verification PIN (`handover_pin VARCHAR(8)`).
    - Cabin engine hour meters (`check_in_meter_hours NUMERIC(10,2)`, `check_out_meter_hours NUMERIC(10,2)`).
    - Keyset indexes on `(fleet_id, start_date, end_date)`, `(renter_id, created_at)`, `(owner_id, created_at)`.
    - `booking_attachments`: Embedded bundled attachment rates.
    - `booking_bank_slips`: Non-custodial wire receipts with SHA-256 duplicate proof.
    - `booking_audit_logs`: Immutable, append-only state transition audit trail with OCC version recording and actor IDs.
    - Seeded dynamic RBAC permissions: `bookings:create`, `bookings:read`, `bookings:confirm`, `bookings:decline`, `bookings:dispute`, `bookings:cancel`, `bookings:checkin`, `bookings:checkout` mapped to standard roles.
  - Added 15 sqlc queries in `backend/internal/platform/database/queries.sql`, compiled to Go in `backend/internal/platform/database/db/` (`models.go`, `queries.sql.go`).
- **OpenAPI 3.1.0 Canonical Contracts & Type Generation:**
  - Updated `contracts/components/schemas/booking.yaml` with 16-state schema, monetary envelopes, and DTOs.
  - Added endpoint contracts: `contracts/endpoints/bookings/payment-evidence.yaml`, `decline-payment.yaml`, `dispute-booking.yaml`, `handover-check-in.yaml`, `return-check-out.yaml`.
  - Re-bundled OpenAPI schema with Redocly (`contracts/openapi.bundled.yaml`).
  - Generated TypeScript DTOs (`admin/src/types/api.generated.ts`) and Go types (`backend/internal/transport/http/generated/api.gen.go`).
- **Go Backend Core Subsystem (`backend/`):**
  - `internal/modules/booking/domain/`:
    - Strict 16-state FSM transition validator (`ValidateTransition`).
    - Deterministic cost quoting engine (`CalculateQuote`) enforcing 50% advance prepayment and duration invariants.
    - Bilateral cancellation matrix (`CalculateCancellation`) enforcing 20%/50% penalty tiers, owner debt accounting, and force majeure (`BK-07`, `BK-08`).
    - Domain entities (`Booking`, `BankSlip`, `AuditLog`) and interfaces (`BookingRepository`, `HoldLeaseManager`).
    - 100% unit test pass rate (`fsm_test.go`, `quote_test.go`).
  - `internal/modules/booking/infrastructure/`:
    - `MemoryHoldLeaseManager`: Volatile Redis/Valkey 60-minute distributed mutex hold leases.
    - `PostgresRepository`: Full `pgx/v5` transactional persistence with OCC conflict detection.
    - `MemoryBookingRepository`: In-memory thread-safe repository for unit/integration testing.
  - `internal/modules/booking/application/`:
    - `BookingService` orchestrating reservation creation with distributed lease locks, bank slip upload, owner confirmation, payment dispute handling, bilateral cancellations, and physical handover/return meter validations.
    - 100% application test pass rate (`booking_service_test.go`).
  - `internal/transport/http/`:
    - Standard library `net/http.ServeMux` REST and RPC colon action endpoints (`POST /api/v1/bookings`, `GET /api/v1/bookings/{id}`, `POST /api/v1/bookings/{id}:confirm-payment`, `POST /api/v1/bookings/{id}:decline-payment`, `POST /api/v1/bookings/{id}:dispute`, `POST /api/v1/bookings/{id}:cancel`, `POST /api/v1/bookings/{id}:handover-check-in`, `POST /api/v1/bookings/{id}:return-check-out`, `POST /api/v1/bookings/{id}/payment-evidence`).
    - Full HTTP integration test suite (`booking_test.go`) verifying auth, self-booking prevention, 409 collision blocking, 422 bad PIN, 201 creation, and meter progression.
  - 100% backend test pass rate across all packages (`go test -count=1 ./...`).
- **Flutter Mobile Client (`mobile/`):**
  - Full-stack localization: English (`app_en.arb`) and Amharic (`app_am.arb`) with zero hardcoded strings.
  - Data layer: `BookingModel`, `BookingAttachmentModel`, `BankSlipModel`, `BookingRemoteDataSourceImpl`, `BookingRepositoryImpl`.
  - Domain layer: `BookingEntity` (with `copyWith`), `BookingAttachmentEntity`, `BankSlipEntity`, use cases (`CreateBookingUseCase`, `GetBookingUseCase`, `SubmitPaymentSlipUseCase`).
  - Presentation layer:
    - `BookingBloc` handling hold reservations, slip uploads, owner confirmations, disputes, cancellations, and handovers.
    - `BookingPriceBreakdownCard` rendering monetary amounts in minor units, daily tariffs, and 50% advance splits.
    - `HoldCountdownTimer` displaying real-time 60-minute reservation expiration.
    - `BookingStatusBadge` supporting 16-state lifecycle with color mapping and bilingual labels.
    - `BookingCheckoutPage` and `BookingDetailPage`.
  - Dependency Injection: Registered in `mobile/lib/core/di/injection.dart` (`GetIt`).
  - Comprehensive unit and widget tests: 19 tests in `test/features/booking/` (100% pass rate).
  - 100% mobile test pass rate (59/59 tests passing); 0 issues on `fvm flutter analyze`.

### Completed (Partition 05-06: Active Rental Management, Shift Telemetry, Breakdown & Extension Architecture)
- **Database Schema & Migrations (`migrations/`):**
  - Created `migrations/000008_create_active_rental_telemetry_tables.up.sql` and `.down.sql`:
    - `handover_inspections`: Pre-rental inspection with 5 mandatory photo angles (`CABIN`, `CHASSIS`, `HYDRAULICS`, `TRACKS_TIRES`, `HOUR_METER`), baseline meter reading, fluid levels/safety glass checks, site GPS location, signed confirmation (`AR-01`, `D-05-06-002`).
    - `handover_challenges`: Dual-party 6-digit random PIN generation with 15-minute TTL, max 3 attempts, contractor GPS capture, touch signature verification (`AR-02`).
    - `rental_shifts`: Shift telemetry logging with 8.0 hr daily shift quota baseline, start/end timestamps, operator assignment, operating hours meter deltas, excess hours accounting (`AR-03`).
    - `rental_extensions`: Extension requisitions with GiST temporal range availability conflict checking (`CheckFleetAvailabilityConflict`), 50% advance prepayment calculation, owner approval workflow (`AR-04`).
    - `mechanical_breakdowns`: Stoppage logging with breakdown photo evidence, automated shift telemetry timer pausing (`ACTIVE_BREAKDOWN_PAUSED`), 4-hour owner service SLA clock, resolution restore workflow (`AR-05`).
    - Seeded dynamic RBAC permissions: `rentals:inspect`, `rentals:challenge`, `rentals:hud`, `rentals:extend`, `rentals:breakdown`, `rentals:resolve` mapped to standard roles.
  - Added queries to `backend/internal/platform/database/queries.sql` (`CheckFleetAvailabilityConflict`, `GetFleetBaselineMeter`), compiled cleanly via `sqlc generate` into `backend/internal/platform/database/db/` (`models.go`, `queries.sql.go`).
- **OpenAPI 3.1.0 Canonical Contracts & Type Generation:**
  - Added schema models in `contracts/components/schemas/rental.yaml`: `HandoverInspectionRequest`, `HandoverConfirmRequest`, `ActiveRentalHUDResponseData`, `RentalExtensionRequest`, `MechanicalBreakdownRequest`, `BreakdownResolutionRequest`.
  - Added 7 endpoint contracts under `contracts/endpoints/rentals/`:
    - `POST /api/v1/bookings/{id}:handover-inspect`
    - `POST /api/v1/bookings/{id}:handover-challenge`
    - `POST /api/v1/bookings/{id}:handover-confirm`
    - `GET /api/v1/bookings/{id}/hud`
    - `POST /api/v1/bookings/{id}:extend`
    - `POST /api/v1/bookings/{id}:breakdown`
    - `POST /api/v1/bookings/{id}:resolve-breakdown`
  - Redocly bundled contract (`contracts/openapi.bundled.yaml`) validated with Spectral CLI (0 errors).
  - Generated TypeScript DTOs (`admin/src/types/api.generated.ts`) and Go types (`backend/internal/transport/http/generated/api.gen.go`).
- **Go Backend Core Subsystem (`backend/`):**
  - Pure standard library `net/http.ServeMux` and `log/slog` with zero third-party web frameworks.
  - `internal/modules/rental/domain/`:
    - Value objects and entities: `HandoverInspection`, `HandoverChallenge`, `RentalShift`, `RentalExtension`, `MechanicalBreakdown`, `ActiveRentalHUD`.
    - Invariant validators: `ValidateMonotonicMeterFloor` (regression prevention), `ValidateHandoverPhotos` (5 mandatory angles), `GenerateChallengePIN` (crypto 6-digit PIN with 15-min TTL), `CalculateShiftPacing`, `CalculateExtensionQuote` (50% advance prepay), `ValidateFaultCategory`.
    - 100% domain unit test pass rate (`invariants_test.go`).
  - `internal/modules/rental/infrastructure/`:
    - `MemoryRentalRepository`: Thread-safe in-memory repository for deterministic testing.
    - `PostgresRentalRepository`: Full `pgx/v5` transactional persistence with availability conflict detection.
  - `internal/modules/rental/application/`:
    - `RentalService` implementing `InspectHandover`, `GenerateHandoverPIN`, `ConfirmHandover`, `GetActiveRentalHUD`, `RequestExtension`, `ReportMechanicalBreakdown`, `ResolveMechanicalBreakdown`.
    - 100% application unit test pass rate (`rental_service_test.go`).
  - `internal/transport/http/`:
    - Handlers in `rental_handler.go` with standard unified response envelopes and RFC 9457 Problem Details error models.
    - Integration tests in `rental_test.go` verifying PIN generation, handover verification, HUD telemetry, extension conflict handling, and breakdown timer pausing/restoration.
  - 100% backend test suite pass rate (`go test -count=1 ./...`).
- **Flutter Mobile Client (`mobile/`):**
  - Full-stack localization: English (`app_en.arb`) and Amharic (`app_am.arb`) with 32 verified keys and zero hardcoded strings (`fvm flutter gen-l10n`).
  - Data layer: `ActiveRentalRemoteDataSourceImpl` and `ActiveRentalRepositoryImpl`.
  - Domain layer: Entities (`ActiveRentalHUD`, `HandoverInspection`, `HandoverChallenge`, `RentalExtension`, `MechanicalBreakdown`), Repository interface, 7 Use Cases.
  - Presentation layer:
    - `ActiveRentalBloc` handling HUD telemetry, inspections, PIN challenges, extensions, and breakdowns.
    - `ActiveRentalHUDView`: High-contrast dark styling (`#111827`) optimized for outdoor sunlight conditions, 56dp minimum touch targets.
    - `ShiftTelemetryProgressCard`: 8.0 hr daily shift progress bar, pacing indicator, and breakdown pause banner.
    - `BreakdownReportModal`: Fault category dropdown, description, DTC code inputs, meter reading, 56dp action button.
    - `ExtensionRequestSheet`: Forward date selector, 50% prepay calculation, 56dp submit button.
  - Registered all dependencies in `mobile/lib/core/di/injection.dart`.
  - 100% mobile test pass rate (72/72 tests passing); 0 issues on `fvm flutter analyze`.
- **Governance & Quality Gates:**
  - `scripts/validate-decision-template.sh`: All architectural decision templates strictly adhere to governance specifications.
  - `scripts/validate-decision-links.sh`: All 19 ADRs and 80 L2 partition decisions have valid bidirectional links.
  - Spectral OpenAPI contract validation: 0 errors on bundled contract.

### Completed (System-Wide Architecture Remediation & Multi-Machine Parity)
- **Multi-Machine SQL Fixtures (`scripts/fixtures/`):**
  - Created idempotent SQL fixtures (`00_base_system_fixtures.sql`, `05_01_auth_fixtures.sql`, `05_02_fleet_fixtures.sql`, `05_03_attachment_fixtures.sql`, `05_04_discovery_fixtures.sql`, `05_05_booking_fixtures.sql`, `05_06_active_rental_fixtures.sql`) using deterministic static UUIDs for dual-machine parity (Windows at office ↔ Mac at home).
  - Created fixture runner scripts (`apply-fixtures.sh` and `apply-fixtures.ps1`).
- **Backend Decomposition (Single-Concept File Separation):**
  - Purged all monolithic `entity.go`, `repository.go`, `invariants.go`, `postgres_repository.go`, `memory_repository.go` across `auth`, `booking`, `machinery`, and `rental` domain modules.
  - Decomposed into 1:1 dedicated entity, value object, invariant, repository interface, and repository adapter files.
  - Purged monolithic `rental_handler.go` (648 lines); decomposed into dedicated single-concept HTTP action handlers (`handover_inspect_handler.go`, `generate_handover_challenge_handler.go`, `confirm_handover_handler.go`, `get_active_rental_hud_handler.go`, `request_rental_extension_handler.go`, `report_mechanical_breakdown_handler.go`, `resolve_mechanical_breakdown_handler.go`, `authenticate_user.go`).
  - Decomposed monolithic `server.go` (2,295 lines down to 88 lines); extracted all 28 action handlers and helpers into dedicated single-concept files (`request_otp_handler.go`, `verify_otp_handler.go`, `refresh_session_handler.go`, `get_current_user_profile_handler.go`, `list_categories_handler.go`, `get_fleet_by_id_handler.go`, `create_fleet_asset_handler.go`, `list_fleets_handler.go`, `create_attachment_handler.go`, `list_attachment_inventory_handler.go`, `link_fleet_attachment_handler.go`, `check_bundled_availability_handler.go`, `search_fleets_handler.go`, `compare_fleets_handler.go`, `get_fleet_availability_handler.go`, `search_fleets_by_radius_handler.go`, `create_booking_handler.go`, `list_bookings_handler.go`, `get_booking_by_id_handler.go`, `submit_booking_payment_evidence_handler.go`, `confirm_booking_payment_handler.go`, `decline_booking_payment_handler.go`, `dispute_booking_payment_handler.go`, `cancel_booking_handler.go`, `handover_check_in_handler.go`, `return_check_out_handler.go`, `dispatch_booking_action_handler.go`, `stub_handlers.go`, `to_generated_booking_helper.go`).
  - 100% backend unit and integration test pass rate maintained.
- **Mobile Architecture & Layer Isolation:**
  - Resolved root `.gitignore` tracking bug (`/data/` vs `data/`), restoring full tracking of Flutter models, data sources, and repository implementations.
  - Decoupled all Flutter data sources (`auth`, `fleet_catalog`, `discovery`, `attachment`, `active_rental`, `booking`) completely from domain entities; strictly accept/return data models/DTOs. Standardized naming to `_data_source.dart`.
  - Normalized mobile presentation architecture: screens isolated in `views/<name>_view.dart`, widgets in `views/widgets/<name>.dart` across all features with backward-compatible aliases in `pages/` and `widgets/`.
  - Decomposed lumped entity files (`booking.dart`, `fleet_asset.dart`, `spec_field_group.dart`, `fleet_compare_matrix.dart`, `fleet_search_filter.dart`, `fleet_search_result.dart`, `fleet_availability.dart`, `attachment.dart`) into dedicated single-concept entity files with clean exports.
  - 100% mobile test pass rate (72 / 72 tests passed).

---

## 2. Next Immediate Tasks (Partition 05-07: Rental Completion, Return Inspection, Damage & Penalty Adjudication Architecture)

### Mandatory Pre-Commit Quality Requirements
Per `AGENTS.md § 1` and `docs/quality/MANUAL_TESTING_PROTOCOL.md`, all partition tasks require:
1. **Automated Unit Tests:** 100% pass on pure domain invariants ($\ge 80\%$ line coverage overall; $\ge 95\%$ on financial/auth).
2. **Automated Integration Tests:** 100% pass on repository transactions and HTTP transport handlers.
3. **Mandatory Live Manual Testing:** Every new endpoint, migration, and UI widget must be manually verified live before commit, with execution proof recorded under `Manual-Tests:` in the commit body.
4. **UI Design-First with Stitch MCP:** Every mobile screen/modal and admin dashboard must be designed using Stitch MCP just-in-time at the start of active requirement development (strictly adhering to design standards and maintaining flow continuity from preceding screens in the user journey), using generated Stitch designs as the direct implementation target.
5. **Strict Granular File Separation ("Single Concept Per File"):** Zero shared or catch-all files (`entity.go`, `models.go`, `invariants.go`, `repository.go`, multi-handlers). Every single entity, value object, invariant, repository interface/impl, use case, route handler, widget, and test must be isolated in its own specifically named file.
6. **Strict Layer & Type Segregation ("No Type/Layer Mix-Ups"):** Absolute isolation across Domain Entities (zero serialization/DB tags), DTOs (`transport/http/dto/`, `data/models/`), DB Models (`infrastructure/persistence/models/`), Repository Interfaces (domain ports accepting/returning entities only), and Repository Implementations (explicit bidirectional mapping). Data sources must never touch domain entities.
7. **Multi-Machine Test Data Parity (Idempotent SQL Fixtures):** All manual test seed data, mock entities, and manual DB inserts must be committed in idempotent SQL scripts under `scripts/fixtures/` using deterministic static UUIDs, allowing instant synchronization at session start across Windows (office) and macOS (home).

### Implementation Tasks

1. **Database Schema & Migrations (`migrations/`):**
   - Create `migrations/000009_create_rental_completion_tables.up.sql` and `.down.sql`:
     - `return_inspections`: Off-hire return inspection with 5 mandatory return photo angles, final engine hour meter, site GPS location, signed off-hire acceptance.
     - `damage_assessments`: Physical damage incident records, itemized repair cost estimation, severity tiering, photo evidence.
     - `rental_settlements`: Final financial settlement ledger calculations (excess shift hours overtime surcharge, fuel deficit charges, security deposit release/deduction, owner payout balance).
     - Dynamic RBAC permissions for rental return and damage adjudication (`rentals:return_inspect`, `rentals:damage_claim`, `rentals:settle`).
   - Generate `sqlc` queries in `backend/internal/platform/database/db/`.
2. **OpenAPI 3.1.0 Canonical Contracts:**
   - Define contracts in `contracts/`:
     - `POST /api/v1/bookings/{id}:return-inspect`
     - `POST /api/v1/bookings/{id}:damage-claim`
     - `POST /api/v1/bookings/{id}:dispute-damage`
     - `POST /api/v1/bookings/{id}:settle`
     - `GET /api/v1/bookings/{id}/settlement-summary`
   - Bundle OpenAPI specs and generate TypeScript DTOs & Go types.
3. **Backend Domain & Application Core (`backend/internal/modules/settlement/` or `rental/`):**
   - Invariant validators: Final meter monotonic validation ($\ge$ check-in meter), overtime excess hour calculation ($> 8.0\text{ hrs/day}$ quota), deposit forfeiture calculation, settlement ledger balance.
   - Comprehensive unit and integration test coverage.
4. **Flutter Mobile Client:**
   - Return inspection checklist with photo capture reticles.
   - Damage reporting and dispute submission view.
   - Final settlement ledger card and escrow release acknowledgment.
   - 100% test coverage and zero analyzer issues.

---

## 3. Active Decisions & Open Blockers

### Active Architectural Decisions
- `ADR-004`: State machine and optimistic concurrency control for booking and settlement lifecycles.
- `ADR-010`: Double-entry minor unit financial ledger and non-custodial escrow settlements.
- `D-05-07-001`: Dual-party photographic handover vs. return image diffing (6 mandatory angles).
- `D-05-07-002`: 72-hour statutory administrative dispute resolution SLA engine.

### Open Questions & Blockers
- No active blockers for Partition 05-07.
- Working branch: `feature/05-07-rental-completion-return-damage`.
- Pre-commit rule active: Automated unit/integration tests + mandatory manual testing verification per `docs/quality/MANUAL_TESTING_PROTOCOL.md`.

