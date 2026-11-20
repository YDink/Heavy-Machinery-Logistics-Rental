package domain

// WinchingPhotoAngle defines the 4 mandatory pre-transit load winching and lashing reticle angles.
type WinchingPhotoAngle string

const (
	WinchingAngleRigProfile   WinchingPhotoAngle = "RIG_PROFILE"
	WinchingAngleChainBinders WinchingPhotoAngle = "CHAIN_BINDERS"
	WinchingAngleBoomPinned   WinchingPhotoAngle = "BOOM_PINNED"
	WinchingAngleRampsLocked  WinchingPhotoAngle = "RAMPS_LOCKED"
)

// WinchingPhoto represents a single verified load-securement photograph.
type WinchingPhoto struct {
	Angle      WinchingPhotoAngle
	StorageURI string
}
