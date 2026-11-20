package domain

import (
	"errors"
	"testing"
)

func TestValidateWinchingPhotos(t *testing.T) {
	validPhotos := []WinchingPhoto{
		{Angle: WinchingAngleRigProfile, StorageURI: "https://cdn.dinkmachinery.com/winch/1.webp"},
		{Angle: WinchingAngleChainBinders, StorageURI: "https://cdn.dinkmachinery.com/winch/2.webp"},
		{Angle: WinchingAngleBoomPinned, StorageURI: "https://cdn.dinkmachinery.com/winch/3.webp"},
		{Angle: WinchingAngleRampsLocked, StorageURI: "https://cdn.dinkmachinery.com/winch/4.webp"},
	}

	if err := ValidateWinchingPhotos(validPhotos); err != nil {
		t.Fatalf("expected valid 4-angle winching photos, got %v", err)
	}

	t.Run("Rejects Missing Angle", func(t *testing.T) {
		err := ValidateWinchingPhotos(validPhotos[:3])
		if !errors.Is(err, ErrInvalidWinchingPhotos) {
			t.Fatalf("expected ErrInvalidWinchingPhotos, got %v", err)
		}
	})

	t.Run("Rejects Blank URI", func(t *testing.T) {
		bad := append([]WinchingPhoto{}, validPhotos...)
		bad[2].StorageURI = "   "
		err := ValidateWinchingPhotos(bad)
		if !errors.Is(err, ErrInvalidWinchingPhotos) {
			t.Fatalf("expected ErrInvalidWinchingPhotos, got %v", err)
		}
	})
}
