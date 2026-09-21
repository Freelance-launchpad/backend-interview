-- Verify activity-reports:1.0.0/014_report_generated_at_column on pg

BEGIN;

SELECT generated_at FROM activity_reports LIMIT 1;

ROLLBACK;
