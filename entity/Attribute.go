package entity

import "gorm.io/datatypes"

type Attribute struct {
	Base
	Name         datatypes.JSON `gorm:"type:jsonb;not null" json:"name"`
	Type         string         `gorm:"type:varchar(50)" json:"type"`
	IsFilterable bool           `gorm:"default:false" json:"is_filterable"`
	IsRequired   bool           `gorm:"default:false" json:"is_required"`
}
