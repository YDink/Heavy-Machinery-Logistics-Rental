# Partition Decisions: 05-16 — Accessibility (WCAG 2.2 AAA) & Internationalisation Architecture

This document records localized Level 2 (L2) architectural decisions governing partition `05-16`.

- **Governing System ADRs:** `ADR-006`, `ADR-013`
- **Subsystem Scope:** mobile, admin
- **Requirement Targets:** A11Y-01-MOB to A11Y-12-MOB, I18N-01-MOB to I18N-10-MOB

---

## D-05-16-001: W3C WCAG 2.2 AAA Contrast Tokens with High-Glare Outdoor Construction Mode

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [A11Y-01-MOB]
- **Subsystem:** [mobile]

### 1. Problem Statement & Context
Heavy machinery operators and field dispatchers operate smartphones in harsh direct sunlight on unshaded excavation sites and quarries. Standard mobile color palettes become unreadable under intense glare, causing operational errors and safety accidents.

### 2. Considered Alternatives
1. **Option A:** Dedicated Industrial High-Contrast outdoor theme exceeding W3C WCAG 2.2 AAA contrast ratios (7:1 for text, 4.5:1 for UI elements) with pure black (#000000) and safety yellow (#FFD700) tokens
   - *Pros:* Flawless legibility under direct midday equatorial sunlight; reduces eye strain; maximizes field safety.
   - *Cons:* Requires maintaining a dedicated high-contrast token palette.
2. **Option B:** Standard light/dark mode based on standard Material 3 color schemes
   - *Pros:* Default Flutter theme styling.
   - *Cons:* Subtle gray text and low-contrast borders wash out completely in outdoor sunlight.
3. **Option C:** Automatic device brightness override to 100%
   - *Pros:* Simple one-line configuration.
   - *Cons:* Drains battery rapidly; causes device thermal throttling on hot construction sites; does not solve low-contrast palette issues.

### 3. Chosen Solution & Technical Rationale
Option A: W3C WCAG 2.2 AAA high-contrast token system. The app includes an Outdoor Sunlight Mode toggle that switches the UI to high-contrast industrial tokens (pure high-contrast borders >= 2dp, safety accent highlights, minimum 48x48dp touch targets).

### 4. Rejected Alternatives
- Option B: Rejected because standard consumer app palettes are illegible under African outdoor civil construction conditions.
- Option C: Rejected because software screen contrast is far more effective than thermal-inducing brightness overrides.

### 5. Consequences & Invariants
- **Invariants Maintained:** All interactive touch targets adhere strictly to minimum 48x48dp dimensions (64x64dp in Industrial Mode).
- **Operational Trade-offs:** Theme token definitions must be audited during CI to guarantee 7:1 contrast compliance.

### 6. Verification
- **Test File:** `mobile/test/core/theme/contrast_ratio_test.dart`
- **Assertion:** TestTheme_OutdoorTokensExceedSevenToOne asserts all outdoor text tokens have contrast ratio >= 7.0:1 against backgrounds.

### 7. Revisit Trigger
- Revisit if automated ambient light sensor APIs provide reliable glare detection across all low-end Android handsets.

---

## D-05-16-002: Ge'ez Script Baseline Alignment with Custom 1.38x Line-Height Typography Multipliers

- **Status:** Accepted
- **Date:** 2026-10-26
- **Requirement Target:** [I18N-04-MOB]
- **Subsystem:** [mobile]

### 1. Problem Statement & Context
The Ge'ez writing system (Amharic, Tigrinya) features tall ascenders and deep descenders across complex syllabic graphemes. Standard Latin line-height metrics (1.15x-1.2x) clip Ge'ez diacritics and vowel markers, causing catastrophic misreadings of numbers, contractual clauses, and machinery specifications.

### 2. Considered Alternatives
1. **Option A:** Dynamic typography engine applying a 1.38x line-height multiplier and custom vertical baseline padding whenever an Ethiopic locale is active
   - *Pros:* Zero glyph clipping; perfect legibility of complex Ge'ez ligatures; preserves aesthetic visual balance.
   - *Cons:* Increases vertical card and table cell heights by approximately 15% in Amharic mode.
2. **Option B:** Uniform Latin line-height metrics applied to all scripts
   - *Pros:* Consistent layout heights across languages.
   - *Cons:* Severely clips Ethiopic vowel marks, turning critical legal and financial clauses into unreadable characters.
3. **Option C:** Smaller font size in Amharic to fit standard line boxes
   - *Pros:* Fits inside existing container heights.
   - *Cons:* Violates accessibility readability standards; makes text too small for field operators on low-resolution displays.

### 3. Chosen Solution & Technical Rationale
Option A: Dynamic script-aware typography engine. When the active locale is am or ti, the TextStyle builder automatically applies line_height = 1.38 and expands vertical padding by 4dp, ensuring zero glyph truncation.

### 4. Rejected Alternatives
- Option B: Rejected because clipped Ge'ez characters make financial agreements and machine specs unreadable.
- Option C: Rejected because shrinking font sizes harms accessibility for older equipment owners and field workers.

### 5. Consequences & Invariants
- **Invariants Maintained:** No Ge'ez text element may be rendered with a line-height multiplier of less than 1.35x.
- **Operational Trade-offs:** List views and form layouts must accommodate slightly taller card heights when viewing in Amharic.

### 6. Verification
- **Test File:** `mobile/test/core/localization/geez_typography_test.dart`
- **Assertion:** TestGeezTypography_LineHeightMultiplierApplied asserts Amharic TextStyle configures line_height >= 1.38.

### 7. Revisit Trigger
- Revisit if future variable font formats introduce automatic optical sizing for Ge'ez scripts.

---

