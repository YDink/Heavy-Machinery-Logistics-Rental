# ADR-003: Clean Architecture Feature Modules over Monolithic Layering

- **Status:** Accepted
- **Date:** 2026-10-26
- **Decision Owner:** Dink Machinery Architecture Review Board
- **Affected Partitions:** 05-02, 05-03, 05-22, 10-data-model
- **Related Requirements:** SYS-01-BE, FL-01-BE, SD-01-BE

---

## 1. Context & Problem Statement

As the backend codebase expands to support complex domain capabilities (dynamic machinery specifications, attachment bundling, booking state machines, and soft-delete governance), layering must prevent database models and SQL queries from leaking into presentation or core business rules. Traditional horizontal layering (controllers, services, repositories) leads to bloated god-classes and tight coupling.

---

## 2. The 5-Point Qualification Evaluation

- [x] **Alternative Test:** Evaluated vertical feature slicing with Clean Architecture, traditional horizontal layering, and Active Record pattern.
- [x] **Surprise Test:** Enforces that domain entities have zero external package dependencies (no SQL, no JSON tags, no HTTP libraries).
- [x] **Constraint Test:** Ensures business logic is testable in total isolation with sub-second unit test execution.
- [x] **Irreversibility Test:** Refactoring leaky domain layers in production requires modifying hundreds of database interactions.
- [x] **Scope Test:** Governs the internal directory layout, package dependency graph, and testing strategy across all backend modules.

---

## 3. Considered Options

### Option 1: Clean Architecture Feature Modules (Vertical Slicing)
- **Summary:** Modules are encapsulated under `internal/modules/<feature>/` with inward dependency flow: `domain` (pure entities & interfaces) $\leftarrow$ `application` (use cases / commands / queries) $\leftarrow$ `infrastructure` (sqlc, Valkey, S3) / `presentation` (HTTP handlers).
- **Pros:**
  - Pure business domain entities with zero external dependencies.
  - Pluggable infrastructure: database or cache drivers can be upgraded without touching domain logic.
  - High unit test speed with mock repository implementations.
- **Cons:**
  - Requires explicit mapper functions to convert between persistence entities and domain entities.

### Option 2: Traditional Horizontal Layering (Global Controllers, Services, Repositories)
- **Summary:** Global packages (`controllers/`, `services/`, `models/`, `repositories/`) containing logic for all domains.
- **Pros:**
  - Familiar to developers coming from legacy MVC frameworks.
- **Cons:**
  - High risk of god-classes (e.g. `BookingService` with thousands of lines).
  - Lack of clear domain boundaries and accidental circular imports.

### Option 3: Active Record / ORM Domain Entities
- **Summary:** Database entities contain persistence methods (e.g. `fleet.Save()`, `booking.Update()`).
- **Pros:**
  - Rapid initial prototyping with minimal boilerplate.
- **Cons:**
  - Tight coupling between domain logic and SQL database schemas.
  - Impossibility of testing domain logic without a live database.

---

## 4. Chosen Solution & Architectural Rationale

We adopt **Option 1: Clean Architecture Feature Modules (Vertical Slicing)**.

Every backend capability resides in a feature-scoped module under `internal/modules/{domain}`:
```
internal/modules/machinery/
├── domain/               # Pure business entities, value objects, repository interfaces
│   ├── fleet.go
│   ├── repository.go     # Interface definition
│   └── errors.go
├── application/          # Use cases, command handlers, query handlers
│   ├── create_fleet.go
│   └── get_fleet.go
└── infrastructure/       # sqlc query implementations, S3 adapters, Valkey caching
    ├── postgres_repo.go  # Implements domain.FleetRepository
    └── mappers.go        # Converts sqlc DB models to domain entities
```

Inward dependency rule is strictly maintained: `domain` never imports `application`, `infrastructure`, or `transport`.

---

## 5. Rejected Alternatives & Technical Justification

- **Option 2 (Traditional Horizontal Layering):** Rejected due to the inevitable evolution of sprawling service files that obscure bounded context boundaries and invite spaghetti dependencies.
- **Option 3 (Active Record):** Rejected because coupling domain entities to database tables violates separation of concerns and impedes unit test determinism.

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - Complete isolation of core marketplace rules from database changes.
  - Unit tests run in milliseconds without spinning up PostgreSQL containers.
  - Easy onboarding: an engineer modifying machinery listing logic touches only `internal/modules/machinery/`.
- **Negative / Technical Debt:**
  - Increased boilerplate for entity-to-DTO mapping.
- **Operational / Observability Impacts:**
  - Trace spans and error codes align cleanly with use case boundaries.

---

## 7. Verification Strategy & Invariant Tests

- **Architecture Boundary Test:** `backend/internal/platform/architecture_test.go` uses Go AST scanning to verify that no file under `internal/modules/*/domain/` imports `database/sql`, `jackc/pgx`, or `net/http`.
- **Repository Interface Compliance:** Each infrastructure adapter compiles against its corresponding domain interface (`var _ domain.FleetRepository = (*PostgresFleetRepository)(nil)`).

---

## 8. Revisit Conditions

- Codebase expands beyond 50 distinct domain modules requiring separate standalone service extraction.
