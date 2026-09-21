-- Verify activity-reports:1.0.0/015_add_formation on pg
BEGIN;

SELECT
  1 / (
    SELECT
      COUNT(*)
    FROM
      activity_types
    WHERE
      type = 'formation'
  );

ROLLBACK;
