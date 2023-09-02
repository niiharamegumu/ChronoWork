package db

import (
	"fmt"
	"log"
	"os"

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

	var dbPath string
	rootPath := os.Getenv("CHRONOWORK_ROOT_PATH")
	if rootPath != "" {
		dbPath = fmt.Sprintf("%s/%s", os.Getenv("CHRONOWORK_ROOT_PATH"), "sqlite.db")
	} else {
		log.Println("CHRONOWORK_ROOT_PATH is not set")
		dbPath = fmt.Sprintf("%s/%s", ".", "sqlite.db")
	}

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	// auto migration for models
	DB.AutoMigrate(
		&models.ChronoWork{},
		&models.ProjectType{},
		&models.Tag{},
		&models.Setting{},
	)

	return nil
}

func CloseDB() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
