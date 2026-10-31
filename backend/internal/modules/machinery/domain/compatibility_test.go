package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
)

func TestCompatibility_RejectsUndersizedHost(t *testing.T) {
	// 20-ton to 26-ton hydraulic breaker (e.g. Soosan SB-81)
	breaker := &domain.Attachment{
		ID:                     uuid.New(),
		AttachmentType:         domain.AttachmentTypeHydraulicBreaker,
		Make:                   "Soosan",
		Model:                  "SB-81",
		MinCarrierWeightTonnes: 18.0,
		MaxCarrierWeightTonnes: 26.0,
		PinDiameterMm:          80.0,
		RequiredFlowMinLmin:    120.0,
		RequiredFlowMaxLmin:    150.0,
		RequiredPressureBar:    180.0,
		StandaloneDailyRateMinor: 1200000,
	}

	// 1. Undersized host (8.0-ton mini/midi excavator)
	err := breaker.ValidateCarrierCompatibility(8.0, 180.0, 130.0, 80.0)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrUndersizedCarrier)

	// 2. Oversized host (45.0-ton quarry excavator)
	err = breaker.ValidateCarrierCompatibility(45.0, 180.0, 130.0, 80.0)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrOversizedCarrier)

	// 3. Mismatched pin diameter (65mm pin on an 80mm breaker ear)
	err = breaker.ValidateCarrierCompatibility(22.0, 180.0, 130.0, 65.0)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrMismatchedPinDiameter)

	// 4. Insufficient hydraulic flow (70 L/min vs 120 L/min min)
	err = breaker.ValidateCarrierCompatibility(22.0, 180.0, 70.0, 80.0)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInsufficientHydraulicFlow)

	// 5. Insufficient hydraulic relief pressure (140 bar vs 180 bar min)
	err = breaker.ValidateCarrierCompatibility(22.0, 140.0, 130.0, 80.0)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInsufficientHydraulicPressure)

	// 6. Compatible 22.0-ton carrier
	err = breaker.ValidateCarrierCompatibility(22.0, 190.0, 135.0, 80.0)
	require.NoError(t, err)
}

func TestAttachment_ValidateBundledRate(t *testing.T) {
	breaker := &domain.Attachment{
		ID:                       uuid.New(),
		StandaloneDailyRateMinor: 1200000, // 12,000 ETB
	}

	// Discounted bundled rate (8,000 ETB) -> valid
	err := breaker.ValidateBundledRate(800000)
	assert.NoError(t, err)

	// Equal rate (12,000 ETB) -> valid
	err = breaker.ValidateBundledRate(1200000)
	assert.NoError(t, err)

	// Exceeds standalone rate (15,000 ETB) -> rejected
	err = breaker.ValidateBundledRate(1500000)
	assert.ErrorIs(t, err, domain.ErrBundledRateExceedsStandalone)

	// Negative rate -> rejected
	err = breaker.ValidateBundledRate(-500)
	assert.Error(t, err)
}

func TestAttachment_ValidateRequiredPhotos(t *testing.T) {
	validPhotos := []domain.AttachmentPhoto{
		{Angle: domain.AttachmentPhotoAngleOverall, StorageURI: "s3://media/overall.webp"},
		{Angle: domain.AttachmentPhotoAngleBracket, StorageURI: "s3://media/bracket.webp"},
		{Angle: domain.AttachmentPhotoAnglePins, StorageURI: "s3://media/pins.webp"},
		{Angle: domain.AttachmentPhotoAngleToolTip, StorageURI: "s3://media/tooltip.webp"},
	}

	err := domain.ValidateAttachmentPhotos(validPhotos)
	assert.NoError(t, err)

	// Missing TOOL_TIP
	missingPhotos := []domain.AttachmentPhoto{
		{Angle: domain.AttachmentPhotoAngleOverall, StorageURI: "s3://media/overall.webp"},
		{Angle: domain.AttachmentPhotoAngleBracket, StorageURI: "s3://media/bracket.webp"},
		{Angle: domain.AttachmentPhotoAnglePins, StorageURI: "s3://media/pins.webp"},
	}

	err = domain.ValidateAttachmentPhotos(missingPhotos)
	assert.ErrorIs(t, err, domain.ErrIncompleteAttachmentPhotos)
}
