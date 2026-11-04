package domain

import (
	"time"

	"github.com/google/uuid"
)

type BankSlip struct {
	ID                   uuid.UUID
	BookingID            uuid.UUID
	UploaderID           uuid.UUID
	StorageURI           string
	FileHashSHA256       string
	TransactionReference *string
	AmountMinor          int64
	Currency             string
	Status               BankSlipStatus
	RejectionReasonCode  *string
	RejectionNotes       *string
	CreatedAt            time.Time
	AuditedAt            *time.Time
}
