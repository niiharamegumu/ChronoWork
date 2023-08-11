package models

import (
	"time"

	"gorm.io/gorm"
)

type ChronoWork struct {
	gorm.Model
	Title         string    `gorm:"size:255" json:"title"`
	Description   string    `gorm:"size:255" json:"description"`
	ProjectTypeID uint      `json:"project_type_id"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`

	ProjectType ProjectType `gorm:"foreignkey:ProjectTypeID"`
}
