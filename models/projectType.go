package models

import (
	"gorm.io/gorm"
)

type ProjectType struct {
	gorm.Model
	Name string `gorm:"size:255" json:"name"`
	Tags []Tag  `gorm:"many2many:project_type_tags;" json:"tags"`
}
