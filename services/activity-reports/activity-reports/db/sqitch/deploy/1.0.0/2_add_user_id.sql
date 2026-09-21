-- Deploy activity-reports:1.0.0/2_add_user_id to pg

BEGIN;

ALTER TABLE
    activity_reports
ADD COLUMN IF NOT EXISTS user_id TEXT NOT NULL,
ADD CONSTRAINT unique_report          UNIQUE (user_id, year, month);

COMMIT;
