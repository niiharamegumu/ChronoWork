package widgets

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/niiharamegumu/ChronoWork/db"
	"github.com/niiharamegumu/ChronoWork/models"
	"github.com/niiharamegumu/ChronoWork/pkg"
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
)

type Work struct {
	Table *tview.Table
}

func NewWork() *Work {
	work := &Work{
		Table: tview.NewTable().
			SetBorders(true).
			SetSelectable(true, false).
			SetFixed(1, 1).
			SetBordersColor(tview.Styles.BorderColor),
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

func GenerateInitWork(startTime, endTime time.Time) *Work {
	work := NewWork()
	var chronoWorks []models.ChronoWork
	var result *gorm.DB

	result = db.DB.
		Preload("ProjectType").
		Preload("ProjectType.Tags").
		Order("created_at desc").
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

	work.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 't':
				work.goToTop()
			case 'b':
				work.goToBottom()
			}
		}
		return event
	})

	return work
}
