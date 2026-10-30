-- Identity & Authentication Tables (REG-01-BE to REG-09-BE, AUT-01-BE to AUT-11-BE)

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number VARCHAR(20) NOT NULL UNIQUE,
    country_code VARCHAR(2) NOT NULL REFERENCES country_configs(country_code),
    full_name VARCHAR(150) NOT NULL,
    primary_role VARCHAR(30) NOT NULL DEFAULT 'RENTER',
    account_status VARCHAR(30) NOT NULL DEFAULT 'UNVERIFIED',
    kyc_status VARCHAR(30) NOT NULL DEFAULT 'UNVERIFIED',
    password_hash VARCHAR(255),
    token_version INT NOT NULL DEFAULT 1,
    phone_verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_users_primary_role CHECK (primary_role IN ('RENTER', 'OWNER', 'OPERATOR', 'LOGISTICS', 'ADMIN')),
    CONSTRAINT chk_users_account_status CHECK (account_status IN ('UNVERIFIED', 'ACTIVE', 'SUSPENDED_SECURITY_FREEZE', 'DELINQUENT_LOCKED', 'TERMINATED')),
    CONSTRAINT chk_users_kyc_status CHECK (kyc_status IN ('UNVERIFIED', 'PENDING', 'VERIFIED', 'REJECTED'))
);

CREATE INDEX idx_users_phone_number ON users(phone_number);
CREATE INDEX idx_users_country_code ON users(country_code);
CREATE INDEX idx_users_account_status ON users(account_status);
CREATE INDEX idx_users_kyc_status ON users(kyc_status);

CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();

-- Statutory Terms & Privacy Acceptance Ledger (REG-08, REG-09)
CREATE TABLE user_legal_agreements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    terms_version VARCHAR(20) NOT NULL,
    privacy_version VARCHAR(20) NOT NULL,
    accepted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    client_ip VARCHAR(45) NOT NULL,
    client_user_agent TEXT NOT NULL
);

CREATE INDEX idx_user_legal_agreements_user_id ON user_legal_agreements(user_id);
