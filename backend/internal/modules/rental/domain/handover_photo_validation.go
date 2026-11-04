package domain

import (
	"fmt"
	"strings"
)

var ValidInspectionAngles = map[string]bool{
	"METER_CLOSEUP":   true,
	"FRONT_ELEVATION": true,
	"REAR_ELEVATION":  true,
	"LEFT_TRACK":      true,
	"RIGHT_TRACK":     true,
	"CABIN_INTERIOR":  true,
}

// ValidateHandoverPhotos asserts that inspection photos are provided and use valid angles.
func ValidateHandoverPhotos(photos []InspectionPhoto) error {
	if len(photos) == 0 {
		return fmt.Errorf("at least 1 reticle inspection photo is required")
	}
	for _, p := range photos {
		if !ValidInspectionAngles[strings.ToUpper(p.Angle)] {
			return fmt.Errorf("invalid inspection photo angle: %s", p.Angle)
		}
		if strings.TrimSpace(p.StorageURI) == "" {
			return fmt.Errorf("inspection photo storage URI cannot be empty")
		}
	}
	return nil
}
