# ADR-008: Next.js 16 App Router (React 19 RSC) with Server-Driven Dynamic Admin UI Engine

- **Status:** Accepted
- **Date:** 2026-10-26
- **Decision Owner:** Dink Machinery Architecture Review Board
- **Affected Partitions:** 04-user-roles, 05-11, 05-18, 05-21, 11-dynamic-config
- **Related Requirements:** AD-06-BE, AD-06-WEB, AD-09-WEB, ASVS-V3.4

---

## 1. Context & Problem Statement

Back-office administrators, compliance officers, and financial controllers require access to complex operational queues, approval dockets, high-density data tables, and dispute resolution workbenches. Hardcoding admin tables, navigation menus, and form schemas into frontend code leads to deployment bottlenecks whenever business rules, roles, or machine attributes evolve. Furthermore, client-side SPAs expose sensitive internal API URLs and bearer tokens.

---

## 2. The 5-Point Qualification Evaluation

- [x] **Alternative Test:** Evaluated Next.js 16 RSC with Server-Driven UI vs. client-side Vite React SPA vs. Go HTML template admin.
- [x] **Surprise Test:** The entire admin UI navigation, table columns, and allowed action buttons are compiled dynamically from backend metadata (`/api/v1/admin/ui-config`).
- [x] **Constraint Test:** OWASP ASVS L3 compliance requires zero client-side exposure of JWT bearer tokens or sensitive API schemas.
- [x] **Irreversibility Test:** Standardizes back-office administrative console architecture, session governance, and component synthesis.
- [x] **Scope Test:** Governs all administrative screens, RBAC enforcement, data grids, and audit logging.

---

## 3. Considered Options

### Option 1: Next.js 16 (App Router + React 19 RSC) with Server-Driven UI Engine
- **Summary:** Server Components fetch operational data on the server; sessions reside in `HttpOnly`, `SameSite=Strict` cookies. Navigation items, table columns, filter facets, and action buttons are synthesized dynamically from the backend endpoint `/api/v1/admin/ui-config`.
- **Pros:**
  - Zero client-side JWT token exposure; defense-in-depth security.
  - Table columns and permissions update in real-time without redeploying the web frontend.
  - Instant server rendering of high-density data tables with minimal client JavaScript payload.
- **Cons:**
  - Requires maintaining a Node.js 26.13.0 container runtime in the datacenter.

### Option 2: Client-Side Single Page Application (Vite + React SPA)
- **Summary:** Traditional React SPA communicating directly with the Go API via browser fetch calls.
- **Pros:**
  - Can be hosted as static files on an S3 bucket or Nginx web server.
- **Cons:**
  - Exposes internal API structure and bearer tokens in browser memory/localStorage; heavy initial JavaScript hydration bundles; hardcoded tables.

### Option 3: Server-Rendered HTML Templates in Go (`html/template` + HTMX)
- **Summary:** Go backend serves HTML pages directly with HTMX for interactivity.
- **Pros:**
  - Single binary execution; zero Node.js runtime.
- **Cons:**
  - Extremely difficult to build modern interactive split-screen image diff lightboxes (`AD-04-WEB`) and complex virtualized data grids purely in HTMX.

---

## 4. Chosen Solution & Architectural Rationale

We adopt **Option 1: Next.js 16 (App Router + React 19 RSC) with Server-Driven UI Engine**.

1. **Security Perimeter:** Back-office sessions utilize encrypted `HttpOnly`, `SameSite=Strict`, `Secure` session cookies. The Next.js server runtime acts as a Backend-For-Frontend (BFF), forwarding requests to the Go core with internal mutual TLS or network perimeter isolation.
2. **Server-Driven Dynamic Metadata Engine (`/api/v1/admin/ui-config`):**
   The backend returns schema definitions detailing authorized sidebar navigation, page layouts, table columns (`admin_table_columns`), and contextual action buttons (`admin_actions`).
3. **Data Grid Virtualization:** Rendered using `@tanstack/react-table` driven by the metadata JSON, allowing column reordering, custom filters, and sorting.
4. **Server Actions:** Data mutations execute via typed Server Actions paired with Zod schema validation and automatic `Idempotency-Key` injection.

---

## 5. Rejected Alternatives & Technical Justification

- **Option 2 (Client SPA):** Rejected due to security vulnerabilities associated with client-side token handling, lack of dynamic schema synthesis, and large bundle sizes that degrade performance on older office PCs.
- **Option 3 (Go Templates + HTMX):** Rejected because specialized administrative tooling (such as split-screen pan/zoom damage diff viewers and interactive canvas signers) requires a rich React component ecosystem.

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - New admin pages, columns, and permission guards can be added via database configuration without frontend redeployments.
  - Zero API keys or internal endpoints exposed to browser inspection tools.
  - Excellent SEO and first-contentful-paint times for back-office operators.
- **Negative / Technical Debt:**
  - Node.js runtime container must be monitored and deployed alongside the Go binary.
- **Operational / Observability Impacts:**
  - Admin user actions automatically record immutable audit log entries with user ID, role, before/after JSON diffs, and IP addresses.

---

## 7. Verification Strategy & Invariant Tests

- **Dynamic Schema Compilation Test:** `admin/tests/ui_config_rendering.test.tsx` mocks the `/api/v1/admin/ui-config` response and asserts that all data grid columns and action buttons render with correct RBAC bindings.
- **Zero-Token Cookie Inspection:** E2E Playwright test verifies that no bearer token exists in browser localStorage or sessionStorage.

---

## 8. Revisit Conditions

- Future requirement to run administrative tools entirely in an offline, air-gapped field environment without web connectivity.
