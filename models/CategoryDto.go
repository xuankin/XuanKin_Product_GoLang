package models

import "github.com/google/uuid"

type CreateCategoryRequest struct {
	Name     map[string]interface{} `json:"name" binding:"required"`
	ParentID *uuid.UUID             `json:"parent_id"`
}

type UpdateCategoryRequest struct {
	Name     map[string]interface{} `json:"name"`
	ParentID *uuid.UUID             `json:"parent_id"`
}

type CategoryResponse struct {
	ID       uuid.UUID              `json:"id"`
	Name     map[string]interface{} `json:"name"`
	ParentID *uuid.UUID             `json:"parent_id,omitempty"`
}
