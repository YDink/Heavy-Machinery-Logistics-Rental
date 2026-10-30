package domain

import "errors"

var (
	ErrInvalidPhoneNumber      = errors.New("invalid E.164 phone number")
	ErrUnsupportedCountryCode   = errors.New("unsupported country code")
	ErrUserNotFound             = errors.New("user not found")
	ErrUserAlreadyExists        = errors.New("user already registered with this phone number")
	ErrAccountSuspended         = errors.New("account is suspended or locked")
	ErrOTPExpired               = errors.New("verification code has expired")
	ErrOTPInvalid               = errors.New("invalid verification code")
	ErrOTPMaxAttemptsExceeded   = errors.New("maximum verification attempts exceeded")
	ErrRateLimitExceeded        = errors.New("rate limit exceeded; please wait before retrying")
	ErrInvalidToken             = errors.New("invalid or malformed token")
	ErrTokenExpired             = errors.New("token has expired")
	ErrTokenRevoked             = errors.New("token has been revoked")
	ErrTokenFamilyCompromised   = errors.New("refresh token replay detected; token family compromised")
	ErrInvalidCredentials       = errors.New("invalid phone number or password")
)
