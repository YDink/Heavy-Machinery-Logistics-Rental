-- ==============================================================================
-- 05_05_booking_fixtures.sql
-- Partition 05-05: Real-Time Availability & Booking Lifecycle State Machine
-- Strictly Idempotent: ON CONFLICT DO UPDATE / DO NOTHING
-- Static Deterministic UUIDs
-- ==============================================================================

-- 1. Active Booking (Confirmed, Handed Over & Operational)
INSERT INTO bookings (
    id, booking_reference, renter_id, owner_id, fleet_id,
    start_date, end_date, hire_period, duration_days,
    daily_rate_applied_minor, base_hire_fee_minor, attachments_fee_minor,
    gross_rental_amount_minor, prepayment_amount_minor, currency,
    site_location_address, status, handover_pin, check_in_meter_hours,
    confirmed_at, activated_at, version, created_at, updated_at
) VALUES (
    '55555555-5555-4555-8555-555555555555',
    'BK-2026-1104-001',
    '11111111-1111-4111-8111-111111111111',
    '22222222-2222-4222-8222-222222222222',
    '33333333-3333-4333-8333-333333333333',
    CURRENT_DATE - INTERVAL '2 days',
    CURRENT_DATE + INTERVAL '5 days',
    daterange(CURRENT_DATE - INTERVAL '2 days', CURRENT_DATE + INTERVAL '5 days', '[)'),
    7,
    3500000,
    24500000,
    4550000,
    29050000,
    29050000,
    'ETB',
    'Addis Ababa Bole Lemi National Industrial Park Phase II',
    'ACTIVE',
    '842913',
    1250.50,
    NOW() - INTERVAL '3 days',
    NOW() - INTERVAL '2 days',
    1,
    NOW() - INTERVAL '4 days',
    NOW()
) ON CONFLICT (id) DO UPDATE SET
    status = 'ACTIVE',
    check_in_meter_hours = 1250.50,
    version = 1;

-- 2. Pending Payment Booking (Awaiting Bank Slip Upload / Hold Verification)
INSERT INTO bookings (
    id, booking_reference, renter_id, owner_id, fleet_id,
    start_date, end_date, hire_period, duration_days,
    daily_rate_applied_minor, base_hire_fee_minor, attachments_fee_minor,
    gross_rental_amount_minor, prepayment_amount_minor, currency,
    site_location_address, status, expires_at,
    version, created_at, updated_at
) VALUES (
    '55555555-5555-4555-8555-555555555556',
    'BK-2026-1104-002',
    '11111111-1111-4111-8111-111111111111',
    '22222222-2222-4222-8222-222222222222',
    '33333333-3333-4333-8333-333333333333',
    CURRENT_DATE + INTERVAL '10 days',
    CURRENT_DATE + INTERVAL '14 days',
    daterange(CURRENT_DATE + INTERVAL '10 days', CURRENT_DATE + INTERVAL '14 days', '[)'),
    4,
    3500000,
    14000000,
    0,
    14000000,
    14000000,
    'ETB',
    'Akaki-Kality Sub-City Highway Expansion Site',
    'PENDING_PAYMENT',
    NOW() + INTERVAL '18 hours',
    1,
    NOW() - INTERVAL '1 hour',
    NOW()
) ON CONFLICT (id) DO UPDATE SET
    status = 'PENDING_PAYMENT',
    version = 1;
