-- ==============================================================================
-- 05_07_rental_completion_fixtures.sql
-- Partition 05-07: Rental Completion, Off-Hire Return Inspection, Damage & Settlement
-- Strictly Idempotent: ON CONFLICT DO UPDATE / DO NOTHING
-- Static Deterministic UUIDs
-- ==============================================================================

-- 1. Verified Return Reticle Physical Inspection (RC-02)
INSERT INTO return_inspections (
    id, booking_id, inspector_id, check_out_meter_hours,
    fluid_levels_ok, safety_glass_intact, accelerometer_steady,
    notes, photos_json, created_at
) VALUES (
    '77777777-7777-4777-8777-000000000001',
    '55555555-5555-4555-8555-555555555555',
    '33333333-3333-4333-8333-333333333333',
    1312.20,
    true,
    true,
    true,
    'Return inspection completed: cabin hour meter checked at 1312.20h, implement intact, fluid levels verified.',
    '[
        {"angle": "FRONT_IMPLEMENT", "storageUri": "https://cdn.dinkmachinery.com/fixtures/inspections/cat320d_return_front.jpg"},
        {"angle": "REAR_COUNTERWEIGHT_ENGINE", "storageUri": "https://cdn.dinkmachinery.com/fixtures/inspections/cat320d_return_rear.jpg"},
        {"angle": "LEFT_SIDE_TRACKS", "storageUri": "https://cdn.dinkmachinery.com/fixtures/inspections/cat320d_return_left.jpg"},
        {"angle": "RIGHT_SIDE_TRACKS", "storageUri": "https://cdn.dinkmachinery.com/fixtures/inspections/cat320d_return_right.jpg"},
        {"angle": "HOUR_METER", "storageUri": "https://cdn.dinkmachinery.com/fixtures/inspections/cat320d_return_meter.jpg"},
        {"angle": "CABIN_INTERIOR", "storageUri": "https://cdn.dinkmachinery.com/fixtures/inspections/cat320d_return_cabin.jpg"}
    ]'::jsonb,
    NOW() - INTERVAL '6 hours'
) ON CONFLICT (id) DO NOTHING;

-- 2. Accepted Damage Claim Docket (RC-04, DSP-01)
INSERT INTO damage_claim_dockets (
    id, booking_id, claimant_id, component_category,
    incident_narrative, repair_quote_uri, claimed_amount_minor,
    accepted_amount_minor, currency, damage_photos_json, status,
    created_at, updated_at
) VALUES (
    '77777777-7777-4777-8777-000000000002',
    '55555555-5555-4555-8555-555555555555',
    '22222222-2222-4222-8222-222222222222',
    'OPERATOR_CABIN',
    'Side cabin window sustained spiderweb cracking during trench stone excavation.',
    'https://cdn.dinkmachinery.com/fixtures/quotes/cat320d_glass_quote.pdf',
    850000,
    850000,
    'ETB',
    '[{"angle": "CAB_GLASS", "storageUri": "https://cdn.dinkmachinery.com/fixtures/damage/cat320d_cab_crack.jpg"}]'::jsonb,
    'ACCEPTED',
    NOW() - INTERVAL '5 hours',
    NOW() - INTERVAL '3 hours'
) ON CONFLICT (id) DO NOTHING;

-- 3. Final Balance Reconciliation Settlement Statement (RC-05)
INSERT INTO rental_settlements (
    id, booking_id, base_machine_fee_minor, bundled_attachments_fee_minor,
    contractual_allowed_hours, actual_operating_hours, credited_breakdown_hours,
    billable_excess_hours, excess_hourly_rate_minor, excess_hours_surcharge_minor,
    accepted_damage_minor, gross_liabilities_minor, prepayment_settled_minor,
    net_balance_minor, currency, status, created_at, updated_at
) VALUES (
    '77777777-7777-4777-8777-000000000003',
    '55555555-5555-4555-8555-555555555555',
    15750000,
    5600000,
    56.00,
    61.70,
    0.00,
    5.70,
    280000,
    1596000,
    850000,
    23796000,
    10675000,
    13121000,
    'ETB',
    'PENDING_PAYMENT',
    NOW() - INTERVAL '2 hours',
    NOW() - INTERVAL '2 hours'
) ON CONFLICT (booking_id) DO UPDATE SET
    net_balance_minor = 13121000;

-- 4. Mutual Double-Blind Reputation Scorecards (RC-07)
-- Contractor rating of Fleet Owner
INSERT INTO booking_reviews (
    id, booking_id, author_id, target_id, role,
    rating, feedback, dimension_ratings_json, is_published, created_at
) VALUES (
    '77777777-7777-4777-8777-000000000004',
    '55555555-5555-4555-8555-555555555555',
    '11111111-1111-4111-8111-111111111111',
    '22222222-2222-4222-8222-222222222222',
    'RENTER',
    4.5,
    'Machine was in excellent working order with responsive hydraulic pressure. Prompt handover coordination.',
    '{"machine_roadworthiness": 5.0, "catalog_accuracy": 4.5, "operator_support": 4.0, "settlement_fairness": 4.5}'::jsonb,
    false,
    NOW() - INTERVAL '1 hour'
) ON CONFLICT (booking_id, author_id) DO NOTHING;

-- Owner rating of Contractor
INSERT INTO booking_reviews (
    id, booking_id, author_id, target_id, role,
    rating, feedback, dimension_ratings_json, is_published, created_at
) VALUES (
    '77777777-7777-4777-8777-000000000005',
    '55555555-5555-4555-8555-555555555555',
    '22222222-2222-4222-8222-222222222222',
    '11111111-1111-4111-8111-111111111111',
    'OWNER',
    5.0,
    'Punctual settlement, safe site operations, and clear communication throughout the hire duration.',
    '{"financial_promptness": 5.0, "equipment_care": 4.8, "site_support": 5.0, "communication": 5.0}'::jsonb,
    false,
    NOW() - INTERVAL '30 minutes'
) ON CONFLICT (booking_id, author_id) DO NOTHING;
