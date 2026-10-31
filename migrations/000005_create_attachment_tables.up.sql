-- Attachment Registration, Compatibility & Inventory Architecture Schema (Partition 05-03, SRS §05.03, ADR-003, ADR-005, ADR-018, ADR-019)

CREATE EXTENSION IF NOT EXISTS "btree_gist";

-- 1. Attachments Core Inventory Entity (SRS §2.2, §3.1, ATT-01 to ATT-04, ATT-09 to ATT-11)
CREATE TABLE attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    attachment_type VARCHAR(50) NOT NULL,
    make VARCHAR(64) NOT NULL,
    model VARCHAR(64) NOT NULL,
    sub_type_size VARCHAR(100) NOT NULL,
    serial_number VARCHAR(100) NOT NULL,
    operating_weight_kg NUMERIC(10, 2) NOT NULL,
    min_carrier_weight_tonnes NUMERIC(6, 2) NOT NULL,
    max_carrier_weight_tonnes NUMERIC(6, 2) NOT NULL,
    required_flow_min_lmin NUMERIC(8, 2) NOT NULL,
    required_flow_max_lmin NUMERIC(8, 2) NOT NULL,
    required_pressure_bar NUMERIC(8, 2) NOT NULL,
    pin_diameter_mm NUMERIC(8, 2) NOT NULL,
    pin_spread_mm NUMERIC(8, 2) NOT NULL,
    standalone_daily_rate_minor BIGINT NOT NULL,
    security_deposit_minor BIGINT NOT NULL DEFAULT 0,
    min_rental_days INT NOT NULL DEFAULT 1,
    payment_evidence_window_hours INT NOT NULL DEFAULT 24,
    status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_attachment_owner_serial UNIQUE (owner_id, serial_number),
    CONSTRAINT chk_attachment_status CHECK (status IN (
        'DRAFT', 'ACTIVE', 'RENTED', 'MAINTENANCE_HOLD', 'INACTIVE', 'DISPUTED_DAMAGE', 'DECOMMISSIONED'
    )),
    CONSTRAINT chk_carrier_weights CHECK (min_carrier_weight_tonnes <= max_carrier_weight_tonnes),
    CONSTRAINT chk_flow_bounds CHECK (required_flow_min_lmin <= required_flow_max_lmin),
    CONSTRAINT chk_standalone_daily_rate CHECK (standalone_daily_rate_minor >= 0)
);

CREATE INDEX idx_attachments_owner_id ON attachments(owner_id);
CREATE INDEX idx_attachments_status ON attachments(status);
CREATE INDEX idx_attachments_type ON attachments(attachment_type);

CREATE TRIGGER trg_attachments_updated_at
BEFORE UPDATE ON attachments
FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();

-- 2. Compatible Host Machinery Subcategories (SRS §3.2, ATT-05)
CREATE TABLE attachment_compatible_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attachment_id UUID NOT NULL REFERENCES attachments(id) ON DELETE CASCADE,
    category_code VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_attachment_category UNIQUE (attachment_id, category_code)
);

CREATE INDEX idx_attachment_compat_attachment_id ON attachment_compatible_categories(attachment_id);
CREATE INDEX idx_attachment_compat_cat_code ON attachment_compatible_categories(category_code);

-- 3. Attachment Conditioning Photos (SRS §3.3, ATT-07, ATT-08)
CREATE TABLE attachment_photos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attachment_id UUID NOT NULL REFERENCES attachments(id) ON DELETE CASCADE,
    angle VARCHAR(32) NOT NULL,
    storage_uri TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_attachment_photo_angle UNIQUE (attachment_id, angle),
    CONSTRAINT chk_attachment_photo_angle CHECK (angle IN ('OVERALL', 'BRACKET', 'PINS', 'TOOL_TIP'))
);

CREATE INDEX idx_attachment_photos_attachment_id ON attachment_photos(attachment_id);

-- 4. Bundled Attachment Overrides per Fleet Asset (SRS §4.2, FL-32, FL-33)
CREATE TABLE fleet_bundled_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fleet_id UUID NOT NULL REFERENCES fleets(id) ON DELETE CASCADE,
    attachment_id UUID NOT NULL REFERENCES attachments(id) ON DELETE RESTRICT,
    bundled_daily_rate_minor BIGINT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_fleet_attachment_pair UNIQUE (fleet_id, attachment_id),
    CONSTRAINT chk_bundled_daily_rate CHECK (bundled_daily_rate_minor >= 0)
);

CREATE INDEX idx_fleet_bundled_fleet_id ON fleet_bundled_attachments(fleet_id);
CREATE INDEX idx_fleet_bundled_attachment_id ON fleet_bundled_attachments(attachment_id);

CREATE TRIGGER trg_fleet_bundled_attachments_updated_at
BEFORE UPDATE ON fleet_bundled_attachments
FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();

-- 5. Attachment Temporal Range Exclusion & Reservations (SRS §5.2, FL-34, D-05-03-001)
CREATE TABLE attachment_reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attachment_id UUID NOT NULL REFERENCES attachments(id) ON DELETE CASCADE,
    reference_id UUID,
    hire_period DATERANGE NOT NULL,
    daily_rate_applied_minor BIGINT NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'HOLD_PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT no_overlapping_attachment_reservations EXCLUDE USING gist (
        attachment_id WITH =,
        hire_period WITH &&
    ) WHERE (status IN ('HOLD_PENDING', 'PAYMENT_EVIDENCE_UPLOADED', 'CONFIRMED', 'ACTIVE', 'RETURN_INITIATED'))
);

CREATE INDEX idx_attachment_reservations_gist ON attachment_reservations USING gist (attachment_id, hire_period);
CREATE INDEX idx_attachment_reservations_status ON attachment_reservations(status);

-- 6. Dynamic RBAC Seeds for Attachments (ADR-018)
INSERT INTO permissions (code, module, action, description) VALUES
('attachments:create', 'attachments', 'create', 'Register new standalone attachments'),
('attachments:read', 'attachments', 'read', 'List and view attachment details and inventory'),
('attachments:update', 'attachments', 'update', 'Update attachment technical details, pricing and status'),
('attachments:delete', 'attachments', 'delete', 'Decommission attachment inventory'),
('attachments:link', 'attachments', 'link', 'Link attachments to host machinery fleets with bundled rates')
ON CONFLICT (code) DO NOTHING;

-- Map permissions to roles
-- super_admin: All permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.code = 'super_admin' AND p.code LIKE 'attachments:%'
ON CONFLICT DO NOTHING;

-- admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.code = 'admin' AND p.code LIKE 'attachments:%'
ON CONFLICT DO NOTHING;

-- owner
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.code = 'owner' AND p.code IN (
    'attachments:create', 'attachments:read', 'attachments:update', 'attachments:delete', 'attachments:link'
)
ON CONFLICT DO NOTHING;

-- renter
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.code = 'renter' AND p.code IN ('attachments:read')
ON CONFLICT DO NOTHING;
