package infrastructure

import (
	"context"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
)

type MemoryMachineryRepository struct {
	mu                     sync.RWMutex
	categories             map[uuid.UUID]domain.Category
	specGroups             map[uuid.UUID][]domain.SpecFieldGroup
	fleets                 map[uuid.UUID]domain.FleetAsset
	photos                 map[uuid.UUID][]domain.FleetPhoto
	specValues             map[uuid.UUID][]domain.FleetSpecValue
	specIndexed            map[uuid.UUID]domain.FleetSpecIndexed
	attachments            map[uuid.UUID]*domain.Attachment
	attachmentPhotos       map[uuid.UUID][]domain.AttachmentPhoto
	attachmentCategories   map[uuid.UUID][]string
	fleetBundled           map[uuid.UUID]map[uuid.UUID]*domain.FleetBundledAttachment
	attachmentReservations map[uuid.UUID][]*domain.AttachmentReservation
	fleetAvailability      map[uuid.UUID][]*domain.DateRangeState
	fleetBlackouts         map[uuid.UUID][]*domain.DateRangeReason
	ownerAggregates        map[uuid.UUID]*domain.OwnerAggregate
	fleetInsurances        map[uuid.UUID]*domain.FleetInsurance
}

var _ domain.CategoryRepository = (*MemoryMachineryRepository)(nil)
var _ domain.SpecDefinitionRepository = (*MemoryMachineryRepository)(nil)
var _ domain.FleetRepository = (*MemoryMachineryRepository)(nil)
var _ domain.FleetMediaRepository = (*MemoryMachineryRepository)(nil)
var _ domain.FleetSpecRepository = (*MemoryMachineryRepository)(nil)
var _ domain.AttachmentRepository = (*MemoryMachineryRepository)(nil)
var _ domain.DiscoveryRepository = (*MemoryMachineryRepository)(nil)

func NewMemoryMachineryRepository() *MemoryMachineryRepository {
	repo := &MemoryMachineryRepository{
		categories:             make(map[uuid.UUID]domain.Category),
		specGroups:             make(map[uuid.UUID][]domain.SpecFieldGroup),
		fleets:                 make(map[uuid.UUID]domain.FleetAsset),
		photos:                 make(map[uuid.UUID][]domain.FleetPhoto),
		specValues:             make(map[uuid.UUID][]domain.FleetSpecValue),
		specIndexed:            make(map[uuid.UUID]domain.FleetSpecIndexed),
		attachments:            make(map[uuid.UUID]*domain.Attachment),
		attachmentPhotos:       make(map[uuid.UUID][]domain.AttachmentPhoto),
		attachmentCategories:   make(map[uuid.UUID][]string),
		fleetBundled:           make(map[uuid.UUID]map[uuid.UUID]*domain.FleetBundledAttachment),
		attachmentReservations: make(map[uuid.UUID][]*domain.AttachmentReservation),
		fleetAvailability:      make(map[uuid.UUID][]*domain.DateRangeState),
		fleetBlackouts:         make(map[uuid.UUID][]*domain.DateRangeReason),
		ownerAggregates:        make(map[uuid.UUID]*domain.OwnerAggregate),
		fleetInsurances:        make(map[uuid.UUID]*domain.FleetInsurance),
	}
	repo.seedDefaults()
	return repo
}

func (m *MemoryMachineryRepository) seedDefaults() {
	rootID := uuid.MustParse("019a8f00-0001-7000-8000-000000000001")
	subID := uuid.MustParse("019a8f00-0002-7000-8000-000000000001")

	m.categories[rootID] = domain.Category{
		ID:           rootID,
		Code:         "EXCAVATION",
		NameEN:       "Excavation",
		NameAM:       "የቁፋሮ ማሽነሪዎች",
		IconURL:      "https://cdn.dinkmachinery.com/icons/excavation.svg",
		IsActive:     true,
		DisplayOrder: 1,
	}

	m.categories[subID] = domain.Category{
		ID:           subID,
		ParentID:     &rootID,
		Code:         "EXCAVATOR_MINI",
		NameEN:       "Mini Excavator (1-5 tons)",
		NameAM:       "አነስተኛ ኤክስካቫተር (1-5 ቶን)",
		IconURL:      "https://cdn.dinkmachinery.com/icons/sub/mini_excavator.svg",
		IsActive:     true,
		DisplayOrder: 1,
	}

	groupID := uuid.MustParse("019a8f00-0003-7000-8000-000000000001")
	m.specGroups[subID] = []domain.SpecFieldGroup{
		{
			ID:            groupID,
			CategoryID:    subID,
			SectionNameEN: "Weights & Dimensions",
			SectionNameAM: "የክብደትና ስፋት ዝርዝሮች",
			DisplayOrder:  1,
			Fields: []domain.SpecFieldDefinition{
				{
					ID:                  uuid.New(),
					GroupID:             groupID,
					FieldKey:            "operating_weight_kg",
					LabelEN:             "Operating Weight",
					LabelAM:             "የስራ ክብደት",
					FieldType:           "DECIMAL",
					IsRequired:          true,
					IsLogisticsCritical: true,
					IsComparable:        true,
				},
				{
					ID:                  uuid.New(),
					GroupID:             groupID,
					FieldKey:            "transport_length_mm",
					LabelEN:             "Transport Length",
					LabelAM:             "የማጓጓዣ ርዝመት",
					FieldType:           "INTEGER",
					IsRequired:          true,
					IsLogisticsCritical: true,
					IsComparable:        true,
				},
			},
		},
	}
}

// Haversine distance in meters
func haversineDistanceMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000.0 // meters
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0

	lat1Rad := lat1 * math.Pi / 180.0
	lat2Rad := lat2 * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// Category Repository
func (m *MemoryMachineryRepository) ListAll(ctx context.Context) ([]domain.Category, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []domain.Category
	for _, c := range m.categories {
		list = append(list, c)
	}
	return list, nil
}

func (m *MemoryMachineryRepository) ListRoots(ctx context.Context) ([]domain.Category, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []domain.Category
	for _, c := range m.categories {
		if c.ParentID == nil {
			list = append(list, c)
		}
	}
	return list, nil
}

func (m *MemoryMachineryRepository) ListSubcategories(ctx context.Context, parentID uuid.UUID) ([]domain.Category, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []domain.Category
	for _, c := range m.categories {
		if c.ParentID != nil && *c.ParentID == parentID {
			list = append(list, c)
		}
	}
	return list, nil
}

func (m *MemoryMachineryRepository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cat, ok := m.categories[id]
	if !ok {
		return nil, domain.ErrCategoryNotFound
	}
	return &cat, nil
}

// Spec Definition Repository
func (m *MemoryMachineryRepository) GetCategorySpecDefinitions(ctx context.Context, categoryID uuid.UUID) ([]domain.SpecFieldGroup, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	groups, ok := m.specGroups[categoryID]
	if !ok {
		return []domain.SpecFieldGroup{}, nil
	}
	return groups, nil
}

// Fleet Repository
func (m *MemoryMachineryRepository) CreateFleet(ctx context.Context, fleet *domain.FleetAsset) (*domain.FleetAsset, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if fleet.ID == uuid.Nil {
		fleet.ID = uuid.New()
	}
	fleet.CreatedAt = time.Now()
	fleet.UpdatedAt = time.Now()
	m.fleets[fleet.ID] = *fleet
	return fleet, nil
}

func (m *MemoryMachineryRepository) GetFleetByID(ctx context.Context, id uuid.UUID) (*domain.FleetAsset, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	fleet, ok := m.fleets[id]
	if !ok {
		return nil, domain.ErrFleetNotFound
	}

	fleet.Photos = m.photos[id]
	fleet.SpecValues = m.specValues[id]
	if spec, ok := m.specIndexed[id]; ok {
		fleet.LogisticsSpec = &spec
	}
	return &fleet, nil
}

func (m *MemoryMachineryRepository) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]domain.FleetAsset, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []domain.FleetAsset
	for _, f := range m.fleets {
		if f.OwnerID == ownerID {
			list = append(list, f)
		}
	}
	return list, nil
}

func (m *MemoryMachineryRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.FleetStatus) (*domain.FleetAsset, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	fleet, ok := m.fleets[id]
	if !ok {
		return nil, domain.ErrFleetNotFound
	}
	fleet.Status = status
	fleet.UpdatedAt = time.Now()
	m.fleets[id] = fleet
	return &fleet, nil
}

// Media Repository
func (m *MemoryMachineryRepository) SavePhotos(ctx context.Context, fleetID uuid.UUID, photos []domain.FleetPhoto) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.photos[fleetID] = append(m.photos[fleetID], photos...)
	return nil
}

func (m *MemoryMachineryRepository) ListPhotos(ctx context.Context, fleetID uuid.UUID) ([]domain.FleetPhoto, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.photos[fleetID], nil
}

// Spec Repository
func (m *MemoryMachineryRepository) SaveSpecValues(ctx context.Context, fleetID uuid.UUID, values []domain.FleetSpecValue) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.specValues[fleetID] = append(m.specValues[fleetID], values...)
	return nil
}

func (m *MemoryMachineryRepository) UpsertSpecIndexed(ctx context.Context, indexed *domain.FleetSpecIndexed) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	indexed.UpdatedAt = time.Now()
	m.specIndexed[indexed.FleetID] = *indexed
	return nil
}

func (m *MemoryMachineryRepository) GetSpecIndexed(ctx context.Context, fleetID uuid.UUID) (*domain.FleetSpecIndexed, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	spec, ok := m.specIndexed[fleetID]
	if !ok {
		return nil, nil
	}
	return &spec, nil
}

// Attachment Repository
func (m *MemoryMachineryRepository) CreateAttachment(
	ctx context.Context,
	attachment *domain.Attachment,
	compatibleCategories []string,
	photos []domain.AttachmentPhoto,
) (*domain.Attachment, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if attachment.ID == uuid.Nil {
		attachment.ID = uuid.New()
	}
	attachment.CreatedAt = time.Now()
	attachment.UpdatedAt = time.Now()

	cp := *attachment
	m.attachments[attachment.ID] = &cp
	m.attachmentCategories[attachment.ID] = compatibleCategories
	m.attachmentPhotos[attachment.ID] = photos

	return attachment, nil
}

func (m *MemoryMachineryRepository) GetAttachmentByID(ctx context.Context, id uuid.UUID) (*domain.Attachment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	att, ok := m.attachments[id]
	if !ok {
		return nil, domain.ErrAttachmentNotFound
	}
	cp := *att
	return &cp, nil
}

func (m *MemoryMachineryRepository) ListAttachmentsByOwner(ctx context.Context, ownerID uuid.UUID, status string) ([]*domain.Attachment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []*domain.Attachment
	for _, a := range m.attachments {
		if a.OwnerID == ownerID {
			if status == "" || string(a.Status) == status {
				cp := *a
				list = append(list, &cp)
			}
		}
	}
	return list, nil
}

func (m *MemoryMachineryRepository) UpdateAttachmentStatus(ctx context.Context, id uuid.UUID, status domain.AttachmentStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	att, ok := m.attachments[id]
	if !ok {
		return domain.ErrAttachmentNotFound
	}
	att.Status = status
	att.UpdatedAt = time.Now()
	return nil
}

func (m *MemoryMachineryRepository) LinkAttachmentToFleet(
	ctx context.Context,
	fleetID uuid.UUID,
	attachmentID uuid.UUID,
	bundledDailyRateMinor int64,
) (*domain.FleetBundledAttachment, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.fleetBundled[fleetID]; !ok {
		m.fleetBundled[fleetID] = make(map[uuid.UUID]*domain.FleetBundledAttachment)
	}

	link := &domain.FleetBundledAttachment{
		ID:                    uuid.New(),
		FleetID:               fleetID,
		AttachmentID:          attachmentID,
		BundledDailyRateMinor: bundledDailyRateMinor,
		IsActive:              true,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	m.fleetBundled[fleetID][attachmentID] = link
	return link, nil
}

func (m *MemoryMachineryRepository) ListBundledAttachmentsByFleet(ctx context.Context, fleetID uuid.UUID) ([]*domain.FleetBundledAttachment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []*domain.FleetBundledAttachment
	if links, ok := m.fleetBundled[fleetID]; ok {
		for _, l := range links {
			if l.IsActive {
				cp := *l
				list = append(list, &cp)
			}
		}
	}
	return list, nil
}

func (m *MemoryMachineryRepository) CheckAttachmentAvailability(ctx context.Context, attachmentID uuid.UUID, startDate, endDate time.Time) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	reservations, ok := m.attachmentReservations[attachmentID]
	if !ok {
		return true, nil
	}

	for _, r := range reservations {
		if r.Status == "CANCELLED" || r.Status == "COMPLETED" {
			continue
		}
		// Overlap check
		if startDate.Before(r.EndDate) && r.StartDate.Before(endDate) {
			return false, nil
		}
	}

	return true, nil
}

func (m *MemoryMachineryRepository) CreateReservation(ctx context.Context, res *domain.AttachmentReservation) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if reservations, ok := m.attachmentReservations[res.AttachmentID]; ok {
		for _, r := range reservations {
			if r.Status == "CANCELLED" || r.Status == "COMPLETED" {
				continue
			}
			if res.StartDate.Before(r.EndDate) && r.StartDate.Before(res.EndDate) {
				return domain.ErrAttachmentScheduleConflict
			}
		}
	}

	cp := *res
	m.attachmentReservations[res.AttachmentID] = append(m.attachmentReservations[res.AttachmentID], &cp)
	return nil
}

// Discovery Repository Implementation (Partition 05-04, SRS §05.04)

func (m *MemoryMachineryRepository) SearchFleetsSpatial(ctx context.Context, filter domain.FleetSearchFilter) ([]domain.FleetSearchResult, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	radiusMeters := filter.RadiusKm * 1000.0
	countryCode := filter.CountryCode
	if countryCode == "" {
		countryCode = "ET"
	}

	var candidates []domain.FleetSearchResult

	for _, f := range m.fleets {
		if f.Status != domain.FleetStatusActive {
			continue
		}

		// Distance check
		dist := haversineDistanceMeters(filter.Latitude, filter.Longitude, f.Latitude, f.Longitude)
		if dist > radiusMeters {
			continue
		}

		// Subcategory filter
		if filter.SubcategoryID != nil && f.SubcategoryID != *filter.SubcategoryID {
			continue
		}

		// Rate filters
		if filter.MinRateMinor != nil && f.DailyRateMinor < *filter.MinRateMinor {
			continue
		}
		if filter.MaxRateMinor != nil && f.DailyRateMinor > *filter.MaxRateMinor {
			continue
		}

		// Attribute filters
		if filter.FuelType != nil && string(f.FuelType) != *filter.FuelType {
			continue
		}
		if filter.OperatorIncluded != nil && f.OperatorIncluded != *filter.OperatorIncluded {
			continue
		}
		if filter.SelfMoveCapable != nil && f.SelfMoveCapable != *filter.SelfMoveCapable {
			continue
		}

		// Query keyword search
		if filter.Query != nil && *filter.Query != "" {
			q := strings.ToLower(strings.TrimSpace(*filter.Query))
			makeLower := strings.ToLower(f.Make)
			modelLower := strings.ToLower(f.Model)
			addrLower := strings.ToLower(f.AddressText)
			if !strings.Contains(makeLower, q) && !strings.Contains(modelLower, q) && !strings.Contains(addrLower, q) {
				continue
			}
		}

		// Temporal blackout / reservation exclusion
		if filter.StartDate != nil && filter.EndDate != nil {
			conflict := false
			if resList, ok := m.fleetAvailability[f.ID]; ok {
				for _, r := range resList {
					if filter.StartDate.Before(r.EndDate) && r.StartDate.Before(*filter.EndDate) {
						conflict = true
						break
					}
				}
			}
			if conflict {
				continue
			}

			if blackouts, ok := m.fleetBlackouts[f.ID]; ok {
				for _, b := range blackouts {
					if filter.StartDate.Before(b.EndDate) && b.StartDate.Before(*filter.EndDate) {
						conflict = true
						break
					}
				}
			}
			if conflict {
				continue
			}
		}

		// Category details
		subcatNameEN := ""
		subcatNameAM := ""
		if subcat, ok := m.categories[f.SubcategoryID]; ok {
			subcatNameEN = subcat.NameEN
			subcatNameAM = subcat.NameAM
		}

		// Primary photo
		var primaryPhoto *string
		if photos, ok := m.photos[f.ID]; ok && len(photos) > 0 {
			primaryPhoto = &photos[0].StorageURI
		}

		// Owner aggregate ratings
		ownerRating := 5.00
		ownerCompleted := 0
		ownerVerified := false
		if agg, ok := m.ownerAggregates[f.OwnerID]; ok {
			ownerRating = agg.AverageRating
			ownerCompleted = agg.TotalCompletedRentals
			ownerVerified = agg.IsVerified
		}

		// Key specs
		var weight *float64
		var payload *float64
		var power *float64
		if spec, ok := m.specIndexed[f.ID]; ok {
			weight = spec.OperatingWeightKg
			payload = spec.PayloadCapacityKg
			power = spec.GrossEnginePowerHp
		}

		score := domain.CalculateRecommendationScore(ownerRating, dist, radiusMeters, f.ManufactureYear, 1.0)

		candidates = append(candidates, domain.FleetSearchResult{
			FleetID:             f.ID,
			SubcategoryID:       f.SubcategoryID,
			SubcategoryNameEN:   subcatNameEN,
			SubcategoryNameAM:   subcatNameAM,
			Make:                f.Make,
			Model:               f.Model,
			ManufactureYear:     f.ManufactureYear,
			DailyRateMinor:      f.DailyRateMinor,
			ExcessHourRateMinor: f.ExcessHourRateMinor,
			DistanceMeters:      math.Round(dist*10) / 10,
			AvailabilityStatus:  domain.AvailabilityStatusAvailableNow,
			SelfMoveCapable:     f.SelfMoveCapable,
			OperatorIncluded:    f.OperatorIncluded,
			FuelType:            string(f.FuelType),
			PrimaryPhotoURI:     primaryPhoto,
			Owner: domain.FleetSearchOwner{
				OwnerID:               f.OwnerID,
				TradingName:           "Authorized Fleet Partner",
				AverageRating:         ownerRating,
				TotalCompletedRentals: ownerCompleted,
				IsVerified:            ownerVerified,
			},
			KeySpecs: domain.FleetSearchKeySpecs{
				OperatingWeightKg: weight,
				PayloadCapacityKg: payload,
				GrossPowerHp:      power,
			},
			RecommendationScore: score,
		})
	}

	domain.SortSearchResults(candidates, filter.SortBy)

	totalMatches := len(candidates)
	start := filter.Offset
	if start > totalMatches {
		return []domain.FleetSearchResult{}, totalMatches, nil
	}
	end := start + filter.Limit
	if end > totalMatches {
		end = totalMatches
	}

	return candidates[start:end], totalMatches, nil
}

func (m *MemoryMachineryRepository) GetFleetsForComparison(ctx context.Context, fleetIDs []uuid.UUID) ([]domain.FleetAsset, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []domain.FleetAsset
	for _, id := range fleetIDs {
		f, ok := m.fleets[id]
		if !ok {
			return nil, domain.ErrFleetNotFound
		}
		f.Photos = m.photos[id]
		result = append(result, f)
	}
	return result, nil
}

func (m *MemoryMachineryRepository) GetComparableSpecsForFleets(ctx context.Context, fleetIDs []uuid.UUID) ([]domain.SpecFieldDefinition, []domain.FleetSpecValue, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(fleetIDs) == 0 {
		return nil, nil, nil
	}

	firstFleet, ok := m.fleets[fleetIDs[0]]
	if !ok {
		return nil, nil, domain.ErrFleetNotFound
	}

	var defs []domain.SpecFieldDefinition
	if groups, ok := m.specGroups[firstFleet.SubcategoryID]; ok {
		for _, g := range groups {
			for _, f := range g.Fields {
				if f.IsComparable {
					defs = append(defs, f)
				}
			}
		}
	}

	var values []domain.FleetSpecValue
	for _, id := range fleetIDs {
		if vals, ok := m.specValues[id]; ok {
			values = append(values, vals...)
		}
	}

	return defs, values, nil
}

func (m *MemoryMachineryRepository) GetFleetAvailability(ctx context.Context, fleetID uuid.UUID, start, end time.Time) (*domain.FleetAvailabilityInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.fleets[fleetID]; !ok {
		return nil, domain.ErrFleetNotFound
	}

	var blackouts []domain.DateRangeReason
	if bList, ok := m.fleetBlackouts[fleetID]; ok {
		for _, b := range bList {
			if start.Before(b.EndDate) && b.StartDate.Before(end) {
				blackouts = append(blackouts, *b)
			}
		}
	}

	var reserved []domain.DateRangeState
	if rList, ok := m.fleetAvailability[fleetID]; ok {
		for _, r := range rList {
			if start.Before(r.EndDate) && r.StartDate.Before(end) {
				reserved = append(reserved, *r)
			}
		}
	}

	isAvailable := len(blackouts) == 0 && len(reserved) == 0
	var reason *string
	if !isAvailable {
		r := "RESERVED"
		if len(blackouts) > 0 {
			r = "MAINTENANCE_BLACKOUT"
		}
		reason = &r
	}

	return &domain.FleetAvailabilityInfo{
		FleetID:           fleetID,
		IsAvailable:       isAvailable,
		UnavailableReason: reason,
		BlackoutRanges:    blackouts,
		ReservedRanges:    reserved,
	}, nil
}

func (m *MemoryMachineryRepository) CreateFleetAvailability(ctx context.Context, fleetID uuid.UUID, start, end time.Time, state string, refID *uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	item := &domain.DateRangeState{
		StartDate: start,
		EndDate:   end,
		State:     state,
	}
	m.fleetAvailability[fleetID] = append(m.fleetAvailability[fleetID], item)
	return nil
}

func (m *MemoryMachineryRepository) UpsertOwnerAggregate(ctx context.Context, agg *domain.OwnerAggregate) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	cp := *agg
	cp.UpdatedAt = time.Now()
	m.ownerAggregates[agg.OwnerID] = &cp
	return nil
}

func (m *MemoryMachineryRepository) GetOwnerAggregate(ctx context.Context, ownerID uuid.UUID) (*domain.OwnerAggregate, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	agg, ok := m.ownerAggregates[ownerID]
	if !ok {
		return nil, nil
	}
	cp := *agg
	return &cp, nil
}

func (m *MemoryMachineryRepository) UpsertFleetInsurance(ctx context.Context, ins *domain.FleetInsurance) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	cp := *ins
	cp.UpdatedAt = time.Now()
	m.fleetInsurances[ins.FleetID] = &cp
	return nil
}

func (m *MemoryMachineryRepository) GetFleetInsurance(ctx context.Context, fleetID uuid.UUID) (*domain.FleetInsurance, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ins, ok := m.fleetInsurances[fleetID]
	if !ok {
		return nil, nil
	}
	cp := *ins
	return &cp, nil
}

// SeedCategory inserts or updates a category for testing or bootstrap.
func (m *MemoryMachineryRepository) SeedCategory(c domain.Category) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.categories[c.ID] = c
	return nil
}

// SeedSpecDefinitions inserts or updates spec field groups for a category.
func (m *MemoryMachineryRepository) SeedSpecDefinitions(categoryID uuid.UUID, groups []domain.SpecFieldGroup) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.specGroups[categoryID] = groups
	return nil
}

// SeedBlackout records a scheduled blackout window for testing.
func (m *MemoryMachineryRepository) SeedBlackout(fleetID uuid.UUID, start, end time.Time, reason string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.fleetBlackouts[fleetID] = append(m.fleetBlackouts[fleetID], &domain.DateRangeReason{
		StartDate: start,
		EndDate:   end,
		Reason:    reason,
	})
}
