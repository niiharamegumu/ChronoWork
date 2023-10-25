package models

import (
	"gorm.io/gorm"
)

type Setting struct {
	gorm.Model
	RelativeDate       uint   `gorm:"default:0" json:"relative_date"`
	PersonDay          uint   `gorm:"default:8" json:"person_day"`
	DisplayAsPersonDay bool   `gorm:"default:1" json:"display_as_person_day"`
	DownloadPath       string `gorm:"default:./" json:"download_path"`
}

func (s *Setting) GetSetting(db *gorm.DB) error {
	if result := db.FirstOrCreate(s); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Setting) UpdateSetting(db *gorm.DB, setting Setting) error {
	dataMap := map[string]any{
		"relative_date":         setting.RelativeDate,
		"person_day":            setting.PersonDay,
		"display_as_person_day": setting.DisplayAsPersonDay,
		"download_path":         setting.DownloadPath,
	}
	if result := db.Model(s).Select(
		"relative_date",
		"person_day",
		"display_as_person_day",
		"download_path").
		Updates(dataMap); result.Error != nil {
		return result.Error
	}
	return nil
}
