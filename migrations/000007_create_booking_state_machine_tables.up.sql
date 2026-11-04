-- Booking, Reservation & Lifecycle State Machine Schema (Partition 05-05, SRS §05.05, ADR-001, ADR-002, ADR-004, D-05-05-001, D-05-05-002)

-- 1. Canonical Bookings Entity
CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_reference VARCHAR(32) NOT NULL UNIQUE,
    renter_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    fleet_id UUID NOT NULL REFERENCES fleets(id) ON DELETE RESTRICT,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    hire_period DATERANGE NOT NULL,
    duration_days INT NOT NULL,
    daily_rate_applied_minor BIGINT NOT NULL,
    base_hire_fee_minor BIGINT NOT NULL,
    attachments_fee_minor BIGINT NOT NULL DEFAULT 0,
    gross_rental_amount_minor BIGINT NOT NULL,
    prepayment_amount_minor BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ETB',
    site_location_address TEXT NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING_PAYMENT',
    rejection_reason_code VARCHAR(64),
    rejection_notes TEXT,
    cancellation_reason TEXT,
    cancelled_by UUID REFERENCES users(id),
    penalty_amount_minor BIGINT NOT NULL DEFAULT 0,
    refund_amount_minor BIGINT NOT NULL DEFAULT 0,
    expires_at TIMESTAMPTZ,
    confirmed_at TIMESTAMPTZ,
    activated_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    handover_pin VARCHAR(6),
    check_in_meter_hours NUMERIC(10, 2),
    check_out_meter_hours NUMERIC(10, 2),
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_bookings_status CHECK (status IN (
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
        'RETURN_INITIATED',
        'PENDING_FINAL_PAYMENT',
        'COMPLETED',
        'RATED',
        'CANCELLED'
    )),
    CONSTRAINT chk_bookings_dates CHECK (end_date > start_date AND duration_days > 0),
    CONSTRAINT chk_bookings_rates CHECK (
        daily_rate_applied_minor > 0 AND 
        base_hire_fee_minor > 0 AND 
        gross_rental_amount_minor > 0 AND 
        prepayment_amount_minor > 0 AND
        penalty_amount_minor >= 0 AND
        refund_amount_minor >= 0
    )
);

CREATE INDEX IF NOT EXISTS idx_bookings_renter_id ON bookings(renter_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_bookings_owner_id ON bookings(owner_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_bookings_fleet_id ON bookings(fleet_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_bookings_reference ON bookings(booking_reference);
CREATE INDEX IF NOT EXISTS idx_bookings_expires_at ON bookings(expires_at) WHERE status = 'PENDING_PAYMENT';
CREATE INDEX IF NOT EXISTS idx_bookings_deleted_at ON bookings(deleted_at) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_bookings_updated_at
BEFORE UPDATE ON bookings
FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();

-- 2. Selected Bundled Attachments per Booking (SRS §3.2, BK-03)
CREATE TABLE IF NOT EXISTS booking_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    attachment_id UUID NOT NULL REFERENCES attachments(id) ON DELETE RESTRICT,
    bundled_daily_rate_minor BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_booking_attachment UNIQUE (booking_id, attachment_id),
    CONSTRAINT chk_booking_attachment_rate CHECK (bundled_daily_rate_minor >= 0)
);

CREATE INDEX IF NOT EXISTS idx_booking_attachments_booking_id ON booking_attachments(booking_id);
CREATE INDEX IF NOT EXISTS idx_booking_attachments_attachment_id ON booking_attachments(attachment_id);

-- 3. Non-Custodial Bank Wire Slips & Payment Evidence (SRS §4.1, BK-07, BK-08)
CREATE TABLE IF NOT EXISTS booking_bank_slips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    uploader_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    storage_uri TEXT NOT NULL,
    file_hash_sha256 VARCHAR(64) NOT NULL,
    transaction_reference VARCHAR(100),
    amount_minor BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ETB',
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING_AUDIT',
    rejection_reason_code VARCHAR(64),
    rejection_notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    audited_at TIMESTAMPTZ,
    CONSTRAINT chk_booking_bank_slips_status CHECK (status IN ('PENDING_AUDIT', 'VERIFIED', 'REJECTED')),
    CONSTRAINT chk_booking_bank_slips_amount CHECK (amount_minor > 0)
);

CREATE INDEX IF NOT EXISTS idx_booking_bank_slips_booking_id ON booking_bank_slips(booking_id);
CREATE INDEX IF NOT EXISTS idx_booking_bank_slips_file_hash ON booking_bank_slips(file_hash_sha256);
CREATE INDEX IF NOT EXISTS idx_booking_bank_slips_status ON booking_bank_slips(status);

-- 4. Immutable Booking State Machine Audit Logs (OWASP ASVS §4.0.3 Non-Repudiation)
CREATE TABLE IF NOT EXISTS booking_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    actor_id UUID REFERENCES users(id),
    action VARCHAR(64) NOT NULL,
    from_status VARCHAR(32) NOT NULL,
    to_status VARCHAR(32) NOT NULL,
    reason TEXT,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_booking_audit_logs_booking_id ON booking_audit_logs(booking_id, created_at DESC);

-- 5. Dynamic RBAC Permissions for Bookings (ADR-018)
INSERT INTO permissions (code, resource, action, description) VALUES
('bookings:create', 'bookings', 'create', 'Create rental requisition and hold lease'),
('bookings:read', 'bookings', 'read', 'View booking details and history'),
('bookings:update', 'bookings', 'update', 'Upload payment evidence, update site logistics or meter status'),
('bookings:confirm_payment', 'bookings', 'confirm_payment', 'Verify bank transfer slip and confirm booking'),
('bookings:decline_payment', 'bookings', 'decline_payment', 'Decline unverified bank transfer evidence'),
('bookings:dispute', 'bookings', 'dispute', 'Raise formal dispute regarding payment evidence decline'),
('bookings:cancel', 'bookings', 'cancel', 'Cancel pending or confirmed booking subject to contractual penalty'),
('bookings:arbitrate', 'bookings', 'arbitrate', 'Arbitrate disputed booking settlements')
ON CONFLICT (code) DO NOTHING;

-- Map permissions to roles
-- super_admin & admin: all booking permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.code IN ('super_admin', 'admin') 
  AND p.code IN (
      'bookings:create', 'bookings:read', 'bookings:update', 
      'bookings:confirm_payment', 'bookings:decline_payment', 
      'bookings:dispute', 'bookings:cancel', 'bookings:arbitrate'
  )
ON CONFLICT DO NOTHING;

-- owner: read, confirm_payment, decline_payment, cancel, update
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.code = 'owner' 
  AND p.code IN (
      'bookings:read', 'bookings:update', 
      'bookings:confirm_payment', 'bookings:decline_payment', 
      'bookings:cancel'
  )
ON CONFLICT DO NOTHING;

-- renter: create, read, update, dispute, cancel
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.code = 'renter' 
  AND p.code IN (
      'bookings:create', 'bookings:read', 'bookings:update', 
      'bookings:dispute', 'bookings:cancel'
  )
ON CONFLICT DO NOTHING;
