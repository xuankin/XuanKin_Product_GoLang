package entity

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AttributeValue struct {
	Base
	AttributeID uuid.UUID      `gorm:"type:uuid;not null" json:"attribute_id"`
	Attribute   Attribute      `gorm:"foreignKey:AttributeID" json:"attribute"`
	Value       datatypes.JSON `gorm:"type:jsonb;not null" json:"value"`
}
