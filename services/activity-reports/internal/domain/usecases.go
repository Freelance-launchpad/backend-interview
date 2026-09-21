package domain

import (
	"context"

	"github.com/Freelance-launchpad/backend-interview/common/jdate"
	"github.com/google/uuid"
)

type ActivityReportService interface {
	CreateActivityReportsV3(ctx context.Context, offerID uuid.UUID, reports []Report) error

	GetMissingActivityReports(ctx context.Context, offerID uuid.UUID) ([]jdate.YearMonth, error)
}
