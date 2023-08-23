package models

import (
	"math"
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

func (c *ChronoWork) StartTrackingChronoWork(db *gorm.DB) error {
	result := db.Model(c).Updates(map[string]interface{}{
		"start_time":  time.Now(),
		"end_time":    time.Time{},
		"is_tracking": true,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (c *ChronoWork) StopTrackingChronoWorks(db *gorm.DB) error {
	result := db.Model(c).Where("is_tracking = ?", true).Updates(map[string]interface{}{
		"end_time":    time.Now(),
		"is_tracking": false,
	})
	if result.Error != nil {
		return result.Error
	}
	result = db.Model(c).Where("is_tracking = ?", false).Updates(map[string]interface{}{
		"total_seconds": gorm.Expr("total_seconds + ?", int(math.Floor(c.EndTime.Sub(c.StartTime).Seconds()))),
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (c *ChronoWork) FindInRangeByTime(db *gorm.DB, startTime, endTime time.Time) ([]ChronoWork, error) {
	var chronoWorks []ChronoWork
	result := db.
		Preload("ProjectType").
		Preload("Tag").
		Order("created_at desc").
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

func FindChronoWork(db *gorm.DB, id uint) (ChronoWork, error) {
	var chronoWork ChronoWork
	result := db.Preload("ProjectType").Preload("Tag").First(&chronoWork, id)
	if result.Error != nil {
		return chronoWork, result.Error
	}
	return chronoWork, nil
}

func FindTrackingChronoWorks(db *gorm.DB) ([]ChronoWork, error) {
	var chronoWorks []ChronoWork
	result := db.Find(&chronoWorks, "is_tracking = ?", true)
	if result.Error != nil {
		return nil, result.Error
	}
	return chronoWorks, nil
}

func DeleteChronoWork(db *gorm.DB, id uint) error {
	result := db.Unscoped().Delete(&ChronoWork{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
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
