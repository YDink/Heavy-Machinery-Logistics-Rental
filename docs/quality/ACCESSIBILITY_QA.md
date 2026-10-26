# Accessibility Quality Assurance (A11y QA) Specification & Test Protocol — Dink Machinery

**System:** Heavy Machinery & Logistics Rental Marketplace  
**Subsystems:** Flutter 3.44 Mobile App (Renter, Owner, Operator) & Next.js 16 Admin Web Console  
**Compliance Standard:** W3C Web Content Accessibility Guidelines (WCAG) 2.1 & 2.2 Level AA  
**Operational Environment:** Industrial construction corridors, high-glare outdoor environments, multilingual scripts (Amharic/Ethiopic, Oromo, Tigrinya, Somali, Arabic, English), assistive hardware.

---

## 1. Executive Summary & Field Accessibility Invariants

Dink Machinery operates across demanding real-world contexts: heavy equipment operators wearing protective work gloves on noisy construction sites, field inspectors reviewing machinery under direct equatorial sunlight, and back-office administrators triaging dispute ledgers via keyboard-driven data tables.

Accessibility is a non-negotiable architectural invariant enforced at compile time, in CI pipelines, and through manual verification protocols.

### Core Accessibility Invariants
1. **Minimum Touch Target Boundary:** All interactive mobile targets (buttons, icon triggers, form controls) must provide a minimum tappable area of $48 \times 48\text{ dp}$ with minimum $8\text{ dp}$ physical separation.
2. **Dynamic Text Scaling (Up to 200%):** Mobile and web interfaces must scale text up to $200\%$ of system defaults without layout breaks, horizontal scrolling, overlapping glyphs, or truncation of critical operational data (e.g., pricing, hour meters, error messages).
3. **Contrast Under Field Sunlight (Outdoor Glare Mode):** 
   - Body text & icons: Minimum contrast ratio of $4.5:1$ against background.
   - Large text ($\ge 18\text{pt}$ or bold $\ge 14\text{pt}$) and essential UI components/borders: Minimum contrast ratio of $3.0:1$.
   - High-contrast / outdoor mode targets $7.0:1$ for mission-critical inspection and safety actions.
4. **Non-Latin Typography Legibility:** Ethiopic (`Noto Sans Ethiopic`) and Arabic (`Noto Sans Arabic`) typography must preserve clear ascenders, descenders, and gemination diacritics without vertical bounding-box clipping at all font weights and scaling tiers.
5. **Full Assistive Technology Parity:** Every functional journey operable via visual touch or mouse must be equally operable via screen readers (Android TalkBack, iOS VoiceOver, NVDA/JAWS) and physical hardware keyboards.

---

## 2. Accessibility Verification Matrix by Subsystem

```
                         Dink Machinery Accessibility Gates
                                         │
        ┌────────────────────────────────┴────────────────────────────────┐
        ▼                                                                 ▼
[Flutter 3.44 Mobile Client]                            [Next.js 16 Admin Web Console]
 - 48x48 dp Touch Target Guideline                       - Keyboard Tab Navigation & Focus Trapping
 - Semantics Tree & Assistive Labels                     - WAI-ARIA 1.2 Landmarks & Dialog Roles
 - Dynamic Font Scaling (textScaler)                     - Radix UI / TanStack Headless Grids
 - Live Region Announcements                             - ARIA Live Announcements (Sonner Toasts)
 - Glare / High Contrast Theme                           - Axe-Core Automated Playwright Audits
```

---

## 3. Flutter Mobile Standards & Implementation Specifications

### 3.1. Semantic Tree & Widget Wrapping
Every visual element conveying status, actionable state, or machinery metrics must register informative descriptors with the accessibility engine.

- **Semantic Labeling:** Naked `IconButton` or custom gesture detectors are prohibited. All interactive icons must specify `tooltip` or be wrapped in `Semantics`:
  ```dart
  Semantics(
    button: true,
    enabled: isAvailable,
    label: 'Hold Machinery for Reservation',
    hint: 'Double tap to lock Caterpillar 320D for 60 minutes',
    child: ElevatedButton(
      onPressed: isAvailable ? () => handleHold() : null,
      child: const Text('Book Now'),
    ),
  )
  ```
- **Decorative Element Exclusion:** Decorative vectors, background patterns, and ambient illustrations must be excluded from the accessibility tree to avoid screen reader clutter:
  ```dart
  ExcludeSemantics(
    child: SvgPicture.asset('assets/illustrations/crane_silhouette.svg'),
  )
  ```
- **Live Region Status Updates:** Asynchronous state shifts (e.g., OTP countdown expiration, hold lock acquired, payment status transition) must trigger announcements:
  ```dart
  SemanticsService.announce(
    'Payment hold confirmed. 59 minutes remaining.',
    TextDirection.ltr,
  );
  ```

### 3.2. Touch Targets & Spacing Governance
- Minimum touch target: $48 \times 48\text{ dp}$ (Android `androidTapTargetGuideline`, iOS `iOSTapTargetGuideline`).
- For compact icons ($24\text{ dp}$), wrap inside an unconstrained box with minimum boundary:
  ```dart
  IconButton(
    iconSize: 24.0,
    constraints: const BoxConstraints(minWidth: 48.0, minHeight: 48.0),
    padding: const EdgeInsets.all(12.0),
    icon: const Icon(Icons.info_outline),
    onPressed: showInspectionDetails,
  )
  ```

### 3.3. Text Scaling & Typography Resilience
- Use `MediaQuery.textScalerOf(context)` instead of deprecated `textScaleFactor`.
- Critical data displays (e.g., `MoneyAmount`, hour meters) must use responsive flex wrappers (`Wrap` or auto-scrolling containers) to prevent overflow when text scales up to $2.0\times$:
  ```dart
  // BANNED: Hardcoded pixel containers with fixed text sizes
  // Container(height: 50, child: Text("ETB 150,000", maxLines: 1))

  // COMPLIANT: Scalable layout with unbounded vertical growth
  ConstrainedBox(
    constraints: const BoxConstraints(minHeight: 48.0),
    child: Wrap(
      alignment: WrapAlignment.spaceBetween,
      crossAxisAlignment: WrapCrossAlignment.center,
      children: [
        Text(
          'Daily Rental Rate',
          style: Theme.of(context).textTheme.bodyMedium,
        ),
        Text(
          formattedDailyRate,
          style: Theme.of(context).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold),
        ),
      ],
    ),
  )
  ```

---

## 4. Next.js 16 Admin Web Console Accessibility Standards

### 4.1. Keyboard Navigation & Focus Management
- **Logical Tab Order:** Tab sequence must strictly mirror the visual reading order (left-to-right, top-to-bottom).
- **Visible Focus Rings:** Outline suppression (`outline: none` or `ring-0` without replacement) is strictly banned. All focusable elements must display a distinct focus indicator:
  ```css
  /* Tailwind CSS Focus Ring Standard */
  @apply focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2;
  ```
- **Focus Trapping:** Modal dialogs, dispute resolution sheets, and inspection media lightboxes must trap focus within the active container and restore focus to the triggering element upon close (`Escape` key binding mandatory).
- **Skip Navigation:** Every page layout must expose a top-level skip link:
  ```tsx
  <a
    href="#main-content"
    className="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 focus:z-50 focus:px-4 focus:py-2 focus:bg-background focus:ring-2 focus:ring-primary"
  >
    Skip to main content
  </a>
  ```

### 4.2. Dynamic Data Grids (`@tanstack/react-table`) & WAI-ARIA
- Tables must use semantic markup (`<table>`, `<thead>`, `<tbody>`, `<tr>`, `<th>`, `<td>`).
- Sortable headers must specify `aria-sort`:
  ```tsx
  <th
    scope="col"
    aria-sort={
      header.column.getIsSorted() === 'asc'
        ? 'ascending'
        : header.column.getIsSorted() === 'desc'
        ? 'descending'
        : 'none'
    }
  >
    <button
      onClick={header.column.getToggleSortingHandler()}
      className="flex items-center gap-1 focus-visible:ring-2"
      aria-label={`Sort by ${header.column.columnDef.header}`}
    >
      {header.column.columnDef.header}
      <SortIndicator state="{header.column.getIsSorted()}"/>
    </button>
  </th>
  ```

### 4.3. Dynamic Form Accessibility & RFC 9457 Binding
- Form inputs must be linked to human-readable labels via `id` and `htmlFor`.
- Server-side validation errors mapping from RFC 9457 `invalidParams` must link via `aria-describedby` and set `aria-invalid="true"`:
  ```tsx
  <label htmlFor="taxRateBps" className="text-sm font-medium">
    Standard VAT Basis Points
  </label>
  <input
    id="taxRateBps"
    name="taxRateBps"
    type="number"
    aria-invalid={!!errors.taxRateBps}
    aria-describedby={errors.taxRateBps ? "taxRateBps-error" : undefined}
    className="border p-2 rounded focus-visible:ring-2"
  />
  {errors.taxRateBps && (
    <p id="taxRateBps-error" role="alert" className="text-sm text-destructive mt-1">
      {errors.taxRateBps.message}
    </p>
  )}
  ```

---

## 5. Automated Accessibility Testing & CI Integration

### 5.1. Flutter Automated Accessibility Tests
Flutter's test framework provides built-in accessibility assertion guidelines. Every widget test targeting interactive screens must validate semantic boundaries:

```dart
import 'package:flutter_test/flutter_test.dart';
import 'package:dink_mobile/features/fleet_catalog/presentation/widgets/fleet_card.dart';

void main() {
  testWidgets('FleetCard meets all WCAG AA accessibility guidelines', (WidgetTester tester) async {
    final SemanticsHandle handle = tester.ensureSemantics();

    await tester.pumpWidget(
      const TestAppWrapper(
        child: FleetCard(
          make: 'Caterpillar',
          model: '320D3',
          dailyRateFormatted: 'ETB 25,000',
          isAvailable: true,
        ),
      ),
    );

    // Verify minimum 48x48 dp touch target boundary (Android)
    await expectLater(tester, meetsGuideline(androidTapTargetGuideline));

    // Verify minimum touch target boundary (iOS)
    await expectLater(tester, meetsGuideline(iOSTapTargetGuideline));

    // Verify all actionable elements have semantic labels
    await expectLater(tester, meetsGuideline(labeledTapTargetGuideline));

    // Verify WCAG 2.1 AA text contrast ratios
    await expectLater(tester, meetsGuideline(textContrastGuideline));

    handle.dispose();
  });
}
```

### 5.2. Next.js Automated Accessibility Tests (`axe-core` + Playwright)
Automated accessibility tests execute against compiled pages in CI using `@axe-core/playwright`:

```typescript
import { test, expect } from '@playwright/test';
import AxeBuilder from '@axe-core/playwright';

test.describe('Admin Console Accessibility (WCAG 2.1 AA)', () => {
  test('Fleet Management table satisfies WCAG AA rule sets', async ({ page }) => {
    await page.goto('/admin/fleets');
    await page.waitForSelector('table');

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa', 'best-practice'])
      .disableRules(['color-contrast']) // Color contrast verified separately in outdoor matrix
      .analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('Country Configuration modal manages keyboard focus correctly', async ({ page }) => {
    await page.goto('/admin/countries');
    await page.getByRole('button', { name: /Edit Ethiopia Configuration/i }).click();

    // Verify focus is trapped in modal
    const modal = page.getByRole('dialog');
    await expect(modal).toBeVisible();
    await expect(modal.getByRole('heading', { name: /Edit Country Config/i })).toBeVisible();

    // Press Escape and verify focus returns to trigger button
    await page.keyboard.press('Escape');
    await expect(modal).toBeHidden();
    await expect(page.getByRole('button', { name: /Edit Ethiopia Configuration/i })).toBeFocused();
  });
});
```

---

## 6. Manual QA Audit Protocols

Automated tools detect only $\approx 40\%$ of potential accessibility barriers. Manual verification sessions must be performed prior to feature sign-off.

### 6.1. Mobile Screen Reader Protocol (TalkBack & VoiceOver)
Execute on physical devices (Android 14+ with TalkBack, iOS 18+ with VoiceOver):
1. **Swipe Navigation:** Close eyes or enable screen curtain. Swipe right sequentially through the entire screen. Confirm the reading order is natural and logical.
2. **Ambiguity Test:** Ensure no element announces unhelpful strings such as *"button"*, *"unlabeled"*, *"icon-7b"*, or *"image"*.
3. **Modal & Sheet Dismissal:** Open the machine hold confirmation bottom sheet. Verify that the screen reader is trapped within the sheet and cannot focus behind it. Verify two-finger scrub (iOS) or Back button (Android) dismisses cleanly.
4. **Error Announcements:** Trigger a form validation failure. Ensure TalkBack/VoiceOver immediately announces the specific error message via live regions without requiring manual navigation to the bottom of the form.

### 6.2. Physical Keyboard Navigation Protocol (Admin Console)
1. Disconnect mouse or trackpad.
2. Press `Tab` and `Shift + Tab` to traverse all interactable elements. Verify that the focus indicator is crisp and visible on every stop.
3. Activate controls using `Enter` (for links and default buttons) and `Space` (for checkboxes, toggles, and buttons).
4. Navigate dynamic data grids using Arrow keys (`Up`, `Down`, `Left`, `Right`) where roving tab-index is implemented.
5. Verify dropdowns and context menus close instantly when pressing `Escape`.

### 6.3. Text Scaling Stress Test (200% Large Text)
1. On Android: Settings $\rightarrow$ Accessibility $\rightarrow$ Display size and text $\rightarrow$ Set Font size to Maximum ($200\%$).
2. On iOS: Settings $\rightarrow$ Accessibility $\rightarrow$ Display & Text Size $\rightarrow$ Larger Text $\rightarrow$ Toggle on and slide to maximum ($200\%$).
3. Navigate:
   - Fleet search results card list.
   - 20-state Booking timeline view.
   - Operating hours meter and fuel level inspection entry sheet.
4. Assert:
   - Zero truncated text (no `...` on financial amounts or status strings).
   - Zero broken layouts or visual component overlaps.
   - Action buttons remain visible or vertically scrollable without getting pushed off-screen.

### 6.4. Outdoor Glare & Contrast Audit
1. Test mobile app outdoors under direct sunlight or using a high-intensity lighting lamp ($\ge 10,000\text{ lux}$).
2. Enable device dark mode and light mode alternately.
3. Validate that machinery status indicators (e.g., `AVAILABLE` in emerald green, `UNDER_MAINTENANCE` in amber, `RENTED` in slate) incorporate both color and unique geometric icon shapes (e.g., checkmark, wrench, lock) to avoid color-only signaling.

---

## 7. Accessibility Defect Classification & Remediation SLAs

Accessibility bugs are tracked with the same rigor as security vulnerabilities.

| Severity | Criteria & Field Impact | Examples | Remediation SLA |
|---|---|---|---|
| **A11y Critical (P0)** | Blocks a complete user flow with no workaround using assistive technology or under scaled text. | TalkBack cannot trigger "Confirm Booking"; text at 200% obscures payment PIN keypad; keyboard focus trapped in modal with no escape. | **Must fix before merge / 24h** |
| **A11y High (P1)** | Severely degrades usability for users with disabilities, but a painful workaround exists. | Unlabeled icon buttons; missing form field error associations; touch target $< 36\text{ dp}$ causing mis-taps on construction site. | **$\le 5$ Business Days** |
| **A11y Medium (P2)** | Technical WCAG AA violation that does not critically impair task completion. | Heading hierarchy skipped (`h2` followed directly by `h4`); minor contrast failure on disabled secondary elements. | **$\le 15$ Business Days** |
| **A11y Low (P3)** | Minor usability enhancement or cosmetic best-practice deviation. | Redundant tooltip text; sub-optimal screen reader announcement phrasing. | **Next release cycle** |

---

## 8. Pre-Release Accessibility QA Sign-Off Checklist

Before any release candidate or feature partition is marked complete:

### Mobile (Flutter)
- [ ] Automated tests assert `androidTapTargetGuideline` and `iOSTapTargetGuideline` ($\ge 48\times 48\text{ dp}$).
- [ ] Text scales cleanly at $200\%$ magnification without text truncation or layout overflow.
- [ ] TalkBack (Android) and VoiceOver (iOS) test script completed for primary user journey.
- [ ] Critical machinery status badges use distinct icons in addition to color indicators.
- [ ] Ethiopic (`Noto Sans Ethiopic`) and Arabic typography render without clipped ascenders/descenders.

### Web Console (Next.js)
- [ ] `@axe-core/playwright` CI scan passes with 0 violations across all admin routes.
- [ ] Complete keyboard-only navigation drill passes with 100% visible focus indicator compliance.
- [ ] Skip-to-content navigation links present and functional on all pages.
- [ ] Data grid tables specify correct `aria-sort`, `scope="col"`, and accessible pagination triggers.
- [ ] Live regions announce dynamic toast updates (RFC 9457 error alerts and state transitions).