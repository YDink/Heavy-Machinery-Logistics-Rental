-- Country Configurations
-- name: GetCountryConfig :one
SELECT * FROM country_configs
WHERE country_code = $1 AND is_enabled = true;

-- name: ListEnabledCountryConfigs :many
SELECT * FROM country_configs
WHERE is_enabled = true
ORDER BY country_code ASC;

-- Users & Authentication (REG-01-BE to REG-09-BE, AUT-01-BE to AUT-11-BE)
-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByPhoneNumber :one
SELECT * FROM users
WHERE phone_number = $1;

-- name: CreateUser :one
INSERT INTO users (
    phone_number,
    country_code,
    full_name,
    primary_role,
    account_status,
    kyc_status,
    password_hash
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: UpdateUserPhoneVerified :one
UPDATE users
SET account_status = 'ACTIVE',
    phone_verified_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserAccountStatus :one
UPDATE users
SET account_status = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserKYCStatus :one
UPDATE users
SET kyc_status = $2
WHERE id = $1
RETURNING *;

-- name: IncrementUserTokenVersion :one
UPDATE users
SET token_version = token_version + 1
WHERE id = $1
RETURNING token_version;

-- Legal Agreements
-- name: RecordUserLegalAgreement :one
INSERT INTO user_legal_agreements (
    user_id,
    terms_version,
    privacy_version,
    client_ip,
    client_user_agent
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: ListUserLegalAgreements :many
SELECT * FROM user_legal_agreements
WHERE user_id = $1
ORDER BY accepted_at DESC;


-- Dynamic RBAC Queries (ADR-018)
-- name: GetUserPermissions :many
SELECT DISTINCT p.code
FROM permissions p
JOIN role_permissions rp ON rp.permission_id = p.id
JOIN user_roles ur ON ur.role_id = rp.role_id
WHERE ur.user_id = $1
ORDER BY p.code ASC;

-- name: UserHasPermission :one
SELECT EXISTS (
    SELECT 1
    FROM permissions p
    JOIN role_permissions rp ON rp.permission_id = p.id
    JOIN user_roles ur ON ur.role_id = rp.role_id
    WHERE ur.user_id = $1 AND p.code = $2
);

-- name: GetUserRoles :many
SELECT r.*
FROM roles r
JOIN user_roles ur ON ur.role_id = r.id
WHERE ur.user_id = $1
ORDER BY r.code ASC;

-- name: AssignRoleToUser :exec
INSERT INTO user_roles (user_id, role_id, assigned_by)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, role_id) DO NOTHING;

-- name: AssignRoleByCodeToUser :exec
INSERT INTO user_roles (user_id, role_id, assigned_by)
SELECT $1, r.id, $3
FROM roles r
WHERE r.code = $2
ON CONFLICT (user_id, role_id) DO NOTHING;

-- name: RemoveRoleFromUser :exec
DELETE FROM user_roles
WHERE user_id = $1 AND role_id = $2;

-- name: ListAllRoles :many
SELECT * FROM roles
ORDER BY code ASC;

-- name: ListAllPermissions :many
SELECT * FROM permissions
ORDER BY resource ASC, action ASC;

-- Dynamic Localization Queries (ADR-019, SRS §05.16)
-- name: UpdateUserPreferredLanguage :one
UPDATE users
SET preferred_language = $2
WHERE id = $1
RETURNING *;

-- name: GetTranslationsByLanguage :many
SELECT 
    k.full_key,
    t.translated_text,
    t.language_code,
    n.namespace_key
FROM localization_translations t
JOIN localization_keys k ON k.id = t.key_id
JOIN localization_namespaces n ON n.id = k.namespace_id
WHERE t.language_code = $1
ORDER BY k.full_key ASC;

-- name: GetTranslationsBundleByNamespace :many
SELECT 
    k.full_key,
    t.translated_text,
    t.language_code
FROM localization_translations t
JOIN localization_keys k ON k.id = t.key_id
JOIN localization_namespaces n ON n.id = k.namespace_id
WHERE n.namespace_key = $1 AND t.language_code = $2
ORDER BY k.full_key ASC;

-- name: UpsertTranslation :one
INSERT INTO localization_translations (
    key_id,
    language_code,
    translated_text,
    is_verified
) VALUES (
    $1, $2, $3, $4
) ON CONFLICT (key_id, language_code) DO UPDATE SET
    translated_text = EXCLUDED.translated_text,
    is_verified = EXCLUDED.is_verified,
    updated_at = NOW()
RETURNING *;


-- Category Taxonomy Queries (Partition 05-02, FL-01)
-- name: ListAllCategories :many
SELECT * FROM categories
WHERE is_active = true
ORDER BY display_order ASC, name_en ASC;

-- name: ListRootCategories :many
SELECT * FROM categories
WHERE parent_id IS NULL AND is_active = true
ORDER BY display_order ASC, name_en ASC;

-- name: ListSubcategories :many
SELECT * FROM categories
WHERE parent_id = $1 AND is_active = true
ORDER BY display_order ASC, name_en ASC;

-- name: GetCategoryByID :one
SELECT * FROM categories
WHERE id = $1;

-- Spec Field Groups & Definitions (FL-07 to FL-15)
-- name: ListSpecGroupsByCategory :many
SELECT * FROM spec_field_groups
WHERE category_id = $1
ORDER BY display_order ASC;

-- name: ListSpecDefinitionsByGroup :many
SELECT * FROM spec_field_definitions
WHERE group_id = $1
ORDER BY display_order ASC;

-- name: GetCategorySpecDefinitions :many
SELECT 
    g.id AS group_id,
    g.section_name_en,
    g.section_name_am,
    g.display_order AS group_display_order,
    d.id AS definition_id,
    d.field_key,
    d.label_en,
    d.label_am,
    d.field_type,
    d.unit_symbol,
    d.is_required,
    d.is_logistics_critical,
    d.is_comparable,
    d.is_visible_to_renter,
    d.min_value,
    d.max_value,
    d.enum_values_json,
    d.display_order AS definition_display_order
FROM spec_field_groups g
JOIN spec_field_definitions d ON d.group_id = g.id
WHERE g.category_id = $1
ORDER BY g.display_order ASC, d.display_order ASC;

-- Fleet Core CRUD & Management (FL-01 to FL-06, FL-16 to FL-24)
-- name: CreateFleet :one
INSERT INTO fleets (
    owner_id,
    subcategory_id,
    make,
    model,
    manufacture_year,
    plate_number,
    chassis_vin,
    baseline_meter_hours,
    fuel_type,
    self_move_capable,
    daily_rate_minor,
    excess_hour_rate_minor,
    min_rental_days,
    payment_timeout_hours,
    operator_included,
    status,
    location,
    address_text
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
    ST_SetSRID(ST_MakePoint($17, $18), 4326)::geography, $19
) RETURNING 
    id, owner_id, subcategory_id, make, model, manufacture_year, plate_number, 
    chassis_vin, baseline_meter_hours, fuel_type, self_move_capable, 
    daily_rate_minor, excess_hour_rate_minor, min_rental_days, payment_timeout_hours, 
    operator_included, status, address_text, version, created_at, updated_at,
    ST_Y(location::geometry) AS latitude,
    ST_X(location::geometry) AS longitude;

-- name: GetFleetByID :one
SELECT 
    f.id, f.owner_id, f.subcategory_id, f.make, f.model, f.manufacture_year, f.plate_number, 
    f.chassis_vin, f.baseline_meter_hours, f.fuel_type, f.self_move_capable, 
    f.daily_rate_minor, f.excess_hour_rate_minor, f.min_rental_days, f.payment_timeout_hours, 
    f.operator_included, f.status, f.address_text, f.version, f.created_at, f.updated_at,
    ST_Y(f.location::geometry) AS latitude,
    ST_X(f.location::geometry) AS longitude,
    c.name_en AS subcategory_name_en,
    c.name_am AS subcategory_name_am
FROM fleets f
JOIN categories c ON c.id = f.subcategory_id
WHERE f.id = $1;

-- name: ListFleetsByOwner :many
SELECT 
    f.id, f.owner_id, f.subcategory_id, f.make, f.model, f.manufacture_year, f.plate_number, 
    f.chassis_vin, f.baseline_meter_hours, f.fuel_type, f.self_move_capable, 
    f.daily_rate_minor, f.excess_hour_rate_minor, f.min_rental_days, f.payment_timeout_hours, 
    f.operator_included, f.status, f.address_text, f.version, f.created_at, f.updated_at,
    ST_Y(f.location::geometry) AS latitude,
    ST_X(f.location::geometry) AS longitude
FROM fleets f
WHERE f.owner_id = $1
ORDER BY f.created_at DESC;

-- name: UpdateFleetStatus :one
UPDATE fleets
SET status = $2,
    version = version + 1
WHERE id = $1
RETURNING *;

-- Fleet Photos (SRS §3.3 FL-05)
-- name: CreateFleetPhoto :one
INSERT INTO fleet_photos (
    fleet_id,
    angle,
    storage_uri,
    file_hash,
    display_order
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: ListFleetPhotos :many
SELECT * FROM fleet_photos
WHERE fleet_id = $1
ORDER BY display_order ASC;

-- Fleet Spec Values & Indexing (SRS §4.3 FL-14)
-- name: UpsertFleetSpecValue :one
INSERT INTO fleet_spec_values (
    fleet_id,
    spec_definition_id,
    value_string,
    value_numeric,
    value_boolean
) VALUES (
    $1, $2, $3, $4, $5
) ON CONFLICT (fleet_id, spec_definition_id) DO UPDATE SET
    value_string = EXCLUDED.value_string,
    value_numeric = EXCLUDED.value_numeric,
    value_boolean = EXCLUDED.value_boolean
RETURNING *;

-- name: ListFleetSpecValues :many
SELECT 
    v.id, v.fleet_id, v.spec_definition_id, v.value_string, v.value_numeric, v.value_boolean,
    d.field_key, d.label_en, d.label_am, d.unit_symbol, d.field_type
FROM fleet_spec_values v
JOIN spec_field_definitions d ON d.id = v.spec_definition_id
WHERE v.fleet_id = $1
ORDER BY d.display_order ASC;

-- name: UpsertFleetSpecIndexed :one
INSERT INTO fleet_spec_indexed (
    fleet_id,
    operating_weight_kg,
    transport_length_mm,
    transport_width_mm,
    transport_height_mm,
    payload_capacity_kg,
    gross_engine_power_hp,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW()
) ON CONFLICT (fleet_id) DO UPDATE SET
    operating_weight_kg = EXCLUDED.operating_weight_kg,
    transport_length_mm = EXCLUDED.transport_length_mm,
    transport_width_mm = EXCLUDED.transport_width_mm,
    transport_height_mm = EXCLUDED.transport_height_mm,
    payload_capacity_kg = EXCLUDED.payload_capacity_kg,
    gross_engine_power_hp = EXCLUDED.gross_engine_power_hp,
    updated_at = NOW()
RETURNING *;

-- name: GetFleetSpecIndexed :one
SELECT * FROM fleet_spec_indexed
WHERE fleet_id = $1;


-- Attachments Inventory & Lifecycle (Partition 05-03, SRS §05.03, ATT-01 to ATT-14)
-- name: CreateAttachment :one
INSERT INTO attachments (
    owner_id,
    attachment_type,
    make,
    model,
    sub_type_size,
    serial_number,
    operating_weight_kg,
    min_carrier_weight_tonnes,
    max_carrier_weight_tonnes,
    required_flow_min_lmin,
    required_flow_max_lmin,
    required_pressure_bar,
    pin_diameter_mm,
    pin_spread_mm,
    standalone_daily_rate_minor,
    security_deposit_minor,
    min_rental_days,
    payment_evidence_window_hours,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
) RETURNING *;

-- name: GetAttachmentByID :one
SELECT * FROM attachments
WHERE id = $1;

-- name: ListAttachmentsByOwner :many
SELECT * FROM attachments
WHERE owner_id = $1 AND (status = $2 OR $2 = '')
ORDER BY created_at DESC;

-- name: UpdateAttachmentStatus :one
UPDATE attachments
SET status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CreateAttachmentCompatibleCategory :one
INSERT INTO attachment_compatible_categories (
    attachment_id,
    category_code
) VALUES (
    $1, $2
) ON CONFLICT (attachment_id, category_code) DO NOTHING
RETURNING *;

-- name: ListAttachmentCompatibleCategories :many
SELECT category_code FROM attachment_compatible_categories
WHERE attachment_id = $1;

-- name: CreateAttachmentPhoto :one
INSERT INTO attachment_photos (
    attachment_id,
    angle,
    storage_uri
) VALUES (
    $1, $2, $3
) ON CONFLICT (attachment_id, angle) DO UPDATE SET
    storage_uri = EXCLUDED.storage_uri
RETURNING *;

-- name: ListAttachmentPhotos :many
SELECT * FROM attachment_photos
WHERE attachment_id = $1;

-- name: UpsertFleetBundledAttachment :one
INSERT INTO fleet_bundled_attachments (
    fleet_id,
    attachment_id,
    bundled_daily_rate_minor,
    is_active
) VALUES (
    $1, $2, $3, $4
) ON CONFLICT (fleet_id, attachment_id) DO UPDATE SET
    bundled_daily_rate_minor = EXCLUDED.bundled_daily_rate_minor,
    is_active = EXCLUDED.is_active,
    updated_at = NOW()
RETURNING *;

-- name: ListBundledAttachmentsByFleet :many
SELECT 
    fba.id AS link_id,
    fba.fleet_id,
    fba.attachment_id,
    fba.bundled_daily_rate_minor,
    fba.is_active,
    a.attachment_type,
    a.make,
    a.model,
    a.sub_type_size,
    a.serial_number,
    a.standalone_daily_rate_minor,
    a.status
FROM fleet_bundled_attachments fba
JOIN attachments a ON a.id = fba.attachment_id
WHERE fba.fleet_id = $1 AND fba.is_active = true;

-- name: CheckAttachmentOverlap :one
SELECT EXISTS (
    SELECT 1 FROM attachment_reservations
    WHERE attachment_id = $1
      AND hire_period && $2::daterange
      AND status IN ('HOLD_PENDING', 'PAYMENT_EVIDENCE_UPLOADED', 'CONFIRMED', 'ACTIVE', 'RETURN_INITIATED')
);

-- name: CreateAttachmentReservation :one
INSERT INTO attachment_reservations (
    attachment_id,
    reference_id,
    hire_period,
    daily_rate_applied_minor,
    status
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;
