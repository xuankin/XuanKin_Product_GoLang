package models

import "github.com/google/uuid"

type CreateMediaRequest struct {
	ProductID    *uuid.UUID `form:"product_id"`
	VariantID    *uuid.UUID `form:"variant_id"`
	OptionID     *uuid.UUID `form:"option_id"`
	ThumbnailURL string     `form:"thumbnail_url"`
	IsPrimary    bool       `form:"is_primary"`
	SortOrder    int        `form:"sort_order"`
}
type MediaResponse struct {
	ID           uuid.UUID `json:"id"`
	Type         string    `json:"type"`
	URL          string    `json:"url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	IsPrimary    bool      `json:"is_primary"`
	SortOrder    int       `json:"sort_order"`
}
