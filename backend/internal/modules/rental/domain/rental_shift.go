package domain

import (
	"time"

	"github.com/google/uuid"
)

type RentalShift struct {
	ID               uuid.UUID
	BookingID        uuid.UUID
	OperatorID       *uuid.UUID
	ShiftDate        time.Time
	StartTime        time.Time
	EndTime          *time.Time
	StartMeterHours  float64
	EndMeterHours    *float64
	AllocatedHours   float64
	AccumulatedHours float64
	ExcessHours      float64
	IsPaused         bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
