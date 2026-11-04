package domain

import (
	"errors"
	"testing"
)

func TestValidateMonotonicMeterFloor(t *testing.T) {
	t.Run("Valid meter reading equal to baseline", func(t *testing.T) {
		err := ValidateMonotonicMeterFloor(1250.0, 1250.0)
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
	})

	t.Run("Valid meter reading higher than baseline", func(t *testing.T) {
		err := ValidateMonotonicMeterFloor(1255.5, 1250.0)
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
	})

	t.Run("Rejects regressive meter reading", func(t *testing.T) {
		err := ValidateMonotonicMeterFloor(1249.9, 1250.0)
		if err == nil {
			t.Fatal("expected error for meter regression, got nil")
		}
		if !errors.Is(err, ErrMeterReadingRegression) {
			t.Fatalf("expected ErrMeterReadingRegression, got: %v", err)
		}
	})
}
