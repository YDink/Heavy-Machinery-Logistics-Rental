-- Logistics Transport Integration, Sub-Bookings & Reverse Cascades Schema
-- Partition 05-08, SRS §05.08, ADR-001, ADR-005, ADR-007, D-05-08-001, D-05-08-002

-- 1. Extend bookings status constraint with logistics sub-machine states
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
    'DISPATCHED',
    'LOADING_INSPECTION',
    'IN_TRANSIT',
    'TRANSSHIPMENT_REQUIRED',
    'AWAITING_FLEET_READINESS',
    'PENDING_ADDITIONAL_PAYMENT',
    'READY_FOR_PICKUP',
    'ACTIVE',
    'ACTIVE_BREAKDOWN_PAUSED',
    'ACTIVE_EXTENSION_HOLD',
    'RETURN_INITIATED',
    'DAMAGE_PENDING_REVIEW',
    'DAMAGE_DISPUTED',
    'PENDING_FINAL_PAYMENT',
    'COMPLETED',
    'RATED',
    'CANCELLED',
    'CANCELLED_CARRIER_DEFAULT'
));

-- 2. Extend bookings with parent-child hierarchy and freight commission tracking
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS parent_booking_id UUID NULL REFERENCES bookings(id) ON DELETE RESTRICT;
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS booking_kind VARCHAR(24) NOT NULL DEFAULT 'MACHINERY';
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS platform_commission_rate_bps INT NOT NULL DEFAULT 1000;
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS platform_commission_minor BIGINT NOT NULL DEFAULT 0;
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS abnormal_load_escort_required BOOLEAN NOT NULL DEFAULT false;

ALTER TABLE bookings DROP CONSTRAINT IF EXISTS chk_bookings_kind;
ALTER TABLE bookings ADD CONSTRAINT chk_bookings_kind CHECK (booking_kind IN ('MACHINERY', 'LOGISTICS_SUB'));

ALTER TABLE bookings DROP CONSTRAINT IF EXISTS chk_bookings_commission;
ALTER TABLE bookings ADD CONSTRAINT chk_bookings_commission CHECK (
    platform_commission_rate_bps >= 0 AND platform_commission_rate_bps <= 10000 AND platform_commission_minor >= 0
);

CREATE INDEX IF NOT EXISTS idx_bookings_parent_hierarchy
ON bookings (parent_booking_id)
WHERE parent_booking_id IS NOT NULL;

-- 3. Extend fleet_spec_indexed with lowboy trailer dimensions and tare weight (D-05-08-002, SRS §3.3)
ALTER TABLE fleet_spec_indexed ADD COLUMN IF NOT EXISTS deck_well_length_mm NUMERIC(10, 2);
ALTER TABLE fleet_spec_indexed ADD COLUMN IF NOT EXISTS deck_width_mm NUMERIC(10, 2);
ALTER TABLE fleet_spec_indexed ADD COLUMN IF NOT EXISTS outriggers_extension_mm NUMERIC(10, 2) NOT NULL DEFAULT 0;
ALTER TABLE fleet_spec_indexed ADD COLUMN IF NOT EXISTS deck_height_ground_mm NUMERIC(10, 2);
ALTER TABLE fleet_spec_indexed ADD COLUMN IF NOT EXISTS tare_weight_kg NUMERIC(10, 2);
ALTER TABLE fleet_spec_indexed ADD COLUMN IF NOT EXISTS axle_count INT;

-- 4. Extend country_configs with GVM & corridor clearance limits (D-05-08-002, SRS §3.3, §3.4)
ALTER TABLE country_configs ADD COLUMN IF NOT EXISTS max_gvm_kg NUMERIC(10, 2) NOT NULL DEFAULT 56000.00;
ALTER TABLE country_configs ADD COLUMN IF NOT EXISTS max_corridor_height_mm INT NOT NULL DEFAULT 4500;
ALTER TABLE country_configs ADD COLUMN IF NOT EXISTS wide_load_width_mm INT NOT NULL DEFAULT 3000;

-- 5. Logistics Sub-Booking Dispatches (Pathway A: In-Platform Lowboy Sub-Booking, FL-31, §5.1)
CREATE TABLE IF NOT EXISTS logistics_dispatches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sub_booking_id UUID NOT NULL UNIQUE REFERENCES bookings(id) ON DELETE CASCADE,
    parent_booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE RESTRICT,
    logistics_fleet_id UUID NOT NULL REFERENCES fleets(id) ON DELETE RESTRICT,
    carrier_owner_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    pickup_datetime TIMESTAMPTZ NOT NULL,
    origin_latitude DOUBLE PRECISION NOT NULL,
    origin_longitude DOUBLE PRECISION NOT NULL,
    origin_depot_name TEXT NOT NULL,
    destination_latitude DOUBLE PRECISION NOT NULL,
    destination_longitude DOUBLE PRECISION NOT NULL,
    destination_site_name TEXT NOT NULL,
    freight_hire_fee_minor BIGINT NOT NULL,
    platform_commission_bps INT NOT NULL DEFAULT 250,
    commission_owed_minor BIGINT NOT NULL,
    abnormal_load_escort_required BOOLEAN NOT NULL DEFAULT false,
    escort_fee_minor BIGINT NOT NULL DEFAULT 0,
    winching_photos_json JSONB NOT NULL DEFAULT '[]'::jsonb,
    dispatched_at TIMESTAMPTZ,
    loaded_at TIMESTAMPTZ,
    in_transit_at TIMESTAMPTZ,
    arrived_at TIMESTAMPTZ,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_logistics_dispatches_fees CHECK (
        freight_hire_fee_minor > 0 AND
        platform_commission_bps >= 0 AND
        commission_owed_minor >= 0 AND
        escort_fee_minor >= 0
    )
);

CREATE INDEX IF NOT EXISTS idx_logistics_dispatches_parent ON logistics_dispatches(parent_booking_id);
CREATE INDEX IF NOT EXISTS idx_logistics_dispatches_carrier ON logistics_dispatches(carrier_owner_id, created_at DESC);

CREATE TRIGGER trg_logistics_dispatches_updated_at
BEFORE UPDATE ON logistics_dispatches
FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();

-- 6. External Private Transport Declarations (Pathway B, SRS §4.3)
CREATE TABLE IF NOT EXISTS external_transports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_booking_id UUID NOT NULL UNIQUE REFERENCES bookings(id) ON DELETE CASCADE,
    registered_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    carrier_name VARCHAR(120) NOT NULL,
    driver_full_name VARCHAR(120) NOT NULL,
    driver_phone_number VARCHAR(24) NOT NULL,
    prime_mover_plate_number VARCHAR(32) NOT NULL,
    estimated_arrival_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_external_transports_parent ON external_transports(parent_booking_id);

CREATE TRIGGER trg_external_transports_updated_at
BEFORE UPDATE ON external_transports
FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();

-- 7. Mid-Transit Breakdown & Emergency Roadside Transshipment (SRS §5.3)
CREATE TABLE IF NOT EXISTS logistics_transshipments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sub_booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE RESTRICT,
    parent_booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE RESTRICT,
    reported_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    breakdown_latitude DOUBLE PRECISION NOT NULL,
    breakdown_longitude DOUBLE PRECISION NOT NULL,
    failure_category VARCHAR(64) NOT NULL,
    failure_description TEXT NOT NULL,
    photos_json JSONB NOT NULL DEFAULT '[]'::jsonb,
    hot_swap_sla_deadline TIMESTAMPTZ NOT NULL,
    transshipment_bounty_minor BIGINT NOT NULL,
    replacement_fleet_id UUID REFERENCES fleets(id) ON DELETE SET NULL,
    replacement_carrier_id UUID REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'OPEN',
    accepted_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_transshipment_failure_category CHECK (failure_category IN (
        'PRIME_MOVER_ENGINE_BLOWN',
        'TRAILER_AXLE_FAILURE',
        'HYDRAULIC_RAMP_FAILURE',
        'CORRIDOR_BLOCKADE',
        'OTHER'
    )),
    CONSTRAINT chk_transshipment_status CHECK (status IN (
        'OPEN', 'ACCEPTED', 'COMPLETED', 'EXPIRED', 'FAILED'
    )),
    CONSTRAINT chk_transshipment_bounty CHECK (transshipment_bounty_minor >= 0)
);

CREATE INDEX IF NOT EXISTS idx_logistics_transshipments_sub ON logistics_transshipments(sub_booking_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_logistics_transshipments_parent ON logistics_transshipments(parent_booking_id, created_at DESC);

CREATE TRIGGER trg_logistics_transshipments_updated_at
BEFORE UPDATE ON logistics_transshipments
FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();

-- 8. Automated Reverse Failure Cascade & Force Majeure Exoneration Ledger (SRS §6, FLOW-18)
CREATE TABLE IF NOT EXISTS logistics_reverse_cascades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_booking_id UUID NOT NULL UNIQUE REFERENCES bookings(id) ON DELETE RESTRICT,
    sub_booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE RESTRICT,
    triggered_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    resolution VARCHAR(64) NOT NULL DEFAULT 'FORCE_MAJEURE_CANCELLATION',
    reason TEXT NOT NULL,
    renter_refund_due_minor BIGINT NOT NULL,
    owner_penalty_levied_minor BIGINT NOT NULL DEFAULT 0,
    owner_commission_restitution_minor BIGINT NOT NULL,
    carrier_default_penalty_minor BIGINT NOT NULL DEFAULT 0,
    carrier_rating_deduction NUMERIC(3, 2) NOT NULL DEFAULT 0.50,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_reverse_cascade_amounts CHECK (
        renter_refund_due_minor >= 0 AND
        owner_penalty_levied_minor >= 0 AND
        owner_commission_restitution_minor >= 0 AND
        carrier_default_penalty_minor >= 0
    )
);

-- 9. Dynamic RBAC Permissions for Logistics Transport Integration (ADR-018)
INSERT INTO permissions (code, resource, action, description) VALUES
('logistics:match', 'logistics', 'match', 'Query dimensionally compatible lowboy carriers for a confirmed booking'),
('logistics:sub_booking', 'logistics', 'sub_booking', 'Instantiate linked in-platform logistics sub-booking'),
('logistics:external_transport', 'logistics', 'external_transport', 'Declare external private transport mobilization details'),
('logistics:dispatch', 'logistics', 'dispatch', 'Execute lowboy dispatch, winching inspection, and transit milestones'),
('logistics:transshipment', 'logistics', 'transshipment', 'Report roadside breakdown and coordinate emergency transshipment'),
('logistics:reverse_cascade', 'logistics', 'reverse_cascade', 'Trigger automated Force Majeure reverse failure cascade')
ON CONFLICT (code) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.code IN ('super_admin', 'admin')
  AND p.code IN (
      'logistics:match',
      'logistics:sub_booking',
      'logistics:external_transport',
      'logistics:dispatch',
      'logistics:transshipment',
      'logistics:reverse_cascade'
  )
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.code IN ('owner', 'fleet_manager')
  AND p.code IN (
      'logistics:match',
      'logistics:sub_booking',
      'logistics:external_transport',
      'logistics:dispatch',
      'logistics:transshipment'
  )
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.code = 'operator'
  AND p.code IN ('logistics:dispatch', 'logistics:transshipment')
ON CONFLICT DO NOTHING;
