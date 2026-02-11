package models

import (
	"github.com/google/uuid"
	"time"
)

type StockMovementResponse struct {
	ID           uuid.UUID `json:"id"`
	InventoryID  uuid.UUID `json:"inventory_id"`
	ChangeAmount int       `json:"change_amount"`
	Type         string    `json:"type"`
	Reason       string    `json:"reason"`
	CreatedAt    time.Time `json:"created_at"`
}
type InitialStockRequest struct {
	WarehouseID uuid.UUID `json:"warehouse_id" binding:"required"`
	Quantity    int       `json:"quantity" binding:"required,min=0"`
}
