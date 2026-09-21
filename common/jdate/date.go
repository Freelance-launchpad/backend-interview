package jdate

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Freelance-launchpad/backend-interview/common/jcal"
	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"

	"gopkg.in/yaml.v3"
)

const FrenchDateFormat = "02/01/2006"

// Date represents a data.
type Date time.Time

// New returns a new Date.
func New(year int, month time.Month, day int) Date {
	return Date(time.Date(
		year,
		month,
		day,
		0, 0, 0, 0, time.UTC,
	))
}

// NewFromTime returns a new Date from a time.Time.
func NewFromTime(time time.Time) Date {
	return New(time.Year(), time.Month(), time.Day())
}

// LastDayOfTheMonth returns a date which is the last day of the month.
func LastDayOfTheMonth(year, month int) Date {
	return New(year, time.Month(month+1), 0)
}

// Parse parses a string and returns a date.
func Parse(stringValue string) (Date, error) {
	date, err := time.Parse(time.RFC3339Nano, stringValue)
	if err != nil {
		if _, ok := err.(*time.ParseError); ok {
			date, err = time.Parse(time.DateOnly, stringValue)
			if err != nil {
				return Date{}, err
			}
		} else {
			return Date{}, err
		}
	}

	return Date(time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0, 0, 0, 0, time.UTC,
	)), nil
}

// UnmarshalJSON implements JSON Marshaller interface.
func (d *Date) UnmarshalJSON(data []byte) error {
	var stringValue string

	err := json.Unmarshal(data, &stringValue)
	if err != nil {
		return err
	}

	date, err := Parse(stringValue)
	if err != nil {
		return err
	}

	*d = date

	return nil
}

// UnmarshalParam is an implementation of Gin's binding.BindUnmarshaler.
// It's implemented explicitly in jgin/param_binding.go, to make sure gin is not imported in jdate.
func (d *Date) UnmarshalParam(param string) error {
	date, err := Parse(param)
	if err != nil {
		return err
	}
	*d = Date(date)
	return nil
}

// MarshalJSON implements JSON Marshaller interface.
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format(time.RFC3339Nano))
}

// Compare returns an integer comparing two dates.
func (d Date) Compare(other Date) int {
	return time.Time(d).Compare((time.Time)(other))
}

// Before reports whether the date d is before others.
func (d Date) Before(other Date) bool {
	return time.Time(d).Before((time.Time)(other))
}

// BeforeTime reports whether the date d is before others.
func (d Date) BeforeTime(other time.Time) bool {
	otherUTC := time.Date(other.Year(), other.Month(), other.Day(), 0, 0, 0, 0, time.UTC)
	return time.Time(d).Before(otherUTC)
}

// After reports whether the date d is after others.
func (d Date) After(other Date) bool {
	return time.Time(d).After((time.Time)(other))
}

// AfterTime reports whether the date d is after others.
func (d Date) AfterTime(other time.Time) bool {
	otherUTC := time.Date(other.Year(), other.Month(), other.Day(), 0, 0, 0, 0, time.UTC)
	return time.Time(d).After(otherUTC)
}

// AddDate returns the date corresponding to adding the given number of years, months, and days to d.
func (d Date) AddDate(years int, months int, days int) Date {
	return Date(time.Time(d).AddDate(years, months, days))
}

// IsZero reports whether t represents the zero time instant,
// January 1, year 1, 00:00:00 UTC.
func (d Date) IsZero() bool {
	return time.Time(d).IsZero()
}

// Year returns the year of the date.
func (d Date) Year() int {
	return time.Time(d).Year()
}

// Month returns the month of the date.
func (d Date) Month() time.Month {
	return time.Time(d).Month()
}

// Day returns the day of the date.
func (d Date) Day() int {
	return time.Time(d).Day()
}

// Sub returns the duration d - t.
func (d Date) Sub(t time.Time) time.Duration {
	return time.Time(d).Sub(t)
}

// Format returns a textual representation of the date value formatted according to layout.
func (d Date) Format(layout string) string {
	return time.Time(d).Format(layout)
}

// Scan implements the Scanner interface.
func (d *Date) Scan(value any) error {
	if value == nil {
		return nil
	}

	if v, ok := value.(time.Time); ok {
		*d = New(v.Year(), v.Month(), v.Day())

		return nil
	}

	if v, ok := value.(string); ok {
		t, err := Parse(v)
		if err != nil {
			return fmt.Errorf("failed to scan Date value: %v", value)
		}

		*d = t

		return nil
	}

	return fmt.Errorf("failed to scan Date value: %v", value)
}

// Value implements the driver Valuer interface.
func (d Date) Value() (driver.Value, error) {
	return time.Time(d), nil
}

// Equal reports whether the date d is equal to others.
func (d Date) Equal(other Date) bool {
	return time.Time(d).Equal((time.Time)(other))
}

// Time returns the time.Time representation of the date.
func (d Date) Time() time.Time {
	return time.Time(d)
}

// AddWorkingDays add working days to a date.
func (d Date) AddWorkingDays(days int) Date {
	return Date(jcal.GetNextNWorkingDay(time.Time(d), days))
}

// Unix returns the number of seconds elapsed since January 1, 1970 UTC.
func (d Date) Unix() int64 {
	return time.Time(d).Unix()
}

func (d Date) Weekday() time.Weekday {
	return d.Time().Weekday()
}

// Today is a helper function that returns today's date.
func Today() Date { return NewFromTime(time.Now()) }

var _ yaml.Unmarshaler = &Date{}

func (d *Date) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err != nil {
		return err
	}

	date, err := Parse(s)
	if err != nil {
		return err
	}

	*d = date
	return nil
}

var _ yaml.Marshaler = Date{}

func (d Date) MarshalYAML() (any, error) {
	return time.Time(d).Format(time.RFC3339Nano), nil
}

// YearMonth represents a year/month tuple.
type YearMonth struct {
	Year  int        `json:"year"`
	Month time.Month `json:"month"`
}

var ErrInvalidateYearMonth = jerror.NewValidationError("INVALID_YEAR_MONTH")

func (ym YearMonth) String() string {
	return fmt.Sprintf("%04d-%02d", ym.Year, ym.Month)
}

func (ym YearMonth) Validate() error {
	if ym.Year < 0 || ym.Month < 1 || ym.Month > 12 {
		return ErrInvalidateYearMonth
	}
	return nil
}

// FirstDay returns the first day of the month.
func (ym YearMonth) FirstDay() Date {
	return New(ym.Year, ym.Month, 1)
}

// LastDay returns the last day of the month.
func (ym YearMonth) LastDay() Date {
	return New(ym.Year, ym.Month+1, 0)
}

// ToTimeFrame converts a DateTuple into a jcal.TimeFrame.
func (ym YearMonth) ToTimeFrame() jcal.TimeFrame {
	return jcal.TimeFrame{
		StartDate: time.Date(ym.Year, ym.Month, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(ym.Year, ym.Month+1, 0, 0, 0, 0, 0, time.UTC),
	}
}

var _ yaml.Unmarshaler = &YearMonth{}

// UnmarshalYAML Cast string into DateTuple during yaml Unmarshal.
func (ym *YearMonth) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err != nil {
		return err
	}

	elem := strings.Split(s, "/")
	if len(elem) != 2 || len(elem[0]) != 2 || len(elem[1]) != 4 {
		return fmt.Errorf("unable to parse: %s, must be mm/yyyy", s)
	}

	month, err := strconv.Atoi(elem[0])
	if err != nil {
		return fmt.Errorf("unable to parse: %s, mm must be a valid number", s)
	}
	year, err := strconv.Atoi(elem[1])
	if err != nil {
		return fmt.Errorf("unable to parse: %s, yyyy must be a valid number", s)
	}

	if month > 12 || month < 0 || year < 0 {
		return fmt.Errorf("unable to parse: %s, mm/yyyy must be a valid date", s)
	}

	*ym = YearMonth{Year: year, Month: time.Month(month)}

	return nil
}

var _ yaml.Marshaler = YearMonth{}

// cast DateTuple type into yaml string (during yaml marshal)
func (ym YearMonth) MarshalYAML() (any, error) {
	return fmt.Sprintf("%02d/%04d", ym.Month, ym.Year), nil
}

func (ym YearMonth) After(other YearMonth) bool {
	return time.Date(ym.Year, ym.Month, 1, 0, 0, 0, 0, time.UTC).
		After(time.Date(other.Year, other.Month, 1, 0, 0, 0, 0, time.UTC))
}

func (ym YearMonth) AfterOrEqual(other YearMonth) bool {
	return time.Date(ym.Year, ym.Month, 1, 0, 0, 0, 0, time.UTC).
		After(time.Date(other.Year, other.Month, 0, 0, 0, 0, 0, time.UTC))
}

func (ym YearMonth) Before(other YearMonth) bool {
	return time.Date(ym.Year, ym.Month, 1, 0, 0, 0, 0, time.UTC).
		Before(time.Date(other.Year, other.Month, 1, 0, 0, 0, 0, time.UTC))
}

func (ym YearMonth) BeforeOrEqual(other YearMonth) bool {
	return time.Date(ym.Year, ym.Month, 1, 0, 0, 0, 0, time.UTC).
		Before(time.Date(other.Year, other.Month+1, 1, 0, 0, 0, 0, time.UTC))
}

func (ym YearMonth) Equal(other YearMonth) bool {
	return ym.Month == other.Month && ym.Year == other.Year
}

func (ym YearMonth) Previous() YearMonth {
	tuple := YearMonth{
		Year:  ym.Year,
		Month: ym.Month - 1,
	}

	if tuple.Month < 1 {
		tuple.Year -= 1
		tuple.Month = 12
	}

	return tuple
}

func (ym YearMonth) Next() YearMonth {
	tuple := YearMonth{
		Year:  ym.Year,
		Month: ym.Month + 1,
	}

	if tuple.Month > 12 {
		tuple.Year += 1
		tuple.Month = 1
	}

	return tuple
}

func (ym YearMonth) AddMonth(offset int) YearMonth {
	newDate := ym.FirstDay().AddDate(0, offset, 0)
	return YearMonth{Year: newDate.Year(), Month: newDate.Month()}
}

func (ym YearMonth) Sub(other YearMonth) int {
	return (ym.Year-other.Year)*12 + int(ym.Month-other.Month)
}
