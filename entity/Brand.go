package entity

import "gorm.io/datatypes"

type Brand struct {
	Base
	Name     datatypes.JSON `gorm:"type:jsonb;not null" json:"name"`
	Logo     string         `gorm:"size:255" json:"logo"`
	Products []Product      `gorm:"foreignKey:BrandID" json:"products,omitempty"`
}
