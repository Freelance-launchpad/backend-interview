-- Revert activity-reports:1.0.0/014_report_generated_at_column from pg

BEGIN;

ALTER TABLE activity_reports
    DROP COLUMN IF EXISTS generated_at;

COMMIT;
