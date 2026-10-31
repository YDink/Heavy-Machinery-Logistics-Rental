package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
)

type RegisterAttachmentInput struct {
	OwnerID                     uuid.UUID
	AttachmentType              domain.AttachmentType
	Make                        string
	Model                       string
	SubTypeSize                 string
	SerialNumber                string
	OperatingWeightKg           float64
	MinCarrierWeightTonnes      float64
	MaxCarrierWeightTonnes      float64
	RequiredFlowMinLmin         float64
	RequiredFlowMaxLmin         float64
	RequiredPressureBar         float64
	PinDiameterMm               float64
	PinSpreadMm                 float64
	CompatibleSubcategories     []string
	StandaloneDailyRateMinor    int64
	SecurityDepositMinor        int64
	MinRentalDays               int
	PaymentEvidenceWindowHours  int
	Photos                      []domain.AttachmentPhoto
}

type LinkAttachmentInput struct {
	FleetID               uuid.UUID
	AttachmentID          uuid.UUID
	BundledDailyRateMinor int64
}

type CheckBundledAvailabilityInput struct {
	FleetID       uuid.UUID
	AttachmentIDs []uuid.UUID
	StartDate     time.Time
	EndDate       time.Time
}

type AttachmentAvailabilityItem struct {
	AttachmentID          uuid.UUID
	IsAvailable           bool
	DailyRateAppliedMinor int64
}

type CheckBundledAvailabilityOutput struct {
	FleetAvailable              bool
	AllAttachmentsAvailable     bool
	AttachmentStatuses          []AttachmentAvailabilityItem
	TotalCombinedDailyRateMinor int64
}

type AttachmentService struct {
	attachmentRepo domain.AttachmentRepository
	fleetRepo      domain.FleetRepository
}

func NewAttachmentService(
	attachmentRepo domain.AttachmentRepository,
	fleetRepo domain.FleetRepository,
) *AttachmentService {
	return &AttachmentService{
		attachmentRepo: attachmentRepo,
		fleetRepo:      fleetRepo,
	}
}

func (s *AttachmentService) RegisterAttachment(ctx context.Context, input RegisterAttachmentInput) (*domain.Attachment, error) {
	if input.Make == "" || input.Model == "" || input.SerialNumber == "" {
		return nil, errors.New("make, model, and serial number are required")
	}

	if input.MinCarrierWeightTonnes > input.MaxCarrierWeightTonnes {
		return nil, errors.New("minimum carrier weight cannot exceed maximum carrier weight")
	}

	if input.RequiredFlowMinLmin > input.RequiredFlowMaxLmin {
		return nil, errors.New("minimum hydraulic flow cannot exceed maximum hydraulic flow")
	}

	if err := domain.ValidateAttachmentPhotos(input.Photos); err != nil {
		return nil, err
	}

	attachment := &domain.Attachment{
		ID:                         uuid.New(),
		OwnerID:                    input.OwnerID,
		AttachmentType:             input.AttachmentType,
		Make:                       input.Make,
		Model:                      input.Model,
		SubTypeSize:                input.SubTypeSize,
		SerialNumber:               input.SerialNumber,
		OperatingWeightKg:          input.OperatingWeightKg,
		MinCarrierWeightTonnes:     input.MinCarrierWeightTonnes,
		MaxCarrierWeightTonnes:     input.MaxCarrierWeightTonnes,
		RequiredFlowMinLmin:        input.RequiredFlowMinLmin,
		RequiredFlowMaxLmin:        input.RequiredFlowMaxLmin,
		RequiredPressureBar:        input.RequiredPressureBar,
		PinDiameterMm:              input.PinDiameterMm,
		PinSpreadMm:                input.PinSpreadMm,
		StandaloneDailyRateMinor:   input.StandaloneDailyRateMinor,
		SecurityDepositMinor:       input.SecurityDepositMinor,
		MinRentalDays:              input.MinRentalDays,
		PaymentEvidenceWindowHours: input.PaymentEvidenceWindowHours,
		Status:                     domain.AttachmentStatusActive,
		CompatibleCategories:       input.CompatibleSubcategories,
		Photos:                     input.Photos,
		CreatedAt:                  time.Now().UTC(),
		UpdatedAt:                  time.Now().UTC(),
	}

	return s.attachmentRepo.CreateAttachment(ctx, attachment, input.CompatibleSubcategories, input.Photos)
}

func (s *AttachmentService) GetAttachment(ctx context.Context, id uuid.UUID) (*domain.Attachment, error) {
	return s.attachmentRepo.GetAttachmentByID(ctx, id)
}

func (s *AttachmentService) ListOwnerInventory(ctx context.Context, ownerID uuid.UUID, status string) ([]*domain.Attachment, error) {
	return s.attachmentRepo.ListAttachmentsByOwner(ctx, ownerID, status)
}

func (s *AttachmentService) LinkAttachmentToFleet(ctx context.Context, ownerID uuid.UUID, input LinkAttachmentInput) (*domain.FleetBundledAttachment, error) {
	attachment, err := s.attachmentRepo.GetAttachmentByID(ctx, input.AttachmentID)
	if err != nil {
		return nil, err
	}

	if attachment.OwnerID != ownerID {
		return nil, domain.ErrUnauthorizedAttachmentOwner
	}

	fleet, err := s.fleetRepo.GetFleetByID(ctx, input.FleetID)
	if err != nil {
		return nil, err
	}

	if fleet.OwnerID != ownerID {
		return nil, domain.ErrUnauthorizedFleetOwner
	}

	// Validate bundled rate invariant: bundled <= standalone (FL-33)
	if err := attachment.ValidateBundledRate(input.BundledDailyRateMinor); err != nil {
		return nil, err
	}

	return s.attachmentRepo.LinkAttachmentToFleet(ctx, input.FleetID, input.AttachmentID, input.BundledDailyRateMinor)
}

func (s *AttachmentService) ListBundledAttachmentsForFleet(ctx context.Context, fleetID uuid.UUID) ([]*domain.FleetBundledAttachment, error) {
	return s.attachmentRepo.ListBundledAttachmentsByFleet(ctx, fleetID)
}

func (s *AttachmentService) CheckBundledAvailability(ctx context.Context, input CheckBundledAvailabilityInput) (*CheckBundledAvailabilityOutput, error) {
	if !input.EndDate.After(input.StartDate) {
		return nil, errors.New("end date must be after start date")
	}

	fleet, err := s.fleetRepo.GetFleetByID(ctx, input.FleetID)
	if err != nil {
		return nil, err
	}

	fleetAvailable := fleet.Status == domain.FleetStatusActive

	allAttachmentsAvailable := true
	var attachmentStatuses []AttachmentAvailabilityItem
	var totalCombinedDailyRateMinor int64

	if fleetAvailable {
		totalCombinedDailyRateMinor += fleet.DailyRateMinor
	}

	bundledMap := make(map[uuid.UUID]*domain.FleetBundledAttachment)
	bundledList, err := s.attachmentRepo.ListBundledAttachmentsByFleet(ctx, input.FleetID)
	if err == nil {
		for _, b := range bundledList {
			bundledMap[b.AttachmentID] = b
		}
	}

	for _, attID := range input.AttachmentIDs {
		isAvailable, err := s.attachmentRepo.CheckAttachmentAvailability(ctx, attID, input.StartDate, input.EndDate)
		if err != nil || !isAvailable {
			isAvailable = false
			allAttachmentsAvailable = false
		}

		rateMinor := int64(0)
		if bundled, ok := bundledMap[attID]; ok {
			rateMinor = bundled.BundledDailyRateMinor
		} else {
			att, err := s.attachmentRepo.GetAttachmentByID(ctx, attID)
			if err == nil && att != nil {
				rateMinor = att.StandaloneDailyRateMinor
			}
		}

		if isAvailable {
			totalCombinedDailyRateMinor += rateMinor
		}

		attachmentStatuses = append(attachmentStatuses, AttachmentAvailabilityItem{
			AttachmentID:          attID,
			IsAvailable:           isAvailable,
			DailyRateAppliedMinor: rateMinor,
		})
	}

	return &CheckBundledAvailabilityOutput{
		FleetAvailable:              fleetAvailable,
		AllAttachmentsAvailable:     allAttachmentsAvailable,
		AttachmentStatuses:          attachmentStatuses,
		TotalCombinedDailyRateMinor: totalCombinedDailyRateMinor,
	}, nil
}

// ReserveAttachment creates a temporal lease hold for an attachment (D-05-03-001).
func (s *AttachmentService) ReserveAttachment(ctx context.Context, attachmentID uuid.UUID, referenceID *uuid.UUID, startDate, endDate time.Time, dailyRateMinor int64) error {
	isAvail, err := s.attachmentRepo.CheckAttachmentAvailability(ctx, attachmentID, startDate, endDate)
	if err != nil {
		return err
	}
	if !isAvail {
		return domain.ErrAttachmentScheduleConflict
	}

	res := &domain.AttachmentReservation{
		ID:                    uuid.New(),
		AttachmentID:          attachmentID,
		ReferenceID:           referenceID,
		StartDate:             startDate,
		EndDate:               endDate,
		DailyRateAppliedMinor: dailyRateMinor,
		Status:                "HOLD_PENDING",
		CreatedAt:             time.Now().UTC(),
	}

	return s.attachmentRepo.CreateReservation(ctx, res)
}
