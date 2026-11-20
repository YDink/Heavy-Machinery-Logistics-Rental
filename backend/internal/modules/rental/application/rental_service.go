package application

import (
	"context"
	"fmt"
	"time"

	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	"github.com/dinkmachinery/dink-core/internal/modules/rental/domain"
	"github.com/google/uuid"
)

// BookingProvider decouples RentalService from the concrete BookingRepository.
type BookingProvider interface {
	GetBookingByID(ctx context.Context, id uuid.UUID) (*bookingDomain.Booking, error)
	UpdateBookingStatusOCC(ctx context.Context, booking *bookingDomain.Booking, expectedVersion int) error
	UpdateBookingHandover(ctx context.Context, bookingID uuid.UUID, checkInMeter float64, expectedVersion int) error
}

// LogisticsDeliveryChecker enforces D-05-08-001: parent machinery cannot enter ACTIVE until linked transport finishes delivery.
type LogisticsDeliveryChecker interface {
	EnsureParentReadyForActivation(ctx context.Context, parentBookingID uuid.UUID) error
}

type RentalService struct {
	rentalRepo       domain.RentalRepository
	bookingRepo      BookingProvider
	logisticsChecker LogisticsDeliveryChecker
}

func NewRentalService(rentalRepo domain.RentalRepository, bookingRepo BookingProvider) *RentalService {
	return &RentalService{
		rentalRepo:  rentalRepo,
		bookingRepo: bookingRepo,
	}
}

// SetLogisticsChecker attaches an optional logistics delivery verification guard.
func (s *RentalService) SetLogisticsChecker(checker LogisticsDeliveryChecker) {
	s.logisticsChecker = checker
}

// HandoverInspectCommand carries inputs for conducting pre-dispatch physical inspection.
type HandoverInspectCommand struct {
	BookingID         uuid.UUID
	InspectorID       uuid.UUID
	CheckInMeterHours float64
	FluidLevelsOk     bool
	SafetyGlassIntact bool
	Notes             string
	Photos            []domain.InspectionPhoto
}

// InspectHandover executes physical inspection recording per AR-01 & D-05-06-002.
func (s *RentalService) InspectHandover(ctx context.Context, cmd HandoverInspectCommand) (*domain.HandoverInspection, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, cmd.BookingID)
	if err != nil {
		return nil, fmt.Errorf("booking retrieval failed: %w", err)
	}

	// Handover inspection can only happen on confirmed bookings awaiting mobilization
	if booking.Status != bookingDomain.StatusConfirmed && booking.Status != bookingDomain.StatusReadyForPickup {
		return nil, fmt.Errorf("cannot inspect booking in state %s (must be CONFIRMED or READY_FOR_PICKUP)", booking.Status)
	}

	// Verify monotonic meter floor against certified baseline (AR-01)
	baselineMeter, err := s.rentalRepo.GetFleetBaselineMeter(ctx, booking.FleetID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fleet baseline meter: %w", err)
	}

	if err := domain.ValidateMonotonicMeterFloor(cmd.CheckInMeterHours, baselineMeter); err != nil {
		return nil, err
	}

	// Validate reticle photos
	if err := domain.ValidateHandoverPhotos(cmd.Photos); err != nil {
		return nil, err
	}

	inspection := &domain.HandoverInspection{
		ID:                uuid.New(),
		BookingID:         cmd.BookingID,
		InspectorID:       cmd.InspectorID,
		CheckInMeterHours: cmd.CheckInMeterHours,
		FluidLevelsOk:     cmd.FluidLevelsOk,
		SafetyGlassIntact: cmd.SafetyGlassIntact,
		Notes:             cmd.Notes,
		Photos:            cmd.Photos,
		CreatedAt:         time.Now().UTC(),
	}

	if err := s.rentalRepo.CreateHandoverInspection(ctx, inspection); err != nil {
		return nil, fmt.Errorf("failed to record handover inspection: %w", err)
	}

	return inspection, nil
}

// GenerateHandoverPIN generates a 6-digit challenge PIN on the owner app with 15m TTL (AR-02).
func (s *RentalService) GenerateHandoverPIN(ctx context.Context, bookingID uuid.UUID, callerUserID uuid.UUID) (*domain.HandoverChallenge, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, fmt.Errorf("booking retrieval failed: %w", err)
	}

	// Only owner can generate handover challenge PIN
	if booking.OwnerID != callerUserID {
		return nil, bookingDomain.ErrUnauthorizedAction
	}

	if booking.Status != bookingDomain.StatusConfirmed && booking.Status != bookingDomain.StatusReadyForPickup {
		return nil, fmt.Errorf("cannot generate handover PIN for booking in status %s", booking.Status)
	}

	pin, err := domain.GenerateChallengePIN()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	challenge := &domain.HandoverChallenge{
		ID:           uuid.New(),
		BookingID:    bookingID,
		ChallengePIN: pin,
		ExpiresAt:    now.Add(domain.HandoverChallengeTTL),
		AttemptCount: 0,
		MaxAttempts:  domain.HandoverChallengeMaxAttempts,
		IsVerified:   false,
		CreatedAt:    now,
	}

	if err := s.rentalRepo.CreateHandoverChallenge(ctx, challenge); err != nil {
		return nil, fmt.Errorf("failed to persist handover challenge: %w", err)
	}

	return challenge, nil
}

// HandoverConfirmCommand carries inputs for verifying PIN and activating rental.
type HandoverConfirmCommand struct {
	BookingID            uuid.UUID
	CallerUserID         uuid.UUID
	PIN                  string
	TouchSignatureBase64 string
	ContractorLat        *float64
	ContractorLng        *float64
}

// ConfirmHandover verifies the dual challenge PIN and transitions the booking to ACTIVE (AR-02).
func (s *RentalService) ConfirmHandover(ctx context.Context, cmd HandoverConfirmCommand) (*bookingDomain.Booking, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, cmd.BookingID)
	if err != nil {
		return nil, fmt.Errorf("booking retrieval failed: %w", err)
	}

	// Verify caller is renter or assigned operator
	isAuthorized := booking.RenterID == cmd.CallerUserID || (booking.AssignedOperatorID != nil && *booking.AssignedOperatorID == cmd.CallerUserID)
	if !isAuthorized {
		return nil, bookingDomain.ErrUnauthorizedAction
	}

	if booking.Status != bookingDomain.StatusConfirmed && booking.Status != bookingDomain.StatusReadyForPickup {
		return nil, fmt.Errorf("cannot activate booking from status %s", booking.Status)
	}

	if s.logisticsChecker != nil {
		if err := s.logisticsChecker.EnsureParentReadyForActivation(ctx, booking.ID); err != nil {
			return nil, err
		}
	}

	challenge, err := s.rentalRepo.GetActiveHandoverChallenge(ctx, cmd.BookingID)
	if err != nil {
		return nil, domain.ErrInvalidChallengePIN
	}

	if challenge.IsVerified {
		return nil, fmt.Errorf("%w: challenge has already been verified", domain.ErrInvalidChallengePIN)
	}

	if challenge.AttemptCount >= domain.HandoverChallengeMaxAttempts {
		return nil, domain.ErrPINLockout
	}

	now := time.Now().UTC()
	if now.After(challenge.ExpiresAt) {
		return nil, domain.ErrPINExpired
	}

	// Compare PIN
	if challenge.ChallengePIN != cmd.PIN {
		updated, incErr := s.rentalRepo.IncrementHandoverChallengeAttempts(ctx, challenge.ID)
		if incErr == nil && updated != nil && updated.AttemptCount >= domain.HandoverChallengeMaxAttempts {
			return nil, domain.ErrPINLockout
		}
		return nil, domain.ErrInvalidChallengePIN
	}

	// Mark verified
	if err := s.rentalRepo.VerifyHandoverChallenge(ctx, challenge.ID, cmd.TouchSignatureBase64, cmd.ContractorLat, cmd.ContractorLng); err != nil {
		return nil, fmt.Errorf("failed to verify challenge: %w", err)
	}

	// Determine starting meter reading
	startMeter := 0.0
	inspection, err := s.rentalRepo.GetHandoverInspection(ctx, booking.ID)
	if err == nil && inspection != nil {
		startMeter = inspection.CheckInMeterHours
	} else {
		baseline, _ := s.rentalRepo.GetFleetBaselineMeter(ctx, booking.FleetID)
		startMeter = baseline
	}

	// Transition booking to ACTIVE
	booking.Status = bookingDomain.StatusActive
	booking.ActivatedAt = &now
	booking.CheckInMeterHours = &startMeter
	if err := s.bookingRepo.UpdateBookingStatusOCC(ctx, booking, booking.Version); err != nil {
		return nil, fmt.Errorf("failed to transition booking to ACTIVE: %w", err)
	}

	// Initialize Day 1 RentalShift with 8.0 hr quota (AR-03)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	day1Shift := &domain.RentalShift{
		ID:               uuid.New(),
		BookingID:        booking.ID,
		OperatorID:       booking.AssignedOperatorID,
		ShiftDate:        today,
		StartTime:        now,
		StartMeterHours:  startMeter,
		AllocatedHours:   domain.StandardShiftHoursPerDay,
		AccumulatedHours: 0.0,
		ExcessHours:      0.0,
		IsPaused:         false,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	_ = s.rentalRepo.CreateRentalShift(ctx, day1Shift)

	return booking, nil
}

// GetActiveRentalHUD compiles telemetry, shift progression, and permitted actions (AR-03, AR-04).
func (s *RentalService) GetActiveRentalHUD(ctx context.Context, bookingID uuid.UUID, callerUserID uuid.UUID) (*domain.ActiveRentalHUD, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, fmt.Errorf("booking retrieval failed: %w", err)
	}

	isAuthorized := booking.RenterID == callerUserID || booking.OwnerID == callerUserID ||
		(booking.AssignedOperatorID != nil && *booking.AssignedOperatorID == callerUserID)
	if !isAuthorized {
		return nil, bookingDomain.ErrUnauthorizedAction
	}

	activeShift, _ := s.rentalRepo.GetActiveRentalShift(ctx, bookingID)
	shifts, _ := s.rentalRepo.ListRentalShifts(ctx, bookingID)

	var totalExcess float64
	var accumToday float64
	shiftPaused := false

	for _, sh := range shifts {
		totalExcess += sh.ExcessHours
	}

	if activeShift != nil {
		accumToday = activeShift.AccumulatedHours
		shiftPaused = activeShift.IsPaused
	}

	activeBreakdown, _ := s.rentalRepo.GetActiveMechanicalBreakdown(ctx, bookingID)
	if activeBreakdown != nil {
		shiftPaused = true
	}

	// Calculate rental day
	rentalDay := 1
	if booking.ActivatedAt != nil {
		daysPassed := int(time.Since(*booking.ActivatedAt).Hours() / 24)
		rentalDay = daysPassed + 1
		if rentalDay > booking.DurationDays {
			rentalDay = booking.DurationDays
		}
	}

	baselineMeter := 0.0
	inspection, err := s.rentalRepo.GetHandoverInspection(ctx, booking.ID)
	if err == nil && inspection != nil {
		baselineMeter = inspection.CheckInMeterHours
	} else {
		bm, _ := s.rentalRepo.GetFleetBaselineMeter(ctx, booking.FleetID)
		baselineMeter = bm
	}

	hud := &domain.ActiveRentalHUD{
		BookingID:             booking.ID,
		Status:                string(booking.Status),
		FleetID:               booking.FleetID,
		FleetTitle:            booking.BookingReference,
		PlateNumber:           "ETH-MACH-01",
		OperatorID:            booking.AssignedOperatorID,
		CurrentRentalDay:      rentalDay,
		TotalRentalDays:       booking.DurationDays,
		HandoverBaselineHours: baselineMeter,
		AccumulatedHoursToday: accumToday,
		AllowedHoursPerDay:    domain.StandardShiftHoursPerDay,
		ExcessHoursTotal:      totalExcess,
		ShiftTimerPaused:      shiftPaused,
		CanExtend:             booking.Status == bookingDomain.StatusActive && callerUserID == booking.RenterID,
		CanReportBreakdown:    booking.Status == bookingDomain.StatusActive && activeBreakdown == nil,
		CanInitiateReturn:     booking.Status == bookingDomain.StatusActive && callerUserID == booking.RenterID,
	}

	return hud, nil
}

// RentalExtensionCommand carries parameters for requesting a forward extension.
type RentalExtensionCommand struct {
	BookingID  uuid.UUID
	RenterID   uuid.UUID
	NewEndDate time.Time
}

// RequestExtension computes 50% prepay quote and submits an extension (AR-04).
func (s *RentalService) RequestExtension(ctx context.Context, cmd RentalExtensionCommand) (*domain.RentalExtension, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, cmd.BookingID)
	if err != nil {
		return nil, fmt.Errorf("booking retrieval failed: %w", err)
	}

	if booking.RenterID != cmd.RenterID {
		return nil, bookingDomain.ErrUnauthorizedAction
	}

	if booking.Status != bookingDomain.StatusActive && booking.Status != bookingDomain.StatusActiveExtensionHold {
		return nil, domain.ErrBookingNotActive
	}

	// Pre-flight forward availability check
	hasConflict, err := s.rentalRepo.CheckForwardBookingConflict(ctx, booking.FleetID, booking.ID, booking.EndDate, cmd.NewEndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to check fleet forward availability: %w", err)
	}
	if hasConflict {
		return nil, domain.ErrExtensionDateUnavailable
	}

	// Calculate incremental gross and prepayment
	bundledDailyMinor := int64(0)
	if booking.DurationDays > 0 {
		bundledDailyMinor = booking.AttachmentsFeeMinor / int64(booking.DurationDays)
	}

	extDays, incGross, incPrepay, err := domain.CalculateExtensionQuote(booking.EndDate, cmd.NewEndDate, booking.DailyRateAppliedMinor, bundledDailyMinor)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	ext := &domain.RentalExtension{
		ID:                         uuid.New(),
		BookingID:                  booking.ID,
		RequestedByID:              cmd.RenterID,
		OriginalEndDate:            booking.EndDate,
		NewEndDate:                 cmd.NewEndDate,
		ExtensionDays:              extDays,
		IncrementalGrossMinor:      incGross,
		IncrementalPrepaymentMinor: incPrepay,
		Currency:                   booking.Currency,
		Status:                     "REQUESTED",
		CreatedAt:                  now,
		UpdatedAt:                  now,
	}

	if err := s.rentalRepo.CreateRentalExtension(ctx, ext); err != nil {
		return nil, fmt.Errorf("failed to record extension request: %w", err)
	}

	// Put booking on extension hold
	booking.Status = bookingDomain.StatusActiveExtensionHold
	_ = s.bookingRepo.UpdateBookingStatusOCC(ctx, booking, booking.Version)

	return ext, nil
}

// BreakdownReportCommand carries breakdown details.
type BreakdownReportCommand struct {
	BookingID              uuid.UUID
	ReporterID             uuid.UUID
	FaultCategory          string
	FaultDescription       string
	DiagnosticTroubleCodes []string
	CurrentHourMeter       float64
	IncidentPhotos         []string
}

// ReportMechanicalBreakdown records breakdown incident, pauses shift timer, and starts 4h SLA (AR-05).
func (s *RentalService) ReportMechanicalBreakdown(ctx context.Context, cmd BreakdownReportCommand) (*domain.MechanicalBreakdown, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, cmd.BookingID)
	if err != nil {
		return nil, fmt.Errorf("booking retrieval failed: %w", err)
	}

	isAuthorized := booking.RenterID == cmd.ReporterID ||
		(booking.AssignedOperatorID != nil && *booking.AssignedOperatorID == cmd.ReporterID)
	if !isAuthorized {
		return nil, bookingDomain.ErrUnauthorizedAction
	}

	if booking.Status != bookingDomain.StatusActive && booking.Status != bookingDomain.StatusActiveExtensionHold {
		return nil, domain.ErrBookingNotActive
	}

	// Check if already paused
	existing, _ := s.rentalRepo.GetActiveMechanicalBreakdown(ctx, cmd.BookingID)
	if existing != nil {
		return nil, domain.ErrBookingAlreadyPaused
	}

	if err := domain.ValidateFaultCategory(cmd.FaultCategory); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	breakdown := &domain.MechanicalBreakdown{
		ID:                     uuid.New(),
		BookingID:              booking.ID,
		ReportedByID:           cmd.ReporterID,
		FaultCategory:          cmd.FaultCategory,
		FaultDescription:       cmd.FaultDescription,
		DiagnosticTroubleCodes: cmd.DiagnosticTroubleCodes,
		CurrentHourMeter:       cmd.CurrentHourMeter,
		IncidentPhotos:         cmd.IncidentPhotos,
		Status:                 "ACTIVE_BREAKDOWN_PAUSED",
		SLADeadline:            now.Add(domain.BreakdownOwnerSLADuration),
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	if err := s.rentalRepo.CreateMechanicalBreakdown(ctx, breakdown); err != nil {
		return nil, fmt.Errorf("failed to record mechanical breakdown: %w", err)
	}

	// Pause active shift
	activeShift, err := s.rentalRepo.GetActiveRentalShift(ctx, booking.ID)
	if err == nil && activeShift != nil && !activeShift.IsPaused {
		_ = s.rentalRepo.UpdateRentalShiftProgress(ctx, activeShift.ID, activeShift.AccumulatedHours, activeShift.ExcessHours, &cmd.CurrentHourMeter, true)
	}

	// Transition booking status to ACTIVE_BREAKDOWN_PAUSED
	booking.Status = bookingDomain.StatusActiveBreakdownPaused
	_ = s.bookingRepo.UpdateBookingStatusOCC(ctx, booking, booking.Version)

	return breakdown, nil
}

// BreakdownResolveCommand carries resolution information.
type BreakdownResolveCommand struct {
	BookingID       uuid.UUID
	BreakdownID     uuid.UUID
	ResolverID      uuid.UUID
	MechanicName    string
	MechanicPhone   string
	ResolutionNotes string
}

// ResolveMechanicalBreakdown clears breakdown, resumes shift telemetry, and restores ACTIVE status (AR-05).
func (s *RentalService) ResolveMechanicalBreakdown(ctx context.Context, cmd BreakdownResolveCommand) (*domain.MechanicalBreakdown, error) {
	var breakdown *domain.MechanicalBreakdown
	var err error

	if cmd.BreakdownID != uuid.Nil {
		breakdown, err = s.rentalRepo.GetMechanicalBreakdown(ctx, cmd.BreakdownID)
	} else if cmd.BookingID != uuid.Nil {
		breakdown, err = s.rentalRepo.GetActiveMechanicalBreakdown(ctx, cmd.BookingID)
	} else {
		return nil, domain.ErrBreakdownNotFound
	}

	if err != nil || breakdown == nil {
		return nil, domain.ErrBreakdownNotFound
	}

	booking, err := s.bookingRepo.GetBookingByID(ctx, breakdown.BookingID)
	if err != nil {
		return nil, fmt.Errorf("booking retrieval failed: %w", err)
	}

	// Owner or authorized dispatcher can resolve
	if booking.OwnerID != cmd.ResolverID {
		return nil, bookingDomain.ErrUnauthorizedAction
	}

	if breakdown.Status != "ACTIVE_BREAKDOWN_PAUSED" && breakdown.Status != "MECHANIC_DISPATCHED" {
		return nil, domain.ErrActiveBreakdownStillOpen
	}

	if err := s.rentalRepo.ResolveMechanicalBreakdown(ctx, breakdown.ID, cmd.MechanicName, cmd.MechanicPhone, cmd.ResolutionNotes); err != nil {
		return nil, fmt.Errorf("failed to resolve breakdown: %w", err)
	}

	// Resume active shift
	activeShift, err := s.rentalRepo.GetActiveRentalShift(ctx, booking.ID)
	if err == nil && activeShift != nil && activeShift.IsPaused {
		_ = s.rentalRepo.UpdateRentalShiftProgress(ctx, activeShift.ID, activeShift.AccumulatedHours, activeShift.ExcessHours, nil, false)
	}

	// Restore booking status to ACTIVE
	booking.Status = bookingDomain.StatusActive
	_ = s.bookingRepo.UpdateBookingStatusOCC(ctx, booking, booking.Version)

	breakdown.Status = "RESOLVED"
	now := time.Now().UTC()
	breakdown.ResolvedAt = &now
	breakdown.MechanicName = cmd.MechanicName
	breakdown.MechanicPhone = cmd.MechanicPhone
	breakdown.ResolutionNotes = cmd.ResolutionNotes

	return breakdown, nil
}


// Return Inspection & Completion (Partition 05-07)

// InitiateReturn halts shift telemetry and transitions booking to RETURN_INITIATED.
func (s *RentalService) InitiateReturn(ctx context.Context, bookingID uuid.UUID, actorID uuid.UUID) (*bookingDomain.Booking, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, fmt.Errorf("booking retrieval failed: %w", err)
	}

	if booking.Status != bookingDomain.StatusActive && booking.Status != bookingDomain.StatusActiveBreakdownPaused {
		return nil, fmt.Errorf("cannot initiate return on booking in state %s (must be ACTIVE)", booking.Status)
	}

	// End active shift if running
	activeShift, err := s.rentalRepo.GetActiveRentalShift(ctx, booking.ID)
	if err == nil && activeShift != nil {
		_ = s.rentalRepo.EndRentalShift(ctx, activeShift.ID, activeShift.StartMeterHours+activeShift.AccumulatedHours, activeShift.AccumulatedHours, activeShift.ExcessHours)
	}

	booking.Status = bookingDomain.StatusReturnInitiated
	if err := s.bookingRepo.UpdateBookingStatusOCC(ctx, booking, booking.Version); err != nil {
		return nil, fmt.Errorf("failed to update booking to RETURN_INITIATED: %w", err)
	}

	return booking, nil
}

// ReturnInspectCommand carries inputs for conducting 6-angle off-hire inspection.
type ReturnInspectCommand struct {
	BookingID           uuid.UUID
	InspectorID         uuid.UUID
	CheckOutMeterHours  float64
	FluidLevelsOk       bool
	SafetyGlassIntact   bool
	AccelerometerSteady bool
	Notes               string
	Photos              []domain.InspectionPhoto
}

// InspectReturn records the 6-angle return inspection and checks non-regressive monotonic meter floor.
func (s *RentalService) InspectReturn(ctx context.Context, cmd ReturnInspectCommand) (*domain.ReturnInspection, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, cmd.BookingID)
	if err != nil {
		return nil, fmt.Errorf("booking retrieval failed: %w", err)
	}

	if booking.Status != bookingDomain.StatusReturnInitiated {
		return nil, fmt.Errorf("cannot inspect return on booking in state %s (must be RETURN_INITIATED)", booking.Status)
	}

	if err := domain.ValidateReturnPhotos(cmd.Photos); err != nil {
		return nil, err
	}

	checkInMeter := 0.0
	if booking.CheckInMeterHours != nil {
		checkInMeter = *booking.CheckInMeterHours
	}

	if err := domain.ValidateReturnMeterMonotonicity(checkInMeter, cmd.CheckOutMeterHours); err != nil {
		return nil, err
	}

	inspection := &domain.ReturnInspection{
		ID:                  uuid.New(),
		BookingID:           cmd.BookingID,
		InspectorID:         cmd.InspectorID,
		CheckOutMeterHours:  cmd.CheckOutMeterHours,
		FluidLevelsOK:       cmd.FluidLevelsOk,
		SafetyGlassIntact:   cmd.SafetyGlassIntact,
		AccelerometerSteady: cmd.AccelerometerSteady,
		Notes:               cmd.Notes,
		Photos:              cmd.Photos,
		CreatedAt:           time.Now().UTC(),
	}

	if err := s.rentalRepo.CreateReturnInspection(ctx, inspection); err != nil {
		return nil, fmt.Errorf("failed to persist return inspection: %w", err)
	}

	booking.Status = bookingDomain.StatusDamagePendingReview
	booking.CheckOutMeterHours = &cmd.CheckOutMeterHours
	_ = s.bookingRepo.UpdateBookingStatusOCC(ctx, booking, booking.Version)

	return inspection, nil
}

// DamageClaimCommand carries inputs for lodging post-rental structural damage dockets.
type DamageClaimCommand struct {
	BookingID          uuid.UUID
	ClaimantID         uuid.UUID
	ComponentCategory  string
	IncidentNarrative  string
	RepairQuoteURI     string
	ClaimedAmountMinor int64
	Currency           string
	DamagePhotos       []string
}

// LodgeDamageClaim files an itemized damage compensation claim docket.
func (s *RentalService) LodgeDamageClaim(ctx context.Context, cmd DamageClaimCommand) (*domain.DamageClaimDocket, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, cmd.BookingID)
	if err != nil {
		return nil, fmt.Errorf("booking retrieval failed: %w", err)
	}

	if booking.OwnerID != cmd.ClaimantID {
		return nil, bookingDomain.ErrUnauthorizedAction
	}

	docket := &domain.DamageClaimDocket{
		ID:                  uuid.New(),
		BookingID:           cmd.BookingID,
		ClaimantID:          cmd.ClaimantID,
		ComponentCategory:   cmd.ComponentCategory,
		IncidentNarrative:   cmd.IncidentNarrative,
		RepairQuoteURI:      cmd.RepairQuoteURI,
		ClaimedAmountMinor:  cmd.ClaimedAmountMinor,
		AcceptedAmountMinor: 0,
		Currency:            cmd.Currency,
		DamagePhotos:        cmd.DamagePhotos,
		Status:              "PENDING_REVIEW",
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
	}

	if err := s.rentalRepo.CreateDamageClaimDocket(ctx, docket); err != nil {
		return nil, fmt.Errorf("failed to lodge damage claim: %w", err)
	}

	booking.DamageClaimSubmitted = true
	_ = s.bookingRepo.UpdateBookingStatusOCC(ctx, booking, booking.Version)

	return docket, nil
}

// RespondDamageClaimCommand carries contractor action on lodged damage claim.
type RespondDamageClaimCommand struct {
	BookingID           uuid.UUID
	ClaimID             uuid.UUID
	ActorID             uuid.UUID
	Action              string
	AcceptedAmountMinor int64
	DisputeReason       string
}

// RespondToDamageClaim records contractor acceptance or dispute of filed damages.
func (s *RentalService) RespondToDamageClaim(ctx context.Context, cmd RespondDamageClaimCommand) (*domain.DamageClaimDocket, error) {
	docket, err := s.rentalRepo.GetDamageClaimDocket(ctx, cmd.ClaimID)
	if err != nil {
		return nil, err
	}

	booking, err := s.bookingRepo.GetBookingByID(ctx, cmd.BookingID)
	if err != nil {
		return nil, fmt.Errorf("booking retrieval failed: %w", err)
	}

	if booking.RenterID != cmd.ActorID {
		return nil, bookingDomain.ErrUnauthorizedAction
	}

	if cmd.Action == "ACCEPT" {
		if err := s.rentalRepo.RespondToDamageClaim(ctx, cmd.ClaimID, "ACCEPTED", cmd.AcceptedAmountMinor, ""); err != nil {
			return nil, err
		}
		booking.AcceptedDamageMinor = cmd.AcceptedAmountMinor
		_ = s.bookingRepo.UpdateBookingStatusOCC(ctx, booking, booking.Version)
		docket.Status = "ACCEPTED"
		docket.AcceptedAmountMinor = cmd.AcceptedAmountMinor
	} else {
		if err := s.rentalRepo.RespondToDamageClaim(ctx, cmd.ClaimID, "DISPUTED", 0, cmd.DisputeReason); err != nil {
			return nil, err
		}
		booking.Status = bookingDomain.StatusDamageDisputed
		_ = s.bookingRepo.UpdateBookingStatusOCC(ctx, booking, booking.Version)
		docket.Status = "DISPUTED"
		docket.DisputeReason = cmd.DisputeReason
	}

	return docket, nil
}

// ReconcileFinalSettlement compiles the deterministic reconciliation statement per RC-05.
func (s *RentalService) ReconcileFinalSettlement(ctx context.Context, bookingID uuid.UUID) (*domain.RentalSettlement, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, fmt.Errorf("booking retrieval failed: %w", err)
	}

	checkIn := 0.0
	if booking.CheckInMeterHours != nil {
		checkIn = *booking.CheckInMeterHours
	}
	checkOut := checkIn
	if booking.CheckOutMeterHours != nil {
		checkOut = *booking.CheckOutMeterHours
	}

	// Calculate credited breakdown hours
	breakdownCreditHours := 0.0
	breakdowns, _ := s.rentalRepo.ListMechanicalBreakdowns(ctx, bookingID)
	for _, b := range breakdowns {
		if b.Status == "RESOLVED" && b.ResolvedAt != nil {
			duration := b.ResolvedAt.Sub(b.CreatedAt).Hours()
			if duration > 0 {
				breakdownCreditHours += duration
			}
		}
	}

	hourlyRateMinor := booking.DailyRateAppliedMinor / 8
	excessQuote := domain.CalculateExcessHours(booking.DurationDays, checkIn, checkOut, breakdownCreditHours, hourlyRateMinor)
	reconStmt := domain.ReconcileFinalSettlement(
		booking.BaseHireFeeMinor,
		booking.AttachmentsFeeMinor,
		excessQuote.SurchargeAmountMinor,
		booking.AcceptedDamageMinor,
		booking.PrepaymentAmountMinor,
	)

	settlement := &domain.RentalSettlement{
		ID:                         uuid.New(),
		BookingID:                  bookingID,
		BaseMachineFeeMinor:        booking.BaseHireFeeMinor,
		BundledAttachmentsFeeMinor: booking.AttachmentsFeeMinor,
		ContractualAllowedHours:    excessQuote.AllowedHours,
		ActualOperatingHours:       excessQuote.ActualGrossHours,
		CreditedBreakdownHours:     excessQuote.CreditedBreakdownHours,
		BillableExcessHours:        excessQuote.BillableExcessHours,
		ExcessHourlyRateMinor:      hourlyRateMinor,
		ExcessHoursSurchargeMinor:  excessQuote.SurchargeAmountMinor,
		AcceptedDamageMinor:        booking.AcceptedDamageMinor,
		GrossLiabilitiesMinor:      reconStmt.GrossLiabilitiesMinor,
		PrepaymentSettledMinor:     booking.PrepaymentAmountMinor,
		NetBalanceMinor:            reconStmt.NetBalanceMinor,
		Currency:                   booking.Currency,
		Status:                     reconStmt.SettlementStatus,
		CreatedAt:                  time.Now().UTC(),
		UpdatedAt:                  time.Now().UTC(),
	}

	if err := s.rentalRepo.UpsertRentalSettlement(ctx, settlement); err != nil {
		return nil, fmt.Errorf("failed to persist settlement statement: %w", err)
	}

	if reconStmt.NetBalanceMinor > 0 {
		booking.Status = bookingDomain.StatusPendingFinalPayment
	} else if reconStmt.NetBalanceMinor == 0 {
		booking.Status = bookingDomain.StatusCompleted
	}
	booking.ExcessHoursSurchargeMinor = excessQuote.SurchargeAmountMinor
	booking.FinalNetBalanceMinor = reconStmt.NetBalanceMinor
	_ = s.bookingRepo.UpdateBookingStatusOCC(ctx, booking, booking.Version)

	return settlement, nil
}

// SubmitReviewCommand carries inputs for submitting double-blind mutual scorecard.
type SubmitReviewCommand struct {
	BookingID        uuid.UUID
	AuthorID         uuid.UUID
	TargetID         uuid.UUID
	Role             string
	Rating           float64
	Feedback         string
	DimensionRatings map[string]float64
}

// SubmitBookingReview records an isolated double-blind review, unblinding upon bilateral submission.
func (s *RentalService) SubmitBookingReview(ctx context.Context, cmd SubmitReviewCommand) (*domain.BookingReview, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, cmd.BookingID)
	if err != nil {
		return nil, fmt.Errorf("booking retrieval failed: %w", err)
	}

	if cmd.TargetID == uuid.Nil {
		if cmd.AuthorID == booking.RenterID {
			cmd.TargetID = booking.OwnerID
			cmd.Role = "RENTER"
		} else if cmd.AuthorID == booking.OwnerID {
			cmd.TargetID = booking.RenterID
			cmd.Role = "OWNER"
		} else {
			return nil, bookingDomain.ErrUnauthorizedAction
		}
	}

	if err := domain.ValidateBookingReview(cmd.AuthorID, cmd.TargetID, cmd.Rating, cmd.Role); err != nil {
		return nil, err
	}

	review := &domain.BookingReview{
		ID:               uuid.New(),
		BookingID:        cmd.BookingID,
		AuthorID:         cmd.AuthorID,
		TargetID:         cmd.TargetID,
		Role:             cmd.Role,
		Rating:           cmd.Rating,
		Feedback:         cmd.Feedback,
		DimensionRatings: cmd.DimensionRatings,
		IsPublished:      false,
		CreatedAt:        time.Now().UTC(),
	}

	if err := s.rentalRepo.CreateBookingReview(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to persist booking review: %w", err)
	}

	// Check if counterparty has also reviewed
	existingReviews, err := s.rentalRepo.GetBookingReviews(ctx, cmd.BookingID)
	if err == nil && len(existingReviews) >= 2 {
		_ = s.rentalRepo.PublishBookingReviews(ctx, cmd.BookingID)
		booking.Status = bookingDomain.StatusRated
		_ = s.bookingRepo.UpdateBookingStatusOCC(ctx, booking, booking.Version)
		review.IsPublished = true
	}

	return review, nil
}
