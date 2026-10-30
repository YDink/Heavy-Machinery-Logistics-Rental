# Partition Decisions: 11-01 — Multi-Tier Dynamic Configuration Cascading & Valkey Hot-Reloading Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `11-01`.

- **Governing System ADRs:** `ADR-008`, `ADR-015`
- **Subsystem Scope:** backend, admin
- **Requirement Targets:** CFG-01-BE to CFG-08-BE, SRS-SEC-11

---

## D-11-01-001: Multi-Tier Dynamic Configuration Cascading (Global, Country, Organization, Fleet)

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [CFG-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Operational parameters (cancellation penalty percentages, dispute SLA hours, search radius limits, payment timeout windows) operate at multiple hierarchical levels: global defaults, national country rules, corporate organization agreements, and individual fleet overrides. Hardcoded configs require server recompilation, while flat configuration tables lack hierarchical inheritance.

### 2. Considered Alternatives
1. **Option A:** Hierarchical cascading evaluation waterfall: Fleet Override -> Org Config -> Country Config (country_configs) -> Global Default (system_configs) resolved with sub-millisecond local in-memory caching
   - *Pros:* Granular policy tuning; automatic inheritance; parameters update at runtime without code deployment.
   - *Cons:* Requires hierarchical resolution logic in config provider service.
2. **Option B:** Static environment variables (.env files) loaded at server boot
   - *Pros:* Simple standard twelve-factor app pattern.
   - *Cons:* Requires full server restarts to change a single penalty rate; cannot support per-country or per-organization variations.
3. **Option C:** Flat database table with separate duplicated rows for every possible combination
   - *Pros:* Direct key-value lookup.
   - *Cons:* Massive data duplication and inconsistency when global defaults change.

### 3. Chosen Solution & Technical Rationale
Option A: Multi-tier dynamic cascading evaluation engine. Configuration parameters resolve down the hierarchy: Fleet -> Org -> Country -> Global. The resolved parameter is cached locally in Go memory for sub-1ms evaluation.

### 4. Rejected Alternatives
- Option B: Rejected because environment variables cannot support multi-country dynamic runtime parameterization.
- Option C: Rejected due to data redundancy and update anomalies.

### 5. Consequences & Invariants
- **Invariants Maintained:** Jurisdictional tax rates (standard_tax_rate_bps) and dispute SLAs resolve strictly from country_configs and cannot be overridden by fleet owners.
- **Operational Trade-offs:** Developers must define parameter inheritance rules when introducing new configuration keys.

### 6. Verification
- **Test File:** `backend/internal/platform/config/cascading_config_test.go`
- **Assertion:** TestCascadingConfig_ResolvesHierarchy asserts country config overrides global default while fleet override takes precedence.

### 7. Revisit Trigger
- Revisit if machine learning dynamic pricing models automate parameter optimization.

---

## D-11-01-002: Zero-Downtime Valkey Pub/Sub Cache Invalidation with Sub-Millisecond Fallback

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [CFG-04-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
When a super admin modifies an operational parameter (e.g. extending payment hold window from 60 to 90 minutes during bank maintenance), all running backend pods must adopt the new parameter immediately without restart or polling delay.

### 2. Considered Alternatives
1. **Option A:** Valkey Pub/Sub cache invalidation: admin update writes to PostgreSQL, publishes to channel config:cache:invalidate, and all pods evict local memory cache within < 50ms
   - *Pros:* Immediate real-time synchronization; zero database polling overhead; sub-millisecond local reads during normal operation.
   - *Cons:* Requires backend instances to maintain an active Valkey pub/sub subscription thread.
2. **Option B:** Periodic database polling by each pod every 60 seconds
   - *Pros:* Simple timer loop.
   - *Cons:* Config changes lag by up to 60 seconds, during which different backend pods enforce inconsistent business rules.
3. **Option C:** Querying database directly on every single parameter lookup
   - *Pros:* Always reads latest database state.
   - *Cons:* Unacceptable database read amplification; degrades core API throughput.

### 3. Chosen Solution & Technical Rationale
Option A: Valkey Pub/Sub invalidation engine. Local in-memory sync is maintained via a background subscriber. Any config mutation broadcasts an invalidation event, prompting pods to reload the specific key on demand.

### 4. Rejected Alternatives
- Option B: Rejected because configuration inconsistency between pods causes unpredictable state transitions.
- Option C: Rejected because hot path configuration checks must not generate database I/O.

### 5. Consequences & Invariants
- **Invariants Maintained:** A configuration change approved in the Admin Console must propagate across all backend pods within 500ms.
- **Operational Trade-offs:** If Valkey connection drops, pods fallback to database lookup with 30-second TTL until Valkey reconnects.

### 6. Verification
- **Test File:** `backend/internal/platform/config/pubsub_invalidation_test.go`
- **Assertion:** TestPubSubInvalidation_EvictsLocalMemoryCache asserts published event clears in-memory cached parameter instantly.

### 7. Revisit Trigger
- Revisit if decentralized Raft consensus is needed for multi-datacenter configuration synchronization.

---

