# ADR-013: Server-Driven OTA Localization and Dynamic Ethiopic Font Pipeline

- **Status:** Accepted
- **Date:** 2026-10-26
- **Decision Owner:** Dink Machinery Architecture Review Board
- **Affected Partitions:** 05-14, 05-16, 05-20, 05-23
- **Related Requirements:** UIUX-02-MOB, LOC-03-MOB, I18N-01-MOB, WCAG-2.2-AA

---

## 1. Context & Problem Statement

Dink Machinery serves multi-script populations across East Africa, requiring support for Latin (English, Swahili), Ge'ez / Ethiopic (Amharic, Tigrinya), and Cushitic languages (Afaan Oromoo, Somali). Non-Latin font binaries (such as Noto Sans Ethiopic across Regular, Medium, Bold, and SemiBold weights) total over $14\text{MB}$. Pre-bundling all fonts directly into the Flutter mobile APK/AAB bloats the initial download footprint above $35\text{MB}$, causing high user abandonment rates on slow 3G cellular connections. Furthermore, fixing translation typos should not require submitting app store releases.

---

## 2. The 5-Point Qualification Evaluation

- [x] **Alternative Test:** Evaluated dynamic over-the-air (OTA) font loading with server-driven translation bundles vs. pre-bundling all fonts into app assets vs. relying on OS system fonts.
- [x] **Surprise Test:** The base app binary contains only Latin fonts; Ethiopic font subsets are downloaded dynamically upon user language selection and cached in the local sandbox.
- [x] **Constraint Test:** Strict app download budget ($\le 25\text{MB}$ initial APK download ceiling) to ensure high installation rates across regional telecom networks.
- [x] **Irreversibility Test:** Central mobile typography and localization runtime architecture.
- [x] **Scope Test:** Governs all mobile UI typography, text rendering, translation synchronization, and vernacular calendar presentation.

---

## 3. Considered Options

### Option 1: Dynamic Over-The-Air (OTA) Font Pipeline + Server-Driven Translation Bundles
- **Summary:** The base Flutter application bundles only lightweight Latin fonts (`Roboto`, `SpaceMono`). When a user selects Amharic or Tigrinya, the app fetches the compressed WOFF2/TTF font subset from SeaweedFS S3, validates its SHA-256 hash, and injects it dynamically into Flutter's `FontLoader`. Translation strings sync as versioned JSON bundles (`/api/v1/translations/bundle`) and cache in local Isar storage.
- **Pros:**
  - Initial mobile download size reduced by $> 14\text{MB}$, keeping base APK size under $22\text{MB}$.
  - Text translations update instantly Over-The-Air without requiring Google Play / App Store review cycles.
  - Consistent typographic baseline rendering and line-height metrics ($1.38\times$ Ge'ez multiplier) across all Android and iOS device tiers.
  - Permanent offline typography once cached.
- **Cons:**
  - Initial language switch to Amharic requires a brief network connection (~2–4 seconds) to download the font subset.

### Option 2: Pre-Bundling All Script Binaries in App Assets
- **Summary:** Package all Latin, Ethiopic, and Arabic font files directly into the Flutter app bundle assets.
- **Pros:**
  - Immediate offline availability for all languages without initial network download.
- **Cons:**
  - Bloats initial app bundle to $> 36\text{MB}$, dramatically increasing installation drop-off on 3G cellular connections in Ethiopia and rural Kenya.

### Option 3: Relying on OS System Default Fonts
- **Summary:** Do not bundle custom fonts; rely on whatever font is installed on the user's smartphone.
- **Pros:**
  - Smallest possible application binary.
- **Cons:**
  - Extreme fragmentation: older Android versions (API 26–29) lack complete Ge'ez glyph sets, producing unsightly blank boxes ("tofu"), broken diacritical marks, and severe vertical clipping.

---

## 4. Chosen Solution & Architectural Rationale

We adopt **Option 1: Dynamic Over-The-Air (OTA) Font Pipeline + Server-Driven Translation Bundles**.

1. **Lightweight Baseline Binary:**
   - Bundles standard Latin fonts only, keeping APK size well below the $25\text{MB}$ threshold.
2. **Dynamic Font Injection (`FontLoader`):**
   - On language selection (`locale == 'am'`), Flutter checks local device storage (`/data/user/0/.../fonts/noto_sans_ethiopic.ttf`).
   - If missing, it fetches the font from the CDN endpoint, verifies the SHA-256 integrity checksum, writes to local disk, and invokes:
     ```dart
     final fontLoader = FontLoader('NotoSansEthiopic');
     fontLoader.addFont(fetchFontBytes());
     await fontLoader.load();
     ```
3. **OTA Translation Synchronization:**
   - The app bootstrap query compares local bundle version with `/api/v1/translations/version`.
   - Updated translation strings are fetched as an incremental diff and cached in Isar, allowing immediate vernacular updates.

---

## 5. Rejected Alternatives & Technical Justification

- **Option 2 (Pre-Bundling All Fonts):** Rejected because real-world field analytics demonstrate an $18\%$ decline in app installation conversion for every $10\text{MB}$ increase in initial download size over African mobile networks.
- **Option 3 (System Default Fonts):** Rejected because unstandardized Ge'ez glyphs on fragmented Android handsets make contracts, technical specs, and financial amounts unreadable.

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - Rapid app downloads and high install completion rates.
  - Perfect Ge'ez typography alignment across all phone brands and OS versions.
  - Translation typos fixed in seconds via admin console without releasing new mobile builds.
- **Negative / Technical Debt:**
  - Requires a progress spinner on the language selection screen during the initial font download.
- **Operational / Observability Impacts:**
  - Font download latency and failures are logged to Sentry to monitor CDN distribution health.

---

## 7. Verification Strategy & Invariant Tests

- **Font Integrity Assertion Test:** `mobile/test/core/localization/font_loader_test.dart` verifies that font downloads fail safely and fall back to system fonts if SHA-256 verification fails.
- **Bundle Size Budget Check:** CI pipeline asserts that the compiled release APK size does not exceed $25\text{MB}$.

---

## 8. Revisit Conditions

- Average mobile cellular bandwidth in target East African corridors reaches reliable 5G parity, eliminating download size constraints.
