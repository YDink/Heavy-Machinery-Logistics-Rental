package domain

// ValidateBaselineMeter asserts that the cabin mechanical engine hour meter reading is non-negative.
func ValidateBaselineMeter(hours float64) error {
	if hours < 0.0 {
		return ErrInvalidBaselineMeter
	}
	return nil
}
