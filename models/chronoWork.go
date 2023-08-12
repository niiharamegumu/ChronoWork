package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type ChronoWork struct {
	gorm.Model
	Title         string    `gorm:"size:255; required" json:"title"`
	ProjectTypeID uint      `json:"project_type_id"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	IsTracking    bool      `json:"is_tracking"`
	TotalSeconds  int       `json:"total_seconds"`

	ProjectType ProjectType `gorm:"foreignkey:ProjectTypeID"`
}

func (c *ChronoWork) GetCombinedTagNames() string {
	tagNames := []string{}
	for _, tag := range c.ProjectType.Tags {
		tagNames = append(tagNames, tag.Name)
	}
	return strings.Join(tagNames, " | ")
}
