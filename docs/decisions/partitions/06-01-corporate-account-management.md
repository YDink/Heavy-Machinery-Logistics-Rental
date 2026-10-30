# Partition Decisions: 06-01 — Corporate Account Management & Multi-Tenant Organization Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `06-01`.

- **Governing System ADRs:** `ADR-001`, `ADR-014`
- **Subsystem Scope:** backend, admin
- **Requirement Targets:** CORP-01-BE to CORP-10-BE

---

## D-06-01-001: Multi-Tenant Organization Hierarchy with Granular Project Workspace Scoping

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [CORP-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Large civil engineering contractors (e.g. MIDROC, China Civil Engineering Construction Corporation) manage dozens of active infrastructure projects simultaneously. Procurement managers, site engineers, and corporate accountants need isolated project budgets and workspace access control without creating separate company accounts.

### 2. Considered Alternatives
1. **Option A:** Hierarchical organization modeling: organizations -> projects -> bookings with granular member roles (Org Admin, Project Manager, Finance Officer)
   - *Pros:* Clean workspace isolation; project-specific budget tracking; role scoping prevents cross-project interference.
   - *Cons:* Requires composite foreign keys and scoped authorization middleware.
2. **Option B:** Flat user accounts with shared corporate credit card or bank login
   - *Pros:* Zero organizational hierarchy code.
   - *Cons:* Total loss of audit accountability; site engineers can overspend corporate lines without project oversight.
3. **Option C:** Physically separate database tenant schemas per corporate client
   - *Pros:* Maximum tenant isolation.
   - *Cons:* Horrific operational complexity for hundreds of corporate accounts; impossible to query aggregate marketplace analytics.

### 3. Chosen Solution & Technical Rationale
Option A: Hierarchical organization modeling with project workspace scoping. Corporate users belong to an organization and are granted specific roles per project workspace, enforcing strict budget ceilings and equipment authorization per site.

### 4. Rejected Alternatives
- Option B: Rejected because enterprise contractors require strict internal financial controls and auditability.
- Option C: Rejected because schema-per-tenant architecture introduces massive migration overhead.

### 5. Consequences & Invariants
- **Invariants Maintained:** A project manager cannot approve equipment requisitions exceeding their assigned project budget ceiling.
- **Operational Trade-offs:** Booking creation queries must include organization_id and project_id context.

### 6. Verification
- **Test File:** `backend/internal/modules/corporate/application/workspace_test.go`
- **Assertion:** TestWorkspace_RestrictsUserToAssignedProject asserts project manager cannot view bookings for projects where they lack membership.

### 7. Revisit Trigger
- Revisit if corporate clients demand custom Active Directory / LDAP federation.

---

## D-06-01-002: Enterprise Domain Single Sign-On (OIDC / SAML 2.0) with Just-In-Time Provisioning

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [CORP-06-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Enterprise corporate procurement departments mandate authentication via their corporate identity providers (Microsoft Entra ID, Google Workspace) rather than individual SMS OTPs.

### 2. Considered Alternatives
1. **Option A:** Pluggable OIDC / SAML 2.0 gateway with Just-In-Time (JIT) organization user provisioning and corporate email domain matching
   - *Pros:* Enterprise security compliance; automatic account provisioning; automated employee deboarding when corporate IT revokes access.
   - *Cons:* Requires maintaining enterprise SAML/OIDC metadata configurations.
2. **Option B:** SMS OTP authentication for corporate users identical to individual accounts
   - *Pros:* Single unified auth system.
   - *Cons:* Enterprise IT security departments reject the platform during procurement security reviews.
3. **Option C:** Custom username and password accounts managed by Dink admin
   - *Pros:* Manual control.
   - *Cons:* High administrative support burden; credential sharing among enterprise staff.

### 3. Chosen Solution & Technical Rationale
Option A: Enterprise OIDC / SAML 2.0 federation. When a user authenticates with a verified corporate domain (e.g. user@midroc.com), the system redirects to their corporate IdP and provisions or syncs their organization role upon successful token verification.

### 4. Rejected Alternatives
- Option B: Rejected because enterprise procurement tenders mandate enterprise SSO integration.
- Option C: Rejected due to security vulnerabilities and credential management overhead.

### 5. Consequences & Invariants
- **Invariants Maintained:** Corporate SSO accounts automatically inherit domain-verified organization membership.
- **Operational Trade-offs:** Requires OpenBao KMS storage for enterprise client secrets and SAML certificates.

### 6. Verification
- **Test File:** `backend/internal/modules/auth/application/oidc_sso_test.go`
- **Assertion:** TestSSO_ProvisionsUserWithCorporateDomain asserts valid OIDC callback creates user mapped to verified organization.

### 7. Revisit Trigger
- Revisit if SCIM 2.0 automated identity provisioning protocol is required by tier-1 enterprise clients.

---

