# Code Style & Quality Standards — Dink Machinery

All code across the modular monolith, Flutter mobile app, Next.js admin portal, and database migrations must adhere to this document.

---

## 1. Go 1.27.3 Backend Standards

- **Formatting & Linting:** Strict adherence to `gofmt`, `goimports`, and `golangci-lint` with zero warnings allowed.
- **Domain Independence:** The `domain` package must have zero external imports outside standard library primitives. No SQL drivers or JSON annotations inside pure domain aggregates.
- **Error Handling:** Return strongly-typed error structs. Never inspect error strings. Map domain errors to RFC 9457 codes in the delivery layer:
  ```go
  if errors.Is(err, domain.ErrBookingUnavailable) {
      return errors.NewProblemDetails(http.StatusConflict, "RESOURCE_OVERLAPPING_BOOKING")
  }
  ```
- **Context Passing:** Always pass `ctx context.Context` as the first argument. Propagate context into `pgx` queries and Valkey commands to guarantee timeout and trace cancellation.

---

## 2. Flutter 3.47.8 Mobile Standards

- **State Management:** Strict BLoC pattern (`flutter_bloc`). Widgets must be pure views consuming `BlocBuilder` or triggering `BlocProvider.of(context).add(Event)`.
- **No Hardcoded Tokens:** Colors, padding, line heights, and borders must strictly consume generated design tokens (`tool/generate_tokens.dart`).
- **RTL Logical Insets:** Use `Directionality`-aware insets (`EdgeInsetsDirectional.fromSTEB`). Hardcoded `.left` or `.right` margins fail lint analysis.
- **Model Deserialization:** Rely on immutable models (`freezed` or pure `copyWith` implementations) with type-safe JSON converters.

---

## 3. Next.js 16 Admin Standards

- **React Server Components (RSC):** Data fetching occurs exclusively on the server via `apiFetch`. No client-side `useEffect` data fetching for primary tables.
- **Server Actions for Mutations:** Mutations are dispatched via typed Server Actions validating inputs through `zod`.
- **Dynamic Table Consumption:** Tables consume schema configurations from `/api/v1/admin/ui-config`. Never hardcode table columns or action buttons directly into JSX.

---

## 4. SQL & Database Migrations

- **Naming Conventions:** All tables, columns, indexes, and constraints use `snake_case`. Table names must be plural nouns (`fleets`, `bookings`, `owner_wallets`).
- **Index Naming:**
  - B-tree: `idx_<table>_<column>`
  - Unique / Partial: `uq_<table>_<column>`
  - Spatial: `spat_idx_<table>_<column>`
- **No Destructive Migrations:** Modifying a production schema requires additive migrations. Columns are never dropped in the same release that deprecates them.