package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type DiscoveryRepository interface {
	SearchFleetsSpatial(ctx context.Context, filter FleetSearchFilter) ([]FleetSearchResult, int, error)
	GetFleetsForComparison(ctx context.Context, fleetIDs []uuid.UUID) ([]FleetAsset, error)
	GetComparableSpecsForFleets(ctx context.Context, fleetIDs []uuid.UUID) ([]SpecFieldDefinition, []FleetSpecValue, error)
	GetFleetAvailability(ctx context.Context, fleetID uuid.UUID, start, end time.Time) (*FleetAvailabilityInfo, error)
	CreateFleetAvailability(ctx context.Context, fleetID uuid.UUID, start, end time.Time, state string, refID *uuid.UUID) error
	UpsertOwnerAggregate(ctx context.Context, agg *OwnerAggregate) error
	GetOwnerAggregate(ctx context.Context, ownerID uuid.UUID) (*OwnerAggregate, error)
	UpsertFleetInsurance(ctx context.Context, ins *FleetInsurance) error
	GetFleetInsurance(ctx context.Context, fleetID uuid.UUID) (*FleetInsurance, error)
}
