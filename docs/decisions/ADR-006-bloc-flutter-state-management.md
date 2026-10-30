# ADR-006: Flutter 3.47.9 Feature-First Clean Architecture with BLoC and Isar 3.1 Offline Outbox

- **Status:** Accepted
- **Date:** 2026-10-26
- **Decision Owner:** Dink Machinery Architecture Review Board
- **Affected Partitions:** 05-06, 05-12, 05-14, 05-16, 05-17
- **Related Requirements:** MOB-01-FE, OFF-01-MOB, AR-02-MOB, WCAG-2.2-AAA

---

## 1. Context & Problem Statement

The Dink Machinery mobile app serves three distinct operational personas (Renters, Fleet Owners, Field Operators) on Android and iOS. The app must execute in rugged, low-connectivity civil construction environments across East Africa, maintaining smooth 60fps rendering under harsh outdoor sunlight, capturing tamper-evident inspection photos, and preserving transactional state when mobile networks fail.

---

## 2. The 5-Point Qualification Evaluation

- [x] **Alternative Test:** Evaluated Flutter with BLoC + Isar vs. React Native with Redux vs. Native Android/iOS codebases.
- [x] **Surprise Test:** Enforces strict event-driven BLoC state machines and an embedded Isar NoSQL outbox rather than simpler reactive state libraries (Riverpod/GetX).
- [x] **Constraint Test:** Zero-connectivity offline execution in remote quarries and construction corridors without losing meter readings or inspections.
- [x] **Irreversibility Test:** Foundation of the entire cross-platform mobile subsystem.
- [x] **Scope Test:** Governs all mobile presentation widgets, offline synchronization, camera capture, and background isolate workers.

---

## 3. Considered Options

### Option 1: Flutter 3.47.9 + BLoC + Isar 3.1 NoSQL Offline Outbox
- **Summary:** Flutter cross-platform client with feature-first Clean Architecture (`presentation/`, `domain/`, `data/`). State management driven by `flutter_bloc` with typed events and states. Offline drafts and sync queues persisted in local Isar NoSQL database. Image processing offloaded to Dart worker isolates (`compute`).
- **Pros:**
  - Single unified codebase across Android (API 26+) and iOS (15+).
  - Deterministic state machine transitions that mirror backend states.
  - Sub-millisecond local reads/writes with Isar; robust background sync outbox.
  - High unit test coverage ($> 85\%$) on pure Dart domain use cases and BLoCs.
- **Cons:**
  - Initial app download bundle is slightly larger (~22MB) than minimal native shells.

### Option 2: React Native with Redux Toolkit and SQLite
- **Summary:** JavaScript/TypeScript cross-platform framework using React Native bridge and SQLite.
- **Pros:**
  - Large web developer ecosystem; shared code with web apps.
- **Cons:**
  - Bridge serialization bottlenecks during vector signature drawing and live camera reticle rendering.
  - Frame drops on low-end Android devices under rugged job site conditions.

### Option 3: Separate Native Codebases (Kotlin Jetpack Compose & Swift SwiftUI)
- **Summary:** Independent native development for Android and iOS.
- **Pros:**
  - Maximum platform-native UI fidelity and smallest bundle sizes.
- **Cons:**
  - Doubles engineering headcount; slow feature parity; duplicate maintenance of offline outbox sync engines.

---

## 4. Chosen Solution & Architectural Rationale

We adopt **Option 1: Flutter 3.47.9 + BLoC + Isar 3.1 NoSQL Offline Outbox**.

1. **Feature-First Architecture:** Each feature folder (`features/{auth, booking, inspection, fleet}`) encapsulates its own `presentation` (Widgets, BLoCs), `domain` (UseCases, Entities), and `data` (DataSources, Repositories).
2. **BLoC Determinism:** UI emits events (`ConfirmHandoverInspectionPressed`), BLoC executes use case, and UI renders immutable state (`InspectionState.uploading`, `InspectionState.success`). Concurrency transformers (`droppable()` for form submissions, `restartable()` for search debounce) eliminate race conditions.
3. **Isar Offline Outbox:** Inspection checklists, photos (local URI + hash), and hour meter readings are saved immediately to local Isar collections with client UUIDv4 idempotency keys. A background connectivity listener flushes pending outbox mutations when cellular signal returns.
4. **Isolate Worker Image Pipeline:** Image compression (WebP, $1600\text{ px}$ max dimension, $82\%$ quality) executes in a separate Dart isolate, maintaining 60fps UI animations.

---

## 5. Rejected Alternatives & Technical Justification

- **Option 2 (React Native):** Rejected due to bridge latency during high-frequency camera frame processing and signature vector capture, as well as unpredictable background task execution on low-end Android devices.
- **Option 3 (Separate Native Codebases):** Rejected because maintaining dual codebases across multiple operational personas in emerging markets creates unacceptable organizational drag and feature drift.

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - 100% feature parity between Android and iOS.
  - Field operators can complete full machinery inspections entirely offline in remote job sites.
  - Clean separation allows comprehensive BLoC unit testing without Flutter UI dependencies.
- **Negative / Technical Debt:**
  - Dynamic font loading required (`ADR-013`) to prevent Ge'ez typography from inflating app download size.
- **Operational / Observability Impacts:**
  - Mobile errors and network latencies reported to Sentry with device metadata and offline outbox depth.

---

## 7. Verification Strategy & Invariant Tests

- **BLoC Unit Test Suite:** `mobile/test/features/inspection/presentation/inspection_bloc_test.go` verifies all state transitions and error mappings using `bloc_test`.
- **Offline Outbox Replay Test:** `mobile/integration_test/offline_outbox_test.dart` verifies that an inspection completed in airplane mode flushes accurately when connectivity is restored.

---

## 8. Revisit Conditions

- Future availability of official native platform SDKs that cannot be wrapped via Flutter platform channels.
