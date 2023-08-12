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

func (c *ChronoWork) FindInRangeByTime(db *gorm.DB, startTime, endTime time.Time) ([]ChronoWork, error) {
	var chronoWorks []ChronoWork
	result := db.
		Preload("ProjectType").
		Preload("ProjectType.Tags").
		Order("updated_at desc").
		Order("id desc").
		Find(
			&chronoWorks,
			"created_at >= ? AND created_at <= ?",
			startTime,
			endTime,
		)
	if result.Error != nil {
		return nil, result.Error
	}
	return chronoWorks, nil
}
