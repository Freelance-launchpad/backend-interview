-- Verify activity-reports:1.0.0/3_add_other_activity_frequency on pg

BEGIN;

SELECT
   other_activity_frequency
FROM
    activity_reports;

ROLLBACK;
