package domain

import (
	"time"

	"github.com/google/uuid"
)

// TransshipmentFailureCategory defines allowable in-transit carrier failure categories.
type TransshipmentFailureCategory string

const (
	FailurePrimeMoverEngineBlown TransshipmentFailureCategory = "PRIME_MOVER_ENGINE_BLOWN"
	FailureTrailerAxleFailure    TransshipmentFailureCategory = "TRAILER_AXLE_FAILURE"
	FailureHydraulicRampFailure  TransshipmentFailureCategory = "HYDRAULIC_RAMP_FAILURE"
	FailureCorridorBlockade      TransshipmentFailureCategory = "CORRIDOR_BLOCKADE"
	FailureOther                 TransshipmentFailureCategory = "OTHER"
)

// TransshipmentStatus defines the lifecycle of an emergency roadside transshipment request.
type TransshipmentStatus string

const (
	TransshipmentOpen      TransshipmentStatus = "OPEN"
	TransshipmentAccepted  TransshipmentStatus = "ACCEPTED"
	TransshipmentCompleted TransshipmentStatus = "COMPLETED"
	TransshipmentExpired   TransshipmentStatus = "EXPIRED"
	TransshipmentFailed    TransshipmentStatus = "FAILED"
)

// LogisticsTransshipment represents a mid-transit breakdown alert and hot-swap SLA docket.
type LogisticsTransshipment struct {
	ID                       uuid.UUID
	SubBookingID             uuid.UUID
	ParentBookingID          uuid.UUID
	ReportedBy               uuid.UUID
	BreakdownLocation        GeoPoint
	FailureCategory          TransshipmentFailureCategory
	FailureDescription       string
	Photos                   []string
	HotSwapSLADeadline       time.Time
	TransshipmentBountyMinor int64
	Currency                 string
	ReplacementFleetID       *uuid.UUID
	ReplacementCarrierID     *uuid.UUID
	Status                   TransshipmentStatus
	AcceptedAt               *time.Time
	CompletedAt              *time.Time
	CreatedAt                time.Time
	UpdatedAt                time.Time
}
