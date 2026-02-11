package models

import (
	"github.com/google/uuid"
	"time"
)

type EsProductIndex struct {
	ID                uuid.UUID              `json:"id"`
	Name              map[string]interface{} `json:"name"`
	Slug              string                 `json:"slug"`
	Description       map[string]interface{} `json:"description"`
	CategoryID        uuid.UUID              `json:"category_id"`
	CategoryName      map[string]interface{} `json:"category_name"`
	BrandID           uuid.UUID              `json:"brand_id"`
	BrandName         map[string]interface{} `json:"brand_name"`
	Status            string                 `json:"status"`
	PrimaryImage      string                 `json:"primary_image"`
	MinPrice          float64                `json:"min_price"`
	MaxPrice          float64                `json:"max_price"`
	AttributesSummary []EsAttributeSummary   `json:"attributes_summary"`
	CreatedAt         time.Time              `json:"created_at"`
}

type EsAttributeSummary struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}
