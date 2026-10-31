package domain

import "errors"

var (
	ErrInvalidBaselineMeter       = errors.New("baseline hour meter reading must be greater than or equal to 0.00")
	ErrIncompletePhotoWalkaround  = errors.New("fleet registration mandates all 6 required angles: FRONT, REAR, LEFT, RIGHT, CABIN, METER")
	ErrInvalidSpecValue           = errors.New("dynamic specification value is invalid or out of allowed bounds")
	ErrCategoryNotFound           = errors.New("machinery category was not found")
	ErrFleetNotFound              = errors.New("fleet asset was not found")
	ErrUnauthorizedFleetOwner     = errors.New("caller is not the authorized owner of this fleet asset")
	ErrMalwareDetected            = errors.New("binary payload failed ClamAV antivirus stream inspection")
	ErrMissingRequiredSpec        = errors.New("mandatory specification field is missing")
)
