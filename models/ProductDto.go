package models

import "github.com/google/uuid"

type CreateProductRequest struct {
	Name        map[string]interface{}    `json:"name" binding:"required"`
	Description map[string]interface{}    `json:"description"`
	CategoryID  uuid.UUID                 `json:"category_id" binding:"required"`
	BrandID     uuid.UUID                 `json:"brand_id" binding:"required"`
	Status      string                    `json:"status" binding:"oneof=ACTIVE INACTIVE"`
	Variants    []ProductVariantRequest   `json:"variants"`
	Slug        string                    `json:"slug"`
	Attributes  []ProductAttributeRequest `json:"attributes"`
}
type UpdateProductRequest struct {
	Name        map[string]interface{} `json:"name"`
	Description map[string]interface{} `json:"description"`
	CategoryID  uuid.UUID              `json:"category_id"`
	BrandID     uuid.UUID              `json:"brand_id"`
	Status      string                 `json:"status"`
}
type ProductResponse struct {
	ID          uuid.UUID                `json:"id"`
	Name        map[string]interface{}   `json:"name"`
	Slug        string                   `json:"slug"`
	Description map[string]interface{}   `json:"description"`
	Status      string                   `json:"status"`
	Category    CategoryResponse         `json:"category"`
	Media       []MediaResponse          `json:"media,omitempty"`
	Brand       BrandResponse            `json:"brand"`
	Variants    []VariantResponse        `json:"variants,omitempty"`
	Attributes  []ProductAttributeDetail `json:"attributes,omitempty"`
}
type ProductAttributeDetail struct {
	AttributeID   uuid.UUID              `json:"attribute_id"`
	AttributeName map[string]interface{} `json:"attribute_name"`
	Values        []interface{}          `json:"values"`
}
type ProductAttributeRequest struct {
	AttributeID uuid.UUID     `json:"attribute_id" binding:"required"`
	Values      []interface{} `json:"values" binding:"required,gt=0"`
}
