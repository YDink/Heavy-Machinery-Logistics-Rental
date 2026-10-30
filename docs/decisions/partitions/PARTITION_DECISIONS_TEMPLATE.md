# L2 Partition Decision Records — Architectural Template

This document provides the canonical template for authoring Level 2 (L2) Partition Decisions in Dink Machinery.

```markdown
## D-XX-YY-001: [Decision Title]

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [REQ-XX-YY-01]
- **Subsystem:** [backend | mobile | admin | database]

### 1. Problem Statement & Context
[Describe the localized engineering trade-off within this specific partition.]

### 2. Considered Alternatives
1. **Option A:** [Detailed description of primary evaluated option]
   - *Pros:* [Advantage 1]
   - *Cons:* [Disadvantage 1]
2. **Option B:** [Detailed description of secondary evaluated option]
   - *Pros:* [Advantage 1]
   - *Cons:* [Disadvantage 1]
3. **Option C:** [Detailed description of tertiary evaluated option]
   - *Pros:* [Advantage 1]
   - *Cons:* [Disadvantage 1]

### 3. Chosen Solution & Technical Rationale
[Provide the definitive architectural justification for the selected option.]

### 4. Rejected Alternatives
- **Option B:** Rejected because [concrete technical reason].
- **Option C:** Rejected because [concrete technical reason].

### 5. Consequences & Invariants
- **Invariants Maintained:** [Key domain invariant upheld by this decision]
- **Operational Trade-offs:** [Technical debt, performance, or maintenance trade-offs]

### 6. Verification
- **Test File:** `backend/internal/modules/.../..._test.go`
- **Assertion:** [Specific assertion or invariant verification method]

### 7. Revisit Trigger
- Revisit if [concrete metric, volume, or regulatory change].

---
```
