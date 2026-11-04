package domain

import (
	"context"

	"github.com/google/uuid"
)

type FleetSpecRepository interface {
	SaveSpecValues(ctx context.Context, fleetID uuid.UUID, values []FleetSpecValue) error
	UpsertSpecIndexed(ctx context.Context, indexed *FleetSpecIndexed) error
	GetSpecIndexed(ctx context.Context, fleetID uuid.UUID) (*FleetSpecIndexed, error)
}
