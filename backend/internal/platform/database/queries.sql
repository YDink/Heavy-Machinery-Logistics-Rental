-- name: GetCountryConfig :one
SELECT * FROM country_configs
WHERE country_code = $1 AND is_enabled = true;

-- name: ListEnabledCountryConfigs :many
SELECT * FROM country_configs
WHERE is_enabled = true
ORDER BY country_code ASC;
