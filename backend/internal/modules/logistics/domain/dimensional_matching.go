package domain

import "fmt"

const (
	DynamicLashingSafetyFactor = 1.10
	DefaultMaxCorridorHeightMm = 4500
	DefaultMaxGvmKg            = 56000.0
	DefaultWideLoadWidthMm     = 3000
)

// CountryLogisticsLimits holds jurisdictional highway clearance and mass caps.
type CountryLogisticsLimits struct {
	CountryCode         string
	DefaultCurrency     string
	MaxGvmKg            float64
	MaxCorridorHeightMm int
	WideLoadWidthMm     int
}

// ComputeMinimumRequiredPayloadKg calculates required trailer payload including the 10% dynamic lashing safety margin (FL-29).
func ComputeMinimumRequiredPayloadKg(operatingWeightKg float64) float64 {
	return operatingWeightKg * DynamicLashingSafetyFactor
}

// EvaluateDimensionalCompatibility checks payload capacity, deck dimensions, corridor bridge height, and EAC GVM safety caps (D-05-08-002, SRS §3.3).
func EvaluateDimensionalCompatibility(
	machine MachineLogisticsSpecs,
	trailer MatchedCarrier,
	limits CountryLogisticsLimits,
) (abnormalLoadEscortRequired bool, err error) {
	if limits.MaxCorridorHeightMm <= 0 {
		limits.MaxCorridorHeightMm = DefaultMaxCorridorHeightMm
	}
	if limits.MaxGvmKg <= 0 {
		limits.MaxGvmKg = DefaultMaxGvmKg
	}
	if limits.WideLoadWidthMm <= 0 {
		limits.WideLoadWidthMm = DefaultWideLoadWidthMm
	}

	minRequiredPayload := ComputeMinimumRequiredPayloadKg(machine.OperatingWeightKg)
	if trailer.PayloadCapacityKg < minRequiredPayload {
		return false, fmt.Errorf(
			"%w: trailer capacity %.0f kg < required %.0f kg (machine %.0f kg * 1.10)",
			ErrTrailerPayloadInsufficient, trailer.PayloadCapacityKg, minRequiredPayload, machine.OperatingWeightKg,
		)
	}

	if trailer.DeckWellLengthMm < machine.TransportLengthMm {
		return false, fmt.Errorf(
			"%w: deck length %d mm < machine transport length %d mm",
			ErrTrailerDeckTooShort, trailer.DeckWellLengthMm, machine.TransportLengthMm,
		)
	}

	effectiveDeckWidth := trailer.DeckWidthMm + trailer.OutriggersExtensionMm
	if effectiveDeckWidth < machine.TransportWidthMm {
		return false, fmt.Errorf(
			"%w: effective deck width %d mm < machine transport width %d mm",
			ErrTrailerDeckTooNarrow, effectiveDeckWidth, machine.TransportWidthMm,
		)
	}

	combinedHeight := trailer.DeckHeightGroundMm + machine.TransportHeightMm
	if combinedHeight > limits.MaxCorridorHeightMm {
		return false, fmt.Errorf(
			"%w: combined height %d mm > corridor limit %d mm",
			ErrCorridorHeightExceeded, combinedHeight, limits.MaxCorridorHeightMm,
		)
	}

	totalGvm := machine.OperatingWeightKg + trailer.TareWeightKg
	if totalGvm > limits.MaxGvmKg {
		return false, fmt.Errorf(
			"%w: combined GVM %.0f kg > regional limit %.0f kg",
			ErrGVMExceeded, totalGvm, limits.MaxGvmKg,
		)
	}

	abnormalLoadEscortRequired = machine.TransportWidthMm > limits.WideLoadWidthMm
	return abnormalLoadEscortRequired, nil
}
