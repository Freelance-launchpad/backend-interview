-- Deploy activity-reports:1.0.0/3_add_green_reports to pg

BEGIN;

ALTER TABLE
    activity_reports
RENAME COLUMN prospection_days TO prospection;

ALTER TABLE
    activity_reports
RENAME COLUMN formation_days TO formation;

CREATE TABLE activity_types (
    type  TEXT PRIMARY KEY
);

INSERT INTO activity_types (type) VALUES
    ('blue'),
    ('green');

ALTER TABLE
    activity_reports
ADD COLUMN IF NOT EXISTS activity_type TEXT REFERENCES activity_types (type) NOT NULL DEFAULT 'blue';

ALTER TABLE
    activity_reports
ADD COLUMN IF NOT EXISTS other_activity DECIMAL DEFAULT 0;

ALTER TABLE
    report_items
RENAME COLUMN nb_days TO amount;

COMMIT;
