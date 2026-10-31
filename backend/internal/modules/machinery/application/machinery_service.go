package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
)

type CreateFleetPhotoInput struct {
	Angle      domain.PhotoAngle
	StorageURI string
	FileHash   string
}

type CreateFleetSpecValueInput struct {
	FieldKey     string
	ValueString  *string
	ValueNumeric *float64
	ValueBoolean *bool
}

type CreateFleetInput struct {
	OwnerID             uuid.UUID
	SubcategoryID       uuid.UUID
	Make                string
	Model               string
	ManufactureYear     int
	PlateNumber         string
	ChassisVIN          string
	BaselineMeterHours  float64
	FuelType            domain.FuelType
	SelfMoveCapable     bool
	DailyRateMinor      int64
	ExcessHourRateMinor int64
	MinRentalDays       int
	PaymentTimeoutHours int
	OperatorIncluded    bool
	Latitude            float64
	Longitude           float64
	AddressText         string
	Photos              []CreateFleetPhotoInput
	SpecValues          []CreateFleetSpecValueInput
}

type CreateFleetOutput struct {
	FleetID           uuid.UUID
	Status            domain.FleetStatus
	LogisticsIndexed  bool
	InferredTransport string
}

type MachineryService struct {
	categoryRepo domain.CategoryRepository
	specDefRepo  domain.SpecDefinitionRepository
	fleetRepo    domain.FleetRepository
	mediaRepo    domain.FleetMediaRepository
	specRepo     domain.FleetSpecRepository
}

func NewMachineryService(
	categoryRepo domain.CategoryRepository,
	specDefRepo domain.SpecDefinitionRepository,
	fleetRepo domain.FleetRepository,
	mediaRepo domain.FleetMediaRepository,
	specRepo domain.FleetSpecRepository,
) *MachineryService {
	return &MachineryService{
		categoryRepo: categoryRepo,
		specDefRepo:  specDefRepo,
		fleetRepo:    fleetRepo,
		mediaRepo:    mediaRepo,
		specRepo:     specRepo,
	}
}

// CreateFleet registers a new machinery asset with photos, dynamic EAV specs, and logistics indexing.
func (s *MachineryService) CreateFleet(ctx context.Context, input CreateFleetInput) (*CreateFleetOutput, error) {
	// 1. Invariant: Validate baseline meter reading (FL-04)
	if err := domain.ValidateBaselineMeter(input.BaselineMeterHours); err != nil {
		return nil, err
	}

	// 2. Invariant: Validate mandatory 6-angle walkaround photos (FL-05)
	var photos []domain.FleetPhoto
	for i, p := range input.Photos {
		photos = append(photos, domain.FleetPhoto{
			Angle:        p.Angle,
			StorageURI:   p.StorageURI,
			FileHash:     p.FileHash,
			DisplayOrder: i,
		})
	}
	if err := domain.ValidateWalkaroundPhotos(photos); err != nil {
		return nil, err
	}

	// 3. Verify subcategory existence
	category, err := s.categoryRepo.GetCategoryByID(ctx, input.SubcategoryID)
	if err != nil || category == nil {
		return nil, domain.ErrCategoryNotFound
	}

	// 4. Construct and persist core Fleet entity
	fleet := &domain.FleetAsset{
		OwnerID:             input.OwnerID,
		SubcategoryID:       input.SubcategoryID,
		Make:                input.Make,
		Model:               input.Model,
		ManufactureYear:     input.ManufactureYear,
		PlateNumber:         input.PlateNumber,
		ChassisVIN:          input.ChassisVIN,
		BaselineMeterHours:  input.BaselineMeterHours,
		FuelType:            input.FuelType,
		SelfMoveCapable:     input.SelfMoveCapable,
		DailyRateMinor:      input.DailyRateMinor,
		ExcessHourRateMinor: input.ExcessHourRateMinor,
		MinRentalDays:       input.MinRentalDays,
		PaymentTimeoutHours: input.PaymentTimeoutHours,
		OperatorIncluded:    input.OperatorIncluded,
		Status:              domain.FleetStatusPendingApproval,
		Latitude:            input.Latitude,
		Longitude:           input.Longitude,
		AddressText:         input.AddressText,
	}

	createdFleet, err := s.fleetRepo.CreateFleet(ctx, fleet)
	if err != nil {
		return nil, fmt.Errorf("failed to commit fleet asset: %w", err)
	}

	// 5. Commit walkaround photos
	for i := range photos {
		photos[i].FleetID = createdFleet.ID
	}
	if err := s.mediaRepo.SavePhotos(ctx, createdFleet.ID, photos); err != nil {
		return nil, fmt.Errorf("failed to save fleet photos: %w", err)
	}

	// 6. Resolve dynamic spec definitions and commit values
	specGroups, err := s.specDefRepo.GetCategorySpecDefinitions(ctx, input.SubcategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve spec definitions: %w", err)
	}

	defMap := make(map[string]domain.SpecFieldDefinition)
	for _, group := range specGroups {
		for _, field := range group.Fields {
			defMap[field.FieldKey] = field
		}
	}

	var specValues []domain.FleetSpecValue
	var indexedSpec domain.FleetSpecIndexed
	indexedSpec.FleetID = createdFleet.ID

	for _, sv := range input.SpecValues {
		def, exists := defMap[sv.FieldKey]
		if !exists {
			continue
		}

		specValues = append(specValues, domain.FleetSpecValue{
			FleetID:      createdFleet.ID,
			SpecID:       def.ID,
			FieldKey:     sv.FieldKey,
			ValueString:  sv.ValueString,
			ValueNumeric: sv.ValueNumeric,
			ValueBoolean: sv.ValueBoolean,
		})

		// Dual-Persistence: Extract Logistics-Critical Dimensions (FL-14, D-05-02-002)
		if def.IsLogisticsCritical && sv.ValueNumeric != nil {
			switch sv.FieldKey {
			case "operating_weight_kg":
				indexedSpec.OperatingWeightKg = sv.ValueNumeric
			case "transport_length_mm":
				indexedSpec.TransportLengthMm = sv.ValueNumeric
			case "transport_width_mm":
				indexedSpec.TransportWidthMm = sv.ValueNumeric
			case "transport_height_mm":
				indexedSpec.TransportHeightMm = sv.ValueNumeric
			case "payload_capacity_kg":
				indexedSpec.PayloadCapacityKg = sv.ValueNumeric
			case "gross_power_hp":
				indexedSpec.GrossEnginePowerHp = sv.ValueNumeric
			}
		}
	}

	if len(specValues) > 0 {
		if err := s.specRepo.SaveSpecValues(ctx, createdFleet.ID, specValues); err != nil {
			return nil, fmt.Errorf("failed to commit dynamic specs: %w", err)
		}
	}

	// Commit logistics performance index
	_ = s.specRepo.UpsertSpecIndexed(ctx, &indexedSpec)

	inferredClass := ""
	if indexedSpec.OperatingWeightKg != nil {
		inferredClass = domain.InferLogisticsTransportClass(*indexedSpec.OperatingWeightKg)
	}

	return &CreateFleetOutput{
		FleetID:           createdFleet.ID,
		Status:            createdFleet.Status,
		LogisticsIndexed:  true,
		InferredTransport: inferredClass,
	}, nil
}

// ListCategories returns active taxonomy nodes.
func (s *MachineryService) ListCategories(ctx context.Context) ([]domain.Category, error) {
	return s.categoryRepo.ListAll(ctx)
}

// GetCategorySpecDefinitions returns the EAV field schema for dynamic mobile form compilation.
func (s *MachineryService) GetCategorySpecDefinitions(ctx context.Context, categoryID uuid.UUID) ([]domain.SpecFieldGroup, error) {
	return s.specDefRepo.GetCategorySpecDefinitions(ctx, categoryID)
}

// GetFleetAsset fetches full details for a registered machine.
func (s *MachineryService) GetFleetAsset(ctx context.Context, fleetID uuid.UUID) (*domain.FleetAsset, error) {
	fleet, err := s.fleetRepo.GetFleetByID(ctx, fleetID)
	if err != nil {
		return nil, err
	}
	photos, _ := s.mediaRepo.ListPhotos(ctx, fleetID)
	fleet.Photos = photos

	specIndexed, _ := s.specRepo.GetSpecIndexed(ctx, fleetID)
	fleet.LogisticsSpec = specIndexed

	return fleet, nil
}
