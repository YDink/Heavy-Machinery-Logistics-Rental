package domain

import (
	"context"

	"github.com/google/uuid"
)

type FleetMediaRepository interface {
	SavePhotos(ctx context.Context, fleetID uuid.UUID, photos []FleetPhoto) error
	ListPhotos(ctx context.Context, fleetID uuid.UUID) ([]FleetPhoto, error)
}
