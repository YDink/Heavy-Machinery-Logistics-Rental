package domain

import (
	"time"

	"github.com/google/uuid"
)

type BookingAttachment struct {
	ID                    uuid.UUID
	BookingID             uuid.UUID
	AttachmentID          uuid.UUID
	AttachmentType        string
	Make                  string
	Model                 string
	SerialNumber          string
	BundledDailyRateMinor int64
	CreatedAt             time.Time
}
