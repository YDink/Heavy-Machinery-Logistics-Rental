-- Revert Active Rental Management, Shift Telemetry, Breakdown & Extension Schema

DROP TABLE IF EXISTS mechanical_breakdowns CASCADE;
DROP TABLE IF EXISTS rental_extensions CASCADE;
DROP TABLE IF EXISTS rental_shifts CASCADE;
DROP TABLE IF EXISTS handover_challenges CASCADE;
DROP TABLE IF EXISTS handover_inspections CASCADE;

ALTER TABLE bookings DROP COLUMN IF EXISTS assigned_operator_id;

DELETE FROM role_permissions WHERE permission_id IN (
    SELECT id FROM permissions WHERE code LIKE 'rentals:%'
);

DELETE FROM permissions WHERE code LIKE 'rentals:%';
