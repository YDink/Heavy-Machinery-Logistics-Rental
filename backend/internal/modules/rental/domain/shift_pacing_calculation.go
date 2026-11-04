package domain

const StandardShiftHoursPerDay = 8.00

// CalculateShiftPacing computes cumulative shift pacing ratio against the 8.0-hour daily standard allocation.
func CalculateShiftPacing(accumulatedHours float64) float64 {
	if StandardShiftHoursPerDay <= 0 {
		return 0.0
	}
	return accumulatedHours / StandardShiftHoursPerDay
}
