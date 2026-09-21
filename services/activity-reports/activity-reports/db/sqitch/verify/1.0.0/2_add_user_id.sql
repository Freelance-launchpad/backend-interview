-- Verify activity-reports:1.0.0/2_add_user_id on pg

BEGIN;

SELECT user_id FROM activity_reports;

ROLLBACK;
