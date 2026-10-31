package domain

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
)

type AttachmentType string

const (
	AttachmentTypeHydraulicBreaker AttachmentType = "HYDRAULIC_BREAKER"
	AttachmentTypeStandardBucket   AttachmentType = "STANDARD_BUCKET"
	AttachmentTypeTrenchingBucket  AttachmentType = "TRENCHING_BUCKET"
	AttachmentTypeAuger            AttachmentType = "AUGER"
	AttachmentTypeRipperTooth      AttachmentType = "RIPPER_TOOTH"
	AttachmentTypeSortingGrapple   AttachmentType = "SORTING_GRAPPLE"
	AttachmentTypeCrusher          AttachmentType = "CRUSHER"
)

type AttachmentStatus string

const (
	AttachmentStatusDraft           AttachmentStatus = "DRAFT"
	AttachmentStatusActive          AttachmentStatus = "ACTIVE"
	AttachmentStatusRented          AttachmentStatus = "RENTED"
	AttachmentStatusMaintenanceHold AttachmentStatus = "MAINTENANCE_HOLD"
	AttachmentStatusInactive        AttachmentStatus = "INACTIVE"
	AttachmentStatusDisputedDamage  AttachmentStatus = "DISPUTED_DAMAGE"
	AttachmentStatusDecommissioned  AttachmentStatus = "DECOMMISSIONED"
)

type AttachmentPhotoAngle string

const (
	AttachmentPhotoAngleOverall AttachmentPhotoAngle = "OVERALL"
	AttachmentPhotoAngleBracket AttachmentPhotoAngle = "BRACKET"
	AttachmentPhotoAnglePins    AttachmentPhotoAngle = "PINS"
	AttachmentPhotoAngleToolTip AttachmentPhotoAngle = "TOOL_TIP"
)

var MandatoryAttachmentPhotoAngles = []AttachmentPhotoAngle{
	AttachmentPhotoAngleOverall,
	AttachmentPhotoAngleBracket,
	AttachmentPhotoAnglePins,
	AttachmentPhotoAngleToolTip,
}

var (
	ErrIncompleteAttachmentPhotos    = errors.New("attachment registration mandates 4 required physical angles: OVERALL, BRACKET, PINS, TOOL_TIP")
	ErrUndersizedCarrier             = errors.New("host carrier weight class is below minimum required tool threshold")
	ErrOversizedCarrier              = errors.New("host carrier weight class exceeds maximum permitted tool threshold")
	ErrMismatchedPinDiameter         = errors.New("mounting pin diameter does not match host carrier coupler")
	ErrInsufficientHydraulicFlow     = errors.New("host machine auxiliary hydraulic flow is insufficient for attachment")
	ErrInsufficientHydraulicPressure = errors.New("host machine operating relief pressure is insufficient for attachment")
	ErrBundledRateExceedsStandalone  = errors.New("bundled daily rate cannot exceed standalone inventory daily rate")
	ErrAttachmentScheduleConflict    = errors.New("attachment is already reserved for overlapping dates")
	ErrAttachmentNotFound            = errors.New("attachment not found")
	ErrUnauthorizedAttachmentOwner   = errors.New("caller is not the authorized owner of this attachment")
)

type Attachment struct {
	ID                          uuid.UUID
	OwnerID                     uuid.UUID
	AttachmentType              AttachmentType
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
	StandaloneDailyRateMinor    int64
	SecurityDepositMinor        int64
	MinRentalDays               int
	PaymentEvidenceWindowHours  int
	Status                      AttachmentStatus
	CompatibleCategories        []string
	Photos                      []AttachmentPhoto
	CreatedAt                   time.Time
	UpdatedAt                   time.Time
}

type AttachmentPhoto struct {
	ID           uuid.UUID
	AttachmentID uuid.UUID
	Angle        AttachmentPhotoAngle
	StorageURI   string
	CreatedAt    time.Time
}

type FleetBundledAttachment struct {
	ID                    uuid.UUID
	FleetID               uuid.UUID
	AttachmentID          uuid.UUID
	BundledDailyRateMinor int64
	IsActive              bool
	AttachmentType        AttachmentType
	Make                  string
	Model                 string
	SubTypeSize           string
	SerialNumber          string
	StandaloneDailyRate   int64
	Status                AttachmentStatus
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type AttachmentReservation struct {
	ID                    uuid.UUID
	AttachmentID          uuid.UUID
	ReferenceID           *uuid.UUID
	StartDate             time.Time
	EndDate               time.Time
	DailyRateAppliedMinor int64
	Status                string
	CreatedAt             time.Time
}

// ValidateAttachmentPhotos validates all 4 mandatory photos are provided (ATT-07).
func ValidateAttachmentPhotos(photos []AttachmentPhoto) error {
	seen := make(map[AttachmentPhotoAngle]bool)
	for _, p := range photos {
		seen[p.Angle] = true
	}
	for _, angle := range MandatoryAttachmentPhotoAngles {
		if !seen[angle] {
			return ErrIncompleteAttachmentPhotos
		}
	}
	return nil
}

// ValidateCarrierCompatibility enforces mechanical & hydraulic compatibility rules (D-05-03-002, ATT-05, ATT-06).
func (a *Attachment) ValidateCarrierCompatibility(carrierWeightTonnes, carrierPressureBar, carrierFlowLmin, pinDiameterMm float64) error {
	// Carrier Weight Class validation (tolerance: +/- 10%)
	weightToleranceMin := a.MinCarrierWeightTonnes * 0.90
	weightToleranceMax := a.MaxCarrierWeightTonnes * 1.10

	if carrierWeightTonnes < weightToleranceMin {
		return ErrUndersizedCarrier
	}
	if carrierWeightTonnes > weightToleranceMax {
		return ErrOversizedCarrier
	}

	// Pin diameter verification (tolerance +/- 2.0 mm)
	if pinDiameterMm > 0 && math.Abs(a.PinDiameterMm-pinDiameterMm) > 2.0 {
		return ErrMismatchedPinDiameter
	}

	// Hydraulic flow verification (must meet required minimum flow)
	if carrierFlowLmin > 0 && carrierFlowLmin < a.RequiredFlowMinLmin {
		return ErrInsufficientHydraulicFlow
	}

	// Hydraulic pressure verification (must meet required operating pressure)
	if carrierPressureBar > 0 && carrierPressureBar < a.RequiredPressureBar {
		return ErrInsufficientHydraulicPressure
	}

	return nil
}

// ValidateBundledRate enforces that bundled rate cannot exceed standalone inventory daily rate (FL-33).
func (a *Attachment) ValidateBundledRate(bundledDailyRateMinor int64) error {
	if bundledDailyRateMinor < 0 {
		return errors.New("bundled daily rate cannot be negative")
	}
	if bundledDailyRateMinor > a.StandaloneDailyRateMinor {
		return ErrBundledRateExceedsStandalone
	}
	return nil
}
