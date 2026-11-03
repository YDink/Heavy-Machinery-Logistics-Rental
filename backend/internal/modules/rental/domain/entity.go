package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrMeterReadingRegression     = errors.New("check-in meter reading cannot be less than machine certified baseline (monotonic floor violation)")
	ErrInvalidChallengePIN        = errors.New("handover challenge PIN does not match")
	ErrPINExpired                 = errors.New("handover challenge PIN has expired (15-minute TTL exceeded)")
	ErrPINLockout                 = errors.New("maximum verification attempts exceeded; challenge locked out")
	ErrInvalidExtensionDate       = errors.New("requested extension date must be strictly after the current end date")
	ErrExtensionDateUnavailable   = errors.New("equipment is already reserved during the requested forward extension period")
	ErrInvalidFaultCategory       = errors.New("unrecognized mechanical fault category")
	ErrBookingNotActive           = errors.New("booking is not currently active")
	ErrBookingAlreadyPaused       = errors.New("booking is already paused for an ongoing breakdown")
	ErrBreakdownNotFound          = errors.New("mechanical breakdown record not found")
	ErrActiveBreakdownStillOpen   = errors.New("cannot resolve breakdown: active incident not in repairable state")
)

type InspectionPhoto struct {
	Angle      string `json:"angle"`
	StorageURI string `json:"storageUri"`
}

type HandoverInspection struct {
	ID                uuid.UUID         `json:"id"`
	BookingID         uuid.UUID         `json:"bookingId"`
	InspectorID       uuid.UUID         `json:"inspectorId"`
	CheckInMeterHours float64           `json:"checkInMeterHours"`
	FluidLevelsOk     bool              `json:"fluidLevelsOk"`
	SafetyGlassIntact bool              `json:"safetyGlassIntact"`
	Notes             string            `json:"notes,omitempty"`
	Photos            []InspectionPhoto `json:"photos"`
	CreatedAt         time.Time         `json:"createdAt"`
}

type HandoverChallenge struct {
	ID                   uuid.UUID  `json:"id"`
	BookingID            uuid.UUID  `json:"bookingId"`
	ChallengePIN         string     `json:"challengePin"`
	ExpiresAt            time.Time  `json:"expiresAt"`
	AttemptCount         int        `json:"attemptCount"`
	MaxAttempts          int        `json:"maxAttempts"`
	IsVerified           bool       `json:"isVerified"`
	VerifiedAt           *time.Time `json:"verifiedAt,omitempty"`
	TouchSignatureBase64 string     `json:"touchSignatureBase64,omitempty"`
	ContractorLat        *float64   `json:"contractorLat,omitempty"`
	ContractorLng        *float64   `json:"contractorLng,omitempty"`
	CreatedAt            time.Time  `json:"createdAt"`
}

type RentalShift struct {
	ID               uuid.UUID  `json:"id"`
	BookingID        uuid.UUID  `json:"bookingId"`
	OperatorID       *uuid.UUID `json:"operatorId,omitempty"`
	ShiftDate        time.Time  `json:"shiftDate"`
	StartTime        time.Time  `json:"startTime"`
	EndTime          *time.Time `json:"endTime,omitempty"`
	StartMeterHours  float64    `json:"startMeterHours"`
	EndMeterHours    *float64   `json:"endMeterHours,omitempty"`
	AllocatedHours   float64    `json:"allocatedHours"`
	AccumulatedHours float64    `json:"accumulatedHours"`
	ExcessHours      float64    `json:"excessHours"`
	IsPaused         bool       `json:"isPaused"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

type RentalExtension struct {
	ID                         uuid.UUID  `json:"id"`
	BookingID                  uuid.UUID  `json:"bookingId"`
	RequestedByID              uuid.UUID  `json:"requestedById"`
	OriginalEndDate            time.Time  `json:"originalEndDate"`
	NewEndDate                 time.Time  `json:"newEndDate"`
	ExtensionDays              int        `json:"extensionDays"`
	IncrementalGrossMinor      int64      `json:"incrementalGrossMinor"`
	IncrementalPrepaymentMinor int64      `json:"incrementalPrepaymentMinor"`
	Currency                   string     `json:"currency"`
	Status                     string     `json:"status"`
	OwnerDecidedAt             *time.Time `json:"ownerDecidedAt,omitempty"`
	OwnerDecisionNotes         string     `json:"ownerDecisionNotes,omitempty"`
	SlipStorageURI             string     `json:"slipStorageUri,omitempty"`
	SlipFileHash               string     `json:"slipFileHash,omitempty"`
	ConfirmedAt                *time.Time `json:"confirmedAt,omitempty"`
	CreatedAt                  time.Time  `json:"createdAt"`
	UpdatedAt                  time.Time  `json:"updatedAt"`
}

type MechanicalBreakdown struct {
	ID                     uuid.UUID  `json:"id"`
	BookingID              uuid.UUID  `json:"bookingId"`
	ReportedByID           uuid.UUID  `json:"reportedById"`
	FaultCategory          string     `json:"faultCategory"`
	FaultDescription       string     `json:"faultDescription"`
	DiagnosticTroubleCodes []string   `json:"diagnosticTroubleCodes"`
	CurrentHourMeter       float64    `json:"currentHourMeter"`
	IncidentPhotos         []string   `json:"incidentPhotos"`
	Status                 string     `json:"status"`
	SLADeadline            time.Time  `json:"slaDeadline"`
	MechanicName           string     `json:"mechanicName,omitempty"`
	MechanicPhone          string     `json:"mechanicPhone,omitempty"`
	ResolutionNotes        string     `json:"resolutionNotes,omitempty"`
	ResolvedAt             *time.Time `json:"resolvedAt,omitempty"`
	CreatedAt              time.Time  `json:"createdAt"`
	UpdatedAt              time.Time  `json:"updatedAt"`
}

type ActiveRentalHUD struct {
	BookingID             uuid.UUID `json:"bookingId"`
	Status                string    `json:"status"`
	FleetID               uuid.UUID `json:"fleetId"`
	FleetTitle            string    `json:"fleetTitle"`
	PlateNumber           string    `json:"plateNumber"`
	OperatorID            *uuid.UUID `json:"operatorId,omitempty"`
	OperatorFullName      string    `json:"operatorFullName,omitempty"`
	OperatorPhoneNumber   string    `json:"operatorPhoneNumber,omitempty"`
	OperatorLicenseClass  string    `json:"operatorLicenseClass,omitempty"`
	CurrentRentalDay      int       `json:"currentRentalDay"`
	TotalRentalDays       int       `json:"totalRentalDays"`
	HandoverBaselineHours float64   `json:"handoverBaselineHours"`
	AccumulatedHoursToday float64   `json:"accumulatedHoursToday"`
	AllowedHoursPerDay    float64   `json:"allowedHoursPerDay"`
	ExcessHoursTotal      float64   `json:"excessHoursTotal"`
	ShiftTimerPaused      bool      `json:"shiftTimerPaused"`
	CanExtend             bool      `json:"canExtend"`
	CanReportBreakdown    bool      `json:"canReportBreakdown"`
	CanInitiateReturn     bool      `json:"canInitiateReturn"`
}
