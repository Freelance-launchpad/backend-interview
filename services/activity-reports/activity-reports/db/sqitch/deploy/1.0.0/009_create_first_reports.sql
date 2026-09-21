-- Deploy activity-reports:1.0.0/009_create_first_reports to pg

BEGIN;

WITH first_report AS (
  SELECT
    offer_id,
    MIN(TO_DATE(year || '-' || month || '-01', 'YYYY-MM-DD')) as first_report_period,
    MIN(activity_type) as activity_type
  FROM public.activity_reports
  GROUP BY offer_id
)
INSERT INTO public.activity_reports (month, year, prospection, vacation_days, days_away, formation, created_at, activity_type, other_activity, other_activity_frequency, offer_id)
SELECT
  (EXTRACT(MONTH FROM pmd.first_report_period - INTERVAL '1 month'))::INTEGER,
  (EXTRACT(YEAR FROM pmd.first_report_period - INTERVAL '1 month'))::INTEGER,
  0,
  0,
  0,
  0,
  NOW(),
  pmd.activity_type,
  NULL,
  NULL,
  pmd.offer_id
FROM first_report pmd;

COMMIT;
