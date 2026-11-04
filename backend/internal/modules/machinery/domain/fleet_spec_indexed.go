package domain

import (
	"time"

	"github.com/google/uuid"
)

// FleetSpecIndexed represents denormalized logistics-critical dimensions.
type FleetSpecIndexed struct {
	FleetID            uuid.UUID
	OperatingWeightKg  *float64
	TransportLengthMm  *float64
	TransportWidthMm   *float64
	TransportHeightMm  *float64
	PayloadCapacityKg  *float64
	GrossEnginePowerHp *float64
	UpdatedAt          time.Time
}
