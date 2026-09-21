package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Freelance-launchpad/backend-interview/common/jdate"
	"github.com/Freelance-launchpad/backend-interview/common/jentity"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/domain"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/pkg"
	transactor "github.com/Thiht/transactor/stdlib"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_activityReportStore_CreateActivityReport(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	report := domain.Report{
		OfferID:             uuid.New(),
		Year:                2025,
		Month:               5,
		ActivityType:        jentity.EntityBlue,
		WeekType:            new(pkg.WeekTypeTueSat),
		FirstDayAsExtraRest: new(true),
		WorkDuration:        12.34,
		ProspectionDuration: 34.56,
		FormationDuration:   56.78,
		Unit:                domain.UnitDays,
		VacationDays:        2,
		DaysAway:            1,
	}

	tests := map[string]struct {
		expect  func(sqlmock.Sqlmock)
		wantErr error
	}{
		"nominal": {
			expect: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(queryCreateActivityReport).
					WithArgs(
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
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		"conflict": {
			expect: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(queryCreateActivityReport).
					WithArgs(
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
					).
					WillReturnError(&pq.Error{Code: "23505"})
			},
			wantErr: domain.ErrReportAlreadyExists,
		},
		"error": {
			expect: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(queryCreateActivityReport).
					WithArgs(
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
					).
					WillReturnError(assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			require.NoError(t, err)
			defer db.Close()
			_, dbGetter := transactor.NewFakeTransactor(db)

			if tt.expect != nil {
				tt.expect(mock)
			}

			s := activityReportStore{
				dbGetter: dbGetter,
			}

			err = s.CreateActivityReport(ctx, report)
			assert.ErrorIs(t, err, tt.wantErr)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_activityReportStore_GetActivityReportsByOfferID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	offerID := uuid.New()

	report := domain.Report{
		ID:                  uuid.New(),
		OfferID:             offerID,
		Year:                2025,
		Month:               5,
		ActivityType:        jentity.EntityBlue,
		WorkDuration:        12.34,
		ProspectionDuration: 34.56,
		FormationDuration:   56.78,
		Unit:                domain.UnitDays,
		VacationDays:        2,
		DaysAway:            1,
		CreatedAt:           time.Now(),
	}

	tests := map[string]struct {
		expect  func(sqlmock.Sqlmock)
		want    []domain.Report
		wantErr error
	}{
		"nominal": {
			expect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(queryGetActivityReportsByOfferID).
					WithArgs(offerID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "offer_id", "year", "month", "activity_type",
						"week_type", "first_day_as_extra_rest",
						"work_duration", "prospection_duration", "formation_duration", "unit",
						"vacation_days", "days_away", "other_activity", "other_activity_frequency", "created_at"}).
						AddRow(
							report.ID, report.OfferID, report.Year, report.Month, report.ActivityType,
							report.WeekType, report.FirstDayAsExtraRest,
							report.WorkDuration, report.ProspectionDuration, report.FormationDuration, report.Unit,
							report.VacationDays, report.DaysAway, report.OtherActivity, report.OtherActivityFrequency, report.CreatedAt,
						))
			},
			want: []domain.Report{
				report,
			},
		},
		"error": {
			expect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(queryGetActivityReportsByOfferID).
					WithArgs(offerID).
					WillReturnError(assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			require.NoError(t, err)
			defer db.Close()
			_, dbGetter := transactor.NewFakeTransactor(db)

			if tt.expect != nil {
				tt.expect(mock)
			}

			s := activityReportStore{
				dbGetter: dbGetter,
			}

			got, err := s.GetActivityReportsByOfferID(ctx, offerID)
			assert.ElementsMatch(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_activityReportStore_GetActivityReportsToGenerate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	offerID := uuid.New()

	report := domain.Report{
		ID:                  uuid.New(),
		OfferID:             offerID,
		Year:                2025,
		Month:               5,
		ActivityType:        jentity.EntityBlue,
		WorkDuration:        12.34,
		ProspectionDuration: 34.56,
		FormationDuration:   56.78,
		Unit:                domain.UnitDays,
		VacationDays:        2,
		DaysAway:            1,
		CreatedAt:           time.Now(),
	}

	tests := map[string]struct {
		expect  func(sqlmock.Sqlmock)
		want    []domain.Report
		wantErr error
	}{
		"nominal": {
			expect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(queryGetActivityReportsToGenerate).
					WillReturnRows(sqlmock.NewRows([]string{"id", "offer_id", "year", "month", "activity_type",
						"week_type", "first_day_as_extra_rest",
						"work_duration", "prospection_duration", "formation_duration", "unit",
						"vacation_days", "days_away", "other_activity", "other_activity_frequency", "created_at"}).
						AddRow(
							report.ID, report.OfferID, report.Year, report.Month, report.ActivityType,
							report.WeekType, report.FirstDayAsExtraRest,
							report.WorkDuration, report.ProspectionDuration, report.FormationDuration, report.Unit,
							report.VacationDays, report.DaysAway, report.OtherActivity, report.OtherActivityFrequency, report.CreatedAt,
						))
			},
			want: []domain.Report{
				report,
			},
		},
		"error": {
			expect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(queryGetActivityReportsToGenerate).
					WillReturnError(assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			require.NoError(t, err)
			defer db.Close()
			_, dbGetter := transactor.NewFakeTransactor(db)

			if tt.expect != nil {
				tt.expect(mock)
			}

			s := activityReportStore{
				dbGetter: dbGetter,
			}

			got, err := s.GetActivityReportsToGenerate(ctx)
			assert.ElementsMatch(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_activityReportStore_GetMissingActivityReportsByOfferID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	offerID := uuid.New()

	from := jdate.YearMonth{Year: 2025, Month: 1}
	to := jdate.YearMonth{Year: 2025, Month: 9}

	query := fmt.Sprintf(queryGetMissingActivityReportsByOfferID, "2025-01-01", "2025-09-01")

	tests := map[string]struct {
		expect  func(sqlmock.Sqlmock)
		want    []jdate.YearMonth
		wantErr error
	}{
		"nominal": {
			expect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).
					WithArgs(offerID).
					WillReturnRows(sqlmock.NewRows([]string{"year", "month"}).
						AddRow(2025, 8))
			},
			want: []jdate.YearMonth{
				{Year: 2025, Month: 8},
			},
		},
		"error": {
			expect: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).
					WithArgs(offerID).
					WillReturnError(assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			require.NoError(t, err)
			defer db.Close()
			_, dbGetter := transactor.NewFakeTransactor(db)

			if tt.expect != nil {
				tt.expect(mock)
			}

			s := activityReportStore{
				dbGetter: dbGetter,
			}

			got, err := s.GetMissingActivityReportsByOfferID(ctx, offerID, from, to)
			assert.ElementsMatch(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_activityReportStore_SetReportAsGenerated(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	id := uuid.New()

	tests := map[string]struct {
		expect  func(sqlmock.Sqlmock)
		wantErr error
	}{
		"nominal": {
			expect: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(querySetReportAsGenerated).
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		"error": {
			expect: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(querySetReportAsGenerated).
					WithArgs(id).
					WillReturnError(assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			require.NoError(t, err)
			defer db.Close()
			_, dbGetter := transactor.NewFakeTransactor(db)

			if tt.expect != nil {
				tt.expect(mock)
			}

			s := activityReportStore{
				dbGetter: dbGetter,
			}

			err = s.SetReportAsGenerated(ctx, id)
			assert.ErrorIs(t, err, tt.wantErr)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
