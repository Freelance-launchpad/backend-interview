-- Revert activity-reports:1.0.0/1_create_activity-report_table from pg

BEGIN;

DROP TABLE activity_reports;

DROP TABLE report_items;

COMMIT;
