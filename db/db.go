package db

import (
	"fmt"

	"github.com/niiharamegumu/ChronoWork/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func ConnectDB() error {
	if DB != nil {
		return nil
	}

	// TODO: Retrieve the database path from an environment variable
	dbPath := fmt.Sprintf("%s/%s", "./", "sqlite.db")
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	// auto migration for models
	DB.AutoMigrate(
		&models.ChronoWork{},
		&models.ProjectType{},
		&models.Tag{},
	)

	return nil
}
