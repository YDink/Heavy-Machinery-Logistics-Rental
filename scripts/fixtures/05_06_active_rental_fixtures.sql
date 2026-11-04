-- ==============================================================================
-- 05_06_active_rental_fixtures.sql
-- Partition 05-06: Active Rental Operations, Handover Inspection, Shift Telemetry & Breakdowns
-- Strictly Idempotent: ON CONFLICT DO UPDATE / DO NOTHING
-- Static Deterministic UUIDs
-- ==============================================================================

-- 1. Verified Handover Physical Inspection (AR-01)
INSERT INTO handover_inspections (
    id, booking_id, inspector_id, check_in_meter_hours,
    fluid_levels_ok, safety_glass_intact, notes, photos_json, created_at
) VALUES (
    '66666666-6666-4666-8666-000000000001',
    '55555555-5555-4555-8555-555555555555',
    '22222222-2222-4222-8222-222222222222',
    1250.50,
    true,
    true,
    'Machine baseline inspected: all hydraulic hoses intact, oil levels at nominal capacity.',
    '[{"angle": "FRONT", "storageUri": "https://cdn.dinkmachinery.com/fixtures/inspections/cat320d_checkin_front.jpg"}, {"angle": "METER", "storageUri": "https://cdn.dinkmachinery.com/fixtures/inspections/cat320d_checkin_meter.jpg"}]'::jsonb,
    NOW() - INTERVAL '2 days'
) ON CONFLICT (id) DO NOTHING;

-- 2. Verified Handover Cryptographic Challenge PIN (AR-02)
INSERT INTO handover_challenges (
    id, booking_id, challenge_pin, expires_at, attempt_count,
    max_attempts, is_verified, verified_at, touch_signature_base64,
    contractor_lat, contractor_lng, created_at
) VALUES (
    '66666666-6666-4666-8666-000000000002',
    '55555555-5555-4555-8555-555555555555',
    '842913',
    NOW() + INTERVAL '10 minutes',
    1,
    3,
    true,
    NOW() - INTERVAL '2 days',
    'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==',
    8.9806,
    38.7578,
    NOW() - INTERVAL '2 days'
) ON CONFLICT (id) DO NOTHING;

-- 3. Completed Past Shift Telemetry (Day 1)
INSERT INTO rental_shifts (
    id, booking_id, operator_id, shift_date, shift_start_time, shift_end_time,
    accumulated_engine_hours, excess_hours, start_meter_hours, end_meter_hours,
    is_paused, created_at, updated_at
) VALUES (
    '66666666-6666-4666-8666-000000000003',
    '55555555-5555-4555-8555-555555555555',
    '33333333-3333-4333-8333-333333333333',
    CURRENT_DATE - INTERVAL '1 day',
    NOW() - INTERVAL '26 hours',
    NOW() - INTERVAL '16 hours',
    9.50,
    1.50,
    1250.50,
    1260.00,
    false,
    NOW() - INTERVAL '26 hours',
    NOW() - INTERVAL '16 hours'
) ON CONFLICT (id) DO NOTHING;

-- 4. Active Shift Telemetry (Current Shift In Progress)
INSERT INTO rental_shifts (
    id, booking_id, operator_id, shift_date, shift_start_time, shift_end_time,
    accumulated_engine_hours, excess_hours, start_meter_hours, end_meter_hours,
    is_paused, created_at, updated_at
) VALUES (
    '66666666-6666-4666-8666-000000000004',
    '55555555-5555-4555-8555-555555555555',
    '33333333-3333-4333-8333-333333333333',
    CURRENT_DATE,
    NOW() - INTERVAL '4 hours',
    NULL,
    3.75,
    0.00,
    1260.00,
    NULL,
    false,
    NOW() - INTERVAL '4 hours',
    NOW()
) ON CONFLICT (id) DO UPDATE SET
    accumulated_engine_hours = 3.75;
