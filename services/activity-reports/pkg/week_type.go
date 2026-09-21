package pkg

import (
	"time"

	"github.com/Freelance-launchpad/backend-interview/common/jdate"
)

type WeekType string

const (
	WeekTypeMonFri WeekType = "mon-fri"
	WeekTypeTueSat WeekType = "tue-sat"
)

// UniqueRestDay returns the rest day that's unique to each week type.
// It can return Saturday or Monday.
// Sunday is not returned because it's a rest day for both week types.
func (wt WeekType) UniqueRestDay() time.Weekday {
	if wt == WeekTypeMonFri {
		return time.Saturday
	}
	return time.Monday
}

// Alternate returns the alternative week type configuration.
func (wt WeekType) Alternate() WeekType {
	if wt == WeekTypeMonFri {
		return WeekTypeTueSat
	}
	return WeekTypeMonFri
}

func (wt WeekType) IsRestDay(day jdate.Date) bool {
	if day.Weekday() == time.Sunday {
		return true
	}
	return day.Weekday() == wt.UniqueRestDay()
}
