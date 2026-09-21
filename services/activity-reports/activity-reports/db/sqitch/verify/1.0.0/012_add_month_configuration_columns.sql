-- Verify activity-reports:1.0.0/012_add_month_configuration_columns on pg

BEGIN;

SELECT week_type, first_day_as_extra_rest FROM activity_reports LIMIT 1;

ROLLBACK;
