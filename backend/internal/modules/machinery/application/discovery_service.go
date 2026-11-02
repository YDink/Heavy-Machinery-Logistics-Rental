package application

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/dinkmachinery/dink-core/internal/modules/machinery/domain"
)

type HoldLeaseChecker interface {
	GetHoldLease(ctx context.Context, fleetID uuid.UUID) (bool, int)
}

type MemoryHoldLeaseChecker struct {
	mu    sync.RWMutex
	holds map[uuid.UUID]time.Time
}

func NewMemoryHoldLeaseChecker() *MemoryHoldLeaseChecker {
	return &MemoryHoldLeaseChecker{
		holds: make(map[uuid.UUID]time.Time),
	}
}

func (m *MemoryHoldLeaseChecker) SetHold(fleetID uuid.UUID, ttl time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.holds[fleetID] = time.Now().Add(ttl)
}

func (m *MemoryHoldLeaseChecker) ReleaseHold(fleetID uuid.UUID) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.holds, fleetID)
}

func (m *MemoryHoldLeaseChecker) GetHoldLease(ctx context.Context, fleetID uuid.UUID) (bool, int) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	exp, ok := m.holds[fleetID]
	if !ok {
		return false, 0
	}
	remaining := int(time.Until(exp).Seconds())
	if remaining <= 0 {
		return false, 0
	}
	return true, remaining
}

type FleetSearchResultSet struct {
	TotalMatches int
	Latitude     float64
	Longitude    float64
	RadiusKm     float64
	Results      []domain.FleetSearchResult
}

type DiscoveryService struct {
	discoveryRepo domain.DiscoveryRepository
	categoryRepo  domain.CategoryRepository
	leaseChecker  HoldLeaseChecker
}

func NewDiscoveryService(
	discoveryRepo domain.DiscoveryRepository,
	categoryRepo domain.CategoryRepository,
	leaseChecker HoldLeaseChecker,
) *DiscoveryService {
	if leaseChecker == nil {
		leaseChecker = NewMemoryHoldLeaseChecker()
	}
	return &DiscoveryService{
		discoveryRepo: discoveryRepo,
		categoryRepo:  categoryRepo,
		leaseChecker:  leaseChecker,
	}
}

// SearchFleets executes PostGIS spatial discovery, multi-filtering, and in-flight lease interception (SR-01 to SR-05, SR-07).
func (s *DiscoveryService) SearchFleets(ctx context.Context, filter domain.FleetSearchFilter) (*FleetSearchResultSet, error) {
	if err := domain.ValidateSearchFilter(&filter); err != nil {
		return nil, err
	}

	results, totalCount, err := s.discoveryRepo.SearchFleetsSpatial(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// Interrogate in-flight hold leases in Valkey / memory (SR-07, BK-01)
	for i := range results {
		isHeld, _ := s.leaseChecker.GetHoldLease(ctx, results[i].FleetID)
		if isHeld {
			results[i].AvailabilityStatus = domain.AvailabilityStatusHeldPendingWire
		}
	}

	return &FleetSearchResultSet{
		TotalMatches: totalCount,
		Latitude:     filter.Latitude,
		Longitude:    filter.Longitude,
		RadiusKm:     filter.RadiusKm,
		Results:      results,
	}, nil
}

// CompareFleets performs taxonomic homogeneity validation and dynamic EAV pivoting for 2-3 fleets (SR-08, FL-15).
func (s *DiscoveryService) CompareFleets(ctx context.Context, fleetIDs []uuid.UUID) (*domain.FleetCompareResult, error) {
	fleets, err := s.discoveryRepo.GetFleetsForComparison(ctx, fleetIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve fleets for comparison: %w", err)
	}

	if err := domain.ValidateComparisonSet(fleets); err != nil {
		return nil, err
	}

	subcat, err := s.categoryRepo.GetCategoryByID(ctx, fleets[0].SubcategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve subcategory: %w", err)
	}

	defs, values, err := s.discoveryRepo.GetComparableSpecsForFleets(ctx, fleetIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve comparable specs: %w", err)
	}

	matrix := domain.BuildComparisonMatrix(fleets, *subcat, defs, values)
	return matrix, nil
}

// GetFleetAvailability inspects booked ranges, maintenance blackouts, and active hold leases (SR-07).
func (s *DiscoveryService) GetFleetAvailability(ctx context.Context, fleetID uuid.UUID, start, end time.Time) (*domain.FleetAvailabilityInfo, error) {
	if end.Before(start) {
		return nil, domain.ErrInvalidDateRange
	}

	info, err := s.discoveryRepo.GetFleetAvailability(ctx, fleetID, start, end)
	if err != nil {
		return nil, err
	}

	// Interrogate active in-flight lease
	isHeld, ttl := s.leaseChecker.GetHoldLease(ctx, fleetID)
	if isHeld {
		info.IsAvailable = false
		reason := "HELD_BY_ANOTHER_RENTER"
		info.UnavailableReason = &reason
		info.HoldExpiresInSeconds = &ttl
	}

	return info, nil
}
