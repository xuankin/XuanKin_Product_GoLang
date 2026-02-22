package entity

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ProductAttributeValue struct {
	Base
	ProductAttributeID uuid.UUID      `gorm:"type:uuid;not null" json:"product_attribute_id"`
	Value              datatypes.JSON `gorm:"type:jsonb;not null" json:"value"`
}
