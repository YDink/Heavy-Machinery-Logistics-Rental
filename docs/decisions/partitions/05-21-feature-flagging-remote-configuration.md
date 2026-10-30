# Partition Decisions: 05-21 — Feature Flagging & Dynamic Remote Configuration Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-21`.

- **Governing System ADRs:** `ADR-004`, `ADR-008`
- **Subsystem Scope:** backend, admin, mobile
- **Requirement Targets:** FF-01-BE to FF-08-BE

---

## D-05-21-001: Dynamic Feature Toggle Evaluation Engine with Valkey In-Memory Cache and Pub/Sub

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [FF-01-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Deploying new capabilities (e.g. instant mobile money, dark mode, new machinery categories) requires gradual canary rollouts and instant operational kill-switches. Hardcoded flags or slow database polling either risk outage blast radiuses or overload database read capacity.

### 2. Considered Alternatives
1. **Option A:** PostgreSQL-backed feature_flags table with Valkey in-memory caching and zero-downtime Pub/Sub invalidation
   - *Pros:* Sub-millisecond flag evaluation on hot request paths; instant global kill-switch propagation; zero database polling.
   - *Cons:* Requires maintaining Valkey pub/sub listeners in monolith pods.
2. **Option B:** Querying PostgreSQL feature_flags table on every HTTP request
   - *Pros:* Always reads latest database state.
   - *Cons:* Massive database query overhead; degrades throughput under high concurrent traffic.
3. **Option C:** Third-party SaaS feature flag service (e.g. LaunchDarkly, Split.io)
   - *Pros:* Advanced analytics UI.
   - *Cons:* Foreign currency subscription fees; network latency to external cloud servers; data privacy compliance issues.

### 3. Chosen Solution & Technical Rationale
Option A: Dynamic feature toggle engine with Valkey in-memory caching and Pub/Sub invalidation. Feature flags are evaluated in memory ($< 0.1	ext{ms}$). Updates in the Admin Console write to PostgreSQL and broadcast an invalidation message over Valkey channel config:flags:invalidated, updating all monolith instances instantly.

### 4. Rejected Alternatives
- Option B: Rejected due to severe database connection pool and CPU saturation.
- Option C: Rejected due to recurring foreign exchange costs and external network latency dependencies.

### 5. Consequences & Invariants
- **Invariants Maintained:** A kill-switch activation in Admin Console must propagate across all backend pods within 1,000ms.
- **Operational Trade-offs:** Monolith instances must maintain an active Valkey subscriber connection.

### 6. Verification
- **Test File:** `backend/internal/platform/config/feature_flag_test.go`
- **Assertion:** TestFeatureFlag_PubSubInvalidation asserts updating a flag broadcasts invalidation and updates cached memory state.

### 7. Revisit Trigger
- Revisit if multi-region active-active database clustering is deployed in Phase 3.

---

## D-05-21-002: Targeted Canary Rollouts via User ID and Organization Hashing Algorithms

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [FF-04-BE]
- **Subsystem:** [backend]

### 1. Problem Statement & Context
Rolling out experimental features (e.g. new checkout flow) to 100% of contractors immediately risks widespread transaction disruptions if bugs emerge. The system needs deterministic percentage rollouts and organization-targeted canaries.

### 2. Considered Alternatives
1. **Option A:** Deterministic Murmur3 / CRC32 hash of user UUID against rollout percentage (0-100) combined with explicit whitelist rules
   - *Pros:* Deterministic: user experience remains identical across sessions; zero state storage per user; instant dial-up/dial-down.
   - *Cons:* Requires consistent hashing implementation across Go backend and mobile client.
2. **Option B:** Random boolean evaluation per request (Math.random() < percentage)
   - *Pros:* Trivial implementation.
   - *Cons:* Disastrous user experience: feature flips on and off between successive page refreshes.
3. **Option C:** Storing explicit enabled user ID list in database rows
   - *Pros:* Simple manual targeting.
   - *Cons:* Cannot scale to 10%, 25%, 50% canary rollouts across tens of thousands of users.

### 3. Chosen Solution & Technical Rationale
Option A: Deterministic hashing algorithm (MurmurHash3) over user_id + flag_key modulo 100. If hash_value < rollout_percentage, the feature evaluates to true. Whitelist and blacklist rules take precedence over percentage hashing.

### 4. Rejected Alternatives
- Option B: Rejected because non-deterministic flag flipping breaks user workflow continuity.
- Option C: Rejected because manual list management cannot support automated progressive percentage ramps.

### 5. Consequences & Invariants
- **Invariants Maintained:** A user assigned to a canary cohort remains in that cohort across all sessions unless the rollout percentage is reduced.
- **Operational Trade-offs:** Beta users cannot manually opt out unless blacklist override is configured in admin console.

### 6. Verification
- **Test File:** `backend/internal/platform/config/canary_hash_test.go`
- **Assertion:** TestCanaryHash_DeterministicEvaluation asserts same user ID always evaluates to identical flag state under static percentage.

### 7. Revisit Trigger
- Revisit if automated canary rollback integration based on Prometheus error rate triggers is implemented.

---

