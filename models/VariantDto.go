package models

import "github.com/google/uuid"

type CreateVariantRequest struct {
	ProductID     uuid.UUID             `json:"product_id" binding:"required"`
	SKU           string                `json:"sku" binding:"required"`
	Price         float64               `json:"price" binding:"required,gt=0"`
	SalePrice     float64               `json:"sale_price"`
	Weight        float64               `json:"weight"`
	AttributeIds  []uuid.UUID           `json:"attribute_value_ids"`
	InitialStocks []InitialStockRequest `json:"initial_stocks"`
}

type UpdateVariantRequest struct {
	Price     *float64 `json:"price"`
	SalePrice *float64 `json:"sale_price"`
	Weight    *float64 `json:"weight"`
	Status    string   `json:"status"`
}

type VariantResponse struct {
	ID         uuid.UUID                `json:"id"`
	SKU        string                   `json:"sku"`
	Price      float64                  `json:"price"`
	SalePrice  float64                  `json:"sale_price"`
	Weight     float64                  `json:"weight"`
	Status     string                   `json:"status"`
	Attributes []AttributeValueResponse `json:"attributes"`
}
