package domain

import (
	"time"

	"github.com/google/uuid"
)

// FleetAsset represents the primary machinery asset entity.
type FleetAsset struct {
	ID                  uuid.UUID
	OwnerID             uuid.UUID
	SubcategoryID       uuid.UUID
	Make                string
	Model               string
	ManufactureYear     int
	PlateNumber         string
	ChassisVIN          string
	BaselineMeterHours  float64
	FuelType            FuelType
	SelfMoveCapable     bool
	DailyRateMinor      int64
	ExcessHourRateMinor int64
	MinRentalDays       int
	PaymentTimeoutHours int
	OperatorIncluded    bool
	Status              FleetStatus
	Latitude            float64
	Longitude           float64
	AddressText         string
	Version             int
	CreatedAt           time.Time
	UpdatedAt           time.Time

	// Attached Collections
	Photos        []FleetPhoto
	SpecValues    []FleetSpecValue
	LogisticsSpec *FleetSpecIndexed
}
