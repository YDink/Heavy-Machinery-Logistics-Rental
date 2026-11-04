-- ==============================================================================
-- 05_04_discovery_fixtures.sql
-- Partition 05-04: Fleet Discovery, Owner Performance Metrics & Insurance
-- Strictly Idempotent: ON CONFLICT DO UPDATE / DO NOTHING
-- Static Deterministic UUIDs
-- ==============================================================================

-- 1. Owner Aggregate Performance Metrics
INSERT INTO owner_aggregates (
    owner_id, average_rating, total_rentals_count, total_completed_rentals, completion_rate, is_verified, updated_at
) VALUES (
    '22222222-2222-4222-8222-222222222222',
    4.95,
    42,
    41,
    0.9762,
    true,
    NOW()
) ON CONFLICT (owner_id) DO UPDATE SET
    average_rating = 4.95,
    is_verified = true;

-- 2. Comprehensive Insurance Policy for CAT 320D
INSERT INTO fleet_insurances (
    id, fleet_id, policy_number, provider_name, coverage_type, status, expires_at, certificate_uri, created_at, updated_at
) VALUES (
    '33333333-eeee-4333-8333-000000000001',
    '33333333-3333-4333-8333-333333333333',
    'POL-NYALA-2026-9921',
    'Nyala Insurance S.C.',
    'COMPREHENSIVE_PLANT_ALL_RISKS',
    'ACTIVE',
    CURRENT_DATE + INTERVAL '180 days',
    'https://cdn.dinkmachinery.com/insurance/cat320d_cert.pdf',
    NOW() - INTERVAL '30 days',
    NOW()
) ON CONFLICT (id) DO UPDATE SET
    status = 'ACTIVE',
    expires_at = CURRENT_DATE + INTERVAL '180 days';
