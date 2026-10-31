package infrastructure

import (
	"context"
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
}

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

// CategoryRepository implementation
func (m *MemoryMachineryRepository) ListAll(ctx context.Context) ([]domain.Category, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []domain.Category
	for _, c := range m.categories {
		if c.IsActive {
			list = append(list, c)
		}
	}
	return list, nil
}

func (m *MemoryMachineryRepository) ListRoots(ctx context.Context) ([]domain.Category, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []domain.Category
	for _, c := range m.categories {
		if c.ParentID == nil && c.IsActive {
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
		if c.ParentID != nil && *c.ParentID == parentID && c.IsActive {
			list = append(list, c)
		}
	}
	return list, nil
}

func (m *MemoryMachineryRepository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if c, exists := m.categories[id]; exists {
		return &c, nil
	}
	return nil, domain.ErrCategoryNotFound
}

// SpecDefinitionRepository implementation
func (m *MemoryMachineryRepository) GetCategorySpecDefinitions(ctx context.Context, categoryID uuid.UUID) ([]domain.SpecFieldGroup, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.specGroups[categoryID], nil
}

// FleetRepository implementation
func (m *MemoryMachineryRepository) CreateFleet(ctx context.Context, fleet *domain.FleetAsset) (*domain.FleetAsset, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	fleet.ID = uuid.New()
	fleet.Version = 1
	fleet.CreatedAt = time.Now()
	fleet.UpdatedAt = time.Now()

	m.fleets[fleet.ID] = *fleet
	return fleet, nil
}

func (m *MemoryMachineryRepository) GetFleetByID(ctx context.Context, id uuid.UUID) (*domain.FleetAsset, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if f, exists := m.fleets[id]; exists {
		return &f, nil
	}
	return nil, domain.ErrFleetNotFound
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
	f, exists := m.fleets[id]
	if !exists {
		return nil, domain.ErrFleetNotFound
	}
	f.Status = status
	f.Version++
	f.UpdatedAt = time.Now()
	m.fleets[id] = f
	return &f, nil
}

// FleetMediaRepository implementation
func (m *MemoryMachineryRepository) SavePhotos(ctx context.Context, fleetID uuid.UUID, photos []domain.FleetPhoto) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.photos[fleetID] = photos
	return nil
}

func (m *MemoryMachineryRepository) ListPhotos(ctx context.Context, fleetID uuid.UUID) ([]domain.FleetPhoto, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.photos[fleetID], nil
}

// FleetSpecRepository implementation
func (m *MemoryMachineryRepository) SaveSpecValues(ctx context.Context, fleetID uuid.UUID, values []domain.FleetSpecValue) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.specValues[fleetID] = values
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
	if sp, exists := m.specIndexed[fleetID]; exists {
		return &sp, nil
	}
	return nil, nil
}

// AttachmentRepository implementation
func (m *MemoryMachineryRepository) CreateAttachment(ctx context.Context, attachment *domain.Attachment, compatibleCategories []string, photos []domain.AttachmentPhoto) (*domain.Attachment, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	cp := *attachment
	cp.CompatibleCategories = compatibleCategories
	cp.Photos = photos
	m.attachments[cp.ID] = &cp
	m.attachmentPhotos[cp.ID] = photos
	m.attachmentCategories[cp.ID] = compatibleCategories
	return &cp, nil
}

func (m *MemoryMachineryRepository) GetAttachmentByID(ctx context.Context, id uuid.UUID) (*domain.Attachment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	att, exists := m.attachments[id]
	if !exists {
		return nil, domain.ErrAttachmentNotFound
	}
	cp := *att
	cp.Photos = m.attachmentPhotos[id]
	cp.CompatibleCategories = m.attachmentCategories[id]
	return &cp, nil
}

func (m *MemoryMachineryRepository) ListAttachmentsByOwner(ctx context.Context, ownerID uuid.UUID, status string) ([]*domain.Attachment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*domain.Attachment
	for _, att := range m.attachments {
		if att.OwnerID == ownerID {
			if status == "" || string(att.Status) == status {
				cp := *att
				cp.Photos = m.attachmentPhotos[att.ID]
				cp.CompatibleCategories = m.attachmentCategories[att.ID]
				result = append(result, &cp)
			}
		}
	}
	return result, nil
}

func (m *MemoryMachineryRepository) UpdateAttachmentStatus(ctx context.Context, id uuid.UUID, status domain.AttachmentStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	att, exists := m.attachments[id]
	if !exists {
		return domain.ErrAttachmentNotFound
	}
	att.Status = status
	att.UpdatedAt = time.Now()
	return nil
}

func (m *MemoryMachineryRepository) LinkAttachmentToFleet(ctx context.Context, fleetID uuid.UUID, attachmentID uuid.UUID, bundledDailyRateMinor int64) (*domain.FleetBundledAttachment, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	att, exists := m.attachments[attachmentID]
	if !exists {
		return nil, domain.ErrAttachmentNotFound
	}

	if _, exists := m.fleetBundled[fleetID]; !exists {
		m.fleetBundled[fleetID] = make(map[uuid.UUID]*domain.FleetBundledAttachment)
	}

	link := &domain.FleetBundledAttachment{
		ID:                    uuid.New(),
		FleetID:               fleetID,
		AttachmentID:          attachmentID,
		BundledDailyRateMinor: bundledDailyRateMinor,
		IsActive:              true,
		AttachmentType:        att.AttachmentType,
		Make:                  att.Make,
		Model:                 att.Model,
		SubTypeSize:           att.SubTypeSize,
		SerialNumber:          att.SerialNumber,
		StandaloneDailyRate:   att.StandaloneDailyRateMinor,
		Status:                att.Status,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	m.fleetBundled[fleetID][attachmentID] = link
	return link, nil
}

func (m *MemoryMachineryRepository) ListBundledAttachmentsByFleet(ctx context.Context, fleetID uuid.UUID) ([]*domain.FleetBundledAttachment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*domain.FleetBundledAttachment
	if links, exists := m.fleetBundled[fleetID]; exists {
		for _, l := range links {
			if l.IsActive {
				cp := *l
				result = append(result, &cp)
			}
		}
	}
	return result, nil
}

func (m *MemoryMachineryRepository) CheckAttachmentAvailability(ctx context.Context, attachmentID uuid.UUID, startDate, endDate time.Time) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	att, exists := m.attachments[attachmentID]
	if !exists {
		return false, domain.ErrAttachmentNotFound
	}

	if att.Status != domain.AttachmentStatusActive && att.Status != domain.AttachmentStatusRented {
		return false, nil
	}

	// Check reservations overlap: start < res.End && res.Start < end
	if reservations, ok := m.attachmentReservations[attachmentID]; ok {
		for _, r := range reservations {
			if r.Status == "CANCELLED" || r.Status == "COMPLETED" {
				continue
			}
			if startDate.Before(r.EndDate) && r.StartDate.Before(endDate) {
				return false, nil
			}
		}
	}

	return true, nil
}

func (m *MemoryMachineryRepository) CreateReservation(ctx context.Context, res *domain.AttachmentReservation) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Concurrency safety check: check overlap
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
