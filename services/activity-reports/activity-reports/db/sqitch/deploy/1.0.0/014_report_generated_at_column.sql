-- Deploy activity-reports:1.0.0/014_report_generated_at_column to pg

BEGIN;

ALTER TABLE activity_reports
    ADD COLUMN IF NOT EXISTS generated_at TIMESTAMPTZ;

UPDATE activity_reports SET generated_at = created_at WHERE generated_at IS NULL;

COMMIT;
