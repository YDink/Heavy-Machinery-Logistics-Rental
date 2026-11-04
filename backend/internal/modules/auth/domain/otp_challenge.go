package domain

import "time"

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
