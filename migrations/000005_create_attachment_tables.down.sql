-- Rollback Attachment Registration, Compatibility & Inventory Architecture Schema

DELETE FROM role_permissions WHERE permission_id IN (
    SELECT id FROM permissions WHERE code LIKE 'attachments:%'
);

DELETE FROM permissions WHERE code LIKE 'attachments:%';

DROP TRIGGER IF EXISTS trg_fleet_bundled_attachments_updated_at ON fleet_bundled_attachments;
DROP TRIGGER IF EXISTS trg_attachments_updated_at ON attachments;

DROP TABLE IF EXISTS attachment_reservations;
DROP TABLE IF EXISTS fleet_bundled_attachments;
DROP TABLE IF EXISTS attachment_photos;
DROP TABLE IF EXISTS attachment_compatible_categories;
DROP TABLE IF EXISTS attachments;
