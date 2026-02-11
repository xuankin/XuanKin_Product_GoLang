package models

import "github.com/google/uuid"

type CreateBrandRequest struct {
	Name map[string]interface{} `json:"name" binding:"required"`
	Logo string                 `json:"logo"`
}

type UpdateBrandRequest struct {
	Name map[string]interface{} `json:"name"`
	Logo string                 `json:"logo"`
}

type BrandResponse struct {
	ID   uuid.UUID              `json:"id"`
	Name map[string]interface{} `json:"name"`
	Logo string                 `json:"logo"`
}
