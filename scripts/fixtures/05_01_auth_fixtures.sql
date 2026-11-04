-- ==============================================================================
-- 05_01_auth_fixtures.sql
-- Partition 05-01: Authentication, Legal Terms & RBAC User Role Assignments
-- Strictly Idempotent: ON CONFLICT DO UPDATE / DO NOTHING
-- Static Deterministic UUIDs
-- ==============================================================================

-- 1. Ensure Roles Are Assigned to Test Users
INSERT INTO user_roles (user_id, role_id, assigned_at)
SELECT '00000000-0000-4000-8000-000000000001', id, NOW()
FROM roles WHERE code IN ('super_admin', 'admin')
ON CONFLICT (user_id, role_id) DO NOTHING;

INSERT INTO user_roles (user_id, role_id, assigned_at)
SELECT '11111111-1111-4111-8111-111111111111', id, NOW()
FROM roles WHERE code = 'renter'
ON CONFLICT (user_id, role_id) DO NOTHING;

INSERT INTO user_roles (user_id, role_id, assigned_at)
SELECT '22222222-2222-4222-8222-222222222222', id, NOW()
FROM roles WHERE code = 'owner'
ON CONFLICT (user_id, role_id) DO NOTHING;

INSERT INTO user_roles (user_id, role_id, assigned_at)
SELECT '33333333-3333-4333-8333-333333333333', id, NOW()
FROM roles WHERE code = 'operator'
ON CONFLICT (user_id, role_id) DO NOTHING;

-- 2. Statutory Legal Agreement Acceptance Records
INSERT INTO user_legal_agreements (
    id, user_id, terms_version, privacy_version, accepted_at, client_ip, client_user_agent
) VALUES (
    '11111111-aaaa-4111-8111-111111111111',
    '11111111-1111-4111-8111-111111111111',
    'v1.0.0',
    'v1.0.0',
    NOW() - INTERVAL '30 days',
    '127.0.0.1',
    'DinkMobile/1.0.0 (Android 14; Mobile)'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO user_legal_agreements (
    id, user_id, terms_version, privacy_version, accepted_at, client_ip, client_user_agent
) VALUES (
    '22222222-aaaa-4222-8222-222222222222',
    '22222222-2222-4222-8222-222222222222',
    'v1.0.0',
    'v1.0.0',
    NOW() - INTERVAL '30 days',
    '127.0.0.1',
    'DinkMobile/1.0.0 (iOS 18.2; Mobile)'
) ON CONFLICT (id) DO NOTHING;
