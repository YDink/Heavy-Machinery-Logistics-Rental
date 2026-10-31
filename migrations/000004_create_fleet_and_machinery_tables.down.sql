-- Rollback Fleet and Machinery Tables

DROP TABLE IF EXISTS fleet_blackouts CASCADE;
DROP TABLE IF EXISTS fleet_operator_assignments CASCADE;
DROP TABLE IF EXISTS owner_operator_memberships CASCADE;
DROP TABLE IF EXISTS fleet_spec_indexed CASCADE;
DROP TABLE IF EXISTS fleet_spec_values CASCADE;
DROP TABLE IF EXISTS fleet_photos CASCADE;

DROP TRIGGER IF EXISTS trg_fleets_updated_at ON fleets;
DROP TABLE IF EXISTS fleets CASCADE;

DROP TRIGGER IF EXISTS trg_spec_field_definitions_updated_at ON spec_field_definitions;
DROP TABLE IF EXISTS spec_field_definitions CASCADE;

DROP TRIGGER IF EXISTS trg_spec_field_groups_updated_at ON spec_field_groups;
DROP TABLE IF EXISTS spec_field_groups CASCADE;

DROP TRIGGER IF EXISTS trg_categories_updated_at ON categories;
DROP TABLE IF EXISTS categories CASCADE;
