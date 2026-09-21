-- Deploy activity-reports:1.0.0/005_add_offer_id to pg

BEGIN;

ALTER TABLE activity_reports ADD COLUMN offer_id UUID NULL;

COMMIT;
