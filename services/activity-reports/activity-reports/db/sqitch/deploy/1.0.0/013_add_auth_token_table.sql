-- Deploy activity-reports:1.0.0/013_add_auth_token_table to pg

BEGIN;

CREATE TABLE IF NOT EXISTS auth_token (
  name        TEXT        PRIMARY KEY,
	token       TEXT        NOT NULL,
  expires_at  TIMESTAMPTZ NOT NULL
);

COMMIT;
