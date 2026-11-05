package domain

// ValidateReturnMeterMonotonicity asserts that the check-out hour meter reading is greater than or equal to check-in.
func ValidateReturnMeterMonotonicity(checkInMeter, checkOutMeter float64) error {
	if checkOutMeter < checkInMeter {
		return ErrReturnMeterRegression
	}
	return nil
}
