package domain

import "github.com/google/uuid"

// SpecFieldDefinition defines dynamic EAV schema metadata.
type SpecFieldDefinition struct {
	ID                  uuid.UUID
	GroupID             uuid.UUID
	FieldKey            string
	LabelEN             string
	LabelAM             string
	FieldType           string // 'STRING', 'DECIMAL', 'INTEGER', 'ENUM', 'BOOLEAN'
	UnitSymbol          *string
	IsRequired          bool
	IsLogisticsCritical bool
	IsComparable        bool
	IsVisibleToRenter   bool
	MinValue            *float64
	MaxValue            *float64
	EnumValuesJSON      []string
	DisplayOrder        int
}
