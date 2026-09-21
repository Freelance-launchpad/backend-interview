-- Verify activity-reports:1.0.0/013_add_auth_token_table on pg

BEGIN;

SELECT name, token, expires_at FROM auth_token LIMIT 1;

ROLLBACK;
