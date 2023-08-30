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

	mainTitle := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("Today's Work").SetTextColor(tcell.ColorPurple)
	timer := widgets.NewTimer()
	err = timer.CheckActiveTracking(tui)
	if err != nil {
		return err
	}
	work := widgets.NewWork()
	work, err = work.GenerateInitWork(tui)
	if err != nil {
		return err
	}
	form := widgets.NewForm()
	form = form.GenerateInitForm(tui, work)

	// add page
	setting := widgets.NewSetting()
	setting.GenerateInitSetting(tui)
	tui.SetMainPage("setting", setting.Form, false)
	if err = tui.SetWidget("settingForm", setting.Form); err != nil {
		return err
	}

	menu := widgets.NewMenu()
	menu = menu.GenerateInitMenu(tui, work, setting)

	tui.SetHeader(header, false)
	tui.SetMenu(menu.List, false)
	tui.SetWork(mainTitle, form.Form, timer.Wrapper, work.Table, true) // default focus
	work.TableCapture(tui, form, timer)
	form.FormCapture(tui)

	tui.GlobalKeyActions()
	if err = tui.App.SetRoot(tui.Grid, true).EnableMouse(true).Run(); err != nil {
		return err
	}
	return nil
}
