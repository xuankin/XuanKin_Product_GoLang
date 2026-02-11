package models

import "github.com/google/uuid"

type CreateAttributeRequest struct {
	Name map[string]interface{} `json:"name" binding:"required"`
}
type UpdateAttributeRequest struct {
	Name map[string]interface{} `json:"name" binding:"required"`
}
type AttributeResponse struct {
	ID     uuid.UUID                `json:"id"`
	Name   map[string]interface{}   `json:"name"`
	Values []AttributeValueResponse `json:"values,omitempty"`
}
