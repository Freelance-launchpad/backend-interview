package pkg

import (
	"github.com/Freelance-launchpad/backend-interview/common/jdate"
)

type RestDaysResponse struct {
	WeekType            WeekType     `json:"week_type"`
	FirstDayAsExtraRest bool         `json:"first_day_as_extra_rest"`
	RestDays            []jdate.Date `json:"rest_days"`
}
