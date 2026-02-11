package models

import "github.com/google/uuid"

type CreateWarehouseRequest struct {
	Name    map[string]interface{} `json:"name" binding:"required"`
	Address string                 `json:"address"`
	Phone   string                 `json:"phone"`
}
type UpdateWarehouseRequest struct {
	Name    map[string]interface{} `json:"name" binding:"required"`
	Address string                 `json:"address"`
	Phone   string                 `json:"phone"`
	Status  string                 `json:"status"`
}
type WarehouseResponse struct {
	ID      uuid.UUID              `json:"id"`
	Name    map[string]interface{} `json:"name"`
	Address string                 `json:"address"`
	Phone   string                 `json:"phone"`
	Status  string                 `json:"status"`
}

type UpdateInventoryRequest struct {
	VariantID   uuid.UUID `json:"variant_id" binding:"required"`
	WarehouseID uuid.UUID `json:"warehouse_id" binding:"required"`
	Amount      int       `json:"amount" binding:"required"`
	Type        string    `json:"type" binding:"required,oneof=IN OUT ADJUST"`
	Reason      string    `json:"reason"`
}

type InventoryResponse struct {
	ID               uuid.UUID         `json:"id"`
	VariantID        uuid.UUID         `json:"variant_id"`
	Warehouse        WarehouseResponse `json:"warehouse"`
	Quantity         int               `json:"available_quantity"`
	ReservedQuantity int               `json:"reserved_quantity"`
}
