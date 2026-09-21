-- Revert activity-reports:1.0.0/015_add_formation from pg
BEGIN;

DELETE FROM activity_types
WHERE
  type = 'formation';

COMMIT;
