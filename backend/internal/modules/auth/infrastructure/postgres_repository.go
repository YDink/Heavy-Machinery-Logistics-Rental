package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/dinkmachinery/dink-core/internal/modules/auth/domain"
	"github.com/dinkmachinery/dink-core/internal/platform/database/db"
)

type PostgresUserRepository struct {
	queries *db.Queries
}

func NewPostgresUserRepository(queries *db.Queries) *PostgresUserRepository {
	return &PostgresUserRepository{queries: queries}
}

func toDomainUser(u db.User) (*domain.User, error) {
	phone, err := domain.ParsePhoneNumber(u.PhoneNumber, u.CountryCode)
	if err != nil {
		return nil, fmt.Errorf("failed to parse stored phone number: %w", err)
	}

	var phoneVerifiedAt *time.Time
	if u.PhoneVerifiedAt.Valid {
		phoneVerifiedAt = &u.PhoneVerifiedAt.Time
	}

	return &domain.User{
		ID:              uuid.UUID(u.ID.Bytes),
		PhoneNumber:     phone,
		CountryCode:     u.CountryCode,
		FullName:        u.FullName,
		PrimaryRole:     domain.UserRole(u.PrimaryRole),
		AccountStatus:   domain.AccountStatus(u.AccountStatus),
		KYCStatus:       domain.KYCStatus(u.KycStatus),
		PasswordHash:    u.PasswordHash,
		TokenVersion:    u.TokenVersion,
		PhoneVerifiedAt: phoneVerifiedAt,
		CreatedAt:       u.CreatedAt.Time,
		UpdatedAt:       u.UpdatedAt.Time,
	}, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	pgID := pgtype.UUID{Bytes: id, Valid: true}
	row, err := r.queries.GetUserByID(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("query GetUserByID failed: %w", err)
	}
	return toDomainUser(row)
}

func (r *PostgresUserRepository) GetByPhoneNumber(ctx context.Context, phone string) (*domain.User, error) {
	row, err := r.queries.GetUserByPhoneNumber(ctx, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("query GetUserByPhoneNumber failed: %w", err)
	}
	return toDomainUser(row)
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	params := db.CreateUserParams{
		PhoneNumber:   user.PhoneNumber.String(),
		CountryCode:   user.CountryCode,
		FullName:      user.FullName,
		PrimaryRole:   string(user.PrimaryRole),
		AccountStatus: string(user.AccountStatus),
		KycStatus:     string(user.KYCStatus),
		PasswordHash:  user.PasswordHash,
	}

	row, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("query CreateUser failed: %w", err)
	}
	return toDomainUser(row)
}

func (r *PostgresUserRepository) MarkPhoneVerified(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	pgID := pgtype.UUID{Bytes: id, Valid: true}
	row, err := r.queries.UpdateUserPhoneVerified(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("query UpdateUserPhoneVerified failed: %w", err)
	}
	return toDomainUser(row)
}

func (r *PostgresUserRepository) UpdateAccountStatus(ctx context.Context, id uuid.UUID, status domain.AccountStatus) (*domain.User, error) {
	params := db.UpdateUserAccountStatusParams{
		ID:            pgtype.UUID{Bytes: id, Valid: true},
		AccountStatus: string(status),
	}
	row, err := r.queries.UpdateUserAccountStatus(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("query UpdateUserAccountStatus failed: %w", err)
	}
	return toDomainUser(row)
}

func (r *PostgresUserRepository) UpdateKYCStatus(ctx context.Context, id uuid.UUID, status domain.KYCStatus) (*domain.User, error) {
	params := db.UpdateUserKYCStatusParams{
		ID:        pgtype.UUID{Bytes: id, Valid: true},
		KycStatus: string(status),
	}
	row, err := r.queries.UpdateUserKYCStatus(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("query UpdateUserKYCStatus failed: %w", err)
	}
	return toDomainUser(row)
}

func (r *PostgresUserRepository) IncrementTokenVersion(ctx context.Context, id uuid.UUID) (int32, error) {
	pgID := pgtype.UUID{Bytes: id, Valid: true}
	version, err := r.queries.IncrementUserTokenVersion(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, domain.ErrUserNotFound
		}
		return 0, fmt.Errorf("query IncrementUserTokenVersion failed: %w", err)
	}
	return version, nil
}
