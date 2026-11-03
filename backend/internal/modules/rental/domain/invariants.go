package domain

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"
)

const (
	StandardShiftHoursPerDay = 8.00
	HandoverChallengeMaxAttempts = 3
	HandoverChallengeTTL = 15 * time.Minute
	BreakdownOwnerSLADuration = 4 * time.Hour
)

var ValidFaultCategories = map[string]bool{
	"ENGINE":            true,
	"HYDRAULIC_CIRCUIT": true,
	"ELECTRICAL":        true,
	"UNDERCARRIAGE":     true,
	"STRUCTURAL":        true,
	"BOOM_BUCKET":       true,
	"OTHER":             true,
}

var ValidInspectionAngles = map[string]bool{
	"METER_CLOSEUP":   true,
	"FRONT_ELEVATION": true,
	"REAR_ELEVATION":  true,
	"LEFT_TRACK":      true,
	"RIGHT_TRACK":     true,
	"CABIN_INTERIOR":  true,
}

// ValidateMonotonicMeterFloor asserts that the check-in meter reading cannot regress below the machine certified baseline.
func ValidateMonotonicMeterFloor(checkInMeter, baselineMeter float64) error {
	if checkInMeter < baselineMeter {
		return fmt.Errorf("%w: check-in reading (%.2f hrs) is lower than certified baseline (%.2f hrs)",
			ErrMeterReadingRegression, checkInMeter, baselineMeter)
	}
	return nil
}

// ValidateHandoverPhotos asserts that inspection photos are provided and use valid angles.
func ValidateHandoverPhotos(photos []InspectionPhoto) error {
	if len(photos) == 0 {
		return fmt.Errorf("at least 1 reticle inspection photo is required")
	}
	for _, p := range photos {
		if !ValidInspectionAngles[strings.ToUpper(p.Angle)] {
			return fmt.Errorf("invalid inspection photo angle: %s", p.Angle)
		}
		if strings.TrimSpace(p.StorageURI) == "" {
			return fmt.Errorf("inspection photo storage URI cannot be empty")
		}
	}
	return nil
}

// GenerateChallengePIN generates a cryptographically random 6-digit numeric PIN (e.g. "819204").
func GenerateChallengePIN() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", fmt.Errorf("failed to generate random PIN: %w", err)
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// CalculateShiftPacing returns the operational pacing ratio: accumulatedToday / 8.0.
func CalculateShiftPacing(accumulatedToday float64) float64 {
	if accumulatedToday <= 0 {
		return 0.0
	}
	return accumulatedToday / StandardShiftHoursPerDay
}

// CalculateExtensionQuote computes the rate differential for extending an active rental.
func CalculateExtensionQuote(currentEndDate, newEndDate time.Time, dailyRateMinor, bundledAttachmentsDailyMinor int64) (extensionDays int, incrementalGross, incrementalPrepayment int64, err error) {
	currentNorm := time.Date(currentEndDate.Year(), currentEndDate.Month(), currentEndDate.Day(), 0, 0, 0, 0, time.UTC)
	newNorm := time.Date(newEndDate.Year(), newEndDate.Month(), newEndDate.Day(), 0, 0, 0, 0, time.UTC)

	if !newNorm.After(currentNorm) {
		return 0, 0, 0, fmt.Errorf("%w: new end date %s must be strictly after current end date %s",
			ErrInvalidExtensionDate, newEndDate.Format("2006-01-02"), currentEndDate.Format("2006-01-02"))
	}

	days := int(newNorm.Sub(currentNorm).Hours() / 24)
	if days <= 0 {
		return 0, 0, 0, fmt.Errorf("%w: extension days must be at least 1", ErrInvalidExtensionDate)
	}

	totalDailyMinor := dailyRateMinor + bundledAttachmentsDailyMinor
	incrementalGross = int64(days) * totalDailyMinor
	incrementalPrepayment = incrementalGross / 2 // Strict 50% integer minor units

	return days, incrementalGross, incrementalPrepayment, nil
}

// ValidateFaultCategory asserts that a reported mechanical fault category is standard.
func ValidateFaultCategory(category string) error {
	catUpper := strings.ToUpper(strings.TrimSpace(category))
	if !ValidFaultCategories[catUpper] {
		return fmt.Errorf("%w: %s", ErrInvalidFaultCategory, category)
	}
	return nil
}
