-- Active Rental Management, Shift Telemetry, Breakdown & Extension Schema (Partition 05-06, SRS §05.06, ADR-006, ADR-012, D-05-06-001, D-05-06-002)

-- 1. Update bookings status constraint to include active operational sub-states
ALTER TABLE bookings DROP CONSTRAINT IF EXISTS chk_bookings_status;
ALTER TABLE bookings ADD CONSTRAINT chk_bookings_status CHECK (status IN (
    'DRAFT',
    'PENDING_PAYMENT',
    'PAYMENT_EVIDENCE_UPLOADED',
    'PAYMENT_DECLINED',
    'DISPUTE_OPENED',
    'DISPUTE_RESOLVED_RENTER',
    'DISPUTE_RESOLVED_OWNER',
    'CONFIRMED',
    'AWAITING_FLEET_READINESS',
    'PENDING_ADDITIONAL_PAYMENT',
    'READY_FOR_PICKUP',
    'ACTIVE',
    'ACTIVE_BREAKDOWN_PAUSED',
    'ACTIVE_EXTENSION_HOLD',
    'RETURN_INITIATED',
    'PENDING_FINAL_PAYMENT',
    'COMPLETED',
    'RATED',
    'CANCELLED'
));

-- Add assigned operator column to bookings
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS assigned_operator_id UUID REFERENCES users(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_bookings_assigned_operator ON bookings(assigned_operator_id) WHERE assigned_operator_id IS NOT NULL;

-- 2. Handover Inspections Entity (AR-01)
CREATE TABLE IF NOT EXISTS handover_inspections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE RESTRICT,
    inspector_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    check_in_meter_hours NUMERIC(10, 2) NOT NULL,
    fluid_levels_ok BOOLEAN NOT NULL DEFAULT true,
    safety_glass_intact BOOLEAN NOT NULL DEFAULT true,
    notes TEXT,
    photos_json JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_handover_inspections_meter CHECK (check_in_meter_hours >= 0)
);

CREATE INDEX IF NOT EXISTS idx_handover_inspections_booking ON handover_inspections(booking_id, created_at DESC);

-- 3. Handover Cryptographic Challenge PINs (AR-02)
CREATE TABLE IF NOT EXISTS handover_challenges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE RESTRICT,
    challenge_pin VARCHAR(6) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    attempt_count INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 3,
    is_verified BOOLEAN NOT NULL DEFAULT false,
    verified_at TIMESTAMPTZ,
    touch_signature_base64 TEXT,
    contractor_lat DOUBLE PRECISION,
    contractor_lng DOUBLE PRECISION,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_handover_challenges_booking ON handover_challenges(booking_id, created_at DESC);

-- 4. Daily Shift Telemetry Tracking (AR-03, 8.0 hrs/day baseline)
CREATE TABLE IF NOT EXISTS rental_shifts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE RESTRICT,
    operator_id UUID REFERENCES users(id) ON DELETE SET NULL,
    shift_date DATE NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    start_meter_hours NUMERIC(10, 2) NOT NULL,
    end_meter_hours NUMERIC(10, 2),
    allocated_hours NUMERIC(6, 2) NOT NULL DEFAULT 8.00,
    accumulated_hours NUMERIC(6, 2) NOT NULL DEFAULT 0.00,
    excess_hours NUMERIC(6, 2) NOT NULL DEFAULT 0.00,
    is_paused BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_rental_shifts_meter CHECK (end_meter_hours IS NULL OR end_meter_hours >= start_meter_hours)
);

CREATE INDEX IF NOT EXISTS idx_rental_shifts_booking_date ON rental_shifts(booking_id, shift_date DESC);

-- 5. Rental Duration Extensions (AR-04)
CREATE TABLE IF NOT EXISTS rental_extensions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE RESTRICT,
    requested_by_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    original_end_date DATE NOT NULL,
    new_end_date DATE NOT NULL,
    extension_days INT NOT NULL,
    incremental_gross_minor BIGINT NOT NULL,
    incremental_prepayment_minor BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ETB',
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING_OWNER',
    owner_decided_at TIMESTAMPTZ,
    owner_decision_notes TEXT,
    slip_storage_uri VARCHAR(512),
    slip_file_hash VARCHAR(64),
    confirmed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_rental_extensions_dates CHECK (new_end_date > original_end_date AND extension_days > 0),
    CONSTRAINT chk_rental_extensions_amounts CHECK (incremental_gross_minor > 0 AND incremental_prepayment_minor > 0),
    CONSTRAINT chk_rental_extensions_status CHECK (status IN (
        'PENDING_OWNER',
        'OWNER_APPROVED',
        'OWNER_DECLINED',
        'PAYMENT_EVIDENCE_UPLOADED',
        'CONFIRMED',
        'EXPIRED',
        'CANCELLED'
    ))
);

CREATE INDEX IF NOT EXISTS idx_rental_extensions_booking ON rental_extensions(booking_id, created_at DESC);

-- 6. Mechanical Breakdowns & Stoppage Logs (AR-05)
CREATE TABLE IF NOT EXISTS mechanical_breakdowns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE RESTRICT,
    reported_by_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    fault_category VARCHAR(64) NOT NULL,
    fault_description TEXT NOT NULL,
    diagnostic_trouble_codes TEXT[] NOT NULL DEFAULT '{}',
    current_hour_meter NUMERIC(10, 2) NOT NULL,
    incident_photos_json JSONB NOT NULL DEFAULT '[]',
    status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE_BREAKDOWN_PAUSED',
    sla_deadline TIMESTAMPTZ NOT NULL,
    mechanic_name VARCHAR(128),
    mechanic_phone VARCHAR(32),
    resolution_notes TEXT,
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_mechanical_breakdowns_category CHECK (fault_category IN (
        'ENGINE',
        'HYDRAULIC_CIRCUIT',
        'ELECTRICAL',
        'UNDERCARRIAGE',
        'STRUCTURAL',
        'BOOM_BUCKET',
        'OTHER'
    )),
    CONSTRAINT chk_mechanical_breakdowns_status CHECK (status IN (
        'ACTIVE_BREAKDOWN_PAUSED',
        'MECHANIC_DISPATCHED',
        'REPAIRED',
        'RESOLVED'
    ))
);

CREATE INDEX IF NOT EXISTS idx_mechanical_breakdowns_booking ON mechanical_breakdowns(booking_id, created_at DESC);

-- 7. Seed Dynamic RBAC Permissions
INSERT INTO permissions (name, description) VALUES
    ('rentals:inspect', 'Submit physical handover reticle inspection and hour meter reading'),
    ('rentals:confirm', 'Confirm dual challenge PIN and sign handover agreement'),
    ('rentals:hud', 'View active rental HUD telemetry and shift progression'),
    ('rentals:extend', 'Request or approve rental duration extensions'),
    ('rentals:breakdown', 'Report mechanical breakdowns and pause shift timers'),
    ('rentals:resolve', 'Certify mechanical repairs and resume rental shifts')
ON CONFLICT (name) DO NOTHING;

-- Map permissions to roles
INSERT INTO role_permissions (role, permission_id)
SELECT 'renter', id FROM permissions WHERE name IN ('rentals:confirm', 'rentals:hud', 'rentals:extend', 'rentals:breakdown')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role, permission_id)
SELECT 'owner', id FROM permissions WHERE name IN ('rentals:inspect', 'rentals:hud', 'rentals:extend', 'rentals:breakdown', 'rentals:resolve')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role, permission_id)
SELECT 'operator', id FROM permissions WHERE name IN ('rentals:inspect', 'rentals:hud', 'rentals:breakdown')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role, permission_id)
SELECT 'admin', id FROM permissions WHERE name LIKE 'rentals:%'
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role, permission_id)
SELECT 'super_admin', id FROM permissions WHERE name LIKE 'rentals:%'
ON CONFLICT DO NOTHING;
