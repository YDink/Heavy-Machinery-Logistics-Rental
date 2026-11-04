package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents the core user entity in the identity domain.
type User struct {
	ID              uuid.UUID
	PhoneNumber     PhoneNumber
	CountryCode     string
	FullName        string
	PrimaryRole     UserRole
	AccountStatus   AccountStatus
	KYCStatus       KYCStatus
	PasswordHash    *string
	TokenVersion    int32
	PhoneVerifiedAt *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
