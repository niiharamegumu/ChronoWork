package widgets

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/niiharamegumu/chronowork/db"
	"github.com/niiharamegumu/chronowork/models"
	"github.com/niiharamegumu/chronowork/service"
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
	s.Form.AddInputField("Show Relative Date(0:Today Only) : ", fmt.Sprint(setting.RelativeDate), 20, nil, nil).
		AddInputField("Person Day : ", fmt.Sprint(setting.PersonDay), 20, nil, nil).
		AddCheckbox("Display As Person Day : ", setting.DisplayAsPersonDay, nil).
		AddButton("Save", func() {
			s.update()
			s.ReStore(tui)
			tui.SetFocus("menu")
		}).
		AddButton("Cancel", func() {
			s.ReStore(tui)
			tui.SetFocus("menu")
		})
}

func (s *Setting) ReStore(tui *service.TUI) {
	s.Form.Clear(true)
	s.GenerateInitSetting(tui)
}

func (s *Setting) update() {
	relativeDate := s.Form.GetFormItemByLabel("Show Relative Date(0:Today Only) : ").(*tview.InputField).GetText()
	personDay := s.Form.GetFormItemByLabel("Person Day : ").(*tview.InputField).GetText()
	displayAsPersonDay := s.Form.GetFormItemByLabel("Display As Person Day : ").(*tview.Checkbox).IsChecked()

	var dateInt, personDayInt int
	var err error
	if dateInt, err = strconv.Atoi(relativeDate); err != nil {
		log.Println(err)
		return
	}
	if personDayInt, err = strconv.Atoi(personDay); err != nil {
		log.Println(err)
		return
	}

	var setting models.Setting
	if err = setting.GetSetting(db.DB); err != nil {
		log.Println(err)
		return
	}
	new := models.Setting{
		RelativeDate:       uint(dateInt),
		PersonDay:          uint(personDayInt),
		DisplayAsPersonDay: displayAsPersonDay,
	}
	if err = setting.UpdateSetting(db.DB, new); err != nil {
		log.Println(err)
		return
	}
}
