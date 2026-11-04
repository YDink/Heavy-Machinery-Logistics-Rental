package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

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
