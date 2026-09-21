package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Freelance-launchpad/backend-interview/common/jdate"
	"github.com/Freelance-launchpad/backend-interview/common/jhttp/httpmock"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_client_GetRestDays(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	offerID := uuid.New()
	period := jdate.YearMonth{
		Year:  2025,
		Month: 10,
	}

	response := pkg.RestDaysResponse{
		WeekType:            pkg.WeekTypeMonFri,
		FirstDayAsExtraRest: true,
		RestDays:            []jdate.Date{jdate.New(2025, 10, 5)},
	}
	responseJson, err := json.Marshal(response)
	require.NoError(t, err)

	tests := map[string]struct {
		client  *httpmock.Client
		want    pkg.RestDaysResponse
		wantErr string
	}{
		"nominal": {
			client: httpmock.New(t).WithCall(
				http.MethodGet,
				fmt.Sprintf("/activity-reports/v3/admin/offers/%s/work-schedules/rest-days", offerID.String()),
				httpmock.ExpectQueryParam("year", "2025"),
				httpmock.ExpectQueryParam("month", "10"),
				httpmock.ReturnStatus(http.StatusOK),
				httpmock.ReturnBody(string(responseJson)),
			),
			want: response,
		},
		"error": {
			client: httpmock.New(t).WithCall(
				http.MethodGet,
				fmt.Sprintf("/activity-reports/v3/admin/offers/%s/work-schedules/rest-days", offerID.String()),
				httpmock.ExpectQueryParam("year", "2025"),
				httpmock.ExpectQueryParam("month", "10"),
				httpmock.ReturnStatus(http.StatusInternalServerError),
			),
			wantErr: "[500]: bad status, expected: 200, with body: ",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			client := client{Doer: tt.client}

			got, err := client.GetRestDays(ctx, offerID, period)
			assert.Equal(t, tt.want, got)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
