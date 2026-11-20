package domain

import (
	"time"

	"github.com/google/uuid"
)

// LogisticsDispatch represents an in-platform lowboy sub-booking mobilization leg.
type LogisticsDispatch struct {
	ID                         uuid.UUID
	SubBookingID               uuid.UUID
	ParentBookingID            uuid.UUID
	LogisticsFleetID           uuid.UUID
	CarrierOwnerID             uuid.UUID
	PickupDateTime             time.Time
	Origin                     GeoPoint
	Destination                GeoPoint
	FreightHireFeeMinor        int64
	PlatformCommissionBps      int
	CommissionOwedMinor        int64
	Currency                   string
	AbnormalLoadEscortRequired bool
	EscortFeeMinor             int64
	WinchingPhotos             []WinchingPhoto
	DispatchedAt               *time.Time
	LoadedAt                   *time.Time
	InTransitAt                *time.Time
	ArrivedAt                  *time.Time
	Version                    int
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}
