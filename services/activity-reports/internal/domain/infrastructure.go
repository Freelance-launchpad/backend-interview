package domain

import (
	"context"

	"github.com/Freelance-launchpad/backend-interview/common/jdate"
	"github.com/Freelance-launchpad/backend-interview/common/js3/v2"
	"github.com/google/uuid"
)

type ActivityReportStore interface {
	CreateActivityReport(ctx context.Context, report Report) error
	GetActivityReportsByOfferID(ctx context.Context, offerID uuid.UUID) ([]Report, error)
	GetActivityReportsToGenerate(ctx context.Context) ([]Report, error)
	GetMissingActivityReportsByOfferID(ctx context.Context, offerID uuid.UUID, from, to jdate.YearMonth) ([]jdate.YearMonth, error)
	SetReportAsGenerated(ctx context.Context, id uuid.UUID) error
}

// ActivityReportGenerator generate files for activity report.
type ActivityReportGenerator interface {
	// GenerateActivityReport generate a file for the given activity report.
	GenerateActivityReport(ctx context.Context, report Report) (js3.FileOutput, error)
}
