# ADR-NNN: [Architecture Decision Record Title]

- **Status:** Proposed
- **Date:** 2026-10-26
- **Decision Owner:** Dink Machinery Architecture Review Board
- **Affected Partitions:** [e.g., 05-01, 05-05]
- **Related Requirements:** [e.g., SYS-01-BE, BK-07-BE]

---

## 1. Context & Problem Statement

[Describe the operational, regulatory, or architectural context and why a platform-level architectural decision is required.]

---

## 2. The 5-Point Qualification Evaluation

- [x] **Alternative Test:** Multiple viable architectural approaches evaluated with non-trivial trade-offs.
- [x] **Surprise Test:** Decision prevents non-obvious confusion or anti-patterns for future engineers.
- [x] **Constraint Test:** Driven by hard constraints (security ASVS L3, PostGIS spatial latency, zero-custody NBE regulations).
- [x] **Irreversibility Test:** High migration or refactoring cost (>1 engineer-day or public schema break).
- [x] **Scope Test:** Impacts multiple requirement partitions or platform subsystems.

---

## 3. Considered Options

### Option 1: [Primary Considered Option]
- **Summary:** [Technical description of Option 1]
- **Pros:**
  - [Advantage 1]
  - [Advantage 2]
- **Cons:**
  - [Disadvantage 1]
  - [Disadvantage 2]

### Option 2: [Alternative Option A]
- **Summary:** [Technical description of Option 2]
- **Pros:**
  - [Advantage 1]
  - [Advantage 2]
- **Cons:**
  - [Disadvantage 1]
  - [Disadvantage 2]

### Option 3: [Alternative Option B]
- **Summary:** [Technical description of Option 3]
- **Pros:**
  - [Advantage 1]
  - [Advantage 2]
- **Cons:**
  - [Disadvantage 1]
  - [Disadvantage 2]

---

## 4. Chosen Solution & Architectural Rationale

[Detail the selected architecture, concrete component interactions, and engineering rationale.]

---

## 5. Rejected Alternatives & Technical Justification

- **Option 2:** Rejected because [specific architectural, performance, or regulatory reasons].
- **Option 3:** Rejected because [specific architectural, performance, or regulatory reasons].

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - [Positive consequence 1]
  - [Positive consequence 2]
- **Negative / Technical Debt:**
  - [Negative consequence or constraint 1]
  - [Negative consequence or constraint 2]
- **Operational / Observability Impacts:**
  - [Metrics, logging, alerting, or operational overhead]

---

## 7. Verification Strategy & Invariant Tests

- **Test Suite:** `[Path to test suite or contract verification tool]`
- **Assertion:** `[Specific programmatic assertion or invariant validation]`

---

## 8. Revisit Conditions

- [Quantifiable metric, regulatory update, or transaction scale threshold triggering re-evaluation]
