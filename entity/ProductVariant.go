package entity

import "github.com/google/uuid"

type ProductVariant struct {
	Base
	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	SKU       string    `gorm:"size:100;uniqueIndex;not null" json:"sku"`
	Price     float64   `gorm:"type:decimal(15,2)" json:"price"`
	SalePrice float64   `gorm:"type:decimal(15,2)" json:"sale_price"`
	Weight    float64   `gorm:"type:decimal(10,2)" json:"weight"`
	Status    string    `gorm:"type:varchar(20);default:'ACTIVE'" json:"status"`

	Inventories []Inventory      `gorm:"foreignKey:VariantID" json:"inventories"`
	Attributes  []AttributeValue `gorm:"many2many:variant_attribute_values;" json:"attributes"`
}
