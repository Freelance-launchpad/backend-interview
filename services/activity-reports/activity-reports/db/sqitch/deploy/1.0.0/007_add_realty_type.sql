-- Deploy activity-reports:1.0.0/007_add_realty_type to pg

BEGIN;

INSERT INTO activity_types (type) VALUES ('realty');

COMMIT;
