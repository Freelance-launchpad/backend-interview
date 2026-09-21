-- Verify activity-reports:1.0.0/005_add_offer_id on pg

BEGIN;

SELECT offer_id FROM activity_reports LIMIT 1;

ROLLBACK;
