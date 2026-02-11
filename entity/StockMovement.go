package entity

import "github.com/google/uuid"

type StockMovement struct {
	Base
	InventoryID  uuid.UUID `gorm:"type:uuid;not null" json:"inventory_id"`
	ChangeAmount int       `gorm:"not null" json:"change_amount"`
	Type         string    `gorm:"type:varchar(20)" json:"type"`
	Reason       string    `gorm:"size:255" json:"reason"`
}
