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

func (p *ProjectType) GetTagNames() []string {
	var tagNames []string
	for _, tag := range p.Tags {
		tagNames = append(tagNames, tag.Name)
	}
	return tagNames
}
