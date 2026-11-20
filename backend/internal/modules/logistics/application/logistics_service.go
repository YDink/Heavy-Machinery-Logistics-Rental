package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	"github.com/dinkmachinery/dink-core/internal/modules/logistics/domain"
	machineryDomain "github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
	"github.com/google/uuid"
)

// BookingPort decouples LogisticsService from concrete booking repositories.
type BookingPort interface {
	GetBookingByID(ctx context.Context, id uuid.UUID) (*bookingDomain.Booking, error)
	UpdateBookingStatusOCC(ctx context.Context, b *bookingDomain.Booking, expectedVersion int) error
	CreateAuditLog(ctx context.Context, log *bookingDomain.AuditLog) error
}

// LogisticsService orchestrates dimensional matching, sub-bookings, dispatch, transshipment, and reverse cascades.
type LogisticsService struct {
	repo        domain.LogisticsRepository
	bookingRepo BookingPort
}

func NewLogisticsService(repo domain.LogisticsRepository, bookingRepo BookingPort) *LogisticsService {
	return &LogisticsService{
		repo:        repo,
		bookingRepo: bookingRepo,
	}
}

type MatchCarriersCommand struct {
	ParentBookingID uuid.UUID
	CallerUserID    uuid.UUID
	SearchRadiusKm  int
	PickupDate      time.Time
}

type MatchCarriersResult struct {
	ParentBookingID            uuid.UUID
	SelfMoveCapable            bool
	RecommendedTrailerClass    string
	AbnormalLoadEscortRequired bool
	MachineSpecs               domain.MachineLogisticsSpecs
	RequiredMinimumPayloadKg   float64
	AvailableLogisticsFleets   []domain.MatchedCarrier
}

// MatchCarriers evaluates machine transport dimensions and discovers compatible lowboy trailers (FL-28, FL-29, SRS §3.3).
func (s *LogisticsService) MatchCarriers(ctx context.Context, cmd MatchCarriersCommand) (*MatchCarriersResult, error) {
	parent, err := s.bookingRepo.GetBookingByID(ctx, cmd.ParentBookingID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrLogisticsBookingNotFound, err)
	}

	if parent.OwnerID != cmd.CallerUserID && parent.RenterID != cmd.CallerUserID {
		return nil, domain.ErrUnauthorizedLogisticsAction
	}

	profile, err := s.repo.GetFleetLogisticsProfile(ctx, parent.FleetID)
	if err != nil {
		return nil, fmt.Errorf("failed to load fleet logistics profile: %w", err)
	}

	limits, err := s.repo.GetCountryLogisticsLimits(ctx, profile.CountryCode)
	if err != nil {
		return nil, fmt.Errorf("failed to load country logistics limits: %w", err)
	}

	if profile.SelfMoveCapable {
		return &MatchCarriersResult{
			ParentBookingID:            parent.ID,
			SelfMoveCapable:            true,
			RecommendedTrailerClass:    "Self-Move Capable (Roadway Direct)",
			AbnormalLoadEscortRequired: false,
			MachineSpecs:               profile.MachineSpecs,
			RequiredMinimumPayloadKg:   0,
			AvailableLogisticsFleets:   []domain.MatchedCarrier{},
		}, nil
	}

	radiusMeters := float64(cmd.SearchRadiusKm) * 1000.0
	if radiusMeters <= 0 {
		radiusMeters = 75000.0
	}

	matches, err := s.repo.MatchCompatibleLogisticsFleets(
		ctx,
		profile.Location,
		radiusMeters,
		profile.CountryCode,
		parent.OwnerID,
		profile.MachineSpecs,
		limits,
		25,
	)
	if err != nil {
		return nil, err
	}

	recommendedClass := machineryDomain.InferLogisticsTransportClass(profile.MachineSpecs.OperatingWeightKg)
	abnormalRequired := profile.MachineSpecs.TransportWidthMm > limits.WideLoadWidthMm

	return &MatchCarriersResult{
		ParentBookingID:            parent.ID,
		SelfMoveCapable:            false,
		RecommendedTrailerClass:    recommendedClass,
		AbnormalLoadEscortRequired: abnormalRequired,
		MachineSpecs:               profile.MachineSpecs,
		RequiredMinimumPayloadKg:   domain.ComputeMinimumRequiredPayloadKg(profile.MachineSpecs.OperatingWeightKg),
		AvailableLogisticsFleets:   matches,
	}, nil
}

type CreateSubBookingCommand struct {
	ParentBookingID  uuid.UUID
	CallerUserID     uuid.UUID
	LogisticsFleetID uuid.UUID
	PickupDateTime   time.Time
	Origin           domain.GeoPoint
	Destination      domain.GeoPoint
}

type SubBookingResult struct {
	SubBookingID               uuid.UUID
	ParentBookingID            uuid.UUID
	Status                     bookingDomain.BookingStatus
	FreightHireFeeMinor        int64
	PlatformCommissionBps      int
	CommissionOwedMinor        int64
	Currency                   string
	AbnormalLoadEscortRequired bool
	EscortFeeMinor             int64
	CarrierFullName            string
	CarrierPhoneNumber         string
}

// CreateSubBooking instantiates a linked child logistics sub-booking with 2.5% freight commission discount (FL-31, SRS §4.2, §4.4).
func (s *LogisticsService) CreateSubBooking(ctx context.Context, cmd CreateSubBookingCommand) (*SubBookingResult, error) {
	parent, err := s.bookingRepo.GetBookingByID(ctx, cmd.ParentBookingID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrLogisticsBookingNotFound, err)
	}

	if parent.OwnerID != cmd.CallerUserID {
		return nil, domain.ErrUnauthorizedLogisticsAction
	}

	if parent.Status != bookingDomain.StatusConfirmed && parent.Status != bookingDomain.StatusAwaitingFleetReadiness {
		return nil, fmt.Errorf("%w: parent booking must be CONFIRMED or AWAITING_FLEET_READINESS (got %s)", domain.ErrInvalidBookingStateForLogistics, parent.Status)
	}

	if existing, err := s.repo.GetLogisticsDispatchByParentBookingID(ctx, parent.ID); err == nil && existing != nil {
		return nil, domain.ErrSubBookingAlreadyExists
	}

	machineProfile, err := s.repo.GetFleetLogisticsProfile(ctx, parent.FleetID)
	if err != nil {
		return nil, fmt.Errorf("failed to load machine profile: %w", err)
	}
	if machineProfile.SelfMoveCapable {
		return nil, domain.ErrSelfMoveNotRequired
	}

	carrierProfile, err := s.repo.GetFleetLogisticsProfile(ctx, cmd.LogisticsFleetID)
	if err != nil {
		return nil, fmt.Errorf("failed to load logistics carrier fleet: %w", err)
	}
	if carrierProfile.RootCategoryCode != "LOGISTICS" {
		return nil, fmt.Errorf("%w: selected fleet is not in LOGISTICS category", domain.ErrInvalidBookingStateForLogistics)
	}

	limits, err := s.repo.GetCountryLogisticsLimits(ctx, machineProfile.CountryCode)
	if err != nil {
		return nil, err
	}

	abnormalEscort, err := domain.EvaluateDimensionalCompatibility(machineProfile.MachineSpecs, carrierProfile.TrailerProfile, limits)
	if err != nil {
		return nil, err
	}

	// Standard mobilization leg quoted at 2 days (outbound + return/transit) or minimum 1 day
	durationDays := 2
	freightHireFeeMinor := carrierProfile.DailyRateMinor * int64(durationDays)
	commissionOwedMinor := domain.CalculatePassThroughFreightCommission(freightHireFeeMinor)

	var escortFeeMinor int64
	if abnormalEscort {
		escortFeeMinor = 250000 // 2,500.00 ETB escort permit & pilot vehicle allowance
	}

	currency := parent.Currency
	if currency == "" {
		currency = limits.DefaultCurrency
	}

	subRef := fmt.Sprintf("BK-SUB-%s", strings.ToUpper(uuid.NewString()[:8]))
	startDate := cmd.PickupDateTime.UTC().Truncate(24 * time.Hour)
	endDate := startDate.AddDate(0, 0, durationDays)

	dispatch := &domain.LogisticsDispatch{
		ParentBookingID:            parent.ID,
		LogisticsFleetID:           carrierProfile.FleetID,
		CarrierOwnerID:             carrierProfile.OwnerID,
		PickupDateTime:             cmd.PickupDateTime.UTC(),
		Origin:                     cmd.Origin,
		Destination:                cmd.Destination,
		FreightHireFeeMinor:        freightHireFeeMinor,
		PlatformCommissionBps:      domain.PassThroughFreightCommissionBps,
		CommissionOwedMinor:        commissionOwedMinor,
		Currency:                   currency,
		AbnormalLoadEscortRequired: abnormalEscort,
		EscortFeeMinor:             escortFeeMinor,
	}

	subBookingID, err := s.repo.CreateLogisticsSubBookingAndDispatch(
		ctx,
		subRef,
		parent.ID,
		parent.OwnerID,
		carrierProfile.OwnerID,
		carrierProfile.FleetID,
		startDate,
		endDate,
		durationDays,
		carrierProfile.DailyRateMinor,
		freightHireFeeMinor,
		domain.PassThroughFreightCommissionBps,
		commissionOwedMinor,
		abnormalEscort,
		currency,
		cmd.Destination.LocationName,
		dispatch,
	)
	if err != nil {
		return nil, err
	}

	_ = s.bookingRepo.CreateAuditLog(ctx, &bookingDomain.AuditLog{
		BookingID:  parent.ID,
		ActorID:    &cmd.CallerUserID,
		Action:     "LOGISTICS_SUB_BOOKING_CREATED",
		FromStatus: parent.Status,
		ToStatus:   parent.Status,
		Reason:     "Linked in-platform lowboy sub-booking instantiated",
		Payload: map[string]any{
			"subBookingId":        subBookingID.String(),
			"logisticsFleetId":    carrierProfile.FleetID.String(),
			"commissionRateBps":   domain.PassThroughFreightCommissionBps,
			"commissionOwedMinor": commissionOwedMinor,
		},
		CreatedAt: time.Now().UTC(),
	})

	return &SubBookingResult{
		SubBookingID:               subBookingID,
		ParentBookingID:            parent.ID,
		Status:                     bookingDomain.StatusConfirmed,
		FreightHireFeeMinor:        freightHireFeeMinor,
		PlatformCommissionBps:      domain.PassThroughFreightCommissionBps,
		CommissionOwedMinor:        commissionOwedMinor,
		Currency:                   currency,
		AbnormalLoadEscortRequired: abnormalEscort,
		EscortFeeMinor:             escortFeeMinor,
		CarrierFullName:            carrierProfile.OwnerFullName,
		CarrierPhoneNumber:         carrierProfile.OwnerPhoneNumber,
	}, nil
}

type RegisterExternalTransportCommand struct {
	ParentBookingID       uuid.UUID
	CallerUserID          uuid.UUID
	CarrierName           string
	DriverFullName        string
	DriverPhoneNumber     string
	PrimeMoverPlateNumber string
	EstimatedArrivalAt    time.Time
}

// RegisterExternalTransport records private/off-platform transport details and unlocks READY_FOR_PICKUP (SRS §4.3).
func (s *LogisticsService) RegisterExternalTransport(ctx context.Context, cmd RegisterExternalTransportCommand) (*domain.ExternalTransport, bookingDomain.BookingStatus, error) {
	parent, err := s.bookingRepo.GetBookingByID(ctx, cmd.ParentBookingID)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", domain.ErrLogisticsBookingNotFound, err)
	}

	if parent.OwnerID != cmd.CallerUserID {
		return nil, "", domain.ErrUnauthorizedLogisticsAction
	}

	if parent.Status != bookingDomain.StatusConfirmed && parent.Status != bookingDomain.StatusAwaitingFleetReadiness && parent.Status != bookingDomain.StatusReadyForPickup {
		return nil, "", fmt.Errorf("%w: cannot register external transport in state %s", domain.ErrInvalidBookingStateForLogistics, parent.Status)
	}

	if dispatch, err := s.repo.GetLogisticsDispatchByParentBookingID(ctx, parent.ID); err == nil && dispatch != nil {
		sub, subErr := s.bookingRepo.GetBookingByID(ctx, dispatch.SubBookingID)
		if subErr == nil && sub.Status != bookingDomain.StatusCancelled && sub.Status != bookingDomain.StatusCancelledCarrierDefault {
			return nil, "", domain.ErrExternalTransportBlocked
		}
	}

	ext, err := s.repo.UpsertExternalTransport(ctx, &domain.ExternalTransport{
		ParentBookingID:       parent.ID,
		RegisteredBy:          cmd.CallerUserID,
		CarrierName:           cmd.CarrierName,
		DriverFullName:        cmd.DriverFullName,
		DriverPhoneNumber:     cmd.DriverPhoneNumber,
		PrimeMoverPlateNumber: cmd.PrimeMoverPlateNumber,
		EstimatedArrivalAt:    cmd.EstimatedArrivalAt.UTC(),
	})
	if err != nil {
		return nil, "", err
	}

	prevStatus := parent.Status
	if parent.Status != bookingDomain.StatusReadyForPickup {
		if err := bookingDomain.ValidateTransition(parent.Status, bookingDomain.StatusReadyForPickup); err != nil {
			return nil, "", err
		}
		parent.Status = bookingDomain.StatusReadyForPickup
		if err := s.bookingRepo.UpdateBookingStatusOCC(ctx, parent, parent.Version); err != nil {
			return nil, "", err
		}
	}

	_ = s.bookingRepo.CreateAuditLog(ctx, &bookingDomain.AuditLog{
		BookingID:  parent.ID,
		ActorID:    &cmd.CallerUserID,
		Action:     "EXTERNAL_TRANSPORT_REGISTERED",
		FromStatus: prevStatus,
		ToStatus:   parent.Status,
		Reason:     "Owner declared private/external transport mobilization",
		Payload: map[string]any{
			"driverName":  cmd.DriverFullName,
			"plateNumber": cmd.PrimeMoverPlateNumber,
		},
		CreatedAt: time.Now().UTC(),
	})

	return ext, parent.Status, nil
}

type AdvanceDispatchStageCommand struct {
	SubBookingID   uuid.UUID
	CallerUserID   uuid.UUID
	TargetStage    bookingDomain.BookingStatus
	WinchingPhotos []domain.WinchingPhoto
	Notes          string
}

type DispatchStageResult struct {
	SubBookingID        uuid.UUID
	ParentBookingID     uuid.UUID
	SubBookingStatus    bookingDomain.BookingStatus
	ParentBookingStatus bookingDomain.BookingStatus
	WinchingPhotosCount int
	UpdatedAt           time.Time
}

// AdvanceDispatchStage progresses a logistics sub-booking through DISPATCHED, LOADING_INSPECTION, IN_TRANSIT, and READY_FOR_PICKUP (SRS §5.1, §5.2).
func (s *LogisticsService) AdvanceDispatchStage(ctx context.Context, cmd AdvanceDispatchStageCommand) (*DispatchStageResult, error) {
	sub, err := s.bookingRepo.GetBookingByID(ctx, cmd.SubBookingID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrLogisticsBookingNotFound, err)
	}

	if sub.ParentBookingID == nil || sub.BookingKind != "LOGISTICS_SUB" {
		return nil, fmt.Errorf("%w: target booking is not a logistics sub-booking", domain.ErrInvalidDispatchStage)
	}

	dispatch, err := s.repo.GetLogisticsDispatchBySubBookingID(ctx, sub.ID)
	if err != nil {
		return nil, err
	}

	if sub.OwnerID != cmd.CallerUserID && dispatch.CarrierOwnerID != cmd.CallerUserID {
		return nil, domain.ErrUnauthorizedLogisticsAction
	}

	if err := bookingDomain.ValidateTransition(sub.Status, cmd.TargetStage); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidDispatchStage, err)
	}

	if cmd.TargetStage == bookingDomain.StatusLoadingInspection {
		if err := domain.ValidateWinchingPhotos(cmd.WinchingPhotos); err != nil {
			return nil, err
		}
		dispatch, err = s.repo.UpdateLogisticsDispatchLoading(ctx, sub.ID, cmd.WinchingPhotos)
		if err != nil {
			return nil, err
		}
	} else {
		dispatch, err = s.repo.UpdateLogisticsDispatchMilestone(ctx, sub.ID, string(cmd.TargetStage))
		if err != nil {
			return nil, err
		}
	}

	prevSubStatus := sub.Status
	sub.Status = cmd.TargetStage
	if err := s.bookingRepo.UpdateBookingStatusOCC(ctx, sub, sub.Version); err != nil {
		return nil, err
	}

	parent, err := s.bookingRepo.GetBookingByID(ctx, *sub.ParentBookingID)
	if err != nil {
		return nil, err
	}

	// Upon carrier arrival at contractor site (READY_FOR_PICKUP), unlock parent booking for handover
	if cmd.TargetStage == bookingDomain.StatusReadyForPickup &&
		(parent.Status == bookingDomain.StatusConfirmed || parent.Status == bookingDomain.StatusAwaitingFleetReadiness) {
		if err := bookingDomain.ValidateTransition(parent.Status, bookingDomain.StatusReadyForPickup); err == nil {
			parent.Status = bookingDomain.StatusReadyForPickup
			if err := s.bookingRepo.UpdateBookingStatusOCC(ctx, parent, parent.Version); err != nil {
				return nil, err
			}
		}
	}

	_ = s.bookingRepo.CreateAuditLog(ctx, &bookingDomain.AuditLog{
		BookingID:  sub.ID,
		ActorID:    &cmd.CallerUserID,
		Action:     "LOGISTICS_DISPATCH_STAGE",
		FromStatus: prevSubStatus,
		ToStatus:   sub.Status,
		Reason:     cmd.Notes,
		Payload: map[string]any{
			"parentBookingId": parent.ID.String(),
			"targetStage":     string(cmd.TargetStage),
		},
		CreatedAt: time.Now().UTC(),
	})

	return &DispatchStageResult{
		SubBookingID:        sub.ID,
		ParentBookingID:     parent.ID,
		SubBookingStatus:    sub.Status,
		ParentBookingStatus: parent.Status,
		WinchingPhotosCount: len(dispatch.WinchingPhotos),
		UpdatedAt:           dispatch.UpdatedAt,
	}, nil
}

type TransshipmentAlertCommand struct {
	SubBookingID       uuid.UUID
	CallerUserID       uuid.UUID
	BreakdownLocation  domain.GeoPoint
	FailureCategory    domain.TransshipmentFailureCategory
	FailureDescription string
	Photos             []string
}

type TransshipmentAlertResult struct {
	Transshipment            *domain.LogisticsTransshipment
	SubBookingStatus         bookingDomain.BookingStatus
	ParentBookingStatus      bookingDomain.BookingStatus
	CandidateCarriersCount   int
	CandidateCarriers        []domain.MatchedCarrier
}

// ReportTransshipmentAlert freezes the parent booking, opens a 60-minute hot-swap window, and searches within 50km (SRS §5.3).
func (s *LogisticsService) ReportTransshipmentAlert(ctx context.Context, cmd TransshipmentAlertCommand) (*TransshipmentAlertResult, error) {
	sub, err := s.bookingRepo.GetBookingByID(ctx, cmd.SubBookingID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrLogisticsBookingNotFound, err)
	}
	if sub.ParentBookingID == nil {
		return nil, fmt.Errorf("%w: booking has no parent booking", domain.ErrInvalidDispatchStage)
	}
	if sub.OwnerID != cmd.CallerUserID && sub.RenterID != cmd.CallerUserID {
		return nil, domain.ErrUnauthorizedLogisticsAction
	}

	if err := bookingDomain.ValidateTransition(sub.Status, bookingDomain.StatusTransshipmentRequired); err != nil {
		return nil, err
	}

	parent, err := s.bookingRepo.GetBookingByID(ctx, *sub.ParentBookingID)
	if err != nil {
		return nil, err
	}

	sub.Status = bookingDomain.StatusTransshipmentRequired
	if err := s.bookingRepo.UpdateBookingStatusOCC(ctx, sub, sub.Version); err != nil {
		return nil, err
	}

	if parent.Status == bookingDomain.StatusConfirmed || parent.Status == bookingDomain.StatusReadyForPickup {
		parent.Status = bookingDomain.StatusAwaitingFleetReadiness
		if err := s.bookingRepo.UpdateBookingStatusOCC(ctx, parent, parent.Version); err != nil {
			return nil, err
		}
	}

	bounty := domain.CalculateTransshipmentBounty(sub.BaseHireFeeMinor)
	deadline := time.Now().UTC().Add(60 * time.Minute)
	trans, err := s.repo.CreateLogisticsTransshipment(ctx, &domain.LogisticsTransshipment{
		SubBookingID:             sub.ID,
		ParentBookingID:          parent.ID,
		ReportedBy:               cmd.CallerUserID,
		BreakdownLocation:        cmd.BreakdownLocation,
		FailureCategory:          cmd.FailureCategory,
		FailureDescription:       cmd.FailureDescription,
		Photos:                   cmd.Photos,
		HotSwapSLADeadline:       deadline,
		TransshipmentBountyMinor: bounty,
		Currency:                 sub.Currency,
		Status:                   domain.TransshipmentOpen,
	})
	if err != nil {
		return nil, err
	}

	machineProfile, err := s.repo.GetFleetLogisticsProfile(ctx, parent.FleetID)
	var candidates []domain.MatchedCarrier
	if err == nil {
		limits, _ := s.repo.GetCountryLogisticsLimits(ctx, machineProfile.CountryCode)
		candidates, _ = s.repo.MatchCompatibleLogisticsFleets(
			ctx,
			cmd.BreakdownLocation,
			50000.0,
			machineProfile.CountryCode,
			sub.OwnerID,
			machineProfile.MachineSpecs,
			limits,
			10,
		)
	}

	return &TransshipmentAlertResult{
		Transshipment:          trans,
		SubBookingStatus:       sub.Status,
		ParentBookingStatus:    parent.Status,
		CandidateCarriersCount: len(candidates),
		CandidateCarriers:      candidates,
	}, nil
}

type AcceptTransshipmentCommand struct {
	SubBookingID       uuid.UUID
	CallerUserID       uuid.UUID
	ReplacementFleetID uuid.UUID
}

// AcceptTransshipment assigns a compatible replacement lowboy within the 60-minute SLA window and resumes IN_TRANSIT (SRS §5.3).
func (s *LogisticsService) AcceptTransshipment(ctx context.Context, cmd AcceptTransshipmentCommand) (*domain.LogisticsTransshipment, bookingDomain.BookingStatus, error) {
	sub, err := s.bookingRepo.GetBookingByID(ctx, cmd.SubBookingID)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", domain.ErrLogisticsBookingNotFound, err)
	}
	if sub.ParentBookingID == nil {
		return nil, "", domain.ErrLogisticsBookingNotFound
	}

	openTrans, err := s.repo.GetOpenTransshipmentBySubBookingID(ctx, sub.ID)
	if err != nil {
		return nil, "", err
	}

	if time.Now().UTC().After(openTrans.HotSwapSLADeadline) {
		_ = s.repo.ExpireLogisticsTransshipment(ctx, openTrans.ID)
		return nil, "", domain.ErrTransshipmentWindowExpired
	}

	parent, err := s.bookingRepo.GetBookingByID(ctx, *sub.ParentBookingID)
	if err != nil {
		return nil, "", err
	}
	machineProfile, err := s.repo.GetFleetLogisticsProfile(ctx, parent.FleetID)
	if err != nil {
		return nil, "", err
	}
	replacementProfile, err := s.repo.GetFleetLogisticsProfile(ctx, cmd.ReplacementFleetID)
	if err != nil {
		return nil, "", err
	}
	limits, err := s.repo.GetCountryLogisticsLimits(ctx, machineProfile.CountryCode)
	if err != nil {
		return nil, "", err
	}

	if _, err := domain.EvaluateDimensionalCompatibility(machineProfile.MachineSpecs, replacementProfile.TrailerProfile, limits); err != nil {
		return nil, "", err
	}

	accepted, err := s.repo.AcceptLogisticsTransshipment(ctx, openTrans.ID, replacementProfile.FleetID, replacementProfile.OwnerID)
	if err != nil {
		return nil, "", err
	}

	sub.Status = bookingDomain.StatusInTransit
	if err := s.bookingRepo.UpdateBookingStatusOCC(ctx, sub, sub.Version); err != nil {
		return nil, "", err
	}

	return accepted, sub.Status, nil
}

type ReverseCascadeCommand struct {
	ParentBookingID uuid.UUID
	CallerUserID    uuid.UUID
	Reason          string
}

type ReverseCascadeResult struct {
	Cascade      *domain.LogisticsReverseCascade
	ParentStatus bookingDomain.BookingStatus
	SubStatus    bookingDomain.BookingStatus
}

// TriggerReverseCascade executes Force Majeure exoneration, 100% renter refund, 0% owner penalty, commission restitution, and carrier rating deduction (SRS §6, FLOW-18).
func (s *LogisticsService) TriggerReverseCascade(ctx context.Context, cmd ReverseCascadeCommand) (*ReverseCascadeResult, error) {
	parent, err := s.bookingRepo.GetBookingByID(ctx, cmd.ParentBookingID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrLogisticsBookingNotFound, err)
	}

	if existing, err := s.repo.GetLogisticsReverseCascadeByParentBookingID(ctx, parent.ID); err == nil && existing != nil {
		return nil, domain.ErrReverseCascadeAlreadyExecuted
	}

	dispatch, err := s.repo.GetLogisticsDispatchByParentBookingID(ctx, parent.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: no linked logistics sub-booking found for parent booking", domain.ErrLogisticsBookingNotFound)
	}

	sub, err := s.bookingRepo.GetBookingByID(ctx, dispatch.SubBookingID)
	if err != nil {
		return nil, err
	}

	if err := bookingDomain.ValidateTransition(sub.Status, bookingDomain.StatusCancelledCarrierDefault); err != nil {
		return nil, err
	}
	if err := bookingDomain.ValidateTransition(parent.Status, bookingDomain.StatusCancelled); err != nil {
		return nil, err
	}

	reason := strings.TrimSpace(cmd.Reason)
	if reason == "" {
		reason = "Force Majeure reverse failure cascade triggered due to lowboy carrier default / transshipment SLA failure"
	}

	parentCommission := parent.PlatformCommissionMinor
	if parentCommission == 0 {
		parentCommission = (parent.GrossRentalAmountMinor * 1000) / 10000
	}

	financials := domain.CalculateReverseCascadeFinancials(
		parent.PrepaymentAmountMinor,
		parentCommission,
		dispatch.FreightHireFeeMinor,
	)

	// 1. Cancel Sub-Booking with CANCELLED_CARRIER_DEFAULT
	sub.Status = bookingDomain.StatusCancelledCarrierDefault
	sub.CancellationReason = &reason
	sub.CancelledBy = &cmd.CallerUserID
	sub.PenaltyAmountMinor = financials.CarrierDefaultPenaltyMinor
	if err := s.bookingRepo.UpdateBookingStatusOCC(ctx, sub, sub.Version); err != nil {
		return nil, err
	}

	// 2. Cancel Parent Booking under Force Majeure (100% renter refund, 0% owner penalty)
	parent.Status = bookingDomain.StatusCancelled
	parent.CancellationReason = &reason
	parent.CancelledBy = &cmd.CallerUserID
	parent.PenaltyAmountMinor = financials.OwnerPenaltyLeviedMinor
	parent.RefundAmountMinor = financials.RenterRefundDueMinor
	if err := s.bookingRepo.UpdateBookingStatusOCC(ctx, parent, parent.Version); err != nil {
		return nil, err
	}

	// 3. Persist Reverse Cascade Ledger & Release GiST Availability + Deduct Carrier Rating
	cascade, err := s.repo.CreateLogisticsReverseCascade(ctx, &domain.LogisticsReverseCascade{
		ParentBookingID:                 parent.ID,
		SubBookingID:                    sub.ID,
		TriggeredBy:                     cmd.CallerUserID,
		Resolution:                      "FORCE_MAJEURE_CANCELLATION",
		Reason:                          reason,
		RenterRefundDueMinor:            financials.RenterRefundDueMinor,
		OwnerPenaltyLeviedMinor:         financials.OwnerPenaltyLeviedMinor,
		OwnerCommissionRestitutionMinor: financials.OwnerCommissionRestitutionMinor,
		CarrierDefaultPenaltyMinor:      financials.CarrierDefaultPenaltyMinor,
		CarrierRatingDeduction:          financials.CarrierRatingDeduction,
		Currency:                        parent.Currency,
	})
	if err != nil {
		return nil, err
	}

	_ = s.repo.ReleaseAvailabilityAndDeductCarrierRating(ctx, parent.ID, sub.ID, dispatch.CarrierOwnerID, financials.CarrierRatingDeduction)

	_ = s.bookingRepo.CreateAuditLog(ctx, &bookingDomain.AuditLog{
		BookingID:  parent.ID,
		ActorID:    &cmd.CallerUserID,
		Action:     "FORCE_MAJEURE_REVERSE_CASCADE",
		FromStatus: bookingDomain.StatusAwaitingFleetReadiness,
		ToStatus:   bookingDomain.StatusCancelled,
		Reason:     reason,
		Payload: map[string]any{
			"subBookingId":                   sub.ID.String(),
			"renterRefundDueMinor":           financials.RenterRefundDueMinor,
			"ownerPenaltyLeviedMinor":        0,
			"ownerCommissionRestitutionMinor": financials.OwnerCommissionRestitutionMinor,
		},
		CreatedAt: time.Now().UTC(),
	})

	return &ReverseCascadeResult{
		Cascade:      cascade,
		ParentStatus: parent.Status,
		SubStatus:    sub.Status,
	}, nil
}

type LogisticsOverviewResult struct {
	ParentBookingID     uuid.UUID
	ParentBookingStatus bookingDomain.BookingStatus
	SelfMoveCapable     bool
	LogisticsMode       string
	Dispatch            *domain.LogisticsDispatch
	SubBooking          *bookingDomain.Booking
	CarrierContact      *domain.FleetLogisticsProfile
	ExternalTransport   *domain.ExternalTransport
	ActiveTransshipment *domain.LogisticsTransshipment
	ReverseCascade      *domain.LogisticsReverseCascade
}

// GetLogisticsOverview returns the current mobilization topology for a parent booking.
func (s *LogisticsService) GetLogisticsOverview(ctx context.Context, parentBookingID, callerUserID uuid.UUID) (*LogisticsOverviewResult, error) {
	parent, err := s.bookingRepo.GetBookingByID(ctx, parentBookingID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrLogisticsBookingNotFound, err)
	}

	profile, err := s.repo.GetFleetLogisticsProfile(ctx, parent.FleetID)
	if err != nil {
		return nil, err
	}

	res := &LogisticsOverviewResult{
		ParentBookingID:     parent.ID,
		ParentBookingStatus: parent.Status,
		SelfMoveCapable:     profile.SelfMoveCapable,
		LogisticsMode:       "UNASSIGNED",
	}

	if profile.SelfMoveCapable {
		res.LogisticsMode = "SELF_MOVE"
		return res, nil
	}

	dispatch, dErr := s.repo.GetLogisticsDispatchByParentBookingID(ctx, parent.ID)
	if dErr == nil && dispatch != nil {
		res.LogisticsMode = "PLATFORM_SUB_BOOKING"
		res.Dispatch = dispatch
		if sub, err := s.bookingRepo.GetBookingByID(ctx, dispatch.SubBookingID); err == nil {
			res.SubBooking = sub
		}
		if carrierProf, err := s.repo.GetFleetLogisticsProfile(ctx, dispatch.LogisticsFleetID); err == nil {
			res.CarrierContact = carrierProf
		}
		if trans, err := s.repo.GetOpenTransshipmentBySubBookingID(ctx, dispatch.SubBookingID); err == nil {
			res.ActiveTransshipment = trans
		}
		if cascade, err := s.repo.GetLogisticsReverseCascadeByParentBookingID(ctx, parent.ID); err == nil {
			res.ReverseCascade = cascade
		}

		// Customer isolation: ensure caller is participant (renter, owner, or carrier)
		isCarrier := dispatch.CarrierOwnerID == callerUserID
		if parent.OwnerID != callerUserID && parent.RenterID != callerUserID && !isCarrier {
			return nil, domain.ErrUnauthorizedLogisticsAction
		}
		return res, nil
	}

	if parent.OwnerID != callerUserID && parent.RenterID != callerUserID {
		return nil, domain.ErrUnauthorizedLogisticsAction
	}

	if ext, err := s.repo.GetExternalTransportByParentBookingID(ctx, parent.ID); err == nil && ext != nil {
		res.LogisticsMode = "EXTERNAL_TRANSPORT"
		res.ExternalTransport = ext
	}

	return res, nil
}

// EnsureParentReadyForActivation enforces D-05-08-001: parent machinery cannot enter ACTIVE while linked logistics sub-booking is still in transit.
func (s *LogisticsService) EnsureParentReadyForActivation(ctx context.Context, parentBookingID uuid.UUID) error {
	dispatch, err := s.repo.GetLogisticsDispatchByParentBookingID(ctx, parentBookingID)
	if err != nil || dispatch == nil {
		return nil // No in-platform sub-booking; external or self-move
	}

	sub, err := s.bookingRepo.GetBookingByID(ctx, dispatch.SubBookingID)
	if err != nil {
		return nil
	}

	if sub.Status != bookingDomain.StatusReadyForPickup && sub.Status != bookingDomain.StatusCompleted && sub.Status != bookingDomain.StatusActive {
		return fmt.Errorf("%w: sub-booking %s is in state %s", domain.ErrParentBookingBlockedUntilDelivery, sub.ID, sub.Status)
	}
	return nil
}
