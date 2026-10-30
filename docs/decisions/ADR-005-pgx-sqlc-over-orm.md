# ADR-005: Type-Safe SQL Compilation (sqlc + pgx/v5) on PostgreSQL 18.4 + PostGIS 3.4 over ORMs

- **Status:** Accepted
- **Date:** 2026-10-26
- **Decision Owner:** Dink Machinery Architecture Review Board
- **Affected Partitions:** 05-02, 05-03, 05-04, 05-12, 10-data-model
- **Related Requirements:** SYS-01-BE, SRCH-01-BE, FL-07-BE, SEC-08-BE

---

## 1. Context & Problem Statement

Dink Machinery requires spatial radius discovery queries (`ST_DWithin` on WGS 84 spheroids), complex JSONB querying for dynamic equipment specs, and strictly typed double-entry ledger mutations. Traditional heavyweight Object-Relational Mappers (ORMs) introduce reflection overhead, obscure execution plans, suffer from N+1 query traps, and lack native support for advanced PostGIS geometry operators.

---

## 2. The 5-Point Qualification Evaluation

- [x] **Alternative Test:** Evaluated `sqlc` + `pgx/v5`, GORM, Ent Go, and raw `database/sql` queries.
- [x] **Surprise Test:** Eliminates ORM magic in favor of pure compile-time type-safe Go code generated directly from SQL files.
- [x] **Constraint Test:** PostGIS GiST spatial indexing and JSONB path operations require explicit SQL optimization and deterministic execution plans.
- [x] **Irreversibility Test:** Fundamental data access layer architecture across all backend persistence modules.
- [x] **Scope Test:** Governs all database migrations, queries, transactions, and data models across the entire platform.

---

## 3. Considered Options

### Option 1: `sqlc` + `jackc/pgx/v5` on PostgreSQL 18.4 + PostGIS 3.4
- **Summary:** Write plain SQL queries in `internal/platform/database/queries/*.sql`. `sqlc` compiles them into type-safe Go repository code using `pgx/v5` native binary protocol encoding.
- **Pros:**
  - Zero runtime reflection; maximum Go execution throughput.
  - Full access to advanced PostgreSQL features: PostGIS spatial operators, `jsonb_path_ops`, window functions, and CTEs.
  - Compile-time SQL validation prevents syntax errors from reaching production.
- **Cons:**
  - Dynamic query construction (e.g. multi-faceted search with 10 optional filters) requires structured query composition or specialized query templates.

### Option 2: Heavyweight ORM (GORM)
- **Summary:** Global ORM with automated struct-to-table mapping and hook lifecycles.
- **Pros:**
  - Rapid prototyping with automated CRUD operations.
- **Cons:**
  - Heavy reflection overhead; unpredictable SQL generation; awkward PostGIS syntax; high risk of N+1 query latency spikes.

### Option 3: Ent Go Entity Framework
- **Summary:** Graph-based schema definition in Go with code-generated client.
- **Pros:**
  - Type-safe query building and graph traversals.
- **Cons:**
  - Steep learning curve; rigid schema DSL; awkward custom SQL and spatial PostGIS integration.

---

## 4. Chosen Solution & Architectural Rationale

We adopt **Option 1: `sqlc` + `jackc/pgx/v5` on PostgreSQL 18.4 + PostGIS 3.4**.

1. **Native Performance & Binary Protocol:** `pgx/v5` communicates directly over PostgreSQL's binary protocol, eliminating string parsing overhead.
2. **Type Safety Without Reflection:** `sqlc` validates SQL queries against the schema DDL at build time and emits type-safe structs:
   ```sql
   -- name: SearchFleetsByRadius :many
   SELECT id, owner_id, title, hourly_rate_minor, currency,
          ST_Distance(location, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) AS distance_meters
   FROM fleets
   WHERE ST_DWithin(location, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, $3)
     AND status = 'ACTIVE'
     AND deleted_at IS NULL
   ORDER BY distance_meters ASC
   LIMIT $4;
   ```
3. **Explicit Connection Pooling:** Managed via `pgxpool.Pool` with strict health check monitoring and query cancellation contexts.

---

## 5. Rejected Alternatives & Technical Justification

- **Option 2 (GORM):** Rejected because ORM reflection overhead wastes CPU resources, hides real database execution plans, and provides poor native PostGIS integration.
- **Option 3 (Ent Go):** Rejected because writing complex financial double-entry triggers and spatial PostGIS joins is far more maintainable in standard SQL than in proprietary Go DSLs.

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - SQL queries are completely transparent and easily analyzed via `EXPLAIN (ANALYZE, BUFFERS)`.
  - Compile-time verification catches column typo bugs before code review.
  - Sub-15ms spatial queries and sub-millisecond simple lookups.
- **Negative / Technical Debt:**
  - Dynamic faceted search queries require well-crafted SQL templates or structured query builders (`squirrel` or dedicated query handlers).
- **Operational / Observability Impacts:**
  - All query execution times are tracked via OpenTelemetry pgx trace spans.

---

## 7. Verification Strategy & Invariant Tests

- **SQL Compilation CI Check:** CI executes `sqlc generate` and validates zero compiler errors or warnings.
- **Spatial Execution Benchmark:** `backend/internal/modules/machinery/infrastructure/spatial_benchmark_test.go` verifies that radius queries over 50,000 synthetic machinery listings complete in $< 20\text{ms}$.

---

## 8. Revisit Conditions

- Introduction of graph-oriented operational features where relational models and PostGIS become an impediment to query performance.
