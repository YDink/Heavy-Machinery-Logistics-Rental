package domain

import (
	"time"

	"github.com/google/uuid"
)

type MechanicalBreakdown struct {
	ID                     uuid.UUID
	BookingID              uuid.UUID
	ReportedByID           uuid.UUID
	FaultCategory          string
	FaultDescription       string
	DiagnosticTroubleCodes []string
	CurrentHourMeter       float64
	IncidentPhotos         []string
	Status                 string
	SLADeadline            time.Time
	MechanicName           string
	MechanicPhone          string
	ResolutionNotes        string
	ResolvedAt             *time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
}
