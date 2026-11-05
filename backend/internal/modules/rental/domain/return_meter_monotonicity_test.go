package domain_test

import (
	"testing"

	"github.com/dinkmachinery/dink-core/internal/modules/rental/domain"
)

func TestValidateReturnMeterMonotonicity(t *testing.T) {
	t.Run("valid equal or higher meter readings pass", func(t *testing.T) {
		if err := domain.ValidateReturnMeterMonotonicity(1250.50, 1250.50); err != nil {
			t.Errorf("expected equal readings to pass, got error: %v", err)
		}
		if err := domain.ValidateReturnMeterMonotonicity(1250.50, 1312.20); err != nil {
			t.Errorf("expected higher reading to pass, got error: %v", err)
		}
	})

	t.Run("regressive lower meter reading fails with ErrReturnMeterRegression", func(t *testing.T) {
		err := domain.ValidateReturnMeterMonotonicity(1250.50, 1240.00)
		if err != domain.ErrReturnMeterRegression {
			t.Fatalf("expected ErrReturnMeterRegression, got: %v", err)
		}
	})
}
