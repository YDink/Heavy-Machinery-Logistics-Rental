package infrastructure

import (
	"context"
	"math"
	"sort"
	"sync"
	"time"

	bookingDomain "github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	"github.com/dinkmachinery/dink-core/internal/modules/logistics/domain"
	"github.com/google/uuid"
)

// BookingCreatorReader abstracts memory booking persistence for sub-booking creation in unit/HTTP tests.
type BookingCreatorReader interface {
	CreateBooking(ctx context.Context, b *bookingDomain.Booking) error
}

type MemoryLogisticsRepository struct {
	mu              sync.RWMutex
	bookingRepo     BookingCreatorReader
	profiles        map[uuid.UUID]*domain.FleetLogisticsProfile
	countryLimits   map[string]domain.CountryLogisticsLimits
	dispatchesBySub map[uuid.UUID]*domain.LogisticsDispatch
	dispatchesByPar map[uuid.UUID]*domain.LogisticsDispatch
	externalByPar   map[uuid.UUID]*domain.ExternalTransport
	transByID       map[uuid.UUID]*domain.LogisticsTransshipment
	transBySub      map[uuid.UUID]*domain.LogisticsTransshipment
	cascadesByPar   map[uuid.UUID]*domain.LogisticsReverseCascade
	releasedRefs    map[uuid.UUID]bool
	ratingDeduction map[uuid.UUID]float64
}

func NewMemoryLogisticsRepository(bookingRepo BookingCreatorReader) *MemoryLogisticsRepository {
	return &MemoryLogisticsRepository{
		bookingRepo:     bookingRepo,
		profiles:        make(map[uuid.UUID]*domain.FleetLogisticsProfile),
		countryLimits:   map[string]domain.CountryLogisticsLimits{
			"ET": {CountryCode: "ET", DefaultCurrency: "ETB", MaxGvmKg: 56000.0, MaxCorridorHeightMm: 4500, WideLoadWidthMm: 3000},
			"KE": {CountryCode: "KE", DefaultCurrency: "KES", MaxGvmKg: 56000.0, MaxCorridorHeightMm: 4500, WideLoadWidthMm: 3000},
		},
		dispatchesBySub: make(map[uuid.UUID]*domain.LogisticsDispatch),
		dispatchesByPar: make(map[uuid.UUID]*domain.LogisticsDispatch),
		externalByPar:   make(map[uuid.UUID]*domain.ExternalTransport),
		transByID:       make(map[uuid.UUID]*domain.LogisticsTransshipment),
		transBySub:      make(map[uuid.UUID]*domain.LogisticsTransshipment),
		cascadesByPar:   make(map[uuid.UUID]*domain.LogisticsReverseCascade),
		releasedRefs:    make(map[uuid.UUID]bool),
		ratingDeduction: make(map[uuid.UUID]float64),
	}
}

func (r *MemoryLogisticsRepository) SeedFleetProfile(p *domain.FleetLogisticsProfile) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cpy := *p
	r.profiles[p.FleetID] = &cpy
}

func (r *MemoryLogisticsRepository) GetFleetLogisticsProfile(ctx context.Context, fleetID uuid.UUID) (*domain.FleetLogisticsProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.profiles[fleetID]
	if !ok {
		return nil, domain.ErrLogisticsBookingNotFound
	}
	cpy := *p
	return &cpy, nil
}

func (r *MemoryLogisticsRepository) GetCountryLogisticsLimits(ctx context.Context, countryCode string) (domain.CountryLogisticsLimits, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if l, ok := r.countryLimits[countryCode]; ok {
		return l, nil
	}
	return domain.CountryLogisticsLimits{
		CountryCode:         countryCode,
		DefaultCurrency:     "ETB",
		MaxGvmKg:            domain.DefaultMaxGvmKg,
		MaxCorridorHeightMm: domain.DefaultMaxCorridorHeightMm,
		WideLoadWidthMm:     domain.DefaultWideLoadWidthMm,
	}, nil
}

func haversineMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const earthR = 6371000.0
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0)*math.Sin(dLon/2)*math.Sin(dLon/2)
	return earthR * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func (r *MemoryLogisticsRepository) MatchCompatibleLogisticsFleets(
	ctx context.Context,
	origin domain.GeoPoint,
	radiusMeters float64,
	countryCode string,
	excludeOwnerID uuid.UUID,
	specs domain.MachineLogisticsSpecs,
	limits domain.CountryLogisticsLimits,
	limit int,
) ([]domain.MatchedCarrier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matches []domain.MatchedCarrier
	for _, p := range r.profiles {
		if p.RootCategoryCode != "LOGISTICS" || p.Status != "ACTIVE" {
			continue
		}
		if p.CountryCode != countryCode || p.OwnerID == excludeOwnerID {
			continue
		}
		dist := haversineMeters(origin.Latitude, origin.Longitude, p.Location.Latitude, p.Location.Longitude)
		if dist > radiusMeters {
			continue
		}
		candidate := p.TrailerProfile
		candidate.DistanceFromDepotMeters = dist
		if candidate.Currency == "" {
			candidate.Currency = limits.DefaultCurrency
		}
		if _, err := domain.EvaluateDimensionalCompatibility(specs, candidate, limits); err != nil {
			continue
		}
		matches = append(matches, candidate)
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].DistanceFromDepotMeters < matches[j].DistanceFromDepotMeters
	})
	if limit > 0 && len(matches) > limit {
		matches = matches[:limit]
	}
	return matches, nil
}

func (r *MemoryLogisticsRepository) CreateLogisticsSubBookingAndDispatch(
	ctx context.Context,
	subBookingRef string,
	parentBookingID uuid.UUID,
	renterOwnerID uuid.UUID,
	carrierOwnerID uuid.UUID,
	logisticsFleetID uuid.UUID,
	startDate time.Time,
	endDate time.Time,
	durationDays int,
	dailyRateMinor int64,
	freightHireFeeMinor int64,
	commissionBps int,
	commissionOwedMinor int64,
	abnormalLoad bool,
	currency string,
	siteAddress string,
	dispatch *domain.LogisticsDispatch,
) (uuid.UUID, error) {
	now := time.Now().UTC()
	subID := uuid.New()
	subBooking := &bookingDomain.Booking{
		ID:                         subID,
		BookingReference:           subBookingRef,
		RenterID:                   renterOwnerID,
		OwnerID:                    carrierOwnerID,
		FleetID:                    logisticsFleetID,
		ParentBookingID:            &parentBookingID,
		BookingKind:                "LOGISTICS_SUB",
		StartDate:                  startDate,
		EndDate:                    endDate,
		DurationDays:               durationDays,
		DailyRateAppliedMinor:      dailyRateMinor,
		BaseHireFeeMinor:           freightHireFeeMinor,
		GrossRentalAmountMinor:     freightHireFeeMinor,
		PrepaymentAmountMinor:      freightHireFeeMinor,
		PlatformCommissionRateBps:  commissionBps,
		PlatformCommissionMinor:    commissionOwedMinor,
		AbnormalLoadEscortRequired: abnormalLoad,
		Currency:                   currency,
		SiteLocationAddress:        siteAddress,
		Status:                     bookingDomain.StatusConfirmed,
		ConfirmedAt:                &now,
		Version:                    1,
		CreatedAt:                  now,
		UpdatedAt:                  now,
	}
	if r.bookingRepo != nil {
		if err := r.bookingRepo.CreateBooking(ctx, subBooking); err != nil {
			return uuid.Nil, err
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	dispatch.ID = uuid.New()
	dispatch.SubBookingID = subID
	dispatch.Version = 1
	dispatch.CreatedAt = now
	dispatch.UpdatedAt = now
	cpy := *dispatch
	r.dispatchesBySub[subID] = &cpy
	r.dispatchesByPar[parentBookingID] = &cpy
	return subID, nil
}

func (r *MemoryLogisticsRepository) GetLogisticsDispatchBySubBookingID(ctx context.Context, subBookingID uuid.UUID) (*domain.LogisticsDispatch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.dispatchesBySub[subBookingID]
	if !ok {
		return nil, domain.ErrLogisticsBookingNotFound
	}
	cpy := *d
	return &cpy, nil
}

func (r *MemoryLogisticsRepository) GetLogisticsDispatchByParentBookingID(ctx context.Context, parentBookingID uuid.UUID) (*domain.LogisticsDispatch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.dispatchesByPar[parentBookingID]
	if !ok {
		return nil, domain.ErrLogisticsBookingNotFound
	}
	cpy := *d
	return &cpy, nil
}

func (r *MemoryLogisticsRepository) UpdateLogisticsDispatchLoading(ctx context.Context, subBookingID uuid.UUID, photos []domain.WinchingPhoto) (*domain.LogisticsDispatch, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	d, ok := r.dispatchesBySub[subBookingID]
	if !ok {
		return nil, domain.ErrLogisticsBookingNotFound
	}
	now := time.Now().UTC()
	d.WinchingPhotos = append([]domain.WinchingPhoto(nil), photos...)
	d.LoadedAt = &now
	d.Version++
	d.UpdatedAt = now
	cpy := *d
	return &cpy, nil
}

func (r *MemoryLogisticsRepository) UpdateLogisticsDispatchMilestone(ctx context.Context, subBookingID uuid.UUID, stage string) (*domain.LogisticsDispatch, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	d, ok := r.dispatchesBySub[subBookingID]
	if !ok {
		return nil, domain.ErrLogisticsBookingNotFound
	}
	now := time.Now().UTC()
	switch stage {
	case "DISPATCHED":
		d.DispatchedAt = &now
	case "LOADING_INSPECTION":
		d.LoadedAt = &now
	case "IN_TRANSIT":
		d.InTransitAt = &now
	case "READY_FOR_PICKUP":
		d.ArrivedAt = &now
	}
	d.Version++
	d.UpdatedAt = now
	cpy := *d
	return &cpy, nil
}

func (r *MemoryLogisticsRepository) UpsertExternalTransport(ctx context.Context, ext *domain.ExternalTransport) (*domain.ExternalTransport, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	if existing, ok := r.externalByPar[ext.ParentBookingID]; ok {
		ext.ID = existing.ID
		ext.CreatedAt = existing.CreatedAt
	} else {
		ext.ID = uuid.New()
		ext.CreatedAt = now
	}
	ext.UpdatedAt = now
	cpy := *ext
	r.externalByPar[ext.ParentBookingID] = &cpy
	return &cpy, nil
}

func (r *MemoryLogisticsRepository) GetExternalTransportByParentBookingID(ctx context.Context, parentBookingID uuid.UUID) (*domain.ExternalTransport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ext, ok := r.externalByPar[parentBookingID]
	if !ok {
		return nil, domain.ErrLogisticsBookingNotFound
	}
	cpy := *ext
	return &cpy, nil
}

func (r *MemoryLogisticsRepository) CreateLogisticsTransshipment(ctx context.Context, t *domain.LogisticsTransshipment) (*domain.LogisticsTransshipment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	t.ID = uuid.New()
	t.CreatedAt = now
	t.UpdatedAt = now
	cpy := *t
	r.transByID[t.ID] = &cpy
	r.transBySub[t.SubBookingID] = &cpy
	return &cpy, nil
}

func (r *MemoryLogisticsRepository) GetOpenTransshipmentBySubBookingID(ctx context.Context, subBookingID uuid.UUID) (*domain.LogisticsTransshipment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.transBySub[subBookingID]
	if !ok || (t.Status != domain.TransshipmentOpen && t.Status != domain.TransshipmentAccepted) {
		return nil, domain.ErrTransshipmentNotOpen
	}
	cpy := *t
	return &cpy, nil
}

func (r *MemoryLogisticsRepository) AcceptLogisticsTransshipment(ctx context.Context, transshipmentID, replacementFleetID, replacementCarrierID uuid.UUID) (*domain.LogisticsTransshipment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.transByID[transshipmentID]
	if !ok {
		return nil, domain.ErrTransshipmentNotOpen
	}
	now := time.Now().UTC()
	t.ReplacementFleetID = &replacementFleetID
	t.ReplacementCarrierID = &replacementCarrierID
	t.Status = domain.TransshipmentAccepted
	t.AcceptedAt = &now
	t.UpdatedAt = now
	cpy := *t
	return &cpy, nil
}

func (r *MemoryLogisticsRepository) ExpireLogisticsTransshipment(ctx context.Context, transshipmentID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.transByID[transshipmentID]
	if !ok {
		return domain.ErrTransshipmentNotOpen
	}
	t.Status = domain.TransshipmentExpired
	t.UpdatedAt = time.Now().UTC()
	return nil
}

func (r *MemoryLogisticsRepository) CreateLogisticsReverseCascade(ctx context.Context, rc *domain.LogisticsReverseCascade) (*domain.LogisticsReverseCascade, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rc.ID = uuid.New()
	rc.CreatedAt = time.Now().UTC()
	cpy := *rc
	r.cascadesByPar[rc.ParentBookingID] = &cpy
	return &cpy, nil
}

func (r *MemoryLogisticsRepository) GetLogisticsReverseCascadeByParentBookingID(ctx context.Context, parentBookingID uuid.UUID) (*domain.LogisticsReverseCascade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rc, ok := r.cascadesByPar[parentBookingID]
	if !ok {
		return nil, domain.ErrLogisticsBookingNotFound
	}
	cpy := *rc
	return &cpy, nil
}

func (r *MemoryLogisticsRepository) ReleaseAvailabilityAndDeductCarrierRating(ctx context.Context, parentBookingID, subBookingID, carrierOwnerID uuid.UUID, ratingDeduction float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.releasedRefs[parentBookingID] = true
	r.releasedRefs[subBookingID] = true
	r.ratingDeduction[carrierOwnerID] += ratingDeduction
	return nil
}

func (r *MemoryLogisticsRepository) IsAvailabilityReleased(refID uuid.UUID) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.releasedRefs[refID]
}

func (r *MemoryLogisticsRepository) CarrierRatingDeductionTotal(carrierOwnerID uuid.UUID) float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.ratingDeduction[carrierOwnerID]
}
