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
	conflicts         map[uuid.UUID]bool
	returnInspections map[uuid.UUID]*domain.ReturnInspection
	damageClaims      map[uuid.UUID]*domain.DamageClaimDocket
	settlements       map[uuid.UUID]*domain.RentalSettlement
	reviews           map[uuid.UUID][]*domain.BookingReview
}

func NewMemoryRentalRepository() *MemoryRentalRepository {
	return &MemoryRentalRepository{
		inspections:    make(map[uuid.UUID]*domain.HandoverInspection),
		challenges:     make(map[uuid.UUID]*domain.HandoverChallenge),
		shifts:         make(map[uuid.UUID]*domain.RentalShift),
		extensions:     make(map[uuid.UUID]*domain.RentalExtension),
		breakdowns:     make(map[uuid.UUID]*domain.MechanicalBreakdown),
		fleetBaselines: make(map[uuid.UUID]float64),
		conflicts:         make(map[uuid.UUID]bool),
		returnInspections: make(map[uuid.UUID]*domain.ReturnInspection),
		damageClaims:      make(map[uuid.UUID]*domain.DamageClaimDocket),
		settlements:       make(map[uuid.UUID]*domain.RentalSettlement),
		reviews:           make(map[uuid.UUID][]*domain.BookingReview),
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


// Return Inspection (Partition 05-07)
func (m *MemoryRentalRepository) CreateReturnInspection(ctx context.Context, inspection *domain.ReturnInspection) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if inspection.ID == uuid.Nil {
		inspection.ID = uuid.New()
	}
	inspection.CreatedAt = time.Now().UTC()
	m.returnInspections[inspection.BookingID] = inspection
	return nil
}

func (m *MemoryRentalRepository) GetReturnInspection(ctx context.Context, bookingID uuid.UUID) (*domain.ReturnInspection, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	insp, ok := m.returnInspections[bookingID]
	if !ok {
		return nil, fmt.Errorf("return inspection for booking %s not found", bookingID)
	}
	return insp, nil
}

// Damage Claims (Partition 05-07)
func (m *MemoryRentalRepository) CreateDamageClaimDocket(ctx context.Context, claim *domain.DamageClaimDocket) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if claim.ID == uuid.Nil {
		claim.ID = uuid.New()
	}
	now := time.Now().UTC()
	claim.CreatedAt = now
	claim.UpdatedAt = now
	m.damageClaims[claim.ID] = claim
	return nil
}

func (m *MemoryRentalRepository) GetDamageClaimDocket(ctx context.Context, claimID uuid.UUID) (*domain.DamageClaimDocket, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	claim, ok := m.damageClaims[claimID]
	if !ok {
		return nil, domain.ErrDamageClaimNotFound
	}
	return claim, nil
}

func (m *MemoryRentalRepository) GetDamageClaimDocketByBooking(ctx context.Context, bookingID uuid.UUID) (*domain.DamageClaimDocket, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var latest *domain.DamageClaimDocket
	for _, c := range m.damageClaims {
		if c.BookingID == bookingID {
			if latest == nil || c.CreatedAt.After(latest.CreatedAt) {
				latest = c
			}
		}
	}
	if latest == nil {
		return nil, domain.ErrDamageClaimNotFound
	}
	return latest, nil
}

func (m *MemoryRentalRepository) ListDamageClaimDockets(ctx context.Context, bookingID uuid.UUID) ([]domain.DamageClaimDocket, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []domain.DamageClaimDocket
	for _, c := range m.damageClaims {
		if c.BookingID == bookingID {
			list = append(list, *c)
		}
	}
	return list, nil
}

func (m *MemoryRentalRepository) RespondToDamageClaim(ctx context.Context, claimID uuid.UUID, status string, acceptedMinor int64, disputeReason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	claim, ok := m.damageClaims[claimID]
	if !ok {
		return domain.ErrDamageClaimNotFound
	}
	now := time.Now().UTC()
	claim.Status = status
	claim.AcceptedAmountMinor = acceptedMinor
	claim.DisputeReason = disputeReason
	claim.ContractorRespondedAt = &now
	claim.UpdatedAt = now
	return nil
}

func (m *MemoryRentalRepository) ArbitrateDamageClaim(ctx context.Context, claimID uuid.UUID, status string, acceptedMinor int64, arbitrationNotes string, arbitratedBy uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	claim, ok := m.damageClaims[claimID]
	if !ok {
		return domain.ErrDamageClaimNotFound
	}
	now := time.Now().UTC()
	claim.Status = status
	claim.AcceptedAmountMinor = acceptedMinor
	claim.ArbitrationNotes = arbitrationNotes
	claim.ArbitratedBy = &arbitratedBy
	claim.ArbitratedAt = &now
	claim.UpdatedAt = now
	return nil
}

func (m *MemoryRentalRepository) WaiveDamageClaim(ctx context.Context, claimID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	claim, ok := m.damageClaims[claimID]
	if !ok {
		return domain.ErrDamageClaimNotFound
	}
	now := time.Now().UTC()
	claim.Status = "WAIVED"
	claim.AcceptedAmountMinor = 0
	claim.UpdatedAt = now
	return nil
}

// Rental Settlement (Partition 05-07)
func (m *MemoryRentalRepository) UpsertRentalSettlement(ctx context.Context, settlement *domain.RentalSettlement) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if settlement.ID == uuid.Nil {
		settlement.ID = uuid.New()
	}
	now := time.Now().UTC()
	settlement.UpdatedAt = now
	if settlement.CreatedAt.IsZero() {
		settlement.CreatedAt = now
	}
	m.settlements[settlement.BookingID] = settlement
	return nil
}

func (m *MemoryRentalRepository) GetRentalSettlement(ctx context.Context, bookingID uuid.UUID) (*domain.RentalSettlement, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	settlement, ok := m.settlements[bookingID]
	if !ok {
		return nil, fmt.Errorf("settlement for booking %s not found", bookingID)
	}
	return settlement, nil
}

func (m *MemoryRentalRepository) UpdateRentalSettlementPaymentEvidence(ctx context.Context, bookingID uuid.UUID, slipURI, slipHash string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	settlement, ok := m.settlements[bookingID]
	if !ok {
		return fmt.Errorf("settlement for booking %s not found", bookingID)
	}
	now := time.Now().UTC()
	settlement.FinalSlipStorageURI = slipURI
	settlement.FinalSlipFileHash = slipHash
	settlement.Status = "PAYMENT_EVIDENCE_UPLOADED"
	settlement.UpdatedAt = now
	return nil
}

func (m *MemoryRentalRepository) ConfirmRentalSettlement(ctx context.Context, bookingID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	settlement, ok := m.settlements[bookingID]
	if !ok {
		return fmt.Errorf("settlement for booking %s not found", bookingID)
	}
	now := time.Now().UTC()
	settlement.Status = "SETTLED"
	settlement.SettledAt = &now
	settlement.UpdatedAt = now
	return nil
}

// Reviews (Partition 05-07)
func (m *MemoryRentalRepository) CreateBookingReview(ctx context.Context, review *domain.BookingReview) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if review.ID == uuid.Nil {
		review.ID = uuid.New()
	}
	review.CreatedAt = time.Now().UTC()
	m.reviews[review.BookingID] = append(m.reviews[review.BookingID], review)
	return nil
}

func (m *MemoryRentalRepository) GetBookingReviews(ctx context.Context, bookingID uuid.UUID) ([]domain.BookingReview, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []domain.BookingReview
	for _, r := range m.reviews[bookingID] {
		list = append(list, *r)
	}
	return list, nil
}

func (m *MemoryRentalRepository) PublishBookingReviews(ctx context.Context, bookingID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now().UTC()
	for _, r := range m.reviews[bookingID] {
		r.IsPublished = true
		r.PublishedAt = &now
	}
	return nil
}

func (m *MemoryRentalRepository) ListPublishedReviewsForTarget(ctx context.Context, targetID uuid.UUID) ([]domain.BookingReview, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []domain.BookingReview
	for _, reviews := range m.reviews {
		for _, r := range reviews {
			if r.TargetID == targetID && r.IsPublished {
				list = append(list, *r)
			}
		}
	}
	return list, nil
}

func (m *MemoryRentalRepository) UpdateBookingCompletionSummary(ctx context.Context, bookingID uuid.UUID, status string, checkOutMeter float64, damageClaimSubmitted, damageClaimWaived bool, acceptedDamageMinor, excessSurchargeMinor, finalNetBalanceMinor int64, expectedVersion int) error {
	// In-memory stub succeeds unconditionally
	return nil
}
