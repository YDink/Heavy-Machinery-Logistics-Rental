package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/database/db"
)

type PostgresMachineryRepository struct {
	queries *db.Queries
}

func NewPostgresMachineryRepository(queries *db.Queries) *PostgresMachineryRepository {
	return &PostgresMachineryRepository{queries: queries}
}

// ==================== Category Repository ====================

func toDomainCategory(c db.Category) domain.Category {
	var parentID *uuid.UUID
	if c.ParentID.Valid {
		id := uuid.UUID(c.ParentID.Bytes)
		parentID = &id
	}
	return domain.Category{
		ID:           uuid.UUID(c.ID.Bytes),
		ParentID:     parentID,
		Code:         c.Code,
		NameEN:       c.NameEn,
		NameAM:       c.NameAm,
		IconURL:      c.IconUrl,
		IsActive:     c.IsActive,
		DisplayOrder: int(c.DisplayOrder),
	}
}

func (r *PostgresMachineryRepository) ListAll(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.queries.ListAllCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	cats := make([]domain.Category, len(rows))
	for i, row := range rows {
		cats[i] = toDomainCategory(row)
	}
	return cats, nil
}

func (r *PostgresMachineryRepository) ListRoots(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.queries.ListRootCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list root categories: %w", err)
	}
	cats := make([]domain.Category, len(rows))
	for i, row := range rows {
		cats[i] = toDomainCategory(row)
	}
	return cats, nil
}

func (r *PostgresMachineryRepository) ListSubcategories(ctx context.Context, parentID uuid.UUID) ([]domain.Category, error) {
	pgID := pgtype.UUID{Bytes: parentID, Valid: true}
	rows, err := r.queries.ListSubcategories(ctx, pgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list subcategories: %w", err)
	}
	cats := make([]domain.Category, len(rows))
	for i, row := range rows {
		cats[i] = toDomainCategory(row)
	}
	return cats, nil
}

func (r *PostgresMachineryRepository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	pgID := pgtype.UUID{Bytes: id, Valid: true}
	row, err := r.queries.GetCategoryByID(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("query GetCategoryByID failed: %w", err)
	}
	c := toDomainCategory(row)
	return &c, nil
}

// ==================== Spec Definition Repository ====================

func (r *PostgresMachineryRepository) GetCategorySpecDefinitions(ctx context.Context, categoryID uuid.UUID) ([]domain.SpecFieldGroup, error) {
	pgID := pgtype.UUID{Bytes: categoryID, Valid: true}
	rows, err := r.queries.GetCategorySpecDefinitions(ctx, pgID)
	if err != nil {
		return nil, fmt.Errorf("failed to query spec definitions: %w", err)
	}

	groupsMap := make(map[uuid.UUID]*domain.SpecFieldGroup)
	var groupOrder []uuid.UUID

	for _, row := range rows {
		gID := uuid.UUID(row.GroupID.Bytes)
		group, exists := groupsMap[gID]
		if !exists {
			group = &domain.SpecFieldGroup{
				ID:            gID,
				CategoryID:    categoryID,
				SectionNameEN: row.SectionNameEn,
				SectionNameAM: row.SectionNameAm,
				DisplayOrder:  int(row.GroupDisplayOrder),
				Fields:        []domain.SpecFieldDefinition{},
			}
			groupsMap[gID] = group
			groupOrder = append(groupOrder, gID)
		}

		var minVal, maxVal *float64
		if row.MinValue.Valid {
			f, err := row.MinValue.Float64Value()
			if err == nil {
				v := f.Float64
				minVal = &v
			}
		}
		if row.MaxValue.Valid {
			f, err := row.MaxValue.Float64Value()
			if err == nil {
				v := f.Float64
				maxVal = &v
			}
		}

		var enumVals []string
		if len(row.EnumValuesJson) > 0 {
			_ = json.Unmarshal(row.EnumValuesJson, &enumVals)
		}

		field := domain.SpecFieldDefinition{
			ID:                  uuid.UUID(row.DefinitionID.Bytes),
			GroupID:             gID,
			FieldKey:            row.FieldKey,
			LabelEN:             row.LabelEn,
			LabelAM:             row.LabelAm,
			FieldType:           row.FieldType,
			UnitSymbol:          row.UnitSymbol,
			IsRequired:          row.IsRequired,
			IsLogisticsCritical: row.IsLogisticsCritical,
			IsComparable:        row.IsComparable,
			IsVisibleToRenter:  row.IsVisibleToRenter,
			MinValue:            minVal,
			MaxValue:            maxVal,
			EnumValuesJSON:     enumVals,
			DisplayOrder:       int(row.DefinitionDisplayOrder),
		}

		group.Fields = append(group.Fields, field)
	}

	groups := make([]domain.SpecFieldGroup, len(groupOrder))
	for i, id := range groupOrder {
		groups[i] = *groupsMap[id]
	}
	return groups, nil
}

// ==================== Fleet Repository ====================

func toNumeric(val float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(fmt.Sprintf("%.2f", val))
	return n
}

func fromNumeric(n pgtype.Numeric) *float64 {
	if !n.Valid {
		return nil
	}
	f, err := n.Float64Value()
	if err != nil {
		return nil
	}
	v := f.Float64
	return &v
}

func toFloat(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int64:
		return float64(v)
	case int:
		return float64(v)
	default:
		return 0.0
	}
}

func (r *PostgresMachineryRepository) CreateFleet(ctx context.Context, f *domain.FleetAsset) (*domain.FleetAsset, error) {
	params := db.CreateFleetParams{
		OwnerID:             pgtype.UUID{Bytes: f.OwnerID, Valid: true},
		SubcategoryID:       pgtype.UUID{Bytes: f.SubcategoryID, Valid: true},
		Make:                f.Make,
		Model:               f.Model,
		ManufactureYear:     int32(f.ManufactureYear),
		PlateNumber:         f.PlateNumber,
		ChassisVin:          f.ChassisVIN,
		BaselineMeterHours:  toNumeric(f.BaselineMeterHours),
		FuelType:            string(f.FuelType),
		SelfMoveCapable:     f.SelfMoveCapable,
		DailyRateMinor:      f.DailyRateMinor,
		ExcessHourRateMinor: f.ExcessHourRateMinor,
		MinRentalDays:       int32(f.MinRentalDays),
		PaymentTimeoutHours: int32(f.PaymentTimeoutHours),
		OperatorIncluded:    f.OperatorIncluded,
		Status:              string(f.Status),
		StMakepoint:         f.Longitude,
		StMakepoint_2:       f.Latitude,
		AddressText:         f.AddressText,
	}

	row, err := r.queries.CreateFleet(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to insert fleet: %w", err)
	}

	meterVal := 0.0
	if m := fromNumeric(row.BaselineMeterHours); m != nil {
		meterVal = *m
	}

	f.ID = uuid.UUID(row.ID.Bytes)
	f.BaselineMeterHours = meterVal
	f.Latitude = toFloat(row.Latitude)
	f.Longitude = toFloat(row.Longitude)
	f.Version = int(row.Version)
	f.CreatedAt = row.CreatedAt.Time
	f.UpdatedAt = row.UpdatedAt.Time

	return f, nil
}

func (r *PostgresMachineryRepository) GetFleetByID(ctx context.Context, id uuid.UUID) (*domain.FleetAsset, error) {
	pgID := pgtype.UUID{Bytes: id, Valid: true}
	row, err := r.queries.GetFleetByID(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrFleetNotFound
		}
		return nil, fmt.Errorf("query GetFleetByID failed: %w", err)
	}

	meterVal := 0.0
	if m := fromNumeric(row.BaselineMeterHours); m != nil {
		meterVal = *m
	}

	return &domain.FleetAsset{
		ID:                  uuid.UUID(row.ID.Bytes),
		OwnerID:             uuid.UUID(row.OwnerID.Bytes),
		SubcategoryID:       uuid.UUID(row.SubcategoryID.Bytes),
		Make:                row.Make,
		Model:               row.Model,
		ManufactureYear:     int(row.ManufactureYear),
		PlateNumber:         row.PlateNumber,
		ChassisVIN:          row.ChassisVin,
		BaselineMeterHours:  meterVal,
		FuelType:            domain.FuelType(row.FuelType),
		SelfMoveCapable:     row.SelfMoveCapable,
		DailyRateMinor:      row.DailyRateMinor,
		ExcessHourRateMinor: row.ExcessHourRateMinor,
		MinRentalDays:       int(row.MinRentalDays),
		PaymentTimeoutHours: int(row.PaymentTimeoutHours),
		OperatorIncluded:    row.OperatorIncluded,
		Status:              domain.FleetStatus(row.Status),
		Latitude:            toFloat(row.Latitude),
		Longitude:           toFloat(row.Longitude),
		AddressText:         row.AddressText,
		Version:             int(row.Version),
		CreatedAt:           row.CreatedAt.Time,
		UpdatedAt:           row.UpdatedAt.Time,
	}, nil
}

func (r *PostgresMachineryRepository) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]domain.FleetAsset, error) {
	pgID := pgtype.UUID{Bytes: ownerID, Valid: true}
	rows, err := r.queries.ListFleetsByOwner(ctx, pgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list fleets by owner: %w", err)
	}

	assets := make([]domain.FleetAsset, len(rows))
	for i, row := range rows {
		meterVal := 0.0
		if m := fromNumeric(row.BaselineMeterHours); m != nil {
			meterVal = *m
		}
		assets[i] = domain.FleetAsset{
			ID:                  uuid.UUID(row.ID.Bytes),
			OwnerID:             uuid.UUID(row.OwnerID.Bytes),
			SubcategoryID:       uuid.UUID(row.SubcategoryID.Bytes),
			Make:                row.Make,
			Model:               row.Model,
			ManufactureYear:     int(row.ManufactureYear),
			PlateNumber:         row.PlateNumber,
			ChassisVIN:          row.ChassisVin,
			BaselineMeterHours:  meterVal,
			FuelType:            domain.FuelType(row.FuelType),
			SelfMoveCapable:     row.SelfMoveCapable,
			DailyRateMinor:      row.DailyRateMinor,
			ExcessHourRateMinor: row.ExcessHourRateMinor,
			MinRentalDays:       int(row.MinRentalDays),
			PaymentTimeoutHours: int(row.PaymentTimeoutHours),
			OperatorIncluded:    row.OperatorIncluded,
			Status:              domain.FleetStatus(row.Status),
			Latitude:            toFloat(row.Latitude),
			Longitude:           toFloat(row.Longitude),
			AddressText:         row.AddressText,
			Version:             int(row.Version),
			CreatedAt:           row.CreatedAt.Time,
			UpdatedAt:           row.UpdatedAt.Time,
		}
	}
	return assets, nil
}

func (r *PostgresMachineryRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.FleetStatus) (*domain.FleetAsset, error) {
	pgID := pgtype.UUID{Bytes: id, Valid: true}
	row, err := r.queries.UpdateFleetStatus(ctx, db.UpdateFleetStatusParams{
		ID:     pgID,
		Status: string(status),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrFleetNotFound
		}
		return nil, fmt.Errorf("failed to update fleet status: %w", err)
	}

	meterVal := 0.0
	if m := fromNumeric(row.BaselineMeterHours); m != nil {
		meterVal = *m
	}

	return &domain.FleetAsset{
		ID:                  uuid.UUID(row.ID.Bytes),
		OwnerID:             uuid.UUID(row.OwnerID.Bytes),
		SubcategoryID:       uuid.UUID(row.SubcategoryID.Bytes),
		Make:                row.Make,
		Model:               row.Model,
		ManufactureYear:     int(row.ManufactureYear),
		PlateNumber:         row.PlateNumber,
		ChassisVIN:          row.ChassisVin,
		BaselineMeterHours:  meterVal,
		FuelType:            domain.FuelType(row.FuelType),
		SelfMoveCapable:     row.SelfMoveCapable,
		DailyRateMinor:      row.DailyRateMinor,
		ExcessHourRateMinor: row.ExcessHourRateMinor,
		MinRentalDays:       int(row.MinRentalDays),
		PaymentTimeoutHours: int(row.PaymentTimeoutHours),
		OperatorIncluded:    row.OperatorIncluded,
		Status:              domain.FleetStatus(row.Status),
		AddressText:         row.AddressText,
		Version:             int(row.Version),
		CreatedAt:           row.CreatedAt.Time,
		UpdatedAt:           row.UpdatedAt.Time,
	}, nil
}

// ==================== Fleet Media Repository ====================

func (r *PostgresMachineryRepository) SavePhotos(ctx context.Context, fleetID uuid.UUID, photos []domain.FleetPhoto) error {
	pgFleetID := pgtype.UUID{Bytes: fleetID, Valid: true}
	for i, photo := range photos {
		_, err := r.queries.CreateFleetPhoto(ctx, db.CreateFleetPhotoParams{
			FleetID:      pgFleetID,
			Angle:        string(photo.Angle),
			StorageUri:   photo.StorageURI,
			FileHash:     photo.FileHash,
			DisplayOrder: int32(i),
		})
		if err != nil {
			return fmt.Errorf("failed to persist photo %s: %w", photo.Angle, err)
		}
	}
	return nil
}

func (r *PostgresMachineryRepository) ListPhotos(ctx context.Context, fleetID uuid.UUID) ([]domain.FleetPhoto, error) {
	pgFleetID := pgtype.UUID{Bytes: fleetID, Valid: true}
	rows, err := r.queries.ListFleetPhotos(ctx, pgFleetID)
	if err != nil {
		return nil, fmt.Errorf("failed to list fleet photos: %w", err)
	}

	photos := make([]domain.FleetPhoto, len(rows))
	for i, row := range rows {
		photos[i] = domain.FleetPhoto{
			ID:           uuid.UUID(row.ID.Bytes),
			FleetID:      fleetID,
			Angle:        domain.PhotoAngle(row.Angle),
			StorageURI:   row.StorageUri,
			FileHash:     row.FileHash,
			DisplayOrder: int(row.DisplayOrder),
			CreatedAt:    row.CreatedAt.Time,
		}
	}
	return photos, nil
}

// ==================== Fleet Spec Repository ====================

func (r *PostgresMachineryRepository) SaveSpecValues(ctx context.Context, fleetID uuid.UUID, values []domain.FleetSpecValue) error {
	pgFleetID := pgtype.UUID{Bytes: fleetID, Valid: true}
	for _, val := range values {
		var numVal pgtype.Numeric
		if val.ValueNumeric != nil {
			numVal = toNumeric(*val.ValueNumeric)
		}
		_, err := r.queries.UpsertFleetSpecValue(ctx, db.UpsertFleetSpecValueParams{
			FleetID:          pgFleetID,
			SpecDefinitionID: pgtype.UUID{Bytes: val.SpecID, Valid: true},
			ValueString:      val.ValueString,
			ValueNumeric:     numVal,
			ValueBoolean:     val.ValueBoolean,
		})
		if err != nil {
			return fmt.Errorf("failed to upsert spec value %s: %w", val.FieldKey, err)
		}
	}
	return nil
}

func (r *PostgresMachineryRepository) UpsertSpecIndexed(ctx context.Context, indexed *domain.FleetSpecIndexed) error {
	pgFleetID := pgtype.UUID{Bytes: indexed.FleetID, Valid: true}
	var weight, length, width, height, payload, power pgtype.Numeric
	if indexed.OperatingWeightKg != nil {
		weight = toNumeric(*indexed.OperatingWeightKg)
	}
	if indexed.TransportLengthMm != nil {
		length = toNumeric(*indexed.TransportLengthMm)
	}
	if indexed.TransportWidthMm != nil {
		width = toNumeric(*indexed.TransportWidthMm)
	}
	if indexed.TransportHeightMm != nil {
		height = toNumeric(*indexed.TransportHeightMm)
	}
	if indexed.PayloadCapacityKg != nil {
		payload = toNumeric(*indexed.PayloadCapacityKg)
	}
	if indexed.GrossEnginePowerHp != nil {
		power = toNumeric(*indexed.GrossEnginePowerHp)
	}

	_, err := r.queries.UpsertFleetSpecIndexed(ctx, db.UpsertFleetSpecIndexedParams{
		FleetID:            pgFleetID,
		OperatingWeightKg:  weight,
		TransportLengthMm:  length,
		TransportWidthMm:   width,
		TransportHeightMm:  height,
		PayloadCapacityKg:  payload,
		GrossEnginePowerHp: power,
	})
	if err != nil {
		return fmt.Errorf("failed to upsert fleet spec indexed: %w", err)
	}
	return nil
}

func (r *PostgresMachineryRepository) GetSpecIndexed(ctx context.Context, fleetID uuid.UUID) (*domain.FleetSpecIndexed, error) {
	pgFleetID := pgtype.UUID{Bytes: fleetID, Valid: true}
	row, err := r.queries.GetFleetSpecIndexed(ctx, pgFleetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get fleet spec indexed: %w", err)
	}

	return &domain.FleetSpecIndexed{
		FleetID:            fleetID,
		OperatingWeightKg:  fromNumeric(row.OperatingWeightKg),
		TransportLengthMm:  fromNumeric(row.TransportLengthMm),
		TransportWidthMm:   fromNumeric(row.TransportWidthMm),
		TransportHeightMm:  fromNumeric(row.TransportHeightMm),
		PayloadCapacityKg:  fromNumeric(row.PayloadCapacityKg),
		GrossEnginePowerHp: fromNumeric(row.GrossEnginePowerHp),
		UpdatedAt:          row.UpdatedAt.Time,
	}, nil
}
