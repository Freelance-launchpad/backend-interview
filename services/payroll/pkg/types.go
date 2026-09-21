package pkg

import (
	"time"

	"github.com/Freelance-launchpad/backend-interview/common/jdb"
	"github.com/Freelance-launchpad/backend-interview/common/jentity"

	"github.com/google/uuid"
)

// PayrollPeriodStart is used because query params doesn't apply prefix for nested struct
// and we cannot use year query param for end and start payroll period
type PayrollPeriodStart struct {
	Year  int `form:"payroll_period_start_year"`
	Month int `form:"payroll_period_start_month"`
}

// PayrollPeriodEnd is used because query params doesn't apply prefix for nested struct
// and we cannot use year query param for end and start payroll period
type PayrollPeriodEnd struct {
	Year  int `form:"payroll_period_end_year"`
	Month int `form:"payroll_period_end_month"`
}

type SalaryFiltersQuery struct {
	jdb.QueryPagination

	PayrollPeriodStart *PayrollPeriodStart
	PayrollPeriodEnd   *PayrollPeriodEnd

	Status    []SalaryStatus `form:"status,omitempty"`
	DateStart *time.Time     `form:"start_date,omitempty"`
	DateEnd   *time.Time     `form:"end_date,omitempty"`

	ExpenseID *string `form:"expense_id,omitempty"`
}

type healthCareJSON struct {
	EmployeePart int64 `json:"employee_part"`
	EmployerPart int64 `json:"employer_part"`
	VariablePart int64 `json:"variable_part"`
	Children     bool  `json:"children"`
	Partner      bool  `json:"partner"`
}

// SalaryStatus represents the status of a Salary.
type SalaryStatus string

type SalaryMode string

type SalaryResponse struct {
	ID                   uuid.UUID      `json:"salary_id"`
	OfferID              uuid.UUID      `json:"offer_id"`
	Status               SalaryStatus   `json:"status"`
	Year                 int            `json:"year"`
	Month                int            `json:"month"`
	Revenue              int64          `json:"revenue"`
	SalaryCost           int64          `json:"salary_cost"`
	GrossSalary          int64          `json:"gross"`
	Net                  int64          `json:"net"`
	NetTaxable           int64          `json:"net_taxable"`
	NetAfterTax          int64          `json:"net_after_tax"`
	VariablePart         int64          `json:"variable_part"`
	EmployerContribution int64          `json:"employer_contribution"`
	EmployeeContribution int64          `json:"employee_contribution"`
	KilometricAllowance  int64          `json:"kilometric_allowance"`
	IncomeTaxRate        float64        `json:"income_tax_rate"`
	WorkedDays           float64        `json:"worked_days"`
	WorkedHours          float64        `json:"worked_hours"`
	PayslipFile          *string        `json:"payslip_file"`
	HealthPlan           healthCareJSON `json:"health_plan"`
	ToPay                int64          `json:"amount"`
	Type                 string         `json:"type"`
	ExpensesAmount       int64          `json:"expenses_amount"`
	SalaryMode           *SalaryMode    `json:"salary_mode"`
}

type SalaryResponseList struct {
	Salaries       []SalaryResponse `json:"items"`
	Count          int              `json:"count"`
	SumNetAfterTax *int64           `json:"sum_net_after_tax,omitempty"`
}

type MaxSimulationData struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

type SimulationData struct {
	Gross int64 `json:"gross"`
	Month int   `json:"month"`
	Year  int   `json:"year"`
}

type GetDebtOutput struct {
	Debt int64 `json:"debt"`
}

type GetSeveranceOutput struct {
	Severance int64 `json:"severance"`
}

type SimulationResult struct {
	Type                  jentity.Entity `json:"type"`
	Gross                 int64          `json:"gross"`
	EmployerContributions int64          `json:"employer_contributions"`
	EmployeeContributions int64          `json:"employee_contributions"`
	Net                   int64          `json:"net"`
	IncomeTax             int64          `json:"income_tax"`
	IncomeTaxRate         float64        `json:"income_tax_rate"`
	NetAfterTax           int64          `json:"net_after_tax"`
	MealVouchers          *int64         `json:"meal_vouchers"`
	Debt                  int64          `json:"employee_contributions_debt_amount"`
	ToPay                 int64          `json:"to_pay"`
	Cost                  int64          `json:"salary_cost"`
	TotalRevenue          int64          `json:"total_revenue"`
	RestitutionRate       int            `json:"restitution_rate"`
	MaxExpensesAmount     int64          `json:"max_expenses_amount"`
	DebugData             map[string]any `json:"debug_data,omitempty"`
}
