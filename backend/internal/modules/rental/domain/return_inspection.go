package domain

import (
	"time"

	"github.com/google/uuid"
)

// ReturnInspection represents an off-hire physical inspection executed before machinery demobilization.
type ReturnInspection struct {
	ID                  uuid.UUID
	BookingID           uuid.UUID
	InspectorID         uuid.UUID
	CheckOutMeterHours  float64
	FluidLevelsOK       bool
	SafetyGlassIntact   bool
	AccelerometerSteady bool
	Notes               string
	Photos              []InspectionPhoto
	CreatedAt           time.Time
}
