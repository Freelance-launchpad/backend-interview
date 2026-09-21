// Package jcal is a tiny package used to get useful infomartion about days and calendar events.
package jcal

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	// ErrUnknownCountry is returned when an unknown country is provided to the package.
	ErrUnknownCountry = errors.New("unknown country provided")

	// ErrMonthNotFound is returned when a month is not found.
	ErrMonthNotFound = errors.New("month not found")

	// ErrUnknownLocale is returned when an invalid locale is provided.
	ErrUnknownLocale = errors.New("unknown locale provided")

	// defaultCalendar is a French business calendar
	defaultCal = newFrCalendar()
)

// Default returns the default calendar (France Calendar at the time).
func Default() Calendar {
	return defaultCal
}

func getCountryCalendar(country string) (Calendar, error) {
	country = strings.ToUpper(country)
	switch country {
	case "FR":
		return defaultCal, nil
	}
	return nil, fmt.Errorf("%w: %s", ErrUnknownCountry, country)
}

// GetMaxWorkingDays return the maximum working days for a given month in given country.
// Country should be provided using country code from ISO 3166 - alpha-2 subdivision (https://www.iso.org/obp/ui/#search),
// Example for retrieving working days for April in 2022 in France : GetMaxWorkingDays("FR", 2022, time.April)
// If an unknown country code is provided ErrUnknownCountry is returned.
func GetMaxWorkingDays(country string, y int, m time.Month) (int, error) {

	calendar, err := getCountryCalendar(country)
	if err != nil {
		return 0, err
	}

	return calendar.MaxWorkingDays(y, m), nil
}

var months = map[string][]string{
	"FRA": {"Janvier", "Février", "Mars", "Avril", "Mai", "Juin", "Juillet", "Août", "Septembre", "Octobre", "Novembre", "Décembre"},
	"ENG": {"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"},
}

// MonthToString return a string representation of the given month in the given language.
func MonthToString(m time.Month, locale string) string {
	localizedMonths, ok := months[locale]
	if !ok {
		return ""
	}
	return localizedMonths[m-1]
}

// StringToMonth return a month based on its string representation in the given language.
func StringToMonth(m string, locale string) (time.Month, error) {
	localizedMonths, ok := months[locale]
	if !ok {
		return time.January, ErrMonthNotFound
	}

	m = strings.ToLower(m)

	var month time.Month
	for idx, lm := range localizedMonths {
		if strings.ToLower(lm) == m {
			month = time.Month(idx + 1)
			break
		}
	}

	if month == 0 {
		return month, ErrMonthNotFound
	}

	return month, nil
}

// ToShortDate returns the input time.Time in the given time location as a formatted string using months localization
// Output for ToShortDate(time.Date(2022, 06, 12, 15, 21, 53, 18, time.UTC), time.UTC, "FRA") -> "12 Juin 2022"
func ToShortDate(date time.Time, location time.Location, locale string) (string, error) {
	localizedMonths, ok := months[locale]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrUnknownLocale, locale)
	}

	zonedDate := date.In(&location)

	return fmt.Sprintf("%d %s %d", zonedDate.Day(), localizedMonths[zonedDate.Month()-1], zonedDate.Year()), nil
}

// GetDaysUntilNextWeekday returns the number of days until the next given time.Weekday from the given time.Time.
// Returns 0 if the next given day is the same as today.
func GetDaysUntilNextWeekday(from time.Time, weekday time.Weekday) int {
	return (int(weekday) - int(from.Weekday()) + 7) % 7
}

// GetDaysSinceLastWeekday returns the number of days since the last given time.Weekday from the given time.Time.
// Returns 0 if the last given day is the same as today.
func GetDaysSinceLastWeekday(from time.Time, weekday time.Weekday) int {
	return (int(from.Weekday()) - int(weekday) + 7) % 7
}

// GetNextYearAndMonth returns the year and month from the given year and month when added the given number of months.
func GetNextYearAndMonth(year int, month int, monthsToAdd int) (int, int) {
	from := time.Date(year, time.Month(month), 1, 0, 0, 1, 0, time.UTC)
	to := from.AddDate(0, monthsToAdd, 0)

	return to.Year(), int(to.Month())
}

func GetNextNWorkingDay(date time.Time, daysOffset int) time.Time {
	if daysOffset <= 0 {
		return date
	}

	for daysOffset > 0 || !Default().IsWorkingDay(date) {
		date = date.AddDate(0, 0, 1)
		if Default().IsWorkingDay(date) {
			daysOffset--
		}
	}

	return date
}

func IsWorkingDay(date time.Time) bool {
	return Default().IsWorkingDay(date)
}

func IsWorkTime(date time.Time) bool {
	return Default().IsWorkTime(date)
}

// TimeFrame allows storing a period of time between two time.Time.
type TimeFrame struct {
	StartDate time.Time
	EndDate   time.Time
}

// Diff returns the difference between the start and the end date in days.
// Warning, the function do a subtraction 15/01 - 01/01 = 14.
func (t TimeFrame) Diff() float64 {
	return t.EndDate.Sub((t.StartDate)).Hours() / 24
}

func GetMonthStartEndDate(t time.Time) TimeFrame {
	timeFrame := TimeFrame{
		StartDate: time.Date(
			t.Year(), t.Month(), 1,
			0, 0, 0, 0, t.Location()),
		EndDate: time.Date(
			t.Year(), t.Month()+1, 1,
			0, 0, 0, -1, t.Location()),
	}
	return timeFrame
}
