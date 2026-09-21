-- Revert activity-reports:1.0.0/011_work_schedules_table from pg

BEGIN;

DROP TABLE IF EXISTS work_schedules;

DROP TYPE IF EXISTS week_type;

COMMIT;
