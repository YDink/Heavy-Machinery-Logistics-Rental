package domain

import "errors"

var (
	ErrMeterReadingRegression   = errors.New("check-in meter reading cannot be less than machine certified baseline (monotonic floor violation)")
	ErrInvalidChallengePIN      = errors.New("handover challenge PIN does not match")
	ErrPINExpired               = errors.New("handover challenge PIN has expired (15-minute TTL exceeded)")
	ErrPINLockout               = errors.New("maximum verification attempts exceeded; challenge locked out")
	ErrInvalidExtensionDate     = errors.New("requested extension date must be strictly after the current end date")
	ErrExtensionDateUnavailable = errors.New("equipment is already reserved during the requested forward extension period")
	ErrInvalidFaultCategory     = errors.New("unrecognized mechanical fault category")
	ErrBookingNotActive         = errors.New("booking is not currently active")
	ErrBookingAlreadyPaused     = errors.New("booking is already paused for an ongoing breakdown")
	ErrBreakdownNotFound        = errors.New("mechanical breakdown record not found")
	ErrActiveBreakdownStillOpen = errors.New("cannot resolve breakdown: active incident not in repairable state")
)
