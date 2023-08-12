package db

import (
	"fmt"
	"time"

	"github.com/niiharamegumu/ChronoWork/models"
	"gorm.io/gorm"
)

func CreateTestData(db *gorm.DB) error {
	tags := []models.Tag{
		{Name: "Tag 1"},
		{Name: "Tag 2"},
	}
	if err := db.Create(&tags).Error; err != nil {
		return err
	}

	projectType := models.ProjectType{
		Name: "Sample Project",
		Tags: tags,
	}
	if err := db.Create(&projectType).Error; err != nil {
		return err
	}

	var chronoWorks []models.ChronoWork
	for i := 0; i < 10; i++ {
		chronoWork := models.ChronoWork{
			Title:         fmt.Sprintf("Sample Work %d", i+1),
			ProjectTypeID: projectType.ID,
			StartTime:     time.Now(),
			EndTime:       time.Now().Add(time.Hour),
			TotalSeconds:  6239,
			IsTracking:    false,
		}
		chronoWorks = append(chronoWorks, chronoWork)
	}
	if err := db.Create(&chronoWorks).Error; err != nil {
		return err
	}

	return nil
}
