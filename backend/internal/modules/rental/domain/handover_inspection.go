package domain

import (
	"time"

	"github.com/google/uuid"
)

type HandoverInspection struct {
	ID                uuid.UUID
	BookingID         uuid.UUID
	InspectorID       uuid.UUID
	CheckInMeterHours float64
	FluidLevelsOk     bool
	SafetyGlassIntact bool
	Notes             string
	Photos            []InspectionPhoto
	CreatedAt         time.Time
}
