package domain

type BankSlipStatus string

const (
	SlipStatusPendingAudit BankSlipStatus = "PENDING_AUDIT"
	SlipStatusVerified     BankSlipStatus = "VERIFIED"
	SlipStatusRejected     BankSlipStatus = "REJECTED"
)
