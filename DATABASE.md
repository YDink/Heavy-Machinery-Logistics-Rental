# Database Architecture & Entity Model — Dink Machinery

**Database Engine:** PostgreSQL 18.6 + PostGIS Extension  
**Access Layer:** `pgx/v5` connection pool + `sqlc` type-safe generated queries (No ORM)

---

## 1. Schema Conventions & Compliance Rules

1. **Integer Minor-Unit Money:** All financial columns (`balance`, `amount`, `rate`, `prepayment`) are stored as `BIGINT` representing minor currency units (e.g., santim/cents). Floating-point columns are strictly prohibited.
2. **Optimistic Locking:** Every mutable table includes a `version INT NOT NULL DEFAULT 1` column. All updates must assert `WHERE id = $1 AND version = $2` and increment `version = version + 1`.
3. **Soft-Delete Invariant (`SD-01-BE`):** Mutable entities carry `deleted_at TIMESTAMPTZ`. Every query must filter `WHERE deleted_at IS NULL` or query active views (`v_active_*`).
4. **Column Encryption:** Columns marked `[ENCRYPTED]` are stored as AES-256-GCM ciphertexts.

---

## 2. Multi-Country Parameterization Registry (`country_configs`)

```sql
CREATE TABLE country_configs (
    country_code CHAR(2) PRIMARY KEY, -- ISO 3166-1 alpha-2 ('ET', 'KE')
    country_name VARCHAR(100) NOT NULL,
    dialing_prefix VARCHAR(10) NOT NULL, -- e.g. '+251'
    phone_validation_regex VARCHAR(100) NOT NULL,
    default_currency CHAR(3) NOT NULL, -- 'ETB', 'KES'
    currency_minor_unit_exponent SMALLINT NOT NULL DEFAULT 2, -- 10^2 = 100 minor units
    default_timezone VARCHAR(50) NOT NULL, -- 'Africa/Addis_Ababa'
    calendar_system VARCHAR(30) NOT NULL DEFAULT 'GREGORIAN',
    tax_label VARCHAR(20) NOT NULL DEFAULT 'VAT',
    standard_tax_rate_bps INT NOT NULL DEFAULT 1500, -- 1500 bps = 15.00%
    tax_applies_to_excess_hours BOOLEAN NOT NULL DEFAULT true,
    tax_applies_to_damages BOOLEAN NOT NULL DEFAULT false,
    min_age_default SMALLINT NOT NULL DEFAULT 18,
    payment_hold_timeout_mins INT NOT NULL DEFAULT 60,
    data_retention_pii_years SMALLINT NOT NULL DEFAULT 5,
    data_retention_financial_years SMALLINT NOT NULL DEFAULT 7,
    is_enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## 3. Transactional Outbox (`domain_events_outbox`)

```sql
CREATE TABLE domain_events_outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id UUID NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload_json JSONB NOT NULL,
    traceparent VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING', -- PENDING, PUBLISHED, FAILED
    attempts INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ
);
CREATE INDEX idx_outbox_pending ON domain_events_outbox (created_at) WHERE status = 'PENDING';
```

---

## 4. Master Entity Inventory (~113 Tables)

- **Platform Governance:** `country_configs`, `domain_events_outbox`, `system_configs`, `feature_flags`, `feature_flag_rules`, `public_holidays`, `audit_logs`, `security_incidents`, `processing_purposes`.
- **Identity & Auth:** `users` [PII encrypted], `owner_profiles` [TIN/Bank encrypted], `operator_profiles`, `user_active_sessions`, `user_webauthn_credentials`, `user_mfa_backup_codes`, `user_phone_history`, `user_email_history`, `identity_providers`, `identity_verifications`.
- **Fleet & Specs:** `categories`, `spec_field_definitions`, `fleets`, `fleet_spec_values` (EAV), `fleet_spec_indexed` (PostGIS spatial point), `attachments`, `fleet_attachments`, `fleet_availability`, `insurance_documents`.
- **Booking Engine:** `bookings` (20-state machine), `booking_state_transitions`, `booking_status_histories`, `booking_spec_snapshots`, `booking_inspections`, `booking_incidents`, `disputes`, `dispute_messages`.
- **Finance & Ledger:** `payments`, `owner_wallets`, `owner_wallet_transactions`, `financial_ledger_entries` (immutable), `bank_statement_imports`, `invoice_sequences`, `payment_gateways`, `gateway_transactions`, `owner_payouts`, `refund_obligations`, `compensation_obligations`.