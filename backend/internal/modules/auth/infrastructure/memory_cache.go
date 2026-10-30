package infrastructure

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/domain"
)

// MemoryCache provides a thread-safe in-memory cache for OTP challenges and Token rotation.
type MemoryCache struct {
	mu           sync.RWMutex
	challenges   map[string]*domain.OTPChallenge
	rateLimits   map[string][]time.Time
	tokens       map[string]*domain.TokenFamilyMetadata
	families     map[string][]string // familyID -> list of tokens
	revokedFams  map[string]bool
	blacklisted  map[string]time.Time
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		challenges:  make(map[string]*domain.OTPChallenge),
		rateLimits:  make(map[string][]time.Time),
		tokens:      make(map[string]*domain.TokenFamilyMetadata),
		families:    make(map[string][]string),
		revokedFams: make(map[string]bool),
		blacklisted: make(map[string]time.Time),
	}
}

// SaveChallenge implements domain.OTPRepository.
func (c *MemoryCache) SaveChallenge(ctx context.Context, challenge *domain.OTPChallenge, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.challenges[challenge.SessionID] = challenge
	return nil
}

// GetChallenge implements domain.OTPRepository.
func (c *MemoryCache) GetChallenge(ctx context.Context, sessionID string) (*domain.OTPChallenge, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ch, exists := c.challenges[sessionID]
	if !exists {
		return nil, domain.ErrOTPExpired
	}
	return ch, nil
}

// IncrementAttempts implements domain.OTPRepository.
func (c *MemoryCache) IncrementAttempts(ctx context.Context, sessionID string) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ch, exists := c.challenges[sessionID]
	if !exists {
		return 0, domain.ErrOTPExpired
	}
	ch.Attempts++
	return ch.Attempts, nil
}

// DeleteChallenge implements domain.OTPRepository.
func (c *MemoryCache) DeleteChallenge(ctx context.Context, sessionID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.challenges, sessionID)
	return nil
}

// CheckAndIncrementRateLimit implements domain.OTPRepository.
func (c *MemoryCache) CheckAndIncrementRateLimit(ctx context.Context, phone string, window time.Duration, maxRequests int) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-window)

	timestamps := c.rateLimits[phone]
	var valid []time.Time
	for _, t := range timestamps {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= maxRequests {
		c.rateLimits[phone] = valid
		return false, nil
	}

	valid = append(valid, now)
	c.rateLimits[phone] = valid
	return true, nil
}

// SaveRefreshToken implements domain.TokenStore.
func (c *MemoryCache) SaveRefreshToken(ctx context.Context, token string, familyID string, userID uuid.UUID, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	metadata := &domain.TokenFamilyMetadata{
		FamilyID:    familyID,
		UserID:      userID,
		ActiveToken: token,
		IsRevoked:   c.revokedFams[familyID],
		CreatedAt:   time.Now(),
	}
	c.tokens[token] = metadata
	c.families[familyID] = append(c.families[familyID], token)
	return nil
}

// GetRefreshToken implements domain.TokenStore.
func (c *MemoryCache) GetRefreshToken(ctx context.Context, token string) (*domain.TokenFamilyMetadata, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	metadata, exists := c.tokens[token]
	if !exists {
		return nil, domain.ErrInvalidToken
	}

	if c.revokedFams[metadata.FamilyID] {
		metadata.IsRevoked = true
	}

	return metadata, nil
}

// RotateRefreshToken implements domain.TokenStore.
func (c *MemoryCache) RotateRefreshToken(ctx context.Context, oldToken string, newToken string, familyID string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if meta, exists := c.tokens[oldToken]; exists {
		meta.ActiveToken = newToken
	}

	return nil
}

// RevokeFamily implements domain.TokenStore.
func (c *MemoryCache) RevokeFamily(ctx context.Context, familyID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.revokedFams[familyID] = true
	for _, tok := range c.families[familyID] {
		if meta, exists := c.tokens[tok]; exists {
			meta.IsRevoked = true
		}
	}
	return nil
}

// IsBlacklisted implements domain.TokenStore.
func (c *MemoryCache) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	exp, exists := c.blacklisted[tokenID]
	if !exists {
		return false, nil
	}
	if time.Now().After(exp) {
		return false, nil
	}
	return true, nil
}

// BlacklistToken implements domain.TokenStore.
func (c *MemoryCache) BlacklistToken(ctx context.Context, tokenID string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.blacklisted[tokenID] = time.Now().Add(ttl)
	return nil
}

// MemoryUserRepository provides an in-memory thread-safe implementation of domain.UserRepository.
type MemoryUserRepository struct {
	mu    sync.RWMutex
	users map[uuid.UUID]*domain.User
	byPh  map[string]*domain.User
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make(map[uuid.UUID]*domain.User),
		byPh:  make(map[string]*domain.User),
	}
}

func (r *MemoryUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	return u, nil
}

func (r *MemoryUserRepository) GetByPhoneNumber(ctx context.Context, phone string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, exists := r.byPh[phone]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	return u, nil
}

func (r *MemoryUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	r.byPh[user.PhoneNumber.String()] = user
	return user, nil
}

func (r *MemoryUserRepository) MarkPhoneVerified(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	now := time.Now()
	u.AccountStatus = domain.AccountStatusActive
	u.PhoneVerifiedAt = &now
	return u, nil
}

func (r *MemoryUserRepository) UpdateAccountStatus(ctx context.Context, id uuid.UUID, status domain.AccountStatus) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	u.AccountStatus = status
	return u, nil
}

func (r *MemoryUserRepository) UpdateKYCStatus(ctx context.Context, id uuid.UUID, status domain.KYCStatus) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	u.KYCStatus = status
	return u, nil
}

func (r *MemoryUserRepository) IncrementTokenVersion(ctx context.Context, id uuid.UUID) (int32, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, exists := r.users[id]
	if !exists {
		return 0, domain.ErrUserNotFound
	}
	u.TokenVersion++
	return u.TokenVersion, nil
}
