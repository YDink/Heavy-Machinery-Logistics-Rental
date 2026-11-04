package domain

import "github.com/google/uuid"

type ActiveRentalHUD struct {
	BookingID             uuid.UUID
	Status                string
	FleetID               uuid.UUID
	FleetTitle            string
	PlateNumber           string
	OperatorID            *uuid.UUID
	OperatorFullName      string
	OperatorPhoneNumber   string
	OperatorLicenseClass  string
	CurrentRentalDay      int
	TotalRentalDays       int
	HandoverBaselineHours float64
	AccumulatedHoursToday float64
	AllowedHoursPerDay    float64
	ExcessHoursTotal      float64
	ShiftTimerPaused      bool
	CanExtend             bool
	CanReportBreakdown    bool
	CanInitiateReturn     bool
}
