# Dink Machinery

Heavy machinery and logistics rental aggregator platform launching in Ethiopia (`country_code = 'ET'`) with multi-country parameterization[cite: 35].

---

## 1. Project Overview & Current Phase

The platform architecture spans three specialized components backed by a single hand-authored API contract[cite: 37]:
1. **`backend/`**: Feature-Based Modular Monolith in Go (Clean Architecture, CQRS, `pgx/sqlc`)[cite: 32, 37].
2. **`mobile/`**: Single Flutter codebase supporting Renter, Owner, and Operator roles (BLoC, `isar_community`)[cite: 32, 37].
3. **`admin/`**: Next.js 16 (App Router) database-driven dynamic UI console[cite: 32, 37].

> **CURRENT STATUS: PHASE 0 (DOCS-ONLY)**
> No application code directories (`backend/`, `mobile/`, `admin/`) exist in the repository during Phase 0. Development is strictly gated by `docs/quality/BOOTSTRAP_CHECKLIST.md`. Phase 1 project scaffolding is executed only when the Phase 0 checklist reaches 100% verification. See `docs/BOOTSTRAP.md` for bootstrap execution procedures.

---

## 2. The Five Engineering Disciplines

1. **Tool Portability:** `AGENTS.md` is the single tool-agnostic AI context file. `CLAUDE.md` and `GEMINI.md` are one-line wrappers (`@AGENTS.md`).
2. **Git & Commit Discipline:** Work proceeds on one branch per requirement partition (`feature/<partition-id>`)[cite: 37]. Every commit enforces the **atomic doc-commit rule**: code changes + tests + `CONTEXT.md` staged simultaneously.
3. **Three-Tier Decision Documentation:** Decisions are qualified via the 5-point test and documented as L1 Architecture Decision Records (`docs/decisions/`), L2 Partition Decisions (`docs/decisions/partitions/`), or L3 Inline comments.
4. **QA as an Operating System:** Verification is governed by machine-checkable quality gates (G0 to G4), a non-decreasing coverage ratchet (80% minimum, 95% financial branch coverage)[cite: 24], and zero-PII logging assertions.
5. **Contract-First & Invariant Standardization:** All mutations follow the unified response envelope, cursor pagination, RFC 9457 Problem Details[cite: 24], string/integer minor-unit money arithmetic, and strict in-transaction ledger atomicity (`BK-17-BE`)[cite: 35].

---

## 3. Directory Navigation

```
├── contracts/          # Hand-authored OpenAPI 3.0 contract (Single Source of Truth)[cite: 37]
├── docs/
│   ├── source/         # Canonical monolithic master documents (SRS v7.1, ADD v4.2, Flows, UI Brief)[cite: 36, 37]
│   ├── shared/         # Split requirements, architecture, flows, and design tokens[cite: 37]
│   ├── quality/        # QA Strategy, Quality Gates, Security QA, and Bootstrap Checklist
│   ├── decisions/      # L1 ADRs and L2 partition decision documents[cite: 37]
│   └── components/     # Pre-bootstrap design specifications for backend, mobile, and admin
```

---

## 4. Local Infrastructure Startup

Phase 0 documentation linting and local schema verification require local Docker services[cite: 37]:

```bash
# Start local PostgreSQL, Valkey, MinIO, and ClamAV
docker compose up -d

# Verify container health
docker compose ps
```