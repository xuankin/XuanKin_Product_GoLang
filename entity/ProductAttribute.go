package entity

import "github.com/google/uuid"

type ProductAttribute struct {
	Base
	ProductID   uuid.UUID               `gorm:"type:uuid;not null" json:"product_id"`
	AttributeID uuid.UUID               `gorm:"type:uuid;not null" json:"attribute_id"`
	Attribute   Attribute               `gorm:"foreignKey:AttributeID" json:"attribute"`
	Values      []ProductAttributeValue `gorm:"foreignKey:ProductAttributeID;constraint:OnDelete:CASCADE;" json:"values"`
}
