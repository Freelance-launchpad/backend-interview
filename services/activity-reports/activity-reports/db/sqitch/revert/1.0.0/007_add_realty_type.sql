-- Revert activity-reports:1.0.0/007_add_realty_type from pg

BEGIN;

DELETE FROM activity_types WHERE type = 'realty';

COMMIT;
