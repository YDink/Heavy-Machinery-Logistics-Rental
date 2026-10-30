# Task 00b: Phase 1 Bootstrap — Infrastructure, Database Baseline & Subsystem Initialization

**Task Identifier:** `TASK-00b-BOOTSTRAP-SUBSYSTEMS`  
**Phase:** 1 (Bootstrap & Subsystem Initialization)  
**Target Completion State:** Fully operational development environment ready for Phase 2 Feature Partitions  
**Governing Documents:** `BOOTSTRAP.md`, `TECH_STACK.md`, `CODE_STANDARDS.md`, `API_CONVENTIONS.md`, `docs/quality/BOOTSTRAP_CHECKLIST.md`

---

## 1. Executive Summary & Objective

The objective of Task 00b is the deterministic, automated bootstrapping of all physical runtime subsystems for **Dink Machinery**. Transitioning from Phase 0 (Contract and Documentation Governance) into Phase 1 establishes the executable skeleton of the modular monolith backend, mobile client, administrative console, local container infrastructure, and pre-commit verification hooks.

Per platform governance invariants, **zero domain business logic** (e.g., booking transitions, payment provider integrations, search algorithms) may be authored during this task. Phase 1 is strictly restricted to runtime scaffolding, dependency lockfile pinning, database extension baselines, code generation from canonical contracts, and automated governance activation.

---

## 2. Invariant Rules During Bootstrap

All execution steps must strictly comply with the following architectural invariants:

1. **Zero Domain Logic:** Do not implement business rules or partition features. All endpoints must return static or stubbed platform envelopes (`/healthz`).
2. **Deterministic Lockfile Pinning:** Every scaffolded subsystem must commit its verified lockfile (`backend/go.sum`, `mobile/pubspec.lock`, `admin/pnpm-lock.yaml`).
3. **Contract Derivation:** Transport DTOs, route signatures, and client API types must be compiled directly from `contracts/openapi.yaml` via automated code generators (`oapi-codegen`, `openapi-typescript`).
4. **Zero-PII Static Assurance:** Logging configurations must enforce masking filters from day zero. No unmasked telephone numbers, TINs, or National IDs are permitted in structured outputs.
5. **Universal Envelopes:** Scaffolded endpoints must adhere strictly to the platform JSON success envelope or RFC 9457 Problem Details error envelope.

---

## 3. Subsystem Scaffolding Matrix

| Subsystem | Technology & Version | Directory | Primary Output Artifacts | Verification Command |
|---|---|---|---|---|
| **Local Infrastructure** | Docker Compose / Traefik 3.x | Root | `docker-compose.yml`, local volumes | `docker compose ps` |
| **Primary Persistence** | PostgreSQL 18.6 + PostGIS 3.6.4 | `migrations/` | `000001_init_platform_schema.up.sql` | `psql -c "\dx"` |
| **Ephemeral Cache** | Valkey 9.1.3-alpine | Root | Distributed lock & rate-limit cluster | `redis-cli ping` |
| **Backend Core** | Go 1.27.3 Modular Monolith | `backend/` | `go.mod`, `go.sum`, `cmd/api/main.go`, `sqlc.yaml` | `go test -race ./...` |
| **Cross-Platform Mobile** | Flutter 3.47.9 / Dart 3.9 | `mobile/` | `pubspec.yaml`, `pubspec.lock`, Clean Arch tree | `fvm flutter test` |
| **Admin Web Console** | Next.js 16 (React 19, App Router) | `admin/` | `package.json`, `pnpm-lock.yaml`, Tailwind v4 | `pnpm run build` |
| **Governance Hooks** | Lefthook 1.8+ / Commitlint | Root | `lefthook.yml`, `.github/commitlint.config.js` | `lefthook run pre-commit` |

---

## 4. Execution Plan & Sequential Procedures

### Step 1: Pre-Execution Verification Gate
Ensure that all Phase 0 contracts and validation scripts pass before generating files:
```bash
# Verify canonical OpenAPI contract
pnpm dlx @stoplight/spectral-cli lint contracts/openapi.yaml --ruleset .spectral.yaml

# Verify decision templates and links
bash scripts/validate-decision-template.sh
bash scripts/validate-decision-links.sh
```

---

### Step 2: Infrastructure Provisioning & Health Verification
Spin up local container services and confirm connection readiness:
```bash
# Provision container suite
docker compose up -d postgres valkey seaweedfs clamav traefik

# Wait for healthy daemon readiness
until docker compose exec postgres pg_isready -U postgres; do sleep 1; done
until docker compose exec valkey redis-cli ping | grep -q PONG; do sleep 1; done
```

---

### Step 3: Database Baseline & Multi-Country Configuration (`LOC-01-BE`)
Establish the initial database migration containing spatial extensions, UUID generation, audit triggers, and default country configurations:
```bash
# Apply baseline migration
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/dink_machinery_dev?sslmode=disable" up

# Verify PostGIS, pgcrypto, and uuid-ossp extension activation
docker compose exec postgres psql -U postgres -d dink_machinery_dev -c "SELECT extname, extversion FROM pg_extension WHERE extname IN ('postgis', 'uuid-ossp', 'pgcrypto');"

# Verify country_configs table and initial seeds (Ethiopia ET, Kenya KE)
docker compose exec postgres psql -U postgres -d dink_machinery_dev -c "SELECT country_code, dialing_prefix, default_currency, standard_tax_rate_bps FROM country_configs;"
```

---

### Step 4: Go 1.27.3 Backend Scaffolding & Code Generation
Initialize the modular monolith package hierarchy and compile contract-first transport stubs:
```bash
cd backend

# Initialize Go module and install pinned platform dependencies
go mod init [github.com/dinkmachinery/dink-core](https://github.com/dinkmachinery/dink-core)
go get [github.com/go-chi/chi/v5@v5.1.0](https://github.com/go-chi/chi/v5@v5.1.0)
go get [github.com/jackc/pgx/v5@v5.7.1](https://github.com/jackc/pgx/v5@v5.7.1)
go get [github.com/google/uuid@v1.6.0](https://github.com/google/uuid@v1.6.0)
go get go.uber.org/zap@v1.27.0
go get go.opentelemetry.io/otel@v1.32.0

# Scaffold Clean Architecture module directories
mkdir -p cmd/api \
  internal/platform/config \
  internal/platform/database/db \
  internal/platform/telemetry \
  internal/transport/http/middleware \
  internal/transport/http/generated \
  internal/modules/auth/domain \
  internal/modules/auth/application \
  internal/modules/auth/infrastructure \
  internal/modules/machinery/domain \
  internal/modules/machinery/application \
  internal/modules/machinery/infrastructure \
  internal/modules/booking/domain \
  internal/modules/booking/application \
  internal/modules/booking/infrastructure \
  internal/modules/payment/domain \
  internal/modules/payment/application \
  internal/modules/payment/infrastructure \
  internal/modules/geo/domain \
  internal/modules/geo/application \
  internal/modules/geo/infrastructure \
  internal/modules/admin/domain \
  internal/modules/admin/application \
  internal/modules/admin/infrastructure

# Compile OpenAPI contract into type-safe Chi router interfaces and DTOs
oapi-codegen -generate types,chi-server,spec \
  -package generated \
  -o internal/transport/http/generated/api.gen.go \
  ../contracts/openapi.yaml

# Compile sqlc queries against baseline PostgreSQL schema
sqlc generate

# Validate binary compilation and baseline tests
go build -o bin/api ./cmd/api
go test -v -race ./...

cd ..
```

---

### Step 5: Flutter 3.47.9 Mobile Workspace Scaffolding
Scaffold the mobile client with BLoC state management, Isar local database, and feature-first structuring:
```bash
fvm flutter create --org com.dinkmachinery --project-name dink_mobile mobile
cd mobile

# Install core runtime dependencies
fvm flutter pub add flutter_bloc bloc get_it dio isar isar_flutter_libs \
  freezed_annotation json_annotation flutter_secure_storage geolocator \
  equatable intl

fvm flutter pub add --dev build_runner freezed json_serializable isar_generator

# Scaffold feature-first directory layout
mkdir -p lib/core/network \
  lib/core/theme \
  lib/core/storage \
  lib/core/errors \
  lib/features/auth/presentation \
  lib/features/auth/domain \
  lib/features/auth/data \
  lib/features/fleet_catalog/presentation \
  lib/features/fleet_catalog/domain \
  lib/features/fleet_catalog/data \
  lib/features/booking/presentation \
  lib/features/booking/domain \
  lib/features/booking/data \
  lib/features/inspection/presentation \
  lib/features/inspection/domain \
  lib/features/inspection/data

# Verify analyzer and test suite
fvm flutter analyze
fvm flutter test

cd ..
```

---

### Step 6: Next.js 16 Admin Console Scaffolding & Contract Sync
Initialize the Next.js App Router application with React 19, Tailwind CSS v4, and dynamic table components:
```bash
pnpm create next-app admin \
  --typescript \
  --tailwind \
  --eslint \
  --app \
  --src-dir \
  --import-alias "@/*" \
  --use-pnpm

cd admin

# Install operational UI and form validation dependencies
pnpm add @tanstack/react-table @tanstack/react-query \
  react-hook-form @hookform/resolvers zod lucide-react sonner clsx tailwind-merge

# Compile TypeScript types from canonical OpenAPI specification
pnpm dlx openapi-typescript ../contracts/openapi.yaml -o src/types/api.generated.ts

# Verify production build and lockfile generation
pnpm run build

cd ..
```

---

### Step 7: Git Governance Activation via Lefthook
Install and activate Git hooks to enforce atomic documentation commits, linting, and commit message formatting:
```bash
# Verify lefthook configuration exists
test -f lefthook.yml

# Register Git hooks
lefthook install

# Run manual pre-commit audit drill
lefthook run pre-commit
```

---

## 5. Phase 1 Verification Checklist & Sign-Off Criteria

Before closing Phase 1, confirm that every item below satisfies the requirements in `docs/quality/BOOTSTRAP_CHECKLIST.md`:

- [ ] **Infrastructure Health:** Docker Compose containers running with clean health checks (Postgres, Valkey, SeaweedFS, ClamAV, Traefik).
- [ ] **PostGIS Extension Verification:** `postgis`, `uuid-ossp`, and `pgcrypto` active on `dink_machinery_dev`.
- [ ] **Country Baseline:** `country_configs` seeded with valid `ET` and `KE` records containing integer tax basis points.
- [ ] **Backend Code Generation:** `internal/transport/http/generated/api.gen.go` generated without manual edits; `sqlc` models compiled into `internal/platform/database/db/`.
- [ ] **Backend Compilation:** `go build ./cmd/api` compiles cleanly; zero warnings on `golangci-lint run`.
- [ ] **Mobile Workspace:** Flutter 3.47.9 workspace passes `fvm flutter analyze` and `fvm flutter test`; `pubspec.lock` committed.
- [ ] **Admin Console:** Next.js 16 builds with zero TypeScript errors (`pnpm run build`); `src/types/api.generated.ts` synchronized with `openapi.yaml`.
- [ ] **Zero-PII Compliance:** Static log audit passes with zero unmasked PII log occurrences.
- [ ] **Governance Hooks:** Lefthook pre-commit and commit-msg hooks operational.

---

## 6. Handover Protocol to Phase 2 Feature Partitions

Upon completing all verification criteria:

1. **Update `docs/quality/BOOTSTRAP_CHECKLIST.md`:** Check all Phase 1 completion boxes.
2. **Update `CONTEXT.md`:** Document completion of Phase 1 bootstrap and note the transition to Phase 2.
3. **Commit Phase 1 Artifacts:**
   ```bash
   git add .
   git commit -m "[BOOTSTRAP] feat(platform): complete phase 1 subsystem initialization

   Bootstrapped Go 1.27.3 modular monolith, Flutter 3.47.9 client, Next.js 16 admin,
   PostgreSQL 18.6 + PostGIS baseline schema, Valkey 9.1.3, and Lefthook governance.

   Tests: All subsystem builds and test suites passing.
   Docs: CONTEXT.md, BOOTSTRAP.md, docs/quality/BOOTSTRAP_CHECKLIST.md
   Refs: #0001"
   ```
4. **Create First Feature Partition Branch:**
   ```bash
   git checkout -b feature/05-01-core-identity-auth
   ```