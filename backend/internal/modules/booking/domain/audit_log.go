package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID         uuid.UUID
	BookingID  uuid.UUID
	ActorID    *uuid.UUID
	Action     string
	FromStatus BookingStatus
	ToStatus   BookingStatus
	Reason     string
	Payload    map[string]any
	CreatedAt  time.Time
}
