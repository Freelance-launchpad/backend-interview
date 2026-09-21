-- Revert activity-reports:1.0.0/008_remove_user_id from pg

BEGIN;

ALTER TABLE activity_reports
DROP CONSTRAINT unique_report,
ADD COLUMN IF NOT EXISTS user_id TEXT NOT NULL,
ADD CONSTRAINT unique_report UNIQUE (user_id, year, month);

COMMIT;
