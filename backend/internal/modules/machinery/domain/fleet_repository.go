package domain

import (
	"context"

	"github.com/google/uuid"
)

type FleetRepository interface {
	CreateFleet(ctx context.Context, fleet *FleetAsset) (*FleetAsset, error)
	GetFleetByID(ctx context.Context, id uuid.UUID) (*FleetAsset, error)
	ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]FleetAsset, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status FleetStatus) (*FleetAsset, error)
}
