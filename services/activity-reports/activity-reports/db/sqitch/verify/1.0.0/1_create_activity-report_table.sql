-- Verify activity-reports:1.0.0/1_create_activity-report_table on pg

BEGIN;

SELECT
    id,
    month,
    year,
    prospection_days,
    vacation_days,
    days_away,
    formation_days,
    created_at
FROM
    activity_reports;

SELECT
    id,
    mission_id,
    report_id,
    nb_days,
    progress_report
FROM
    report_items;

ROLLBACK;
