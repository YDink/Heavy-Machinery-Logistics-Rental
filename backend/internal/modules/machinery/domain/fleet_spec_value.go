package domain

import "github.com/google/uuid"

// FleetSpecValue stores dynamic EAV specifications for an asset.
type FleetSpecValue struct {
	ID           uuid.UUID
	FleetID      uuid.UUID
	SpecID       uuid.UUID
	FieldKey     string
	ValueString  *string
	ValueNumeric *float64
	ValueBoolean *bool
}
