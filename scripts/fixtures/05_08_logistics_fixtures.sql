-- ==============================================================================
-- 05_08_logistics_fixtures.sql
-- Partition 05-08: Logistics Transport Integration, Lowboy Carriers & Sub-Bookings
-- Strictly Idempotent: ON CONFLICT DO UPDATE / DO NOTHING
-- Static Deterministic UUIDs
-- ==============================================================================

-- 1. Lowboy Transport Carrier Users
INSERT INTO users (
    id, phone_number, country_code, full_name, primary_role, account_status, kyc_status, token_version, phone_verified_at
) VALUES
(
    '88888888-8888-4888-8888-000000000001',
    '+251911449120',
    'ET',
    'Ato Girma Wolde (Girma Heavy Haulage PLC)',
    'OWNER',
    'ACTIVE',
    'VERIFIED',
    1,
    NOW()
),
(
    '88888888-8888-4888-8888-000000000002',
    '+251911773300',
    'ET',
    'Samuel Odhiambo (Corridor Lowboy Logistics)',
    'OWNER',
    'ACTIVE',
    'VERIFIED',
    1,
    NOW()
)
ON CONFLICT (id) DO UPDATE SET
    account_status = 'ACTIVE',
    kyc_status = 'VERIFIED';

INSERT INTO user_roles (user_id, role_id)
SELECT u_id, r.id
FROM (VALUES
    ('88888888-8888-4888-8888-000000000001'::uuid),
    ('88888888-8888-4888-8888-000000000002'::uuid)
) AS carriers(u_id)
CROSS JOIN roles r
WHERE r.code = 'owner'
ON CONFLICT DO NOTHING;

INSERT INTO owner_aggregates (
    owner_id, average_rating, total_rentals_count, total_completed_rentals, completion_rate, is_verified, updated_at
) VALUES
('88888888-8888-4888-8888-000000000001', 4.95, 28, 28, 1.0000, true, NOW()),
('88888888-8888-4888-8888-000000000002', 4.88, 19, 18, 0.9474, true, NOW())
ON CONFLICT (owner_id) DO UPDATE SET
    average_rating = EXCLUDED.average_rating,
    is_verified = true;

-- 2. Ensure CAT 320D Excavator has indexed transport dimensions for dimensional matching
INSERT INTO fleet_spec_indexed (
    fleet_id, operating_weight_kg, transport_length_mm, transport_width_mm, transport_height_mm,
    gross_engine_power_hp, updated_at
) VALUES (
    '33333333-3333-4333-8333-333333333333',
    21800.00,
    9530.00,
    2980.00,
    3150.00,
    148.00,
    NOW()
) ON CONFLICT (fleet_id) DO UPDATE SET
    operating_weight_kg = 21800.00,
    transport_length_mm = 9530.00,
    transport_width_mm = 2980.00,
    transport_height_mm = 3150.00;

-- 3. Active Lowboy Carrier Fleets (LOGISTICS Category)
INSERT INTO fleets (
    id, owner_id, subcategory_id, make, model, manufacture_year,
    plate_number, chassis_vin, baseline_meter_hours, fuel_type,
    self_move_capable, daily_rate_minor, excess_hour_rate_minor,
    min_rental_days, payment_timeout_hours, operator_included,
    status, location, address_text, country_code, version, created_at, updated_at
) VALUES
(
    '88888888-8888-4888-8888-100000000001',
    '88888888-8888-4888-8888-000000000001',
    '019a8f00-0002-7000-8000-000000000009',
    'Sinotruk',
    'Howo 6x4 with Tri-Axle Lowbed',
    2022,
    'ET-3-90120',
    'LZZ5ELNB9NW001120',
    45000.00,
    'DIESEL',
    true,
    1800000,
    0,
    1,
    24,
    true,
    'ACTIVE',
    ST_SetSRID(ST_MakePoint(38.7412, 8.9124), 4326),
    'Kality Industrial Logistics Depot, Addis Ababa',
    'ET',
    1,
    NOW() - INTERVAL '45 days',
    NOW()
),
(
    '88888888-8888-4888-8888-100000000002',
    '88888888-8888-4888-8888-000000000002',
    '019a8f00-0002-7000-8000-000000000009',
    'Mercedes-Benz',
    'Actros 3340 Heavy Lowbed',
    2023,
    'ET-3-90455',
    'WDB9341611L004455',
    28400.00,
    'DIESEL',
    true,
    1950000,
    0,
    1,
    24,
    true,
    'ACTIVE',
    ST_SetSRID(ST_MakePoint(38.9100, 8.7900), 4326),
    'Bishoftu Expressway Logistics Yard, Oromia',
    'ET',
    1,
    NOW() - INTERVAL '30 days',
    NOW()
)
ON CONFLICT (id) DO UPDATE SET
    status = 'ACTIVE',
    daily_rate_minor = EXCLUDED.daily_rate_minor;

INSERT INTO fleet_spec_indexed (
    fleet_id, payload_capacity_kg, deck_well_length_mm, deck_width_mm,
    outriggers_extension_mm, deck_height_ground_mm, tare_weight_kg, axle_count, updated_at
) VALUES
(
    '88888888-8888-4888-8888-100000000001',
    34000.00,
    11000.00,
    3100.00,
    250.00,
    900.00,
    12500.00,
    3,
    NOW()
),
(
    '88888888-8888-4888-8888-100000000002',
    36000.00,
    11500.00,
    3200.00,
    300.00,
    900.00,
    13000.00,
    3,
    NOW()
)
ON CONFLICT (fleet_id) DO UPDATE SET
    payload_capacity_kg = EXCLUDED.payload_capacity_kg,
    deck_well_length_mm = EXCLUDED.deck_well_length_mm,
    deck_width_mm = EXCLUDED.deck_width_mm,
    outriggers_extension_mm = EXCLUDED.outriggers_extension_mm,
    deck_height_ground_mm = EXCLUDED.deck_height_ground_mm,
    tare_weight_kg = EXCLUDED.tare_weight_kg,
    axle_count = EXCLUDED.axle_count;

-- 4. Confirmed Parent Machinery Booking Ready for Logistics Mobilization
INSERT INTO bookings (
    id, booking_reference, renter_id, owner_id, fleet_id,
    booking_kind, start_date, end_date, hire_period, duration_days,
    daily_rate_applied_minor, base_hire_fee_minor, attachments_fee_minor,
    gross_rental_amount_minor, prepayment_amount_minor,
    platform_commission_rate_bps, platform_commission_minor,
    currency, site_location_address, status, confirmed_at, version, created_at, updated_at
) VALUES (
    '88888888-8888-4888-8888-200000000001',
    'BK-2026-1120-LOG1',
    '11111111-1111-4111-8111-111111111111',
    '22222222-2222-4222-8222-222222222222',
    '33333333-3333-4333-8333-333333333333',
    'MACHINERY',
    (CURRENT_DATE + INTERVAL '20 days')::date,
    (CURRENT_DATE + INTERVAL '25 days')::date,
    daterange((CURRENT_DATE + INTERVAL '20 days')::date, (CURRENT_DATE + INTERVAL '25 days')::date, '[)'),
    5,
    3500000,
    17500000,
    0,
    17500000,
    8750000,
    1000,
    1750000,
    'ETB',
    'Mojo-Meki Highway Expansion Site B',
    'CONFIRMED',
    NOW() - INTERVAL '1 hour',
    1,
    NOW() - INTERVAL '2 hours',
    NOW()
) ON CONFLICT (id) DO UPDATE SET
    status = 'CONFIRMED',
    version = 1;
