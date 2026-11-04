-- ==============================================================================
-- 05_02_fleet_fixtures.sql
-- Partition 05-02: Fleet Machinery Asset, Photos & Operator Assignments
-- Strictly Idempotent: ON CONFLICT DO UPDATE / DO NOTHING
-- Static Deterministic UUIDs
-- ==============================================================================

-- 1. Active Test Machinery Fleet Asset (CAT 320D Crawler Excavator)
INSERT INTO fleets (
    id, owner_id, subcategory_id, make, model, manufacture_year,
    plate_number, chassis_vin, baseline_meter_hours, fuel_type,
    self_move_capable, daily_rate_minor, excess_hour_rate_minor,
    min_rental_days, payment_timeout_hours, operator_included,
    status, location, address_text, version, created_at, updated_at
) VALUES (
    '33333333-3333-4333-8333-333333333333',
    '22222222-2222-4222-8222-222222222222',
    '019a8f00-0002-7000-8000-000000000003',
    'Caterpillar',
    '320D Hydraulic Excavator',
    2021,
    'ET-3-88412',
    'CAT320D2021ETB001',
    1250.50,
    'DIESEL',
    false,
    3500000,
    450000,
    1,
    24,
    true,
    'ACTIVE',
    ST_SetSRID(ST_MakePoint(38.7578, 8.9806), 4326),
    'Bole Airport Cargo Logistics Area, Addis Ababa',
    1,
    NOW() - INTERVAL '60 days',
    NOW()
) ON CONFLICT (id) DO UPDATE SET
    status = 'ACTIVE',
    daily_rate_minor = 3500000,
    excess_hour_rate_minor = 450000,
    baseline_meter_hours = 1250.50;

-- 2. Mandatory 6-Angle Walkaround Photos
INSERT INTO fleet_photos (id, fleet_id, angle, storage_uri, file_hash, display_order) VALUES
('33333333-bbbb-4333-8333-000000000001', '33333333-3333-4333-8333-333333333333', 'FRONT', 'https://cdn.dinkmachinery.com/fixtures/cat320d_front.jpg', 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855', 1),
('33333333-bbbb-4333-8333-000000000002', '33333333-3333-4333-8333-333333333333', 'REAR', 'https://cdn.dinkmachinery.com/fixtures/cat320d_rear.jpg', 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855', 2),
('33333333-bbbb-4333-8333-000000000003', '33333333-3333-4333-8333-333333333333', 'LEFT', 'https://cdn.dinkmachinery.com/fixtures/cat320d_left.jpg', 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855', 3),
('33333333-bbbb-4333-8333-000000000004', '33333333-3333-4333-8333-333333333333', 'RIGHT', 'https://cdn.dinkmachinery.com/fixtures/cat320d_right.jpg', 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855', 4),
('33333333-bbbb-4333-8333-000000000005', '33333333-3333-4333-8333-333333333333', 'CABIN', 'https://cdn.dinkmachinery.com/fixtures/cat320d_cabin.jpg', 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855', 5),
('33333333-bbbb-4333-8333-000000000006', '33333333-3333-4333-8333-333333333333', 'METER', 'https://cdn.dinkmachinery.com/fixtures/cat320d_meter.jpg', 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855', 6)
ON CONFLICT (id) DO NOTHING;

-- 3. Owner Operator Membership & Machine Assignment
INSERT INTO owner_operator_memberships (owner_id, operator_id, status, created_at)
VALUES (
    '22222222-2222-4222-8222-222222222222',
    '33333333-3333-4333-8333-333333333333',
    'ACTIVE',
    NOW() - INTERVAL '30 days'
) ON CONFLICT (owner_id, operator_id) DO NOTHING;

INSERT INTO fleet_operator_assignments (fleet_id, operator_id, shift_type, assigned_at)
VALUES (
    '33333333-3333-4333-8333-333333333333',
    '33333333-3333-4333-8333-333333333333',
    'FULL_TIME',
    NOW() - INTERVAL '30 days'
) ON CONFLICT (fleet_id, operator_id) DO NOTHING;
