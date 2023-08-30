package models

import "gorm.io/gorm"

type Setting struct {
	gorm.Model
	RelativeDate uint `gorm:"default:0" json:"relative_date"`
}

func (s *Setting) GetSetting(db *gorm.DB) error {
	if result := db.FirstOrCreate(s); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Setting) UpdateSetting(db *gorm.DB, relativeDate int) error {
	if result := db.Model(s).Update("relative_date", relativeDate); result.Error != nil {
		return result.Error
	}
	return nil
}
