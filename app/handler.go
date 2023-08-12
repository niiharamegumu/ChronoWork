package app

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/niiharamegumu/ChronoWork/pkg"
	"github.com/niiharamegumu/ChronoWork/service"
	"github.com/niiharamegumu/ChronoWork/widgets"
	"github.com/rivo/tview"
)

var (
	headerText = fmt.Sprintf("%s - %s", " ChronoWork ", time.Now().Format("2006-01-02"))
	tui        = service.NewTUI()
)

func InitialSetting() error {
	header := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(headerText)
	timer := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Timer")

	menu := widgets.NewMenu()
	menuW := menu.GenerateInitMenu(tui)

	mainTitle := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("Today's Work").SetTextColor(tcell.ColorPurple)
	work := widgets.NewWork()
	workW := work.GenerateInitWork(pkg.TodayStartTime(), pkg.TodayEndTime(), tui)

	tui.SetHeader(header, false)
	tui.SetMenu(menuW.List, false)
	tui.SetMain(mainTitle, workW.Form, timer, workW.Table, true) // default focus

	tui.GlobalKeyActions()
	if err := tui.App.SetRoot(tui.Grid, true).EnableMouse(true).Run(); err != nil {
		return err
	}
	return nil
}
