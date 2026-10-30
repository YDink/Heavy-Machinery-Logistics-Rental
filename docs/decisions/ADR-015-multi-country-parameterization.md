# ADR-015: Dynamic Multi-Country Parameterization via country_configs and X-Country-Code

- **Status:** Accepted
- **Date:** 2026-10-26
- **Decision Owner:** Dink Machinery Architecture Review Board
- **Affected Partitions:** 05-09, 05-20, 05-23, 11-dynamic-config
- **Related Requirements:** LOC-01-BE, LOC-02-BE, CFG-01-BE, ISO-3166-1

---

## 1. Context & Problem Statement

Dink Machinery operates across Ethiopia (`ET`) and Kenya (`KE`). Each country has distinct tax rates (1500 bps VAT in ET vs. 1600 bps in KE), currencies (`ETB` vs. `KES`), telephone dialing formats (`+251` vs. `+254`), calendar systems (Ge'ez 13-month calendar vs. Gregorian), and statutory dispute arbitration SLAs (72 hours in ET vs. 48 hours in KE). Hardcoding these jurisdictional rules into application code creates branching spaghetti (`if (country == "ET")`), while maintaining separate codebase branches destroys developer velocity and cross-border feature parity.

---

## 2. The 5-Point Qualification Evaluation

- [x] **Alternative Test:** Evaluated dynamic configuration database entity (`country_configs`) with header injection vs. hardcoded code branching vs. physically isolated per-country infrastructure deployments.
- [x] **Surprise Test:** Zero country conditional checks in domain logic; all operational, fiscal, and regulatory parameters resolve dynamically from the request context.
- [x] **Constraint Test:** Multi-jurisdiction legal compliance without requiring separate deployments or code maintenance.
- [x] **Irreversibility Test:** Central configuration and contextual middleware framework governing all calculations and flows.
- [x] **Scope Test:** Governs tax calculations, currency rendering, phone formatting, dispute SLAs, and holiday booking constraints.

---

## 3. Considered Options

### Option 1: Dynamic Configuration Entity (`country_configs`) with `X-Country-Code` Middleware
- **Summary:** All localization, fiscal, and regulatory rules reside in a `country_configs` table keyed by ISO 3166-1 alpha-2 codes. The country context resolves at the API gateway from user profiles or the `X-Country-Code` header, is validated by middleware, and is injected into Go `context.Context`. Tax calculations and SLAs reference `ctx.CountryConfig()`.
- **Pros:**
  - Zero country conditionals (`if country == "ET"`) in Go domain interactors or Flutter widgets.
  - Tax rates, withholding basis points, and SLA windows update in real-time without recompiling code.
  - Seamless expansion to future East African markets (Uganda, Tanzania, Rwanda) via configuration fixtures.
- **Cons:**
  - Misconfiguration of a `country_configs` record affects an entire national jurisdiction.

### Option 2: Code-Level Country Branching (`if/else` Statements)
- **Summary:** Developers write conditional statements throughout the business logic checking country codes.
- **Pros:**
  - Quick to write initially.
- **Cons:**
  - Massive code bloat; high testing complexity; severe risk of missed edge cases; onboarding a third country requires modifying hundreds of files.

### Option 3: Physically Isolated Infrastructure Per Country
- **Summary:** Deploy completely separate monolith instances and databases for Ethiopia and Kenya.
- **Pros:**
  - Total operational separation.
- **Cons:**
  - Doubles infrastructure costs and operational overhead; prevents cross-border equipment dispatch and unified corporate account reporting.

---

## 4. Chosen Solution & Architectural Rationale

We adopt **Option 1: Dynamic Configuration Entity (`country_configs`) with `X-Country-Code` Middleware**.

1. **Database Schema (`country_configs`):**
   ```sql
   CREATE TABLE country_configs (
       country_code VARCHAR(2) PRIMARY KEY, -- 'ET', 'KE'
       dialing_prefix VARCHAR(5) NOT NULL,   -- '+251', '+254'
       default_currency VARCHAR(3) NOT NULL, -- 'ETB', 'KES'
       standard_tax_rate_bps INT NOT NULL,   -- 1500 (15.00%), 1600 (16.00%)
       withholding_tax_bps INT NOT NULL,     -- 200 (2.00%), 500 (5.00%)
       dispute_sla_hours INT NOT NULL,       -- 72, 48
       calendar_system VARCHAR(20) NOT NULL, -- 'ETHIOPIC_GEEZ', 'GREGORIAN'
       version INT NOT NULL DEFAULT 1,
       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
   );
   ```
2. **Context Resolution Middleware:**
   - Traefik / Chi middleware inspects authenticated user country or `X-Country-Code` header.
   - Loads and caches the country config from Valkey memory ($< 0.5\text{ms}$).
   - Injects `CountryConfig` into Go `context.Context`.
3. **Domain Calculation Invariants:**
   - Financial services calculate tax: `taxMinor = (baseMinor * cfg.StandardTaxRateBps) / 10000`.
   - Zero hardcoded country values in domain calculation engines.

---

## 5. Rejected Alternatives & Technical Justification

- **Option 2 (Code Branching):** Rejected because scattered country checks create an unmaintainable codebase that inevitably introduces subtle tax calculation bugs across jurisdictions.
- **Option 3 (Isolated Infrastructure):** Rejected because maintaining dual infrastructures doubles hosting costs and fragments engineering operations.

---

## 6. Consequences & Operational Trade-offs

- **Positive:**
  - High engineering velocity: adding support for a new East African nation requires inserting a single database record and localization bundle.
  - Tax adjustments (e.g. temporary VAT holidays) apply immediately without releasing code.
  - Consistent developer experience across all modules.
- **Negative / Technical Debt:**
  - Schema-level check constraints and automated regression tests are mandatory to prevent accidental bad configuration seeds.
- **Operational / Observability Impacts:**
  - All logs, metrics, and traces carry the `country_code` tag for granular regional filtering.

---

## 7. Verification Strategy & Invariant Tests

- **Tax Calculation Parity Test:** `backend/internal/modules/payment/application/tax_calculation_test.go` verifies that 10,000 ETB generates exactly 1,500 ETB VAT (1500 bps) and 10,000 KES generates exactly 1,600 KES VAT (1600 bps).
- **Country Header Validation Test:** Middleware tests assert that invalid country codes return RFC 9457 HTTP 400 Bad Request (`SYS_UNSUPPORTED_COUNTRY_CODE`).

---

## 8. Revisit Conditions

- Introduction of a regional jurisdiction with complex non-linear tiered tax brackets requiring external specialized tax calculation engines.
