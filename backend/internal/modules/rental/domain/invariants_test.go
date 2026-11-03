package domain

import (
	"errors"
	"testing"
	"time"
)

func TestValidateMonotonicMeterFloor(t *testing.T) {
	t.Run("Valid meter reading equal to baseline", func(t *testing.T) {
		err := ValidateMonotonicMeterFloor(1250.0, 1250.0)
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
	})

	t.Run("Valid meter reading higher than baseline", func(t *testing.T) {
		err := ValidateMonotonicMeterFloor(1255.5, 1250.0)
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
	})

	t.Run("Rejects regressive meter reading", func(t *testing.T) {
		err := ValidateMonotonicMeterFloor(1249.9, 1250.0)
		if err == nil {
			t.Fatal("expected error for meter regression, got nil")
		}
		if !errors.Is(err, ErrMeterReadingRegression) {
			t.Fatalf("expected ErrMeterReadingRegression, got: %v", err)
		}
	})
}

func TestValidateHandoverPhotos(t *testing.T) {
	t.Run("Valid photos", func(t *testing.T) {
		photos := []InspectionPhoto{
			{Angle: "METER_CLOSEUP", StorageURI: "s3://inspections/meter.webp"},
			{Angle: "FRONT_ELEVATION", StorageURI: "s3://inspections/front.webp"},
			{Angle: "CABIN_INTERIOR", StorageURI: "s3://inspections/cabin.webp"},
		}
		if err := ValidateHandoverPhotos(photos); err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
	})

	t.Run("Rejects empty photos", func(t *testing.T) {
		if err := ValidateHandoverPhotos([]InspectionPhoto{}); err == nil {
			t.Fatal("expected error for empty photos, got nil")
		}
	})

	t.Run("Rejects invalid angle", func(t *testing.T) {
		photos := []InspectionPhoto{
			{Angle: "UNDERWATER_SELFIE", StorageURI: "s3://inspections/test.webp"},
		}
		if err := ValidateHandoverPhotos(photos); err == nil {
			t.Fatal("expected error for invalid angle, got nil")
		}
	})
}

func TestGenerateChallengePIN(t *testing.T) {
	pin, err := GenerateChallengePIN()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pin) != 6 {
		t.Fatalf("expected 6 digits, got %d (%s)", len(pin), pin)
	}
}

func TestCalculateShiftPacing(t *testing.T) {
	t.Run("Under shift quota", func(t *testing.T) {
		pacing := CalculateShiftPacing(4.0)
		if pacing != 0.5 {
			t.Fatalf("expected 0.5, got %.2f", pacing)
		}
	})

	t.Run("Exact shift quota", func(t *testing.T) {
		pacing := CalculateShiftPacing(8.0)
		if pacing != 1.0 {
			t.Fatalf("expected 1.0, got %.2f", pacing)
		}
	})

	t.Run("Excess shift pacing", func(t *testing.T) {
		pacing := CalculateShiftPacing(10.0)
		if pacing != 1.25 {
			t.Fatalf("expected 1.25, got %.2f", pacing)
		}
	})
}

func TestCalculateExtensionQuote(t *testing.T) {
	now := time.Now().UTC()
	currentEnd := now.Add(48 * time.Hour)
	newEnd := currentEnd.Add(72 * time.Hour) // 3 days extension

	t.Run("Valid 3-day extension with 50% prepayment", func(t *testing.T) {
		extDays, grossMinor, prepayMinor, err := CalculateExtensionQuote(currentEnd, newEnd, 2000000, 500000)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if extDays != 3 {
			t.Fatalf("expected 3 days, got %d", extDays)
		}
		// 3 * 2,500,000 = 7,500,000
		if grossMinor != 7500000 {
			t.Fatalf("expected 7500000 gross minor, got %d", grossMinor)
		}
		// 7,500,000 / 2 = 3,750,000
		if prepayMinor != 3750000 {
			t.Fatalf("expected 3750000 prepay minor, got %d", prepayMinor)
		}
	})

	t.Run("Rejects non-forward date", func(t *testing.T) {
		pastDate := currentEnd.Add(-24 * time.Hour)
		_, _, _, err := CalculateExtensionQuote(currentEnd, pastDate, 2000000, 0)
		if err == nil {
			t.Fatal("expected error for past date, got nil")
		}
		if !errors.Is(err, ErrInvalidExtensionDate) {
			t.Fatalf("expected ErrInvalidExtensionDate, got: %v", err)
		}
	})
}

func TestValidateFaultCategory(t *testing.T) {
	t.Run("Valid categories", func(t *testing.T) {
		for _, cat := range []string{"ENGINE", "HYDRAULIC_CIRCUIT", "ELECTRICAL", "UNDERCARRIAGE", "STRUCTURAL", "BOOM_BUCKET", "OTHER"} {
			if err := ValidateFaultCategory(cat); err != nil {
				t.Fatalf("expected valid category %s, got: %v", cat, err)
			}
		}
	})

	t.Run("Invalid category", func(t *testing.T) {
		if err := ValidateFaultCategory("FLAT_TIRE_EMERGENCY"); err == nil {
			t.Fatal("expected error for invalid fault category, got nil")
		}
	})
}
