package domain

import (
	"time"

	"github.com/Freelance-launchpad/backend-interview/common/jentity"
	"github.com/Freelance-launchpad/backend-interview/services/activity-reports/pkg"
	"github.com/google/uuid"
)

type Unit string

const (
	UnitDays  Unit = "days"
	UnitHours Unit = "hours"
)

func EntityToUnit(entity jentity.Entity) Unit {
	if entity == jentity.EntityBlue {
		return UnitDays
	}
	return UnitHours
}

// Report describe the time spent by a user in a month.
// Worked, Prospection and Formation are expressed in days for blue offers and hours for green/realty.
// VacationDays and DaysAway are always reported in days.
type Report struct {
	ID           uuid.UUID
	OfferID      uuid.UUID
	Month        int
	Year         int
	ActivityType jentity.Entity

	WeekType            *pkg.WeekType
	FirstDayAsExtraRest *bool

	WorkDuration        float64
	ProspectionDuration float64
	FormationDuration   float64
	Unit                Unit

	VacationDays           float64
	DaysAway               float64
	OtherActivity          *float64
	OtherActivityFrequency *string

	CreatedAt time.Time
}

type CreateReportInput struct {
	Month                  int      `json:"month"`
	Year                   int      `json:"year"`
	WorkDuration           float64  `json:"work_duration"`
	Prospection            float64  `json:"prospection"`
	Formation              float64  `json:"formation"`
	VacationDays           float64  `json:"vacation_days"`
	DaysAway               float64  `json:"days_away"`
	OtherActivity          *float64 `json:"other_activity_hours"`
	OtherActivityFrequency *string  `json:"other_activity_frequency"`
}

func (r CreateReportInput) ToReport() Report {
	return Report{
		Month:                  r.Month,
		Year:                   r.Year,
		WorkDuration:           r.WorkDuration,
		ProspectionDuration:    r.Prospection,
		FormationDuration:      r.Formation,
		VacationDays:           r.VacationDays,
		DaysAway:               r.DaysAway,
		OtherActivity:          r.OtherActivity,
		OtherActivityFrequency: r.OtherActivityFrequency,
	}
}

type ReportResponse struct {
	ID                  string        `json:"id"`
	Month               int           `json:"month"`
	Year                int           `json:"year"`
	WeekType            *pkg.WeekType `json:"week_type"`
	FirstDayAsExtraRest *bool         `json:"first_day_as_extra_rest"`
	WorkDuration        float64       `json:"work_duration"`
	ProspectionDuration float64       `json:"prospection_duration"`
	FormationDuration   float64       `json:"formation_duration"`
	Unit                Unit          `json:"unit"`
	VacationDays        float64       `json:"vacation_days"`
	DaysAway            float64       `json:"days_away"`
	CreatedAt           string        `json:"created_at"`
}

func ReportToResponse(r Report) ReportResponse {
	return ReportResponse{
		ID:                  r.ID.String(),
		Month:               r.Month,
		Year:                r.Year,
		WeekType:            r.WeekType,
		FirstDayAsExtraRest: r.FirstDayAsExtraRest,
		WorkDuration:        r.WorkDuration,
		ProspectionDuration: r.ProspectionDuration,
		FormationDuration:   r.FormationDuration,
		Unit:                r.Unit,
		VacationDays:        r.VacationDays,
		DaysAway:            r.DaysAway,
		CreatedAt:           r.CreatedAt.Format(time.RFC3339),
	}
}

type MissingReportResponse struct {
	Month int `json:"month"`
	Year  int `json:"year"`
	Days  int `json:"days"`
}
