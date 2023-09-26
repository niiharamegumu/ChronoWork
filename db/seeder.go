package db

import (
	"chronowork/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func CreateTestData(db *gorm.DB) error {
	tags := []models.Tag{
		{Name: "Tag 1"},
		{Name: "Tag 2"},
		{Name: "Tag 3"},
		{Name: "Tag 4"},
	}
	if err := db.Create(&tags).Error; err != nil {
		return err
	}

	projectType1 := models.ProjectType{
		Name: "Sample Project01",
		Tags: []models.Tag{tags[0], tags[1]},
	}
	if err := db.Create(&projectType1).Error; err != nil {
		return err
	}

	projectType2 := models.ProjectType{
		Name: "Sample Project02",
		Tags: []models.Tag{tags[2], tags[3]},
	}
	if err := db.Create(&projectType2).Error; err != nil {
		return err
	}

	var chronoWorks []models.ChronoWork
	for i := 0; i < 2; i++ {
		chronoWork := models.ChronoWork{
			Title:         fmt.Sprintf("Sample %d", i+1),
			ProjectTypeID: projectType1.ID,
			TagID:         tags[0].ID,
			StartTime:     time.Time{},
			EndTime:       time.Time{},
			TotalSeconds:  3600,
			IsTracking:    false,
		}
		chronoWorks = append(chronoWorks, chronoWork)
	}
	if err := db.Create(&chronoWorks).Error; err != nil {
		return err
	}

	chronoWorks = []models.ChronoWork{}
	for i := 0; i < 2; i++ {
		chronoWork := models.ChronoWork{
			Title:         fmt.Sprintf("Work %d", i+1),
			ProjectTypeID: projectType2.ID,
			TagID:         tags[2].ID,
			StartTime:     time.Time{},
			EndTime:       time.Time{},
			TotalSeconds:  3600,
			IsTracking:    false,
		}
		chronoWorks = append(chronoWorks, chronoWork)
	}
	if err := db.Create(&chronoWorks).Error; err != nil {
		return err
	}

	return nil
}
