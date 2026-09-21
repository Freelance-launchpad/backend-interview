package usecases

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/Freelance-launchpad/backend-interview/common/jclock"
	"github.com/Freelance-launchpad/backend-interview/common/jdate"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/internal/domain"
	pkgusers "github.com/Freelance-launchpad/backend-interview/services/users/pkg"
	usersclient "github.com/Freelance-launchpad/backend-interview/services/users/pkg/client"
	"github.com/Thiht/transactor"
	"github.com/google/uuid"
)

type ActivityReportService struct {
	reportStore domain.ActivityReportStore
	transactor  transactor.Transactor

	generator domain.ActivityReportGenerator

	usersClient usersclient.Client

	clock jclock.Clock
}

func NewActivityReportService(ars domain.ActivityReportStore, transactor transactor.Transactor,
	generator domain.ActivityReportGenerator, usersClient usersclient.Client) *ActivityReportService {

	return &ActivityReportService{
		reportStore: ars,
		transactor:  transactor,
		generator:   generator,
		usersClient: usersClient,
		clock:       jclock.RealClock{},
	}
}

func (s *ActivityReportService) CreateActivityReportsV3(ctx context.Context, offerID uuid.UUID, reports []domain.Report) error {
	const errorMsg = "ActivityReportService.CreateActivityReportsV3 has failed"

	offerWithContract, err := s.usersClient.GetOfferWithContract(ctx, offerID)
	if err != nil {
		return fmt.Errorf("%s: %w", errorMsg, err)
	}

	if err := s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		for i := range reports {
			reports[i].OfferID = offerID
			reports[i].ActivityType = offerWithContract.GetEntity()
			reports[i].Unit = domain.EntityToUnit(offerWithContract.GetEntity())

			if err := s.reportStore.CreateActivityReport(ctx, reports[i]); err != nil {
				return fmt.Errorf("%s: %w", errorMsg, err)
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func getFirstPeriod(contract pkgusers.Contract) jdate.YearMonth {
	period := jdate.YearMonth{Month: 3, Year: 2023}

	startDatePeriod := jdate.YearMonth{Month: contract.StartDate.Month(), Year: contract.StartDate.Year()}
	if startDatePeriod.After(period) {
		period = startDatePeriod
	}

	return period
}

func (s *ActivityReportService) GetMissingActivityReports(ctx context.Context, offerID uuid.UUID) ([]jdate.YearMonth, error) {
	const reportFirstDay = 6
	today := jdate.NewFromTime(s.clock.Now().UTC())

	jobContract, err := s.usersClient.GetOfferJobContract(ctx, offerID)
	if err != nil {
		return nil, err
	}
	if jobContract.StartDate.After(today) {
		return []jdate.YearMonth{}, nil
	}

	from := getFirstPeriod(jobContract)

	if today.Before(jdate.New(from.Year, from.Month+1, reportFirstDay)) {
		return []jdate.YearMonth{}, nil
	}

	to := jdate.YearMonth{
		Month: today.Month(),
		Year:  today.Year(),
	}.Previous()

	if today.Day() < reportFirstDay {
		to = to.Previous()
	}

	if endDate := jobContract.EndDate; endDate != nil {
		lastContractMonth := jdate.YearMonth{Year: endDate.Year(), Month: endDate.Month()}
		if to.After(lastContractMonth) {
			to = lastContractMonth
		}
	}

	return s.reportStore.GetMissingActivityReportsByOfferID(ctx, offerID, from, to)
}

func (s *ActivityReportService) generateReport(ctx context.Context, report domain.Report) error {
	file, err := s.generator.GenerateActivityReport(ctx, report)
	if err != nil {
		return err
	}

	if file.ContentLength == 0 {
		return nil
	}

	fileContent, err := io.ReadAll(file.Reader)
	if err != nil {
		return err
	}

	fields := map[string]string{
		"tag":        "activity_report",
		"admin_only": "false",
		"year":       strconv.Itoa(report.Year),
		"month":      strconv.Itoa(report.Month),
	}

	if _, err := s.usersClient.UploadOfferFile(ctx, report.OfferID, fileContent, "activity_report.pdf", fields); err != nil {
		return err
	}

	return nil
}
