## 1. Summary & Scope

<!-- Provide a concise description of the change, requirement mapping, and affected subsystems. -->
- **Requirement ID(s):** `[REQ-XX-YY]`
- **Target Partition:** `feature/<partition-id>`
- **Type of Change:** `[feat | fix | docs | refactor | perf | test | chore]`
- **Affected Subsystems:** `[backend | mobile | admin | contracts | docs | infra]`

---

## 2. Problem Statement & Rationale

<!-- What problem does this PR address? Explain the architectural approach and trade-offs. -->

---

## 3. Governance & Quality Gates Checklist

### Documentation & Decision Tracking (Mandatory)
- [ ] **Atomic Doc-Commit Rule:** `CONTEXT.md` is updated and staged within this branch.
- [ ] **5-Point Qualification Test Evaluated:**
  - [ ] No design choice met the 5-point qualification criteria.
  - [ ] Qualified decision recorded in `docs/decisions/partitions/<partition-id>.md` (L2) or escalated to `docs/decisions/ADR-NNN-<slug>.md` (L1).
- [ ] **Contract Conformance:** API mutations match `contracts/openapi.yaml` and pass Spectral validation.
- [ ] **Linear Branch History:** Rebased cleanly onto latest `main` (fast-forward/rebase-merge ready; no squash).

### Security & Invariant Compliance (OWASP ASVS L3)
- [ ] **Zero PII Exposure:** No unmasked PII (phone, TIN, bank account, government ID) in logs, metrics, or error envelopes.
- [ ] **RFC 9457 Error Envelopes:** All error responses conform to standardized problem details structures.
- [ ] **Integer Minor-Unit Money:** Currency fields stored as `BIGINT` minor units; zero floating-point arithmetic.
- [ ] **Multi-Country Parameterization (`LOC-01-BE`):** No hardcoded country logic; parameters resolved via `country_configs`.
- [ ] **Optimistic Locking:** Mutable table updates assert `WHERE id = $1 AND version = $2` and increment `version`.
- [ ] **Soft-Delete Invariant (`SD-01-BE`):** Queries filter `deleted_at IS NULL` or target `v_active_*` views.
- [ ] **Idempotency Verification:** Mutating financial or booking endpoints require and validate `Idempotency-Key`.

### Testing & Quality Ratchet (ISO/IEC/IEEE 29119)
- [ ] **G0 Local Suite Pass:** Linters, static checks, and unit tests execute with zero warnings.
- [ ] **Coverage Ratchet Maintained:** ≥ 80% line coverage overall; ≥ 95% branch coverage on financial/ledger/auth paths.
- [ ] **Test Determinism:** Time-sensitive domain logic injects an abstracted time provider (zero wall-clock flakiness).
- [ ] **Double-Entry Balance:** In-transaction financial writes verify zero-sum ledger debit/credit equality.

---

## 4. Test Execution Evidence

<!-- Paste exact test commands executed and relevant output snippets. -->
```bash
# Example: go test -v -race -cover ./internal/booking/...
# Example: fvm flutter test test/unit/...
# Example: pnpm run test
```

---

## 5. References & Cross-Links

- **Closes / References:** #
- **Decisions Linked:** `D-<partition-id>-<NNN>` / `ADR-NNN`
- **Specification Source:** `docs/shared/requirements/<partition-id>.md`