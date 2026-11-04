package domain

import (
	"context"

	"github.com/google/uuid"
)

// UserRepository specifies data persistence invariants for User entities.
type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByPhoneNumber(ctx context.Context, phone string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	MarkPhoneVerified(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateAccountStatus(ctx context.Context, id uuid.UUID, status AccountStatus) (*User, error)
	UpdateKYCStatus(ctx context.Context, id uuid.UUID, status KYCStatus) (*User, error)
	IncrementTokenVersion(ctx context.Context, id uuid.UUID) (int32, error)
}
