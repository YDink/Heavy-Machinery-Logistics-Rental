-- Rollback Booking, Reservation & Lifecycle State Machine Schema (Partition 05-05)

-- 1. Remove RBAC Permissions
DELETE FROM role_permissions WHERE permission_id IN (
    SELECT id FROM permissions WHERE module = 'bookings'
);

DELETE FROM permissions WHERE module = 'bookings';

-- 2. Drop Tables
DROP TABLE IF EXISTS booking_audit_logs CASCADE;
DROP TABLE IF EXISTS booking_bank_slips CASCADE;
DROP TABLE IF EXISTS booking_attachments CASCADE;
DROP TABLE IF EXISTS bookings CASCADE;
