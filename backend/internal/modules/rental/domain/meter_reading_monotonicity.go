package domain

import "fmt"

// ValidateMonotonicMeterFloor asserts that the check-in meter reading cannot regress below the machine certified baseline.
func ValidateMonotonicMeterFloor(checkInMeter, baselineMeter float64) error {
	if checkInMeter < baselineMeter {
		return fmt.Errorf("%w: check-in reading (%.2f hrs) is lower than certified baseline (%.2f hrs)",
			ErrMeterReadingRegression, checkInMeter, baselineMeter)
	}
	return nil
}
