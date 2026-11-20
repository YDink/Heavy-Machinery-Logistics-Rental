package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/database/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostgresBookingRepository struct {
	queries *db.Queries
}

func NewPostgresBookingRepository(queries *db.Queries) *PostgresBookingRepository {
	return &PostgresBookingRepository{queries: queries}
}

func toDomainBooking(b db.Booking) *domain.Booking {
	var parentID *uuid.UUID
	if b.ParentBookingID.Valid {
		id := uuid.UUID(b.ParentBookingID.Bytes)
		parentID = &id
	}
	res := &domain.Booking{
		ID:                         uuid.UUID(b.ID.Bytes),
		BookingReference:           b.BookingReference,
		RenterID:                   uuid.UUID(b.RenterID.Bytes),
		OwnerID:                    uuid.UUID(b.OwnerID.Bytes),
		FleetID:                    uuid.UUID(b.FleetID.Bytes),
		ParentBookingID:            parentID,
		BookingKind:                b.BookingKind,
		StartDate:                  b.StartDate.Time,
		EndDate:                    b.EndDate.Time,
		DurationDays:               int(b.DurationDays),
		DailyRateAppliedMinor:      b.DailyRateAppliedMinor,
		BaseHireFeeMinor:           b.BaseHireFeeMinor,
		AttachmentsFeeMinor:        b.AttachmentsFeeMinor,
		GrossRentalAmountMinor:     b.GrossRentalAmountMinor,
		PrepaymentAmountMinor:      b.PrepaymentAmountMinor,
		PlatformCommissionRateBps:  int(b.PlatformCommissionRateBps),
		PlatformCommissionMinor:    b.PlatformCommissionMinor,
		AbnormalLoadEscortRequired: b.AbnormalLoadEscortRequired,
		Currency:                   b.Currency,
		SiteLocationAddress:    b.SiteLocationAddress,
		Status:                 domain.BookingStatus(b.Status),
		RejectionReasonCode:    b.RejectionReasonCode,
		RejectionNotes:         b.RejectionNotes,
		CancellationReason:     b.CancellationReason,
		PenaltyAmountMinor:     b.PenaltyAmountMinor,
		RefundAmountMinor:      b.RefundAmountMinor,
		HandoverPIN:            b.HandoverPin,
		Version:                int(b.Version),
		CreatedAt:              b.CreatedAt.Time,
		UpdatedAt:              b.UpdatedAt.Time,
	}

	if b.AssignedOperatorID.Valid {
		opID := uuid.UUID(b.AssignedOperatorID.Bytes)
		res.AssignedOperatorID = &opID
	}
	if b.CancelledBy.Valid {
		id := uuid.UUID(b.CancelledBy.Bytes)
		res.CancelledBy = &id
	}
	if b.ExpiresAt.Valid {
		t := b.ExpiresAt.Time
		res.ExpiresAt = &t
	}
	if b.ConfirmedAt.Valid {
		t := b.ConfirmedAt.Time
		res.ConfirmedAt = &t
	}
	if b.ActivatedAt.Valid {
		t := b.ActivatedAt.Time
		res.ActivatedAt = &t
	}
	if b.CompletedAt.Valid {
		t := b.CompletedAt.Time
		res.CompletedAt = &t
	}
	if b.CancelledAt.Valid {
		t := b.CancelledAt.Time
		res.CancelledAt = &t
	}
	if b.CheckInMeterHours.Valid {
		val, _ := b.CheckInMeterHours.Float64Value()
		res.CheckInMeterHours = &val.Float64
	}
	if b.CheckOutMeterHours.Valid {
		val, _ := b.CheckOutMeterHours.Float64Value()
		res.CheckOutMeterHours = &val.Float64
	}
	return res
}

func (r *PostgresBookingRepository) CreateBooking(ctx context.Context, b *domain.Booking) error {
	hirePeriod := pgtype.Range[pgtype.Date]{
		Lower:     pgtype.Date{Time: b.StartDate, Valid: true},
		Upper:     pgtype.Date{Time: b.EndDate, Valid: true},
		LowerType: pgtype.Inclusive,
		UpperType: pgtype.Exclusive,
		Valid:     true,
	}

	var expiresAt pgtype.Timestamptz
	if b.ExpiresAt != nil {
		expiresAt = pgtype.Timestamptz{Time: *b.ExpiresAt, Valid: true}
	}

	row, err := r.queries.CreateBooking(ctx, db.CreateBookingParams{
		BookingReference:       b.BookingReference,
		RenterID:               pgtype.UUID{Bytes: b.RenterID, Valid: true},
		OwnerID:                pgtype.UUID{Bytes: b.OwnerID, Valid: true},
		FleetID:                pgtype.UUID{Bytes: b.FleetID, Valid: true},
		StartDate:              pgtype.Date{Time: b.StartDate, Valid: true},
		EndDate:                pgtype.Date{Time: b.EndDate, Valid: true},
		HirePeriod:             hirePeriod,
		DurationDays:           int32(b.DurationDays),
		DailyRateAppliedMinor:  b.DailyRateAppliedMinor,
		BaseHireFeeMinor:       b.BaseHireFeeMinor,
		AttachmentsFeeMinor:    b.AttachmentsFeeMinor,
		GrossRentalAmountMinor: b.GrossRentalAmountMinor,
		PrepaymentAmountMinor:  b.PrepaymentAmountMinor,
		Currency:               b.Currency,
		SiteLocationAddress:    b.SiteLocationAddress,
		Status:                 string(b.Status),
		ExpiresAt:              expiresAt,
		HandoverPin:            b.HandoverPIN,
	})
	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}

	b.ID = uuid.UUID(row.ID.Bytes)
	b.CreatedAt = row.CreatedAt.Time
	b.UpdatedAt = row.UpdatedAt.Time
	b.Version = int(row.Version)
	return nil
}

func (r *PostgresBookingRepository) GetBookingByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	row, err := r.queries.GetBookingByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBookingNotFound
		}
		return nil, fmt.Errorf("failed to get booking by ID: %w", err)
	}
	b := toDomainBooking(row)
	atts, err := r.ListBookingAttachments(ctx, id)
	if err == nil {
		b.Attachments = atts
	}
	return b, nil
}

func (r *PostgresBookingRepository) GetBookingByReference(ctx context.Context, reference string) (*domain.Booking, error) {
	row, err := r.queries.GetBookingByReference(ctx, reference)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBookingNotFound
		}
		return nil, fmt.Errorf("failed to get booking by reference: %w", err)
	}
	b := toDomainBooking(row)
	atts, err := r.ListBookingAttachments(ctx, b.ID)
	if err == nil {
		b.Attachments = atts
	}
	return b, nil
}

func (r *PostgresBookingRepository) ListBookingsByRenter(ctx context.Context, renterID uuid.UUID, status string, limit int) ([]*domain.Booking, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.queries.ListBookingsByRenter(ctx, db.ListBookingsByRenterParams{
		RenterID: pgtype.UUID{Bytes: renterID, Valid: true},
		Column2:  status,
		Limit:    int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list bookings by renter: %w", err)
	}

	res := make([]*domain.Booking, len(rows))
	for i, row := range rows {
		res[i] = toDomainBooking(row)
	}
	return res, nil
}

func (r *PostgresBookingRepository) ListBookingsByOwner(ctx context.Context, ownerID uuid.UUID, status string, limit int) ([]*domain.Booking, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.queries.ListBookingsByOwner(ctx, db.ListBookingsByOwnerParams{
		OwnerID: pgtype.UUID{Bytes: ownerID, Valid: true},
		Column2: status,
		Limit:   int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list bookings by owner: %w", err)
	}

	res := make([]*domain.Booking, len(rows))
	for i, row := range rows {
		res[i] = toDomainBooking(row)
	}
	return res, nil
}

func (r *PostgresBookingRepository) UpdateBookingStatusOCC(ctx context.Context, b *domain.Booking, expectedVersion int) error {
	var cancelledBy pgtype.UUID
	if b.CancelledBy != nil {
		cancelledBy = pgtype.UUID{Bytes: *b.CancelledBy, Valid: true}
	}
	var expiresAt pgtype.Timestamptz
	if b.ExpiresAt != nil {
		expiresAt = pgtype.Timestamptz{Time: *b.ExpiresAt, Valid: true}
	}

	row, err := r.queries.UpdateBookingStatusOCC(ctx, db.UpdateBookingStatusOCCParams{
		ID:                  pgtype.UUID{Bytes: b.ID, Valid: true},
		Status:              string(b.Status),
		RejectionReasonCode: b.RejectionReasonCode,
		RejectionNotes:      b.RejectionNotes,
		CancellationReason:  b.CancellationReason,
		CancelledBy:         cancelledBy,
		PenaltyAmountMinor:  b.PenaltyAmountMinor,
		RefundAmountMinor:   b.RefundAmountMinor,
		ExpiresAt:           expiresAt,
		Version:             int32(expectedVersion),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrConcurrencyConflict
		}
		return fmt.Errorf("failed to update booking status: %w", err)
	}

	b.Version = int(row.Version)
	b.UpdatedAt = row.UpdatedAt.Time
	return nil
}

func (r *PostgresBookingRepository) UpdateBookingHandover(ctx context.Context, bookingID uuid.UUID, checkInMeter float64, expectedVersion int) error {
	var meterNumeric pgtype.Numeric
	_ = meterNumeric.Scan(fmt.Sprintf("%.2f", checkInMeter))

	_, err := r.queries.UpdateBookingHandover(ctx, db.UpdateBookingHandoverParams{
		ID:                 pgtype.UUID{Bytes: bookingID, Valid: true},
		CheckInMeterHours:  meterNumeric,
		Version:            int32(expectedVersion),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrConcurrencyConflict
		}
		return fmt.Errorf("failed to update booking handover: %w", err)
	}
	return nil
}

func (r *PostgresBookingRepository) UpdateBookingReturn(ctx context.Context, bookingID uuid.UUID, status domain.BookingStatus, checkOutMeter float64, expectedVersion int) error {
	var meterNumeric pgtype.Numeric
	_ = meterNumeric.Scan(fmt.Sprintf("%.2f", checkOutMeter))

	_, err := r.queries.UpdateBookingReturn(ctx, db.UpdateBookingReturnParams{
		ID:                 pgtype.UUID{Bytes: bookingID, Valid: true},
		Status:             string(status),
		CheckOutMeterHours: meterNumeric,
		Version:            int32(expectedVersion),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrConcurrencyConflict
		}
		return fmt.Errorf("failed to update booking return: %w", err)
	}
	return nil
}

func (r *PostgresBookingRepository) AddBookingAttachment(ctx context.Context, bookingID, attachmentID uuid.UUID, dailyRateMinor int64) error {
	_, err := r.queries.AddBookingAttachment(ctx, db.AddBookingAttachmentParams{
		BookingID:             pgtype.UUID{Bytes: bookingID, Valid: true},
		AttachmentID:          pgtype.UUID{Bytes: attachmentID, Valid: true},
		BundledDailyRateMinor: dailyRateMinor,
	})
	if err != nil {
		return fmt.Errorf("failed to add booking attachment: %w", err)
	}
	return nil
}

func (r *PostgresBookingRepository) ListBookingAttachments(ctx context.Context, bookingID uuid.UUID) ([]domain.BookingAttachment, error) {
	rows, err := r.queries.ListBookingAttachments(ctx, pgtype.UUID{Bytes: bookingID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list booking attachments: %w", err)
	}

	res := make([]domain.BookingAttachment, len(rows))
	for i, row := range rows {
		res[i] = domain.BookingAttachment{
			ID:                    uuid.UUID(row.ID.Bytes),
			BookingID:             uuid.UUID(row.BookingID.Bytes),
			AttachmentID:          uuid.UUID(row.AttachmentID.Bytes),
			AttachmentType:        row.AttachmentType,
			Make:                  row.Make,
			Model:                 row.Model,
			SerialNumber:          row.SerialNumber,
			BundledDailyRateMinor: row.BundledDailyRateMinor,
			CreatedAt:             row.CreatedAt.Time,
		}
	}
	return res, nil
}

func (r *PostgresBookingRepository) CreateBankSlip(ctx context.Context, slip *domain.BankSlip) error {
	row, err := r.queries.CreateBookingBankSlip(ctx, db.CreateBookingBankSlipParams{
		BookingID:            pgtype.UUID{Bytes: slip.BookingID, Valid: true},
		UploaderID:           pgtype.UUID{Bytes: slip.UploaderID, Valid: true},
		StorageUri:           slip.StorageURI,
		FileHashSha256:       slip.FileHashSHA256,
		TransactionReference: slip.TransactionReference,
		AmountMinor:          slip.AmountMinor,
		Currency:             slip.Currency,
		Status:               string(slip.Status),
	})
	if err != nil {
		return fmt.Errorf("failed to create booking bank slip: %w", err)
	}
	slip.ID = uuid.UUID(row.ID.Bytes)
	slip.CreatedAt = row.CreatedAt.Time
	return nil
}

func (r *PostgresBookingRepository) GetLatestBankSlip(ctx context.Context, bookingID uuid.UUID) (*domain.BankSlip, error) {
	row, err := r.queries.GetLatestBankSlipForBooking(ctx, pgtype.UUID{Bytes: bookingID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSlipNotFound
		}
		return nil, fmt.Errorf("failed to get latest bank slip: %w", err)
	}

	slip := &domain.BankSlip{
		ID:                   uuid.UUID(row.ID.Bytes),
		BookingID:            uuid.UUID(row.BookingID.Bytes),
		UploaderID:           uuid.UUID(row.UploaderID.Bytes),
		StorageURI:           row.StorageUri,
		FileHashSHA256:       row.FileHashSha256,
		TransactionReference: row.TransactionReference,
		AmountMinor:          row.AmountMinor,
		Currency:             row.Currency,
		Status:               domain.BankSlipStatus(row.Status),
		RejectionReasonCode:  row.RejectionReasonCode,
		RejectionNotes:       row.RejectionNotes,
		CreatedAt:            row.CreatedAt.Time,
	}
	if row.AuditedAt.Valid {
		t := row.AuditedAt.Time
		slip.AuditedAt = &t
	}
	return slip, nil
}

func (r *PostgresBookingRepository) UpdateBankSlipStatus(ctx context.Context, slipID uuid.UUID, status domain.BankSlipStatus, reasonCode, notes *string) error {
	_, err := r.queries.UpdateBankSlipStatus(ctx, db.UpdateBankSlipStatusParams{
		ID:                  pgtype.UUID{Bytes: slipID, Valid: true},
		Status:              string(status),
		RejectionReasonCode: reasonCode,
		RejectionNotes:      notes,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrSlipNotFound
		}
		return fmt.Errorf("failed to update bank slip status: %w", err)
	}
	return nil
}

func (r *PostgresBookingRepository) CreateAuditLog(ctx context.Context, l *domain.AuditLog) error {
	var actorID pgtype.UUID
	if l.ActorID != nil {
		actorID = pgtype.UUID{Bytes: *l.ActorID, Valid: true}
	}
	payloadBytes, err := json.Marshal(l.Payload)
	if err != nil {
		payloadBytes = []byte("{}")
	}

	row, err := r.queries.CreateBookingAuditLog(ctx, db.CreateBookingAuditLogParams{
		BookingID:  pgtype.UUID{Bytes: l.BookingID, Valid: true},
		ActorID:    actorID,
		Action:     l.Action,
		FromStatus: string(l.FromStatus),
		ToStatus:   string(l.ToStatus),
		Reason:     &l.Reason,
		Payload:    payloadBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}
	l.ID = uuid.UUID(row.ID.Bytes)
	l.CreatedAt = row.CreatedAt.Time
	return nil
}

func (r *PostgresBookingRepository) ListAuditLogs(ctx context.Context, bookingID uuid.UUID) ([]*domain.AuditLog, error) {
	rows, err := r.queries.ListBookingAuditLogs(ctx, pgtype.UUID{Bytes: bookingID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list audit logs: %w", err)
	}

	res := make([]*domain.AuditLog, len(rows))
	for i, row := range rows {
		var actor *uuid.UUID
		if row.ActorID.Valid {
			id := uuid.UUID(row.ActorID.Bytes)
			actor = &id
		}
		var payload map[string]any
		_ = json.Unmarshal(row.Payload, &payload)

		reasonStr := ""
		if row.Reason != nil {
			reasonStr = *row.Reason
		}

		res[i] = &domain.AuditLog{
			ID:         uuid.UUID(row.ID.Bytes),
			BookingID:  uuid.UUID(row.BookingID.Bytes),
			ActorID:    actor,
			Action:     row.Action,
			FromStatus: domain.BookingStatus(row.FromStatus),
			ToStatus:   domain.BookingStatus(row.ToStatus),
			Reason:     reasonStr,
			Payload:    payload,
			CreatedAt:  row.CreatedAt.Time,
		}
	}
	return res, nil
}

var _ domain.BookingRepository = (*PostgresBookingRepository)(nil)
