package pkg

import "github.com/Freelance-launchpad/backend-interview/common/jerror/v2"

var (
	ErrPendingSalary = jerror.NewBusinessError("PENDING_SALARY")
)
