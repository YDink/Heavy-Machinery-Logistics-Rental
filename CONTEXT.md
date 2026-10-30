# CONTEXT.md — Dink Machinery Active Session State

**Project Phase:** Phase 2 (Partition Feature Development)  
**Active Requirement Partition:** `05-01-core-identity-auth` (Registration, Phone Verification & Authentication)  
**Active Branch:** `feature/05-01-core-identity-auth`  
**Current Stop Point:** Partition 05-01 Core Identity & Authentication Implementation Complete & Certified

---

## 1. Current Progress Snapshot

### Completed (Partition 05-01: Core Identity & Authentication)
- **Database Layer & Schema Invariants (`migrations/` & `backend/`):**
  - Created `migrations/000002_create_identity_and_auth_tables.up.sql` and `.down.sql`:
    - `users` table: E.164 phone uniqueness (`REG-01`), country code foreign key, role constraints, account and KYC status state machines (`REG-03`, `REG-04`).
    - `user_legal_agreements` table: Non-repudiation ledger for Terms & Privacy Policy acceptance (`REG-08`, `REG-09`).
  - Generated sqlc persistence models and queries in `backend/internal/platform/database/db/` (`User`, `UserLegalAgreement`).
- **Backend Domain & Application Core (`backend/internal/modules/auth/`):**
  - `domain/phone_number.go`: E.164 phone normalizer, country resolver (`ET` / `KE`), carrier detection (`Ethio Telecom`, `Safaricom ET`, `Safaricom KE`, `Airtel KE`), and zero-PII masking.
  - `application/password.go`: Argon2id key derivation matching SRS §3.2 ($m=64\text{MB}, t=3, p=4$, salt 16 bytes, key 32 bytes) with constant-time verification.
  - `application/otp_service.go`: 6-digit CSPRNG generation, 300s TTL (`AUT-01`, `AUT-02`), 3-attempt exhaustion invalidation (`AUT-03`), and distributed leaky-bucket rate limiter (`D-05-01-001`).
  - `application/token_service.go`: Asymmetric ES256 JWT access token minting (15m TTL) and Refresh Token Rotation (30d TTL) with instant token family revocation on replay detection (`D-05-01-002`, `ADR-014`).
  - `application/auth_service.go`: Use cases coordinating request OTP, verify OTP, session refresh, and profile retrieval.
  - `infrastructure/`: Thread-safe memory cache and mock SMS gateway for testing; `PostgresUserRepository` for live database access.
- **Backend HTTP Transport & OpenAPI Binding (`backend/internal/transport/http/`):**
  - Implemented `generated.ServerInterface` on `transport.Server` handling `/auth/otp/request`, `/auth/otp/verify`, `/auth/session/refresh`, `/auth/users/me`.
  - Created `internal/platform/response/response.go` formatting universal envelopes and RFC 9457 Problem Details (`AGENTS.md § 2`).
  - Wired into `cmd/api/main.go` on `/api/v1` base route.
  - 100% backend unit and integration test pass rate.
- **Flutter Mobile Auth Subsystem (`mobile/lib/`):**
  - Clean Architecture scaffolding:
    - Data: `ApiClient` with Dio interceptor injecting `X-Country-Code` and `Authorization`, `SecureStorageService` wrapping `flutter_secure_storage`, `AuthRemoteDataSource`, `AuthRepositoryImpl`.
    - Domain: Entities `AuthTokens`, `AuthUser`; use cases `RequestOtpUseCase`, `VerifyOtpUseCase`, `GetCurrentUserUseCase`.
    - Presentation: `AuthBloc` state machine managing `UnauthenticatedState`, `OtpSentState`, `AuthenticatedState`, `AuthErrorState`.
  - 100% mobile test pass rate (10 tests: `app_config_test.dart`, `auth_bloc_test.dart`, `widget_test.dart`); zero warnings on `fvm flutter analyze`.
- **Next.js Admin Console (`admin/`):**
  - Turbopack production build verified cleanly in 719ms; ESLint passed with 0 errors.

---

## 2. Next Immediate Tasks Upon Resuming

1. **Partition 05-01 Mobile UI Pages:**
   - Implement `PhoneInputPage` with dynamic country code selection (+251 / +254) and telco badges.
   - Implement `OtpVerificationPage` with 6-digit reticle and 60-second animated cooldown timer.
2. **Next Partition: `05-02-fleet-registration-machinery-logistics`:**
   - Create feature branch `feature/05-02-fleet-registration-machinery-logistics`.
   - Implement fleet asset schema, SeaweedFS media ingestion, and machinery catalog endpoints.

---

## 3. Active Decisions & Open Blockers

### Decisions Formulated & Enforced
- `ADR-001` through `ADR-017` in `docs/decisions/`.
- `D-05-01-001`: Phone OTP Authentication with Leaky Bucket Rate Limiting and Carrier Failover.
- `D-05-01-002`: Dual-Token ES256 Minting with Valkey Refresh Token Family Rotation.

### Blockers
- None. Partition 05-01 core identity & auth pipeline is fully tested and certified.
