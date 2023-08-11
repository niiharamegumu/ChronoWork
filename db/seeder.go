package db

import (
	"time"

	"github.com/niiharamegumu/ChronoWork/models"
	"gorm.io/gorm"
)

func CreateTestData(db *gorm.DB) error {
	projectType := models.ProjectType{
		Name: "Sample Project",
		Tags: []models.Tag{
			{Name: "Tag 1"},
			{Name: "Tag 2"},
		},
	}
	if err := db.Create(&projectType).Error; err != nil {
		return err
	}

	chronoWork := models.ChronoWork{
		Title:         "Sample Task",
		Description:   "Sample Description",
		ProjectTypeID: projectType.ID,
		StartTime:     time.Now(),
		EndTime:       time.Now().Add(time.Hour),
	}
	if err := db.Create(&chronoWork).Error; err != nil {
		return err
	}

	return nil
}
