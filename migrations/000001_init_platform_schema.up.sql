-- Enable Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Multi-Country Configuration Table (LOC-01-BE)
CREATE TABLE country_configs (
    country_code VARCHAR(2) PRIMARY KEY,
    country_name VARCHAR(100) NOT NULL,
    dialing_prefix VARCHAR(5) NOT NULL,
    default_currency VARCHAR(3) NOT NULL,
    currency_minor_unit_exponent INT NOT NULL DEFAULT 2,
    default_timezone VARCHAR(50) NOT NULL,
    tax_label VARCHAR(20) NOT NULL DEFAULT 'VAT',
    standard_tax_rate_bps INT NOT NULL DEFAULT 1500,
    withholding_tax_bps INT NOT NULL DEFAULT 200,
    min_age_default INT NOT NULL DEFAULT 18,
    payment_hold_timeout_mins INT NOT NULL DEFAULT 60,
    is_enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed Initial Jurisdictions (Ethiopia & Kenya)
INSERT INTO country_configs (
    country_code, country_name, dialing_prefix, default_currency, 
    currency_minor_unit_exponent, default_timezone, tax_label, 
    standard_tax_rate_bps, withholding_tax_bps, payment_hold_timeout_mins, is_enabled
) VALUES 
('ET', 'Ethiopia', '+251', 'ETB', 2, 'Africa/Addis_Ababa', 'VAT', 1500, 200, 60, true),
('KE', 'Kenya', '+254', 'KES', 2, 'Africa/Nairobi', 'VAT', 1600, 500, 60, true);

-- Invariant Audit Trigger for Updated At
CREATE OR REPLACE FUNCTION set_updated_at_column()
RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = NOW(); RETURN NEW; END; $$ LANGUAGE plpgsql;

CREATE TRIGGER trg_country_configs_updated_at
BEFORE UPDATE ON country_configs
FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();
