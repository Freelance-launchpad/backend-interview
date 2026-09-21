-- Verify activity-reports:1.0.0/3_add_green_reports on pg

BEGIN;

SELECT
    prospection, formation, activity_type, other_activity
FROM
    activity_reports;

SELECT
    amount
FROM
    report_items;

ROLLBACK;
