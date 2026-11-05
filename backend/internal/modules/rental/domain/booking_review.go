package domain

import (
	"time"

	"github.com/google/uuid"
)

// BookingReview represents a double-blind mutual reputation evaluation archived into the trust ledger.
type BookingReview struct {
	ID               uuid.UUID
	BookingID        uuid.UUID
	AuthorID         uuid.UUID
	TargetID         uuid.UUID
	Role             string
	Rating           float64
	Feedback         string
	DimensionRatings map[string]float64
	IsPublished      bool
	PublishedAt      *time.Time
	CreatedAt        time.Time
}
