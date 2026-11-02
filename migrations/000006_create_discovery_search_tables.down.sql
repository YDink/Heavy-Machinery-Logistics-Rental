-- Rollback Fleet Discovery, Spatial Search & Dynamic Comparison Schema (Partition 05-04)

DELETE FROM role_permissions WHERE permission_id IN (SELECT id FROM permissions WHERE code IN ('fleets:search', 'fleets:compare'));
DELETE FROM permissions WHERE code IN ('fleets:search', 'fleets:compare');

DROP TABLE IF EXISTS fleet_availability CASCADE;
DROP TABLE IF EXISTS fleet_insurances CASCADE;
DROP TABLE IF EXISTS owner_aggregates CASCADE;

DROP INDEX IF EXISTS idx_fleet_spec_indexed_specs_json;
ALTER TABLE fleet_spec_indexed DROP COLUMN IF EXISTS specs_json;

DROP INDEX IF EXISTS idx_fleets_search_vector_trgm;
DROP INDEX IF EXISTS idx_fleets_country_code;
ALTER TABLE fleets DROP COLUMN IF EXISTS search_vector;
ALTER TABLE fleets DROP COLUMN IF EXISTS country_code;
