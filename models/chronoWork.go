package models

import (
	"time"

	"gorm.io/gorm"
)

type ChronoWork struct {
	gorm.Model
	Title         string    `gorm:"size:255; required" json:"title"`
	ProjectTypeID uint      `json:"project_type_id"`
	TagID         uint      `json:"tag_id"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	IsTracking    bool      `json:"is_tracking"`
	TotalSeconds  int       `json:"total_seconds"`

	ProjectType ProjectType `gorm:"foreignkey:ProjectTypeID"`
	Tag         Tag         `gorm:"foreignkey:TagID"`
}

func (c *ChronoWork) FindInRangeByTime(db *gorm.DB, startTime, endTime time.Time) ([]ChronoWork, error) {
	var chronoWorks []ChronoWork
	result := db.
		Preload("ProjectType").
		Preload("Tag").
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

func CreateChronoWork(db *gorm.DB, title string, projectTypeID, tagID uint) error {
	chronoWork := ChronoWork{
		Title:         title,
		ProjectTypeID: projectTypeID,
		TagID:         tagID,
		StartTime:     time.Time{},
		EndTime:       time.Time{},
		IsTracking:    false,
		TotalSeconds:  0,
	}
	result := db.Create(&chronoWork)
	if result.Error != nil {
		return result.Error
	}
	return nil

}
