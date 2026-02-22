package entity

import "github.com/google/uuid"

type VariantOption struct {
	Base
	VariantID   uuid.UUID            `gorm:"type:uuid;not null" json:"variant_id"`
	SKU         string               `gorm:"size:100;uniqueIndex;not null" json:"sku"`
	Price       float64              `gorm:"type:decimal(15,2)" json:"price"`
	SalePrice   float64              `gorm:"type:decimal(15,2)" json:"sale_price"`
	Weight      float64              `gorm:"type:decimal(10,2)" json:"weight"`
	Status      string               `gorm:"type:varchar(20);default:'ACTIVE'" json:"status"`
	Values      []VariantOptionValue `gorm:"foreignKey:OptionID;constraint:OnDelete:CASCADE;" json:"values"`
	Inventories []Inventory          `gorm:"foreignKey:OptionID;constraint:OnDelete:CASCADE;" json:"inventories"`
	Media       []Media              `gorm:"foreignKey:OptionID;constraint:OnDelete:CASCADE;" json:"media"`
}
