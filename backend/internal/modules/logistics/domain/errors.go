package domain

import "errors"

var (
	ErrLogisticsBookingNotFound         = errors.New("parent booking not found for logistics coordination")
	ErrSelfMoveNotRequired              = errors.New("machine is self-move capable and does not require lowboy logistics")
	ErrTrailerPayloadInsufficient       = errors.New("ERR_TRAILER_PAYLOAD_INSUFFICIENT: lowboy payload capacity is below required 1.10x dynamic lashing margin")
	ErrTrailerDeckTooShort              = errors.New("ERR_TRAILER_DECK_TOO_SHORT: lowboy deck well length is shorter than machine transport length")
	ErrTrailerDeckTooNarrow             = errors.New("ERR_TRAILER_DECK_TOO_NARROW: lowboy deck width plus outriggers is narrower than machine width")
	ErrCorridorHeightExceeded           = errors.New("ERR_CORRIDOR_HEIGHT_EXCEEDED: combined deck height and machine height exceeds regional bridge corridor clearance")
	ErrGVMExceeded                      = errors.New("ERR_GVM_LIMIT_EXCEEDED: combined machine weight and trailer tare weight exceeds national highway GVM ceiling")
	ErrInvalidWinchingPhotos            = errors.New("ERR_INVALID_WINCHING_PHOTOS: all 4 mandatory winching and lashing reticle angles are required")
	ErrTransshipmentWindowExpired       = errors.New("ERR_TRANSSHIPMENT_TIMEOUT: 60-minute emergency roadside transshipment SLA window has expired")
	ErrTransshipmentNotOpen             = errors.New("no open transshipment request exists for this logistics sub-booking")
	ErrReverseCascadeAlreadyExecuted    = errors.New("Force Majeure reverse cascade has already been executed for this booking")
	ErrInvalidDispatchStage             = errors.New("invalid logistics dispatch stage transition")
	ErrParentBookingBlockedUntilDelivery = errors.New("parent machinery booking cannot transition to ACTIVE until linked logistics delivery completes")
	ErrExternalTransportBlocked         = errors.New("cannot register external transport when an active in-platform logistics sub-booking exists")
	ErrSubBookingAlreadyExists          = errors.New("a logistics sub-booking is already active for this parent booking")
	ErrUnauthorizedLogisticsAction      = errors.New("caller is not authorized to perform this logistics action")
	ErrInvalidBookingStateForLogistics  = errors.New("booking state does not permit logistics mobilization")
)
