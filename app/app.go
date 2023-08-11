package app

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

func Execute() {
	if err := InitShow(); err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
}
