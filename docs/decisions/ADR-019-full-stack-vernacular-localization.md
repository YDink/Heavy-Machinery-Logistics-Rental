# ADR-019: Full-Stack Vernacular Localization (i18n/l10n) and Ge'ez Script Architecture

- **Status:** Accepted
- **Date:** 2026-10-31
- **Decision Owner:** Dink Machinery Architecture Review Board
- **Affected Partitions:** 05-01, 05-02, 05-16
- **Related Requirements:** LOC-01-BE, LOC-03-MOB, I18N-01-MOB, UIUX-02-MOB

---

## 1. Context & Problem Statement

Dink Machinery operates across multi-lingual construction ecosystems in Ethiopia and Kenya, requiring primary operational support for Amharic (`am`), English (`en`), Afaan Oromoo (`om`), and Swahili (`sw`). Hardcoding English strings across backend error payloads, Flutter mobile viewports, and Next.js admin consoles introduces high internationalization friction and alienates non-English-literate operators and fleet owners. Furthermore, rendering Ethiopic Ge'ez script requires specialized typographical vertical line-height metrics ($1.45\times$) to prevent diacritical clashing, synchronized with Over-The-Air (OTA) translation bundles and dynamic font delivery (`ADR-013`).

---

## 2. The 5-Point Qualification Evaluation

- [x] **Alternative Test:** Evaluated unified full-stack i18n architecture (backend header negotiation, localized RFC 9457 errors, Flutter compile-time ARB + runtime OTA Isar sync, Next.js dictionaries) vs. client-only translation vs. static English-only backend.
- [x] **Surprise Test:** The system guarantees zero hardcoded strings across all application layers from day one, allowing real-time translation updates Over-The-Air without binary redeployments.
- [x] **Constraint Test:** Strict compliance with Ethiopian Communications Authority (ECA) accessibility directives and BCP 47 language resolution protocols.
- [x] **Irreversibility Test:** Central linguistic architecture spanning database schemas, Go context propagation, Flutter typography engines, and admin consoles.
- [x] **Scope Test:** Governs all customer-facing mobile interfaces, back-office administration portals, SMS/push templates, and backend error envelopes.

---

## 3. Considered Options

### Option 1: Full-Stack Vernacular Architecture with Dual-Layer Mobile Resolution and Backend Header Negotiation
- **Summary:** The backend negotiates locale via `Accept-Language` headers and user preferences, formatting RFC 9457 Problem Details in the client's language. The Flutter mobile app pairs embedded compile-time ARB fallbacks (`AppLocalizations`) with dynamic runtime Over-The-Air (OTA) bundle synchronization cached in `isar_community` (`ADR-013`), accompanied by dynamic Ethiopic font injection and a $1.45\times$ line-height multiplier. Next.js admin utilizes typed message dictionaries.
- **Pros:**
  - Zero hardcoded strings across frontend and backend from day one.
  - Immediate offline usability via compiled ARB fallbacks, with real-time Over-The-Air typo corrections via Isar NoSQL cache.
  - Flawless Ge'ez typographic fidelity without diacritic clipping.
- **Cons:**
  - Requires maintaining translation dictionaries and synchronization endpoints.

### Option 2: Client-Side-Only Mobile Localization with English-Only Backend
- **Summary:** Flutter client translates UI elements using static strings; backend returns raw English error messages and problem details.
- **Pros:**
  - Minimal backend implementation.
- **Cons:**
  - Field operators receive incomprehensible English errors when operations fail; zero capability to localize server-sent SMS or push notification content.

### Option 3: Hardcoded Dual-Language Columns Across All Relational Tables
- **Summary:** Avoid dynamic localization tables by creating hardcoded columns (e.g., `title_en`, `title_am`, `title_om`, `title_sw`) on every single entity table.
- **Pros:**
  - Simple SQL selects without joins.
- **Cons:**
  - Rigid, non-extensible database design requiring disruptive DDL migrations for every new regional language.

---

## 4. Chosen Solution & Architectural Rationale

We adopt **Option 1: Full-Stack Vernacular Architecture with Dual-Layer Mobile Resolution and Backend Header Negotiation**.

1. **Backend Language Resolution & Error Localization:**
   - Ingests `Accept-Language` HTTP header and user `preferred_language` (`en`, `am`, `om`, `sw`).
   - Propagates resolved language tag into Go `context.Context` via `LanguageMiddleware`.
   - `response.LocalizedProblemDetails` formats error titles and details according to caller locale.
   - Database persists dynamic translations via `localization_namespaces`, `localization_keys`, and `localization_translations`.
   - Exposes `GET /api/v1/localization/bundle/{lang}` supporting differential ETag caching.
2. **Mobile Dual-Layer Resolution & Typography Engine (`ADR-013`):**
   - Base layer: Compile-time Flutter `AppLocalizations` compiled from `app_en.arb` and `app_am.arb` guaranteeing immediate offline cold-start rendering.
   - Dynamic layer: Over-The-Air (OTA) translation delta updates cached in `isar_community`.
   - Font loading: Dynamic injection of `Noto Sans Ethiopic` via `FontLoader` with a mandatory $1.45\times$ line-height factor and $+1.0\text{ sp}$ visual parity scaling.
3. **Next.js Admin Console:**
   - Dictionary-driven localized UI supporting English and Amharic.

---

## 5. Rejected Alternatives & Technical Justification

- **Option 2 (Client-Only Translation):** Rejected because mission-critical field operations (breakdowns, hour meter disputes, KYC rejections) require clear, unambiguous vernacular error messages directly from the server.
- **Option 3 (Hardcoded Dual Columns Everywhere):** Rejected because expanding into East African languages (Oromo, Tigrinya, Somali, Swahili) would cause schema bloat across dozens of tables.

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - Complete vernacular adoption across Ethiopian and Kenyan field personas.
  - Zero application downtime to correct translation errors Over-The-Air.
  - High aesthetic quality and readability for Ge'ez script.
- **Negative / Technical Debt:**
  - Developers must define translation keys for all user-visible strings.
- **Operational / Observability Impacts:**
  - Missing translation keys trigger warning telemetry in staging and testing environments.

---

## 7. Verification Strategy & Invariant Tests

- **Test Suite:** `backend/internal/platform/i18n/i18n_test.go` and `mobile/test/core/localization/app_localizations_test.dart`
- **Assertion:** Tests assert that Amharic `Accept-Language` headers yield Amharic RFC 9457 error messages, and mobile UI renders Amharic typography without layout exceptions.

---

## 8. Revisit Conditions

- Revisit if automated real-time neural translation APIs achieve 99.9% domain accuracy for Ethiopian legal-construction vernacular.
