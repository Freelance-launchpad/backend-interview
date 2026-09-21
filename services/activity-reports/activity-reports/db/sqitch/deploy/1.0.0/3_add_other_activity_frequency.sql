-- Deploy activity-reports:1.0.0/3_add_other_activity_frequency to pg

BEGIN;

ALTER TABLE
    activity_reports
ADD COLUMN IF NOT EXISTS other_activity_frequency TEXT;


COMMIT;
