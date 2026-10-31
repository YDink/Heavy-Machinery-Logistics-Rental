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
