-- Rental Completion, Demobilization, Damage Triage & Settlement Reconciliation Schema
-- Partition 05-07, SRS §05.07, ADR-004, ADR-010, D-05-07-001, D-05-07-002

-- 1. Update bookings status constraint to include completion and damage adjudication sub-states
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
    'DAMAGE_PENDING_REVIEW',
    'DAMAGE_DISPUTED',
    'PENDING_FINAL_PAYMENT',
    'COMPLETED',
    'RATED',
    'CANCELLED'
));

-- Add completion, damage claim, and settlement summary tracking to bookings
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS damage_claim_submitted BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS damage_claim_waived BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS accepted_damage_minor BIGINT NOT NULL DEFAULT 0;
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS excess_hours_surcharge_minor BIGINT NOT NULL DEFAULT 0;
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS final_net_balance_minor BIGINT NOT NULL DEFAULT 0;

-- 2. Return Reticle Inspections Entity (RC-02, D-05-07-001)
CREATE TABLE IF NOT EXISTS return_inspections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE RESTRICT,
    inspector_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    check_out_meter_hours NUMERIC(10, 2) NOT NULL,
    fluid_levels_ok BOOLEAN NOT NULL DEFAULT true,
    safety_glass_intact BOOLEAN NOT NULL DEFAULT true,
    accelerometer_steady BOOLEAN NOT NULL DEFAULT true,
    notes TEXT,
    photos_json JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_return_inspections_meter CHECK (check_out_meter_hours >= 0)
);

CREATE INDEX IF NOT EXISTS idx_return_inspections_booking ON return_inspections(booking_id, created_at DESC);

-- 3. Damage Claim Dockets & Arbitration (RC-04, D-05-07-002)
CREATE TABLE IF NOT EXISTS damage_claim_dockets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE RESTRICT,
    claimant_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    component_category VARCHAR(64) NOT NULL,
    incident_narrative TEXT NOT NULL,
    repair_quote_uri VARCHAR(512),
    claimed_amount_minor BIGINT NOT NULL,
    accepted_amount_minor BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'ETB',
    damage_photos_json JSONB NOT NULL DEFAULT '[]',
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING_REVIEW',
    dispute_reason TEXT,
    contractor_responded_at TIMESTAMPTZ,
    arbitration_notes TEXT,
    arbitrated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    arbitrated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_damage_claims_amount CHECK (claimed_amount_minor >= 0 AND accepted_amount_minor >= 0),
    CONSTRAINT chk_damage_claims_category CHECK (component_category IN (
        'GET_IMPLEMENT_TIPS',
        'UNDERCARRIAGE_TRACKS',
        'HYDRAULIC_CIRCUIT',
        'STRUCTURAL_BOOM_ARM_CHASSIS',
        'OPERATOR_CABIN',
        'ENGINE_POWERTRAIN',
        'OTHER'
    )),
    CONSTRAINT chk_damage_claims_status CHECK (status IN (
        'PENDING_REVIEW',
        'ACCEPTED',
        'DISPUTED',
        'RESOLVED_FULL',
        'RESOLVED_ADJUSTED',
        'DISMISSED',
        'WAIVED'
    ))
);

CREATE INDEX IF NOT EXISTS idx_damage_claim_dockets_booking ON damage_claim_dockets(booking_id, created_at DESC);

-- 4. Final Balance Reconciliation Settlements (RC-05, RC-06)
CREATE TABLE IF NOT EXISTS rental_settlements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL UNIQUE REFERENCES bookings(id) ON DELETE RESTRICT,
    base_machine_fee_minor BIGINT NOT NULL,
    bundled_attachments_fee_minor BIGINT NOT NULL DEFAULT 0,
    contractual_allowed_hours NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    actual_operating_hours NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    credited_breakdown_hours NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    billable_excess_hours NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    excess_hourly_rate_minor BIGINT NOT NULL DEFAULT 0,
    excess_hours_surcharge_minor BIGINT NOT NULL DEFAULT 0,
    accepted_damage_minor BIGINT NOT NULL DEFAULT 0,
    gross_liabilities_minor BIGINT NOT NULL,
    prepayment_settled_minor BIGINT NOT NULL,
    net_balance_minor BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ETB',
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING_PAYMENT',
    final_slip_storage_uri VARCHAR(512),
    final_slip_file_hash VARCHAR(64),
    settled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_rental_settlements_status CHECK (status IN (
        'PENDING_PAYMENT',
        'PAYMENT_EVIDENCE_UPLOADED',
        'SETTLED',
        'REFUND_DUE',
        'REFUND_CONFIRMED',
        'ZERO_BALANCE_AUTO_COMPLETED'
    ))
);

CREATE INDEX IF NOT EXISTS idx_rental_settlements_booking ON rental_settlements(booking_id);

-- 5. Mutual Double-Blind Review Engine (RC-07)
CREATE TABLE IF NOT EXISTS booking_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE RESTRICT,
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    target_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    role VARCHAR(16) NOT NULL,
    rating NUMERIC(2, 1) NOT NULL,
    feedback TEXT,
    dimension_ratings_json JSONB NOT NULL DEFAULT '{}',
    is_published BOOLEAN NOT NULL DEFAULT false,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_booking_reviews_rating CHECK (rating >= 1.0 AND rating <= 5.0),
    CONSTRAINT chk_booking_reviews_role CHECK (role IN ('RENTER', 'OWNER', 'OPERATOR')),
    CONSTRAINT uq_booking_reviews_author UNIQUE (booking_id, author_id)
);

CREATE INDEX IF NOT EXISTS idx_booking_reviews_target ON booking_reviews(target_id, is_published, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_booking_reviews_booking ON booking_reviews(booking_id);

-- 6. Seed Dynamic RBAC Permissions
INSERT INTO permissions (code, resource, action, description) VALUES
    ('rentals:return_initiate', 'rentals', 'return_initiate', 'Initiate rental demobilization and stop shift telemetry'),
    ('rentals:return_inspect', 'rentals', 'return_inspect', 'Execute 6-angle return reticle inspection and log checkout meter'),
    ('rentals:damage_claim', 'rentals', 'damage_claim', 'Lodge compensable structural damage claim within 24h statutory window'),
    ('rentals:damage_respond', 'rentals', 'damage_respond', 'Accept or dispute filed damage claim liability'),
    ('rentals:arbitrate', 'rentals', 'arbitrate', 'Adjudicate and issue binding administrative rulings on contested damage claims'),
    ('rentals:settle', 'rentals', 'settle', 'Generate, confirm, or verify final balance reconciliation statements'),
    ('rentals:review', 'rentals', 'review', 'Submit double-blind mutual trust and reputation scorecard')
ON CONFLICT (code) DO NOTHING;

-- Map permissions to roles
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.code = 'renter' AND p.code IN ('rentals:return_initiate', 'rentals:damage_respond', 'rentals:settle', 'rentals:review')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.code = 'owner' AND p.code IN ('rentals:return_inspect', 'rentals:damage_claim', 'rentals:settle', 'rentals:review')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.code = 'operator' AND p.code IN ('rentals:return_inspect')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.code = 'admin' AND p.code IN (
    'rentals:return_initiate',
    'rentals:return_inspect',
    'rentals:damage_claim',
    'rentals:damage_respond',
    'rentals:arbitrate',
    'rentals:settle',
    'rentals:review'
)
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.code = 'super_admin' AND p.code IN (
    'rentals:return_initiate',
    'rentals:return_inspect',
    'rentals:damage_claim',
    'rentals:damage_respond',
    'rentals:arbitrate',
    'rentals:settle',
    'rentals:review'
)
ON CONFLICT DO NOTHING;
