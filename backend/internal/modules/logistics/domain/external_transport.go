package domain

import (
	"time"

	"github.com/google/uuid"
)

// ExternalTransport represents an equipment owner's private/off-platform transport declaration.
type ExternalTransport struct {
	ID                     uuid.UUID
	ParentBookingID        uuid.UUID
	RegisteredBy           uuid.UUID
	CarrierName            string
	DriverFullName         string
	DriverPhoneNumber      string
	PrimeMoverPlateNumber  string
	EstimatedArrivalAt     time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
}
