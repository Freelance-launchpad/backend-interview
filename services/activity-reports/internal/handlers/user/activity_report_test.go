package userhandlers

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/Freelance-launchpad/backend-interview/common/jdate"
	"github.com/Freelance-launchpad/backend-interview/common/jgin"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/domain"
	mocks "github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_GetMissingReports(t *testing.T) {
	t.Parallel()
	offerID := uuid.New()

	periods := []jdate.YearMonth{
		{
			Month: 3,
			Year:  2023,
		},
		{
			Month: 4,
			Year:  2023,
		},
		{
			Month: 5,
			Year:  2023,
		},
	}

	response := `[
		{
			"month": 3,
			"year": 2023,
			"days": 23
		},
		{
			"month": 4,
			"year": 2023,
			"days": 19
		},
		{
			"month": 5,
			"year": 2023,
			"days": 20
		}
	]`

	tests := []struct {
		name       string
		ctxoptions []jgin.ContextOption
		expect     func(m *mocks.ActivityReportService)
		want       string
		wantStatus int
	}{
		{
			name: "nominal",
			ctxoptions: []jgin.ContextOption{
				jgin.WithOfferID(offerID),
			},
			expect: func(m *mocks.ActivityReportService) {
				m.EXPECT().GetMissingActivityReports(mock.Anything, offerID).Return(periods, nil)
			},
			want:       response,
			wantStatus: http.StatusOK,
		},
		{
			name: "service failed",
			ctxoptions: []jgin.ContextOption{
				jgin.WithOfferID(offerID),
			},
			expect: func(m *mocks.ActivityReportService) {
				m.EXPECT().GetMissingActivityReports(mock.Anything, offerID).Return(nil, assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			activityReportService := mocks.NewActivityReportService(t)
			if tt.expect != nil {
				tt.expect(activityReportService)
			}

			r, c := jgin.CreateTestContext(http.MethodGet, tt.ctxoptions...)

			handlers := NewActivityReportHandler(nil, activityReportService)
			handlers.GetMissingActivityReports(c)

			assert.Equal(t, tt.wantStatus, r.Code)
			if tt.want != "" {
				assert.JSONEq(t, tt.want, r.Body.String())
			}
		})
	}
}

func Test_CreateActivityReports(t *testing.T) {
	t.Parallel()

	offerID := uuid.New()

	reports := []domain.Report{
		{
			Month:                  3,
			Year:                   2023,
			WorkDuration:           23.5,
			ProspectionDuration:    0.5,
			FormationDuration:      1.5,
			VacationDays:           1,
			DaysAway:               0,
			OtherActivity:          new(17.5),
			OtherActivityFrequency: new("daily"),
		},
	}

	body := []byte(`[{
		"month": 3,
		"year": 2023,
		"work_duration": 23.5,
		"prospection": 0.5,
		"vacation_days": 1,
		"days_away": 0,
		"formation": 1.5,
		"other_activity_hours": 17.5,
		"other_activity_frequency": "daily"
	}]`)

	tests := []struct {
		name                  string
		offerID               uuid.UUID
		body                  []byte
		activityReportService func(m *mocks.ActivityReportService)
		wantStatus            int
	}{
		{
			name:    "nominal",
			body:    body,
			offerID: offerID,
			activityReportService: func(m *mocks.ActivityReportService) {
				m.EXPECT().CreateActivityReportsV3(mock.Anything, offerID, reports).Return(nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:    "service failed",
			body:    body,
			offerID: offerID,
			activityReportService: func(m *mocks.ActivityReportService) {
				m.EXPECT().CreateActivityReportsV3(mock.Anything, offerID, reports).Return(assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "bad JSON",
			body:       []byte(`{]`),
			offerID:    offerID,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			activityReportService := mocks.NewActivityReportService(t)
			if tt.activityReportService != nil {
				tt.activityReportService(activityReportService)
			}

			r, c := jgin.CreateTestContext(http.MethodPost, jgin.WithOfferID(tt.offerID), jgin.WithBody(io.NopCloser(bytes.NewBuffer(tt.body))))

			h := NewActivityReportHandler(nil, activityReportService)
			h.CreateActivityReports(c)
			c.Writer.WriteHeaderNow()

			assert.Equal(t, tt.wantStatus, r.Code)
		})
	}
}
