package models

import "github.com/google/uuid"

type CreateMediaRequest struct {
	ProductID    uuid.UUID  `json:"product_id" binding:"required"`
	VariantID    *uuid.UUID `json:"variant_id"`
	OptionID     *uuid.UUID `json:"option_id"`
	Type         string     `json:"type" binding:"required,oneof=IMAGE VIDEO"`
	URL          string     `form:"url"`
	ThumbnailURL string     `json:"thumbnail_url"`
	IsPrimary    bool       `json:"is_primary"`
	SortOrder    int        `json:"sort_order"`
}

type MediaResponse struct {
	ID           uuid.UUID `json:"id"`
	Type         string    `json:"type"`
	URL          string    `json:"url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	IsPrimary    bool      `json:"is_primary"`
	SortOrder    int       `json:"sort_order"`
}
