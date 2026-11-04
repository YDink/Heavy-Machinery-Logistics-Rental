package domain

// ValidateWalkaroundPhotos asserts that all 6 required physical angles are present (FL-05).
func ValidateWalkaroundPhotos(photos []FleetPhoto) error {
	anglePresent := make(map[PhotoAngle]bool)
	for _, p := range photos {
		anglePresent[p.Angle] = true
	}

	for _, required := range MandatoryPhotoAngles {
		if !anglePresent[required] {
			return ErrIncompletePhotoWalkaround
		}
	}
	return nil
}
