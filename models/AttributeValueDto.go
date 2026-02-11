package models

import "github.com/google/uuid"

type CreateAttributeValueRequest struct {
	AttributeID uuid.UUID              `json:"attribute_id" binding:"required"`
	Value       map[string]interface{} `json:"value" binding:"required"`
}

type AttributeValueResponse struct {
	ID            uuid.UUID              `json:"id"`
	AttributeID   uuid.UUID              `json:"attribute_id"`
	AttributeName map[string]interface{} `json:"attribute_name,omitempty"`
	Value         map[string]interface{} `json:"value"`
}
