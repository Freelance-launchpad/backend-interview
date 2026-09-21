-- Revert activity-reports:1.0.0/013_add_auth_token_table from pg

BEGIN;

DROP TABLE IF EXISTS auth_token;

COMMIT;
