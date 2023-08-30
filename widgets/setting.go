package widgets

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/niiharamegumu/ChronoWork/db"
	"github.com/niiharamegumu/ChronoWork/models"
	"github.com/niiharamegumu/ChronoWork/service"
	"github.com/rivo/tview"
)

type Setting struct {
	Form *tview.Form
}

func NewSetting() *Setting {
	return &Setting{
		Form: tview.NewForm().
			SetLabelColor(tcell.ColorPurple),
	}
}

func (s *Setting) GenerateInitSetting(tui *service.TUI) {
	var setting models.Setting
	if err := setting.GetSetting(db.DB); err != nil {
		log.Println(err)
		return
	}
	s.Form.AddInputField("Show Relative Date(0:Today Only) : ", fmt.Sprintln(setting.RelativeDate), 20, nil, nil).
		AddButton("Save", func() {
			s.Update()
			s.ReStore(tui)
			tui.SetFocus("menu")
		}).
		AddButton("Cancel", func() {
			tui.ChangeToPage("work")
			tui.SetFocus("menu")
		})
}

func (s *Setting) ReStore(tui *service.TUI) {
	s.Form.Clear(true)
	s.GenerateInitSetting(tui)
}

func (s *Setting) Update() {
	relativeDate := s.Form.GetFormItemByLabel("Show Relative Date(0:Today Only) : ").(*tview.InputField).GetText()

	var dateInt uint64
	var err error
	if dateInt, err = strconv.ParseUint(relativeDate, 10, 64); err != nil {
		log.Println(err)
		return
	}
	var setting models.Setting
	if err := setting.GetSetting(db.DB); err != nil {
		log.Println(err)
		return
	}
	setting.UpdateSetting(db.DB, int(dateInt))
}
