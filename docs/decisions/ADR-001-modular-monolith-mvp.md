# ADR-001: Modular Monolith in Go 1.27.3 over Distributed Microservices

- **Status:** Accepted
- **Date:** 2026-10-26
- **Decision Owner:** Dink Machinery Architecture Review Board
- **Affected Partitions:** 05-05, 05-08, 08-booking-state-machine, 12-tech-stack
- **Related Requirements:** SYS-01-BE, SYS-02-BE, BK-17-BE

---

## 1. Context & Problem Statement

Dink Machinery coordinates heavy machinery rentals, logistics sub-bookings, operator verification, and financial ledgers across Ethiopia and Kenya. The engineering organization requires high development velocity, deterministic deployments, low network serialization overhead, and strict transaction atomicity without the cognitive and operational burdens of distributed microservices.

---

## 2. The 5-Point Qualification Evaluation

- [x] **Alternative Test:** Evaluated microservices over gRPC, unstructured single-package monolith, and modular monolith architecture.
- [x] **Surprise Test:** Prevents common industry bias toward microservices that introduces premature distributed systems complexity.
- [x] **Constraint Test:** Driven by hard requirement for atomic in-transaction double-entry ledger commits across bookings and logistics.
- [x] **Irreversibility Test:** High migration cost (>1 engineer-year to deconstruct or reconstruct module boundaries).
- [x] **Scope Test:** Governs all backend packages, runtime deployments, database transactional perimeters, and testing strategy.

---

## 3. Considered Options

### Option 1: Modular Monolith in Go 1.27.3 with Clean Architecture
- **Summary:** Single deployable static binary with strict internal package boundaries (`internal/modules/{auth, machinery, booking, payment, inspection, geo, admin}`). Packages communicate via exported interfaces or an in-memory typed event bus.
- **Pros:**
  - Single static binary deployment ($< 35\text{MB}$ RSS baseline).
  - Single database transaction (`pgx.Tx`) spans multi-module mutations (e.g. booking state transition + ledger write).
  - Sub-millisecond in-process function invocation with zero network failure modes.
- **Cons:**
  - Requires strict CI linter enforcement to prevent unauthorized cross-package imports.
  - Shared CPU and memory runtime across all functional modules.

### Option 2: Distributed Microservices over gRPC on Kubernetes
- **Summary:** Independent deployable services for Auth, Fleet, Booking, Payment, and Admin communicating via gRPC.
- **Pros:**
  - Independent deployment lifecycles and independent vertical scaling per service.
  - Hard physical process boundaries preventing direct cross-domain coupling.
- **Cons:**
  - Distributed transaction complexity (requires Saga orchestrator or 2PC for booking-ledger atomicity).
  - High operational infrastructure burden in local colocation facilities in Addis Ababa.
  - Network serialization overhead and complex distributed failure modes.

### Option 3: Traditional Unstructured Go Monolith
- **Summary:** Single Go codebase with flat packages or unstructured domain layering.
- **Pros:**
  - Fast initial prototyping with zero architectural scaffolding overhead.
- **Cons:**
  - Tangled circular dependencies and direct cross-domain SQL joins.
  - Rapid architectural degradation and unmaintainable monolithic entanglement.

---

## 4. Chosen Solution & Architectural Rationale

We adopt **Option 1: Modular Monolith in Go 1.27.3**. 

Domain logic is strictly organized into Clean Architecture feature modules under `internal/modules/{domain}`. Each module maintains its own `domain`, `application`, and `infrastructure` layers. Cross-domain interactions execute exclusively via exported Go interfaces or an in-memory typed domain event bus (`eventbus.Bus`). Transactional consistency across bookings, wallet commission debits, and ledger logs is guaranteed within a single `pgx.Tx` transaction boundary.

```
+-------------------------------------------------------------------------+
|                          Traefik 3.x Ingress Gateway                    |
+-------------------------------------------------------------------------+
                                     |
+-------------------------------------------------------------------------+
|                  Go 1.27.3 Modular Monolith (bin/api)                   |
|                                                                         |
|  [Auth Module]     [Machinery Module]     [Booking Module]              |
|  - Domain          - Domain               - Domain                      |
|  - Application     - Application          - Application                 |
|  - Infrastructure  - Infrastructure       - Infrastructure              |
|                                                                         |
|  [Payment Module]  [Inspection Module]    [Admin Module]                |
|  - Domain          - Domain               - Domain                      |
|  - Application     - Application          - Application                 |
|  - Infrastructure  - Infrastructure       - Infrastructure              |
|                                                                         |
|               In-Memory Typed Event Bus (eventbus.Bus)                  |
+-------------------------------------------------------------------------+
                                     |
              +----------------------+----------------------+
              |                                             |
+-----------------------------+               +---------------------------+
| PostgreSQL 18.4 + PostGIS   |               | Valkey 8.0 Cache & Mutex  |
+-----------------------------+               +---------------------------+
```

---

## 5. Rejected Alternatives & Technical Justification

- **Option 2 (Distributed Microservices):** Rejected due to the operational fragility of distributed transactions across emerging-market datacenter environments, high network serialization latency, and unnecessary DevOps cognitive load.
- **Option 3 (Traditional Unstructured Monolith):** Rejected because the absence of compiler-enforced package isolation inevitably causes spaghetti cross-queries and unmaintainable domain coupling.

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - Single binary artifact simplifies Docker packaging and bare-metal colocation deployment.
  - Zero distributed rollback choreography or partial network failures on multi-table mutations.
  - High unit test speed with in-memory SQLite/pgx mocks.
- **Negative / Technical Debt:**
  - A crash or memory leak in one module can terminate the entire application process.
  - Requires strict linting rules (`golangci-lint` with `depguard`) to forbid unauthorized package imports.
- **Operational / Observability Impacts:**
  - Structured module-scoped loggers and OpenTelemetry trace spans tag all operations with `module: [module_name]`.

---

## 7. Verification Strategy & Invariant Tests

- **Static Linter Rule:** `golangci-lint` enforces `depguard` rules ensuring `internal/modules/booking/domain` never imports `internal/modules/payment/infrastructure`.
- **Transaction Atomicity Unit Test:** `backend/internal/modules/booking/application/confirm_booking_test.go` verifies that if a ledger entry fails, the entire booking confirmation rolls back.
- **Binary Footprint Assertion:** CI build step asserts that the compiled static binary size is $\le 45\text{MB}$.

---

## 8. Revisit Conditions

- Backend engineering team scales beyond 35 concurrent contributors working on disjoint domains.
- A single domain module (e.g. high-frequency telematics GPS stream ingestion) demands $\ge 10\times$ horizontal resource scaling relative to the core transactional platform.
