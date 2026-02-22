package entity

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ProductVariant struct {
	Base
	ProductID uuid.UUID       `gorm:"type:uuid;not null" json:"product_id"`
	Code      string          `gorm:"size:100;not null" json:"code"`
	Name      datatypes.JSON  `gorm:"type:jsonb" json:"name"`
	Status    string          `gorm:"type:varchar(20);default:'ACTIVE'" json:"status"`
	Options   []VariantOption `gorm:"foreignKey:VariantID;constraint:OnDelete:CASCADE;" json:"options"`
	Media     []Media         `gorm:"foreignKey:VariantID;constraint:OnDelete:CASCADE;" json:"media"`
}
