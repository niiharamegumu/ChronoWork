package service

import (
	"errors"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TUI struct {
	App      *tview.Application
	Grid     *tview.Grid
	MainPage *tview.Pages
	Widgets  map[string]tview.Primitive
}

func (t *TUI) SetHeader(header tview.Primitive, focus bool) {
	t.Grid.AddItem(header, 0, 0, 1, 3, 0, 0, focus)
	t.Widgets["header"] = header
}

func (t *TUI) SetMenu(menu tview.Primitive, focus bool) {
	t.Grid.AddItem(menu, 1, 0, 1, 1, 0, 0, focus)
	t.Widgets["menu"] = menu
}

func (t *TUI) SetWork(mainTitle, mainForm, mainTimer, mainContent tview.Primitive, focus bool) {
	main := tview.NewGrid().SetRows(1, 10, 0).SetColumns(0, 0).SetBorders(true)
	main.AddItem(mainTitle, 0, 0, 1, 2, 0, 0, false)
	main.AddItem(mainForm, 1, 0, 1, 1, 0, 0, false)
	main.AddItem(mainTimer, 1, 1, 1, 1, 0, 0, false)
	main.AddItem(mainContent, 2, 0, 1, 2, 0, 0, true)
	t.Widgets["mainTitle"] = mainTitle
	t.Widgets["mainTimer"] = mainTimer
	t.Widgets["mainWorkForm"] = mainForm
	t.Widgets["mainWorkContent"] = mainContent

	t.SetMainPage("work", main, true)
	t.Grid.AddItem(t.MainPage, 1, 1, 1, 2, 0, 100, focus)
}

func (t *TUI) SetMainPage(name string, page tview.Primitive, focus bool) {
	t.MainPage.AddPage(name, page, true, focus)
}

func (t *TUI) SetWidget(name string, widget tview.Primitive) error {
	if _, ok := t.Widgets[name]; ok {
		return errors.New("widget already exists")
	}
	t.Widgets[name] = widget
	return nil
}

func (t *TUI) SetModal(modal tview.Primitive) {
	t.App.SetRoot(modal, false)
	t.Widgets["modal"] = modal
}

func (t *TUI) DeleteModal() {
	t.App.SetRoot(t.Grid, true)
	delete(t.Widgets, "modal")
}

func NewTUI() *TUI {
	return &TUI{
		App: tview.NewApplication(),
		Grid: tview.NewGrid().
			SetRows(1, 0).
			SetColumns(15, 0).
			SetBorders(true),
		MainPage: tview.NewPages(),
		Widgets:  make(map[string]tview.Primitive),
	}
}

func (t *TUI) GlobalKeyActions() {
	t.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			t.SetFocus("menu")
		}
		return event
	})
}

func (t *TUI) Quit() {
	t.App.Stop()
}

func (t *TUI) SetFocus(name string) {
	t.App.SetFocus(t.Widgets[name])
}

func (t *TUI) ChangeToPage(name string) {
	t.MainPage.SwitchToPage(name)
}
