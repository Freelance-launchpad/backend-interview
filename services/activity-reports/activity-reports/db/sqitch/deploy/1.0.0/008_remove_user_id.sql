-- Deploy activity-reports:1.0.0/008_remove_user_id to pg

BEGIN;

ALTER TABLE activity_reports
DROP CONSTRAINT unique_report,
DROP COLUMN user_id,
ADD CONSTRAINT unique_report UNIQUE (offer_id, year, month);

COMMIT;
