package entity

import "github.com/google/uuid"

type VariantOptionValue struct {
	Base
	OptionID  uuid.UUID `gorm:"type:uuid;not null" json:"option_id"`
	Name      string    `gorm:"size:100" json:"name"`
	Value     string    `gorm:"size:100" json:"value"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
}
