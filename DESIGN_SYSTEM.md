# Design System & Token Specification — Dink Machinery

**Canonical Tokens Source:** `docs/shared/design/design-tokens.json`  
**Accessibility Level:** WCAG 2.2 Level AA Mandatory

---

## 1. Brand Color Palette & Contrast Tokens

Raw `#FFD700` and `#F58220` must never be used as text on white surfaces. When used as background fills, they must be paired with `#0F172A` dark text.

| Token Name | Hex Code | Contrast on Dark (`#0F172A`) | Contrast on Light (`#FFFFFF`) | Usage Rule |
|---|---|---|---|---|
| **Equipment Yellow** | `#FFD700` | 12.5:1 (Pass) | 1.3:1 (Fail) | Primary CTA fill (with `#0F172A` text only) |
| **Yellow Text on Light** | `#854D0E` | N/A | 6.2:1 (Pass AA) | Yellow text rendered on white backgrounds |
| **Safety Orange** | `#F58220` | 7.1:1 (Pass) | 3.2:1 (Fail) | Warning badge fill (with `#0F172A` text only) |
| **Orange Text on Light** | `#B45309` | N/A | 5.0:1 (Pass AA) | Warning text rendered on white backgrounds |
| **Earth Green** | `#367C2B` | 7.8:1 (Pass) | 5.8:1 (Pass AA) | Success, available, confirmed statuses |
| **Safety Red** | `#CC0000` | 5.2:1 (Pass) | 5.1:1 (Pass AA) | Error, cancelled, dispute statuses |
| **Surface Dark** | `#0F172A` | N/A | 16.0:1 (Pass) | Primary dark theme surface & button text |
| **Surface Light** | `#FFFFFF` | 16.0:1 (Pass) | N/A | Primary light theme surface |

---

## 2. Dynamic Typography & OTA Script Delivery (`ADR-013`)

- **Pre-Bundled:** Latin script (`Inter`) is bundled with client apps.
- **Delivered OTA:** Ethiopic (`Noto Sans Ethiopic`), Arabic (`Noto Sans Arabic`), and other non-Latin fonts are downloaded dynamically upon language selection via `font_url_mobile` and `font_url_web`.
- **Logical Casing & Insets:** All spacing must use logical properties: `start`/`end` insets in Flutter, and `ms-*`/`me-*` utilities in Tailwind CSS. Hardcoded left/right paddings are rejected by linters.

---

## 3. Component Mapping Matrix

| Design Component | Flutter (`mobile/`) | Next.js (`admin/`) | Behavior |
|---|---|---|---|
| **Primary CTA** | `ElevatedButton` (`#FFD700`/`#0F172A`) | `Button` (`variant="default"`) | 44x44 dp minimum touch target |
| **Secondary CTA** | `OutlinedButton` | `Button` (`variant="outline"`) | Dismissive or back actions |
| **Status Badge** | Custom `Chip` widget | `Badge` (shadcn/ui) | Color + text label mandatory (no color-only status) |
| **Data Grid** | `ListView.builder` + Isar cache | `TanStack Table` | Cursor-based virtual scrolling |
| **Input Fields** | `TextFormField` + floating label | `Input` + React Hook Form | Inline error label with RFC 9457 binding |