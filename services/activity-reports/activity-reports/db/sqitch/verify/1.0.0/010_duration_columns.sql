-- Verify activity-reports:1.0.0/010_duration_columns on pg

BEGIN;

SELECT work_duration, unit, prospection_duration, formation_duration
FROM activity_reports LIMIT 1;

ROLLBACK;
