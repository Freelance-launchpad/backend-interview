-- Verify activity-reports:1.0.0/011_work_schedules_table on pg

BEGIN;

SELECT offer_id, created_at, application_date, week_type FROM work_schedules LIMIT 1;

ROLLBACK;
