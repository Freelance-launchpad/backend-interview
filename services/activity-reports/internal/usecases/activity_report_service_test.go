package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/Freelance-launchpad/backend-interview/common/jclock"
	"github.com/Freelance-launchpad/backend-interview/common/jdate"
	jdbtesting "github.com/Freelance-launchpad/backend-interview/common/jdb/testing"
	"github.com/Freelance-launchpad/backend-interview/common/jentity"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/domain"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/mocks"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/pkg"
	pkgusers "github.com/Freelance-launchpad/backend-interview/services/users/pkg"
	pkgusersmocks "github.com/Freelance-launchpad/backend-interview/services/users/pkg/client/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	arg "github.com/stretchr/testify/mock"
)

func TestService_CreateActivityReportsV3(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	offerID := uuid.New()

	report := domain.Report{
		Month:               3,
		Year:                2023,
		ProspectionDuration: 0.5,
		VacationDays:        1,
		DaysAway:            0,
		FormationDuration:   1.5,
	}

	reportFilled := report
	reportFilled.OfferID = offerID
	reportFilled.ActivityType = jentity.EntityGreen
	reportFilled.Unit = domain.UnitHours
	reportFilled.WeekType = new(pkg.WeekTypeTueSat)
	reportFilled.FirstDayAsExtraRest = new(true)

	reports := []domain.Report{report}
	reportsFilled := []domain.Report{reportFilled}

	offerWithContract := pkgusers.OfferWithContract{
		Offer: pkgusers.Offer{
			ID:     offerID,
			UserID: "bibou",
			Type:   jentity.OfferTypeLife,
		},
		Contract: &pkgusers.Contract{
			Entity: jentity.EntityGreen,
		},
	}

	tests := []struct {
		name                string
		offerID             uuid.UUID
		reports             []domain.Report
		activityReportStore func(m *mocks.ActivityReportStore)
		usersClient         func(*pkgusersmocks.Client)
		workScheduleService func(m *mocks.WorkScheduleService)
		wantErr             error
	}{
		{
			name:    "nominal",
			offerID: offerID,
			reports: reports,
			usersClient: func(m *pkgusersmocks.Client) {
				m.EXPECT().GetOfferWithContract(ctx, offerID).Return(offerWithContract, nil)
			},
			activityReportStore: func(m *mocks.ActivityReportStore) {
				m.EXPECT().CreateActivityReport(jdbtesting.Context(ctx), reportsFilled[0]).Return(nil)
			},
		},
		{
			name:    "GetOfferWithContract failed",
			offerID: offerID,
			reports: reports,
			usersClient: func(m *pkgusersmocks.Client) {
				m.EXPECT().GetOfferWithContract(ctx, offerID).Return(pkgusers.OfferWithContract{}, assert.AnError)
			},
			wantErr: assert.AnError,
		},
		{
			name:    "GetMonthConfiguration failed",
			reports: reports,
			offerID: offerID,
			usersClient: func(m *pkgusersmocks.Client) {
				m.EXPECT().GetOfferWithContract(ctx, offerID).Return(offerWithContract, nil)
			},
			wantErr: assert.AnError,
		},
		{
			name:    "CreateActivityReport failed",
			reports: reports,
			offerID: offerID,
			usersClient: func(m *pkgusersmocks.Client) {
				m.EXPECT().GetOfferWithContract(ctx, offerID).Return(offerWithContract, nil)
			},
			activityReportStore: func(m *mocks.ActivityReportStore) {
				m.EXPECT().CreateActivityReport(jdbtesting.Context(ctx), reportsFilled[0]).Return(assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			activityReportStore := mocks.NewActivityReportStore(t)
			if tt.activityReportStore != nil {
				tt.activityReportStore(activityReportStore)
			}
			usersClient := pkgusersmocks.NewClient(t)
			if tt.usersClient != nil {
				tt.usersClient(usersClient)
			}
			workScheduleService := mocks.NewWorkScheduleService(t)
			if tt.workScheduleService != nil {
				tt.workScheduleService(workScheduleService)
			}
			transactor := jdbtesting.NewFakeTransactor()

			s := &ActivityReportService{
				reportStore: activityReportStore,
				usersClient: usersClient,
				transactor:  transactor,
			}

			err := s.CreateActivityReportsV3(context.Background(), tt.offerID, tt.reports)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestService_GetMissingActivityReports(t *testing.T) {
	t.Parallel()
	offerID := uuid.New()

	tests := []struct {
		name        string
		offerID     uuid.UUID
		clock       jclock.Clock
		expectStore func(m *mocks.ActivityReportStore)
		expectUsers func(m *pkgusersmocks.Client)
		want        []jdate.YearMonth
		wantErr     error
	}{
		{
			name:    "contract started 2 months ago - before reporting day - should ignore current and previous month",
			offerID: offerID,
			clock:   jclock.NewFakeClock(time.Date(2023, time.May, 4, 10, 0, 0, 0, time.UTC)),
			expectUsers: func(m *pkgusersmocks.Client) {
				m.EXPECT().GetOfferJobContract(arg.Anything, offerID).Return(pkgusers.Contract{
					StartDate: jdate.New(2023, 3, 15),
				}, nil)
			},
			expectStore: func(m *mocks.ActivityReportStore) {
				m.EXPECT().GetMissingActivityReportsByOfferID(
					arg.Anything,
					offerID,
					jdate.YearMonth{Year: 2023, Month: 3},
					jdate.YearMonth{Year: 2023, Month: 3},
				).Return([]jdate.YearMonth{{Year: 2023, Month: 3}}, nil)
			},
			want: []jdate.YearMonth{{Year: 2023, Month: 3}},
		},
		{
			name:    "contract started 2 months ago - on reporting day - should ignore current month",
			offerID: offerID,
			clock:   jclock.NewFakeClock(time.Date(2023, time.May, 6, 10, 0, 0, 0, time.UTC)),
			expectUsers: func(m *pkgusersmocks.Client) {
				m.EXPECT().GetOfferJobContract(arg.Anything, offerID).Return(pkgusers.Contract{
					StartDate: jdate.New(2023, 3, 15),
				}, nil)
			},
			expectStore: func(m *mocks.ActivityReportStore) {
				m.EXPECT().GetMissingActivityReportsByOfferID(
					arg.Anything,
					offerID,
					jdate.YearMonth{Year: 2023, Month: 3},
					jdate.YearMonth{Year: 2023, Month: 4},
				).Return([]jdate.YearMonth{{Year: 2023, Month: 3}, {Year: 2023, Month: 4}}, nil)
			},
			want: []jdate.YearMonth{{Year: 2023, Month: 3}, {Year: 2023, Month: 4}},
		},
		{
			name:    "contract started 3 months ago - after reporting day - should include current month",
			offerID: offerID,
			clock:   jclock.NewFakeClock(time.Date(2023, time.June, 15, 10, 0, 0, 0, time.UTC)),
			expectUsers: func(m *pkgusersmocks.Client) {
				m.EXPECT().GetOfferJobContract(arg.Anything, offerID).Return(pkgusers.Contract{
					StartDate: jdate.New(2023, 3, 15), // start in March - should include March, April and May
				}, nil)
			},
			expectStore: func(m *mocks.ActivityReportStore) {
				m.EXPECT().GetMissingActivityReportsByOfferID(
					arg.Anything,
					offerID,
					jdate.YearMonth{Year: 2023, Month: 3},
					jdate.YearMonth{Year: 2023, Month: 5},
				).Return([]jdate.YearMonth{{Year: 2023, Month: 3}, {Year: 2023, Month: 4}, {Year: 2023, Month: 5}}, nil)
			},
			want: []jdate.YearMonth{{Year: 2023, Month: 3}, {Year: 2023, Month: 4}, {Year: 2023, Month: 5}},
		},
		{
			name:    "contract started last month - before reporting day - should return empty",
			offerID: offerID,
			clock:   jclock.NewFakeClock(time.Date(2023, time.May, 4, 10, 0, 0, 0, time.UTC)),
			expectUsers: func(m *pkgusersmocks.Client) {
				m.EXPECT().GetOfferJobContract(arg.Anything, offerID).Return(pkgusers.Contract{
					StartDate: jdate.New(2023, 4, 15),
				}, nil)
			},
			want: []jdate.YearMonth{},
		},
		{
			name:    "contract started last month - on reporting day - should include last month",
			offerID: offerID,
			clock:   jclock.NewFakeClock(time.Date(2023, time.May, 6, 10, 0, 0, 0, time.UTC)),
			expectUsers: func(m *pkgusersmocks.Client) {
				m.EXPECT().GetOfferJobContract(arg.Anything, offerID).Return(pkgusers.Contract{
					StartDate: jdate.New(2023, 4, 15),
				}, nil)
			},
			expectStore: func(m *mocks.ActivityReportStore) {
				m.EXPECT().GetMissingActivityReportsByOfferID(
					arg.Anything,
					offerID,
					jdate.YearMonth{Year: 2023, Month: 4},
					jdate.YearMonth{Year: 2023, Month: 4},
				).Return([]jdate.YearMonth{{Year: 2023, Month: 4}}, nil)
			},
			want: []jdate.YearMonth{{Year: 2023, Month: 4}},
		},
		{
			name:    "churner",
			offerID: offerID,
			clock:   jclock.NewFakeClock(time.Date(2023, time.July, 6, 10, 0, 0, 0, time.UTC)),
			expectUsers: func(m *pkgusersmocks.Client) {
				m.EXPECT().GetOfferJobContract(arg.Anything, offerID).Return(pkgusers.Contract{
					StartDate: jdate.New(2023, 4, 15),
					EndDate:   new(jdate.New(2023, 5, 31)),
				}, nil)
			},
			expectStore: func(m *mocks.ActivityReportStore) {
				m.EXPECT().GetMissingActivityReportsByOfferID(
					arg.Anything,
					offerID,
					jdate.YearMonth{Year: 2023, Month: 4},
					jdate.YearMonth{Year: 2023, Month: 5},
				).Return([]jdate.YearMonth{{Year: 2023, Month: 5}}, nil)
			},
			want: []jdate.YearMonth{{Year: 2023, Month: 5}},
		},
		{
			name:    "GetOfferJobContract failed",
			offerID: offerID,
			clock:   jclock.NewFakeClock(time.Date(2023, time.June, 15, 10, 0, 0, 0, time.UTC)),
			expectUsers: func(m *pkgusersmocks.Client) {
				m.EXPECT().GetOfferJobContract(arg.Anything, offerID).Return(pkgusers.Contract{}, assert.AnError)
			},
			wantErr: assert.AnError,
		},
		{
			name:    "GetMissingActivityReportsByOfferID failed",
			offerID: offerID,
			clock:   jclock.NewFakeClock(time.Date(2023, time.June, 15, 10, 0, 0, 0, time.UTC)),
			expectUsers: func(m *pkgusersmocks.Client) {
				m.EXPECT().GetOfferJobContract(arg.Anything, offerID).Return(pkgusers.Contract{
					StartDate: jdate.New(2023, 3, 15),
				}, nil)
			},
			expectStore: func(m *mocks.ActivityReportStore) {
				m.EXPECT().GetMissingActivityReportsByOfferID(
					arg.Anything,
					offerID,
					jdate.YearMonth{Year: 2023, Month: 3},
					jdate.YearMonth{Year: 2023, Month: 5},
				).Return(nil, assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := mocks.NewActivityReportStore(t)
			if tt.expectStore != nil {
				tt.expectStore(store)
			}
			usersClient := pkgusersmocks.NewClient(t)
			if tt.expectUsers != nil {
				tt.expectUsers(usersClient)
			}

			s := &ActivityReportService{
				reportStore: store,
				usersClient: usersClient,
				clock:       tt.clock,
			}

			got, err := s.GetMissingActivityReports(context.Background(), tt.offerID)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
