package application

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	machineryDomain "github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
	"github.com/google/uuid"
)

type FleetProvider interface {
	GetFleetAsset(ctx context.Context, id uuid.UUID) (*machineryDomain.FleetAsset, error)
}

type AttachmentProvider interface {
	GetBundledAttachment(ctx context.Context, fleetID, attachmentID uuid.UUID) (*machineryDomain.FleetBundledAttachment, error)
}

type BookingService struct {
	repo               domain.BookingRepository
	leaseManager       domain.HoldLeaseManager
	fleetProvider      FleetProvider
	attachmentProvider AttachmentProvider
}

func NewBookingService(
	repo domain.BookingRepository,
	leaseManager domain.HoldLeaseManager,
	fleetProvider FleetProvider,
	attachmentProvider AttachmentProvider,
) *BookingService {
	return &BookingService{
		repo:               repo,
		leaseManager:       leaseManager,
		fleetProvider:      fleetProvider,
		attachmentProvider: attachmentProvider,
	}
}

type CreateBookingInput struct {
	RenterID            uuid.UUID
	FleetID             uuid.UUID
	StartDate           time.Time
	EndDate             time.Time
	SiteLocationAddress string
	AttachmentIDs       []uuid.UUID
}

func (s *BookingService) CreateBooking(ctx context.Context, input CreateBookingInput) (*domain.Booking, error) {
	// 1. Fetch Fleet
	fleet, err := s.fleetProvider.GetFleetAsset(ctx, input.FleetID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve fleet asset: %w", err)
	}

	// Invariant: Fleet owner cannot book their own machinery
	if fleet.OwnerID == input.RenterID {
		return nil, domain.ErrSelfBookingProhibited
	}

	// 2. Acquire Valkey 60-Minute Distributed Lease (BK-05)
	leaseTTL := 60 * time.Minute
	acquired, err := s.leaseManager.AcquireHoldLease(ctx, input.FleetID, input.AttachmentIDs, input.RenterID, leaseTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire hold lease: %w", err)
	}
	if !acquired {
		return nil, domain.ErrConcurrentReservationCollision
	}

	// 3. Resolve Bundled Attachments Rates
	var totalAttachmentDailyRate int64
	var bundledAttachments []domain.BookingAttachment
	for _, attID := range input.AttachmentIDs {
		var dailyRate int64
		if s.attachmentProvider != nil {
			bundled, err := s.attachmentProvider.GetBundledAttachment(ctx, input.FleetID, attID)
			if err == nil && bundled != nil {
				dailyRate = bundled.BundledDailyRateMinor
			}
		}
		totalAttachmentDailyRate += dailyRate
		bundledAttachments = append(bundledAttachments, domain.BookingAttachment{
			ID:                    uuid.New(),
			AttachmentID:          attID,
			BundledDailyRateMinor: dailyRate,
			CreatedAt:             time.Now(),
		})
	}

	// 4. Calculate Quote Breakdown (BK-04)
	quote, err := domain.CalculateQuote(
		input.StartDate,
		input.EndDate,
		fleet.DailyRateMinor,
		totalAttachmentDailyRate,
		fleet.MinRentalDays,
		"ETB",
	)
	if err != nil {
		_ = s.leaseManager.ReleaseHoldLease(ctx, input.FleetID, input.AttachmentIDs)
		return nil, err
	}

	// 5. Generate Reference and Handover PIN
	refNum, _ := rand.Int(rand.Reader, big.NewInt(9000))
	refStr := fmt.Sprintf("BK-%04d", refNum.Int64()+1000)

	pinNum, _ := rand.Int(rand.Reader, big.NewInt(900000))
	pinStr := fmt.Sprintf("%06d", pinNum.Int64()+100000)

	now := time.Now()
	expiresAt := now.Add(leaseTTL)

	booking := &domain.Booking{
		ID:                     uuid.New(),
		BookingReference:       refStr,
		RenterID:               input.RenterID,
		OwnerID:                fleet.OwnerID,
		FleetID:                input.FleetID,
		StartDate:              input.StartDate,
		EndDate:                input.EndDate,
		DurationDays:           quote.DurationDays,
		DailyRateAppliedMinor:  quote.DailyRateMinor,
		BaseHireFeeMinor:       quote.BaseHireFeeMinor,
		AttachmentsFeeMinor:    quote.AttachmentsFeeMinor,
		GrossRentalAmountMinor: quote.GrossRentalAmountMinor,
		PrepaymentAmountMinor:  quote.PrepaymentAmountMinor,
		Currency:               quote.Currency,
		SiteLocationAddress:    input.SiteLocationAddress,
		Status:                 domain.StatusPendingPayment,
		ExpiresAt:              &expiresAt,
		HandoverPIN:            &pinStr,
		Attachments:            bundledAttachments,
		Version:                1,
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	if err := s.repo.CreateBooking(ctx, booking); err != nil {
		_ = s.leaseManager.ReleaseHoldLease(ctx, input.FleetID, input.AttachmentIDs)
		return nil, fmt.Errorf("failed to persist booking: %w", err)
	}

	for _, att := range bundledAttachments {
		_ = s.repo.AddBookingAttachment(ctx, booking.ID, att.AttachmentID, att.BundledDailyRateMinor)
	}

	_ = s.repo.CreateAuditLog(ctx, &domain.AuditLog{
		BookingID:  booking.ID,
		ActorID:    &input.RenterID,
		Action:     "CREATE_BOOKING",
		FromStatus: domain.StatusDraft,
		ToStatus:   domain.StatusPendingPayment,
		Reason:     "Initial booking requisition and 60-min Valkey lease hold created",
		Payload: map[string]any{
			"grossRentalAmountMinor": quote.GrossRentalAmountMinor,
			"prepaymentAmountMinor":  quote.PrepaymentAmountMinor,
			"durationDays":           quote.DurationDays,
		},
		CreatedAt: now,
	})

	return booking, nil
}

type SubmitPaymentSlipInput struct {
	BookingID            uuid.UUID
	UploaderID           uuid.UUID
	StorageURI           string
	FileHashSHA256       string
	TransactionReference string
	AmountMinor          int64
	Currency             string
}

func (s *BookingService) SubmitPaymentSlip(ctx context.Context, input SubmitPaymentSlipInput) (*domain.BankSlip, error) {
	b, err := s.repo.GetBookingByID(ctx, input.BookingID)
	if err != nil {
		return nil, err
	}

	if b.RenterID != input.UploaderID {
		return nil, domain.ErrUnauthorizedAction
	}

	if err := domain.ValidateTransition(b.Status, domain.StatusPaymentEvidenceUploaded); err != nil {
		return nil, err
	}

	now := time.Now()
	slip := &domain.BankSlip{
		ID:                   uuid.New(),
		BookingID:            b.ID,
		UploaderID:           input.UploaderID,
		StorageURI:           input.StorageURI,
		FileHashSHA256:       input.FileHashSHA256,
		TransactionReference: &input.TransactionReference,
		AmountMinor:          input.AmountMinor,
		Currency:             input.Currency,
		Status:               domain.SlipStatusPendingAudit,
		CreatedAt:            now,
	}

	if err := s.repo.CreateBankSlip(ctx, slip); err != nil {
		return nil, fmt.Errorf("failed to persist bank slip: %w", err)
	}

	// Extend Valkey lease to 24 hours for owner audit
	var attIDs []uuid.UUID
	for _, att := range b.Attachments {
		attIDs = append(attIDs, att.AttachmentID)
	}
	_ = s.leaseManager.ExtendHoldLease(ctx, b.FleetID, attIDs, 24*time.Hour)

	auditExpiresAt := now.Add(24 * time.Hour)
	b.Status = domain.StatusPaymentEvidenceUploaded
	b.ExpiresAt = &auditExpiresAt

	if err := s.repo.UpdateBookingStatusOCC(ctx, b, b.Version); err != nil {
		return nil, err
	}

	_ = s.repo.CreateAuditLog(ctx, &domain.AuditLog{
		BookingID:  b.ID,
		ActorID:    &input.UploaderID,
		Action:     "UPLOAD_PAYMENT_SLIP",
		FromStatus: domain.StatusPendingPayment,
		ToStatus:   domain.StatusPaymentEvidenceUploaded,
		Reason:     "Bank transfer slip uploaded; lease extended to 24 hours",
		Payload: map[string]any{
			"slipId":         slip.ID,
			"amountMinor":    input.AmountMinor,
			"fileHashSHA256": input.FileHashSHA256,
		},
		CreatedAt: now,
	})

	return slip, nil
}

type ConfirmPaymentInput struct {
	BookingID       uuid.UUID
	ActorID         uuid.UUID
	ExpectedVersion int
	PaymentRef      string
	Notes           string
}

func (s *BookingService) ConfirmPayment(ctx context.Context, input ConfirmPaymentInput) (*domain.Booking, error) {
	b, err := s.repo.GetBookingByID(ctx, input.BookingID)
	if err != nil {
		return nil, err
	}

	if b.OwnerID != input.ActorID {
		return nil, domain.ErrUnauthorizedAction
	}

	if err := domain.ValidateTransition(b.Status, domain.StatusConfirmed); err != nil {
		return nil, err
	}

	now := time.Now()
	b.Status = domain.StatusConfirmed
	b.ConfirmedAt = &now

	if err := s.repo.UpdateBookingStatusOCC(ctx, b, input.ExpectedVersion); err != nil {
		return nil, err
	}

	// Update latest slip to verified
	if slip, err := s.repo.GetLatestBankSlip(ctx, b.ID); err == nil && slip != nil {
		_ = s.repo.UpdateBankSlipStatus(ctx, slip.ID, domain.SlipStatusVerified, nil, &input.Notes)
	}

	_ = s.repo.CreateAuditLog(ctx, &domain.AuditLog{
		BookingID:  b.ID,
		ActorID:    &input.ActorID,
		Action:     "CONFIRM_PAYMENT",
		FromStatus: domain.StatusPaymentEvidenceUploaded,
		ToStatus:   domain.StatusConfirmed,
		Reason:     "Owner verified direct bank transfer receipt in full",
		Payload: map[string]any{
			"paymentReference": input.PaymentRef,
			"notes":            input.Notes,
		},
		CreatedAt: now,
	})

	return b, nil
}

type DeclinePaymentInput struct {
	BookingID       uuid.UUID
	ActorID         uuid.UUID
	ExpectedVersion int
	ReasonCode      string
	RejectionNotes  string
}

func (s *BookingService) DeclinePayment(ctx context.Context, input DeclinePaymentInput) (*domain.Booking, error) {
	b, err := s.repo.GetBookingByID(ctx, input.BookingID)
	if err != nil {
		return nil, err
	}

	if b.OwnerID != input.ActorID {
		return nil, domain.ErrUnauthorizedAction
	}

	if input.ReasonCode == "" {
		return nil, domain.ErrEmptyReason
	}

	if err := domain.ValidateTransition(b.Status, domain.StatusPaymentDeclined); err != nil {
		return nil, err
	}

	now := time.Now()
	b.Status = domain.StatusPaymentDeclined
	b.RejectionReasonCode = &input.ReasonCode
	b.RejectionNotes = &input.RejectionNotes

	if err := s.repo.UpdateBookingStatusOCC(ctx, b, input.ExpectedVersion); err != nil {
		return nil, err
	}

	if slip, err := s.repo.GetLatestBankSlip(ctx, b.ID); err == nil && slip != nil {
		_ = s.repo.UpdateBankSlipStatus(ctx, slip.ID, domain.SlipStatusRejected, &input.ReasonCode, &input.RejectionNotes)
	}

	_ = s.repo.CreateAuditLog(ctx, &domain.AuditLog{
		BookingID:  b.ID,
		ActorID:    &input.ActorID,
		Action:     "DECLINE_PAYMENT",
		FromStatus: domain.StatusPaymentEvidenceUploaded,
		ToStatus:   domain.StatusPaymentDeclined,
		Reason:     input.ReasonCode + ": " + input.RejectionNotes,
		Payload: map[string]any{
			"reasonCode":     input.ReasonCode,
			"rejectionNotes": input.RejectionNotes,
		},
		CreatedAt: now,
	})

	return b, nil
}

type DisputeBookingInput struct {
	BookingID       uuid.UUID
	ActorID         uuid.UUID
	ExpectedVersion int
	DisputeReason   string
	EvidenceURIs    []string
}

func (s *BookingService) DisputeBooking(ctx context.Context, input DisputeBookingInput) (*domain.Booking, error) {
	b, err := s.repo.GetBookingByID(ctx, input.BookingID)
	if err != nil {
		return nil, err
	}

	if b.RenterID != input.ActorID {
		return nil, domain.ErrUnauthorizedAction
	}

	if err := domain.ValidateTransition(b.Status, domain.StatusDisputeOpened); err != nil {
		return nil, err
	}

	now := time.Now()
	b.Status = domain.StatusDisputeOpened

	if err := s.repo.UpdateBookingStatusOCC(ctx, b, input.ExpectedVersion); err != nil {
		return nil, err
	}

	_ = s.repo.CreateAuditLog(ctx, &domain.AuditLog{
		BookingID:  b.ID,
		ActorID:    &input.ActorID,
		Action:     "DISPUTE_BOOKING",
		FromStatus: domain.StatusPaymentDeclined,
		ToStatus:   domain.StatusDisputeOpened,
		Reason:     input.DisputeReason,
		Payload: map[string]any{
			"disputeReason": input.DisputeReason,
			"evidenceUris":  input.EvidenceURIs,
		},
		CreatedAt: now,
	})

	return b, nil
}

type CancelBookingInput struct {
	BookingID       uuid.UUID
	ActorID         uuid.UUID
	ExpectedVersion int
	Reason          string
	IsForceMajeure  bool
}

func (s *BookingService) CancelBooking(ctx context.Context, input CancelBookingInput) (*domain.Booking, error) {
	b, err := s.repo.GetBookingByID(ctx, input.BookingID)
	if err != nil {
		return nil, err
	}

	isRenter := b.RenterID == input.ActorID
	isOwner := b.OwnerID == input.ActorID
	if !isRenter && !isOwner {
		return nil, domain.ErrUnauthorizedAction
	}

	if err := domain.ValidateTransition(b.Status, domain.StatusCancelled); err != nil {
		return nil, err
	}

	now := time.Now()
	res, err := domain.CalculateCancellation(b, isRenter, now, input.IsForceMajeure)
	if err != nil {
		return nil, err
	}

	// Release Valkey Lease
	var attIDs []uuid.UUID
	for _, att := range b.Attachments {
		attIDs = append(attIDs, att.AttachmentID)
	}
	_ = s.leaseManager.ReleaseHoldLease(ctx, b.FleetID, attIDs)

	b.Status = domain.StatusCancelled
	b.CancellationReason = &input.Reason
	b.CancelledBy = &input.ActorID
	b.PenaltyAmountMinor = res.PenaltyMinor
	b.RefundAmountMinor = res.RefundMinor
	b.CancelledAt = &now

	if err := s.repo.UpdateBookingStatusOCC(ctx, b, input.ExpectedVersion); err != nil {
		return nil, err
	}

	_ = s.repo.CreateAuditLog(ctx, &domain.AuditLog{
		BookingID:  b.ID,
		ActorID:    &input.ActorID,
		Action:     "CANCEL_BOOKING",
		FromStatus: b.Status,
		ToStatus:   domain.StatusCancelled,
		Reason:     input.Reason,
		Payload: map[string]any{
			"penaltyMinor":   res.PenaltyMinor,
			"refundMinor":    res.RefundMinor,
			"isOwnerDebt":    res.IsOwnerDebt,
			"isForceMajeure": input.IsForceMajeure,
		},
		CreatedAt: now,
	})

	return b, nil
}

type HandoverCheckInInput struct {
	BookingID           uuid.UUID
	ActorID             uuid.UUID
	ExpectedVersion     int
	HandoverPIN         string
	OperatingHoursMeter float64
	PhotoURIs           []string
}

func (s *BookingService) HandoverCheckIn(ctx context.Context, input HandoverCheckInInput) (*domain.Booking, error) {
	b, err := s.repo.GetBookingByID(ctx, input.BookingID)
	if err != nil {
		return nil, err
	}

	if b.RenterID != input.ActorID {
		return nil, domain.ErrUnauthorizedAction
	}

	if b.HandoverPIN == nil || *b.HandoverPIN != input.HandoverPIN {
		return nil, domain.ErrInvalidChallengePIN
	}

	// Verify meter against baseline
	fleet, err := s.fleetProvider.GetFleetAsset(ctx, b.FleetID)
	if err == nil && fleet != nil {
		if input.OperatingHoursMeter < fleet.BaselineMeterHours {
			return nil, domain.ErrInvalidMeterReading
		}
	}

	if err := domain.ValidateTransition(b.Status, domain.StatusActive); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateBookingHandover(ctx, b.ID, input.OperatingHoursMeter, input.ExpectedVersion); err != nil {
		return nil, err
	}

	_ = s.repo.CreateAuditLog(ctx, &domain.AuditLog{
		BookingID:  b.ID,
		ActorID:    &input.ActorID,
		Action:     "HANDOVER_CHECK_IN",
		FromStatus: b.Status,
		ToStatus:   domain.StatusActive,
		Reason:     "Field handover verified via 6-digit challenge PIN",
		Payload: map[string]any{
			"operatingHoursMeter": input.OperatingHoursMeter,
			"photoCount":          len(input.PhotoURIs),
		},
		CreatedAt: time.Now(),
	})

	return s.repo.GetBookingByID(ctx, b.ID)
}

type ReturnCheckOutInput struct {
	BookingID           uuid.UUID
	ActorID             uuid.UUID
	ExpectedVersion     int
	OperatingHoursMeter float64
	PhotoURIs           []string
	DamageNotes         string
}

func (s *BookingService) ReturnCheckOut(ctx context.Context, input ReturnCheckOutInput) (*domain.Booking, error) {
	b, err := s.repo.GetBookingByID(ctx, input.BookingID)
	if err != nil {
		return nil, err
	}

	if b.OwnerID != input.ActorID {
		return nil, domain.ErrUnauthorizedAction
	}

	if b.CheckInMeterHours != nil && input.OperatingHoursMeter < *b.CheckInMeterHours {
		return nil, domain.ErrInvalidMeterReading
	}

	targetStatus := domain.StatusCompleted
	if input.DamageNotes != "" {
		targetStatus = domain.StatusPendingFinalPayment
	}

	if err := s.repo.UpdateBookingReturn(ctx, b.ID, targetStatus, input.OperatingHoursMeter, input.ExpectedVersion); err != nil {
		return nil, err
	}

	_ = s.repo.CreateAuditLog(ctx, &domain.AuditLog{
		BookingID:  b.ID,
		ActorID:    &input.ActorID,
		Action:     "RETURN_CHECK_OUT",
		FromStatus: b.Status,
		ToStatus:   targetStatus,
		Reason:     "Rental return check-out completed",
		Payload: map[string]any{
			"checkOutMeter": input.OperatingHoursMeter,
			"damageNotes":   input.DamageNotes,
		},
		CreatedAt: time.Now(),
	})

	return s.repo.GetBookingByID(ctx, b.ID)
}

func (s *BookingService) GetBooking(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	return s.repo.GetBookingByID(ctx, id)
}

func (s *BookingService) ListUserBookings(ctx context.Context, userID uuid.UUID, role string, status string, limit int) ([]*domain.Booking, error) {
	if role == "owner" {
		return s.repo.ListBookingsByOwner(ctx, userID, status, limit)
	}
	return s.repo.ListBookingsByRenter(ctx, userID, status, limit)
}
