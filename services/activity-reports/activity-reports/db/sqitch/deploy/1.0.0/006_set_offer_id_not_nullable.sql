-- Deploy activity-reports:1.0.0/006_set_offer_id_not_nullable to pg

BEGIN;

ALTER TABLE activity_reports ALTER COLUMN offer_id SET NOT NULL;

COMMIT;
