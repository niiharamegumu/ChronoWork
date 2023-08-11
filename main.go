package main

import (
	"fmt"
	"os"

	"github.com/niiharamegumu/ChronoWork/db"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

func init() {
	var err error
	dbConn, err = db.ConnectDB()
	if err != nil {
		fmt.Println("database connection error", err)
		os.Exit(1)
	}

}

func main() {
	// seeder
	// if err := db.CreateTestData(dbConn); err != nil {
	// 	fmt.Println("error creating test data", err)
	// 	return
	// }

	fmt.Println("test data created successfully")
}
