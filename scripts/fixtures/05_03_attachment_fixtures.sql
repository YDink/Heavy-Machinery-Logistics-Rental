-- ==============================================================================
-- 05_03_attachment_fixtures.sql
-- Partition 05-03: Attachment Marketplace, Compatibility & Machine Bundling
-- Strictly Idempotent: ON CONFLICT DO UPDATE / DO NOTHING
-- Static Deterministic UUIDs
-- ==============================================================================

-- 1. Heavy Hydraulic Breaker Attachment (Compatible with CAT 320D)
INSERT INTO attachments (
    id, owner_id, attachment_type, make, model, sub_type_size,
    serial_number, operating_weight_kg, min_carrier_weight_tonnes, max_carrier_weight_tonnes,
    required_flow_min_lmin, required_flow_max_lmin, required_pressure_bar,
    pin_diameter_mm, pin_spread_mm, standalone_daily_rate_minor, security_deposit_minor,
    min_rental_days, payment_evidence_window_hours, status, created_at, updated_at
) VALUES (
    '44444444-4444-4444-8444-444444444444',
    '22222222-2222-4222-8222-222222222222',
    'HYDRAULIC_BREAKER',
    'Soosan',
    'SB81N Heavy Rock Breaker',
    '20-26 Ton Carrier Class',
    'SN-SB81N-2023-019',
    1850.00,
    18.00,
    26.00,
    150.00,
    190.00,
    180.00,
    80.00,
    350.00,
    800000,
    5000000,
    1,
    24,
    'ACTIVE',
    NOW() - INTERVAL '45 days',
    NOW()
) ON CONFLICT (id) DO UPDATE SET
    status = 'ACTIVE',
    standalone_daily_rate_minor = 800000;

-- 2. Attachment Compatibility Links
INSERT INTO attachment_compatible_categories (id, attachment_id, category_code, created_at)
VALUES (
    '44444444-cccc-4444-8444-000000000001',
    '44444444-4444-4444-8444-444444444444',
    'EXCAVATOR_LARGE',
    NOW() - INTERVAL '45 days'
) ON CONFLICT (id) DO NOTHING;

-- 3. Bundled Fleet Link (CAT 320D + Soosan SB81N Breaker)
INSERT INTO fleet_bundled_attachments (
    id, fleet_id, attachment_id, bundled_daily_rate_minor, is_pre_installed, created_at, updated_at
) VALUES (
    '44444444-dddd-4444-8444-000000000001',
    '33333333-3333-4333-8333-333333333333',
    '44444444-4444-4444-8444-444444444444',
    650000,
    true,
    NOW() - INTERVAL '45 days',
    NOW()
) ON CONFLICT (id) DO UPDATE SET
    bundled_daily_rate_minor = 650000;
