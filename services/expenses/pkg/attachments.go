package pkg

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type AttachmentType string

const (
	AttachmentTypeProof             AttachmentType = "proof"
	AttachmentTypeReceipt           AttachmentType = "receipt"
	AttachmentTypeCoproperty        AttachmentType = "coproperty"
	AttachmentTypeHouseAssurance    AttachmentType = "houseAssurance"
	AttachmentTypeReceiptAdditional AttachmentType = "receiptAdditional"
)

type AttachmentForm struct {
	Type            AttachmentType        `form:"type" binding:"required"`
	ExpenseCategory ExpenseCategory       `form:"expense_category" binding:"required"`
	File            *multipart.FileHeader `form:"file" binding:"required"`
}

type AdminAttachmentForm struct {
	AttachmentForm
	OfferID string `form:"offer_id" binding:"required"`
}

type AttachmentIDResponse struct {
	ID uuid.UUID `json:"id"`
}

// Attachment is the output struct returned by the API.
type Attachment struct {
	ID       uuid.UUID      `json:"id"`
	Type     AttachmentType `json:"type"`
	Filename string         `json:"filename"`
	URL      string         `json:"url"`
	OCRData  any            `json:"ocr_data"`
}
