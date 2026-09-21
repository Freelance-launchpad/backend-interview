package adminhandlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/domain"
	mocks "github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_ActivityReportHandler_GetActivityReports(t *testing.T) {
	t.Parallel()
	offerID := uuid.New()
	reports := []domain.Report{
		{
			ID:      uuid.New(),
			OfferID: offerID,
		},
	}
	tests := []struct {
		name               string
		paramOfferID       string
		expect             func(m *mocks.ActivityReportStore)
		wantResponse       []domain.ReportResponse
		wantResponseStatus int
	}{
		{
			name:         "nominal",
			paramOfferID: offerID.String(),
			expect: func(m *mocks.ActivityReportStore) {
				m.EXPECT().GetActivityReportsByOfferID(mock.Anything, offerID).Return(reports, nil)
			},
			wantResponse:       []domain.ReportResponse{domain.ReportToResponse(reports[0])},
			wantResponseStatus: http.StatusOK,
		},
		{
			name:               "invalid offer id",
			paramOfferID:       "invalid",
			expect:             func(m *mocks.ActivityReportStore) {},
			wantResponse:       nil,
			wantResponseStatus: http.StatusBadRequest,
		},
		{
			name:         "GetActivityReportsByOfferID failed",
			paramOfferID: offerID.String(),
			expect: func(m *mocks.ActivityReportStore) {
				m.EXPECT().GetActivityReportsByOfferID(mock.Anything, offerID).Return(nil, assert.AnError)
			},
			wantResponse:       nil,
			wantResponseStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mocks.NewActivityReportStore(t)
			if tt.expect != nil {
				tt.expect(s)
			}
			h := NewActivityReportHandler(s)
			router := gin.Default()
			router.GET("/activity-reports/v3/admin/offers/:offer_id/activity-reports", h.GetActivityReports)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/activity-reports/v3/admin/offers/"+tt.paramOfferID+"/activity-reports", nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantResponseStatus, w.Code)
			if tt.wantResponseStatus == http.StatusOK {
				var response []domain.ReportResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResponse, response)
			}
		})
	}
}
