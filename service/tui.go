package service

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TUI struct {
	App     *tview.Application
	Grid    *tview.Grid
	Widgets map[string]tview.Primitive
}

func (t *TUI) SetHeader(header tview.Primitive, focus bool) {
	t.Grid.AddItem(header, 0, 0, 1, 3, 0, 0, focus)
	t.Widgets["header"] = header
}

func (t *TUI) SetMenu(menu tview.Primitive, focus bool) {
	t.Grid.AddItem(menu, 1, 0, 1, 1, 0, 50, focus)
	t.Widgets["menu"] = menu
}

func (t *TUI) SetMain(mainTitle, mainForm, mainTimer, mainContent tview.Primitive, focus bool) {
	main := tview.NewGrid().SetRows(1, 10, 0).SetColumns(0, 0).SetBorders(true)
	main.AddItem(mainTitle, 0, 0, 1, 2, 0, 0, false)
	main.AddItem(mainForm, 1, 0, 1, 1, 0, 0, false)
	main.AddItem(mainTimer, 1, 1, 1, 1, 0, 0, false)
	main.AddItem(mainContent, 2, 0, 1, 2, 0, 0, true)

	t.Grid.AddItem(main, 1, 1, 1, 2, 0, 100, focus)
	t.Widgets["main"] = main
	t.Widgets["mainTitle"] = mainTitle
	t.Widgets["mainForm"] = mainForm
	t.Widgets["mainContent"] = mainContent
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
			SetColumns(30, 0).
			SetBorders(true),
		Widgets: make(map[string]tview.Primitive),
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
