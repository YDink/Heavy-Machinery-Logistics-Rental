-- Fleet Registration, Machinery Taxonomy & Logistics Specification Schema (Partition 05-02, SRS §05.02, ADR-003, ADR-005, ADR-018, ADR-019)

-- 1. Two-Level Category Taxonomy (SRS §3.1 FL-01, SRS Appendix C)
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL UNIQUE,
    name_en VARCHAR(100) NOT NULL,
    name_am VARCHAR(100) NOT NULL,
    icon_url TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_code ON categories(code);
CREATE TRIGGER trg_categories_updated_at
BEFORE UPDATE ON categories
FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();

-- 2. Dynamic Specification Field Groups (SRS §4.1)
CREATE TABLE spec_field_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    section_name_en VARCHAR(100) NOT NULL,
    section_name_am VARCHAR(100) NOT NULL,
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_spec_field_groups_category_id ON spec_field_groups(category_id);
CREATE TRIGGER trg_spec_field_groups_updated_at
BEFORE UPDATE ON spec_field_groups
FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();

-- 3. Dynamic Specification Field Definitions (SRS §4.1, FL-07 to FL-15)
CREATE TABLE spec_field_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL REFERENCES spec_field_groups(id) ON DELETE CASCADE,
    field_key VARCHAR(100) NOT NULL UNIQUE,
    label_en VARCHAR(100) NOT NULL,
    label_am VARCHAR(100) NOT NULL,
    field_type VARCHAR(20) NOT NULL,
    unit_symbol VARCHAR(20),
    is_required BOOLEAN NOT NULL DEFAULT false,
    is_logistics_critical BOOLEAN NOT NULL DEFAULT false,
    is_comparable BOOLEAN NOT NULL DEFAULT true,
    is_visible_to_renter BOOLEAN NOT NULL DEFAULT true,
    min_value NUMERIC(12, 2),
    max_value NUMERIC(12, 2),
    enum_values_json JSONB,
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_spec_field_type CHECK (field_type IN ('STRING', 'DECIMAL', 'INTEGER', 'ENUM', 'BOOLEAN'))
);

CREATE INDEX idx_spec_field_definitions_group_id ON spec_field_definitions(group_id);
CREATE INDEX idx_spec_field_definitions_field_key ON spec_field_definitions(field_key);
CREATE TRIGGER trg_spec_field_definitions_updated_at
BEFORE UPDATE ON spec_field_definitions
FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();

-- 4. Core Fleet Asset Entity (SRS §2.2, §3.1, §5.1, §7.1)
CREATE TABLE fleets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subcategory_id UUID NOT NULL REFERENCES categories(id),
    make VARCHAR(64) NOT NULL,
    model VARCHAR(64) NOT NULL,
    manufacture_year INT NOT NULL,
    plate_number VARCHAR(30) NOT NULL,
    chassis_vin VARCHAR(50) NOT NULL,
    baseline_meter_hours NUMERIC(10, 2) NOT NULL,
    fuel_type VARCHAR(20) NOT NULL,
    self_move_capable BOOLEAN NOT NULL DEFAULT false,
    daily_rate_minor BIGINT NOT NULL,
    excess_hour_rate_minor BIGINT NOT NULL,
    min_rental_days INT NOT NULL DEFAULT 1,
    payment_timeout_hours INT NOT NULL DEFAULT 24,
    operator_included BOOLEAN NOT NULL DEFAULT false,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING_APPROVAL',
    location geography(Point, 4326) NOT NULL,
    address_text TEXT NOT NULL,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_fleets_status CHECK (status IN (
        'DRAFT', 'PENDING_APPROVAL', 'ACTIVE', 'RENTED', 
        'INACTIVE', 'INSURANCE_HOLD', 'DELINQUENT_LOCKED', 'DECOMMISSIONED'
    )),
    CONSTRAINT chk_fleets_fuel_type CHECK (fuel_type IN ('DIESEL', 'ELECTRIC', 'PETROL', 'HYBRID', 'BIODIESEL')),
    CONSTRAINT chk_fleets_baseline_meter CHECK (baseline_meter_hours >= 0.0),
    CONSTRAINT chk_fleets_rates CHECK (daily_rate_minor > 0 AND excess_hour_rate_minor >= 0),
    CONSTRAINT chk_fleets_manufacture_year CHECK (manufacture_year >= 1980)
);

CREATE INDEX idx_fleets_owner_id ON fleets(owner_id);
CREATE INDEX idx_fleets_subcategory_id ON fleets(subcategory_id);
CREATE INDEX idx_fleets_status ON fleets(status);
CREATE INDEX idx_fleets_location ON fleets USING GIST(location);
CREATE TRIGGER trg_fleets_updated_at
BEFORE UPDATE ON fleets
FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();

-- 5. Mandatory 6-Angle Walkaround Photos (SRS §3.3 FL-05, D-05-02-001)
CREATE TABLE fleet_photos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fleet_id UUID NOT NULL REFERENCES fleets(id) ON DELETE CASCADE,
    angle VARCHAR(20) NOT NULL,
    storage_uri TEXT NOT NULL,
    file_hash VARCHAR(64) NOT NULL,
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_fleet_photos_angle CHECK (angle IN ('FRONT', 'REAR', 'LEFT', 'RIGHT', 'CABIN', 'METER'))
);

CREATE INDEX idx_fleet_photos_fleet_id ON fleet_photos(fleet_id);
CREATE UNIQUE INDEX uq_fleet_photos_fleet_angle ON fleet_photos(fleet_id, angle);

-- 6. Dynamic Specification Values (SRS §4.1)
CREATE TABLE fleet_spec_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fleet_id UUID NOT NULL REFERENCES fleets(id) ON DELETE CASCADE,
    spec_definition_id UUID NOT NULL REFERENCES spec_field_definitions(id) ON DELETE CASCADE,
    value_string TEXT,
    value_numeric NUMERIC(12, 2),
    value_boolean BOOLEAN,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_fleet_spec_values_fleet_def UNIQUE (fleet_id, spec_definition_id)
);

CREATE INDEX idx_fleet_spec_values_fleet_id ON fleet_spec_values(fleet_id);
CREATE INDEX idx_fleet_spec_values_def_id ON fleet_spec_values(spec_definition_id);

-- 7. Denormalized Logistics Performance Index (SRS §4.3 FL-14, D-05-02-002)
CREATE TABLE fleet_spec_indexed (
    fleet_id UUID PRIMARY KEY REFERENCES fleets(id) ON DELETE CASCADE,
    operating_weight_kg NUMERIC(10, 2),
    transport_length_mm NUMERIC(10, 2),
    transport_width_mm NUMERIC(10, 2),
    transport_height_mm NUMERIC(10, 2),
    payload_capacity_kg NUMERIC(10, 2),
    gross_engine_power_hp NUMERIC(10, 2),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fleet_spec_indexed_weight ON fleet_spec_indexed(operating_weight_kg);

-- 8. Operator Roster Linking & Fleet Crew Assignments (SRS §6.1, §6.2)
CREATE TABLE owner_operator_memberships (
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    operator_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (owner_id, operator_id),
    CONSTRAINT chk_memberships_status CHECK (status IN ('PENDING_INVITE', 'ACTIVE', 'REVOKED'))
);

CREATE TABLE fleet_operator_assignments (
    fleet_id UUID NOT NULL REFERENCES fleets(id) ON DELETE CASCADE,
    operator_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    shift_type VARCHAR(20) NOT NULL DEFAULT 'DAY_SHIFT',
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (fleet_id, operator_id),
    CONSTRAINT chk_assignments_shift CHECK (shift_type IN ('DAY_SHIFT', 'NIGHT_SHIFT', 'FULL_TIME'))
);

-- 9. Availability Blackout Calendar (SRS §8.1 FL-35)
CREATE TABLE fleet_blackouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fleet_id UUID NOT NULL REFERENCES fleets(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reason TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_fleet_blackouts_dates CHECK (end_date >= start_date)
);

CREATE INDEX idx_fleet_blackouts_fleet_id ON fleet_blackouts(fleet_id);

-- 10. Canonical Taxonomy Seed (SRS Appendix C)
-- Root Categories
INSERT INTO categories (id, parent_id, code, name_en, name_am, icon_url, display_order) VALUES
('019a8f00-0001-7000-8000-000000000001', NULL, 'EXCAVATION', 'Excavation', 'የቁፋሮ ማሽነሪዎች', 'https://cdn.dinkmachinery.com/icons/excavation.svg', 1),
('019a8f00-0001-7000-8000-000000000002', NULL, 'EARTHMOVING', 'Earthmoving', 'የመሬት ስራ ማሽነሪዎች', 'https://cdn.dinkmachinery.com/icons/earthmoving.svg', 2),
('019a8f00-0001-7000-8000-000000000003', NULL, 'LIFTING', 'Lifting & Cranes', 'የማንሻ ማሽነሪዎችና ክሬኖች', 'https://cdn.dinkmachinery.com/icons/lifting.svg', 3),
('019a8f00-0001-7000-8000-000000000004', NULL, 'HAULING', 'Hauling & Tipper', 'የጭነትና የቲፐር ተሽከርካሪዎች', 'https://cdn.dinkmachinery.com/icons/hauling.svg', 4),
('019a8f00-0001-7000-8000-000000000005', NULL, 'LOGISTICS', 'Heavy Logistics & Lowboy', 'የከባድ ማሽነሪ ማጓጓዣ (ሎውቦይ)', 'https://cdn.dinkmachinery.com/icons/logistics.svg', 5);

-- Subcategories for Excavation
INSERT INTO categories (id, parent_id, code, name_en, name_am, icon_url, display_order) VALUES
('019a8f00-0002-7000-8000-000000000001', '019a8f00-0001-7000-8000-000000000001', 'EXCAVATOR_MINI', 'Mini Excavator (1-5 tons)', 'አነስተኛ ኤክስካቫተር (1-5 ቶን)', 'https://cdn.dinkmachinery.com/icons/sub/mini_excavator.svg', 1),
('019a8f00-0002-7000-8000-000000000002', '019a8f00-0001-7000-8000-000000000001', 'EXCAVATOR_MIDI', 'Midi Excavator (5-10 tons)', 'መካከለኛ ኤክስካቫተር (5-10 ቶን)', 'https://cdn.dinkmachinery.com/icons/sub/midi_excavator.svg', 2),
('019a8f00-0002-7000-8000-000000000003', '019a8f00-0001-7000-8000-000000000001', 'EXCAVATOR_LARGE', 'Large Crawler Excavator (10-50 tons)', 'ትልቅ ክራውለር ኤክስካቫተር (10-50 ቶን)', 'https://cdn.dinkmachinery.com/icons/sub/crawler_excavator.svg', 3);

-- Subcategories for Earthmoving
INSERT INTO categories (id, parent_id, code, name_en, name_am, icon_url, display_order) VALUES
('019a8f00-0002-7000-8000-000000000004', '019a8f00-0001-7000-8000-000000000002', 'BULLDOZER', 'Bulldozer', 'ቡልዶዘር', 'https://cdn.dinkmachinery.com/icons/sub/bulldozer.svg', 1),
('019a8f00-0002-7000-8000-000000000005', '019a8f00-0001-7000-8000-000000000002', 'MOTOR_GRADER', 'Motor Grader', 'ሞተር ግሬደር', 'https://cdn.dinkmachinery.com/icons/sub/grader.svg', 2),
('019a8f00-0002-7000-8000-000000000006', '019a8f00-0001-7000-8000-000000000002', 'WHEEL_LOADER', 'Wheel Loader', 'ዊል ሎደር', 'https://cdn.dinkmachinery.com/icons/sub/wheel_loader.svg', 3);

-- Subcategories for Logistics
INSERT INTO categories (id, parent_id, code, name_en, name_am, icon_url, display_order) VALUES
('019a8f00-0002-7000-8000-000000000007', '019a8f00-0001-7000-8000-000000000005', 'LOWBOY_LIGHT', 'Lowboy Trailer (Light, <10 ton)', 'ቀላል ሎውቦይ ተሳቢ (<10 ቶን)', 'https://cdn.dinkmachinery.com/icons/sub/lowboy_light.svg', 1),
('019a8f00-0002-7000-8000-000000000008', '019a8f00-0001-7000-8000-000000000005', 'LOWBOY_MEDIUM', 'Lowboy Trailer (Medium, 10-30 ton)', 'መካከለኛ ሎውቦይ ተሳቢ (10-30 ቶን)', 'https://cdn.dinkmachinery.com/icons/sub/lowboy_medium.svg', 2),
('019a8f00-0002-7000-8000-000000000009', '019a8f00-0001-7000-8000-000000000005', 'LOWBOY_HEAVY', 'Lowboy Trailer (Heavy, 30+ ton)', 'ከባድ ሎውቦይ ተሳቢ (30+ ቶን)', 'https://cdn.dinkmachinery.com/icons/sub/lowboy_heavy.svg', 3);

-- 11. Seed Spec Field Groups & Definitions for Excavators (SRS §4.4 Example A)
INSERT INTO spec_field_groups (id, category_id, section_name_en, section_name_am, display_order) VALUES
('019a8f00-0003-7000-8000-000000000001', '019a8f00-0002-7000-8000-000000000001', 'Weights & Grounding', 'የክብደትና የመሬት ጫና ዝርዝሮች', 1),
('019a8f00-0003-7000-8000-000000000002', '019a8f00-0002-7000-8000-000000000001', 'Engine & Power', 'የሞተርና የጉልበት ዝርዝሮች', 2),
('019a8f00-0003-7000-8000-000000000003', '019a8f00-0002-7000-8000-000000000001', 'Physical Dimensions', 'የማሽኑ አካላዊ ስፋትና ልኬት', 3);

INSERT INTO spec_field_definitions (
    group_id, field_key, label_en, label_am, field_type, unit_symbol, 
    is_required, is_logistics_critical, is_comparable, min_value, max_value, display_order
) VALUES
('019a8f00-0003-7000-8000-000000000001', 'operating_weight_kg', 'Operating Weight', 'የስራ ክብደት', 'DECIMAL', 'kg', true, true, true, 500, 150000, 1),
('019a8f00-0003-7000-8000-000000000002', 'gross_power_hp', 'Gross Engine Power', 'ጠቅላላ የሞተር ጉልበት', 'DECIMAL', 'hp', true, false, true, 10, 2000, 1),
('019a8f00-0003-7000-8000-000000000003', 'transport_length_mm', 'Overall Transport Length', 'የማጓጓዣ ርዝመት', 'INTEGER', 'mm', true, true, true, 1000, 25000, 1),
('019a8f00-0003-7000-8000-000000000003', 'transport_width_mm', 'Overall Transport Width', 'የማጓጓዣ ወርድ', 'INTEGER', 'mm', true, true, true, 500, 6000, 2),
('019a8f00-0003-7000-8000-000000000003', 'transport_height_mm', 'Overall Transport Height', 'የማጓጓዣ ቁመት', 'INTEGER', 'mm', true, true, true, 500, 6000, 3);
