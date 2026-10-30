DROP TRIGGER IF EXISTS trg_country_configs_updated_at ON country_configs;
DROP FUNCTION IF EXISTS set_updated_at_column();
DROP TABLE IF EXISTS country_configs;
DROP EXTENSION IF EXISTS "pgcrypto";
DROP EXTENSION IF EXISTS "postgis";
DROP EXTENSION IF EXISTS "uuid-ossp";
