-- Revert Rental Completion, Demobilization, Damage Triage & Settlement Reconciliation Schema

DROP TABLE IF EXISTS booking_reviews CASCADE;
DROP TABLE IF EXISTS rental_settlements CASCADE;
DROP TABLE IF EXISTS damage_claim_dockets CASCADE;
DROP TABLE IF EXISTS return_inspections CASCADE;

ALTER TABLE bookings DROP COLUMN IF EXISTS damage_claim_submitted;
ALTER TABLE bookings DROP COLUMN IF EXISTS damage_claim_waived;
ALTER TABLE bookings DROP COLUMN IF EXISTS accepted_damage_minor;
ALTER TABLE bookings DROP COLUMN IF EXISTS excess_hours_surcharge_minor;
ALTER TABLE bookings DROP COLUMN IF EXISTS final_net_balance_minor;

-- Restore bookings status constraint to migration 000008 set
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
    'PENDING_FINAL_PAYMENT',
    'COMPLETED',
    'RATED',
    'CANCELLED'
));

DELETE FROM role_permissions WHERE permission_id IN (
    SELECT id FROM permissions WHERE code IN (
        'rentals:return_initiate',
        'rentals:return_inspect',
        'rentals:damage_claim',
        'rentals:damage_respond',
        'rentals:arbitrate',
        'rentals:settle',
        'rentals:review'
    )
);

DELETE FROM permissions WHERE code IN (
    'rentals:return_initiate',
    'rentals:return_inspect',
    'rentals:damage_claim',
    'rentals:damage_respond',
    'rentals:arbitrate',
    'rentals:settle',
    'rentals:review'
);
