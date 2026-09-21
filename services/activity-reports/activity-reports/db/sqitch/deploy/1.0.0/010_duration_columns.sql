-- Deploy activity-reports:1.0.0/010_duration_columns to pg

BEGIN;

CREATE TYPE unit AS ENUM ('days', 'hours');

ALTER TABLE activity_reports
    ADD COLUMN IF NOT EXISTS work_duration DECIMAL NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS unit unit;

UPDATE activity_reports
    SET unit = 'hours' WHERE activity_type = 'green' OR activity_type = 'realty';
UPDATE activity_reports
    SET unit = 'days' WHERE activity_type = 'blue';

ALTER TABLE activity_reports ALTER COLUMN unit SET NOT NULL;

ALTER TABLE activity_reports
    RENAME COLUMN prospection TO prospection_duration;

ALTER TABLE activity_reports
    RENAME COLUMN formation TO formation_duration;

COMMIT;
