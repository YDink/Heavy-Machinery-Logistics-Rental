package domain

type KYCStatus string

const (
	KYCStatusUnverified KYCStatus = "UNVERIFIED"
	KYCStatusPending    KYCStatus = "PENDING"
	KYCStatusVerified   KYCStatus = "VERIFIED"
	KYCStatusRejected   KYCStatus = "REJECTED"
)
