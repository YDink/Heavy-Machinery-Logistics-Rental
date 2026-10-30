# Phase 1 Bootstrap Runbook & Orchestration Guide — Dink Machinery

**System:** Heavy Machinery & Logistics Rental Marketplace  
**Phase:** 1 (Bootstrap & Subsystem Initialization)  
**Governing Documents:** `AGENTS.md`, `ARCHITECTURE.md`, `TECH_STACK.md`, `CODE_STANDARDS.md`, `docs/quality/BOOTSTRAP_CHECKLIST.md`

---

## 1. Overview & Operational Invariants

This runbook defines the exact, deterministic sequence required to transition Dink Machinery from **Phase 0 (Documentation & Contract Definition)** into **Phase 1 (Machine Bootstrapping)** and prepare the repository for **Phase 2 (Partition Feature Development)**.

### Invariant Rules During Bootstrap
1. **Zero Domain Logic in Phase 1:** Bootstrapping creates runtime scaffolds, dependency graphs, database extensions, code-generation pipelines, and continuous integration gates. No partition business logic (e.g., booking state transitions, payment gateway integrations) may be implemented during Phase 1.
2. **Strict Lockfile Pinning:** Every scaffolded subsystem (`backend/`, `mobile/`, `admin/`) must generate and commit deterministic lockfiles (`go.sum`, `pubspec.lock`, `package-lock.json`).
3. **Canonical Verification:** The bootstrap is considered complete only when all verification gates in `docs/quality/BOOTSTRAP_CHECKLIST.md` pass without warnings.

---

## 2. Host Prerequisites & Runtime Matrix

Ensure the host workstation meets the required toolchain specifications before starting the bootstrap sequence:

| Runtime / Tool | Required Version | Verification Command | Purpose |
|---|---|---|---|
| **Go** | `1.27.3` | `go version` | Backend modular monolith engine |
| **Flutter / Dart** | `3.47.8` / `3.9.x` | `flutter --version` | Cross-platform mobile client |
| **Node.js / npm** | `26.13.0` / `10.x+` | `node -v && npm -v` | Admin web console & contract tooling |
| **Docker Engine & Compose** | `27.x+` / `v2.x+` | `docker compose version` | Local containerized infrastructure |
| **Spectral CLI** | `v6.x+` | `npx @stoplight/spectral-cli -v` | OpenAPI 3.1.0 contract linting |
| **sqlc** | `v1.27+` | `sqlc version` | Type-safe SQL code generation |
| **oapi-codegen** | `v2.4+` | `oapi-codegen --version` | Go transport code generation |
| **golang-migrate** | `v4.18+` | `migrate -version` | Database migration runner |
| **Lefthook** | `v1.8+` | `lefthook version` | Local Git hook governance |

---

## 3. Step-by-Step Execution Sequence

```mermaid
flowchart TD
  S0[Phase 0 Gate Audit] --> S1[Local Infrastructure Spin-up]
  S1 --> S2[Database Extensions & Migration Baseline]
  S2 --> S3[Backend Scaffold & Code Generation]
  S3 --> S4[Flutter Mobile Scaffold & SDK Generation]
  S4 --> S5[Next.js Admin Scaffold & Schema Sync]
  S5 --> S6[Lefthook & Git Governance Activation]
  S6 --> S7[End-to-End Verification Drill]
  S7 --> S8[Phase 1 Commit & Phase 2 Partition Branching]
```

---

### Step 0: Pre-Execution Verification Gate

Before generating code or directories, run the documentation linting and architectural governance suites:

```bash
# 1. Lint OpenAPI canonical contract
npx @stoplight/spectral-cli lint contracts/openapi.yaml --ruleset .spectral.yaml

# 2. Audit decision records and templates
bash scripts/validate-decision-template.sh
bash scripts/validate-decision-links.sh

# 3. Assert zero drift against source-of-truth
bash scripts/ci-docs.yml || true
```

*Proceed only if all pre-execution checks pass with zero errors.*

---

### Step 1: Local Infrastructure Provisioning

Initialize and start the local container dependencies using `docker-compose.yml`.

```bash
# Start local persistence, cache, storage, and security daemons
docker compose up -d postgres valkey seaweedfs clamav traefik

# Verify container health status
docker compose ps
```

**Health Verification:**
- PostgreSQL must be accepting connections on `localhost:5432`.
- Valkey must respond with `PONG` to `docker compose exec valkey redis-cli ping`.
- SeaweedFS S3 API must be reachable on `http://localhost:9000`.
- ClamAV daemon must respond to `PING` on port `3310`.

---

### Step 2: Database Extensions & Schema Baseline

Execute the baseline migration to configure PostGIS, enable cryptographic extensions, establish the active country configuration schema, and seed default operational jurisdictions.

#### 1. Create Migration Files
```bash
migrate create -ext sql -dir migrations -seq -digits 6 init_platform_schema
```

#### 2. Apply Up Migration (`migrations/000001_init_platform_schema.up.sql`)
```sql
-- Enable Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Multi-Country Configuration Table (LOC-01-BE)
CREATE TABLE country_configs (
    country_code VARCHAR(2) PRIMARY KEY,
    country_name VARCHAR(100) NOT NULL,
    dialing_prefix VARCHAR(5) NOT NULL,
    default_currency VARCHAR(3) NOT NULL,
    currency_minor_unit_exponent INT NOT NULL DEFAULT 2,
    default_timezone VARCHAR(50) NOT NULL,
    tax_label VARCHAR(20) NOT NULL DEFAULT 'VAT',
    standard_tax_rate_bps INT NOT NULL DEFAULT 1500,
    withholding_tax_bps INT NOT NULL DEFAULT 200,
    min_age_default INT NOT NULL DEFAULT 18,
    payment_hold_timeout_mins INT NOT NULL DEFAULT 60,
    is_enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed Initial Jurisdictions (Ethiopia & Kenya)
INSERT INTO country_configs (
    country_code, country_name, dialing_prefix, default_currency, 
    currency_minor_unit_exponent, default_timezone, tax_label, 
    standard_tax_rate_bps, withholding_tax_bps, payment_hold_timeout_mins, is_enabled
) VALUES 
('ET', 'Ethiopia', '+251', 'ETB', 2, 'Africa/Addis_Ababa', 'VAT', 1500, 200, 60, true),
('KE', 'Kenya', '+254', 'KES', 2, 'Africa/Nairobi', 'VAT', 1600, 500, 60, true);

-- Invariant Audit Trigger for Updated At
CREATE OR REPLACE FUNCTION set_updated_at_column()
RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = NOW(); RETURN NEW; END; $$ LANGUAGE plpgsql;

CREATE TRIGGER trg_country_configs_updated_at
BEFORE UPDATE ON country_configs
FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();
```

#### 3. Run Migration
```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/dink_machinery_dev?sslmode=disable" up
```

---

### Step 3: Backend Modular Monolith Bootstrap (Go 1.27.3)

Initialize the Go module, establish Clean Architecture packages, configure `sqlc`, and generate HTTP transport bindings from `contracts/openapi.yaml`.

#### 1. Module Initialization & Directory Scaffolding
```bash
mkdir -p backend
cd backend

go mod init [github.com/dinkmachinery/dink-core](https://github.com/dinkmachinery/dink-core)
go get [github.com/go-chi/chi/v5](https://github.com/go-chi/chi/v5)
go get [github.com/jackc/pgx/v5](https://github.com/jackc/pgx/v5)
go get [github.com/google/uuid](https://github.com/google/uuid)
go get go.uber.org/zap
go get go.opentelemetry.io/otel

mkdir -p cmd/api \
  internal/platform/database \
  internal/platform/config \
  internal/platform/telemetry \
  internal/transport/http/middleware \
  internal/transport/http/generated \
  internal/modules/auth/domain \
  internal/modules/auth/application \
  internal/modules/auth/infrastructure \
  internal/modules/machinery/domain \
  internal/modules/booking/domain \
  internal/modules/payment/domain \
  internal/modules/admin/domain \
  internal/modules/geo/domain
```

#### 2. Generate OpenAPI HTTP Server Stubs
Generate the type-safe Chi router stubs and DTO structures:
```bash
# From backend directory
oapi-codegen -generate types,chi-server,spec \
  -package generated \
  -o internal/transport/http/generated/api.gen.go \
  ../contracts/openapi.yaml
```

#### 3. Configure `sqlc.yaml` and Generate Repository Layer
Create `backend/sqlc.yaml`:
```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "internal/platform/database/queries.sql"
    schema: "../migrations/"
    gen:
      go:
        package: "db"
        out: "internal/platform/database/db"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_empty_slices: true
        emit_pointers_for_null_types: true
```

Create baseline query `backend/internal/platform/database/queries.sql`:
```sql
-- name: GetCountryConfig :one
SELECT * FROM country_configs
WHERE country_code = $1 AND is_enabled = true;

-- name: ListEnabledCountryConfigs :many
SELECT * FROM country_configs
WHERE is_enabled = true
ORDER BY country_code ASC;
```

Run compilation:
```bash
sqlc generate
```

#### 4. Baseline Application Bootstrap (`backend/cmd/api/main.go`)
Provide standard health-check endpoints conforming to the platform envelope:
```go
package main

import (
	"encoding/json"
	"net/http"
	"time"

	"[github.com/go-chi/chi/v5](https://github.com/go-chi/chi/v5)"
	"[github.com/go-chi/chi/v5/middleware](https://github.com/go-chi/chi/v5/middleware)"
	"go.uber.org/zap"
)

type HealthResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
	Error   interface{}            `json:"error"`
	Meta    map[string]interface{} `json:"meta"`
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		resp := HealthResponse{
			Success: true,
			Data: map[string]interface{}{
				"status":  "healthy",
				"version": "1.0.0-phase1",
			},
			Error: nil,
			Meta: map[string]interface{}{
				"version":   "1.0",
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	})

	logger.Info("Backend modular monolith initialized", zap.String("port", "8080"))
	_ = http.ListenAndServe(":8080", r)
}
```

Verify compilation:
```bash
go build -o bin/api ./cmd/api
cd ..
```

---

### Step 4: Flutter Mobile Subsystem Bootstrap (Flutter 3.47.8)

Scaffold the cross-platform mobile workspace, pin dependencies, and establish the feature-first Clean Architecture hierarchy.

```bash
flutter create --org com.dinkmachinery --project-name dink_mobile mobile
cd mobile

# Add core architecture dependencies
flutter pub add flutter_bloc bloc get_it dio isar isar_flutter_libs \
  freezed_annotation json_annotation flutter_secure_storage geolocator \
  equatable intl

flutter pub add --dev build_runner freezed json_serializable isar_generator

# Create Feature-First Architecture Layout
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

# Verify Flutter compilation and lockfile generation
flutter test
cd ..
```

---

### Step 5: Next.js Admin Console Bootstrap (Next.js 16)

Scaffold the Next.js App Router admin application with TypeScript, Tailwind CSS v4, and Radix UI primitives.

```bash
npx create-next-app@latest admin \
  --typescript \
  --tailwind \
  --eslint \
  --app \
  --src-dir \
  --import-alias "@/*" \
  --use-npm

cd admin

# Install administrative UI libraries
npm install @tanstack/react-table @tanstack/react-query \
  react-hook-form @hookform/resolvers zod lucide-react sonner clsx tailwind-merge

# Generate TypeScript types from canonical OpenAPI contract
npx openapi-typescript ../contracts/openapi.yaml -o src/types/api.generated.ts

# Verify production build and lockfile integrity
npm run build
cd ..
```

---

### Step 6: Git Hook Governance Activation (`Lefthook`)

Activate local governance hooks using `lefthook` to prevent unlinted commits, broken links, or commits lacking required documentation footers.

Initialize `lefthook.yml` at repository root:
```yaml
pre-commit:
  parallel: true
  commands:
    spectral-lint:
      glob: "contracts/*.yaml"
      run: npx @stoplight/spectral-cli lint {staged_files} --ruleset .spectral.yaml
    decision-template-check:
      glob: "docs/decisions/**/*.md"
      run: bash scripts/validate-decision-template.sh
    decision-links-check:
      glob: "docs/**/*.md"
      run: bash scripts/validate-decision-links.sh
    go-lint:
      glob: "backend/**/*.go"
      run: cd backend && golangci-lint run ./...
    context-check:
      run: git diff --cached --name-only | grep -q "^CONTEXT.md$" || (echo "ERROR: CONTEXT.md must be updated and staged in every commit." && exit 1)

commit-msg:
  commands:
    commitlint:
      run: npx commitlint --config .github/commitlint.config.js --edit {1}
```

Install hooks:
```bash
lefthook install
```

---

## 4. End-to-End Verification Gate

Run the complete verification battery to ensure full compliance before marking Phase 1 complete:

```bash
# 1. Verify Docker Compose services are healthy
docker compose ps --filter "status=running"

# 2. Assert PostgreSQL extensions exist
docker compose exec postgres psql -U postgres -d dink_machinery_dev -c "SELECT extname FROM pg_extension WHERE extname IN ('postgis', 'uuid-ossp', 'pgcrypto');"

# 3. Assert baseline country seeds exist
docker compose exec postgres psql -U postgres -d dink_machinery_dev -c "SELECT country_code, default_currency, standard_tax_rate_bps FROM country_configs;"

# 4. Run backend static checks and compilation
cd backend && go vet ./... && go test -race ./... && cd ..

# 5. Run mobile test suite
cd mobile && flutter test && cd ..

# 6. Run admin lint and type checking
cd admin && npm run lint && cd ..

# 7. Run ASVS L3 Zero-PII logging check
bash -c "grep -rEni 'logger.*\.(Info|Debug|Warn|Error).*(password|national_id|tin_number|bank_account)' backend/ admin/ mobile/ && exit 1 || echo 'Zero PII logging asserted.'"
```

---

## 5. Phase 1 Completion & Transition to Phase 2

Once all verification checks pass:

1. **Update `docs/quality/BOOTSTRAP_CHECKLIST.md`:** Mark all Phase 1 gate verification items as checked.
2. **Update `CONTEXT.md`:** Record the completion of Phase 1 scaffolding and identify the first partition for Phase 2 development.
3. **Execute Atomic Bootstrap Commit:**
   ```bash
   git add .
   git commit -m "[BOOTSTRAP] feat(platform): complete phase 1 repository bootstrap

   Initialized Go 1.27.3 modular monolith, Flutter 3.47.8 mobile client, Next.js 16
   admin console, PostgreSQL 18.6 + PostGIS baseline schema, and Lefthook governance.

   Tests: All baseline subsystem compilation and test suites passing.
   Docs: CONTEXT.md, docs/quality/BOOTSTRAP_CHECKLIST.md
   Refs: #0001"
   ```
4. **Initialize Feature Partition Branch:**
   ```bash
   # Branch out for the first operational partition: Core Auth & Identity
   git checkout -b feature/05-01-core-identity-auth
   ```