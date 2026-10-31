-- Rollback Dynamic RBAC and Dynamic Localization Schemas

DROP TRIGGER IF EXISTS trg_localization_translations_updated_at ON localization_translations;
DROP TABLE IF EXISTS localization_translations CASCADE;
DROP TABLE IF EXISTS localization_keys CASCADE;
DROP TRIGGER IF EXISTS trg_localization_namespaces_updated_at ON localization_namespaces;
DROP TABLE IF EXISTS localization_namespaces CASCADE;

DROP TABLE IF EXISTS user_roles CASCADE;
DROP TABLE IF EXISTS role_permissions CASCADE;
DROP TABLE IF EXISTS permissions CASCADE;
DROP TRIGGER IF EXISTS trg_roles_updated_at ON roles;
DROP TABLE IF EXISTS roles CASCADE;

ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_users_preferred_language;
ALTER TABLE users DROP COLUMN IF EXISTS preferred_language;
