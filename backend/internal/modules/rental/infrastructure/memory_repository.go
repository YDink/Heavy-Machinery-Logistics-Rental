package infrastructure

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dinkmachinery/dink-core/internal/modules/rental/domain"
	"github.com/google/uuid"
)

type MemoryRentalRepository struct {
	mu            sync.RWMutex
	inspections   map[uuid.UUID]*domain.HandoverInspection
	challenges    map[uuid.UUID]*domain.HandoverChallenge
	shifts        map[uuid.UUID]*domain.RentalShift
	extensions    map[uuid.UUID]*domain.RentalExtension
	breakdowns    map[uuid.UUID]*domain.MechanicalBreakdown
	fleetBaselines map[uuid.UUID]float64
	conflicts     map[uuid.UUID]bool
}

func NewMemoryRentalRepository() *MemoryRentalRepository {
	return &MemoryRentalRepository{
		inspections:    make(map[uuid.UUID]*domain.HandoverInspection),
		challenges:     make(map[uuid.UUID]*domain.HandoverChallenge),
		shifts:         make(map[uuid.UUID]*domain.RentalShift),
		extensions:     make(map[uuid.UUID]*domain.RentalExtension),
		breakdowns:     make(map[uuid.UUID]*domain.MechanicalBreakdown),
		fleetBaselines: make(map[uuid.UUID]float64),
		conflicts:      make(map[uuid.UUID]bool),
	}
}

func (m *MemoryRentalRepository) SetFleetBaseline(fleetID uuid.UUID, baseline float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.fleetBaselines[fleetID] = baseline
}

func (m *MemoryRentalRepository) SetConflict(fleetID uuid.UUID, hasConflict bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.conflicts[fleetID] = hasConflict
}

func (m *MemoryRentalRepository) CreateHandoverInspection(ctx context.Context, inspection *domain.HandoverInspection) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if inspection.ID == uuid.Nil {
		inspection.ID = uuid.New()
	}
	inspection.CreatedAt = time.Now().UTC()
	m.inspections[inspection.BookingID] = inspection
	return nil
}

func (m *MemoryRentalRepository) GetHandoverInspection(ctx context.Context, bookingID uuid.UUID) (*domain.HandoverInspection, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ins, ok := m.inspections[bookingID]
	if !ok {
		return nil, fmt.Errorf("handover inspection not found")
	}
	return ins, nil
}

func (m *MemoryRentalRepository) CreateHandoverChallenge(ctx context.Context, challenge *domain.HandoverChallenge) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if challenge.ID == uuid.Nil {
		challenge.ID = uuid.New()
	}
	challenge.CreatedAt = time.Now().UTC()
	m.challenges[challenge.BookingID] = challenge
	return nil
}

func (m *MemoryRentalRepository) GetActiveHandoverChallenge(ctx context.Context, bookingID uuid.UUID) (*domain.HandoverChallenge, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ch, ok := m.challenges[bookingID]
	if !ok {
		return nil, fmt.Errorf("active challenge PIN not found")
	}
	if ch.IsVerified {
		return nil, fmt.Errorf("challenge PIN already verified")
	}
	if time.Now().UTC().After(ch.ExpiresAt) {
		return nil, domain.ErrPINExpired
	}
	return ch, nil
}

func (m *MemoryRentalRepository) IncrementHandoverChallengeAttempts(ctx context.Context, challengeID uuid.UUID) (*domain.HandoverChallenge, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, ch := range m.challenges {
		if ch.ID == challengeID {
			ch.AttemptCount++
			return ch, nil
		}
	}
	return nil, fmt.Errorf("challenge not found")
}

func (m *MemoryRentalRepository) VerifyHandoverChallenge(ctx context.Context, challengeID uuid.UUID, touchSignatureBase64 string, lat, lng *float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, ch := range m.challenges {
		if ch.ID == challengeID {
			ch.IsVerified = true
			now := time.Now().UTC()
			ch.VerifiedAt = &now
			ch.TouchSignatureBase64 = touchSignatureBase64
			ch.ContractorLat = lat
			ch.ContractorLng = lng
			return nil
		}
	}
	return fmt.Errorf("challenge not found")
}

func (m *MemoryRentalRepository) CreateRentalShift(ctx context.Context, shift *domain.RentalShift) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if shift.ID == uuid.Nil {
		shift.ID = uuid.New()
	}
	shift.CreatedAt = time.Now().UTC()
	shift.UpdatedAt = time.Now().UTC()
	m.shifts[shift.ID] = shift
	return nil
}

func (m *MemoryRentalRepository) GetActiveRentalShift(ctx context.Context, bookingID uuid.UUID) (*domain.RentalShift, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, s := range m.shifts {
		if s.BookingID == bookingID && s.EndTime == nil {
			return s, nil
		}
	}
	return nil, fmt.Errorf("no active shift found")
}

func (m *MemoryRentalRepository) UpdateRentalShiftProgress(ctx context.Context, shiftID uuid.UUID, accumulatedHours, excessHours float64, endMeter *float64, isPaused bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.shifts[shiftID]
	if !ok {
		return fmt.Errorf("shift not found")
	}
	s.AccumulatedHours = accumulatedHours
	s.ExcessHours = excessHours
	s.EndMeterHours = endMeter
	s.IsPaused = isPaused
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (m *MemoryRentalRepository) EndRentalShift(ctx context.Context, shiftID uuid.UUID, endMeter float64, accumulatedHours, excessHours float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.shifts[shiftID]
	if !ok {
		return fmt.Errorf("shift not found")
	}
	now := time.Now().UTC()
	s.EndTime = &now
	s.EndMeterHours = &endMeter
	s.AccumulatedHours = accumulatedHours
	s.ExcessHours = excessHours
	s.UpdatedAt = now
	return nil
}

func (m *MemoryRentalRepository) ListRentalShifts(ctx context.Context, bookingID uuid.UUID) ([]domain.RentalShift, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []domain.RentalShift
	for _, s := range m.shifts {
		if s.BookingID == bookingID {
			res = append(res, *s)
		}
	}
	return res, nil
}

func (m *MemoryRentalRepository) CreateRentalExtension(ctx context.Context, ext *domain.RentalExtension) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if ext.ID == uuid.Nil {
		ext.ID = uuid.New()
	}
	ext.CreatedAt = time.Now().UTC()
	ext.UpdatedAt = time.Now().UTC()
	m.extensions[ext.ID] = ext
	return nil
}

func (m *MemoryRentalRepository) GetRentalExtension(ctx context.Context, extensionID uuid.UUID) (*domain.RentalExtension, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ext, ok := m.extensions[extensionID]
	if !ok {
		return nil, fmt.Errorf("rental extension not found")
	}
	return ext, nil
}

func (m *MemoryRentalRepository) ListRentalExtensions(ctx context.Context, bookingID uuid.UUID) ([]domain.RentalExtension, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []domain.RentalExtension
	for _, ext := range m.extensions {
		if ext.BookingID == bookingID {
			res = append(res, *ext)
		}
	}
	return res, nil
}

func (m *MemoryRentalRepository) UpdateRentalExtensionStatus(ctx context.Context, extensionID uuid.UUID, status string, notes *string, slipURI, slipHash *string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	ext, ok := m.extensions[extensionID]
	if !ok {
		return fmt.Errorf("rental extension not found")
	}
	ext.Status = status
	if notes != nil {
		ext.OwnerDecisionNotes = *notes
	}
	if slipURI != nil {
		ext.SlipStorageURI = *slipURI
	}
	if slipHash != nil {
		ext.SlipFileHash = *slipHash
	}
	if status == "CONFIRMED" {
		now := time.Now().UTC()
		ext.ConfirmedAt = &now
	}
	ext.UpdatedAt = time.Now().UTC()
	return nil
}

func (m *MemoryRentalRepository) ExpandBookingHirePeriod(ctx context.Context, bookingID uuid.UUID, newEndDate time.Time, additionalDays int, additionalGrossMinor, additionalPrepayMinor int64, expectedVersion int) error {
	return nil
}

func (m *MemoryRentalRepository) CreateMechanicalBreakdown(ctx context.Context, breakdown *domain.MechanicalBreakdown) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if breakdown.ID == uuid.Nil {
		breakdown.ID = uuid.New()
	}
	breakdown.CreatedAt = time.Now().UTC()
	breakdown.UpdatedAt = time.Now().UTC()
	m.breakdowns[breakdown.ID] = breakdown
	return nil
}

func (m *MemoryRentalRepository) GetMechanicalBreakdown(ctx context.Context, breakdownID uuid.UUID) (*domain.MechanicalBreakdown, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	bd, ok := m.breakdowns[breakdownID]
	if !ok {
		return nil, domain.ErrBreakdownNotFound
	}
	return bd, nil
}

func (m *MemoryRentalRepository) GetActiveMechanicalBreakdown(ctx context.Context, bookingID uuid.UUID) (*domain.MechanicalBreakdown, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, bd := range m.breakdowns {
		if bd.BookingID == bookingID && bd.Status == "ACTIVE_BREAKDOWN_PAUSED" {
			return bd, nil
		}
	}
	return nil, nil
}

func (m *MemoryRentalRepository) ListMechanicalBreakdowns(ctx context.Context, bookingID uuid.UUID) ([]domain.MechanicalBreakdown, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []domain.MechanicalBreakdown
	for _, bd := range m.breakdowns {
		if bd.BookingID == bookingID {
			res = append(res, *bd)
		}
	}
	return res, nil
}

func (m *MemoryRentalRepository) ResolveMechanicalBreakdown(ctx context.Context, breakdownID uuid.UUID, mechanicName, mechanicPhone, resolutionNotes string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	bd, ok := m.breakdowns[breakdownID]
	if !ok {
		return domain.ErrBreakdownNotFound
	}
	bd.Status = "RESOLVED"
	bd.MechanicName = mechanicName
	bd.MechanicPhone = mechanicPhone
	bd.ResolutionNotes = resolutionNotes
	now := time.Now().UTC()
	bd.ResolvedAt = &now
	bd.UpdatedAt = now
	return nil
}

func (m *MemoryRentalRepository) GetFleetBaselineMeter(ctx context.Context, fleetID uuid.UUID) (float64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	baseline, ok := m.fleetBaselines[fleetID]
	if !ok {
		return 0.0, nil
	}
	return baseline, nil
}

func (m *MemoryRentalRepository) CheckForwardBookingConflict(ctx context.Context, fleetID uuid.UUID, excludeBookingID uuid.UUID, currentEndDate, newEndDate time.Time) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.conflicts[fleetID], nil
}
