package domain

import (
	"time"

	"github.com/google/uuid"
)

type RentalExtension struct {
	ID                         uuid.UUID
	BookingID                  uuid.UUID
	RequestedByID              uuid.UUID
	OriginalEndDate            time.Time
	NewEndDate                 time.Time
	ExtensionDays              int
	IncrementalGrossMinor      int64
	IncrementalPrepaymentMinor int64
	Currency                   string
	Status                     string
	OwnerDecidedAt             *time.Time
	OwnerDecisionNotes         string
	SlipStorageURI             string
	SlipFileHash               string
	ConfirmedAt                *time.Time
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}
