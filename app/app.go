package app

import (
	"fmt"
	"os"

	"github.com/niiharamegumu/ChronoWork/db"
)

func init() {
	var err error
	err = db.ConnectDB()
	if err != nil {
		fmt.Println("database connection error", err)
		os.Exit(1)
	}
	// ==============================
	// [SEEDER] create test data
	// ==============================
	// if err := db.CreateTestData(db.DB); err != nil {
	// 	fmt.Println("error creating test data", err)
	// }
}

func Execute() {
	if err := InitialSetting(); err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
}
