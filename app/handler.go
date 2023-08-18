package app

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/niiharamegumu/ChronoWork/service"
	"github.com/niiharamegumu/ChronoWork/widgets"
	"github.com/rivo/tview"
)

var (
	headerText = fmt.Sprintf("%s - %s", " ChronoWork ", time.Now().Format("2006-01-02"))
	tui        = service.NewTUI()
)

func InitialSetting() error {
	var err error

	header := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(headerText)

	menu := widgets.NewMenu()
	menu = menu.GenerateInitMenu(tui)

	mainTitle := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("Today's Work").SetTextColor(tcell.ColorPurple)
	timer := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Timer")
	work := widgets.NewWork()
	work, err = work.GenerateInitWork(tui)
	if err != nil {
		return err
	}

	tui.SetHeader(header, false)
	tui.SetMenu(menu.List, false)
	tui.SetMain(mainTitle, work.Form, timer, work.Table, true) // default focus

	work.TableCapture(tui)
	work.FormCapture(tui)

	tui.GlobalKeyActions()
	if err = tui.App.SetRoot(tui.Grid, true).EnableMouse(true).Run(); err != nil {
		return err
	}
	return nil
}
