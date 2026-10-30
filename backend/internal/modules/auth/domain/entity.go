package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleRenter    UserRole = "RENTER"
	RoleOwner     UserRole = "OWNER"
	RoleOperator  UserRole = "OPERATOR"
	RoleLogistics UserRole = "LOGISTICS"
	RoleAdmin     UserRole = "ADMIN"
)

type AccountStatus string

const (
	AccountStatusUnverified              AccountStatus = "UNVERIFIED"
	AccountStatusActive                  AccountStatus = "ACTIVE"
	AccountStatusSuspendedSecurityFreeze AccountStatus = "SUSPENDED_SECURITY_FREEZE"
	AccountStatusDelinquentLocked        AccountStatus = "DELINQUENT_LOCKED"
	AccountStatusTerminated              AccountStatus = "TERMINATED"
)

type KYCStatus string

const (
	KYCStatusUnverified KYCStatus = "UNVERIFIED"
	KYCStatusPending    KYCStatus = "PENDING"
	KYCStatusVerified   KYCStatus = "VERIFIED"
	KYCStatusRejected   KYCStatus = "REJECTED"
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

// TokenPair represents the minted access and refresh token pair (ADR-014).
type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	TokenType        string
	ExpiresInSeconds int
}

// OTPPurpose represents the business context for the OTP challenge.
type OTPPurpose string

const (
	OTPPurposeRegistration  OTPPurpose = "REGISTRATION"
	OTPPurposeLogin         OTPPurpose = "LOGIN"
	OTPPurposePasswordReset OTPPurpose = "PASSWORD_RESET"
	OTPPurposeStepUp        OTPPurpose = "STEP_UP"
)

// OTPChallenge represents an in-flight verification challenge (AUT-01 to AUT-03).
type OTPChallenge struct {
	SessionID string
	Phone     string
	Code      string
	Purpose   OTPPurpose
	Attempts  int
	CreatedAt time.Time
	ExpiresAt time.Time
}
