-- Revert activity-reports:1.0.0/3_add_other_activity_frequency from pg

BEGIN;

ALTER TABLE
    activity_reports
DROP COLUMN IF EXISTS other_activity_frequency;


COMMIT;
