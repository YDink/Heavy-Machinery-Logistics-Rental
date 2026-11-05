package domain

var MandatoryReturnAngles = []string{
	"FRONT_IMPLEMENT",
	"REAR_COUNTERWEIGHT_ENGINE",
	"LEFT_SIDE_TRACKS",
	"RIGHT_SIDE_TRACKS",
	"HOUR_METER",
	"CABIN_INTERIOR",
}

// ValidateReturnPhotos asserts that all 6 mandatory return angles are present and attested.
func ValidateReturnPhotos(photos []InspectionPhoto) error {
	if len(photos) < len(MandatoryReturnAngles) {
		return ErrMissingReturnPhotos
	}

	captured := make(map[string]bool)
	for _, p := range photos {
		if p.StorageURI != "" {
			captured[p.Angle] = true
		}
	}

	for _, required := range MandatoryReturnAngles {
		if !captured[required] {
			return ErrMissingReturnPhotos
		}
	}

	return nil
}
