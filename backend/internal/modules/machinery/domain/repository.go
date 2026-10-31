package domain

import (
	"context"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	ListAll(ctx context.Context) ([]Category, error)
	ListRoots(ctx context.Context) ([]Category, error)
	ListSubcategories(ctx context.Context, parentID uuid.UUID) ([]Category, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*Category, error)
}

type SpecDefinitionRepository interface {
	GetCategorySpecDefinitions(ctx context.Context, categoryID uuid.UUID) ([]SpecFieldGroup, error)
}

type FleetRepository interface {
	CreateFleet(ctx context.Context, fleet *FleetAsset) (*FleetAsset, error)
	GetFleetByID(ctx context.Context, id uuid.UUID) (*FleetAsset, error)
	ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]FleetAsset, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status FleetStatus) (*FleetAsset, error)
}

type FleetMediaRepository interface {
	SavePhotos(ctx context.Context, fleetID uuid.UUID, photos []FleetPhoto) error
	ListPhotos(ctx context.Context, fleetID uuid.UUID) ([]FleetPhoto, error)
}

type FleetSpecRepository interface {
	SaveSpecValues(ctx context.Context, fleetID uuid.UUID, values []FleetSpecValue) error
	UpsertSpecIndexed(ctx context.Context, indexed *FleetSpecIndexed) error
	GetSpecIndexed(ctx context.Context, fleetID uuid.UUID) (*FleetSpecIndexed, error)
}
