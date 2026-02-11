package entity

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Product struct {
	Base
	Name        datatypes.JSON `gorm:"type:jsonb;not null" json:"name"`
	Slug        string         `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	Description datatypes.JSON `gorm:"type:jsonb" json:"description"`
	CategoryID  uuid.UUID      `gorm:"type:uuid;not null" json:"category_id"`
	BrandID     uuid.UUID      `gorm:"type:uuid;not null" json:"brand_id"`
	Status      string         `gorm:"type:varchar(20);default:'ACTIVE'" json:"status"`

	Category Category         `gorm:"foreignKey:CategoryID" json:"category"`
	Brand    Brand            `gorm:"foreignKey:BrandID" json:"brand"`
	Variants []ProductVariant `gorm:"foreignKey:ProductID" json:"variants"`
	Media    []Media          `gorm:"foreignKey:ProductID" json:"media"`
}
