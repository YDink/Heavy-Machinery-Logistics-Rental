package domain

import (
	"time"

	"github.com/google/uuid"
)

// FleetPhoto represents a walkaround photo of the asset.
type FleetPhoto struct {
	ID           uuid.UUID
	FleetID      uuid.UUID
	Angle        PhotoAngle
	StorageURI   string
	FileHash     string
	DisplayOrder int
	CreatedAt    time.Time
}
