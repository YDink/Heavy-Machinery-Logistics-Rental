package domain

import (
	"context"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	ListAll(ctx context.Context) ([]Category, error)
	ListRoots(ctx context.Context) ([]Category, error)
	ListSubcategories(ctx context.Context, parentID uuid.UUID) ([]Category, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*Category, error)
}
