package domain

import "testing"

func TestCalculateShiftPacing(t *testing.T) {
	t.Run("Under shift quota", func(t *testing.T) {
		pacing := CalculateShiftPacing(4.0)
		if pacing != 0.5 {
			t.Fatalf("expected 0.5, got %.2f", pacing)
		}
	})

	t.Run("Exact shift quota", func(t *testing.T) {
		pacing := CalculateShiftPacing(8.0)
		if pacing != 1.0 {
			t.Fatalf("expected 1.0, got %.2f", pacing)
		}
	})

	t.Run("Excess shift pacing", func(t *testing.T) {
		pacing := CalculateShiftPacing(10.0)
		if pacing != 1.25 {
			t.Fatalf("expected 1.25, got %.2f", pacing)
		}
	})
}
