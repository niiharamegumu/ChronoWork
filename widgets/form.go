package widgets

import (
	"fmt"
	"log"
	"strconv"
	
	"github.com/niiharamegumu/chronowork/db"
	"github.com/niiharamegumu/chronowork/models"
	"github.com/niiharamegumu/chronowork/service"
	"github.com/niiharamegumu/chronowork/util/timeutil"
	"github.com/gdamore/tcell/v2"
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
			SetLabelColor(tcell.ColorPurple).
			SetFieldTextColor(tcell.ColorGray).
			SetFieldBackgroundColor(tcell.ColorWhite),
	}
	return form
}

func (f *Form) GenerateInitForm(tui *service.TUI, work *Work) *Form {
	f.ConfigureStoreForm(tui, work)
	return f
}

func (f *Form) FormCapture(tui *service.TUI) {
	f.Form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlB:
			tui.SetFocus("mainWorkContent")
		}
		return event
	})
}

func (f *Form) ResetForm() {
	f.Form.GetFormItemByLabel("Title").(*tview.InputField).SetText("")
	f.Form.GetFormItemByLabel("Project").(*tview.DropDown).SetCurrentOption(0)
	f.Form.GetFormItemByLabel("Tags").(*tview.DropDown).SetOptions([]string{notSelectText}, nil).SetCurrentOption(0)
}

func (f *Form) ConfigureStoreForm(tui *service.TUI, work *Work) {
	f.Form.
		AddInputField("Title", "", 50, nil, nil).
		AddDropDown("Project", append([]string{notSelectText}, models.AllProjectTypeNames(db.DB)...), 0, f.projectDropDownChanged).
		AddDropDown("Tags", []string{notSelectText}, 0, nil).
		AddButton("Store", func() {
			if err := f.store(); err != nil {
				log.Println(err)
				return
			}
			if err := work.ReStoreTable(timeutil.RelativeStartTime(), timeutil.TodayEndTime()); err != nil {
				log.Println(err)
				return
			}
			tui.SetFocus("mainWorkContent")
		}).
		AddButton("Cancel", func() {
			tui.SetFocus("mainWorkContent")
		})
}

func (f *Form) configureUpdateForm(tui *service.TUI, work *Work, chronoWork *models.ChronoWork) {
	projectOptions := append([]string{notSelectText}, models.AllProjectTypeNames(db.DB)...)
	tagsOptions := []string{notSelectText}
	f.Form.AddInputField("Title", chronoWork.Title, 50, nil, nil).
		AddDropDown("Project", projectOptions, 0, f.projectDropDownChanged).
		AddDropDown("Tags", tagsOptions, 0, nil)

	var projectType models.ProjectType
	if chronoWork.ProjectType.Name != "" {
		result := db.DB.Preload("Tags").Where("name = ?", chronoWork.ProjectType.Name).Find(&projectType)
		if result.Error != nil {
			log.Println(result.Error)
			return
		}
		for i, projectOption := range projectOptions {
			if projectOption == chronoWork.ProjectType.Name {
				f.Form.GetFormItemByLabel("Project").(*tview.DropDown).SetCurrentOption(i)
				break
			}
		}

		tagsOptions = append(tagsOptions, projectType.GetTagNames()...)
		f.Form.GetFormItemByLabel("Tags").(*tview.DropDown).SetOptions(tagsOptions, nil)
		if chronoWork.Tag.Name != "" {
			for i, tagOption := range tagsOptions {
				if tagOption == chronoWork.Tag.Name {
					f.Form.GetFormItemByLabel("Tags").(*tview.DropDown).SetCurrentOption(i)
					break
				}
			}
		}
	}

	f.Form.AddButton("Update", func() {
		if err := f.update(chronoWork); err != nil {
			log.Println(err)
			return
		}
		if err := work.ReStoreTable(timeutil.RelativeStartTime(), timeutil.TodayEndTime()); err != nil {
			log.Println(err)
			return
		}
		tui.SetFocus("mainWorkContent")
	}).
		AddButton("Cancel", func() {
			tui.SetFocus("mainWorkContent")
		})
}

func (f *Form) configureTimerForm(tui *service.TUI, work *Work, chronoWork *models.ChronoWork) {
	hour := chronoWork.TotalSeconds / 3600
	minute := (chronoWork.TotalSeconds - hour*3600) / 60
	second := chronoWork.TotalSeconds - hour*3600 - minute*60

	f.Form.AddInputField("Hour(0-)", fmt.Sprint(hour), 20, nil, nil).
		AddInputField("Minute(0-59)", fmt.Sprint(minute), 20, nil, nil).
		AddInputField("Second(0-59)", fmt.Sprint(second), 20, nil, nil).
		AddButton("Reset", func() {
			if err := f.resetTimer(chronoWork); err != nil {
				log.Println(err)
				return
			}
			if err := work.ReStoreTable(timeutil.RelativeStartTime(), timeutil.TodayEndTime()); err != nil {
				log.Println(err)
				return
			}
			tui.SetFocus("mainWorkContent")
		}).
		AddButton("Cancel", func() {
			tui.SetFocus("mainWorkContent")
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

	if _, err := models.CreateChronoWork(db.DB, title, projectTypeID, tagID); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (f *Form) update(chronoWork *models.ChronoWork) error {
	title := f.Form.GetFormItemByLabel("Title").(*tview.InputField).GetText()
	_, projectVal := f.Form.GetFormItemByLabel("Project").(*tview.DropDown).GetCurrentOption()
	_, tagVal := f.Form.GetFormItemByLabel("Tags").(*tview.DropDown).GetCurrentOption()

	if title == "" {
		return nil
	}
	var projectTypeID uint = 0
	var tagID uint = 0
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
	if err := chronoWork.UpdateChronoWork(db.DB, title, projectTypeID, tagID); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (f *Form) resetTimer(chronoWork *models.ChronoWork) error {
	hour := f.Form.GetFormItemByLabel("Hour(0-)").(*tview.InputField).GetText()
	minute := f.Form.GetFormItemByLabel("Minute(0-59)").(*tview.InputField).GetText()
	second := f.Form.GetFormItemByLabel("Second(0-59)").(*tview.InputField).GetText()

	var hourInt, minuteInt, secondInt uint64
	var err error
	if hourInt, err = strconv.ParseUint(hour, 10, 16); err != nil {
		return err
	}
	if minuteInt, err = strconv.ParseUint(minute, 10, 16); err != nil {
		return err
	}
	if secondInt, err = strconv.ParseUint(second, 10, 16); err != nil {
		return err
	}
	totalSeconds := int(hourInt*3600 + minuteInt*60 + secondInt)
	err = chronoWork.UpdateChronoWorkTotalSeconds(db.DB, totalSeconds)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
