package infrastructure

import (
	"context"
	"sync"
	"time"

	"github.com/dinkmachinery/dink-core/internal/modules/booking/domain"
	"github.com/google/uuid"
)

type leaseEntry struct {
	holderID  uuid.UUID
	expiresAt time.Time
}

// MemoryHoldLeaseManager provides an in-memory implementation of domain.HoldLeaseManager for local testing and zero-dep runtime.
type MemoryHoldLeaseManager struct {
	mu     sync.RWMutex
	leases map[string]leaseEntry
}

func NewMemoryHoldLeaseManager() *MemoryHoldLeaseManager {
	return &MemoryHoldLeaseManager{
		leases: make(map[string]leaseEntry),
	}
}

func (m *MemoryHoldLeaseManager) AcquireHoldLease(ctx context.Context, fleetID uuid.UUID, attachmentIDs []uuid.UUID, renterID uuid.UUID, ttl time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	// Clean up expired keys
	for k, v := range m.leases {
		if now.After(v.expiresAt) {
			delete(m.leases, k)
		}
	}

	fleetKey := "hold:fleet:" + fleetID.String()
	if entry, exists := m.leases[fleetKey]; exists && now.Before(entry.expiresAt) {
		return false, nil
	}

	for _, attID := range attachmentIDs {
		attKey := "hold:attachment:" + attID.String()
		if entry, exists := m.leases[attKey]; exists && now.Before(entry.expiresAt) {
			return false, nil
		}
	}

	// All locks available, acquire atomically
	expires := now.Add(ttl)
	m.leases[fleetKey] = leaseEntry{holderID: renterID, expiresAt: expires}
	for _, attID := range attachmentIDs {
		attKey := "hold:attachment:" + attID.String()
		m.leases[attKey] = leaseEntry{holderID: renterID, expiresAt: expires}
	}

	return true, nil
}

func (m *MemoryHoldLeaseManager) ExtendHoldLease(ctx context.Context, fleetID uuid.UUID, attachmentIDs []uuid.UUID, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	expires := now.Add(ttl)

	fleetKey := "hold:fleet:" + fleetID.String()
	if entry, exists := m.leases[fleetKey]; exists {
		entry.expiresAt = expires
		m.leases[fleetKey] = entry
	}

	for _, attID := range attachmentIDs {
		attKey := "hold:attachment:" + attID.String()
		if entry, exists := m.leases[attKey]; exists {
			entry.expiresAt = expires
			m.leases[attKey] = entry
		}
	}

	return nil
}

func (m *MemoryHoldLeaseManager) ReleaseHoldLease(ctx context.Context, fleetID uuid.UUID, attachmentIDs []uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.leases, "hold:fleet:"+fleetID.String())
	for _, attID := range attachmentIDs {
		delete(m.leases, "hold:attachment:"+attID.String())
	}

	return nil
}

func (m *MemoryHoldLeaseManager) IsFleetOnHold(ctx context.Context, fleetID uuid.UUID) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, exists := m.leases["hold:fleet:"+fleetID.String()]
	if !exists {
		return false, nil
	}

	if time.Now().After(entry.expiresAt) {
		return false, nil
	}

	return true, nil
}

var _ domain.HoldLeaseManager = (*MemoryHoldLeaseManager)(nil)
