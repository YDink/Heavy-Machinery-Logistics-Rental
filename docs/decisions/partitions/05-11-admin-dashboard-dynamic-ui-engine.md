# Partition Decisions: 05-11 — Admin Dashboard, Dynamic UI Engine & RBAC Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-11`.

- **Governing System ADRs:** `ADR-002`, `ADR-008`
- **Subsystem Scope:** backend, admin
- **Requirement Targets:** AD-01-BE to AD-14-BE, AD-01-WEB to AD-14-WEB

---

## D-05-11-001: Dynamic Schema Compilation from /api/v1/admin/ui-config with Redis Caching

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [AD-06-BE]
- **Subsystem:** [admin, backend]

### 1. Problem Statement & Context
Back-office administrative consoles for machinery operations, dispute arbitrations, and finance require frequent updates to table columns, filters, and operational buttons. Deploying new frontend builds for simple column or permission changes slows down operational agility.

### 2. Considered Alternatives
1. **Option A:** Server-driven UI configuration engine (/api/v1/admin/ui-config) compiling navigation, data grid columns, and action buttons dynamically from database metadata
   - *Pros:* Zero web redeployments to alter tables or permissions; server-side RBAC scoping; dynamic TanStack table synthesis.
   - *Cons:* Requires careful schema design and caching in Valkey.
2. **Option B:** Hardcoded React table schemas and action handlers in Next.js code
   - *Pros:* Fast initial frontend coding.
   - *Cons:* Every column addition, rename, or permission tweak requires a full frontend CI/CD deployment cycle.
3. **Option C:** Third-party SaaS admin builder (e.g. Retool, Forest Admin)
   - *Pros:* No custom admin console code.
   - *Cons:* Severe data sovereignty violations (Proclamation No. 1321/2024); recurring foreign currency subscription costs.

### 3. Chosen Solution & Technical Rationale
Option A: Server-driven dynamic UI engine. The Next.js 16 back-office fetches UI schema definitions from /api/v1/admin/ui-config, cached in Valkey. Layouts, tables, and buttons compile dynamically matching the authenticated administrator role permissions.

### 4. Rejected Alternatives
- Option B: Rejected because hardcoded admin tables create deployment bottlenecks for operational teams.
- Option C: Rejected due to statutory data residency prohibitions and foreign currency payment bans.

### 5. Consequences & Invariants
- **Invariants Maintained:** An administrator only receives UI configurations and action buttons for which their role holds verified permissions.
- **Operational Trade-offs:** Dynamic column rendering requires robust fallback formatters for dates, currencies, and status badges.

### 6. Verification
- **Test File:** `admin/tests/dynamic_ui_engine_test.tsx`
- **Assertion:** TestDynamicUI_RendersColumnsFromSchema asserts TanStack table columns match backend schema definition dynamically.

### 7. Revisit Trigger
- Revisit if administrative requirements demand complex custom canvas workflows beyond tabular data grids.

---

## D-05-11-002: Granular RBAC Matrix Evaluation with Immutable Audit Journaling

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [AD-07-BE]
- **Subsystem:** [backend, database]

### 1. Problem Statement & Context
Administrative personnel hold differing authority levels (Super Admin, Operations Manager, Finance Controller, Dispute Arbitrator, Support Agent). Unauthorized state overrides or ledger adjustments pose severe financial fraud risks.

### 2. Considered Alternatives
1. **Option A:** Granular permissions matrix (role_permissions) evaluated in transport middleware paired with automated before/after JSON diff audit logging
   - *Pros:* Strict least-privilege security; tamper-evident audit journal of all administrative mutations; ASVS Level 3 compliance.
   - *Cons:* Small latency overhead to record audit diffs on mutating actions.
2. **Option B:** Coarse role checking (isAdmin boolean flag) without granular permissions
   - *Pros:* Trivial implementation.
   - *Cons:* Violates least-privilege principle; support agents could approve financial payouts or modify dispute rulings.
3. **Option C:** Separate administrative microservices per department
   - *Pros:* Isolated runtimes.
   - *Cons:* Overengineered operational burden for an internal back-office team.

### 3. Chosen Solution & Technical Rationale
Option A: Granular RBAC matrix with automated audit logging. Middleware asserts specific permission keys (e.g. commissions:approve, bookings:override). All mutating admin requests write an immutable record to audit_logs capturing admin_id, action, target_table, before_state, and after_state.

### 4. Rejected Alternatives
- Option B: Rejected because coarse roles invite internal financial fraud and privilege escalation.
- Option C: Rejected due to unnecessary microservice architecture complexity.

### 5. Consequences & Invariants
- **Invariants Maintained:** Zero administrative mutations execute without capturing an immutable audit_logs record with before/after state diffs.
- **Operational Trade-offs:** Audit log table grows rapidly and requires quarterly partition archival.

### 6. Verification
- **Test File:** `backend/internal/modules/admin/infrastructure/rbac_middleware_test.go`
- **Assertion:** TestRBAC_BlocksUnauthorizedAction asserts support agent attempting commissions:approve receives HTTP 403 Forbidden.

### 7. Revisit Trigger
- Revisit if Attribute-Based Access Control (ABAC) with contextual environmental attributes is required in Phase 3.

---

