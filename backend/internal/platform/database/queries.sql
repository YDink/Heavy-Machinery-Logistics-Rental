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


-- Discovery, Spatial Search & Dynamic Comparison (Partition 05-04, SRS §05.04, SR-01 to SR-08, D-05-04-001, D-05-04-002)

-- name: SearchFleetsSpatial :many
SELECT 
    f.id,
    f.owner_id,
    f.subcategory_id,
    f.make,
    f.model,
    f.manufacture_year,
    f.plate_number,
    f.baseline_meter_hours,
    f.fuel_type,
    f.self_move_capable,
    f.daily_rate_minor,
    f.excess_hour_rate_minor,
    f.min_rental_days,
    f.operator_included,
    f.status,
    f.address_text,
    f.country_code,
    ST_Y(f.location::geometry) AS latitude,
    ST_X(f.location::geometry) AS longitude,
    ST_Distance(f.location, ST_SetSRID(ST_MakePoint(sqlc.arg(longitude)::float8, sqlc.arg(latitude)::float8), 4326)::geography) AS distance_meters,
    c.name_en AS subcategory_name_en,
    c.name_am AS subcategory_name_am,
    COALESCE(oa.average_rating, 5.00) AS owner_rating,
    COALESCE(oa.total_completed_rentals, 0) AS owner_completed_rentals,
    COALESCE(oa.is_verified, false) AS owner_verified,
    fsi.operating_weight_kg,
    fsi.payload_capacity_kg,
    fsi.gross_engine_power_hp,
    (SELECT storage_uri FROM fleet_photos fp WHERE fp.fleet_id = f.id ORDER BY fp.display_order ASC LIMIT 1) AS primary_photo_uri
FROM fleets f
JOIN categories c ON c.id = f.subcategory_id
LEFT JOIN owner_aggregates oa ON oa.owner_id = f.owner_id
LEFT JOIN fleet_spec_indexed fsi ON fsi.fleet_id = f.id
WHERE f.status = 'ACTIVE'
  AND f.country_code = sqlc.arg(country_code)
  AND ST_DWithin(f.location, ST_SetSRID(ST_MakePoint(sqlc.arg(longitude)::float8, sqlc.arg(latitude)::float8), 4326)::geography, sqlc.arg(radius_meters)::float8)
  AND (sqlc.narg(subcategory_id)::uuid IS NULL OR f.subcategory_id = sqlc.narg(subcategory_id))
  AND (sqlc.narg(min_rate_minor)::bigint IS NULL OR f.daily_rate_minor >= sqlc.narg(min_rate_minor))
  AND (sqlc.narg(max_rate_minor)::bigint IS NULL OR f.daily_rate_minor <= sqlc.narg(max_rate_minor))
  AND (sqlc.narg(fuel_type)::text IS NULL OR f.fuel_type = sqlc.narg(fuel_type))
  AND (sqlc.narg(operator_included)::boolean IS NULL OR f.operator_included = sqlc.narg(operator_included))
  AND (sqlc.narg(self_move_capable)::boolean IS NULL OR f.self_move_capable = sqlc.narg(self_move_capable))
  AND (sqlc.narg(query_text)::text IS NULL OR f.search_vector ILIKE '%' || sqlc.narg(query_text) || '%')
  AND (
      sqlc.narg(start_date)::date IS NULL OR sqlc.narg(end_date)::date IS NULL OR (
          NOT EXISTS (
              SELECT 1 FROM fleet_availability fa 
              WHERE fa.fleet_id = f.id 
                AND fa.hire_period && daterange(sqlc.narg(start_date)::date, sqlc.narg(end_date)::date, '[]')
                AND fa.state IN ('RESERVED', 'MAINTENANCE', 'RENTED')
          )
          AND NOT EXISTS (
              SELECT 1 FROM fleet_blackouts fb 
              WHERE fb.fleet_id = f.id 
                AND daterange(fb.start_date, fb.end_date, '[]') && daterange(sqlc.narg(start_date)::date, sqlc.narg(end_date)::date, '[]')
          )
      )
  )
ORDER BY distance_meters ASC
LIMIT sqlc.arg(query_limit)::int OFFSET sqlc.arg(query_offset)::int;

-- name: GetComparableSpecsForFleets :many
SELECT 
    d.field_key,
    d.label_en,
    d.label_am,
    d.unit_symbol,
    d.field_type,
    d.display_order,
    v.fleet_id,
    v.value_string,
    v.value_numeric,
    v.value_boolean
FROM spec_field_definitions d
JOIN fleet_spec_values v ON v.spec_definition_id = d.id
WHERE d.is_comparable = true
  AND v.fleet_id = ANY(sqlc.arg(fleet_ids)::uuid[])
ORDER BY d.display_order ASC;

-- name: GetFleetsForComparison :many
SELECT 
    f.id,
    f.subcategory_id,
    f.make,
    f.model,
    f.manufacture_year,
    f.daily_rate_minor,
    f.excess_hour_rate_minor,
    f.operator_included,
    f.fuel_type,
    f.self_move_capable,
    c.name_en AS subcategory_name_en,
    c.name_am AS subcategory_name_am,
    (SELECT storage_uri FROM fleet_photos fp WHERE fp.fleet_id = f.id ORDER BY fp.display_order ASC LIMIT 1) AS primary_photo_uri
FROM fleets f
JOIN categories c ON c.id = f.subcategory_id
WHERE f.id = ANY(sqlc.arg(fleet_ids)::uuid[]) AND f.status = 'ACTIVE';

-- name: ListFleetAvailabilityRanges :many
SELECT id, fleet_id, hire_period, state, created_at
FROM fleet_availability
WHERE fleet_id = $1
  AND hire_period && daterange($2::date, $3::date, '[]');

-- name: ListFleetBlackoutsInRange :many
SELECT id, fleet_id, start_date, end_date, reason
FROM fleet_blackouts
WHERE fleet_id = $1
  AND daterange(start_date, end_date, '[]') && daterange($2::date, $3::date, '[]');

-- name: CreateFleetAvailability :one
INSERT INTO fleet_availability (
    fleet_id,
    hire_period,
    state,
    reference_id
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: UpsertOwnerAggregate :one
INSERT INTO owner_aggregates (
    owner_id,
    average_rating,
    total_rentals_count,
    total_completed_rentals,
    completion_rate,
    is_verified,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, NOW()
) ON CONFLICT (owner_id) DO UPDATE SET
    average_rating = EXCLUDED.average_rating,
    total_rentals_count = EXCLUDED.total_rentals_count,
    total_completed_rentals = EXCLUDED.total_completed_rentals,
    completion_rate = EXCLUDED.completion_rate,
    is_verified = EXCLUDED.is_verified,
    updated_at = NOW()
RETURNING *;

-- name: GetOwnerAggregate :one
SELECT * FROM owner_aggregates
WHERE owner_id = $1;

-- name: UpsertFleetInsurance :one
INSERT INTO fleet_insurances (
    fleet_id,
    policy_number,
    provider_name,
    coverage_type,
    status,
    expires_at,
    certificate_uri,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW()
) RETURNING *;

-- name: GetFleetInsuranceByFleetID :one
SELECT * FROM fleet_insurances
WHERE fleet_id = $1 AND status = 'ACTIVE' AND expires_at >= CURRENT_DATE
ORDER BY expires_at DESC
LIMIT 1;


-- ==================== Booking State Machine (Partition 05-05) ====================

-- name: CreateBooking :one
INSERT INTO bookings (
    booking_reference,
    renter_id,
    owner_id,
    fleet_id,
    start_date,
    end_date,
    hire_period,
    duration_days,
    daily_rate_applied_minor,
    base_hire_fee_minor,
    attachments_fee_minor,
    gross_rental_amount_minor,
    prepayment_amount_minor,
    currency,
    site_location_address,
    status,
    expires_at,
    handover_pin,
    version
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, 1
) RETURNING *;

-- name: GetBookingByID :one
SELECT * FROM bookings
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetBookingByReference :one
SELECT * FROM bookings
WHERE booking_reference = $1 AND deleted_at IS NULL;

-- name: ListBookingsByRenter :many
SELECT * FROM bookings
WHERE renter_id = $1 
  AND ($2::text = '' OR status = $2)
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $3;

-- name: ListBookingsByOwner :many
SELECT * FROM bookings
WHERE owner_id = $1 
  AND ($2::text = '' OR status = $2)
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $3;

-- name: UpdateBookingStatusOCC :one
UPDATE bookings
SET status = $2,
    version = version + 1,
    rejection_reason_code = COALESCE($3, rejection_reason_code),
    rejection_notes = COALESCE($4, rejection_notes),
    cancellation_reason = COALESCE($5, cancellation_reason),
    cancelled_by = COALESCE($6, cancelled_by),
    penalty_amount_minor = COALESCE($7, penalty_amount_minor),
    refund_amount_minor = COALESCE($8, refund_amount_minor),
    confirmed_at = CASE WHEN $2 = 'CONFIRMED' THEN NOW() ELSE confirmed_at END,
    activated_at = CASE WHEN $2 = 'ACTIVE' THEN NOW() ELSE activated_at END,
    completed_at = CASE WHEN $2 = 'COMPLETED' THEN NOW() ELSE completed_at END,
    cancelled_at = CASE WHEN $2 = 'CANCELLED' THEN NOW() ELSE cancelled_at END,
    expires_at = COALESCE($9, expires_at),
    updated_at = NOW()
WHERE id = $1 AND version = $10 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateBookingHandover :one
UPDATE bookings
SET status = 'ACTIVE',
    check_in_meter_hours = $2,
    activated_at = NOW(),
    version = version + 1,
    updated_at = NOW()
WHERE id = $1 AND version = $3 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateBookingReturn :one
UPDATE bookings
SET status = $2,
    check_out_meter_hours = $3,
    completed_at = CASE WHEN $2 = 'COMPLETED' THEN NOW() ELSE completed_at END,
    version = version + 1,
    updated_at = NOW()
WHERE id = $1 AND version = $4 AND deleted_at IS NULL
RETURNING *;

-- name: AddBookingAttachment :one
INSERT INTO booking_attachments (
    booking_id,
    attachment_id,
    bundled_daily_rate_minor
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: ListBookingAttachments :many
SELECT ba.*, a.attachment_type, a.make, a.model, a.serial_number
FROM booking_attachments ba
JOIN attachments a ON a.id = ba.attachment_id
WHERE ba.booking_id = $1;

-- name: CreateBookingBankSlip :one
INSERT INTO booking_bank_slips (
    booking_id,
    uploader_id,
    storage_uri,
    file_hash_sha256,
    transaction_reference,
    amount_minor,
    currency,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetLatestBankSlipForBooking :one
SELECT * FROM booking_bank_slips
WHERE booking_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: UpdateBankSlipStatus :one
UPDATE booking_bank_slips
SET status = $2,
    rejection_reason_code = $3,
    rejection_notes = $4,
    audited_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CreateBookingAuditLog :one
INSERT INTO booking_audit_logs (
    booking_id,
    actor_id,
    action,
    from_status,
    to_status,
    reason,
    payload
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: ListBookingAuditLogs :many
SELECT * FROM booking_audit_logs
WHERE booking_id = $1
ORDER BY created_at ASC;


-- Active Rental Management, Shift Telemetry, Breakdown & Extension (Partition 05-06, SRS §05.06, RENT-01 to RENT-12, D-05-06-001, D-05-06-002)

-- name: CreateHandoverInspection :one
INSERT INTO handover_inspections (
    booking_id,
    inspector_id,
    check_in_meter_hours,
    fluid_levels_ok,
    safety_glass_intact,
    notes,
    photos_json
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetHandoverInspectionByBookingID :one
SELECT * FROM handover_inspections
WHERE booking_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: CreateHandoverChallenge :one
INSERT INTO handover_challenges (
    booking_id,
    challenge_pin,
    expires_at,
    attempt_count,
    max_attempts,
    is_verified
) VALUES (
    $1, $2, $3, 0, $4, false
) RETURNING *;

-- name: GetActiveHandoverChallenge :one
SELECT * FROM handover_challenges
WHERE booking_id = $1 AND is_verified = false AND expires_at > NOW()
ORDER BY created_at DESC
LIMIT 1;

-- name: IncrementHandoverChallengeAttempts :one
UPDATE handover_challenges
SET attempt_count = attempt_count + 1
WHERE id = $1
RETURNING *;

-- name: VerifyHandoverChallenge :one
UPDATE handover_challenges
SET is_verified = true,
    verified_at = NOW(),
    touch_signature_base64 = $2,
    contractor_lat = $3,
    contractor_lng = $4
WHERE id = $1
RETURNING *;

-- name: CreateRentalShift :one
INSERT INTO rental_shifts (
    booking_id,
    operator_id,
    shift_date,
    start_time,
    start_meter_hours,
    allocated_hours,
    accumulated_hours,
    excess_hours,
    is_paused
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetActiveRentalShift :one
SELECT * FROM rental_shifts
WHERE booking_id = $1 AND end_time IS NULL
ORDER BY created_at DESC
LIMIT 1;

-- name: UpdateRentalShiftProgress :one
UPDATE rental_shifts
SET accumulated_hours = $2,
    excess_hours = $3,
    end_meter_hours = $4,
    is_paused = $5,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: EndRentalShift :one
UPDATE rental_shifts
SET end_time = NOW(),
    end_meter_hours = $2,
    accumulated_hours = $3,
    excess_hours = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ListRentalShiftsByBookingID :many
SELECT * FROM rental_shifts
WHERE booking_id = $1
ORDER BY shift_date DESC, created_at DESC;

-- name: CreateRentalExtension :one
INSERT INTO rental_extensions (
    booking_id,
    requested_by_id,
    original_end_date,
    new_end_date,
    extension_days,
    incremental_gross_minor,
    incremental_prepayment_minor,
    currency,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetRentalExtensionByID :one
SELECT * FROM rental_extensions
WHERE id = $1;

-- name: ListRentalExtensionsByBookingID :many
SELECT * FROM rental_extensions
WHERE booking_id = $1
ORDER BY created_at DESC;

-- name: UpdateRentalExtensionStatus :one
UPDATE rental_extensions
SET status = $2,
    owner_decided_at = COALESCE($3, owner_decided_at),
    owner_decision_notes = COALESCE($4, owner_decision_notes),
    slip_storage_uri = COALESCE($5, slip_storage_uri),
    slip_file_hash = COALESCE($6, slip_file_hash),
    confirmed_at = CASE WHEN $2 = 'CONFIRMED' THEN NOW() ELSE confirmed_at END,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ExpandBookingHirePeriod :one
UPDATE bookings
SET end_date = $2,
    hire_period = daterange(start_date, $2, '[)'),
    duration_days = duration_days + $3,
    gross_rental_amount_minor = gross_rental_amount_minor + $4,
    prepayment_amount_minor = prepayment_amount_minor + $5,
    version = version + 1,
    updated_at = NOW()
WHERE id = $1 AND version = $6 AND deleted_at IS NULL
RETURNING *;

-- name: CreateMechanicalBreakdown :one
INSERT INTO mechanical_breakdowns (
    booking_id,
    reported_by_id,
    fault_category,
    fault_description,
    diagnostic_trouble_codes,
    current_hour_meter,
    incident_photos_json,
    status,
    sla_deadline
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetMechanicalBreakdownByID :one
SELECT * FROM mechanical_breakdowns
WHERE id = $1;

-- name: GetActiveMechanicalBreakdownByBookingID :one
SELECT * FROM mechanical_breakdowns
WHERE booking_id = $1 AND status IN ('ACTIVE_BREAKDOWN_PAUSED', 'MECHANIC_DISPATCHED')
ORDER BY created_at DESC
LIMIT 1;

-- name: ListMechanicalBreakdownsByBookingID :many
SELECT * FROM mechanical_breakdowns
WHERE booking_id = $1
ORDER BY created_at DESC;

-- name: ResolveMechanicalBreakdown :one
UPDATE mechanical_breakdowns
SET status = 'RESOLVED',
    mechanic_name = $2,
    mechanic_phone = $3,
    resolution_notes = $4,
    resolved_at = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetFleetBaselineMeter :one
SELECT id, baseline_meter_hours FROM fleets
WHERE id = $1;

-- name: CheckFleetAvailabilityConflict :one
SELECT EXISTS (
    SELECT 1 FROM bookings
    WHERE fleet_id = sqlc.arg(fleet_id)
      AND id != sqlc.arg(exclude_booking_id)
      AND hire_period && sqlc.arg(hire_period)::daterange
      AND status IN ('HOLD_PENDING', 'PAYMENT_EVIDENCE_UPLOADED', 'CONFIRMED', 'ACTIVE', 'RETURN_INITIATED', 'ACTIVE_BREAKDOWN_PAUSED', 'ACTIVE_EXTENSION_HOLD')
      AND deleted_at IS NULL
);


-- Rental Completion, Return Inspection & Settlement Reconciliation (Partition 05-07)
-- name: CreateReturnInspection :one
INSERT INTO return_inspections (
    booking_id,
    inspector_id,
    check_out_meter_hours,
    fluid_levels_ok,
    safety_glass_intact,
    accelerometer_steady,
    notes,
    photos_json
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetReturnInspectionByBookingID :one
SELECT * FROM return_inspections
WHERE booking_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: CreateDamageClaimDocket :one
INSERT INTO damage_claim_dockets (
    booking_id,
    claimant_id,
    component_category,
    incident_narrative,
    repair_quote_uri,
    claimed_amount_minor,
    accepted_amount_minor,
    currency,
    damage_photos_json,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetDamageClaimDocketByID :one
SELECT * FROM damage_claim_dockets
WHERE id = $1;

-- name: GetDamageClaimDocketByBookingID :one
SELECT * FROM damage_claim_dockets
WHERE booking_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: ListDamageClaimDocketsByBookingID :many
SELECT * FROM damage_claim_dockets
WHERE booking_id = $1
ORDER BY created_at DESC;

-- name: RespondToDamageClaim :one
UPDATE damage_claim_dockets
SET status = $2,
    accepted_amount_minor = $3,
    dispute_reason = $4,
    contractor_responded_at = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ArbitrateDamageClaim :one
UPDATE damage_claim_dockets
SET status = $2,
    accepted_amount_minor = $3,
    arbitration_notes = $4,
    arbitrated_by = $5,
    arbitrated_at = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: WaiveDamageClaim :one
UPDATE damage_claim_dockets
SET status = 'WAIVED',
    accepted_amount_minor = 0,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpsertRentalSettlement :one
INSERT INTO rental_settlements (
    booking_id,
    base_machine_fee_minor,
    bundled_attachments_fee_minor,
    contractual_allowed_hours,
    actual_operating_hours,
    credited_breakdown_hours,
    billable_excess_hours,
    excess_hourly_rate_minor,
    excess_hours_surcharge_minor,
    accepted_damage_minor,
    gross_liabilities_minor,
    prepayment_settled_minor,
    net_balance_minor,
    currency,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
ON CONFLICT (booking_id) DO UPDATE
SET base_machine_fee_minor = EXCLUDED.base_machine_fee_minor,
    bundled_attachments_fee_minor = EXCLUDED.bundled_attachments_fee_minor,
    contractual_allowed_hours = EXCLUDED.contractual_allowed_hours,
    actual_operating_hours = EXCLUDED.actual_operating_hours,
    credited_breakdown_hours = EXCLUDED.credited_breakdown_hours,
    billable_excess_hours = EXCLUDED.billable_excess_hours,
    excess_hourly_rate_minor = EXCLUDED.excess_hourly_rate_minor,
    excess_hours_surcharge_minor = EXCLUDED.excess_hours_surcharge_minor,
    accepted_damage_minor = EXCLUDED.accepted_damage_minor,
    gross_liabilities_minor = EXCLUDED.gross_liabilities_minor,
    prepayment_settled_minor = EXCLUDED.prepayment_settled_minor,
    net_balance_minor = EXCLUDED.net_balance_minor,
    status = EXCLUDED.status,
    updated_at = NOW()
RETURNING *;

-- name: GetRentalSettlementByBookingID :one
SELECT * FROM rental_settlements
WHERE booking_id = $1;

-- name: UpdateRentalSettlementFinalPaymentEvidence :one
UPDATE rental_settlements
SET final_slip_storage_uri = $2,
    final_slip_file_hash = $3,
    status = 'PAYMENT_EVIDENCE_UPLOADED',
    updated_at = NOW()
WHERE booking_id = $1
RETURNING *;

-- name: ConfirmRentalSettlement :one
UPDATE rental_settlements
SET status = 'SETTLED',
    settled_at = NOW(),
    updated_at = NOW()
WHERE booking_id = $1
RETURNING *;

-- name: CreateBookingReview :one
INSERT INTO booking_reviews (
    booking_id,
    author_id,
    target_id,
    role,
    rating,
    feedback,
    dimension_ratings_json,
    is_published,
    published_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetBookingReviewsByBookingID :many
SELECT * FROM booking_reviews
WHERE booking_id = $1
ORDER BY created_at ASC;

-- name: PublishBookingReviews :exec
UPDATE booking_reviews
SET is_published = true,
    published_at = NOW()
WHERE booking_id = $1;

-- name: ListPublishedReviewsForTarget :many
SELECT * FROM booking_reviews
WHERE target_id = $1 AND is_published = true
ORDER BY created_at DESC;

-- name: UpdateBookingCompletionSummary :one
UPDATE bookings
SET status = $2,
    check_out_meter_hours = $3,
    damage_claim_submitted = $4,
    damage_claim_waived = $5,
    accepted_damage_minor = $6,
    excess_hours_surcharge_minor = $7,
    final_net_balance_minor = $8,
    version = version + 1,
    updated_at = NOW()
WHERE id = $1 AND version = $9 AND deleted_at IS NULL
RETURNING *;
