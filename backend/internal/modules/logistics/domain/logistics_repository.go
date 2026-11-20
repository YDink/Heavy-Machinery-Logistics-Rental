package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// FleetLogisticsProfile carries the dimensional, operational, and ownership metadata of a fleet unit.
type FleetLogisticsProfile struct {
	FleetID           uuid.UUID
	OwnerID           uuid.UUID
	Make              string
	Model             string
	RootCategoryCode  string
	SelfMoveCapable   bool
	DailyRateMinor    int64
	Status            string
	CountryCode       string
	AddressText       string
	Location          GeoPoint
	OwnerFullName     string
	OwnerPhoneNumber  string
	MachineSpecs      MachineLogisticsSpecs
	TrailerProfile    MatchedCarrier
}

// LogisticsRepository defines the pure domain persistence port for Partition 05-08.
type LogisticsRepository interface {
	GetFleetLogisticsProfile(ctx context.Context, fleetID uuid.UUID) (*FleetLogisticsProfile, error)
	GetCountryLogisticsLimits(ctx context.Context, countryCode string) (CountryLogisticsLimits, error)
	MatchCompatibleLogisticsFleets(
		ctx context.Context,
		origin GeoPoint,
		radiusMeters float64,
		countryCode string,
		excludeOwnerID uuid.UUID,
		specs MachineLogisticsSpecs,
		limits CountryLogisticsLimits,
		limit int,
	) ([]MatchedCarrier, error)

	CreateLogisticsSubBookingAndDispatch(
		ctx context.Context,
		subBookingRef string,
		parentBookingID uuid.UUID,
		renterOwnerID uuid.UUID,
		carrierOwnerID uuid.UUID,
		logisticsFleetID uuid.UUID,
		startDate time.Time,
		endDate time.Time,
		durationDays int,
		dailyRateMinor int64,
		freightHireFeeMinor int64,
		commissionBps int,
		commissionOwedMinor int64,
		abnormalLoad bool,
		currency string,
		siteAddress string,
		dispatch *LogisticsDispatch,
	) (uuid.UUID, error)

	GetLogisticsDispatchBySubBookingID(ctx context.Context, subBookingID uuid.UUID) (*LogisticsDispatch, error)
	GetLogisticsDispatchByParentBookingID(ctx context.Context, parentBookingID uuid.UUID) (*LogisticsDispatch, error)
	UpdateLogisticsDispatchLoading(ctx context.Context, subBookingID uuid.UUID, photos []WinchingPhoto) (*LogisticsDispatch, error)
	UpdateLogisticsDispatchMilestone(ctx context.Context, subBookingID uuid.UUID, stage string) (*LogisticsDispatch, error)

	UpsertExternalTransport(ctx context.Context, ext *ExternalTransport) (*ExternalTransport, error)
	GetExternalTransportByParentBookingID(ctx context.Context, parentBookingID uuid.UUID) (*ExternalTransport, error)

	CreateLogisticsTransshipment(ctx context.Context, t *LogisticsTransshipment) (*LogisticsTransshipment, error)
	GetOpenTransshipmentBySubBookingID(ctx context.Context, subBookingID uuid.UUID) (*LogisticsTransshipment, error)
	AcceptLogisticsTransshipment(ctx context.Context, transshipmentID, replacementFleetID, replacementCarrierID uuid.UUID) (*LogisticsTransshipment, error)
	ExpireLogisticsTransshipment(ctx context.Context, transshipmentID uuid.UUID) error

	CreateLogisticsReverseCascade(ctx context.Context, rc *LogisticsReverseCascade) (*LogisticsReverseCascade, error)
	GetLogisticsReverseCascadeByParentBookingID(ctx context.Context, parentBookingID uuid.UUID) (*LogisticsReverseCascade, error)
	ReleaseAvailabilityAndDeductCarrierRating(ctx context.Context, parentBookingID, subBookingID, carrierOwnerID uuid.UUID, ratingDeduction float64) error
}
