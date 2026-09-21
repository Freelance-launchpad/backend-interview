-- Deploy activity-reports:1.0.0/012_add_month_configuration_columns to pg

BEGIN;

ALTER TABLE activity_reports
  ADD COLUMN IF NOT EXISTS week_type week_type,
  ADD COLUMN IF NOT EXISTS first_day_as_extra_rest BOOLEAN;

COMMIT;
