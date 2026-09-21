package pkg

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/Freelance-launchpad/backend-interview/common/jentity"
	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"

	"github.com/google/uuid"
)

type ExpenseStatus string

const (
	ExpenseStatusInProgress               ExpenseStatus = "in_progress"
	ExpenseStatusNeedUpdate               ExpenseStatus = "need_update"
	ExpenseStatusValidated                ExpenseStatus = "validated"
	ExpenseStatusRefused                  ExpenseStatus = "refused"
	ExpenseStatusDeleted                  ExpenseStatus = "deleted"
	ExpenseStatusCompleted                ExpenseStatus = "completed"
	ExpenseStatusNonRefundable            ExpenseStatus = "non_refundable"
	ExpenseStatusToReview                 ExpenseStatus = "to_review"
	ExpenseStatusAutoValidationInProgress ExpenseStatus = "auto_validation_in_progress"
)

var AllExpenseStatuses = []ExpenseStatus{
	ExpenseStatusInProgress,
	ExpenseStatusNeedUpdate,
	ExpenseStatusValidated,
	ExpenseStatusRefused,
	ExpenseStatusDeleted,
	ExpenseStatusCompleted,
	ExpenseStatusNonRefundable,
	ExpenseStatusToReview,
	ExpenseStatusAutoValidationInProgress,
}

func (s ExpenseStatus) IsValid() bool {
	return slices.Contains(AllExpenseStatuses, s)
}

func (s ExpenseStatus) IsEditable() bool {
	return s == ExpenseStatusInProgress || s == ExpenseStatusNeedUpdate
}

type RequestSaveExpense struct {
	ID                    *uuid.UUID      `json:"id" swaggerignore:"true"`
	OfferID               uuid.UUID       `json:"offer_id" swaggerignore:"true"`
	Status                ExpenseStatus   `json:"status" swaggerignore:"true"`
	Name                  string          `json:"name"`
	Category              ExpenseCategory `json:"category"`
	Amount                int64           `json:"amount"`
	AmountWoVat           *int64          `json:"amount_wo_vat"`
	Client                *uuid.UUID      `json:"client"`
	Mission               *string         `json:"mission"`
	Prospecting           bool            `json:"prospecting"`
	AdditionalInformation string          `json:"additional_information"` // TODO: static typing
	Date                  time.Time       `json:"date"`
	DuplicatedExpenseID   *uuid.UUID      `json:"duplicated_expense_id"`
	Attachments           []Attachment    `json:"attachments"`
	// Region and Currency are derived server-side from the offer entity, never sent by the client.
	Region   jentity.Region `json:"-" swaggerignore:"true"`
	Currency string         `json:"-" swaggerignore:"true"`
}

func (e *RequestSaveExpense) ComputeTransportAmount() error {
	// In the case of annual transport subscription, the amount is divided by 12 and the user needs to re-declare it every month.
	// The initial declared amount is stored in the additional info.
	// Note: when updating a transport subscription, the client needs to set amount=additional_information.declared_amount if it's unchanged,
	// otherwise we have no way of knowing if the amount was updated.
	if e.Category == ExpenseCategoryTransportSubscription {
		var additionalInfo AdditionalInformation
		if err := json.Unmarshal([]byte(e.AdditionalInformation), &additionalInfo); err != nil {
			return err
		}

		if additionalInfo.TransportSubscription == nil {
			return fmt.Errorf("transport subscription additional information is nil")
		}

		additionalInfo.TransportSubscription.DeclaredAmount = e.Amount
		if additionalInfo.TransportSubscription.Period == TransportSubscriptionPeriodYearly {
			e.Amount = int64(math.Round(float64(e.Amount) / 12))
		}

		additionalInfoBytes, err := json.Marshal(additionalInfo)
		if err != nil {
			return err
		}
		e.AdditionalInformation = string(additionalInfoBytes)
	}

	return nil
}

func (e RequestSaveExpense) GetAttachmentIDs() []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(e.Attachments))
	for _, attachment := range e.Attachments {
		if attachment.ID != uuid.Nil {
			ids = append(ids, attachment.ID)
		}
	}
	return ids
}

type RequestSaveExpenseAdmin struct {
	RequestSaveExpense
	RefusalReason           *int    `json:"refusal_reason"`
	AdditionalRefusalReason *string `json:"additional_refusal_reason"`
	Refunded                int64   `json:"-"`
}

type Expense struct {
	ID                      uuid.UUID             `json:"id"`
	OfferID                 uuid.UUID             `json:"offer_id"`
	Name                    string                `json:"name"`
	Category                ExpenseCategory       `json:"category"`
	Icon                    string                `json:"icon"`
	Status                  ExpenseStatus         `json:"status"`
	Amount                  int64                 `json:"amount"`
	AmountWoVat             *int64                `json:"amount_wo_vat"`
	Currency                string                `json:"currency"`
	Client                  *uuid.UUID            `json:"client"`
	Mission                 *string               `json:"mission"`
	AdditionalInformation   string                `json:"additional_information"`
	Attachments             []Attachment          `json:"attachments"`
	RefusalReason           *int                  `json:"refusal_reason"`
	AdditionalRefusalReason *string               `json:"additional_refusal_reason"`
	Date                    time.Time             `json:"date"`
	CreatedAt               time.Time             `json:"created_at"`
	StatusUpdatedAt         *time.Time            `json:"status_updated_at"`
	Refunded                int64                 `json:"refunded"`
	Prospecting             bool                  `json:"-"`
	Locked                  bool                  `json:"locked"`
	DuplicatedExpenseID     *uuid.UUID            `json:"duplicated_expense_id"`
	OCRValidationResults    []OCRValidationResult `json:"ocr_validation_result"`
	AutoValidationEligible  bool                  `json:"auto_validation_eligible"`
	ReimbursementRate       int                   `json:"reimbursement_rate"`
}

func (e Expense) ToRequestSaveExpense() RequestSaveExpense {
	return RequestSaveExpense{
		ID:                    &e.ID,
		OfferID:               e.OfferID,
		Status:                e.Status,
		Name:                  e.Name,
		Category:              e.Category,
		Amount:                e.Amount,
		AmountWoVat:           e.AmountWoVat,
		Client:                e.Client,
		Mission:               e.Mission,
		Prospecting:           e.Prospecting,
		AdditionalInformation: e.AdditionalInformation,
		Date:                  e.Date,
		DuplicatedExpenseID:   e.DuplicatedExpenseID,
		Attachments:           e.Attachments,
	}
}

func (e Expense) ToRequestSaveExpenseAdmin() RequestSaveExpenseAdmin {
	return RequestSaveExpenseAdmin{
		RequestSaveExpense:      e.ToRequestSaveExpense(),
		RefusalReason:           e.RefusalReason,
		AdditionalRefusalReason: e.AdditionalRefusalReason,
		Refunded:                e.Refunded,
	}
}

func (e Expense) GetAttachmentIDs() []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(e.Attachments))
	for _, attachment := range e.Attachments {
		if attachment.ID != uuid.Nil {
			ids = append(ids, attachment.ID)
		}
	}
	return ids
}

func (e Expense) ToSalaryExpense() (_ SalaryExpense, err error) {
	defer jerror.Wrap(&err, "with id", e.ID.String())
	se := SalaryExpense{
		ID:                e.ID,
		Name:              e.Name,
		Category:          e.Category,
		Icon:              e.Icon,
		BaseAmount:        e.Amount,
		BaseAmountWoVat:   e.AmountWoVat,
		Amount:            e.Amount,
		Refunded:          e.Refunded,
		ReimbursementRate: e.ReimbursementRate,
		Date:              e.Date,
	}

	if e.Category == ExpenseCategoryIK {
		var ikInfo AdditionalInformation
		if err := json.Unmarshal([]byte(e.AdditionalInformation), &ikInfo); err != nil {
			return SalaryExpense{}, fmt.Errorf("ik json.Unmarshal error: %w", err)
		}
		if ikInfo.IK == nil {
			return SalaryExpense{}, fmt.Errorf("ik additional information is nil")
		}
	}

	return se, nil
}

type ExpenseFilter struct {
	OfferIDs []uuid.UUID `form:"offer_ids"`

	Offset         *uint64 `form:"offset"`
	Limit          *uint64 `form:"limit"`
	OrderBy        *string `form:"sort"`
	Desc           bool
	OwnerEmail     *string            `form:"owner_email"`
	ExpensesIDs    *[]uuid.UUID       `form:"expenses_ids"`
	DateStart      *time.Time         `form:"date_start"`
	DateEnd        *time.Time         `form:"date_end"`
	MaxExpenseDate *time.Time         `form:"expense_date"`
	Category       *[]ExpenseCategory `form:"category"`
	Client         *uuid.UUID         `form:"client"`
	Mission        *uuid.UUID         `form:"mission"`
	Status         []ExpenseStatus    `form:"status"`
	NotStatus      []ExpenseStatus    `form:"not_status"`
}

type ListExpensesResponse struct {
	Expenses []Expense `json:"items"`
	Count    uint64    `json:"count"`
}

type SalaryExpense struct {
	ID              uuid.UUID       `json:"id"`
	Name            string          `json:"name"`
	Category        ExpenseCategory `json:"category"`
	Icon            string          `json:"icon"`
	BaseAmount      int64           `json:"base_amount"`
	BaseAmountWoVat *int64          `json:"amount_wo_vat"`
	// Amount is the computed amount to refund after applying the rules
	Amount            int64     `json:"amount"`
	ReimbursementRate int       `json:"reimbursement_rate"`
	Refunded          int64     `json:"-"`
	Date              time.Time `json:"date"`
}

const (
	nticAmortizationMonths          = 36
	nticAmortizationAmountThreshold = 500_00
)

// NTIC expenses over 500€ are amortized over 36 months.
func (e SalaryExpense) isAmortizedNTIC() bool {
	if e.Category != ExpenseCategoryNTICGear {
		return false
	}

	nticAmount := e.BaseAmount
	if e.BaseAmountWoVat != nil {
		nticAmount = *e.BaseAmountWoVat
	}

	return nticAmount >= nticAmortizationAmountThreshold
}

func (e SalaryExpense) GetReimbursementRate() int {
	if e.isAmortizedNTIC() {
		return 100
	}

	return e.ReimbursementRate
}

func (e SalaryExpense) GetAmountToRefund() int64 {
	total := e.TotalAmountToRefund()

	if e.isAmortizedNTIC() {
		toRefund := int64(math.Round(float64(e.BaseAmount) / float64(nticAmortizationMonths)))
		return min(toRefund, total-e.Refunded)
	}

	return total
}

func (e SalaryExpense) TotalAmountToRefund() int64 {
	if e.isAmortizedNTIC() {
		return e.BaseAmount
	}

	return int64(math.Round(float64(e.BaseAmount) * float64(e.ReimbursementRate) / 100))
}

// Completed compare the refunded amount to the maximum refunded amount allowed for the expense
func (e SalaryExpense) Completed() bool {
	if slices.Contains([]ExpenseCategory{ExpenseCategoryHousingRenters, ExpenseCategoryHousingTenants}, e.Category) {
		return true
	}

	return e.Refunded >= e.TotalAmountToRefund()
}

type OCRValidationResult struct {
	Name              string         `json:"name"`
	Result            *bool          `json:"result"`
	RecommendedStatus *ExpenseStatus `json:"recommended_status"`
	IsPrevailingRule  *bool          `json:"is_prevailing_rule"`
}
