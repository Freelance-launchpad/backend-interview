package client

import (
	"context"
	"net/url"

	"github.com/Freelance-launchpad/backend-interview/common/jdate"
	"github.com/Freelance-launchpad/backend-interview/common/jhttp/v2"
	"github.com/Freelance-launchpad/backend-interview/services/payroll/pkg"
	"github.com/google/uuid"
)

type Client interface {
	GetOfferSalaries(ctx context.Context, offerID uuid.UUID, filters pkg.SalaryFiltersQuery) ([]pkg.SalaryResponse, int, error)
	DeleteSalary(ctx context.Context, salaryID uuid.UUID) error
	InitSdtc(ctx context.Context, offerID uuid.UUID, churnDate jdate.Date, salaryID *uuid.UUID) error
	CancelSdtc(ctx context.Context, offerID uuid.UUID) error

	SimulateSalary(ctx context.Context, offerID uuid.UUID, input pkg.SimulationData) (pkg.SimulationResult, error)
	RecomputeSalaryRules(ctx context.Context, salaryID uuid.UUID) error

	GetDebt(ctx context.Context, offerID uuid.UUID) (int64, error)
	GetSeverance(ctx context.Context, offerID uuid.UUID) (int64, error)
}

type client struct {
	client jhttp.Doer
	url    url.URL
}

func New(httpClient jhttp.Doer, uri url.URL) Client {
	return (Client)(nil)
}
