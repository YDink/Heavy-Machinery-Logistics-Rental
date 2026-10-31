package application_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/dinkmachinery/dink-core/internal/modules/machinery/application"
	"github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
	"github.com/dinkmachinery/dink-core/internal/modules/machinery/infrastructure"
)

func setupTestMachineryService() (*application.MachineryService, *infrastructure.MemoryMachineryRepository, uuid.UUID) {
	repo := infrastructure.NewMemoryMachineryRepository()
	categoryID := uuid.New()

	// Seed test category
	_ = repo.SeedCategory(domain.Category{
		ID:           categoryID,
		NameEN:       "Crawler Excavator",
		NameAM:       "የክራውለር ኤክስካቫተር",
		Code:         "EXC-CRW",
		IsActive:     true,
		DisplayOrder: 1,
	})

	unit := "kg"
	// Seed test spec group & field definition
	_ = repo.SeedSpecDefinitions(categoryID, []domain.SpecFieldGroup{
		{
			ID:            uuid.New(),
			CategoryID:    categoryID,
			SectionNameEN: "Dimensions & Weight",
			SectionNameAM: "የክብደትና ስፋት ዝርዝሮች",
			DisplayOrder:  1,
			Fields: []domain.SpecFieldDefinition{
				{
					ID:                  uuid.New(),
					FieldKey:            "operating_weight_kg",
					LabelEN:             "Operating Weight",
					LabelAM:             "የስራ ክብደት",
					FieldType:           "DECIMAL",
					UnitSymbol:          &unit,
					IsRequired:          true,
					IsLogisticsCritical: true,
				},
			},
		},
	})

	svc := application.NewMachineryService(repo, repo, repo, repo, repo)
	return svc, repo, categoryID
}

func complete6Photos() []application.CreateFleetPhotoInput {
	return []application.CreateFleetPhotoInput{
		{Angle: domain.PhotoAngleFront, StorageURI: "s3://fleet/front.jpg", FileHash: "hash1"},
		{Angle: domain.PhotoAngleRear, StorageURI: "s3://fleet/rear.jpg", FileHash: "hash2"},
		{Angle: domain.PhotoAngleLeft, StorageURI: "s3://fleet/left.jpg", FileHash: "hash3"},
		{Angle: domain.PhotoAngleRight, StorageURI: "s3://fleet/right.jpg", FileHash: "hash4"},
		{Angle: domain.PhotoAngleCabin, StorageURI: "s3://fleet/cabin.jpg", FileHash: "hash5"},
		{Angle: domain.PhotoAngleMeter, StorageURI: "s3://fleet/meter.jpg", FileHash: "hash6"},
	}
}

func TestMachineryService_CreateFleet_Success(t *testing.T) {
	svc, _, categoryID := setupTestMachineryService()
	ctx := context.Background()

	weightKg := 18500.0 // 18.5 tons -> CLASS_2_RIGID_3AXLE_24T
	powerHp := 150.0

	input := application.CreateFleetInput{
		OwnerID:             uuid.New(),
		SubcategoryID:       categoryID,
		Make:                "Caterpillar",
		Model:               "320D",
		ManufactureYear:     2022,
		PlateNumber:         "ET-3-45678",
		ChassisVIN:          "CAT320D0099887766",
		BaselineMeterHours:  1250.5,
		FuelType:            domain.FuelTypeDiesel,
		SelfMoveCapable:     false,
		DailyRateMinor:      4500000, // 45,000.00 ETB
		ExcessHourRateMinor: 500000,
		MinRentalDays:       3,
		PaymentTimeoutHours: 24,
		OperatorIncluded:    true,
		Latitude:            9.0123,
		Longitude:           38.7654,
		AddressText:         "Bole Sub-City, Addis Ababa",
		Photos:              complete6Photos(),
		SpecValues: []application.CreateFleetSpecValueInput{
			{
				FieldKey:     "operating_weight_kg",
				ValueNumeric: &weightKg,
			},
			{
				FieldKey:     "gross_engine_power_hp",
				ValueNumeric: &powerHp,
			},
		},
	}

	out, err := svc.CreateFleet(ctx, input)
	if err != nil {
		t.Fatalf("expected successful fleet creation, got error: %v", err)
	}

	if out.FleetID == uuid.Nil {
		t.Errorf("expected non-nil FleetID")
	}

	if out.Status != domain.FleetStatusPendingApproval {
		t.Errorf("expected initial status PENDING_APPROVAL, got %s", out.Status)
	}

	if !out.LogisticsIndexed {
		t.Errorf("expected LogisticsIndexed to be true")
	}

	expectedClass := "Two-Axle Lowboy Trailer (25t Payload Capacity)"
	if out.InferredTransport != expectedClass {
		t.Errorf("expected InferredTransport %s, got %s", expectedClass, out.InferredTransport)
	}

	// Verify retrieval through service
	asset, err := svc.GetFleetAsset(ctx, out.FleetID)
	if err != nil {
		t.Fatalf("failed to retrieve created fleet asset: %v", err)
	}

	if len(asset.Photos) != 6 {
		t.Errorf("expected 6 walkaround photos attached, got %d", len(asset.Photos))
	}

	if asset.LogisticsSpec == nil {
		t.Fatalf("expected non-nil LogisticsSpec")
	}

	if asset.LogisticsSpec.OperatingWeightKg == nil || *asset.LogisticsSpec.OperatingWeightKg != weightKg {
		t.Errorf("expected LogisticsSpec operating weight %v, got %v", weightKg, asset.LogisticsSpec.OperatingWeightKg)
	}
}

func TestMachineryService_CreateFleet_InvalidBaselineMeter(t *testing.T) {
	svc, _, categoryID := setupTestMachineryService()
	ctx := context.Background()

	input := application.CreateFleetInput{
		OwnerID:            uuid.New(),
		SubcategoryID:      categoryID,
		BaselineMeterHours: -5.0, // Invalid negative reading
		Photos:             complete6Photos(),
	}

	_, err := svc.CreateFleet(ctx, input)
	if err != domain.ErrInvalidBaselineMeter {
		t.Errorf("expected ErrInvalidBaselineMeter, got %v", err)
	}
}

func TestMachineryService_CreateFleet_IncompleteWalkaroundPhotos(t *testing.T) {
	svc, _, categoryID := setupTestMachineryService()
	ctx := context.Background()

	// Only 2 photos provided
	photos := []application.CreateFleetPhotoInput{
		{Angle: domain.PhotoAngleFront, StorageURI: "s3://fleet/front.jpg"},
		{Angle: domain.PhotoAngleRear, StorageURI: "s3://fleet/rear.jpg"},
	}

	input := application.CreateFleetInput{
		OwnerID:            uuid.New(),
		SubcategoryID:      categoryID,
		BaselineMeterHours: 100.0,
		Photos:             photos,
	}

	_, err := svc.CreateFleet(ctx, input)
	if err != domain.ErrIncompletePhotoWalkaround {
		t.Errorf("expected ErrIncompletePhotoWalkaround, got %v", err)
	}
}

func TestMachineryService_CreateFleet_CategoryNotFound(t *testing.T) {
	svc, _, _ := setupTestMachineryService()
	ctx := context.Background()

	nonExistentCategoryID := uuid.New()

	input := application.CreateFleetInput{
		OwnerID:            uuid.New(),
		SubcategoryID:      nonExistentCategoryID,
		BaselineMeterHours: 100.0,
		Photos:             complete6Photos(),
	}

	_, err := svc.CreateFleet(ctx, input)
	if err != domain.ErrCategoryNotFound {
		t.Errorf("expected ErrCategoryNotFound, got %v", err)
	}
}

func TestMachineryService_TaxonomyAndSpecQueries(t *testing.T) {
	svc, _, categoryID := setupTestMachineryService()
	ctx := context.Background()

	categories, err := svc.ListCategories(ctx)
	if err != nil {
		t.Fatalf("unexpected error listing categories: %v", err)
	}
	if len(categories) == 0 {
		t.Errorf("expected at least 1 category, got %d", len(categories))
	}

	specGroups, err := svc.GetCategorySpecDefinitions(ctx, categoryID)
	if err != nil {
		t.Fatalf("unexpected error getting spec defs: %v", err)
	}
	if len(specGroups) != 1 {
		t.Errorf("expected 1 spec group, got %d", len(specGroups))
	}
	if len(specGroups[0].Fields) != 1 {
		t.Errorf("expected 1 field definition in group, got %d", len(specGroups[0].Fields))
	}
}
