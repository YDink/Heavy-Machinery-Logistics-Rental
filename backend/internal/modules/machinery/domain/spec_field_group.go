package domain

import "github.com/google/uuid"

// SpecFieldGroup groups related technical specification attributes.
type SpecFieldGroup struct {
	ID            uuid.UUID
	CategoryID    uuid.UUID
	SectionNameEN string
	SectionNameAM string
	DisplayOrder  int
	Fields        []SpecFieldDefinition
}
