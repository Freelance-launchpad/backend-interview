-- Revert activity-reports:1.0.0/009_create_first_reports from pg

BEGIN;

WITH first_report AS (
  SELECT
    id,
    MIN(TO_DATE(year || '-' || month || '-01', 'YYYY-MM-DD')) as first_report_period
  FROM public.activity_reports
  GROUP BY offer_id
)
DELETE FROM public.activity_reports WHERE id IN (SELECT id FROM first_report);

COMMIT;
