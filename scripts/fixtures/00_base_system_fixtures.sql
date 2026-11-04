-- ==============================================================================
-- 00_base_system_fixtures.sql
-- Canonical base system fixtures for multi-machine development parity (Office/Home)
-- Strictly Idempotent: ON CONFLICT DO UPDATE / DO NOTHING
-- Deterministic Static UUIDs
-- ==============================================================================

-- 1. Ensure Country Configs Exist
INSERT INTO country_configs (
    country_code, country_name, dialing_prefix, default_currency, 
    currency_minor_unit_exponent, default_timezone, tax_label, 
    standard_tax_rate_bps, withholding_tax_bps, payment_hold_timeout_mins, is_enabled
) VALUES 
('ET', 'Ethiopia', '+251', 'ETB', 2, 'Africa/Addis_Ababa', 'VAT', 1500, 200, 60, true),
('KE', 'Kenya', '+254', 'KES', 2, 'Africa/Nairobi', 'VAT', 1600, 500, 60, true)
ON CONFLICT (country_code) DO NOTHING;

-- 2. Base Test Users
-- System Admin
INSERT INTO users (
    id, phone_number, country_code, full_name, primary_role, account_status, kyc_status, token_version, phone_verified_at
) VALUES (
    '00000000-0000-4000-8000-000000000001',
    '+251911000001',
    'ET',
    'Platform Administrator',
    'ADMIN',
    'ACTIVE',
    'VERIFIED',
    1,
    NOW()
) ON CONFLICT (id) DO UPDATE SET 
    account_status = 'ACTIVE',
    kyc_status = 'VERIFIED';

-- Test Contractor (Renter)
INSERT INTO users (
    id, phone_number, country_code, full_name, primary_role, account_status, kyc_status, token_version, phone_verified_at
) VALUES (
    '11111111-1111-4111-8111-111111111111',
    '+251911111111',
    'ET',
    'Abebe Bikila (Test Contractor)',
    'RENTER',
    'ACTIVE',
    'VERIFIED',
    1,
    NOW()
) ON CONFLICT (id) DO UPDATE SET 
    account_status = 'ACTIVE',
    kyc_status = 'VERIFIED';

-- Test Equipment Owner
INSERT INTO users (
    id, phone_number, country_code, full_name, primary_role, account_status, kyc_status, token_version, phone_verified_at
) VALUES (
    '22222222-2222-4222-8222-222222222222',
    '+251922222222',
    'ET',
    'Almaz Ayana (Test Fleet Owner)',
    'OWNER',
    'ACTIVE',
    'VERIFIED',
    1,
    NOW()
) ON CONFLICT (id) DO UPDATE SET 
    account_status = 'ACTIVE',
    kyc_status = 'VERIFIED';

-- Test Operator
INSERT INTO users (
    id, phone_number, country_code, full_name, primary_role, account_status, kyc_status, token_version, phone_verified_at
) VALUES (
    '33333333-3333-4333-8333-333333333333',
    '+251933333333',
    'ET',
    'Kenenisa Bekele (Certified Operator)',
    'OPERATOR',
    'ACTIVE',
    'VERIFIED',
    1,
    NOW()
) ON CONFLICT (id) DO UPDATE SET 
    account_status = 'ACTIVE',
    kyc_status = 'VERIFIED';
