package domain

import "github.com/google/uuid"

// Category represents a node in the hierarchical machinery taxonomy.
type Category struct {
	ID           uuid.UUID
	ParentID     *uuid.UUID
	Code         string
	NameEN       string
	NameAM       string
	IconURL      string
	IsActive     bool
	DisplayOrder int
}
