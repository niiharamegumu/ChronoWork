package models

import (
	"gorm.io/gorm"
)

type ProjectType struct {
	gorm.Model
	Name string `gorm:"size:255; unique; not null" json:"name"`
	Tags []Tag  `gorm:"many2many:project_type_tags;" json:"tags"`
}

func AllProjectTypeNames(db *gorm.DB) []string {
	var projectTypes []ProjectType
	db.Find(&projectTypes)
	var projectTypeNames []string
	for _, projectType := range projectTypes {
		projectTypeNames = append(projectTypeNames, projectType.Name)
	}
	return projectTypeNames
}

func AllProjectTypeWithTags(db *gorm.DB) []ProjectType {
	var projectTypes []ProjectType
	db.Preload("Tags").Find(&projectTypes)
	return projectTypes
}

func (p *ProjectType) GetTagNames() []string {
	var tagNames []string
	for _, tag := range p.Tags {
		tagNames = append(tagNames, tag.Name)
	}
	return tagNames
}

func FindProjectTypeByName(db *gorm.DB, name string) (*ProjectType, error) {
	var projectType ProjectType
	result := db.
		Preload("Tags").
		Where("name = ?", name).Find(&projectType)
	if result.Error != nil {
		return nil, result.Error
	}
	return &projectType, nil
}

func FindProjectTypeByID(db *gorm.DB, id uint) (*ProjectType, error) {
	var projectType ProjectType
	result := db.
		Preload("Tags").
		First(&projectType, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &projectType, nil
}

func DeleteProjectType(db *gorm.DB, id uint) error {
	result := db.Unscoped().Delete(&ProjectType{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
