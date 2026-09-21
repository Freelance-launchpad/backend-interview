package pkg

import (
	"time"

	"github.com/Freelance-launchpad/backend-interview/common/jkafka"
	"github.com/google/uuid"
)

const TopicSalaries = "salaries"

// SalariesEvent is sent on the `salaries` topic on Kafka.
type SalariesEvent struct {
	jkafka.Event
	State SalaryState `json:"state"`
}

type SalaryState struct {
	ID                   uuid.UUID    `json:"id"`
	OfferID              uuid.UUID    `json:"offer_id"`
	StateID              uuid.UUID    `json:"state_id"`
	Timestamp            time.Time    `json:"timestamp"`
	Status               SalaryStatus `json:"status"`
	PaidAt               *time.Time   `json:"paid_at"`
	Year                 int          `json:"year"`
	Month                int          `json:"month"`
	HRID                 *string      `json:"hr_id"`
	Gross                int64        `json:"gross"`
	Net                  int64        `json:"net"`
	NetTaxable           int64        `json:"net_taxable"`
	FixedPart            int64        `json:"fixed_part"`
	VariablePart         int64        `json:"variable_part"`
	EmployerContribution int64        `json:"employer_contribution"`
	EmployeeContribution int64        `json:"employee_contribution"`
	WorkedDays           int          `json:"worked_days"`
	WorkedHours          float64      `json:"worked_hours"`
	Balance              int64        `json:"balance"`
	IncomeTaxRate        float64      `json:"income_tax_rate"`
	NetAfterTax          int64        `json:"net_after_tax"`
	ToPay                int64        `json:"to_pay"`
	PayslipFile          *string      `json:"payslip_file"`
	MealVouchers         *int         `json:"meal_vouchers"`
	MealVouchersAmount   int64        `json:"meal_vouchers_amount"`
	ForceOpen            bool         `json:"force_open"`
	SalaryMode           *SalaryMode  `json:"salary_mode"`
	Misc                 int64        `json:"misc"`
	HealthPlan           *HealthPlan  `json:"health_plan"`
	IsSDTC               bool         `json:"is_sdtc"`
}

type HealthPlan struct {
	EmployerPart int64 `json:"employer_part"`
	EmployeePart int64 `json:"employee_part"`
	VariablePart int64 `json:"variable_part"`
	Children     *bool `json:"children"`
	Partner      *bool `json:"partner"`
}
