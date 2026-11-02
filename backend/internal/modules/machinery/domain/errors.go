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

	// Discovery & Comparison Errors (Partition 05-04, SRS §05.04, SR-03, SR-08)
	ErrNonComparableCategories    = errors.New("fleets belong to different subcategories and cannot be compared")
	ErrInvalidComparisonCount    = errors.New("comparison requires between 2 and 3 fleets")
	ErrSpatialRadiusOutOfBounds   = errors.New("search radius must be between 1 and 500 km")
	ErrInvalidDateRange           = errors.New("start date must be before or equal to end date")
	ErrInvalidRateRange           = errors.New("minimum daily rate cannot exceed maximum daily rate")
)
