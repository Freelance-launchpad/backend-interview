package jcal

import (
	"time"

	"github.com/rickar/cal/v2"
	"github.com/rickar/cal/v2/fr"
)

// franceCalendar implemts the Calendar interface for France.
type franceCalendar struct {
	calendar *cal.BusinessCalendar
}

func newFrCalendar() *franceCalendar {
	frCal := franceCalendar{
		calendar: cal.NewBusinessCalendar(),
	}
	var holidays []*cal.Holiday
	for _, h := range fr.Holidays {
		if h == fr.LundiDePentecôte {
			// discard "Lundi de pentecote" because this is french solidarity day
			continue
		}
		holidays = append(holidays, h)
	}

	frCal.calendar.AddHoliday(holidays...)
	frCal.calendar.SetWorkHours(time.Hour*8, time.Hour*18)
	return &frCal
}

// Country returns the country code from ISO 3166 - alpha-2 subdivision (https://www.iso.org/obp/ui/#search) for the current calendar.
func (f *franceCalendar) Country() string {
	return "FR"
}

// MaxWorkingDays return the maximum working days for a given month in France.
// Example for retrieving working days for April in 2022 : GetMaxWorkingDays(2022, time.April)
func (f *franceCalendar) MaxWorkingDays(y int, m time.Month) int {
	return f.calendar.WorkdaysInMonth(y, m)
}

// WorkingDaysLeft return the number of working days left in the month in France.
// Example for retrieving working days left at the moment : GetWorkingDaysLeft( time.Now())
// Example with a specific date (2017-05-14): GetWorkingDaysLeft(time.Date(2017, time.May, 14, 12, 0, 0, 0, time.UTC))
func (f *franceCalendar) WorkingDaysLeft(date time.Time) int {
	return f.calendar.WorkdaysRemain(date)
}

// WorkingDaysInRange returns the number of working days in the given range of date (inclusive).
func (f *franceCalendar) WorkingDaysInRange(start, end time.Time) int {
	return f.calendar.WorkdaysInRange(start, end)
}

// IsHoliday tells if the given date is a holiday or not.
func (f *franceCalendar) IsHoliday(date time.Time) bool {
	isHoliday, _, _ := f.calendar.IsHoliday(date)
	return isHoliday
}

// IsWorkingDay tells if the given date is a working day or not.
func (f *franceCalendar) IsWorkingDay(date time.Time) bool {
	return f.calendar.IsWorkday(date)
}

// IsWorkTime tells if the given time is a working time or not.
func (f *franceCalendar) IsWorkTime(t time.Time) bool {
	return f.calendar.IsWorkTime(t)
}
