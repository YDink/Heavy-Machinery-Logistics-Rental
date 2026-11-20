-- Revert Logistics Transport Integration, Sub-Bookings & Reverse Cascades Schema

DROP TABLE IF EXISTS logistics_reverse_cascades CASCADE;
DROP TABLE IF EXISTS logistics_transshipments CASCADE;
DROP TABLE IF EXISTS external_transports CASCADE;
DROP TABLE IF EXISTS logistics_dispatches CASCADE;

ALTER TABLE country_configs DROP COLUMN IF EXISTS max_gvm_kg;
ALTER TABLE country_configs DROP COLUMN IF EXISTS max_corridor_height_mm;
ALTER TABLE country_configs DROP COLUMN IF EXISTS wide_load_width_mm;

ALTER TABLE fleet_spec_indexed DROP COLUMN IF EXISTS deck_well_length_mm;
ALTER TABLE fleet_spec_indexed DROP COLUMN IF EXISTS deck_width_mm;
ALTER TABLE fleet_spec_indexed DROP COLUMN IF EXISTS outriggers_extension_mm;
ALTER TABLE fleet_spec_indexed DROP COLUMN IF EXISTS deck_height_ground_mm;
ALTER TABLE fleet_spec_indexed DROP COLUMN IF EXISTS tare_weight_kg;
ALTER TABLE fleet_spec_indexed DROP COLUMN IF EXISTS axle_count;

DROP INDEX IF EXISTS idx_bookings_parent_hierarchy;
ALTER TABLE bookings DROP CONSTRAINT IF EXISTS chk_bookings_kind;
ALTER TABLE bookings DROP CONSTRAINT IF EXISTS chk_bookings_commission;
ALTER TABLE bookings DROP COLUMN IF EXISTS parent_booking_id;
ALTER TABLE bookings DROP COLUMN IF EXISTS booking_kind;
ALTER TABLE bookings DROP COLUMN IF EXISTS platform_commission_rate_bps;
ALTER TABLE bookings DROP COLUMN IF EXISTS platform_commission_minor;
ALTER TABLE bookings DROP COLUMN IF EXISTS abnormal_load_escort_required;

-- Restore bookings status constraint to migration 000009 set
ALTER TABLE bookings DROP CONSTRAINT IF EXISTS chk_bookings_status;
ALTER TABLE bookings ADD CONSTRAINT chk_bookings_status CHECK (status IN (
    'DRAFT',
    'PENDING_PAYMENT',
    'PAYMENT_EVIDENCE_UPLOADED',
    'PAYMENT_DECLINED',
    'DISPUTE_OPENED',
    'DISPUTE_RESOLVED_RENTER',
    'DISPUTE_RESOLVED_OWNER',
    'CONFIRMED',
    'AWAITING_FLEET_READINESS',
    'PENDING_ADDITIONAL_PAYMENT',
    'READY_FOR_PICKUP',
    'ACTIVE',
    'ACTIVE_BREAKDOWN_PAUSED',
    'ACTIVE_EXTENSION_HOLD',
    'RETURN_INITIATED',
    'DAMAGE_PENDING_REVIEW',
    'DAMAGE_DISPUTED',
    'PENDING_FINAL_PAYMENT',
    'COMPLETED',
    'RATED',
    'CANCELLED'
));

DELETE FROM role_permissions WHERE permission_id IN (
    SELECT id FROM permissions WHERE code IN (
        'logistics:match',
        'logistics:sub_booking',
        'logistics:external_transport',
        'logistics:dispatch',
        'logistics:transshipment',
        'logistics:reverse_cascade'
    )
);

DELETE FROM permissions WHERE code IN (
    'logistics:match',
    'logistics:sub_booking',
    'logistics:external_transport',
    'logistics:dispatch',
    'logistics:transshipment',
    'logistics:reverse_cascade'
);
