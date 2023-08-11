package main

import (
	"github.com/niiharamegumu/ChronoWork/app"
)

func main() {
	// seeder
	// if err := db.CreateTestData(dbConn); err != nil {
	// 	fmt.Println("error creating test data", err)
	// 	return
	// }
	app.Execute()
}
