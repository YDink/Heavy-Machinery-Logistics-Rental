package domain

import (
	"fmt"
	"strings"
	"time"
)

const BreakdownOwnerSLADuration = 4 * time.Hour

var ValidFaultCategories = map[string]bool{
	"ENGINE":            true,
	"HYDRAULIC_CIRCUIT": true,
	"ELECTRICAL":        true,
	"UNDERCARRIAGE":     true,
	"STRUCTURAL":        true,
	"BOOM_BUCKET":       true,
	"OTHER":             true,
}

// ValidateFaultCategory validates mechanical breakdown categories.
func ValidateFaultCategory(category string) error {
	if !ValidFaultCategories[strings.ToUpper(category)] {
		return fmt.Errorf("%w: unrecognized fault category %s", ErrInvalidFaultCategory, category)
	}
	return nil
}
