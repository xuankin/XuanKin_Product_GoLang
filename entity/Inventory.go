package entity

import "github.com/google/uuid"

type Inventory struct {
	Base
	OptionID         uuid.UUID `gorm:"type:uuid;not null" json:"option_id"`
	WarehouseID      uuid.UUID `gorm:"type:uuid;not null" json:"warehouse_id"`
	Quantity         int       `gorm:"default:0" json:"quantity"`
	ReservedQuantity int       `gorm:"default:0" json:"reserved_quantity"`
	Warehouse        Warehouse `gorm:"foreignKey:WarehouseID"`
}
