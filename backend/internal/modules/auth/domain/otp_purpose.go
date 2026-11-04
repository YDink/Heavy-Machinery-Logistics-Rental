package domain

// OTPPurpose represents the business context for the OTP challenge.
type OTPPurpose string

const (
	OTPPurposeRegistration  OTPPurpose = "REGISTRATION"
	OTPPurposeLogin         OTPPurpose = "LOGIN"
	OTPPurposePasswordReset OTPPurpose = "PASSWORD_RESET"
	OTPPurposeStepUp        OTPPurpose = "STEP_UP"
)
