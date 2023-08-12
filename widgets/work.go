package widgets

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/niiharamegumu/ChronoWork/db"
	"github.com/niiharamegumu/ChronoWork/models"
	"github.com/niiharamegumu/ChronoWork/pkg"
	"github.com/niiharamegumu/ChronoWork/service"
	"github.com/rivo/tview"
	"gorm.io/gorm"
)

var (
	workHeader = []string{
		"ID",
		"TotalTime",
		"Title",
		"Project",
		"Tags",
	}
	notSelectText = "Not Select"
)

type Work struct {
	Table *tview.Table
	Form  *tview.Form
}

func NewWork() *Work {
	work := &Work{
		Table: tview.NewTable().
			SetBorders(true).
			SetSelectable(true, false).
			SetFixed(1, 1).
			SetBordersColor(tview.Styles.BorderColor),
		Form: tview.NewForm().
			SetButtonBackgroundColor(tcell.ColorPurple).
			SetLabelColor(tcell.ColorPurple),
	}
	for i, header := range workHeader {
		work.Table.
			SetCell(0, i,
				tview.
					NewTableCell(header).
					SetAlign(tview.AlignCenter).
					SetTextColor(tcell.ColorPurple).
					SetSelectable(false).
					SetExpansion(1),
			)
	}
	return work
}

func (w *Work) goToTop() {
	w.Table.ScrollToBeginning().Select(1, 0)
}

func (w *Work) goToBottom() {
	w.Table.ScrollToEnd().Select(w.Table.GetRowCount()-1, 0)
}

func (w *Work) tableCapture(tui *service.TUI) {
	w.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 't':
				w.goToTop()
			case 'b':
				w.goToBottom()
			case 'a':
				tui.SetFocus("mainForm")
			}
		}
		return event
	})
}

func (w *Work) formCapture(tui *service.TUI) {
	w.Form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			w.resetForm()
			tui.SetFocus("mainContent")
		}
		return event
	})
}

func (w *Work) resetForm() {
	w.Form.
		GetFormItemByLabel("Title").(*tview.InputField).
		SetText("")
	w.Form.
		GetFormItemByLabel("Project").(*tview.DropDown).
		SetCurrentOption(0)
	w.Form.
		GetFormItemByLabel("Tags").(*tview.DropDown).
		SetOptions([]string{notSelectText}, nil).
		SetCurrentOption(0)
}

func GenerateInitWork(startTime, endTime time.Time, tui *service.TUI) *Work {
	work := NewWork()
	var chronoWorks []models.ChronoWork
	var result *gorm.DB

	result = db.DB.
		Preload("ProjectType").
		Preload("ProjectType.Tags").
		Order("updated_at desc").
		Order("id desc").
		Find(
			&chronoWorks,
			"created_at >= ? AND created_at <= ?",
			startTime,
			endTime,
		)
	if result.Error != nil {
		fmt.Println(result.Error)
		return work
	}
	for i, chronoWork := range chronoWorks {
		// ID
		work.Table.SetCell(i+1, 0,
			tview.
				NewTableCell(fmt.Sprintf("%d", chronoWork.ID)).
				SetAlign(tview.AlignCenter).
				SetExpansion(1))
		// TotalTime
		work.Table.SetCell(i+1, 1,
			tview.
				NewTableCell(pkg.FormatTime(chronoWork.TotalSeconds)).
				SetAlign(tview.AlignCenter).
				SetExpansion(1))
		// Title
		work.Table.SetCell(i+1, 2,
			tview.
				NewTableCell(chronoWork.Title).
				SetAlign(tview.AlignCenter).
				SetExpansion(1))
		// Project
		work.Table.SetCell(i+1, 3,
			tview.
				NewTableCell(chronoWork.ProjectType.Name).
				SetAlign(tview.AlignCenter).
				SetExpansion(1))
		// Tags
		work.Table.SetCell(i+1, 4,
			tview.
				NewTableCell(chronoWork.GetCombinedTagNames()).
				SetAlign(tview.AlignCenter).
				SetExpansion(1))
	}
	work.tableCapture(tui)

	tagsOptions := append([]string{notSelectText})
	work.Form.
		AddInputField("Title", "", 50, nil, nil).
		AddDropDown("Project", append([]string{notSelectText}, models.AllProjectTypeNames(db.DB)...), 0, func(option string, optionIndex int) {
			if work.Form.GetFormItemByLabel("Tags") == nil {
				return
			}
			var projectType models.ProjectType
			result = db.DB.Preload("Tags").Where("name = ?", option).Find(&projectType)
			if result.Error != nil {
				tagsOptions = []string{notSelectText}
			} else {
				tagsOptions = append([]string{notSelectText}, projectType.GetTagNames()...)
			}
			work.Form.
				GetFormItemByLabel("Tags").(*tview.DropDown).
				SetOptions(tagsOptions, nil).
				SetCurrentOption(0)
		}).
		AddDropDown("Tags", tagsOptions, 0, nil).
		AddButton("Save", nil).
		AddButton("Cancel", func() {
			work.resetForm()
			tui.SetFocus("mainContent")
		})
	work.formCapture(tui)

	return work
}
