-- Verify activity-reports:1.0.0/007_add_realty_type on pg

BEGIN;

SELECT 1 / (SELECT COUNT(*) FROM activity_types WHERE type = 'realty');

ROLLBACK;
