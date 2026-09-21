package pkg

import (
	"github.com/google/uuid"
)

type LockExpensesQuery struct {
	OfferID     uuid.UUID   `json:"offer_id"`
	ExpensesIDs []uuid.UUID `json:"expenses_ids"`
}

// RefundedExpense represents a refunded expense with its refunded amount taken into consideration.
type RefundedExpense struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	Refunded int64     `json:"refunded_amount"`
}

type RefundedExpenses struct {
	Expenses []RefundedExpense `json:"refunded_expenses"`
}

type NonRefundableParams struct {
	OfferID uuid.UUID `json:"offer_id"`
}
