package app

import (
	"fmt"
	"log"
	"os"

	"github.com/niiharamegumu/ChronoWork/db"
)

func init() {
	err := db.ConnectDB()
	if err != nil {
		fmt.Println("database connection error", err)
		os.Exit(1)
	}
	log.Println("database connection success")
}

func Execute() {
	defer func() {
		if err := db.CloseDB(); err != nil {
			log.Println("error closing database", err)
		}
		log.Println("database connection closed")
	}()

	// ==============================
	// [SEEDER] create test data
	// ==============================
	// if err := db.CreateTestData(db.DB); err != nil {
	// 	log.Println("error creating test data", err)
	// }

	if err := InitialSetting(); err != nil {
		log.Println("error", err)
		os.Exit(1)
	}
}
