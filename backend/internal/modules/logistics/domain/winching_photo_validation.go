package domain

import (
	"fmt"
	"strings"
)

var requiredWinchingAngles = []WinchingPhotoAngle{
	WinchingAngleRigProfile,
	WinchingAngleChainBinders,
	WinchingAngleBoomPinned,
	WinchingAngleRampsLocked,
}

// ValidateWinchingPhotos enforces that all 4 mandatory winching and lashing reticle angles are present (SRS §5.1).
func ValidateWinchingPhotos(photos []WinchingPhoto) error {
	seen := make(map[WinchingPhotoAngle]bool, len(requiredWinchingAngles))
	for _, p := range photos {
		if strings.TrimSpace(p.StorageURI) == "" {
			return fmt.Errorf("%w: empty storageUri for angle %s", ErrInvalidWinchingPhotos, p.Angle)
		}
		seen[p.Angle] = true
	}

	for _, angle := range requiredWinchingAngles {
		if !seen[angle] {
			return fmt.Errorf("%w: missing mandatory angle %s", ErrInvalidWinchingPhotos, angle)
		}
	}

	return nil
}
