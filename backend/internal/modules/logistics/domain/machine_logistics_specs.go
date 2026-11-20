package domain

// MachineLogisticsSpecs captures the indexed physical weight and transport dimensions of a machine.
type MachineLogisticsSpecs struct {
	OperatingWeightKg float64
	TransportLengthMm int
	TransportWidthMm  int
	TransportHeightMm int
}
