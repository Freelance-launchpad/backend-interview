-- Revert activity-reports:1.0.0/3_add_green_reports from pg

BEGIN;

ALTER TABLE
    activity_reports
RENAME COLUMN prospection TO prospection_days;

ALTER TABLE
    activity_reports
RENAME COLUMN formation TO formation_days;

DROP TABLE activity_types;

ALTER TABLE
    activity_reports
DROP COLUMN IF EXISTS activity_type;

ALTER TABLE
    activity_reports
DROP COLUMN IF EXISTS other_activity;

ALTER TABLE
    report_items
RENAME COLUMN amount TO nb_days;

COMMIT;