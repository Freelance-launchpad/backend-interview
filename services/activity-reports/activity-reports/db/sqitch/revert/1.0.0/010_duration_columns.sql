-- Revert activity-reports:1.0.0/010_duration_columns from pg

BEGIN;

ALTER TABLE activity_reports
    DROP COLUMN IF EXISTS work_duration,
    DROP COLUMN IF EXISTS unit;

DROP TYPE unit;

ALTER TABLE activity_reports
    RENAME COLUMN prospection_duration TO prospection;

ALTER TABLE activity_reports
        RENAME COLUMN formation_duration TO formation;

COMMIT;
