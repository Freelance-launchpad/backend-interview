-- Revert activity-reports:1.0.0/005_add_offer_id from pg

BEGIN;

ALTER TABLE activity_reports DROP COLUMN offer_id;

COMMIT;
