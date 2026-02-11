package entity

import "github.com/google/uuid"

type Media struct {
	Base
	ProductID    uuid.UUID  `gorm:"type:uuid;not null" json:"product_id"`
	VariantID    *uuid.UUID `gorm:"type:uuid" json:"variant_id"`
	Type         string     `gorm:"type:varchar(20)" json:"type"`
	URL          string     `gorm:"size:255;not null" json:"url"`
	ThumbnailURL string     `gorm:"size:255" json:"thumbnail_url"`
	IsPrimary    bool       `gorm:"default:false" json:"is_primary"`
	SortOrder    int        `gorm:"default:0" json:"sort_order"`
}
