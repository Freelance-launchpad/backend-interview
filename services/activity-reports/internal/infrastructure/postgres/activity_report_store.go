package postgres

import (
	"context"
	"fmt"

	"github.com/Freelance-launchpad/backend-interview/common/jdate"
	"github.com/Freelance-launchpad/backend-interview/common/jdb"
	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/domain"
	transactor "github.com/Thiht/transactor/stdlib"
	"github.com/google/uuid"
)

type activityReportStore struct {
	dbGetter transactor.DBGetter
}

func NewActivityReportStore(dbGetter transactor.DBGetter) domain.ActivityReportStore {
	return &activityReportStore{
		dbGetter: dbGetter,
	}
}

const queryCreateActivityReport = `
INSERT INTO activity_reports (
	offer_id,
	year,
	month,
	activity_type,
	week_type,
	first_day_as_extra_rest,
	work_duration,
	prospection_duration,
	formation_duration,
	unit,
	vacation_days,
	days_away,
	other_activity,
	other_activity_frequency
)
VALUES
	($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`

func (s *activityReportStore) CreateActivityReport(ctx context.Context, report domain.Report) (err error) {
	defer jerror.Wrap(&err)

	_, err = s.dbGetter(ctx).ExecContext(ctx, queryCreateActivityReport,
		report.OfferID,
		report.Year,
		report.Month,
		report.ActivityType,
		report.WeekType,
		report.FirstDayAsExtraRest,
		report.WorkDuration,
		report.ProspectionDuration,
		report.FormationDuration,
		report.Unit,
		report.VacationDays,
		report.DaysAway,
		report.OtherActivity,
		report.OtherActivityFrequency,
	)
	if err != nil {
		if jdb.IsErrorUniqueViolation(err) {
			return domain.ErrReportAlreadyExists
		}
		return err
	}
	return nil
}

const queryGetActivityReportsByOfferID = `
SELECT
	id,
	offer_id,
	year,
	month,
	activity_type,
	week_type,
	first_day_as_extra_rest,
	work_duration,
	prospection_duration,
	formation_duration,
	unit,
	vacation_days,
	days_away,
	other_activity,
	other_activity_frequency,
	created_at
FROM
	activity_reports
WHERE
	offer_id = $1
ORDER BY
	year, month`

func (s *activityReportStore) GetActivityReportsByOfferID(ctx context.Context, offerID uuid.UUID) (reports []domain.Report, err error) {
	defer jerror.Wrap(&err)

	rows, err := s.dbGetter(ctx).QueryContext(ctx, queryGetActivityReportsByOfferID, offerID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var report domain.Report
		if err := rows.Scan(
			&report.ID,
			&report.OfferID,
			&report.Year,
			&report.Month,
			&report.ActivityType,
			&report.WeekType,
			&report.FirstDayAsExtraRest,
			&report.WorkDuration,
			&report.ProspectionDuration,
			&report.FormationDuration,
			&report.Unit,
			&report.VacationDays,
			&report.DaysAway,
			&report.OtherActivity,
			&report.OtherActivityFrequency,
			&report.CreatedAt,
		); err != nil {
			return nil, err
		}

		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reports, nil
}

const queryGetActivityReportsToGenerate = `
SELECT
	id,
	offer_id,
	year,
	month,
	activity_type,
	week_type,
	first_day_as_extra_rest,
	work_duration,
	prospection_duration,
	formation_duration,
	unit,
	vacation_days,
	days_away,
	other_activity,
	other_activity_frequency,
	created_at
FROM
	activity_reports
WHERE
	generated_at IS NULL
FOR UPDATE`

func (s *activityReportStore) GetActivityReportsToGenerate(ctx context.Context) (reports []domain.Report, err error) {
	defer jerror.Wrap(&err)

	rows, err := s.dbGetter(ctx).QueryContext(ctx, queryGetActivityReportsToGenerate)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var report domain.Report
		if err := rows.Scan(
			&report.ID,
			&report.OfferID,
			&report.Year,
			&report.Month,
			&report.ActivityType,
			&report.WeekType,
			&report.FirstDayAsExtraRest,
			&report.WorkDuration,
			&report.ProspectionDuration,
			&report.FormationDuration,
			&report.Unit,
			&report.VacationDays,
			&report.DaysAway,
			&report.OtherActivity,
			&report.OtherActivityFrequency,
			&report.CreatedAt,
		); err != nil {
			return nil, err
		}

		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reports, nil
}

const queryGetMissingActivityReportsByOfferID = `
SELECT
	extract(year from period.i) as year,
	extract(month from period.i) as month
FROM
	generate_series('%s', '%s', INTERVAL '1 MONTH') period(i)
EXCEPT
	(SELECT year, month FROM activity_reports WHERE offer_id = $1)
ORDER BY year, month`

func (s *activityReportStore) GetMissingActivityReportsByOfferID(ctx context.Context, offerID uuid.UUID, from, to jdate.YearMonth) (missingPeriods []jdate.YearMonth, err error) {
	defer jerror.Wrap(&err)

	start := fmt.Sprintf("%d-%02d-01", from.Year, from.Month)
	end := fmt.Sprintf("%d-%02d-01", to.Year, to.Month)
	query := fmt.Sprintf(queryGetMissingActivityReportsByOfferID, start, end)

	rows, err := s.dbGetter(ctx).QueryContext(ctx, query, offerID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var period jdate.YearMonth
		if err := rows.Scan(&period.Year, &period.Month); err != nil {
			return nil, err
		}
		missingPeriods = append(missingPeriods, period)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return missingPeriods, nil
}

const querySetReportAsGenerated = `
UPDATE
	activity_reports
SET
	generated_at = NOW()
WHERE
	id = $1`

func (s *activityReportStore) SetReportAsGenerated(ctx context.Context, id uuid.UUID) (err error) {
	defer jerror.Wrap(&err)

	if _, err := s.dbGetter(ctx).ExecContext(ctx, querySetReportAsGenerated, id); err != nil {
		return err
	}

	return nil
}
