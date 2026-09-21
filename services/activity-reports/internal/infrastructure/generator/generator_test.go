package generator

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	gotenbergmocks "github.com/Freelance-launchpad/backend-interview/common/clients/gotenberg/mocks"
	"github.com/Freelance-launchpad/backend-interview/common/jclock"
	"github.com/Freelance-launchpad/backend-interview/common/jentity"
	"github.com/Freelance-launchpad/backend-interview/common/js3/v2"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/domain"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/pkg"
	pkgusers "github.com/Freelance-launchpad/backend-interview/services/users/pkg"
	pkgusersmocks "github.com/Freelance-launchpad/backend-interview/services/users/pkg/client/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	arg "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPDFGenerator_GenerateActivityReport(t *testing.T) {
	t.Parallel()
	clock := jclock.NewFakeClock(time.Now())

	template, err := os.ReadFile("../../../activity-reports-cronjobs/resources/template_activity_report.html")
	require.NoError(t, err)

	userID := "bibou"
	reportID := uuid.New()
	report := domain.Report{
		ID:                     reportID,
		ActivityType:           jentity.EntityBlue,
		OfferID:                uuid.New(),
		Month:                  3,
		Year:                   2023,
		WeekType:               new(pkg.WeekTypeTueSat),
		FirstDayAsExtraRest:    new(true),
		ProspectionDuration:    0.5,
		VacationDays:           5,
		DaysAway:               0,
		FormationDuration:      1.5,
		WorkDuration:           200,
		Unit:                   domain.UnitDays,
		OtherActivity:          new(2.5),
		OtherActivityFrequency: new("weekly"),
	}
	user := pkgusers.User{
		ID:                   userID,
		FirstName:            "Raymond",
		LastName:             "Le Chien",
		SocialSecurityNumber: new("123456789"),
	}
	pdf := strings.NewReader("foobarbaz")

	tests := []struct {
		name                           string
		report                         domain.Report
		htmlToPDFConverterExpectations func(*gotenbergmocks.Client)
		userExpectations               func(*pkgusersmocks.Client)
		wantErr                        error
		want                           js3.FileOutput
	}{
		{
			name:   "nominal",
			report: report,
			userExpectations: func(m *pkgusersmocks.Client) {
				m.EXPECT().GetUserByOfferID(arg.Anything, report.OfferID).Return(user, nil)
			},
			htmlToPDFConverterExpectations: func(m *gotenbergmocks.Client) {
				m.EXPECT().HTMLToPDF(arg.Anything, arg.Anything, arg.Anything).Return(int64(42), pdf, nil)
			},
			want: js3.FileOutput{
				ContentLength: int64(42),
				ContentType:   "application/pdf",
				FileName:      "rapport_activité_Mars_2023.pdf",
				Reader:        io.NopCloser(strings.NewReader("foobarbaz")),
			},
		},
		{
			name:   "GetUser failed",
			report: report,
			userExpectations: func(m *pkgusersmocks.Client) {
				m.EXPECT().GetUserByOfferID(arg.Anything, report.OfferID).Return(user, assert.AnError)
			},
			wantErr: assert.AnError,
		},
		{
			name:   "Converter failed",
			report: report,
			userExpectations: func(m *pkgusersmocks.Client) {
				m.EXPECT().GetUserByOfferID(arg.Anything, report.OfferID).Return(user, nil)
			},
			htmlToPDFConverterExpectations: func(m *gotenbergmocks.Client) {
				m.EXPECT().HTMLToPDF(arg.Anything, arg.Anything, arg.Anything).Return(int64(42), pdf, assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pdfGenerator := gotenbergmocks.NewClient(t)
			if tt.htmlToPDFConverterExpectations != nil {
				tt.htmlToPDFConverterExpectations(pdfGenerator)
			}
			usersClient := pkgusersmocks.NewClient(t)
			if tt.userExpectations != nil {
				tt.userExpectations(usersClient)
			}
			service, err := NewPDFGenerator(template, usersClient, pdfGenerator)
			require.NoError(t, err)
			service.clock = clock

			got, err := service.GenerateActivityReport(context.Background(), tt.report)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
