-- Dynamic RBAC & Full-Stack Vernacular Localization Schema (ADR-018, ADR-019, SRS §4 Appendix B, SRS §05.16)

-- 1. User Preferred Language (ADR-019, LOC-01-BE)
ALTER TABLE users ADD COLUMN preferred_language VARCHAR(5) NOT NULL DEFAULT 'en';
ALTER TABLE users ADD CONSTRAINT chk_users_preferred_language CHECK (preferred_language IN ('en', 'am', 'om', 'sw'));

-- 2. Dynamic RBAC Relational Schema (ADR-018)
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_system BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_roles_updated_at
BEFORE UPDATE ON roles
FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();

CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100) NOT NULL UNIQUE,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    assigned_by UUID REFERENCES users(id) ON DELETE SET NULL,
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);

-- 3. Dynamic Localization Schema (ADR-019, SRS §05.16 §2.2)
CREATE TABLE localization_namespaces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    namespace_key VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_localization_namespaces_updated_at
BEFORE UPDATE ON localization_namespaces
FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();

CREATE TABLE localization_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    namespace_id UUID NOT NULL REFERENCES localization_namespaces(id) ON DELETE CASCADE,
    full_key VARCHAR(150) NOT NULL UNIQUE,
    description TEXT,
    parameter_schema_json JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_localization_keys_namespace_id ON localization_keys(namespace_id);
CREATE INDEX idx_localization_keys_full_key ON localization_keys(full_key);

CREATE TABLE localization_translations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key_id UUID NOT NULL REFERENCES localization_keys(id) ON DELETE CASCADE,
    language_code VARCHAR(5) NOT NULL,
    translated_text TEXT NOT NULL,
    audio_storage_uri TEXT,
    is_verified BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_localization_translations_key_lang UNIQUE (key_id, language_code),
    CONSTRAINT chk_localization_translations_lang CHECK (language_code IN ('en', 'am', 'om', 'sw'))
);

CREATE INDEX idx_localization_translations_lang ON localization_translations(language_code);
CREATE TRIGGER trg_localization_translations_updated_at
BEFORE UPDATE ON localization_translations
FOR EACH ROW EXECUTE FUNCTION set_updated_at_column();

-- 4. Canonical Roles Seed (SRS §4.1 & Appendix B)
INSERT INTO roles (code, name, description, is_system) VALUES
('super_admin', 'Super Administrator', 'Full system governance, manages roles, configs, and emergency state overrides.', true),
('admin', 'Administrator', 'Operational management: approvals, disputes, user management, and analytics.', true),
('support_agent', 'Support Agent', 'Handles disputes and customer support inquiries.', true),
('finance_admin', 'Finance Administrator', 'Manages commission tracking, owner payouts, wallet credit limits, and financial reports.', true),
('owner', 'Fleet Owner', 'Equipment owner who registers fleets, manages operators, sets tariffs, and verifies payments.', true),
('renter', 'Renter', 'Contractor or individual who browses, reserves, and pays for machinery rentals.', true),
('operator', 'Certified Operator', 'Certified field operator operating machinery on project sites.', true),
('org_admin', 'Organization Administrator', 'Full control over corporate settings, members, and approval hierarchies.', true),
('corp_requester', 'Corporate Requester', 'Corporate team member authorized to initiate machinery rental requests.', true),
('corp_approver', 'Corporate Approver', 'Corporate officer authorized to approve or reject internal rental requests.', true),
('fleet_manager', 'Fleet Manager', 'Corporate member managing assigned organizational machinery fleets.', true);

-- 5. Canonical Permissions Seed (SRS Appendix B)
INSERT INTO permissions (code, resource, action, description) VALUES
('dashboard:read', 'dashboard', 'read', 'View administrative operational dashboard'),
('fleets:list', 'fleets', 'list', 'View equipment listings catalog'),
('fleets:create', 'fleets', 'create', 'Register new machinery fleet assets'),
('fleets:approve', 'fleets', 'approve', 'Approve or reject fleet listings in compliance queue'),
('fleets:update', 'fleets', 'update', 'Suspend or activate fleet listings'),
('users:read', 'users', 'read', 'View user profile and identity records'),
('users:update', 'users', 'update', 'Suspend or reactivate user accounts'),
('users:approve', 'users', 'approve', 'Verify user trade documents and national identity records'),
('disputes:list', 'disputes', 'list', 'View customer disputes and support tickets'),
('disputes:resolve', 'disputes', 'resolve', 'Issue legally binding dispute rulings'),
('commissions:read', 'commissions', 'read', 'View commission ledgers and revenue reports'),
('commissions:approve', 'commissions', 'approve', 'Confirm commission payouts and wallet adjustments'),
('templates:manage', 'templates', 'manage', 'Manage dynamic spec templates and category taxonomies'),
('roles:manage', 'roles', 'manage', 'Manage roles and permission grants'),
('config:manage', 'config', 'manage', 'Edit system parameters and multi-country configurations'),
('notifications:manage', 'notifications', 'manage', 'Manage notification templates and delivery gateways'),
('reports:read', 'reports', 'read', 'View advanced platform analytics and financial exports'),
('audit:read', 'audit', 'read', 'View immutable audit log entries'),
('reviews:manage', 'reviews', 'manage', 'Flag and moderate customer ratings and reviews'),
('faq:manage', 'faq', 'manage', 'Manage multilingual FAQ articles and help guides'),
('bookings:create', 'bookings', 'create', 'Create rental bookings'),
('bookings:override', 'bookings', 'override', 'Force administrative state transitions on in-flight bookings'),
('localization:manage', 'localization', 'manage', 'Create and publish dynamic OTA translation bundles');

-- 6. Seed Role Permissions Mapping (Appendix B Matrix)
-- super_admin: All permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p WHERE r.code = 'super_admin';

-- admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.code = 'admin' AND p.code IN (
    'dashboard:read', 'fleets:list', 'fleets:approve', 'fleets:update',
    'users:read', 'users:update', 'users:approve', 'disputes:list', 'disputes:resolve',
    'commissions:read', 'templates:manage', 'reports:read', 'audit:read',
    'reviews:manage', 'faq:manage', 'localization:manage'
);

-- support_agent
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.code = 'support_agent' AND p.code IN ('fleets:list', 'users:read', 'disputes:list');

-- finance_admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.code = 'finance_admin' AND p.code IN (
    'dashboard:read', 'commissions:read', 'commissions:approve', 'reports:read'
);

-- owner
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.code = 'owner' AND p.code IN ('fleets:create', 'fleets:list');

-- renter
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.code = 'renter' AND p.code IN ('fleets:list', 'bookings:create');

-- 7. Seed Baseline Localization Namespaces & Keys
INSERT INTO localization_namespaces (namespace_key, description) VALUES
('auth', 'Authentication and registration viewports'),
('errors', 'Universal RFC 9457 problem details error codes'),
('common', 'Universal buttons, status tags, and measurement units');

-- Auth Keys
INSERT INTO localization_keys (namespace_id, full_key, description)
SELECT n.id, k.full_key, k.descr FROM localization_namespaces n, (VALUES
    ('auth.title_phone_input', 'Phone input screen title'),
    ('auth.subtitle_phone_input', 'Phone input screen description'),
    ('auth.phone_number_label', 'Phone number input field label'),
    ('auth.send_otp_button', 'Send verification code button'),
    ('auth.carrier_badge_prefix', 'Carrier label prefix'),
    ('auth.terms_agreement_notice', 'Legal terms acceptance notice'),
    ('auth.title_otp_verify', 'OTP verification screen heading'),
    ('auth.subtitle_otp_verify', 'OTP verification instructions'),
    ('auth.verify_button', 'Verify OTP code button'),
    ('auth.resend_otp_button', 'Resend OTP button'),
    ('auth.cooldown_remaining', 'Timer label for resend cooldown')
) AS k(full_key, descr) WHERE n.namespace_key = 'auth';

-- Error Keys
INSERT INTO localization_keys (namespace_id, full_key, description)
SELECT n.id, k.full_key, k.descr FROM localization_namespaces n, (VALUES
    ('errors.malformed_json_title', 'Malformed JSON error title'),
    ('errors.malformed_json_detail', 'Malformed JSON error details'),
    ('errors.invalid_phone_title', 'Invalid phone number title'),
    ('errors.invalid_phone_detail', 'Invalid phone number details'),
    ('errors.rate_limit_title', 'Rate limit exceeded title'),
    ('errors.rate_limit_detail', 'Rate limit exceeded details'),
    ('errors.forbidden_permission_title', 'Insufficient permission title'),
    ('errors.forbidden_permission_detail', 'Insufficient permission details'),
    ('errors.unauthorized_title', 'Unauthorized title'),
    ('errors.unauthorized_detail', 'Unauthorized details')
) AS k(full_key, descr) WHERE n.namespace_key = 'errors';

-- Seed English Translations
INSERT INTO localization_translations (key_id, language_code, translated_text, is_verified)
SELECT k.id, 'en', t.text, true FROM localization_keys k
JOIN (VALUES
    ('auth.title_phone_input', 'Enter your phone number'),
    ('auth.subtitle_phone_input', 'We will send a 6-digit SMS verification code to verify your account.'),
    ('auth.phone_number_label', 'Phone Number'),
    ('auth.send_otp_button', 'Send Verification Code'),
    ('auth.carrier_badge_prefix', 'Carrier'),
    ('auth.terms_agreement_notice', 'By proceeding, you agree to our Terms of Service & Privacy Policy.'),
    ('auth.title_otp_verify', 'Verify your phone'),
    ('auth.subtitle_otp_verify', 'Enter the 6-digit verification code sent to'),
    ('auth.verify_button', 'Verify & Continue'),
    ('auth.resend_otp_button', 'Resend Code'),
    ('auth.cooldown_remaining', 'Resend code in'),
    ('errors.malformed_json_title', 'Bad Request'),
    ('errors.malformed_json_detail', 'Request body could not be parsed as valid JSON.'),
    ('errors.invalid_phone_title', 'Invalid Phone Number'),
    ('errors.invalid_phone_detail', 'The provided phone number is invalid for the selected country.'),
    ('errors.rate_limit_title', 'Too Many Requests'),
    ('errors.rate_limit_detail', 'Rate limit exceeded. Please wait before retrying.'),
    ('errors.forbidden_permission_title', 'Forbidden'),
    ('errors.forbidden_permission_detail', 'You do not have the required permissions to perform this action.'),
    ('errors.unauthorized_title', 'Unauthorized'),
    ('errors.unauthorized_detail', 'Authentication credentials are missing, invalid, or expired.')
) AS t(full_key, text) ON k.full_key = t.full_key;

-- Seed Amharic Translations (Verified Idiomatic Ge'ez)
INSERT INTO localization_translations (key_id, language_code, translated_text, is_verified)
SELECT k.id, 'am', t.text, true FROM localization_keys k
JOIN (VALUES
    ('auth.title_phone_input', 'ስልክ ቁጥርዎን ያስገቡ'),
    ('auth.subtitle_phone_input', 'መለያዎን ለማረጋገጥ ባለ 6 አሃዝ የማረጋገጫ ኮድ በኤስኤምኤስ (SMS) እንልካለን።'),
    ('auth.phone_number_label', 'ስልክ ቁጥር'),
    ('auth.send_otp_button', 'የማረጋገጫ ኮድ ላክ'),
    ('auth.carrier_badge_prefix', 'አገልግሎት አቅራቢ'),
    ('auth.terms_agreement_notice', 'በመቀጠልዎ በአገልግሎት ውል እና የግላዊነት መመሪያችን ተስማምተዋል።'),
    ('auth.title_otp_verify', 'ስልክዎን ያረጋግጡ'),
    ('auth.subtitle_otp_verify', 'የተላከውን ባለ 6 አሃዝ የማረጋገጫ ኮድ ያስገቡ ወደ'),
    ('auth.verify_button', 'አረጋግጥና ቀጥል'),
    ('auth.resend_otp_button', 'ኮዱን እንደገና ላክ'),
    ('auth.cooldown_remaining', 'ኮዱን በድጋሚ ለመላክ የቀረው ጊዜ:'),
    ('errors.malformed_json_title', 'ልክ ያልሆነ ጥያቄ'),
    ('errors.malformed_json_detail', 'የተላከው የጥያቄ መረጃ ትክክለኛ የJSON ቅርጸት የለውም።'),
    ('errors.invalid_phone_title', 'ልክ ያልሆነ ስልክ ቁጥር'),
    ('errors.invalid_phone_detail', 'የተሰጠው ስልክ ቁጥር ለተመረጠው ሀገር ትክክል አይደለም።'),
    ('errors.rate_limit_title', 'የሙከራ ገደብ አልፏል'),
    ('errors.rate_limit_detail', 'የተፈቀደው የጥያቄ ገደብ አልፏል። እባክዎን ትንሽ ቆይተው እንደገና ይሞክሩ።'),
    ('errors.forbidden_permission_title', 'የተከለከለ እርምጃ'),
    ('errors.forbidden_permission_detail', 'ይህን ተግባር ለማከናወን የሚያስፈልገው ፈቃድ የለዎትም።'),
    ('errors.unauthorized_title', 'ያልተፈቀደ መዳረሻ'),
    ('errors.unauthorized_detail', 'የመግቢያ ፈቃድዎ አልተገኘም ወይም ጊዜው አልፏል።')
) AS t(full_key, text) ON k.full_key = t.full_key;
