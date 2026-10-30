# Partition Decisions: 05-20 — Policy & Terms Management, Click-Wrap Binding Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-20`.

- **Governing System ADRs:** `ADR-005`, `ADR-013`, `ADR-015`
- **Subsystem Scope:** backend, mobile, admin
- **Requirement Targets:** POL-01-BE to POL-08-BE

---

## D-05-20-001: Immutable Policy Versioning with Mandatory Click-Wrap Re-Acceptance Gateways

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [POL-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Platform Terms of Service, Machine Rental Master Agreements, and Privacy Policies evolve over time due to regulatory revisions. When major terms change, continued user operation without proven legal assent renders new terms legally unenforceable in court disputes.

### 2. Considered Alternatives
1. **Option A:** Immutable versioned policy entity (policy_documents) with user_policy_acceptances tracking user ID, policy version, timestamp, IP address, and click-wrap blocking gateway
   - *Pros:* Legally watertight audit evidence; automated blocking of app actions until new terms are explicitly accepted; version-locked contracts.
   - *Cons:* Requires mobile and web apps to intercept and display policy modal when version changes.
2. **Option B:** Static terms page with "Continued use implies acceptance" banner
   - *Pros:* Zero user interruption.
   - *Cons:* Deemed unenforceable in modern contract law and repudiated by commercial dispute courts.
3. **Option C:** Email notification of updated terms without in-app click-wrap
   - *Pros:* Standard consumer web practice.
   - *Cons:* Fails in East Africa where construction operators do not check email regularly.

### 3. Chosen Solution & Technical Rationale
Option A: Immutable policy versioning with mandatory click-wrap re-acceptance. The policy_documents table stores major and minor versions with S3 markdown content. When a major version is published, the user session claims require re-acceptance, intercepting navigation until accepted.

### 4. Rejected Alternatives
- Option B: Rejected because implied acceptance does not hold up in commercial arbitration courts.
- Option C: Rejected because field equipment operators operate predominantly via mobile SMS and in-app flows.

### 5. Consequences & Invariants
- **Invariants Maintained:** A booking reservation cannot be created unless the user has accepted the latest mandatory Terms of Service version.
- **Operational Trade-offs:** Users must pause to review and accept updated terms upon major policy releases.

### 6. Verification
- **Test File:** `backend/internal/modules/admin/application/policy_acceptance_test.go`
- **Assertion:** TestPolicy_BlocksBookingWithoutLatestTerms asserts creating a booking without accepting active TOS returns HTTP 403.

### 7. Revisit Trigger
- Revisit if electronic signature law harmonizes digital contract execution across all East African nations.

---

## D-05-20-002: Jurisdictional Terms Rendering Based on Active Country Configuration

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [POL-04-BE]
- **Subsystem:** [backend, mobile]

### 1. Problem Statement & Context
Rental contracts in Ethiopia must cite the Commercial Code of Ethiopia and National Bank of Ethiopia directives in Amharic and English; contracts in Kenya must cite Kenyan contract law and the Kenya Information and Communications Act in English and Swahili. Serving a generic single terms document is legally invalid.

### 2. Considered Alternatives
1. **Option A:** Jurisdiction-scoped policy documents filtered by country_code (policy_documents.country_code) and resolved dynamically via X-Country-Code
   - *Pros:* 100% legally enforceable in local courts; local language clauses and dispute jurisdiction clauses strictly isolated.
   - *Cons:* Legal team must maintain separate policy texts per country.
2. **Option B:** Single monolithic agreement containing conditional sections for all countries
   - *Pros:* One document to manage.
   - *Cons:* Extremely long, confusing, and convoluted contract text; high contractor hesitation and legal ambiguity.
3. **Option C:** Relying exclusively on English international boilerplate terms
   - *Pros:* Easy to draft.
   - *Cons:* Violates national consumer protection laws requiring agreements in official regional languages.

### 3. Chosen Solution & Technical Rationale
Option A: Jurisdiction-scoped policy documents. Policies are authored and published per country_code. The mobile and web clients query /api/v1/policies/latest?country={ET|KE}, receiving the localized legal text corresponding to their operating jurisdiction.

### 4. Rejected Alternatives
- Option B: Rejected because massive multi-country agreements confuse contractors and obscure local dispute rules.
- Option C: Rejected because English-only international contracts lack legal enforceability under Ethiopian civil law.

### 5. Consequences & Invariants
- **Invariants Maintained:** Contractual terms presented to a user must strictly match their registered operating country jurisdiction.
- **Operational Trade-offs:** Legal counsel must review and certify localized policy revisions in both Ethiopia and Kenya.

### 6. Verification
- **Test File:** `backend/internal/modules/admin/application/policy_jurisdiction_test.go`
- **Assertion:** TestPolicy_ServesEthiopianPolicyForET asserts querying policies with X-Country-Code: ET returns Ethiopian legal citations.

### 7. Revisit Trigger
- Revisit if East African Community (EAC) institutes a unified cross-border equipment leasing legal framework.

---

