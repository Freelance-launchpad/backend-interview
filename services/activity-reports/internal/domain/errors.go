package domain

import (
	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"
)

var (
	// ErrReportAlreadyExists is returned when somebody try to create a report that already exists.
	ErrReportAlreadyExists = jerror.NewAlreadyExistsError("REPORT_ALREADY_EXISTS")
)
