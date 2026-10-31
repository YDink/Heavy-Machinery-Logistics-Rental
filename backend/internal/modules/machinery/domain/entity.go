package domain

import (
	"time"

	"github.com/google/uuid"
)

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

type FuelType string

const (
	FuelTypeDiesel    FuelType = "DIESEL"
	FuelTypeElectric  FuelType = "ELECTRIC"
	FuelTypePetrol    FuelType = "PETROL"
	FuelTypeHybrid    FuelType = "HYBRID"
	FuelTypeBiodiesel FuelType = "BIODIESEL"
)

type PhotoAngle string

const (
	PhotoAngleFront PhotoAngle = "FRONT"
	PhotoAngleRear  PhotoAngle = "REAR"
	PhotoAngleLeft  PhotoAngle = "LEFT"
	PhotoAngleRight PhotoAngle = "RIGHT"
	PhotoAngleCabin PhotoAngle = "CABIN"
	PhotoAngleMeter PhotoAngle = "METER"
)

var MandatoryPhotoAngles = []PhotoAngle{
	PhotoAngleFront,
	PhotoAngleRear,
	PhotoAngleLeft,
	PhotoAngleRight,
	PhotoAngleCabin,
	PhotoAngleMeter,
}

// Category represents a node in the hierarchical machinery taxonomy.
type Category struct {
	ID           uuid.UUID
	ParentID     *uuid.UUID
	Code         string
	NameEN       string
	NameAM       string
	IconURL      string
	IsActive     bool
	DisplayOrder int
}

// SpecFieldGroup groups related technical specification attributes.
type SpecFieldGroup struct {
	ID             uuid.UUID
	CategoryID     uuid.UUID
	SectionNameEN  string
	SectionNameAM  string
	DisplayOrder   int
	Fields         []SpecFieldDefinition
}

// SpecFieldDefinition defines dynamic EAV schema metadata.
type SpecFieldDefinition struct {
	ID                 uuid.UUID
	GroupID            uuid.UUID
	FieldKey           string
	LabelEN            string
	LabelAM            string
	FieldType          string // 'STRING', 'DECIMAL', 'INTEGER', 'ENUM', 'BOOLEAN'
	UnitSymbol         *string
	IsRequired         bool
	IsLogisticsCritical bool
	IsComparable       bool
	IsVisibleToRenter  bool
	MinValue           *float64
	MaxValue           *float64
	EnumValuesJSON     []string
	DisplayOrder       int
}

// FleetPhoto represents a walkaround photo of the asset.
type FleetPhoto struct {
	ID           uuid.UUID
	FleetID      uuid.UUID
	Angle        PhotoAngle
	StorageURI   string
	FileHash     string
	DisplayOrder int
	CreatedAt    time.Time
}

// FleetSpecValue stores dynamic EAV specifications for an asset.
type FleetSpecValue struct {
	ID               uuid.UUID
	FleetID          uuid.UUID
	SpecID           uuid.UUID
	FieldKey         string
	ValueString      *string
	ValueNumeric     *float64
	ValueBoolean     *bool
}

// FleetSpecIndexed represents denormalized logistics-critical dimensions.
type FleetSpecIndexed struct {
	FleetID            uuid.UUID
	OperatingWeightKg  *float64
	TransportLengthMm  *float64
	TransportWidthMm   *float64
	TransportHeightMm  *float64
	PayloadCapacityKg  *float64
	GrossEnginePowerHp *float64
	UpdatedAt          time.Time
}

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
	Photos              []FleetPhoto
	SpecValues          []FleetSpecValue
	LogisticsSpec       *FleetSpecIndexed
}

// ValidateBaselineMeter asserts that the cabin mechanical engine hour meter reading is non-negative.
func ValidateBaselineMeter(hours float64) error {
	if hours < 0.0 {
		return ErrInvalidBaselineMeter
	}
	return nil
}

// ValidateWalkaroundPhotos asserts that all 6 required physical angles are present (FL-05).
func ValidateWalkaroundPhotos(photos []FleetPhoto) error {
	anglePresent := make(map[PhotoAngle]bool)
	for _, p := range photos {
		anglePresent[p.Angle] = true
	}

	for _, required := range MandatoryPhotoAngles {
		if !anglePresent[required] {
			return ErrIncompletePhotoWalkaround
		}
	}
	return nil
}

// InferLogisticsTransportClass automatically calculates required trailer class per SRS §7.2.
func InferLogisticsTransportClass(operatingWeightKg float64) string {
	if operatingWeightKg < 12000 {
		return "Two-Axle Tilt-Bed / Rigid Flatbed Carrier"
	} else if operatingWeightKg < 24000 {
		return "Two-Axle Lowboy Trailer (25t Payload Capacity)"
	} else if operatingWeightKg < 40000 {
		return "Tri-Axle Heavy Lowboy Trailer (40t Payload Capacity)"
	}
	return "Multi-Axle Modular Hydraulic Platform Trailer (> 60t Capacity)"
}
