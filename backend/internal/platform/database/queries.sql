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
