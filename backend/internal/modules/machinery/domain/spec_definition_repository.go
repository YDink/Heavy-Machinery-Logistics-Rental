package domain

import (
	"context"

	"github.com/google/uuid"
)

type SpecDefinitionRepository interface {
	GetCategorySpecDefinitions(ctx context.Context, categoryID uuid.UUID) ([]SpecFieldGroup, error)
}
