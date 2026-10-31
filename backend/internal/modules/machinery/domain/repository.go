package domain

import (
	"context"
	"time"

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

type AttachmentRepository interface {
	CreateAttachment(ctx context.Context, attachment *Attachment, compatibleCategories []string, photos []AttachmentPhoto) (*Attachment, error)
	GetAttachmentByID(ctx context.Context, id uuid.UUID) (*Attachment, error)
	ListAttachmentsByOwner(ctx context.Context, ownerID uuid.UUID, status string) ([]*Attachment, error)
	UpdateAttachmentStatus(ctx context.Context, id uuid.UUID, status AttachmentStatus) error
	LinkAttachmentToFleet(ctx context.Context, fleetID uuid.UUID, attachmentID uuid.UUID, bundledDailyRateMinor int64) (*FleetBundledAttachment, error)
	ListBundledAttachmentsByFleet(ctx context.Context, fleetID uuid.UUID) ([]*FleetBundledAttachment, error)
	CheckAttachmentAvailability(ctx context.Context, attachmentID uuid.UUID, startDate, endDate time.Time) (bool, error)
	CreateReservation(ctx context.Context, res *AttachmentReservation) error
}
