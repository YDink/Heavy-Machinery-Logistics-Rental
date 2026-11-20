package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/dinkmachinery/dink-core/internal/modules/logistics/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/database/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostgresLogisticsRepository struct {
	queries *db.Queries
}

func NewPostgresLogisticsRepository(queries *db.Queries) *PostgresLogisticsRepository {
	return &PostgresLogisticsRepository{queries: queries}
}

func numericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	f, err := n.Float64Value()
	if err != nil || !f.Valid {
		return 0
	}
	return f.Float64
}

func float64ToNumeric(v float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(fmt.Sprintf("%.2f", v))
	return n
}

func (r *PostgresLogisticsRepository) GetFleetLogisticsProfile(ctx context.Context, fleetID uuid.UUID) (*domain.FleetLogisticsProfile, error) {
	row, err := r.queries.GetFleetLogisticsContext(ctx, pgtype.UUID{Bytes: fleetID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLogisticsBookingNotFound
		}
		return nil, fmt.Errorf("failed to get fleet logistics profile: %w", err)
	}

	axleCount := 3
	if row.AxleCount != nil {
		axleCount = int(*row.AxleCount)
	}

	deckWellLength := int(math.Round(numericToFloat64(row.DeckWellLengthMm)))
	if deckWellLength == 0 {
		deckWellLength = int(math.Round(numericToFloat64(row.TransportLengthMm)))
	}
	deckWidth := int(math.Round(numericToFloat64(row.DeckWidthMm)))
	if deckWidth == 0 {
		deckWidth = int(math.Round(numericToFloat64(row.TransportWidthMm)))
	}
	deckGroundHeight := int(math.Round(numericToFloat64(row.DeckHeightGroundMm)))
	if deckGroundHeight == 0 {
		deckGroundHeight = 900
	}

	return &domain.FleetLogisticsProfile{
		FleetID:          uuid.UUID(row.ID.Bytes),
		OwnerID:          uuid.UUID(row.OwnerID.Bytes),
		Make:             row.Make,
		Model:            row.Model,
		RootCategoryCode: row.RootCategoryCode,
		SelfMoveCapable:  row.SelfMoveCapable,
		DailyRateMinor:   row.DailyRateMinor,
		Status:           row.Status,
		CountryCode:      row.CountryCode,
		AddressText:      row.AddressText,
		Location: domain.GeoPoint{
			Latitude:     row.Latitude,
			Longitude:    row.Longitude,
			LocationName: row.AddressText,
		},
		OwnerFullName:    row.OwnerFullName,
		OwnerPhoneNumber: row.OwnerPhoneNumber,
		MachineSpecs: domain.MachineLogisticsSpecs{
			OperatingWeightKg: numericToFloat64(row.OperatingWeightKg),
			TransportLengthMm: int(math.Round(numericToFloat64(row.TransportLengthMm))),
			TransportWidthMm:  int(math.Round(numericToFloat64(row.TransportWidthMm))),
			TransportHeightMm: int(math.Round(numericToFloat64(row.TransportHeightMm))),
		},
		TrailerProfile: domain.MatchedCarrier{
			FleetID:               uuid.UUID(row.ID.Bytes),
			OwnerID:               uuid.UUID(row.OwnerID.Bytes),
			Make:                  row.Make,
			Model:                 row.Model,
			OwnerTradingName:      row.OwnerFullName,
			OwnerPhoneNumber:      row.OwnerPhoneNumber,
			PayloadCapacityKg:     numericToFloat64(row.PayloadCapacityKg),
			DeckWellLengthMm:      deckWellLength,
			DeckWidthMm:           deckWidth,
			OutriggersExtensionMm: int(math.Round(numericToFloat64(row.OutriggersExtensionMm))),
			DeckHeightGroundMm:    deckGroundHeight,
			TareWeightKg:          numericToFloat64(row.TareWeightKg),
			AxleCount:             axleCount,
			DailyRateMinor:        row.DailyRateMinor,
		},
	}, nil
}

func (r *PostgresLogisticsRepository) GetCountryLogisticsLimits(ctx context.Context, countryCode string) (domain.CountryLogisticsLimits, error) {
	row, err := r.queries.GetCountryLogisticsLimits(ctx, countryCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.CountryLogisticsLimits{
				CountryCode:         countryCode,
				DefaultCurrency:     "ETB",
				MaxGvmKg:            domain.DefaultMaxGvmKg,
				MaxCorridorHeightMm: domain.DefaultMaxCorridorHeightMm,
				WideLoadWidthMm:     domain.DefaultWideLoadWidthMm,
			}, nil
		}
		return domain.CountryLogisticsLimits{}, fmt.Errorf("failed to query country logistics limits: %w", err)
	}
	return domain.CountryLogisticsLimits{
		CountryCode:         row.CountryCode,
		DefaultCurrency:     row.DefaultCurrency,
		MaxGvmKg:            numericToFloat64(row.MaxGvmKg),
		MaxCorridorHeightMm: int(row.MaxCorridorHeightMm),
		WideLoadWidthMm:     int(row.WideLoadWidthMm),
	}, nil
}

func (r *PostgresLogisticsRepository) MatchCompatibleLogisticsFleets(
	ctx context.Context,
	origin domain.GeoPoint,
	radiusMeters float64,
	countryCode string,
	excludeOwnerID uuid.UUID,
	specs domain.MachineLogisticsSpecs,
	limits domain.CountryLogisticsLimits,
	limit int,
) ([]domain.MatchedCarrier, error) {
	if limit <= 0 {
		limit = 20
	}
	minPayload := domain.ComputeMinimumRequiredPayloadKg(specs.OperatingWeightKg)

	rows, err := r.queries.MatchCompatibleLogisticsFleets(ctx, db.MatchCompatibleLogisticsFleetsParams{
		Longitude:           origin.Longitude,
		Latitude:            origin.Latitude,
		CountryCode:         countryCode,
		ExcludeOwnerID:      pgtype.UUID{Bytes: excludeOwnerID, Valid: true},
		RadiusMeters:        radiusMeters,
		MinPayloadKg:        float64ToNumeric(minPayload),
		RequiredLengthMm:    float64ToNumeric(float64(specs.TransportLengthMm)),
		RequiredWidthMm:     float64ToNumeric(float64(specs.TransportWidthMm)),
		MachineHeightMm:     float64ToNumeric(float64(specs.TransportHeightMm)),
		MaxCorridorHeightMm: float64ToNumeric(float64(limits.MaxCorridorHeightMm)),
		MachineWeightKg:     float64ToNumeric(specs.OperatingWeightKg),
		MaxGvmKg:            float64ToNumeric(limits.MaxGvmKg),
		QueryLimit:          int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to match compatible logistics fleets: %w", err)
	}

	result := make([]domain.MatchedCarrier, len(rows))
	for i, row := range rows {
		result[i] = domain.MatchedCarrier{
			FleetID:                 uuid.UUID(row.FleetID.Bytes),
			OwnerID:                 uuid.UUID(row.OwnerID.Bytes),
			Make:                    row.Make,
			Model:                   row.Model,
			OwnerTradingName:        row.OwnerTradingName,
			OwnerPhoneNumber:        row.OwnerPhoneNumber,
			PayloadCapacityKg:       numericToFloat64(row.PayloadCapacityKg),
			DeckWellLengthMm:        int(math.Round(numericToFloat64(row.DeckWellLengthMm))),
			DeckWidthMm:             int(math.Round(numericToFloat64(row.DeckWidthMm))),
			OutriggersExtensionMm:   int(math.Round(numericToFloat64(row.OutriggersExtensionMm))),
			DeckHeightGroundMm:      int(math.Round(numericToFloat64(row.DeckHeightGroundMm))),
			TareWeightKg:            numericToFloat64(row.TareWeightKg),
			AxleCount:               int(row.AxleCount),
			DailyRateMinor:          row.DailyRateMinor,
			Currency:                limits.DefaultCurrency,
			DistanceFromDepotMeters: row.DistanceFromDepotMeters,
			AverageRating:           numericToFloat64(row.AverageRating),
		}
	}
	return result, nil
}

func (r *PostgresLogisticsRepository) CreateLogisticsSubBookingAndDispatch(
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
	dispatch *domain.LogisticsDispatch,
) (uuid.UUID, error) {
	hirePeriod := pgtype.Range[pgtype.Date]{
		Lower:     pgtype.Date{Time: startDate, Valid: true},
		Upper:     pgtype.Date{Time: endDate, Valid: true},
		LowerType: pgtype.Inclusive,
		UpperType: pgtype.Exclusive,
		Valid:     true,
	}

	subBookingRow, err := r.queries.CreateLogisticsSubBooking(ctx, db.CreateLogisticsSubBookingParams{
		BookingReference:           subBookingRef,
		RenterID:                   pgtype.UUID{Bytes: renterOwnerID, Valid: true},
		OwnerID:                    pgtype.UUID{Bytes: carrierOwnerID, Valid: true},
		FleetID:                    pgtype.UUID{Bytes: logisticsFleetID, Valid: true},
		ParentBookingID:            pgtype.UUID{Bytes: parentBookingID, Valid: true},
		StartDate:                  pgtype.Date{Time: startDate, Valid: true},
		EndDate:                    pgtype.Date{Time: endDate, Valid: true},
		HirePeriod:                 hirePeriod,
		DurationDays:               int32(durationDays),
		DailyRateAppliedMinor:      dailyRateMinor,
		BaseHireFeeMinor:           freightHireFeeMinor,
		GrossRentalAmountMinor:     freightHireFeeMinor,
		PrepaymentAmountMinor:      freightHireFeeMinor,
		PlatformCommissionRateBps:  int32(commissionBps),
		PlatformCommissionMinor:    commissionOwedMinor,
		AbnormalLoadEscortRequired: abnormalLoad,
		Currency:                   currency,
		SiteLocationAddress:        siteAddress,
		Status:                     "CONFIRMED",
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create logistics sub-booking: %w", err)
	}

	subBookingID := uuid.UUID(subBookingRow.ID.Bytes)
	dispatchRow, err := r.queries.CreateLogisticsDispatch(ctx, db.CreateLogisticsDispatchParams{
		SubBookingID:               pgtype.UUID{Bytes: subBookingID, Valid: true},
		ParentBookingID:            pgtype.UUID{Bytes: parentBookingID, Valid: true},
		LogisticsFleetID:           pgtype.UUID{Bytes: logisticsFleetID, Valid: true},
		CarrierOwnerID:             pgtype.UUID{Bytes: carrierOwnerID, Valid: true},
		PickupDatetime:             pgtype.Timestamptz{Time: dispatch.PickupDateTime, Valid: true},
		OriginLatitude:             dispatch.Origin.Latitude,
		OriginLongitude:            dispatch.Origin.Longitude,
		OriginDepotName:            dispatch.Origin.LocationName,
		DestinationLatitude:        dispatch.Destination.Latitude,
		DestinationLongitude:       dispatch.Destination.Longitude,
		DestinationSiteName:        dispatch.Destination.LocationName,
		FreightHireFeeMinor:        freightHireFeeMinor,
		PlatformCommissionBps:      int32(commissionBps),
		CommissionOwedMinor:        commissionOwedMinor,
		AbnormalLoadEscortRequired: abnormalLoad,
		EscortFeeMinor:             dispatch.EscortFeeMinor,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create logistics dispatch: %w", err)
	}

	dispatch.ID = uuid.UUID(dispatchRow.ID.Bytes)
	dispatch.SubBookingID = subBookingID
	dispatch.CreatedAt = dispatchRow.CreatedAt.Time
	dispatch.UpdatedAt = dispatchRow.UpdatedAt.Time
	dispatch.Version = int(dispatchRow.Version)
	return subBookingID, nil
}

func toDomainDispatch(row db.LogisticsDispatch) *domain.LogisticsDispatch {
	var photos []domain.WinchingPhoto
	if len(row.WinchingPhotosJson) > 0 {
		_ = json.Unmarshal(row.WinchingPhotosJson, &photos)
	}
	res := &domain.LogisticsDispatch{
		ID:                         uuid.UUID(row.ID.Bytes),
		SubBookingID:               uuid.UUID(row.SubBookingID.Bytes),
		ParentBookingID:            uuid.UUID(row.ParentBookingID.Bytes),
		LogisticsFleetID:           uuid.UUID(row.LogisticsFleetID.Bytes),
		CarrierOwnerID:             uuid.UUID(row.CarrierOwnerID.Bytes),
		PickupDateTime:             row.PickupDatetime.Time,
		Origin:                     domain.GeoPoint{Latitude: row.OriginLatitude, Longitude: row.OriginLongitude, LocationName: row.OriginDepotName},
		Destination:                domain.GeoPoint{Latitude: row.DestinationLatitude, Longitude: row.DestinationLongitude, LocationName: row.DestinationSiteName},
		FreightHireFeeMinor:        row.FreightHireFeeMinor,
		PlatformCommissionBps:      int(row.PlatformCommissionBps),
		CommissionOwedMinor:        row.CommissionOwedMinor,
		Currency:                   "ETB",
		AbnormalLoadEscortRequired: row.AbnormalLoadEscortRequired,
		EscortFeeMinor:             row.EscortFeeMinor,
		WinchingPhotos:             photos,
		Version:                    int(row.Version),
		CreatedAt:                  row.CreatedAt.Time,
		UpdatedAt:                  row.UpdatedAt.Time,
	}
	if row.DispatchedAt.Valid {
		t := row.DispatchedAt.Time
		res.DispatchedAt = &t
	}
	if row.LoadedAt.Valid {
		t := row.LoadedAt.Time
		res.LoadedAt = &t
	}
	if row.InTransitAt.Valid {
		t := row.InTransitAt.Time
		res.InTransitAt = &t
	}
	if row.ArrivedAt.Valid {
		t := row.ArrivedAt.Time
		res.ArrivedAt = &t
	}
	return res
}

func (r *PostgresLogisticsRepository) GetLogisticsDispatchBySubBookingID(ctx context.Context, subBookingID uuid.UUID) (*domain.LogisticsDispatch, error) {
	row, err := r.queries.GetLogisticsDispatchBySubBookingID(ctx, pgtype.UUID{Bytes: subBookingID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLogisticsBookingNotFound
		}
		return nil, fmt.Errorf("failed to get dispatch by sub-booking: %w", err)
	}
	return toDomainDispatch(row), nil
}

func (r *PostgresLogisticsRepository) GetLogisticsDispatchByParentBookingID(ctx context.Context, parentBookingID uuid.UUID) (*domain.LogisticsDispatch, error) {
	row, err := r.queries.GetLogisticsDispatchByParentBookingID(ctx, pgtype.UUID{Bytes: parentBookingID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLogisticsBookingNotFound
		}
		return nil, fmt.Errorf("failed to get dispatch by parent booking: %w", err)
	}
	return toDomainDispatch(row), nil
}

func (r *PostgresLogisticsRepository) UpdateLogisticsDispatchLoading(ctx context.Context, subBookingID uuid.UUID, photos []domain.WinchingPhoto) (*domain.LogisticsDispatch, error) {
	rawPhotos, err := json.Marshal(photos)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal winching photos: %w", err)
	}
	row, err := r.queries.UpdateLogisticsDispatchLoading(ctx, db.UpdateLogisticsDispatchLoadingParams{
		SubBookingID:       pgtype.UUID{Bytes: subBookingID, Valid: true},
		WinchingPhotosJson: rawPhotos,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update dispatch loading inspection: %w", err)
	}
	return toDomainDispatch(row), nil
}

func (r *PostgresLogisticsRepository) UpdateLogisticsDispatchMilestone(ctx context.Context, subBookingID uuid.UUID, stage string) (*domain.LogisticsDispatch, error) {
	row, err := r.queries.UpdateLogisticsDispatchMilestone(ctx, db.UpdateLogisticsDispatchMilestoneParams{
		SubBookingID: pgtype.UUID{Bytes: subBookingID, Valid: true},
		Column2:      stage,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update dispatch milestone: %w", err)
	}
	return toDomainDispatch(row), nil
}

func toDomainExternalTransport(row db.ExternalTransport) *domain.ExternalTransport {
	return &domain.ExternalTransport{
		ID:                    uuid.UUID(row.ID.Bytes),
		ParentBookingID:       uuid.UUID(row.ParentBookingID.Bytes),
		RegisteredBy:          uuid.UUID(row.RegisteredBy.Bytes),
		CarrierName:           row.CarrierName,
		DriverFullName:        row.DriverFullName,
		DriverPhoneNumber:     row.DriverPhoneNumber,
		PrimeMoverPlateNumber: row.PrimeMoverPlateNumber,
		EstimatedArrivalAt:    row.EstimatedArrivalAt.Time,
		CreatedAt:             row.CreatedAt.Time,
		UpdatedAt:             row.UpdatedAt.Time,
	}
}

func (r *PostgresLogisticsRepository) UpsertExternalTransport(ctx context.Context, ext *domain.ExternalTransport) (*domain.ExternalTransport, error) {
	row, err := r.queries.UpsertExternalTransport(ctx, db.UpsertExternalTransportParams{
		ParentBookingID:       pgtype.UUID{Bytes: ext.ParentBookingID, Valid: true},
		RegisteredBy:          pgtype.UUID{Bytes: ext.RegisteredBy, Valid: true},
		CarrierName:           ext.CarrierName,
		DriverFullName:        ext.DriverFullName,
		DriverPhoneNumber:     ext.DriverPhoneNumber,
		PrimeMoverPlateNumber: ext.PrimeMoverPlateNumber,
		EstimatedArrivalAt:    pgtype.Timestamptz{Time: ext.EstimatedArrivalAt, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upsert external transport: %w", err)
	}
	return toDomainExternalTransport(row), nil
}

func (r *PostgresLogisticsRepository) GetExternalTransportByParentBookingID(ctx context.Context, parentBookingID uuid.UUID) (*domain.ExternalTransport, error) {
	row, err := r.queries.GetExternalTransportByParentBookingID(ctx, pgtype.UUID{Bytes: parentBookingID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLogisticsBookingNotFound
		}
		return nil, fmt.Errorf("failed to get external transport: %w", err)
	}
	return toDomainExternalTransport(row), nil
}

func toDomainTransshipment(row db.LogisticsTransshipment) *domain.LogisticsTransshipment {
	var photos []string
	if len(row.PhotosJson) > 0 {
		_ = json.Unmarshal(row.PhotosJson, &photos)
	}
	res := &domain.LogisticsTransshipment{
		ID:                       uuid.UUID(row.ID.Bytes),
		SubBookingID:             uuid.UUID(row.SubBookingID.Bytes),
		ParentBookingID:          uuid.UUID(row.ParentBookingID.Bytes),
		ReportedBy:               uuid.UUID(row.ReportedBy.Bytes),
		BreakdownLocation:        domain.GeoPoint{Latitude: row.BreakdownLatitude, Longitude: row.BreakdownLongitude},
		FailureCategory:          domain.TransshipmentFailureCategory(row.FailureCategory),
		FailureDescription:       row.FailureDescription,
		Photos:                   photos,
		HotSwapSLADeadline:       row.HotSwapSlaDeadline.Time,
		TransshipmentBountyMinor: row.TransshipmentBountyMinor,
		Currency:                 "ETB",
		Status:                   domain.TransshipmentStatus(row.Status),
		CreatedAt:                row.CreatedAt.Time,
		UpdatedAt:                row.UpdatedAt.Time,
	}
	if row.ReplacementFleetID.Valid {
		id := uuid.UUID(row.ReplacementFleetID.Bytes)
		res.ReplacementFleetID = &id
	}
	if row.ReplacementCarrierID.Valid {
		id := uuid.UUID(row.ReplacementCarrierID.Bytes)
		res.ReplacementCarrierID = &id
	}
	if row.AcceptedAt.Valid {
		t := row.AcceptedAt.Time
		res.AcceptedAt = &t
	}
	if row.CompletedAt.Valid {
		t := row.CompletedAt.Time
		res.CompletedAt = &t
	}
	return res
}

func (r *PostgresLogisticsRepository) CreateLogisticsTransshipment(ctx context.Context, t *domain.LogisticsTransshipment) (*domain.LogisticsTransshipment, error) {
	rawPhotos, err := json.Marshal(t.Photos)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transshipment photos: %w", err)
	}
	row, err := r.queries.CreateLogisticsTransshipment(ctx, db.CreateLogisticsTransshipmentParams{
		SubBookingID:             pgtype.UUID{Bytes: t.SubBookingID, Valid: true},
		ParentBookingID:          pgtype.UUID{Bytes: t.ParentBookingID, Valid: true},
		ReportedBy:               pgtype.UUID{Bytes: t.ReportedBy, Valid: true},
		BreakdownLatitude:        t.BreakdownLocation.Latitude,
		BreakdownLongitude:       t.BreakdownLocation.Longitude,
		FailureCategory:          string(t.FailureCategory),
		FailureDescription:       t.FailureDescription,
		PhotosJson:               rawPhotos,
		HotSwapSlaDeadline:       pgtype.Timestamptz{Time: t.HotSwapSLADeadline, Valid: true},
		TransshipmentBountyMinor: t.TransshipmentBountyMinor,
		Status:                   string(t.Status),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create logistics transshipment: %w", err)
	}
	return toDomainTransshipment(row), nil
}

func (r *PostgresLogisticsRepository) GetOpenTransshipmentBySubBookingID(ctx context.Context, subBookingID uuid.UUID) (*domain.LogisticsTransshipment, error) {
	row, err := r.queries.GetOpenTransshipmentBySubBookingID(ctx, pgtype.UUID{Bytes: subBookingID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTransshipmentNotOpen
		}
		return nil, fmt.Errorf("failed to get open transshipment: %w", err)
	}
	return toDomainTransshipment(row), nil
}

func (r *PostgresLogisticsRepository) AcceptLogisticsTransshipment(ctx context.Context, transshipmentID, replacementFleetID, replacementCarrierID uuid.UUID) (*domain.LogisticsTransshipment, error) {
	row, err := r.queries.AcceptLogisticsTransshipment(ctx, db.AcceptLogisticsTransshipmentParams{
		ID:                   pgtype.UUID{Bytes: transshipmentID, Valid: true},
		ReplacementFleetID:   pgtype.UUID{Bytes: replacementFleetID, Valid: true},
		ReplacementCarrierID: pgtype.UUID{Bytes: replacementCarrierID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to accept transshipment: %w", err)
	}
	return toDomainTransshipment(row), nil
}

func (r *PostgresLogisticsRepository) ExpireLogisticsTransshipment(ctx context.Context, transshipmentID uuid.UUID) error {
	_, err := r.queries.ExpireLogisticsTransshipment(ctx, pgtype.UUID{Bytes: transshipmentID, Valid: true})
	return err
}

func toDomainReverseCascade(row db.LogisticsReverseCascade) *domain.LogisticsReverseCascade {
	return &domain.LogisticsReverseCascade{
		ID:                             uuid.UUID(row.ID.Bytes),
		ParentBookingID:                uuid.UUID(row.ParentBookingID.Bytes),
		SubBookingID:                   uuid.UUID(row.SubBookingID.Bytes),
		TriggeredBy:                    uuid.UUID(row.TriggeredBy.Bytes),
		Resolution:                     row.Resolution,
		Reason:                         row.Reason,
		RenterRefundDueMinor:           row.RenterRefundDueMinor,
		OwnerPenaltyLeviedMinor:        row.OwnerPenaltyLeviedMinor,
		OwnerCommissionRestitutionMinor: row.OwnerCommissionRestitutionMinor,
		CarrierDefaultPenaltyMinor:     row.CarrierDefaultPenaltyMinor,
		CarrierRatingDeduction:         numericToFloat64(row.CarrierRatingDeduction),
		Currency:                       "ETB",
		CreatedAt:                      row.CreatedAt.Time,
	}
}

func (r *PostgresLogisticsRepository) CreateLogisticsReverseCascade(ctx context.Context, rc *domain.LogisticsReverseCascade) (*domain.LogisticsReverseCascade, error) {
	row, err := r.queries.CreateLogisticsReverseCascade(ctx, db.CreateLogisticsReverseCascadeParams{
		ParentBookingID:                pgtype.UUID{Bytes: rc.ParentBookingID, Valid: true},
		SubBookingID:                   pgtype.UUID{Bytes: rc.SubBookingID, Valid: true},
		TriggeredBy:                    pgtype.UUID{Bytes: rc.TriggeredBy, Valid: true},
		Resolution:                     rc.Resolution,
		Reason:                         rc.Reason,
		RenterRefundDueMinor:           rc.RenterRefundDueMinor,
		OwnerPenaltyLeviedMinor:        rc.OwnerPenaltyLeviedMinor,
		OwnerCommissionRestitutionMinor: rc.OwnerCommissionRestitutionMinor,
		CarrierDefaultPenaltyMinor:     rc.CarrierDefaultPenaltyMinor,
		CarrierRatingDeduction:         float64ToNumeric(rc.CarrierRatingDeduction),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to record reverse cascade: %w", err)
	}
	return toDomainReverseCascade(row), nil
}

func (r *PostgresLogisticsRepository) GetLogisticsReverseCascadeByParentBookingID(ctx context.Context, parentBookingID uuid.UUID) (*domain.LogisticsReverseCascade, error) {
	row, err := r.queries.GetLogisticsReverseCascadeByParentBookingID(ctx, pgtype.UUID{Bytes: parentBookingID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLogisticsBookingNotFound
		}
		return nil, fmt.Errorf("failed to get reverse cascade: %w", err)
	}
	return toDomainReverseCascade(row), nil
}

func (r *PostgresLogisticsRepository) ReleaseAvailabilityAndDeductCarrierRating(ctx context.Context, parentBookingID, subBookingID, carrierOwnerID uuid.UUID, ratingDeduction float64) error {
	_ = r.queries.ReleaseFleetAvailabilityByReference(ctx, pgtype.UUID{Bytes: parentBookingID, Valid: true})
	_ = r.queries.ReleaseFleetAvailabilityByReference(ctx, pgtype.UUID{Bytes: subBookingID, Valid: true})
	return r.queries.DeductOwnerAggregateRating(ctx, db.DeductOwnerAggregateRatingParams{
		OwnerID: pgtype.UUID{Bytes: carrierOwnerID, Valid: true},
		Column2: float64ToNumeric(ratingDeduction),
	})
}
