package jcal

import "time"

// Calendar defines a calendar attached to a country.
type Calendar interface {

	// Country returns the country code from ISO 3166 - alpha-2 subdivision (https://www.iso.org/obp/ui/#search) for the current calendar.
	Country() string

	// MaxWorkingDays returns the maximum working days for a given month.
	// Example for retrieving working days for April in 2022 : GetMaxWorkingDays(2022, time.April)
	MaxWorkingDays(y int, m time.Month) int

	// WorkingDaysLeft returns the number of working days left in the month.
	// Example for retrieving working days left at the moment : GetWorkingDaysLeft( time.Now())
	// Example with a specific date (2017-05-14): GetWorkingDaysLeft(time.Date(2017, time.May, 14, 12, 0, 0, 0, time.UTC))
	WorkingDaysLeft(date time.Time) int

	// WorkingDaysInRange returns the number of working days in the given range of date (inclusive).
	WorkingDaysInRange(start, end time.Time) int

	// IsHoliday tells if the given date is a holiday or not.
	IsHoliday(date time.Time) bool

	// IsWorkingDay tells if the given date is a working day or not.
	IsWorkingDay(date time.Time) bool

	// IsWorkTime tells if the given time is a working time or not.
	IsWorkTime(t time.Time) bool
}
