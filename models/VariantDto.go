package models

import "github.com/google/uuid"

type VariantOptionRequest struct {
	SKU       string                      `json:"sku" binding:"required"`
	Price     float64                     `json:"price" binding:"required"`
	SalePrice float64                     `json:"sale_price" binding:"required"`
	Weight    float64                     `json:"weight"`
	Values    []VariantOptionValueRequest `json:"values"`

	Inventories []InitialStockRequest `json:"inventories"`
}
type CreateVariantRequest struct {
	ProductID uuid.UUID              `json:"product_id" binding:"required"`
	Code      string                 `json:"code" binding:"required"`
	Name      map[string]interface{} `json:"name" binding:"required"`
	Options   []VariantOptionRequest `json:"options" binding:"required"`
}
type VariantOptionValueRequest struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	SortOrder int    `json:"sort_order"`
}
type ProductVariantRequest struct {
	Code    string                 `json:"code" binding:"required"`
	Name    map[string]interface{} `json:"name"`
	Status  string                 `json:"status"`
	Options []VariantOptionRequest `json:"options" binding:"required"`
}
type VariantOptionResponse struct {
	ID          uuid.UUID                    `json:"id"`
	SKU         string                       `json:"sku"`
	Price       float64                      `json:"price"`
	SalePrice   float64                      `json:"sale_price"`
	Weight      float64                      `json:"weight"`
	Status      string                       `json:"status"`
	Values      []VariantOptionValueResponse `json:"values"`
	Media       []MediaResponse              `json:"media,omitempty"`
	Inventories []InventoryResponse          `json:"inventories,omitempty"`
}
type VariantOptionValueResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	SortOrder int       `json:"sort_order"`
}
type VariantResponse struct {
	ID      uuid.UUID               `json:"id"`
	Code    string                  `json:"code"`
	Name    map[string]interface{}  `json:"name"`
	Status  string                  `json:"status"`
	Options []VariantOptionResponse `json:"options"`
}
type UpdateVariantRequest struct {
	Code   string                 `json:"code"`
	Name   map[string]interface{} `json:"name"`
	Status string                 `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}
type UpdateVariantOptionRequest struct {
	SKU       string  `json:"sku"`
	Price     float64 `json:"price"`
	SalePrice float64 `json:"sale_price"`
	Weight    float64 `json:"weight"`
	Status    string  `json:"status"`
}
