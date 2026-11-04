package domain

// InferLogisticsTransportClass automatically calculates required trailer class per SRS §7.2.
func InferLogisticsTransportClass(operatingWeightKg float64) string {
	if operatingWeightKg < 12000 {
		return "Two-Axle Tilt-Bed / Rigid Flatbed Carrier"
	} else if operatingWeightKg < 24000 {
		return "Two-Axle Lowboy Trailer (25t Payload Capacity)"
	} else if operatingWeightKg < 40000 {
		return "Tri-Axle Heavy Lowboy Trailer (40t Payload Capacity)"
	}
	return "Multi-Axle Modular Hydraulic Platform Trailer (> 60t Capacity)"
}
