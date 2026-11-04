package domain

import "errors"

var (
	ErrInvalidDateRange               = errors.New("start date must be in the future and before end date")
	ErrInvalidDuration                = errors.New("rental duration must meet fleet minimum days")
	ErrConcurrentReservationCollision = errors.New("fleet or attachment is currently reserved or on hold for overlapping dates")
	ErrInvalidStateTransition         = errors.New("invalid booking state transition")
	ErrBookingNotFound                = errors.New("booking not found")
	ErrUnauthorizedAction             = errors.New("user is not authorized to perform action on this booking")
	ErrConcurrencyConflict            = errors.New("booking was updated by another transaction (version conflict)")
	ErrInvalidMeterReading            = errors.New("engine meter hours cannot regress (must be >= baseline check-in)")
	ErrInvalidChallengePIN            = errors.New("invalid 6-digit handover challenge PIN")
	ErrSlipNotFound                   = errors.New("bank transfer payment slip not found")
	ErrInvalidSlipStatus              = errors.New("bank slip status is not eligible for audit")
	ErrEmptyReason                    = errors.New("reason code and explanation cannot be empty")
	ErrSelfBookingProhibited          = errors.New("fleet owner cannot book their own machinery")
)
