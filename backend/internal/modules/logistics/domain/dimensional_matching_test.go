package domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestEvaluateDimensionalCompatibility(t *testing.T) {
	limits := CountryLogisticsLimits{
		CountryCode:         "ET",
		DefaultCurrency:     "ETB",
		MaxGvmKg:            56000.0,
		MaxCorridorHeightMm: 4500,
		WideLoadWidthMm:     3000,
	}

	cat320D := MachineLogisticsSpecs{
		OperatingWeightKg: 21800.0,
		TransportLengthMm: 9530,
		TransportWidthMm:  2980,
		TransportHeightMm: 3150,
	}

	validTrailer := MatchedCarrier{
		FleetID:               uuid.New(),
		PayloadCapacityKg:     34000.0,
		DeckWellLengthMm:      11000,
		DeckWidthMm:           3100,
		OutriggersExtensionMm: 250,
		DeckHeightGroundMm:    900,
		TareWeightKg:          12500.0,
		AxleCount:             3,
	}

	t.Run("Valid Tri-Axle Lowboy for CAT 320D", func(t *testing.T) {
		abnormal, err := EvaluateDimensionalCompatibility(cat320D, validTrailer, limits)
		if err != nil {
			t.Fatalf("expected valid compatibility, got error: %v", err)
		}
		if abnormal {
			t.Fatalf("expected abnormalLoad=false for width 2980mm <= 3000mm")
		}
	})

	t.Run("Rejects Insufficient Payload Margin (1.10x Safety Margin)", func(t *testing.T) {
		undersized := validTrailer
		undersized.PayloadCapacityKg = 22000.0 // > 21800 but < 23980 (1.10x)
		_, err := EvaluateDimensionalCompatibility(cat320D, undersized, limits)
		if !errors.Is(err, ErrTrailerPayloadInsufficient) {
			t.Fatalf("expected ErrTrailerPayloadInsufficient, got %v", err)
		}
	})

	t.Run("Rejects Short Deck Length", func(t *testing.T) {
		shortDeck := validTrailer
		shortDeck.DeckWellLengthMm = 8500
		_, err := EvaluateDimensionalCompatibility(cat320D, shortDeck, limits)
		if !errors.Is(err, ErrTrailerDeckTooShort) {
			t.Fatalf("expected ErrTrailerDeckTooShort, got %v", err)
		}
	})

	t.Run("Rejects Narrow Deck Width", func(t *testing.T) {
		narrowDeck := validTrailer
		narrowDeck.DeckWidthMm = 2500
		narrowDeck.OutriggersExtensionMm = 0
		_, err := EvaluateDimensionalCompatibility(cat320D, narrowDeck, limits)
		if !errors.Is(err, ErrTrailerDeckTooNarrow) {
			t.Fatalf("expected ErrTrailerDeckTooNarrow, got %v", err)
		}
	})

	t.Run("Rejects Corridor Height Exceeded (> 4500mm)", func(t *testing.T) {
		tallDeck := validTrailer
		tallDeck.DeckHeightGroundMm = 1500 // 1500 + 3150 = 4650 > 4500
		_, err := EvaluateDimensionalCompatibility(cat320D, tallDeck, limits)
		if !errors.Is(err, ErrCorridorHeightExceeded) {
			t.Fatalf("expected ErrCorridorHeightExceeded, got %v", err)
		}
	})

	t.Run("Rejects EAC 56-Tonne GVM Limit Exceeded (D-05-08-002)", func(t *testing.T) {
		heavyMachine := cat320D
		heavyMachine.OperatingWeightKg = 45000.0
		heavyTrailer := validTrailer
		heavyTrailer.PayloadCapacityKg = 55000.0
		heavyTrailer.TareWeightKg = 14000.0 // 45000 + 14000 = 59000 > 56000
		_, err := EvaluateDimensionalCompatibility(heavyMachine, heavyTrailer, limits)
		if !errors.Is(err, ErrGVMExceeded) {
			t.Fatalf("expected ErrGVMExceeded, got %v", err)
		}
	})

	t.Run("Flags Abnormal Wide Load Escort Required (> 3000mm)", func(t *testing.T) {
		wideMachine := cat320D
		wideMachine.TransportWidthMm = 3250
		abnormal, err := EvaluateDimensionalCompatibility(wideMachine, validTrailer, limits)
		if err != nil {
			t.Fatalf("expected valid compatibility with outriggers, got %v", err)
		}
		if !abnormal {
			t.Fatalf("expected abnormalLoad=true for width 3250mm > 3000mm")
		}
	})
}
