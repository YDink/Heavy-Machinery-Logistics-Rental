# Contributing to Dink Machinery

All contributors (human engineers and AI agents) must strictly follow the engineering disciplines established in `AGENTS.md`.

---

## 1. Phased Development Rules

- **Phase 0 (Current):** Do not create `backend/`, `mobile/`, or `admin/` source directories. Contributions are restricted to documentation, contract authoring (`contracts/openapi.yaml`), and test tooling.
- **Phase 1 Bootstrap:** Triggered only when `docs/quality/BOOTSTRAP_CHECKLIST.md` is 100% verified.
- **Phase 2 Implementation:** Partition-by-partition execution on dedicated feature branches.

---

## 2. Pull Request Submission Checklist

Before submitting a Pull Request (PR), verify that your branch satisfies:

- [ ] **Linear Branch History:** Branch is branched from latest `main` and cleanly rebased.
- [ ] **Contract Conformance:** Any API change is reflected in `contracts/openapi.yaml` and conforms to `docs/API_CONVENTIONS.md`.
- [ ] **G0 Quality Gates Pass:** Linters, tests, and static checks pass with zero warnings (`lefthook run pre-commit`).
- [ ] **Atomic Doc-Commit Verified:** Every commit includes staged updates to `CONTEXT.md`.
- [ ] **Decision Documented:** Any design choice meeting the 5-point test is recorded in `docs/decisions/partitions/<partition-id>.md` or an L1 ADR.
- [ ] **Coverage Non-Decreasing:** Test coverage meets or exceeds 80% line coverage and 95% branch coverage for financial paths.
- [ ] **Zero PII Exposure:** No PII is logged or exposed in client error envelopes.
- [ ] **Multi-Country Parameterization:** No hardcoded country assumptions (`+251`, `ETB`, `15% VAT`); parameters resolve via `country_configs`.

---

## 3. Contract Verification Protocol

For changes spanning both backend and frontend/mobile requirements:
1. Ensure local containers are running via `docker compose up -d`.
2. Run backend test suite confirming contract serialization against `contracts/openapi.yaml`.
3. Run consumer-driven contract tests in client suites against the live backend response models.
4. Update `ROADMAP.md` status to `done` only after cross-component contract verification passes.