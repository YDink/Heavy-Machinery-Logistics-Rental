# ADR-016: Modular OpenAPI Contract Architecture & Per-Endpoint Specification

- **Status:** Accepted
- **Date:** 2026-10-30
- **Decision Owner:** Dink Machinery Architecture Review Board
- **Affected Partitions:** 09-non-functional, 12-tech-stack
- **Related Requirements:** SYS-02-BE, INF-02-OPS, CONTRACT-01-API

---

## 1. Context & Problem Statement

Dink Machinery enforces a contract-first engineering lifecycle where `contracts/openapi.yaml` serves as the single source of truth for the Go 1.27.3 modular monolith HTTP transport router (`oapi-codegen`), the Next.js 16 admin console TypeScript types (`openapi-typescript`), and the Flutter 3.47.9 Dio HTTP client models.

As the platform scales across its 24 SRS feature partitions (encompassing over 80 discrete endpoints across identity, fleet catalog, geospatial search, 20-state reservation lifecycles, double-entry financial ledger operations, dispute arbitration, and dynamic configuration engines), maintaining a monolithic `openapi.yaml` exceeding 1,600+ lines introduces severe engineering bottlenecks:
1. **Merge Conflicts:** Multiple feature teams and autonomous AI agents modifying the monolithic contract file concurrently on partition branches (`feature/<partition-id>`) suffer frequent, high-risk Git merge conflicts.
2. **Review & Audit Friction:** Code reviews for isolated endpoints require wading through thousands of lines of unrelated schemas, obscuring domain-level boundary modifications.
3. **Partition Traceability:** The monolithic contract fails to mirror the Clean Architecture modularity and SRS partition boundaries enforced across the Go backend and Flutter mobile applications.

The platform requires a modular, per-endpoint contract architecture that preserves 100% strict type generation, pre-commit linting, and RFC 9457 compliance while isolating endpoints into discrete, dedicated specification files.

---

## 2. The 5-Point Qualification Evaluation

- [x] **Alternative Test:** Evaluated monolithic single-file OpenAPI vs. domain-grouped multi-file OpenAPI vs. granular per-endpoint modular OpenAPI paired with deterministic Redocly bundling.
- [x] **Surprise Test:** Clarifies why endpoints are individually isolated into separate YAML files rather than grouped into huge monolithic domain files, and documents the necessity of the Redocly distribution bundling step for code generation.
- [x] **Constraint Test:** Driven by `oapi-codegen` and Stoplight Spectral parser constraints regarding external cross-file schema references and OpenAPI 3.1 `path-item` evaluation.
- [x] **Irreversibility Test:** Modifying contract directory organization and generator build steps impacts all CI validation gates, backend HTTP routers, admin TypeScript bindings, and mobile API clients.
- [x] **Scope Test:** Governs the entire API interface across all backend modules, admin portals, mobile clients, and QA linting gates.

---

## 3. Considered Options

### Option 1: Granular Per-Endpoint Modular Contracts with Redocly Bundling (Chosen)
- **Summary:** Every endpoint resides in a dedicated file (`contracts/endpoints/<domain>/<action>.yaml`) defining its Path Item, operations, local request/response payloads, and references to shared primitives (`contracts/components/`). The root `contracts/openapi.yaml` maps paths to endpoint files via relative `$ref`. A high-speed build step (`@redocly/cli bundle`, executing in $< 25\text{ms}$) resolves all references into a deterministic distribution artifact (`contracts/openapi.bundled.yaml`), which feeds `oapi-codegen`, `openapi-typescript`, and Stoplight Spectral.
- **Pros:**
  - Complete elimination of Git merge conflicts across concurrent partition feature branches.
  - 1:1 mapping between SRS requirement partitions and contract files.
  - Granular code reviews: inspecting a feature PR touches only the specific endpoint file.
  - Zero disruption to downstream toolchains (`oapi-codegen` and `openapi-typescript` consume the bundled file with zero external reference errors).
- **Cons:**
  - Requires maintaining the deterministic bundling build step in pre-commit hooks and CI.

### Option 2: Domain-Grouped Contract Files (e.g., `auth.yaml`, `bookings.yaml`, `fleets.yaml`)
- **Summary:** Group endpoints into 6–8 domain-level files rather than individual endpoint files.
- **Pros:**
  - Fewer total files than granular per-endpoint isolation.
- **Cons:**
  - Still prone to merge conflicts when multiple engineers work on different operations within the same domain (e.g., `create-booking` vs `cancel-booking` vs `upload-inspection`).
  - Fails to achieve atomic requirement-level isolation.

### Option 3: Status Quo Monolithic `openapi.yaml`
- **Summary:** Retain the single monolithic `contracts/openapi.yaml` containing all paths, operations, parameters, and schemas.
- **Pros:**
  - Zero bundling step required; all tools read a single file directly.
- **Cons:**
  - Unmaintainable at scale ($> 5,000$ lines as the full SRS is specified).
  - High probability of merge collisions and regression errors during Phase 2 parallel development.

---

## 4. Chosen Solution & Architectural Rationale

The platform adopts **Option 1: Granular Per-Endpoint Modular Contracts with Redocly Bundling**.

### Directory Structure Standard
```text
contracts/
├── openapi.yaml                          # Root entrypoint: API metadata, servers, tags, security, and path $refs
├── openapi.bundled.yaml                  # Deterministic distribution artifact (consumed by codegen & linters)
├── components/                           # Shared cross-cutting components
│   ├── parameters/                       # Reusable headers and query parameters
│   │   ├── country-code-header.yaml      # X-Country-Code header
│   │   ├── idempotency-key-header.yaml   # Idempotency-Key header
│   │   ├── cursor-query-param.yaml       # Keyset pagination cursor
│   │   └── limit-query-param.yaml        # Page size limit
│   ├── responses/                        # Standard RFC 9457 error responses
│   │   ├── bad-request.yaml              # HTTP 400
│   │   ├── unauthorized.yaml             # HTTP 401
│   │   ├── forbidden.yaml                # HTTP 403
│   │   ├── not-found.yaml                # HTTP 404
│   │   ├── conflict.yaml                 # HTTP 409
│   │   ├── too-many-requests.yaml        # HTTP 429
│   │   └── internal-server-error.yaml    # HTTP 500
│   └── schemas/                          # Shared domain entities and value objects
│       ├── money-amount.yaml             # Integer minor units + ISO 4217 currency
│       ├── response-meta.yaml            # Standard requestId, version, timestamp
│       ├── pagination-meta.yaml          # Keyset cursor pagination envelope
│       ├── problem-details.yaml          # RFC 9457 Problem Details model
│       └── fleet-equipment.yaml          # Core Fleet domain entity
└── endpoints/                            # 1:1 Isolated per-endpoint specifications
    ├── auth/                             # Identity, OTP, session lifecycle
    ├── machinery/                        # Fleet catalog, categories, specifications
    ├── geo/                              # PostGIS geospatial search
    ├── bookings/                         # 20-state reservation lifecycle & inspections
    ├── payments/                         # Double-entry ledger & payment initiation
    ├── wallets/                          # Owner wallet & transactions
    └── admin/                            # Dynamic UI config & country registries
```

### Toolchain Integration Pipeline
1. **Source Authoring:** Developers modify only the specific `contracts/endpoints/<domain>/<action>.yaml` and any shared schemas in `contracts/components/`.
2. **Bundling (`@redocly/cli bundle`):** Compiles the modular tree into `contracts/openapi.bundled.yaml` in $< 25\text{ms}$.
3. **Go Transport Router (`oapi-codegen`):** Compiles `contracts/openapi.bundled.yaml` into `backend/internal/transport/http/generated/api.gen.go`.
4. **Next.js TypeScript Types (`openapi-typescript`):** Compiles `contracts/openapi.bundled.yaml` into `admin/src/types/api.generated.ts`.
5. **Quality Gate Linting (Spectral):** Validates `contracts/openapi.bundled.yaml` against `.spectral.yaml` in pre-commit hooks (`lefthook.yml`) and CI (`ci-docs.yml`).

---

## 5. Rejected Alternatives & Technical Justification

- **Option 2 (Domain-Grouped Files):** Rejected because grouping 15+ booking endpoints into a single `bookings.yaml` reintroduces merge conflicts during parallel sprint execution and prevents atomic requirement mapping.
- **Option 3 (Monolithic Contract):** Rejected due to unsustainable merge conflict risk and human cognitive fatigue navigating 1,600+ lines of YAML.

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - Perfect isolation: parallel branches touching different endpoints never collide in Git.
  - Modular clean reviews: PR diffs are constrained strictly to the modified endpoint specification.
  - Zero tooling breakage: downstream generators receive a clean, fully dereferenced distribution spec.
- **Negative / Technical Debt:**
  - Adds a bundling step (`redocly bundle`) to the local build and pre-commit workflow.
  - Requires developers to maintain relative reference paths (`$ref: '../../components/...'`).
- **Operational / Observability Impacts:**
  - Bundled file `contracts/openapi.bundled.yaml` is checked into version control to guarantee deterministic reproducibility on machines without Node.js installed.

---

## 7. Verification Strategy & Invariant Tests

- **Test Suite:** Stoplight Spectral Linter & Redocly Validator
- **Commands:**
  - `pnpm --dir admin dlx --package=@redocly/cli redocly bundle contracts/openapi.yaml -o contracts/openapi.bundled.yaml`
  - `pnpm dlx @stoplight/spectral-cli lint contracts/openapi.bundled.yaml --ruleset .spectral.yaml`
- **Assertion:** Zero errors on Universal Success Envelope, RFC 9457 Problem Details, Minor-Unit Currency, and Idempotency-Key headers.

---

## 8. Revisit Conditions

- Re-evaluate if OpenAPI 4.0 (Moonwalk) introduces native multi-file namespaces that eliminate the need for pre-codegen bundling.
