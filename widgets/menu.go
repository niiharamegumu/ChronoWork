package widgets

import (
	"github.com/niiharamegumu/ChronoWork/pkg"
	"github.com/niiharamegumu/ChronoWork/service"
	"github.com/rivo/tview"
)

type Menu struct {
	List *tview.List
}

func NewMenu() *Menu {
	return &Menu{
		List: tview.NewList(),
	}
}

func (m *Menu) AddListItem(text string, shortcut rune, selected func()) *Menu {
	m.List.AddItem(text, "", shortcut, selected)
	return m
}

func (m *Menu) GenerateInitMenu(tui *service.TUI, work *Work, setting *Setting, project *Project) *Menu {
	m.AddListItem("Works", 'w', func() {
		work.ReStoreTable(pkg.RelativeStartTime(), pkg.TodayEndTime())
		tui.ChangeToPage("work")
		tui.SetFocus("mainWorkContent")
	})
	m.AddListItem("Projects", 'p', func() {
		tui.ChangeToPage("project")
		tui.SetFocus("projectTable")
	})
	m.AddListItem("Setting", 's', func() {
		setting.ReStore(tui)
		tui.ChangeToPage("setting")
		tui.SetFocus("settingForm")
	})
	m.AddListItem("Quit", 'q', tui.Quit)
	return m
}
