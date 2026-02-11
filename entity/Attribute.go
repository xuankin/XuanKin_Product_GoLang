package entity

import "gorm.io/datatypes"

type Attribute struct {
	Base
	Name   datatypes.JSON   `gorm:"type:jsonb;not null" json:"name"`
	Values []AttributeValue `gorm:"foreignKey:AttributeID" json:"values"`
}
