package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

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

// ==================== Attachment Repository ====================

func fromNumericVal(n pgtype.Numeric) float64 {
	p := fromNumeric(n)
	if p == nil {
		return 0.0
	}
	return *p
}

func toDomainAttachment(a db.Attachment, categories []string, photos []domain.AttachmentPhoto) *domain.Attachment {
	return &domain.Attachment{
		ID:                         uuid.UUID(a.ID.Bytes),
		OwnerID:                    uuid.UUID(a.OwnerID.Bytes),
		AttachmentType:             domain.AttachmentType(a.AttachmentType),
		Make:                       a.Make,
		Model:                      a.Model,
		SubTypeSize:                a.SubTypeSize,
		SerialNumber:               a.SerialNumber,
		OperatingWeightKg:          fromNumericVal(a.OperatingWeightKg),
		MinCarrierWeightTonnes:     fromNumericVal(a.MinCarrierWeightTonnes),
		MaxCarrierWeightTonnes:     fromNumericVal(a.MaxCarrierWeightTonnes),
		RequiredFlowMinLmin:        fromNumericVal(a.RequiredFlowMinLmin),
		RequiredFlowMaxLmin:        fromNumericVal(a.RequiredFlowMaxLmin),
		RequiredPressureBar:        fromNumericVal(a.RequiredPressureBar),
		PinDiameterMm:              fromNumericVal(a.PinDiameterMm),
		PinSpreadMm:                fromNumericVal(a.PinSpreadMm),
		StandaloneDailyRateMinor:   a.StandaloneDailyRateMinor,
		SecurityDepositMinor:       a.SecurityDepositMinor,
		MinRentalDays:              int(a.MinRentalDays),
		PaymentEvidenceWindowHours: int(a.PaymentEvidenceWindowHours),
		Status:                     domain.AttachmentStatus(a.Status),
		CompatibleCategories:       categories,
		Photos:                     photos,
		CreatedAt:                  a.CreatedAt.Time,
		UpdatedAt:                  a.UpdatedAt.Time,
	}
}

func (r *PostgresMachineryRepository) CreateAttachment(
	ctx context.Context,
	attachment *domain.Attachment,
	compatibleCategories []string,
	photos []domain.AttachmentPhoto,
) (*domain.Attachment, error) {
	pgOwnerID := pgtype.UUID{Bytes: attachment.OwnerID, Valid: true}
	created, err := r.queries.CreateAttachment(ctx, db.CreateAttachmentParams{
		OwnerID:                    pgOwnerID,
		AttachmentType:             string(attachment.AttachmentType),
		Make:                       attachment.Make,
		Model:                      attachment.Model,
		SubTypeSize:                attachment.SubTypeSize,
		SerialNumber:               attachment.SerialNumber,
		OperatingWeightKg:          toNumeric(attachment.OperatingWeightKg),
		MinCarrierWeightTonnes:     toNumeric(attachment.MinCarrierWeightTonnes),
		MaxCarrierWeightTonnes:     toNumeric(attachment.MaxCarrierWeightTonnes),
		RequiredFlowMinLmin:        toNumeric(attachment.RequiredFlowMinLmin),
		RequiredFlowMaxLmin:        toNumeric(attachment.RequiredFlowMaxLmin),
		RequiredPressureBar:        toNumeric(attachment.RequiredPressureBar),
		PinDiameterMm:              toNumeric(attachment.PinDiameterMm),
		PinSpreadMm:                toNumeric(attachment.PinSpreadMm),
		StandaloneDailyRateMinor:   attachment.StandaloneDailyRateMinor,
		SecurityDepositMinor:       attachment.SecurityDepositMinor,
		MinRentalDays:              int32(attachment.MinRentalDays),
		PaymentEvidenceWindowHours: int32(attachment.PaymentEvidenceWindowHours),
		Status:                     string(attachment.Status),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create attachment: %w", err)
	}

	attID := created.ID
	for _, cat := range compatibleCategories {
		_, _ = r.queries.CreateAttachmentCompatibleCategory(ctx, db.CreateAttachmentCompatibleCategoryParams{
			AttachmentID: attID,
			CategoryCode: cat,
		})
	}

	persistedPhotos := make([]domain.AttachmentPhoto, 0, len(photos))
	for _, p := range photos {
		ph, err := r.queries.CreateAttachmentPhoto(ctx, db.CreateAttachmentPhotoParams{
			AttachmentID: attID,
			Angle:        string(p.Angle),
			StorageUri:   p.StorageURI,
		})
		if err == nil {
			persistedPhotos = append(persistedPhotos, domain.AttachmentPhoto{
				ID:           uuid.UUID(ph.ID.Bytes),
				AttachmentID: uuid.UUID(attID.Bytes),
				Angle:        p.Angle,
				StorageURI:   ph.StorageUri,
				CreatedAt:    ph.CreatedAt.Time,
			})
		}
	}

	return toDomainAttachment(created, compatibleCategories, persistedPhotos), nil
}

func (r *PostgresMachineryRepository) GetAttachmentByID(ctx context.Context, id uuid.UUID) (*domain.Attachment, error) {
	pgID := pgtype.UUID{Bytes: id, Valid: true}
	a, err := r.queries.GetAttachmentByID(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAttachmentNotFound
		}
		return nil, fmt.Errorf("failed to get attachment: %w", err)
	}

	cats, _ := r.queries.ListAttachmentCompatibleCategories(ctx, pgID)
	photoRows, _ := r.queries.ListAttachmentPhotos(ctx, pgID)
	photos := make([]domain.AttachmentPhoto, len(photoRows))
	for i, ph := range photoRows {
		photos[i] = domain.AttachmentPhoto{
			ID:           uuid.UUID(ph.ID.Bytes),
			AttachmentID: id,
			Angle:        domain.AttachmentPhotoAngle(ph.Angle),
			StorageURI:   ph.StorageUri,
			CreatedAt:    ph.CreatedAt.Time,
		}
	}

	return toDomainAttachment(a, cats, photos), nil
}

func (r *PostgresMachineryRepository) ListAttachmentsByOwner(ctx context.Context, ownerID uuid.UUID, status string) ([]*domain.Attachment, error) {
	pgOwnerID := pgtype.UUID{Bytes: ownerID, Valid: true}
	rows, err := r.queries.ListAttachmentsByOwner(ctx, db.ListAttachmentsByOwnerParams{
		OwnerID: pgOwnerID,
		Status:  status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list attachments by owner: %w", err)
	}

	result := make([]*domain.Attachment, len(rows))
	for i, row := range rows {
		cats, _ := r.queries.ListAttachmentCompatibleCategories(ctx, row.ID)
		photoRows, _ := r.queries.ListAttachmentPhotos(ctx, row.ID)
		photos := make([]domain.AttachmentPhoto, len(photoRows))
		for j, ph := range photoRows {
			photos[j] = domain.AttachmentPhoto{
				ID:           uuid.UUID(ph.ID.Bytes),
				AttachmentID: uuid.UUID(row.ID.Bytes),
				Angle:        domain.AttachmentPhotoAngle(ph.Angle),
				StorageURI:   ph.StorageUri,
				CreatedAt:    ph.CreatedAt.Time,
			}
		}
		result[i] = toDomainAttachment(row, cats, photos)
	}
	return result, nil
}

func (r *PostgresMachineryRepository) UpdateAttachmentStatus(ctx context.Context, id uuid.UUID, status domain.AttachmentStatus) error {
	pgID := pgtype.UUID{Bytes: id, Valid: true}
	_, err := r.queries.UpdateAttachmentStatus(ctx, db.UpdateAttachmentStatusParams{
		ID:     pgID,
		Status: string(status),
	})
	if err != nil {
		return fmt.Errorf("failed to update attachment status: %w", err)
	}
	return nil
}

func (r *PostgresMachineryRepository) LinkAttachmentToFleet(ctx context.Context, fleetID uuid.UUID, attachmentID uuid.UUID, bundledDailyRateMinor int64) (*domain.FleetBundledAttachment, error) {
	pgFleetID := pgtype.UUID{Bytes: fleetID, Valid: true}
	pgAttID := pgtype.UUID{Bytes: attachmentID, Valid: true}

	att, err := r.GetAttachmentByID(ctx, attachmentID)
	if err != nil {
		return nil, err
	}

	link, err := r.queries.UpsertFleetBundledAttachment(ctx, db.UpsertFleetBundledAttachmentParams{
		FleetID:               pgFleetID,
		AttachmentID:          pgAttID,
		BundledDailyRateMinor: bundledDailyRateMinor,
		IsActive:              true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to link fleet attachment: %w", err)
	}

	return &domain.FleetBundledAttachment{
		ID:                    uuid.UUID(link.ID.Bytes),
		FleetID:               fleetID,
		AttachmentID:          attachmentID,
		BundledDailyRateMinor: link.BundledDailyRateMinor,
		IsActive:              link.IsActive,
		AttachmentType:        att.AttachmentType,
		Make:                  att.Make,
		Model:                 att.Model,
		SubTypeSize:           att.SubTypeSize,
		SerialNumber:          att.SerialNumber,
		StandaloneDailyRate:   att.StandaloneDailyRateMinor,
		Status:                att.Status,
		CreatedAt:             link.CreatedAt.Time,
		UpdatedAt:             link.UpdatedAt.Time,
	}, nil
}

func (r *PostgresMachineryRepository) ListBundledAttachmentsByFleet(ctx context.Context, fleetID uuid.UUID) ([]*domain.FleetBundledAttachment, error) {
	pgFleetID := pgtype.UUID{Bytes: fleetID, Valid: true}
	rows, err := r.queries.ListBundledAttachmentsByFleet(ctx, pgFleetID)
	if err != nil {
		return nil, fmt.Errorf("failed to list bundled attachments: %w", err)
	}

	result := make([]*domain.FleetBundledAttachment, len(rows))
	for i, row := range rows {
		result[i] = &domain.FleetBundledAttachment{
			ID:                    uuid.UUID(row.LinkID.Bytes),
			FleetID:               fleetID,
			AttachmentID:          uuid.UUID(row.AttachmentID.Bytes),
			BundledDailyRateMinor: row.BundledDailyRateMinor,
			IsActive:              row.IsActive,
			AttachmentType:        domain.AttachmentType(row.AttachmentType),
			Make:                  row.Make,
			Model:                 row.Model,
			SubTypeSize:           row.SubTypeSize,
			SerialNumber:          row.SerialNumber,
			StandaloneDailyRate:   row.StandaloneDailyRateMinor,
			Status:                domain.AttachmentStatus(row.Status),
		}
	}
	return result, nil
}

func (r *PostgresMachineryRepository) CheckAttachmentAvailability(ctx context.Context, attachmentID uuid.UUID, startDate, endDate time.Time) (bool, error) {
	pgAttID := pgtype.UUID{Bytes: attachmentID, Valid: true}
	hireRange := pgtype.Range[pgtype.Date]{
		Lower:     pgtype.Date{Time: startDate, Valid: true},
		Upper:     pgtype.Date{Time: endDate, Valid: true},
		LowerType: pgtype.Inclusive,
		UpperType: pgtype.Exclusive,
		Valid:     true,
	}

	hasOverlap, err := r.queries.CheckAttachmentOverlap(ctx, db.CheckAttachmentOverlapParams{
		AttachmentID: pgAttID,
		Column2:      hireRange,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check overlap: %w", err)
	}

	return !hasOverlap, nil
}

func (r *PostgresMachineryRepository) CreateReservation(ctx context.Context, res *domain.AttachmentReservation) error {
	pgAttID := pgtype.UUID{Bytes: res.AttachmentID, Valid: true}
	var pgRefID pgtype.UUID
	if res.ReferenceID != nil {
		pgRefID = pgtype.UUID{Bytes: *res.ReferenceID, Valid: true}
	}
	hireRange := pgtype.Range[pgtype.Date]{
		Lower:     pgtype.Date{Time: res.StartDate, Valid: true},
		Upper:     pgtype.Date{Time: res.EndDate, Valid: true},
		LowerType: pgtype.Inclusive,
		UpperType: pgtype.Exclusive,
		Valid:     true,
	}

	_, err := r.queries.CreateAttachmentReservation(ctx, db.CreateAttachmentReservationParams{
		AttachmentID:          pgAttID,
		ReferenceID:           pgRefID,
		HirePeriod:            hireRange,
		DailyRateAppliedMinor: res.DailyRateAppliedMinor,
		Status:                res.Status,
	})
	if err != nil {
		return fmt.Errorf("failed to create attachment reservation: %w", err)
	}
	return nil
}
