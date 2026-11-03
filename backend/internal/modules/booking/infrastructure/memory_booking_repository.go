package infrastructure

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	"github.com/google/uuid"
)

type MemoryBookingRepository struct {
	mu          sync.RWMutex
	bookings    map[uuid.UUID]*domain.Booking
	attachments map[uuid.UUID][]domain.BookingAttachment
	slips       map[uuid.UUID][]*domain.BankSlip
	auditLogs   map[uuid.UUID][]*domain.AuditLog
}

func NewMemoryBookingRepository() *MemoryBookingRepository {
	return &MemoryBookingRepository{
		bookings:    make(map[uuid.UUID]*domain.Booking),
		attachments: make(map[uuid.UUID][]domain.BookingAttachment),
		slips:       make(map[uuid.UUID][]*domain.BankSlip),
		auditLogs:   make(map[uuid.UUID][]*domain.AuditLog),
	}
}

func (r *MemoryBookingRepository) CreateBooking(ctx context.Context, b *domain.Booking) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.bookings[b.ID]; exists {
		return fmt.Errorf("booking with id %s already exists", b.ID)
	}

	cpy := *b
	r.bookings[b.ID] = &cpy
	return nil
}

func (r *MemoryBookingRepository) GetBookingByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	b, ok := r.bookings[id]
	if !ok {
		return nil, domain.ErrBookingNotFound
	}

	cpy := *b
	cpy.Attachments = append([]domain.BookingAttachment(nil), r.attachments[id]...)
	return &cpy, nil
}

func (r *MemoryBookingRepository) GetBookingByReference(ctx context.Context, ref string) (*domain.Booking, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, b := range r.bookings {
		if b.BookingReference == ref {
			cpy := *b
			cpy.Attachments = append([]domain.BookingAttachment(nil), r.attachments[b.ID]...)
			return &cpy, nil
		}
	}
	return nil, domain.ErrBookingNotFound
}

func (r *MemoryBookingRepository) ListBookingsByRenter(ctx context.Context, renterID uuid.UUID, status string, limit int) ([]*domain.Booking, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Booking
	for _, b := range r.bookings {
		if b.RenterID == renterID {
			if status == "" || string(b.Status) == status {
				cpy := *b
				cpy.Attachments = append([]domain.BookingAttachment(nil), r.attachments[b.ID]...)
				result = append(result, &cpy)
			}
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}

func (r *MemoryBookingRepository) ListBookingsByOwner(ctx context.Context, ownerID uuid.UUID, status string, limit int) ([]*domain.Booking, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Booking
	for _, b := range r.bookings {
		if b.OwnerID == ownerID {
			if status == "" || string(b.Status) == status {
				cpy := *b
				cpy.Attachments = append([]domain.BookingAttachment(nil), r.attachments[b.ID]...)
				result = append(result, &cpy)
			}
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}

func (r *MemoryBookingRepository) UpdateBookingStatusOCC(ctx context.Context, b *domain.Booking, expectedVersion int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.bookings[b.ID]
	if !ok {
		return domain.ErrBookingNotFound
	}

	if existing.Version != expectedVersion {
		return domain.ErrConcurrencyConflict
	}

	b.Version = existing.Version + 1
	b.UpdatedAt = time.Now()
	cpy := *b
	r.bookings[b.ID] = &cpy
	return nil
}

func (r *MemoryBookingRepository) UpdateBookingHandover(ctx context.Context, bookingID uuid.UUID, checkInMeter float64, expectedVersion int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.bookings[bookingID]
	if !ok {
		return domain.ErrBookingNotFound
	}

	if existing.Version != expectedVersion {
		return domain.ErrConcurrencyConflict
	}

	now := time.Now()
	existing.Status = domain.StatusActive
	existing.CheckInMeterHours = &checkInMeter
	existing.ActivatedAt = &now
	existing.Version++
	existing.UpdatedAt = now
	return nil
}

func (r *MemoryBookingRepository) UpdateBookingReturn(ctx context.Context, bookingID uuid.UUID, status domain.BookingStatus, checkOutMeter float64, expectedVersion int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.bookings[bookingID]
	if !ok {
		return domain.ErrBookingNotFound
	}

	if existing.Version != expectedVersion {
		return domain.ErrConcurrencyConflict
	}

	now := time.Now()
	existing.Status = status
	existing.CheckOutMeterHours = &checkOutMeter
	if status == domain.StatusCompleted {
		existing.CompletedAt = &now
	}
	existing.Version++
	existing.UpdatedAt = now
	return nil
}

func (r *MemoryBookingRepository) AddBookingAttachment(ctx context.Context, bookingID, attachmentID uuid.UUID, dailyRateMinor int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	item := domain.BookingAttachment{
		ID:                    uuid.New(),
		BookingID:             bookingID,
		AttachmentID:          attachmentID,
		BundledDailyRateMinor: dailyRateMinor,
		CreatedAt:             time.Now(),
	}
	r.attachments[bookingID] = append(r.attachments[bookingID], item)
	return nil
}

func (r *MemoryBookingRepository) ListBookingAttachments(ctx context.Context, bookingID uuid.UUID) ([]domain.BookingAttachment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items, ok := r.attachments[bookingID]
	if !ok {
		return []domain.BookingAttachment{}, nil
	}
	return append([]domain.BookingAttachment(nil), items...), nil
}

func (r *MemoryBookingRepository) CreateBankSlip(ctx context.Context, slip *domain.BankSlip) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cpy := *slip
	r.slips[slip.BookingID] = append(r.slips[slip.BookingID], &cpy)
	return nil
}

func (r *MemoryBookingRepository) GetLatestBankSlip(ctx context.Context, bookingID uuid.UUID) (*domain.BankSlip, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items, ok := r.slips[bookingID]
	if !ok || len(items) == 0 {
		return nil, domain.ErrSlipNotFound
	}
	cpy := *items[len(items)-1]
	return &cpy, nil
}

func (r *MemoryBookingRepository) UpdateBankSlipStatus(ctx context.Context, slipID uuid.UUID, status domain.BankSlipStatus, reasonCode, notes *string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for _, slipList := range r.slips {
		for _, s := range slipList {
			if s.ID == slipID {
				s.Status = status
				s.RejectionReasonCode = reasonCode
				s.RejectionNotes = notes
				s.AuditedAt = &now
				return nil
			}
		}
	}
	return domain.ErrSlipNotFound
}

func (r *MemoryBookingRepository) CreateAuditLog(ctx context.Context, log *domain.AuditLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cpy := *log
	r.auditLogs[log.BookingID] = append(r.auditLogs[log.BookingID], &cpy)
	return nil
}

func (r *MemoryBookingRepository) ListAuditLogs(ctx context.Context, bookingID uuid.UUID) ([]*domain.AuditLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items, ok := r.auditLogs[bookingID]
	if !ok {
		return []*domain.AuditLog{}, nil
	}

	var result []*domain.AuditLog
	for _, l := range items {
		cpy := *l
		result = append(result, &cpy)
	}
	return result, nil
}

var _ domain.BookingRepository = (*MemoryBookingRepository)(nil)
