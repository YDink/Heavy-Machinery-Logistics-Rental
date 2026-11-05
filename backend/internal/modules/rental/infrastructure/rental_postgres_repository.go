package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dinkmachinery/dink-core/internal/modules/rental/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/database/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostgresRentalRepository struct {
	queries *db.Queries
}

func NewPostgresRentalRepository(queries *db.Queries) *PostgresRentalRepository {
	return &PostgresRentalRepository{queries: queries}
}

func floatToNumeric(val float64) pgtype.Numeric {
	var num pgtype.Numeric
	_ = num.Scan(fmt.Sprintf("%.2f", val))
	return num
}

func numericToFloat(num pgtype.Numeric) float64 {
	val, err := num.Float64Value()
	if err != nil || !val.Valid {
		return 0.0
	}
	return val.Float64
}

func (r *PostgresRentalRepository) CreateHandoverInspection(ctx context.Context, inspection *domain.HandoverInspection) error {
	photosBytes, err := json.Marshal(inspection.Photos)
	if err != nil {
		return fmt.Errorf("failed to marshal inspection photos: %w", err)
	}

	var notes *string
	if inspection.Notes != "" {
		notes = &inspection.Notes
	}

	row, err := r.queries.CreateHandoverInspection(ctx, db.CreateHandoverInspectionParams{
		BookingID:         pgtype.UUID{Bytes: inspection.BookingID, Valid: true},
		InspectorID:       pgtype.UUID{Bytes: inspection.InspectorID, Valid: true},
		CheckInMeterHours: floatToNumeric(inspection.CheckInMeterHours),
		FluidLevelsOk:     inspection.FluidLevelsOk,
		SafetyGlassIntact: inspection.SafetyGlassIntact,
		Notes:             notes,
		PhotosJson:        photosBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to insert handover inspection: %w", err)
	}

	inspection.ID = uuid.UUID(row.ID.Bytes)
	inspection.CreatedAt = row.CreatedAt.Time
	return nil
}

func (r *PostgresRentalRepository) GetHandoverInspection(ctx context.Context, bookingID uuid.UUID) (*domain.HandoverInspection, error) {
	row, err := r.queries.GetHandoverInspectionByBookingID(ctx, pgtype.UUID{Bytes: bookingID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("handover inspection not found: %w", err)
	}

	var photos []domain.InspectionPhoto
	_ = json.Unmarshal(row.PhotosJson, &photos)

	var notes string
	if row.Notes != nil {
		notes = *row.Notes
	}

	return &domain.HandoverInspection{
		ID:                uuid.UUID(row.ID.Bytes),
		BookingID:         uuid.UUID(row.BookingID.Bytes),
		InspectorID:       uuid.UUID(row.InspectorID.Bytes),
		CheckInMeterHours: numericToFloat(row.CheckInMeterHours),
		FluidLevelsOk:     row.FluidLevelsOk,
		SafetyGlassIntact: row.SafetyGlassIntact,
		Notes:             notes,
		Photos:            photos,
		CreatedAt:         row.CreatedAt.Time,
	}, nil
}

func (r *PostgresRentalRepository) CreateHandoverChallenge(ctx context.Context, challenge *domain.HandoverChallenge) error {
	row, err := r.queries.CreateHandoverChallenge(ctx, db.CreateHandoverChallengeParams{
		BookingID:    pgtype.UUID{Bytes: challenge.BookingID, Valid: true},
		ChallengePin: challenge.ChallengePIN,
		ExpiresAt:    pgtype.Timestamptz{Time: challenge.ExpiresAt, Valid: true},
		MaxAttempts:  int32(challenge.MaxAttempts),
	})
	if err != nil {
		return fmt.Errorf("failed to create handover challenge: %w", err)
	}

	challenge.ID = uuid.UUID(row.ID.Bytes)
	challenge.CreatedAt = row.CreatedAt.Time
	return nil
}

func (r *PostgresRentalRepository) GetActiveHandoverChallenge(ctx context.Context, bookingID uuid.UUID) (*domain.HandoverChallenge, error) {
	row, err := r.queries.GetActiveHandoverChallenge(ctx, pgtype.UUID{Bytes: bookingID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("active challenge not found: %w", err)
	}

	res := &domain.HandoverChallenge{
		ID:           uuid.UUID(row.ID.Bytes),
		BookingID:    uuid.UUID(row.BookingID.Bytes),
		ChallengePIN: row.ChallengePin,
		ExpiresAt:    row.ExpiresAt.Time,
		AttemptCount: int(row.AttemptCount),
		MaxAttempts:  int(row.MaxAttempts),
		IsVerified:   row.IsVerified,
		CreatedAt:    row.CreatedAt.Time,
	}

	if row.TouchSignatureBase64 != nil {
		res.TouchSignatureBase64 = *row.TouchSignatureBase64
	}
	if row.ContractorLat != nil {
		res.ContractorLat = row.ContractorLat
	}
	if row.ContractorLng != nil {
		res.ContractorLng = row.ContractorLng
	}
	if row.VerifiedAt.Valid {
		t := row.VerifiedAt.Time
		res.VerifiedAt = &t
	}

	return res, nil
}

func (r *PostgresRentalRepository) IncrementHandoverChallengeAttempts(ctx context.Context, challengeID uuid.UUID) (*domain.HandoverChallenge, error) {
	row, err := r.queries.IncrementHandoverChallengeAttempts(ctx, pgtype.UUID{Bytes: challengeID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to increment challenge attempts: %w", err)
	}

	return &domain.HandoverChallenge{
		ID:           uuid.UUID(row.ID.Bytes),
		BookingID:    uuid.UUID(row.BookingID.Bytes),
		ChallengePIN: row.ChallengePin,
		ExpiresAt:    row.ExpiresAt.Time,
		AttemptCount: int(row.AttemptCount),
		MaxAttempts:  int(row.MaxAttempts),
		IsVerified:   row.IsVerified,
		CreatedAt:    row.CreatedAt.Time,
	}, nil
}

func (r *PostgresRentalRepository) VerifyHandoverChallenge(ctx context.Context, challengeID uuid.UUID, touchSignatureBase64 string, lat, lng *float64) error {
	var sig *string
	if touchSignatureBase64 != "" {
		sig = &touchSignatureBase64
	}
	_, err := r.queries.VerifyHandoverChallenge(ctx, db.VerifyHandoverChallengeParams{
		ID:                   pgtype.UUID{Bytes: challengeID, Valid: true},
		TouchSignatureBase64: sig,
		ContractorLat:        lat,
		ContractorLng:        lng,
	})
	return err
}

func (r *PostgresRentalRepository) CreateRentalShift(ctx context.Context, shift *domain.RentalShift) error {
	var opID pgtype.UUID
	if shift.OperatorID != nil {
		opID = pgtype.UUID{Bytes: *shift.OperatorID, Valid: true}
	}

	row, err := r.queries.CreateRentalShift(ctx, db.CreateRentalShiftParams{
		BookingID:        pgtype.UUID{Bytes: shift.BookingID, Valid: true},
		OperatorID:       opID,
		ShiftDate:        pgtype.Date{Time: shift.ShiftDate, Valid: true},
		StartTime:        pgtype.Timestamptz{Time: shift.StartTime, Valid: true},
		StartMeterHours:  floatToNumeric(shift.StartMeterHours),
		AllocatedHours:   floatToNumeric(shift.AllocatedHours),
		AccumulatedHours: floatToNumeric(shift.AccumulatedHours),
		ExcessHours:      floatToNumeric(shift.ExcessHours),
		IsPaused:         shift.IsPaused,
	})
	if err != nil {
		return fmt.Errorf("failed to create rental shift: %w", err)
	}

	shift.ID = uuid.UUID(row.ID.Bytes)
	shift.CreatedAt = row.CreatedAt.Time
	shift.UpdatedAt = row.UpdatedAt.Time
	return nil
}

func (r *PostgresRentalRepository) GetActiveRentalShift(ctx context.Context, bookingID uuid.UUID) (*domain.RentalShift, error) {
	row, err := r.queries.GetActiveRentalShift(ctx, pgtype.UUID{Bytes: bookingID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("active rental shift not found: %w", err)
	}

	res := &domain.RentalShift{
		ID:               uuid.UUID(row.ID.Bytes),
		BookingID:        uuid.UUID(row.BookingID.Bytes),
		ShiftDate:        row.ShiftDate.Time,
		StartTime:        row.StartTime.Time,
		StartMeterHours:  numericToFloat(row.StartMeterHours),
		AllocatedHours:   numericToFloat(row.AllocatedHours),
		AccumulatedHours: numericToFloat(row.AccumulatedHours),
		ExcessHours:      numericToFloat(row.ExcessHours),
		IsPaused:         row.IsPaused,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
	}

	if row.OperatorID.Valid {
		id := uuid.UUID(row.OperatorID.Bytes)
		res.OperatorID = &id
	}
	if row.EndMeterHours.Valid {
		f := numericToFloat(row.EndMeterHours)
		res.EndMeterHours = &f
	}
	if row.EndTime.Valid {
		t := row.EndTime.Time
		res.EndTime = &t
	}

	return res, nil
}

func (r *PostgresRentalRepository) UpdateRentalShiftProgress(ctx context.Context, shiftID uuid.UUID, accumulatedHours, excessHours float64, endMeter *float64, isPaused bool) error {
	var em pgtype.Numeric
	if endMeter != nil {
		em = floatToNumeric(*endMeter)
	}
	_, err := r.queries.UpdateRentalShiftProgress(ctx, db.UpdateRentalShiftProgressParams{
		ID:               pgtype.UUID{Bytes: shiftID, Valid: true},
		AccumulatedHours: floatToNumeric(accumulatedHours),
		ExcessHours:      floatToNumeric(excessHours),
		EndMeterHours:    em,
		IsPaused:         isPaused,
	})
	return err
}

func (r *PostgresRentalRepository) EndRentalShift(ctx context.Context, shiftID uuid.UUID, endMeter float64, accumulatedHours, excessHours float64) error {
	_, err := r.queries.EndRentalShift(ctx, db.EndRentalShiftParams{
		ID:               pgtype.UUID{Bytes: shiftID, Valid: true},
		EndMeterHours:    floatToNumeric(endMeter),
		AccumulatedHours: floatToNumeric(accumulatedHours),
		ExcessHours:      floatToNumeric(excessHours),
	})
	return err
}

func (r *PostgresRentalRepository) ListRentalShifts(ctx context.Context, bookingID uuid.UUID) ([]domain.RentalShift, error) {
	rows, err := r.queries.ListRentalShiftsByBookingID(ctx, pgtype.UUID{Bytes: bookingID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list rental shifts: %w", err)
	}

	var res []domain.RentalShift
	for _, row := range rows {
		s := domain.RentalShift{
			ID:               uuid.UUID(row.ID.Bytes),
			BookingID:        uuid.UUID(row.BookingID.Bytes),
			ShiftDate:        row.ShiftDate.Time,
			StartTime:        row.StartTime.Time,
			StartMeterHours:  numericToFloat(row.StartMeterHours),
			AllocatedHours:   numericToFloat(row.AllocatedHours),
			AccumulatedHours: numericToFloat(row.AccumulatedHours),
			ExcessHours:      numericToFloat(row.ExcessHours),
			IsPaused:         row.IsPaused,
			CreatedAt:        row.CreatedAt.Time,
			UpdatedAt:        row.UpdatedAt.Time,
		}
		if row.OperatorID.Valid {
			id := uuid.UUID(row.OperatorID.Bytes)
			s.OperatorID = &id
		}
		if row.EndMeterHours.Valid {
			f := numericToFloat(row.EndMeterHours)
			s.EndMeterHours = &f
		}
		if row.EndTime.Valid {
			t := row.EndTime.Time
			s.EndTime = &t
		}
		res = append(res, s)
	}
	return res, nil
}

func (r *PostgresRentalRepository) CreateRentalExtension(ctx context.Context, ext *domain.RentalExtension) error {
	row, err := r.queries.CreateRentalExtension(ctx, db.CreateRentalExtensionParams{
		BookingID:                  pgtype.UUID{Bytes: ext.BookingID, Valid: true},
		RequestedByID:              pgtype.UUID{Bytes: ext.RequestedByID, Valid: true},
		OriginalEndDate:            pgtype.Date{Time: ext.OriginalEndDate, Valid: true},
		NewEndDate:                 pgtype.Date{Time: ext.NewEndDate, Valid: true},
		ExtensionDays:              int32(ext.ExtensionDays),
		IncrementalGrossMinor:      ext.IncrementalGrossMinor,
		IncrementalPrepaymentMinor: ext.IncrementalPrepaymentMinor,
		Currency:                   ext.Currency,
		Status:                     ext.Status,
	})
	if err != nil {
		return fmt.Errorf("failed to create rental extension: %w", err)
	}

	ext.ID = uuid.UUID(row.ID.Bytes)
	ext.CreatedAt = row.CreatedAt.Time
	ext.UpdatedAt = row.UpdatedAt.Time
	return nil
}

func (r *PostgresRentalRepository) GetRentalExtension(ctx context.Context, extensionID uuid.UUID) (*domain.RentalExtension, error) {
	row, err := r.queries.GetRentalExtensionByID(ctx, pgtype.UUID{Bytes: extensionID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("rental extension not found: %w", err)
	}

	res := &domain.RentalExtension{
		ID:                         uuid.UUID(row.ID.Bytes),
		BookingID:                  uuid.UUID(row.BookingID.Bytes),
		RequestedByID:              uuid.UUID(row.RequestedByID.Bytes),
		OriginalEndDate:            row.OriginalEndDate.Time,
		NewEndDate:                 row.NewEndDate.Time,
		ExtensionDays:              int(row.ExtensionDays),
		IncrementalGrossMinor:      row.IncrementalGrossMinor,
		IncrementalPrepaymentMinor: row.IncrementalPrepaymentMinor,
		Currency:                   row.Currency,
		Status:                     row.Status,
		CreatedAt:                  row.CreatedAt.Time,
		UpdatedAt:                  row.UpdatedAt.Time,
	}

	if row.OwnerDecidedAt.Valid {
		t := row.OwnerDecidedAt.Time
		res.OwnerDecidedAt = &t
	}
	if row.OwnerDecisionNotes != nil {
		res.OwnerDecisionNotes = *row.OwnerDecisionNotes
	}
	if row.SlipStorageUri != nil {
		res.SlipStorageURI = *row.SlipStorageUri
	}
	if row.SlipFileHash != nil {
		res.SlipFileHash = *row.SlipFileHash
	}
	if row.ConfirmedAt.Valid {
		t := row.ConfirmedAt.Time
		res.ConfirmedAt = &t
	}

	return res, nil
}

func (r *PostgresRentalRepository) ListRentalExtensions(ctx context.Context, bookingID uuid.UUID) ([]domain.RentalExtension, error) {
	rows, err := r.queries.ListRentalExtensionsByBookingID(ctx, pgtype.UUID{Bytes: bookingID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list rental extensions: %w", err)
	}

	var res []domain.RentalExtension
	for _, row := range rows {
		e := domain.RentalExtension{
			ID:                         uuid.UUID(row.ID.Bytes),
			BookingID:                  uuid.UUID(row.BookingID.Bytes),
			RequestedByID:              uuid.UUID(row.RequestedByID.Bytes),
			OriginalEndDate:            row.OriginalEndDate.Time,
			NewEndDate:                 row.NewEndDate.Time,
			ExtensionDays:              int(row.ExtensionDays),
			IncrementalGrossMinor:      row.IncrementalGrossMinor,
			IncrementalPrepaymentMinor: row.IncrementalPrepaymentMinor,
			Currency:                   row.Currency,
			Status:                     row.Status,
			CreatedAt:                  row.CreatedAt.Time,
			UpdatedAt:                  row.UpdatedAt.Time,
		}
		if row.OwnerDecidedAt.Valid {
			t := row.OwnerDecidedAt.Time
			e.OwnerDecidedAt = &t
		}
		if row.OwnerDecisionNotes != nil {
			e.OwnerDecisionNotes = *row.OwnerDecisionNotes
		}
		if row.SlipStorageUri != nil {
			e.SlipStorageURI = *row.SlipStorageUri
		}
		if row.SlipFileHash != nil {
			e.SlipFileHash = *row.SlipFileHash
		}
		if row.ConfirmedAt.Valid {
			t := row.ConfirmedAt.Time
			e.ConfirmedAt = &t
		}
		res = append(res, e)
	}
	return res, nil
}

func (r *PostgresRentalRepository) UpdateRentalExtensionStatus(ctx context.Context, extensionID uuid.UUID, status string, notes *string, slipURI, slipHash *string) error {
	var ownerDecidedAt pgtype.Timestamptz
	if status == "OWNER_APPROVED" || status == "OWNER_DECLINED" {
		ownerDecidedAt = pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true}
	}

	_, err := r.queries.UpdateRentalExtensionStatus(ctx, db.UpdateRentalExtensionStatusParams{
		ID:                 pgtype.UUID{Bytes: extensionID, Valid: true},
		Status:             status,
		OwnerDecidedAt:     ownerDecidedAt,
		OwnerDecisionNotes: notes,
		SlipStorageUri:     slipURI,
		SlipFileHash:       slipHash,
	})
	return err
}

func (r *PostgresRentalRepository) ExpandBookingHirePeriod(ctx context.Context, bookingID uuid.UUID, newEndDate time.Time, additionalDays int, additionalGrossMinor, additionalPrepayMinor int64, expectedVersion int) error {
	_, err := r.queries.ExpandBookingHirePeriod(ctx, db.ExpandBookingHirePeriodParams{
		ID:                     pgtype.UUID{Bytes: bookingID, Valid: true},
		EndDate:                pgtype.Date{Time: newEndDate, Valid: true},
		DurationDays:           int32(additionalDays),
		GrossRentalAmountMinor: additionalGrossMinor,
		PrepaymentAmountMinor:  additionalPrepayMinor,
		Version:                int32(expectedVersion),
	})
	return err
}

func (r *PostgresRentalRepository) CreateMechanicalBreakdown(ctx context.Context, breakdown *domain.MechanicalBreakdown) error {
	photosBytes, err := json.Marshal(breakdown.IncidentPhotos)
	if err != nil {
		return fmt.Errorf("failed to marshal incident photos: %w", err)
	}

	row, err := r.queries.CreateMechanicalBreakdown(ctx, db.CreateMechanicalBreakdownParams{
		BookingID:              pgtype.UUID{Bytes: breakdown.BookingID, Valid: true},
		ReportedByID:           pgtype.UUID{Bytes: breakdown.ReportedByID, Valid: true},
		FaultCategory:          breakdown.FaultCategory,
		FaultDescription:       breakdown.FaultDescription,
		DiagnosticTroubleCodes: breakdown.DiagnosticTroubleCodes,
		CurrentHourMeter:       floatToNumeric(breakdown.CurrentHourMeter),
		IncidentPhotosJson:     photosBytes,
		Status:                 breakdown.Status,
		SlaDeadline:            pgtype.Timestamptz{Time: breakdown.SLADeadline, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to insert mechanical breakdown: %w", err)
	}

	breakdown.ID = uuid.UUID(row.ID.Bytes)
	breakdown.CreatedAt = row.CreatedAt.Time
	breakdown.UpdatedAt = row.UpdatedAt.Time
	return nil
}

func (r *PostgresRentalRepository) GetMechanicalBreakdown(ctx context.Context, breakdownID uuid.UUID) (*domain.MechanicalBreakdown, error) {
	row, err := r.queries.GetMechanicalBreakdownByID(ctx, pgtype.UUID{Bytes: breakdownID, Valid: true})
	if err != nil {
		return nil, domain.ErrBreakdownNotFound
	}

	var photos []string
	_ = json.Unmarshal(row.IncidentPhotosJson, &photos)

	res := &domain.MechanicalBreakdown{
		ID:                     uuid.UUID(row.ID.Bytes),
		BookingID:              uuid.UUID(row.BookingID.Bytes),
		ReportedByID:           uuid.UUID(row.ReportedByID.Bytes),
		FaultCategory:          row.FaultCategory,
		FaultDescription:       row.FaultDescription,
		DiagnosticTroubleCodes: row.DiagnosticTroubleCodes,
		CurrentHourMeter:       numericToFloat(row.CurrentHourMeter),
		IncidentPhotos:         photos,
		Status:                 row.Status,
		SLADeadline:            row.SlaDeadline.Time,
		CreatedAt:              row.CreatedAt.Time,
		UpdatedAt:              row.UpdatedAt.Time,
	}

	if row.MechanicName != nil {
		res.MechanicName = *row.MechanicName
	}
	if row.MechanicPhone != nil {
		res.MechanicPhone = *row.MechanicPhone
	}
	if row.ResolutionNotes != nil {
		res.ResolutionNotes = *row.ResolutionNotes
	}
	if row.ResolvedAt.Valid {
		t := row.ResolvedAt.Time
		res.ResolvedAt = &t
	}

	return res, nil
}

func (r *PostgresRentalRepository) GetActiveMechanicalBreakdown(ctx context.Context, bookingID uuid.UUID) (*domain.MechanicalBreakdown, error) {
	row, err := r.queries.GetActiveMechanicalBreakdownByBookingID(ctx, pgtype.UUID{Bytes: bookingID, Valid: true})
	if err != nil {
		return nil, nil // Return nil if no active breakdown found
	}

	var photos []string
	_ = json.Unmarshal(row.IncidentPhotosJson, &photos)

	return &domain.MechanicalBreakdown{
		ID:                     uuid.UUID(row.ID.Bytes),
		BookingID:              uuid.UUID(row.BookingID.Bytes),
		ReportedByID:           uuid.UUID(row.ReportedByID.Bytes),
		FaultCategory:          row.FaultCategory,
		FaultDescription:       row.FaultDescription,
		DiagnosticTroubleCodes: row.DiagnosticTroubleCodes,
		CurrentHourMeter:       numericToFloat(row.CurrentHourMeter),
		IncidentPhotos:         photos,
		Status:                 row.Status,
		SLADeadline:            row.SlaDeadline.Time,
		CreatedAt:              row.CreatedAt.Time,
		UpdatedAt:              row.UpdatedAt.Time,
	}, nil
}

func (r *PostgresRentalRepository) ListMechanicalBreakdowns(ctx context.Context, bookingID uuid.UUID) ([]domain.MechanicalBreakdown, error) {
	rows, err := r.queries.ListMechanicalBreakdownsByBookingID(ctx, pgtype.UUID{Bytes: bookingID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list mechanical breakdowns: %w", err)
	}

	var res []domain.MechanicalBreakdown
	for _, row := range rows {
		var photos []string
		_ = json.Unmarshal(row.IncidentPhotosJson, &photos)

		b := domain.MechanicalBreakdown{
			ID:                     uuid.UUID(row.ID.Bytes),
			BookingID:              uuid.UUID(row.BookingID.Bytes),
			ReportedByID:           uuid.UUID(row.ReportedByID.Bytes),
			FaultCategory:          row.FaultCategory,
			FaultDescription:       row.FaultDescription,
			DiagnosticTroubleCodes: row.DiagnosticTroubleCodes,
			CurrentHourMeter:       numericToFloat(row.CurrentHourMeter),
			IncidentPhotos:         photos,
			Status:                 row.Status,
			SLADeadline:            row.SlaDeadline.Time,
			CreatedAt:              row.CreatedAt.Time,
			UpdatedAt:              row.UpdatedAt.Time,
		}
		if row.MechanicName != nil {
			b.MechanicName = *row.MechanicName
		}
		if row.MechanicPhone != nil {
			b.MechanicPhone = *row.MechanicPhone
		}
		if row.ResolutionNotes != nil {
			b.ResolutionNotes = *row.ResolutionNotes
		}
		if row.ResolvedAt.Valid {
			t := row.ResolvedAt.Time
			b.ResolvedAt = &t
		}
		res = append(res, b)
	}
	return res, nil
}

func (r *PostgresRentalRepository) ResolveMechanicalBreakdown(ctx context.Context, breakdownID uuid.UUID, mechanicName, mechanicPhone, resolutionNotes string) error {
	var mn, mp, rn *string
	if mechanicName != "" {
		mn = &mechanicName
	}
	if mechanicPhone != "" {
		mp = &mechanicPhone
	}
	if resolutionNotes != "" {
		rn = &resolutionNotes
	}

	_, err := r.queries.ResolveMechanicalBreakdown(ctx, db.ResolveMechanicalBreakdownParams{
		ID:              pgtype.UUID{Bytes: breakdownID, Valid: true},
		MechanicName:    mn,
		MechanicPhone:   mp,
		ResolutionNotes: rn,
	})
	return err
}

func (r *PostgresRentalRepository) GetFleetBaselineMeter(ctx context.Context, fleetID uuid.UUID) (float64, error) {
	row, err := r.queries.GetFleetBaselineMeter(ctx, pgtype.UUID{Bytes: fleetID, Valid: true})
	if err != nil {
		return 0.0, fmt.Errorf("failed to get fleet baseline meter: %w", err)
	}
	return numericToFloat(row.BaselineMeterHours), nil
}

func (r *PostgresRentalRepository) CheckForwardBookingConflict(ctx context.Context, fleetID uuid.UUID, excludeBookingID uuid.UUID, currentEndDate, newEndDate time.Time) (bool, error) {
	// Checks if any conflicting booking exists on the machine
	hirePeriod := pgtype.Range[pgtype.Date]{
		Lower:     pgtype.Date{Time: currentEndDate, Valid: true},
		Upper:     pgtype.Date{Time: newEndDate, Valid: true},
		LowerType: pgtype.Inclusive,
		UpperType: pgtype.Exclusive,
		Valid:     true,
	}

	exists, err := r.queries.CheckFleetAvailabilityConflict(ctx, db.CheckFleetAvailabilityConflictParams{
		FleetID:          pgtype.UUID{Bytes: fleetID, Valid: true},
		ExcludeBookingID: pgtype.UUID{Bytes: excludeBookingID, Valid: true},
		HirePeriod:       hirePeriod,
	})
	if err != nil {
		return false, err
	}
	return exists, nil
}


// Return Inspection (Partition 05-07)
func (r *PostgresRentalRepository) CreateReturnInspection(ctx context.Context, inspection *domain.ReturnInspection) error {
	photosBytes, err := json.Marshal(inspection.Photos)
	if err != nil {
		return fmt.Errorf("failed to marshal return inspection photos: %w", err)
	}

	var notes *string
	if inspection.Notes != "" {
		notes = &inspection.Notes
	}

	row, err := r.queries.CreateReturnInspection(ctx, db.CreateReturnInspectionParams{
		BookingID:           pgtype.UUID{Bytes: inspection.BookingID, Valid: true},
		InspectorID:         pgtype.UUID{Bytes: inspection.InspectorID, Valid: true},
		CheckOutMeterHours:  floatToNumeric(inspection.CheckOutMeterHours),
		FluidLevelsOk:       inspection.FluidLevelsOK,
		SafetyGlassIntact:   inspection.SafetyGlassIntact,
		AccelerometerSteady: inspection.AccelerometerSteady,
		Notes:               notes,
		PhotosJson:          photosBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to insert return inspection: %w", err)
	}

	inspection.ID = row.ID.Bytes
	inspection.CreatedAt = row.CreatedAt.Time
	return nil
}

func (r *PostgresRentalRepository) GetReturnInspection(ctx context.Context, bookingID uuid.UUID) (*domain.ReturnInspection, error) {
	row, err := r.queries.GetReturnInspectionByBookingID(ctx, pgtype.UUID{Bytes: bookingID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get return inspection: %w", err)
	}

	var photos []domain.InspectionPhoto
	if len(row.PhotosJson) > 0 {
		_ = json.Unmarshal(row.PhotosJson, &photos)
	}

	var notes string
	if row.Notes != nil {
		notes = *row.Notes
	}

	return &domain.ReturnInspection{
		ID:                  row.ID.Bytes,
		BookingID:           row.BookingID.Bytes,
		InspectorID:         row.InspectorID.Bytes,
		CheckOutMeterHours:  numericToFloat(row.CheckOutMeterHours),
		FluidLevelsOK:       row.FluidLevelsOk,
		SafetyGlassIntact:   row.SafetyGlassIntact,
		AccelerometerSteady: row.AccelerometerSteady,
		Notes:               notes,
		Photos:              photos,
		CreatedAt:           row.CreatedAt.Time,
	}, nil
}

// Damage Claims (Partition 05-07)
func (r *PostgresRentalRepository) CreateDamageClaimDocket(ctx context.Context, claim *domain.DamageClaimDocket) error {
	photosBytes, err := json.Marshal(claim.DamagePhotos)
	if err != nil {
		return fmt.Errorf("failed to marshal damage photos: %w", err)
	}

	var quoteURI *string
	if claim.RepairQuoteURI != "" {
		quoteURI = &claim.RepairQuoteURI
	}

	row, err := r.queries.CreateDamageClaimDocket(ctx, db.CreateDamageClaimDocketParams{
		BookingID:           pgtype.UUID{Bytes: claim.BookingID, Valid: true},
		ClaimantID:          pgtype.UUID{Bytes: claim.ClaimantID, Valid: true},
		ComponentCategory:   claim.ComponentCategory,
		IncidentNarrative:   claim.IncidentNarrative,
		RepairQuoteUri:      quoteURI,
		ClaimedAmountMinor:  claim.ClaimedAmountMinor,
		AcceptedAmountMinor: claim.AcceptedAmountMinor,
		Currency:            claim.Currency,
		DamagePhotosJson:    photosBytes,
		Status:              claim.Status,
	})
	if err != nil {
		return fmt.Errorf("failed to insert damage claim docket: %w", err)
	}

	claim.ID = row.ID.Bytes
	claim.CreatedAt = row.CreatedAt.Time
	claim.UpdatedAt = row.UpdatedAt.Time
	return nil
}

func (r *PostgresRentalRepository) GetDamageClaimDocket(ctx context.Context, claimID uuid.UUID) (*domain.DamageClaimDocket, error) {
	row, err := r.queries.GetDamageClaimDocketByID(ctx, pgtype.UUID{Bytes: claimID, Valid: true})
	if err != nil {
		return nil, domain.ErrDamageClaimNotFound
	}
	return toDomainDamageClaim(&row), nil
}

func (r *PostgresRentalRepository) GetDamageClaimDocketByBooking(ctx context.Context, bookingID uuid.UUID) (*domain.DamageClaimDocket, error) {
	row, err := r.queries.GetDamageClaimDocketByBookingID(ctx, pgtype.UUID{Bytes: bookingID, Valid: true})
	if err != nil {
		return nil, domain.ErrDamageClaimNotFound
	}
	return toDomainDamageClaim(&row), nil
}

func (r *PostgresRentalRepository) ListDamageClaimDockets(ctx context.Context, bookingID uuid.UUID) ([]domain.DamageClaimDocket, error) {
	rows, err := r.queries.ListDamageClaimDocketsByBookingID(ctx, pgtype.UUID{Bytes: bookingID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list damage claims: %w", err)
	}
	var result []domain.DamageClaimDocket
	for _, row := range rows {
		result = append(result, *toDomainDamageClaim(&row))
	}
	return result, nil
}

func toDomainDamageClaim(row *db.DamageClaimDocket) *domain.DamageClaimDocket {
	var photos []string
	if len(row.DamagePhotosJson) > 0 {
		_ = json.Unmarshal(row.DamagePhotosJson, &photos)
	}

	var quoteURI, disputeReason, arbitrationNotes string
	if row.RepairQuoteUri != nil {
		quoteURI = *row.RepairQuoteUri
	}
	if row.DisputeReason != nil {
		disputeReason = *row.DisputeReason
	}
	if row.ArbitrationNotes != nil {
		arbitrationNotes = *row.ArbitrationNotes
	}

	var contractorRespondedAt, arbitratedAt *time.Time
	if row.ContractorRespondedAt.Valid {
		contractorRespondedAt = &row.ContractorRespondedAt.Time
	}
	if row.ArbitratedAt.Valid {
		arbitratedAt = &row.ArbitratedAt.Time
	}

	var arbitratedBy *uuid.UUID
	if row.ArbitratedBy.Valid {
		u := uuid.UUID(row.ArbitratedBy.Bytes)
		arbitratedBy = &u
	}

	return &domain.DamageClaimDocket{
		ID:                    row.ID.Bytes,
		BookingID:             row.BookingID.Bytes,
		ClaimantID:            row.ClaimantID.Bytes,
		ComponentCategory:     row.ComponentCategory,
		IncidentNarrative:     row.IncidentNarrative,
		RepairQuoteURI:        quoteURI,
		ClaimedAmountMinor:    row.ClaimedAmountMinor,
		AcceptedAmountMinor:   row.AcceptedAmountMinor,
		Currency:              row.Currency,
		DamagePhotos:          photos,
		Status:                row.Status,
		DisputeReason:         disputeReason,
		ContractorRespondedAt: contractorRespondedAt,
		ArbitrationNotes:      arbitrationNotes,
		ArbitratedBy:          arbitratedBy,
		ArbitratedAt:          arbitratedAt,
		CreatedAt:             row.CreatedAt.Time,
		UpdatedAt:             row.UpdatedAt.Time,
	}
}

func (r *PostgresRentalRepository) RespondToDamageClaim(ctx context.Context, claimID uuid.UUID, status string, acceptedMinor int64, disputeReason string) error {
	var reason *string
	if disputeReason != "" {
		reason = &disputeReason
	}

	_, err := r.queries.RespondToDamageClaim(ctx, db.RespondToDamageClaimParams{
		ID:                  pgtype.UUID{Bytes: claimID, Valid: true},
		Status:              status,
		AcceptedAmountMinor: acceptedMinor,
		DisputeReason:       reason,
	})
	return err
}

func (r *PostgresRentalRepository) ArbitrateDamageClaim(ctx context.Context, claimID uuid.UUID, status string, acceptedMinor int64, arbitrationNotes string, arbitratedBy uuid.UUID) error {
	var notes *string
	if arbitrationNotes != "" {
		notes = &arbitrationNotes
	}

	_, err := r.queries.ArbitrateDamageClaim(ctx, db.ArbitrateDamageClaimParams{
		ID:                  pgtype.UUID{Bytes: claimID, Valid: true},
		Status:              status,
		AcceptedAmountMinor: acceptedMinor,
		ArbitrationNotes:    notes,
		ArbitratedBy:        pgtype.UUID{Bytes: arbitratedBy, Valid: true},
	})
	return err
}

func (r *PostgresRentalRepository) WaiveDamageClaim(ctx context.Context, claimID uuid.UUID) error {
	_, err := r.queries.WaiveDamageClaim(ctx, pgtype.UUID{Bytes: claimID, Valid: true})
	return err
}

// Rental Settlement (Partition 05-07)
func (r *PostgresRentalRepository) UpsertRentalSettlement(ctx context.Context, settlement *domain.RentalSettlement) error {
	row, err := r.queries.UpsertRentalSettlement(ctx, db.UpsertRentalSettlementParams{
		BookingID:                  pgtype.UUID{Bytes: settlement.BookingID, Valid: true},
		BaseMachineFeeMinor:        settlement.BaseMachineFeeMinor,
		BundledAttachmentsFeeMinor: settlement.BundledAttachmentsFeeMinor,
		ContractualAllowedHours:    floatToNumeric(settlement.ContractualAllowedHours),
		ActualOperatingHours:       floatToNumeric(settlement.ActualOperatingHours),
		CreditedBreakdownHours:     floatToNumeric(settlement.CreditedBreakdownHours),
		BillableExcessHours:        floatToNumeric(settlement.BillableExcessHours),
		ExcessHourlyRateMinor:      settlement.ExcessHourlyRateMinor,
		ExcessHoursSurchargeMinor:  settlement.ExcessHoursSurchargeMinor,
		AcceptedDamageMinor:        settlement.AcceptedDamageMinor,
		GrossLiabilitiesMinor:      settlement.GrossLiabilitiesMinor,
		PrepaymentSettledMinor:     settlement.PrepaymentSettledMinor,
		NetBalanceMinor:            settlement.NetBalanceMinor,
		Currency:                   settlement.Currency,
		Status:                     settlement.Status,
	})
	if err != nil {
		return fmt.Errorf("failed to upsert rental settlement: %w", err)
	}

	settlement.ID = row.ID.Bytes
	settlement.CreatedAt = row.CreatedAt.Time
	settlement.UpdatedAt = row.UpdatedAt.Time
	return nil
}

func (r *PostgresRentalRepository) GetRentalSettlement(ctx context.Context, bookingID uuid.UUID) (*domain.RentalSettlement, error) {
	row, err := r.queries.GetRentalSettlementByBookingID(ctx, pgtype.UUID{Bytes: bookingID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get rental settlement: %w", err)
	}

	var slipURI, slipHash string
	if row.FinalSlipStorageUri != nil {
		slipURI = *row.FinalSlipStorageUri
	}
	if row.FinalSlipFileHash != nil {
		slipHash = *row.FinalSlipFileHash
	}

	var settledAt *time.Time
	if row.SettledAt.Valid {
		settledAt = &row.SettledAt.Time
	}

	return &domain.RentalSettlement{
		ID:                         row.ID.Bytes,
		BookingID:                  row.BookingID.Bytes,
		BaseMachineFeeMinor:        row.BaseMachineFeeMinor,
		BundledAttachmentsFeeMinor: row.BundledAttachmentsFeeMinor,
		ContractualAllowedHours:    numericToFloat(row.ContractualAllowedHours),
		ActualOperatingHours:       numericToFloat(row.ActualOperatingHours),
		CreditedBreakdownHours:     numericToFloat(row.CreditedBreakdownHours),
		BillableExcessHours:        numericToFloat(row.BillableExcessHours),
		ExcessHourlyRateMinor:      row.ExcessHourlyRateMinor,
		ExcessHoursSurchargeMinor:  row.ExcessHoursSurchargeMinor,
		AcceptedDamageMinor:        row.AcceptedDamageMinor,
		GrossLiabilitiesMinor:      row.GrossLiabilitiesMinor,
		PrepaymentSettledMinor:     row.PrepaymentSettledMinor,
		NetBalanceMinor:            row.NetBalanceMinor,
		Currency:                   row.Currency,
		Status:                     row.Status,
		FinalSlipStorageURI:        slipURI,
		FinalSlipFileHash:          slipHash,
		SettledAt:                  settledAt,
		CreatedAt:                  row.CreatedAt.Time,
		UpdatedAt:                  row.UpdatedAt.Time,
	}, nil
}

func (r *PostgresRentalRepository) UpdateRentalSettlementPaymentEvidence(ctx context.Context, bookingID uuid.UUID, slipURI, slipHash string) error {
	_, err := r.queries.UpdateRentalSettlementFinalPaymentEvidence(ctx, db.UpdateRentalSettlementFinalPaymentEvidenceParams{
		BookingID:           pgtype.UUID{Bytes: bookingID, Valid: true},
		FinalSlipStorageUri: &slipURI,
		FinalSlipFileHash:   &slipHash,
	})
	return err
}

func (r *PostgresRentalRepository) ConfirmRentalSettlement(ctx context.Context, bookingID uuid.UUID) error {
	_, err := r.queries.ConfirmRentalSettlement(ctx, pgtype.UUID{Bytes: bookingID, Valid: true})
	return err
}

// Booking Reviews (Partition 05-07)
func (r *PostgresRentalRepository) CreateBookingReview(ctx context.Context, review *domain.BookingReview) error {
	dimBytes, err := json.Marshal(review.DimensionRatings)
	if err != nil {
		return fmt.Errorf("failed to marshal dimension ratings: %w", err)
	}

	var feedback *string
	if review.Feedback != "" {
		feedback = &review.Feedback
	}

	var publishedAt pgtype.Timestamptz
	if review.PublishedAt != nil {
		publishedAt = pgtype.Timestamptz{Time: *review.PublishedAt, Valid: true}
	}

	row, err := r.queries.CreateBookingReview(ctx, db.CreateBookingReviewParams{
		BookingID:            pgtype.UUID{Bytes: review.BookingID, Valid: true},
		AuthorID:             pgtype.UUID{Bytes: review.AuthorID, Valid: true},
		TargetID:             pgtype.UUID{Bytes: review.TargetID, Valid: true},
		Role:                 review.Role,
		Rating:               floatToNumeric(review.Rating),
		Feedback:             feedback,
		DimensionRatingsJson: dimBytes,
		IsPublished:          review.IsPublished,
		PublishedAt:          publishedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to insert booking review: %w", err)
	}

	review.ID = row.ID.Bytes
	review.CreatedAt = row.CreatedAt.Time
	return nil
}

func (r *PostgresRentalRepository) GetBookingReviews(ctx context.Context, bookingID uuid.UUID) ([]domain.BookingReview, error) {
	rows, err := r.queries.GetBookingReviewsByBookingID(ctx, pgtype.UUID{Bytes: bookingID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get booking reviews: %w", err)
	}

	var result []domain.BookingReview
	for _, row := range rows {
		result = append(result, *toDomainReview(&row))
	}
	return result, nil
}

func (r *PostgresRentalRepository) PublishBookingReviews(ctx context.Context, bookingID uuid.UUID) error {
	return r.queries.PublishBookingReviews(ctx, pgtype.UUID{Bytes: bookingID, Valid: true})
}

func (r *PostgresRentalRepository) ListPublishedReviewsForTarget(ctx context.Context, targetID uuid.UUID) ([]domain.BookingReview, error) {
	rows, err := r.queries.ListPublishedReviewsForTarget(ctx, pgtype.UUID{Bytes: targetID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list published reviews: %w", err)
	}

	var result []domain.BookingReview
	for _, row := range rows {
		result = append(result, *toDomainReview(&row))
	}
	return result, nil
}

func toDomainReview(row *db.BookingReview) *domain.BookingReview {
	var dimensions map[string]float64
	if len(row.DimensionRatingsJson) > 0 {
		_ = json.Unmarshal(row.DimensionRatingsJson, &dimensions)
	}

	var feedback string
	if row.Feedback != nil {
		feedback = *row.Feedback
	}

	var publishedAt *time.Time
	if row.PublishedAt.Valid {
		publishedAt = &row.PublishedAt.Time
	}

	return &domain.BookingReview{
		ID:               row.ID.Bytes,
		BookingID:        row.BookingID.Bytes,
		AuthorID:         row.AuthorID.Bytes,
		TargetID:         row.TargetID.Bytes,
		Role:             row.Role,
		Rating:           numericToFloat(row.Rating),
		Feedback:         feedback,
		DimensionRatings: dimensions,
		IsPublished:      row.IsPublished,
		PublishedAt:      publishedAt,
		CreatedAt:        row.CreatedAt.Time,
	}
}

func (r *PostgresRentalRepository) UpdateBookingCompletionSummary(
	ctx context.Context,
	bookingID uuid.UUID,
	status string,
	checkOutMeter float64,
	damageClaimSubmitted, damageClaimWaived bool,
	acceptedDamageMinor, excessSurchargeMinor, finalNetBalanceMinor int64,
	expectedVersion int,
) error {
	_, err := r.queries.UpdateBookingCompletionSummary(ctx, db.UpdateBookingCompletionSummaryParams{
		ID:                        pgtype.UUID{Bytes: bookingID, Valid: true},
		Status:                    status,
		CheckOutMeterHours:        floatToNumeric(checkOutMeter),
		DamageClaimSubmitted:      damageClaimSubmitted,
		DamageClaimWaived:         damageClaimWaived,
		AcceptedDamageMinor:       acceptedDamageMinor,
		ExcessHoursSurchargeMinor: excessSurchargeMinor,
		FinalNetBalanceMinor:      finalNetBalanceMinor,
		Version:                   int32(expectedVersion),
	})
	if err != nil {
		return fmt.Errorf("failed to update booking completion summary: %w", err)
	}
	return nil
}
