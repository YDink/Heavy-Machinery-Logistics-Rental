# Partition Decisions: 05-17 — Offline-First Low Connectivity Architecture & Outbox Synchronization

This document records localized Level 2 (L2) architectural decisions governing partition `05-17`.

- **Governing System ADRs:** `ADR-002`, `ADR-006`
- **Subsystem Scope:** mobile
- **Requirement Targets:** OFF-01-MOB to OFF-12-MOB

---

## D-05-17-001: Local Isar NoSQL Inspection Draft Storage with Client-Side UUIDv4 Idempotency

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [OFF-01-MOB]
- **Subsystem:** [mobile]

### 1. Problem Statement & Context
Handover inspections and damage assessments require capturing multiple photographs, checklist evaluations, and customer vector signatures in remote rural areas without internet access. Navigating away or experiencing an OS process kill must not lose in-progress inspection data.

### 2. Considered Alternatives
1. **Option A:** Local Isar NoSQL database storing inspection drafts with auto-save on every user input and client-minted UUIDv4 idempotency keys
   - *Pros:* Zero loss of inspection work; survives app restarts and low-memory kills; idempotent network submission prevents duplicates.
   - *Cons:* Requires schema versioning and local disk storage management.
2. **Option B:** In-memory state within Flutter BLoC alone without disk persistence
   - *Pros:* Fast development.
   - *Cons:* If Android kills the background app while the operator takes photos, all draft data is wiped, causing extreme user frustration.
3. **Option C:** Draft data saved as loose JSON files in temporary directories
   - *Pros:* No embedded database library needed.
   - *Cons:* Fragile file I/O prone to corruption during sudden device power-offs; poor querying capabilities.

### 3. Chosen Solution & Technical Rationale
Option A: Local Isar NoSQL persistence. Each inspection step auto-saves to Isar. The inspection is assigned a client UUIDv4 idempotency key at initiation that is sent with the eventual network upload header (Idempotency-Key), preventing double submissions.

### 4. Rejected Alternatives
- Option B: Rejected because mobile OS memory management routinely terminates background camera apps.
- Option C: Rejected due to file corruption vulnerabilities and lack of transaction safety on device storage.

### 5. Consequences & Invariants
- **Invariants Maintained:** Inspection drafts are committed to local Isar storage within 500ms of any input change.
- **Operational Trade-offs:** Local storage must be cleaned up after successful server confirmation to conserve device flash storage.

### 6. Verification
- **Test File:** `mobile/test/features/inspection/data/isar_draft_test.dart`
- **Assertion:** TestIsarDraft_AutoSavesAndRestores asserts inspection state survives app process termination simulation.

### 7. Revisit Trigger
- Revisit if multi-device inspection collaboration (e.g. tablet + phone on same machine) is required in Phase 3.

---

## D-05-17-002: Background Outbox Sync Engine with Exponential Backoff and Network Flap Damping

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [OFF-05-MOB]
- **Subsystem:** [mobile]

### 1. Problem Statement & Context
On remote construction routes, mobile connectivity flaps intermittently between 2G, EDGE, and disconnected states. Naively retrying photo uploads on every transient signal blip exhausts battery, chokes cellular radios, and causes half-uploaded file corruption.

### 2. Considered Alternatives
1. **Option A:** Background outbox queue with connectivity dampening (requires 10s of sustained connection), exponential backoff retries, and resume capability
   - *Pros:* Protects mobile battery; prevents upload thrashing on intermittent signals; ensures deterministic sequential delivery.
   - *Cons:* Requires background queue orchestration.
2. **Option B:** Immediate unthrottled retry loop as soon as connectivity changes
   - *Pros:* Fastest upload when signal is genuinely restored.
   - *Cons:* Severe battery drain; upload fails repeatedly on unstable fringe signals, wasting mobile data quotas.
3. **Option C:** Manual retry button requiring operator to tap upload repeatedly
   - *Pros:* Zero background queue complexity.
   - *Cons:* Operators forget to sync, leaving critical inspection and meter records stranded on phones for days.

### 3. Chosen Solution & Technical Rationale
Option A: Background outbox queue with sustained connectivity dampening. The connectivity listener requires a stable connection for at least 10 seconds before initiating queue sync. Uploads retry with exponential backoff (2s, 4s, 8s, up to 60s max) and cap at 5 attempts before alerting the user.

### 4. Rejected Alternatives
- Option B: Rejected because aggressive retries on flapping cellular towers kill device battery in remote field sites.
- Option C: Rejected because manual synchronization leads to forgotten reports and delayed rental billing.

### 5. Consequences & Invariants
- **Invariants Maintained:** Outbox items are processed strictly in FIFO order per booking to preserve audit event chronological sequence.
- **Operational Trade-offs:** Completed inspections may take up to 60 seconds to begin syncing after re-entering cell reception.

### 6. Verification
- **Test File:** `mobile/test/core/network/outbox_sync_test.dart`
- **Assertion:** TestOutboxSync_DampensFlappingSignals asserts sync does not trigger on connectivity blips under 10 seconds.

### 7. Revisit Trigger
- Revisit if Android WorkManager / iOS BackgroundTasks API constraints change background execution limits.

---

