package models

import (
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name string `gorm:"size:255; required" json:"name"`
}

func AllTagsNmaes(db *gorm.DB) []string {
	var tags []Tag
	db.Find(&tags)
	var tagNames []string
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}
	return tagNames
}
