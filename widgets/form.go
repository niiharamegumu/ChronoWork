package widgets

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/niiharamegumu/ChronoWork/db"
	"github.com/niiharamegumu/ChronoWork/models"
	"github.com/niiharamegumu/ChronoWork/service"
	"github.com/rivo/tview"
)

var notSelectText = "Not Select"

type Form struct {
	Form *tview.Form
}

func NewForm() *Form {
	form := &Form{
		Form: tview.NewForm().
			SetButtonBackgroundColor(tcell.ColorPurple).
			SetLabelColor(tcell.ColorPurple),
	}
	return form
}

func (f *Form) GenerateInitForm(tui *service.TUI, work *Work) *Form {
	f.ConfigureForm(tui, work)
	return f
}

func (f *Form) FormCapture(tui *service.TUI) {
	f.Form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlB:
			tui.SetFocus("mainContent")
		}
		return event
	})
}

func (f *Form) ResetForm() {
	f.Form.GetFormItemByLabel("Title").(*tview.InputField).SetText("")
	f.Form.GetFormItemByLabel("Project").(*tview.DropDown).SetCurrentOption(0)
	f.Form.GetFormItemByLabel("Tags").(*tview.DropDown).SetOptions([]string{notSelectText}, nil).SetCurrentOption(0)
}

func (f *Form) ConfigureForm(tui *service.TUI, work *Work) {
	tagsOptions := append([]string{notSelectText})
	f.Form.
		AddInputField("Title", "", 50, nil, nil).
		AddDropDown("Project", append([]string{notSelectText}, models.AllProjectTypeNames(db.DB)...), 0, f.projectDropDownChanged).
		AddDropDown("Tags", tagsOptions, 0, nil).
		AddButton("Save", func() {
			if err := f.store(); err != nil {
				log.Println(err)
				return
			}
			if err := work.ReStoreTable(); err != nil {
				log.Println(err)
				return
			}
			tui.SetFocus("mainContent")
		}).
		AddButton("Cancel", func() {
			tui.SetFocus("mainContent")
		})
}

func (f *Form) projectDropDownChanged(option string, optionIndex int) {
	var tagsOptions []string
	if f.Form.GetFormItemByLabel("Tags") == nil {
		return
	}
	var projectType models.ProjectType
	result := db.DB.Preload("Tags").Where("name = ?", option).Find(&projectType)
	if result.Error != nil {
		tagsOptions = []string{notSelectText}
	} else {
		tagsOptions = append([]string{notSelectText}, projectType.GetTagNames()...)
	}
	f.Form.GetFormItemByLabel("Tags").(*tview.DropDown).
		SetOptions(tagsOptions, nil).
		SetCurrentOption(0)
}

func (f *Form) store() error {
	title := f.Form.GetFormItemByLabel("Title").(*tview.InputField).GetText()
	_, projectVal := f.Form.GetFormItemByLabel("Project").(*tview.DropDown).GetCurrentOption()
	_, tagVal := f.Form.GetFormItemByLabel("Tags").(*tview.DropDown).GetCurrentOption()

	if title == "" {
		return nil
	}

	var projectTypeID uint
	var tagID uint
	if projectVal != notSelectText {
		projectType, err := models.FindProjectTypeByName(db.DB, projectVal)
		if err != nil {
			log.Println(err)
			return err
		}
		projectTypeID = projectType.ID
		if tagVal != notSelectText {
			for _, tag := range projectType.Tags {
				if tag.Name == tagVal {
					tagID = tag.ID
				}
			}
		}
	}

	if err := models.CreateChronoWork(db.DB, title, projectTypeID, tagID); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
