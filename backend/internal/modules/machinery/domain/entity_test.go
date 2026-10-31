package domain_test

import (
	"testing"

	"github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
)

func TestValidateBaselineMeter(t *testing.T) {
	if err := domain.ValidateBaselineMeter(1450.50); err != nil {
		t.Errorf("expected valid for 1450.50, got %v", err)
	}

	if err := domain.ValidateBaselineMeter(0.0); err != nil {
		t.Errorf("expected valid for 0.0, got %v", err)
	}

	if err := domain.ValidateBaselineMeter(-10.5); err == nil {
		t.Error("expected error for negative meter reading, got nil")
	}
}

func TestValidateWalkaroundPhotos(t *testing.T) {
	// Missing CABIN and METER angles
	incomplete := []domain.FleetPhoto{
		{Angle: domain.PhotoAngleFront},
		{Angle: domain.PhotoAngleRear},
		{Angle: domain.PhotoAngleLeft},
		{Angle: domain.PhotoAngleRight},
	}

	if err := domain.ValidateWalkaroundPhotos(incomplete); err != domain.ErrIncompletePhotoWalkaround {
		t.Errorf("expected ErrIncompletePhotoWalkaround, got %v", err)
	}

	// Complete 6 angles
	complete := []domain.FleetPhoto{
		{Angle: domain.PhotoAngleFront},
		{Angle: domain.PhotoAngleRear},
		{Angle: domain.PhotoAngleLeft},
		{Angle: domain.PhotoAngleRight},
		{Angle: domain.PhotoAngleCabin},
		{Angle: domain.PhotoAngleMeter},
	}

	if err := domain.ValidateWalkaroundPhotos(complete); err != nil {
		t.Errorf("expected complete walkaround to be valid, got %v", err)
	}
}

func TestInferLogisticsTransportClass(t *testing.T) {
	tests := []struct {
		weightKg float64
		expected string
	}{
		{5000, "Two-Axle Tilt-Bed / Rigid Flatbed Carrier"},
		{11999, "Two-Axle Tilt-Bed / Rigid Flatbed Carrier"},
		{12000, "Two-Axle Lowboy Trailer (25t Payload Capacity)"},
		{22000, "Two-Axle Lowboy Trailer (25t Payload Capacity)"},
		{24000, "Tri-Axle Heavy Lowboy Trailer (40t Payload Capacity)"},
		{38000, "Tri-Axle Heavy Lowboy Trailer (40t Payload Capacity)"},
		{40000, "Multi-Axle Modular Hydraulic Platform Trailer (> 60t Capacity)"},
		{65000, "Multi-Axle Modular Hydraulic Platform Trailer (> 60t Capacity)"},
	}

	for _, tt := range tests {
		result := domain.InferLogisticsTransportClass(tt.weightKg)
		if result != tt.expected {
			t.Errorf("InferLogisticsTransportClass(%f) = %q; want %q", tt.weightKg, result, tt.expected)
		}
	}
}
