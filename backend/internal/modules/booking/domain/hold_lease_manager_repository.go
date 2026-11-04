package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type HoldLeaseManager interface {
	AcquireHoldLease(ctx context.Context, fleetID uuid.UUID, attachmentIDs []uuid.UUID, renterID uuid.UUID, ttl time.Duration) (bool, error)
	ExtendHoldLease(ctx context.Context, fleetID uuid.UUID, attachmentIDs []uuid.UUID, ttl time.Duration) error
	ReleaseHoldLease(ctx context.Context, fleetID uuid.UUID, attachmentIDs []uuid.UUID) error
	IsFleetOnHold(ctx context.Context, fleetID uuid.UUID) (bool, error)
}
