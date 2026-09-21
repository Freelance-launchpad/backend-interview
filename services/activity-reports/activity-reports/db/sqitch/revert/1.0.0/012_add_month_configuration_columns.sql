-- Revert activity-reports:1.0.0/012_add_month_configuration_columns from pg

BEGIN;

ALTER TABLE activity_reports
    DROP COLUMN IF EXISTS week_type,
    DROP COLUMN IF EXISTS first_day_as_extra_rest;

COMMIT;
