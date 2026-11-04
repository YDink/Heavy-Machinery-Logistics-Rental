package domain

type AccountStatus string

const (
	AccountStatusUnverified              AccountStatus = "UNVERIFIED"
	AccountStatusActive                  AccountStatus = "ACTIVE"
	AccountStatusSuspendedSecurityFreeze AccountStatus = "SUSPENDED_SECURITY_FREEZE"
	AccountStatusDelinquentLocked        AccountStatus = "DELINQUENT_LOCKED"
	AccountStatusTerminated              AccountStatus = "TERMINATED"
)
