package entity

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Category struct {
	Base
	Name     datatypes.JSON `gorm:"type:jsonb;not null" json:"name"`
	ParentId *uuid.UUID     `gorm:"type:uuid;" json:"parent_id"`
	Products []Product      `gorm:"foreignKey:CategoryID" json:"products"`
	Parent   *Category      `gorm:"foreignKey:ParentId" json:"parent,omitempty"`
	Children []Category     `gorm:"foreignKey:ParentId" json:"children,omitempty"`
}
