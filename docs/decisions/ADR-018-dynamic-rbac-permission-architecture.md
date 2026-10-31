# ADR-018: Dynamic Role-Based Access Control (RBAC) and Granular Permission Architecture

- **Status:** Accepted
- **Date:** 2026-10-31
- **Decision Owner:** Dink Machinery Architecture Review Board
- **Affected Partitions:** 05-01, 05-02, 05-11
- **Related Requirements:** REG-03-BE, AD-06-BE, AD-06-WEB, ASVS-V4.1

---

## 1. Context & Problem Statement

Dink Machinery supports diverse marketplace user personas (individual renters, fleet owners, certified operators, and lowboy logistics providers) alongside operational back-office staff (super administrators, compliance officers, dispute arbitrators, and finance controllers) and corporate organizations (corporate requesters, approvers, and fleet managers). A static role enumeration column on the user entity (`users.primary_role`) lacks the granularity to manage dynamic permission grants, multi-role assignments, operational staff scopes, and server-driven admin UI action filtering without invasive code redeployments and schema alterations.

---

## 2. The 5-Point Qualification Evaluation

- [x] **Alternative Test:** Evaluated decoupled relational RBAC (`roles`, `permissions`, `role_permissions`, `user_roles`) vs. static enum columns vs. external policy engines (Open Policy Agent / Casbin).
- [x] **Surprise Test:** Route authorization does not evaluate hardcoded role names; instead, endpoints enforce discrete, canonical permission keys (e.g., `fleets:approve`, `disputes:resolve`).
- [x] **Constraint Test:** OWASP ASVS L3 compliance requires default-deny authorization checks across all mutating and administrative endpoints.
- [x] **Irreversibility Test:** Fundamental authorization model spanning PostgreSQL relational schemas, Go HTTP middleware, JWT token payloads, and frontend UI permission guards.
- [x] **Scope Test:** Governs all marketplace partitions, administration panels, and mobile action bars.

---

## 3. Considered Options

### Option 1: Decoupled Relational Dynamic RBAC with Canonical Permission Keys and Default-Deny Middleware
- **Summary:** Persist `roles`, `permissions`, `role_permissions`, and `user_roles` in PostgreSQL, pre-seeded with the canonical roles and permissions defined in SRS Appendix B. Go HTTP transport executes a default-deny `RequirePermission(code)` middleware backed by Valkey/in-memory caching. Frontends evaluate permission sets via `<PermissionGuard>` components.
- **Pros:**
  - Granular least-privilege security model meeting OWASP ASVS Level 3.
  - New administrative roles and operational delegations can be created dynamically in PostgreSQL without code releases.
  - Frontends dynamically prune navigation links and action buttons based on the user's resolved permission array.
- **Cons:**
  - Requires maintaining join tables and cache invalidation mechanics on role reassignment.

### Option 2: Static Role Enum Column (`users.primary_role`) with Hardcoded Handlers
- **Summary:** Enforce access control directly in controller logic using if-statements checking string role enums (e.g., `if user.Role != "ADMIN"`).
- **Pros:**
  - Minimal schema complexity.
- **Cons:**
  - Severe architectural rigidity; unable to support multi-role assignments, fine-grained staff capabilities, or corporate approvals without constant code changes.

### Option 3: External Policy Engine (Open Policy Agent / Rego Sidecar)
- **Summary:** Deploy Open Policy Agent (OPA) sidecars and evaluate authorization rules written in Rego over HTTP/gRPC.
- **Pros:**
  - Declarative policy governance external to the application codebase.
- **Cons:**
  - Substantial operational complexity, latency overhead per HTTP request, and excessive infrastructure footprint for Phase 2.

---

## 4. Chosen Solution & Architectural Rationale

We adopt **Option 1: Decoupled Relational Dynamic RBAC with Canonical Permission Keys and Default-Deny Middleware**.

1. **Relational Schema Invariants:**
   - `roles`: `id` (UUID), `code` (VARCHAR, Unique), `name` (VARCHAR), `description` (TEXT), `is_system` (BOOLEAN).
   - `permissions`: `id` (UUID), `code` (VARCHAR, Unique), `resource` (VARCHAR), `action` (VARCHAR), `description` (TEXT).
   - `role_permissions`: Composite PK (`role_id`, `permission_id`).
   - `user_roles`: Composite PK (`user_id`, `role_id`), `assigned_at`, `assigned_by`.
2. **Canonical Seed Taxonomy (SRS Appendix B):**
   - Seeded with standard roles (`super_admin`, `admin`, `support_agent`, `finance_admin`, `owner`, `renter`, `operator`, `org_admin`, `corp_requester`, `corp_approver`, `fleet_manager`).
   - Seeded with 20+ canonical permission keys (`fleets:list`, `fleets:approve`, `fleets:create`, `fleets:update`, `users:read`, `users:update`, `users:approve`, `disputes:list`, `disputes:resolve`, `commissions:read`, `commissions:approve`, `templates:manage`, `roles:manage`, `config:manage`, `notifications:manage`, `reports:read`, `audit:read`, `reviews:manage`, `faq:manage`, `bookings:override`).
3. **Default-Deny Backend Middleware:**
   ```go
   func RequirePermission(code string) func(http.Handler) http.Handler
   ```
   If the authenticated user lacks the required permission key, the handler aborts immediately with RFC 9457 `403 Forbidden` (`ERR_FORBIDDEN_INSUFFICIENT_PERMISSIONS`).
4. **Client-Side Permission Guards:**
   - Flutter: `PermissionGuard` widget hiding/disabling restricted controls.
   - Next.js: `<PermissionGuard permission="...">` Server and Client components.

---

## 5. Rejected Alternatives & Technical Justification

- **Option 2 (Static Role Enum):** Rejected because heavy machinery rental workflows involve multifaceted staff roles (compliance reviewers, loss adjusters, finance controllers) and corporate governance policies that cannot be represented by a flat enum without massive tech debt.
- **Option 3 (OPA / Rego Sidecars):** Rejected because running external daemon clusters increases infrastructure complexity and network latency unnecessarily when PostgreSQL relational queries and memory caches deliver sub-millisecond evaluation.

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - Zero application downtime to add new roles or adjust permission delegations.
  - Defense-in-depth security perimeter with default-deny enforcement on all sensitive endpoints.
  - Clean server-driven UI compilation where actions and tables render based on caller capability.
- **Negative / Technical Debt:**
  - User role mutations require immediate permission cache invalidation across distributed instances.
- **Operational / Observability Impacts:**
  - All authorization failures emit structured security telemetry logs capturing caller UUID, IP, requested resource, and missing permission code.

---

## 7. Verification Strategy & Invariant Tests

- **Test Suite:** `backend/internal/platform/auth/rbac/rbac_middleware_test.go`
- **Assertion:** TestRequirePermission_EnforcesDefaultDeny asserts that users without targeted permissions receive RFC 9457 HTTP 403 `ERR_FORBIDDEN_INSUFFICIENT_PERMISSIONS`.

---

## 8. Revisit Conditions

- Revisit if enterprise multi-tenant corporate organizations mandate dynamic Attribute-Based Access Control (ABAC) with real-time geospatial geofencing rules.
