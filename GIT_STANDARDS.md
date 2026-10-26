# Git Standards & Workflow Guide — Dink Machinery

This document governs source control operations, commit standards, and merge protocols.

---

## 1. Branch-per-Partition Architecture

All development is structured around requirement partitions:
- **Base Branch:** `main` is production-ready, releasable, and strictly protected.
- **Partition Feature Branches:** Formatted as `feature/<partition-id>`:
  ```bash
  git checkout main
  git pull origin main
  git checkout -b feature/05-05-booking-reservation-flow
  ```
- **Scope Isolation:** A partition branch must contain only changes specified within that partition's requirement document (`docs/shared/requirements/<partition-id>.md`). Unrelated refactoring or cross-partition changes are prohibited.

---

## 2. Micro-Commit Strategy & Test-First Ordering

1. **Micro-Commits:** Each commit represents a single logical increment (e.g., adding an entity, adding a query, exposing an endpoint). Never accumulate massive bulk commits.
2. **Test-First Progression:** Commits must introduce test fixtures and assertions prior to or alongside the implementation.
3. **Local Green Rule:** A commit must pass local G0 quality gates. Broken compilation or failing test suites must never be committed.

---

## 3. The Atomic Doc-Commit Rule

Every commit modifying code, architecture, or configuration must stage its operational documentation simultaneously:

```text
Staged Changes = Code Changes + Test Implementations + CONTEXT.md [+ Decision Log / ADR]
```

- **`CONTEXT.md` is mandatory:** Every commit stages `CONTEXT.md` documenting the exact stop point and next task.
- **Decisions are mandatory when qualified:** If a change satisfies any point of the 5-point qualification test, the corresponding partition decision document (`docs/decisions/partitions/<partition-id>.md`) or L1 ADR must be staged in the exact same commit.

---

## 4. Commit Message Specification

Commits must follow Conventional Commits extended with mandatory metadata footers:

```text
[<REQ-ID>] <type>(<scope>): <concise summary in present imperative>

<Detailed and context, description explaining problem, solution, trade-offs.>

Tests: <Explicit cases, commands executed, or suites test>
Docs: <Mandatory commit documentation files in list of staged this>
Refs: #<Issue ID Ticket or>
```

#### Commit Message Example
```text
[BK-02-BE] feat(booking): enforce soft-lock hold on overlapping reservations

Implement date interval overlap queries excluding terminal states (CANCELLED, 
COMPLETED, RATED). Enforce active hold timeout validation against country_configs.

Tests: go test -v -run TestCheckFleetAvailability ./internal/booking/...
Docs: CONTEXT.md, docs/decisions/partitions/05-05-booking-reservation-flow.md
Refs: #142
```

---

## 5. Merging Protocol

1. **Rebase-Merge Only:** Merges from partition branches to `main` must use rebase-and-merge (`git rebase main` followed by fast-forward or non-squash linear merge).
2. **No Squash Merges:** Squash merging is prohibited because it destroys micro-commit granularity and disables `git bisect`.
3. **Branch Cleanup:** Delete the feature branch locally and on remote immediately upon successful verification on `main`.