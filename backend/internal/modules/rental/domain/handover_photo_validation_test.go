package domain

import "testing"

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
