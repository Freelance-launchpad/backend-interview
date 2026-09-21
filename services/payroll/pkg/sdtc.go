package pkg

import (
	"time"

	"github.com/Freelance-launchpad/backend-interview/common/jdate"
	"github.com/google/uuid"
)

type InitSdtcInput struct {
	OfferID   uuid.UUID  `json:"offer_id"`
	ChurnDate jdate.Date `json:"churn_date"`
	SalaryID  *uuid.UUID `json:"salary_id"`
}

type SdtcStatus = string

// Solde de tout compte
type SDTC struct {
	OfferID        uuid.UUID  `json:"offer_id"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Status         SdtcStatus `json:"status"`
	ChurnDate      jdate.Date `json:"churn_date"`
	SalaryAmount   int64      `json:"salary_amount"`
	ExpenseAmount  int64      `json:"expense_amount"`
	TransferAmount int64      `json:"transfer_amount"`
}
