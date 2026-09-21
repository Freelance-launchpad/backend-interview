-- Revert activity-reports:1.0.0/2_add_user_id from pg

BEGIN;

ALTER TABLE
    activity_reports
DROP COLUMN IF EXISTS user_id;

COMMIT;
