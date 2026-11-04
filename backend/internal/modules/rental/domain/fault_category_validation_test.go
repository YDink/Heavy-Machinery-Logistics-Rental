package domain

import (
	"testing"
)

func TestValidateFaultCategory(t *testing.T) {
	t.Run("Valid categories", func(t *testing.T) {
		for _, cat := range []string{"ENGINE", "HYDRAULIC_CIRCUIT", "ELECTRICAL", "UNDERCARRIAGE", "STRUCTURAL", "BOOM_BUCKET", "OTHER"} {
			if err := ValidateFaultCategory(cat); err != nil {
				t.Fatalf("expected valid category %s, got: %v", cat, err)
			}
		}
	})

	t.Run("Invalid category", func(t *testing.T) {
		if err := ValidateFaultCategory("FLAT_TIRE_EMERGENCY"); err == nil {
			t.Fatal("expected error for invalid fault category, got nil")
		}
	})
}
