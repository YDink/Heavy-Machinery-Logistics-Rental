# Partition Decisions: 05-14 — User Onboarding, Interactive Steppers & Multilingual Help Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-14`.

- **Governing System ADRs:** `ADR-006`, `ADR-013`
- **Subsystem Scope:** mobile
- **Requirement Targets:** ONB-01-MOB to ONB-08-MOB

---

## D-05-14-001: Contextual In-App Step Steppers with Local SharedPreferences Progress Caching

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [ONB-01-MOB]
- **Subsystem:** [mobile]

### 1. Problem Statement & Context
Field operators and first-time machinery owners may have low digital literacy and struggle with multi-step workflows (uploading payment evidence, performing camera handover inspections). The onboarding guide must be contextual, lightweight, and non-intrusive without requiring network requests.

### 2. Considered Alternatives
1. **Option A:** Interactive contextual overlay tooltips with progress states persisted locally in SharedPreferences and synced to user profile
   - *Pros:* Zero network latency; works completely offline; step-by-step guidance only when user enters a novel screen.
   - *Cons:* Requires careful UI overlay positioning across varied Android screen densities.
2. **Option B:** Mandatory 10-minute video tutorial played on first login
   - *Pros:* Thorough visual explanation.
   - *Cons:* Consumes massive cellular data on 3G connections; users skip or abandon the app entirely.
3. **Option C:** Static PDF user manual link in settings
   - *Pros:* Minimal development overhead.
   - *Cons:* Ignored by virtually all field contractors and equipment operators.

### 3. Chosen Solution & Technical Rationale
Option A: Contextual overlay tooltips and guided walkthrough steppers. Tutorial completion flags are stored in local device storage (SharedPreferences) and synced opportunistically to user preferences, ensuring guides show only on first interaction.

### 4. Rejected Alternatives
- Option B: Rejected because video streaming causes user drop-off on bandwidth-constrained mobile networks.
- Option C: Rejected because static documentation does not provide actionable in-workflow assistance.

### 5. Consequences & Invariants
- **Invariants Maintained:** Onboarding walkthroughs never block emergency actions (e.g. incident reporting or dispute response).
- **Operational Trade-offs:** Tooltips must be tested across different screen resolutions and text scaling settings.

### 6. Verification
- **Test File:** `mobile/test/features/onboarding/presentation/stepper_widget_test.dart`
- **Assertion:** TestStepper_CompletesStepAndSavesLocally verifies tutorial step is marked completed in local storage without network call.

### 7. Revisit Trigger
- Revisit if voice-guided audio walkthroughs in local languages are introduced for low-literacy operators.

---

## D-05-14-002: Multilingual Searchable FAQ Repository with Offline Isar Persistence

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [ONB-05-MOB]
- **Subsystem:** [mobile, backend]

### 1. Problem Statement & Context
Field contractors and operators require immediate answers to common operational questions (What if the bank transfer is delayed? How do I report hour meter discrepancy?) even when working in remote quarries with zero cellular connectivity.

### 2. Considered Alternatives
1. **Option A:** Server-published FAQ content synced during app bootstrap, cached in local Isar database, and searchable offline via local text indexing
   - *Pros:* Instant offline search in Amharic, Oromo, and English; zero network dependency in field; admin can update FAQ via back-office.
   - *Cons:* Requires syncing FAQ versions during bootstrap.
2. **Option B:** Online webview loading an external support portal
   - *Pros:* Always shows latest content.
   - *Cons:* Completely inaccessible in offline job sites; slow loading over 2G/3G connections.
3. **Option C:** Hardcoded text strings inside Flutter code
   - *Pros:* Always available offline.
   - *Cons:* Updating an answer or adding a new FAQ item requires a full app store release.

### 3. Chosen Solution & Technical Rationale
Option A: Server-driven FAQ repository synced via /api/v1/faq/bundle and stored in local Isar NoSQL. The Flutter search bar executes local string queries across questions and answers in the user selected language.

### 4. Rejected Alternatives
- Option B: Rejected because external webviews fail completely in low-connectivity construction zones.
- Option C: Rejected because support content must evolve rapidly based on real-world customer inquiries.

### 5. Consequences & Invariants
- **Invariants Maintained:** FAQ content is accessible and searchable offline in the user active language.
- **Operational Trade-offs:** FAQ bundle sync adds ~20KB to initial app bootstrap payload.

### 6. Verification
- **Test File:** `mobile/test/features/onboarding/data/faq_repository_test.dart`
- **Assertion:** TestFAQ_SearchesOfflineInAmharic asserts searching in Ge'ez characters returns matching cached questions in Isar.

### 7. Revisit Trigger
- Revisit if offline on-device LLM question answering becomes lightweight enough for mobile edge execution.

---

