-- Deploy activity-reports:1.0.0/015_add_formation to pg
BEGIN;

INSERT INTO
  activity_types (type)
VALUES
  ('formation');

COMMIT;
