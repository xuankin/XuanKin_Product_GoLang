package models

import "github.com/google/uuid"

type CreateAttributeRequest struct {
	Name         map[string]interface{} `json:"name" binding:"required"`
	Type         string                 `json:"type" binding:"required,oneof=TEXT NUMBER BOOLEAN ENUM"`
	IsFilterable bool                   `json:"is_filterable"`
	IsRequired   bool                   `json:"is_required"`
}
type UpdateAttributeRequest struct {
	Name         map[string]interface{} `json:"name"`
	Type         string                 `json:"type" binding:"omitempty,oneof=TEXT NUMBER BOOLEAN ENUM"`
	IsFilterable *bool                  `json:"is_filterable"`
	IsRequired   *bool                  `json:"is_required"`
}
type AttributeResponse struct {
	ID           uuid.UUID              `json:"id"`
	Name         map[string]interface{} `json:"name"`
	Type         string                 `json:"type"`
	IsFilterable bool                   `json:"is_filterable"`
	IsRequired   bool                   `json:"is_required"`
}
