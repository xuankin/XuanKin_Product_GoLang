package entity

import "gorm.io/datatypes"

type Warehouse struct {
	Base
	Name        datatypes.JSON `gorm:"type:jsonb;not null" json:"name"`
	Address     string         `gorm:"type:text" json:"address"`
	Phone       string         `gorm:"size:20" json:"phone"`
	Status      string         `gorm:"type:varchar(20);default:'ACTIVE'" json:"status"`
	Inventories []Inventory    `gorm:"foreignKey:WarehouseID" json:"inventories"`
}
