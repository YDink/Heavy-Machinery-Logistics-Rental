-- Fleet Discovery, Spatial Search & Dynamic Comparison Schema (Partition 05-04, SRS §05.04, ADR-004, ADR-005, D-05-04-001, D-05-04-002)

-- 1. Enable pg_trgm for Full-Text Trigram Fuzzy Matching (SRS §3.2, SR-02)
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- 2. Add Jurisdiction & Generated Trigram Search Vector to fleets
ALTER TABLE fleets ADD COLUMN IF NOT EXISTS country_code VARCHAR(2) NOT NULL DEFAULT 'ET';
ALTER TABLE fleets ADD COLUMN IF NOT EXISTS search_vector TEXT GENERATED ALWAYS AS (make || ' ' || model || ' ' || address_text) STORED;

CREATE INDEX IF NOT EXISTS idx_fleets_country_code ON fleets(country_code);
CREATE INDEX IF NOT EXISTS idx_fleets_search_vector_trgm ON fleets USING gin (search_vector gin_trgm_ops);

-- 3. Dynamic Spec Faceting on fleet_spec_indexed via Pre-Aggregated JSONB (D-05-04-002, SR-03)
ALTER TABLE fleet_spec_indexed ADD COLUMN IF NOT EXISTS specs_json JSONB NOT NULL DEFAULT '{}'::jsonb;
CREATE INDEX IF NOT EXISTS idx_fleet_spec_indexed_specs_json ON fleet_spec_indexed USING gin (specs_json jsonb_path_ops);

-- 4. Owner Trust Ledger & Performance Ratings (SRS §3.4, §6.3, SR-04, SR-06)
CREATE TABLE IF NOT EXISTS owner_aggregates (
    owner_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    average_rating NUMERIC(3, 2) NOT NULL DEFAULT 5.00,
    total_rentals_count INT NOT NULL DEFAULT 0,
    total_completed_rentals INT NOT NULL DEFAULT 0,
    completion_rate NUMERIC(5, 4) NOT NULL DEFAULT 1.0000,
    is_verified BOOLEAN NOT NULL DEFAULT false,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 5. Comprehensive Fleet Insurance Status (SRS §6.2, §6.3, SR-06)
CREATE TABLE IF NOT EXISTS fleet_insurances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fleet_id UUID NOT NULL REFERENCES fleets(id) ON DELETE CASCADE,
    policy_number VARCHAR(100) NOT NULL,
    provider_name VARCHAR(100) NOT NULL,
    coverage_type VARCHAR(50) NOT NULL DEFAULT 'COMPREHENSIVE_PLANT_ALL_RISKS',
    status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',
    expires_at DATE NOT NULL,
    certificate_uri TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_fleet_insurances_status CHECK (status IN ('ACTIVE', 'EXPIRED', 'PENDING_AUDIT'))
);

CREATE INDEX IF NOT EXISTS idx_fleet_insurances_fleet_id ON fleet_insurances(fleet_id);

-- 6. Fleet Availability Temporal Range Exclusion (SRS §5.1, SR-07, BK-01)
CREATE TABLE IF NOT EXISTS fleet_availability (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fleet_id UUID NOT NULL REFERENCES fleets(id) ON DELETE CASCADE,
    hire_period DATERANGE NOT NULL,
    state VARCHAR(32) NOT NULL DEFAULT 'RESERVED',
    reference_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_fleet_avail_state CHECK (state IN ('AVAILABLE', 'RESERVED', 'MAINTENANCE', 'RENTED')),
    CONSTRAINT no_overlapping_fleet_availability EXCLUDE USING gist (
        fleet_id WITH =,
        hire_period WITH &&
    ) WHERE (state IN ('RESERVED', 'MAINTENANCE', 'RENTED'))
);

CREATE INDEX IF NOT EXISTS idx_fleet_avail_gist ON fleet_availability USING gist (fleet_id, hire_period);
CREATE INDEX IF NOT EXISTS idx_fleet_avail_state ON fleet_availability(state);

-- 7. Dynamic RBAC Permissions for Discovery & Comparison
INSERT INTO permissions (code, resource, action, description) VALUES
('fleets:search', 'fleets', 'search', 'Execute spatial radius and multi-variable fleet discovery queries'),
('fleets:compare', 'fleets', 'compare', 'Execute dynamic side-by-side EAV spec comparison')
ON CONFLICT (code) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.code IN ('super_admin', 'admin', 'owner', 'renter') AND p.code IN ('fleets:search', 'fleets:compare')
ON CONFLICT DO NOTHING;
