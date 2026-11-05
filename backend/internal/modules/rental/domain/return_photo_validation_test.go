package domain_test

import (
	"testing"

	"github.com/dinkmachinery/dink-core/internal/modules/rental/domain"
)

func TestValidateReturnPhotos_RequiresAllSixAngles(t *testing.T) {
	validPhotos := []domain.InspectionPhoto{
		{Angle: "FRONT_IMPLEMENT", StorageURI: "https://cdn.example.com/front.jpg"},
		{Angle: "REAR_COUNTERWEIGHT_ENGINE", StorageURI: "https://cdn.example.com/rear.jpg"},
		{Angle: "LEFT_SIDE_TRACKS", StorageURI: "https://cdn.example.com/left.jpg"},
		{Angle: "RIGHT_SIDE_TRACKS", StorageURI: "https://cdn.example.com/right.jpg"},
		{Angle: "HOUR_METER", StorageURI: "https://cdn.example.com/meter.jpg"},
		{Angle: "CABIN_INTERIOR", StorageURI: "https://cdn.example.com/cabin.jpg"},
	}

	err := domain.ValidateReturnPhotos(validPhotos)
	if err != nil {
		t.Fatalf("expected valid 6 angles to pass, got error: %v", err)
	}

	// Missing one angle (CABIN_INTERIOR omitted)
	incompletePhotos := validPhotos[:5]
	err = domain.ValidateReturnPhotos(incompletePhotos)
	if err != domain.ErrMissingReturnPhotos {
		t.Fatalf("expected ErrMissingReturnPhotos when only 5 angles provided, got: %v", err)
	}
}
