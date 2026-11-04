package domain

type FleetStatus string

const (
	FleetStatusDraft            FleetStatus = "DRAFT"
	FleetStatusPendingApproval  FleetStatus = "PENDING_APPROVAL"
	FleetStatusActive           FleetStatus = "ACTIVE"
	FleetStatusRented           FleetStatus = "RENTED"
	FleetStatusInactive         FleetStatus = "INACTIVE"
	FleetStatusInsuranceHold    FleetStatus = "INSURANCE_HOLD"
	FleetStatusDelinquentLocked FleetStatus = "DELINQUENT_LOCKED"
	FleetStatusDecommissioned   FleetStatus = "DECOMMISSIONED"
)
