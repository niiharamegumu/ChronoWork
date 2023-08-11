package app

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/niiharamegumu/ChronoWork/views/layouts"
	"github.com/niiharamegumu/ChronoWork/views/widgets"
	"github.com/rivo/tview"
)

var appName string
var app *tview.Application

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("error loading .env file: %v", err)
	}
	appName = os.Getenv("APP_NAME")
	app = tview.NewApplication()
}

func InitShow() error {
	headerText := fmt.Sprintf("%s - %s", appName, time.Now().Format("2006-01-02"))
	header := widgets.SimpleText(headerText)

	flexRow := layouts.FlexRow()
	flexRow.AddItem(header, 3, 1, false)

	if err := app.SetRoot(flexRow, true).Run(); err != nil {
		return err
	}
	return nil
}
