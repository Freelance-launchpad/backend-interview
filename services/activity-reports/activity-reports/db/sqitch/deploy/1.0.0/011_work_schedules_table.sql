-- Deploy activity-reports:1.0.0/011_work_schedules_table to pg

BEGIN;

CREATE TYPE week_type AS ENUM ('mon-fri', 'tue-sat');

CREATE TABLE IF NOT EXISTS work_schedules (
  offer_id         UUID      NOT NULL,
  created_at       TIMESTAMP NOT NULL DEFAULT NOW(),
  application_date DATE      NOT NULL,
  week_type        week_type NOT NULL,

  PRIMARY KEY (offer_id, created_at)
);

COMMIT;
