package widgets

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/niiharamegumu/ChronoWork/db"
	"github.com/niiharamegumu/ChronoWork/models"
	"github.com/niiharamegumu/ChronoWork/pkg"
	"github.com/niiharamegumu/ChronoWork/service"
	"github.com/rivo/tview"
)

var (
	workHeader = []string{
		"ID",
		"TotalTime",
		"Title",
		"Project",
		"Tags",
		"TRACKING",
	}
)

type Work struct {
	Table *tview.Table
}

func NewWork() *Work {
	work := &Work{
		Table: tview.NewTable().
			SetSelectable(true, false).
			SetFixed(1, 1),
	}
	return work
}

func (w *Work) GenerateInitWork(tui *service.TUI) (*Work, error) {
	w.setHeader()
	if err := w.setBody(pkg.RelativeStartTime(), pkg.TodayEndTime()); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Work) TableCapture(tui *service.TUI, form *Form, timer *Timer) {
	w.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 't':
				// table top
				w.goToTop()
			case 'b':
				// table bottom
				w.goToBottom()
			case 'a':
				// add new work
				form.Form.Clear(true)
				form.ConfigureStoreForm(tui, w)
				tui.SetFocus("mainWorkForm")
			case 'u':
				// update work
				row, _ := w.Table.GetSelection()
				cell := w.Table.GetCell(row, 0)
				if cell.Text == "" {
					break
				}
				id := cell.Text
				if intId, err := strconv.ParseUint(id, 10, 0); err == nil {
					uintId := uint(intId)
					chronoWork, err := models.FindChronoWork(db.DB, uintId)
					if err != nil {
						log.Println(err)
						break
					}
					form.Form.Clear(true)
					form.configureUpdateForm(tui, w, &chronoWork)
					tui.SetFocus("mainWorkForm")
				}
			case 'r':
				// reset timer or update timer
				row, _ := w.Table.GetSelection()
				cell := w.Table.GetCell(row, 0)
				if cell.Text == "" {
					break
				}
				id := cell.Text
				if intId, err := strconv.ParseUint(id, 10, 0); err == nil {
					uintId := uint(intId)
					chronoWork, err := models.FindChronoWork(db.DB, uintId)
					if err != nil {
						log.Println(err)
						break
					}
					if chronoWork.TotalSeconds == 0 {
						break
					}
					form.Form.Clear(true)
					form.configureTimerForm(tui, w, &chronoWork)
					tui.SetFocus("mainWorkForm")
				}
			case 'd':
				// delete work
				row, _ := w.Table.GetSelection()
				cell := w.Table.GetCell(row, 0)
				if cell.Text == "" {
					break
				}
				modal := tview.NewModal().
					SetText("Are you sure you want to delete this work?").
					AddButtons([]string{"Yes", "No"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if buttonLabel == "Yes" {
							id := cell.Text
							if intId, err := strconv.ParseUint(id, 10, 0); err == nil {
								uintId := uint(intId)
								if err := models.DeleteChronoWork(db.DB, uintId); err != nil {
									log.Println(err)
								}
								if err := w.ReStoreTable(pkg.RelativeStartTime(), pkg.TodayEndTime()); err != nil {
									log.Println(err)
								}
							}
						}
						tui.DeleteModal()
						tui.SetFocus("mainWorkContent")
						w.goToTop()
					})
				tui.SetModal(modal)
				tui.SetFocus("modal")
			}
		case tcell.KeyEnter:
			// toggle tracking
			row, _ := w.Table.GetSelection()
			cell := w.Table.GetCell(row, 0)
			if cell.Text == "" {
				break
			}
			id := cell.Text
			if intId, err := strconv.ParseUint(id, 10, 0); err == nil {
				uintId := uint(intId)
				chronoWork, err := models.FindChronoWork(db.DB, uintId)
				if err != nil {
					log.Println(err)
					break
				}
				chronoWorks, err := models.FindTrackingChronoWorks(db.DB)
				if err != nil {
					log.Println(err)
					break
				}
				// if tracking work exists, stop tracking
				if len(chronoWorks) > 0 {
					for _, cw := range chronoWorks {
						if cw.ID != chronoWork.ID {
							if err := cw.StopTrackingChronoWorks(db.DB); err != nil {
								log.Println(err)
								break
							}
						}
					}
				}
				// target tracking work
				if chronoWork.IsTracking {
					chronoWork.StopTrackingChronoWorks(db.DB)
					timer.ResetSetText()
					timer.StopCalculateSeconds()
				} else {
					chronoWork.StartTrackingChronoWork(db.DB)
					timer.SetStartTimer(chronoWork.StartTime)
					timer.SetCalculateSeconds(tui)
					timer.SetTimerText(chronoWork)
				}
				w.ReStoreTable(pkg.RelativeStartTime(), pkg.TodayEndTime())
				w.Table.Select(row, 0)
			}
		}
		return event
	})
}

func (w *Work) ReStoreTable(startTime, endTime time.Time) error {
	w.Table.Clear()
	w.setHeader()
	if err := w.setBody(startTime, endTime); err != nil {
		return err
	}
	return nil
}

func (w *Work) setHeader() {
	for i, header := range workHeader {
		w.Table.
			SetCell(0, i,
				tview.
					NewTableCell(header).
					SetAlign(tview.AlignCenter).
					SetTextColor(tcell.ColorWhite).
					SetBackgroundColor(tcell.ColorPurple).
					SetSelectable(false).
					SetExpansion(1),
			)
	}
}

func (w *Work) setBody(startTime, endTime time.Time) error {
	var chronoWork models.ChronoWork
	var chronoWorks []models.ChronoWork
	var err error

	chronoWorks, err = chronoWork.FindInRangeByTime(db.DB, startTime, endTime)
	if err != nil {
		log.Println(err)
		return err
	}

	activeTrackingChronoWorks, err := models.FindTrackingChronoWorks(db.DB)
	if err != nil {
		log.Println(err)
		return err
	}
	if len(activeTrackingChronoWorks) > 0 {
		for _, activeTrackingChronoWork := range activeTrackingChronoWorks {
			isInclude := false
			for _, cw := range chronoWorks {
				if cw.ID == activeTrackingChronoWork.ID {
					isInclude = true
					break
				}
			}
			if !isInclude {
				chronoWorks = append(chronoWorks, activeTrackingChronoWork)
			}
		}
	}

	date := time.Now().Format("2006/01/02")
	rowCount := 1
	for _, chronoWork := range chronoWorks {
		targetDate := chronoWork.CreatedAt.Format("2006/01/02")
		if date != targetDate {
			// blank row * 2
			for i := 0; i < 2; i++ {
				w.Table.SetCell(rowCount, 0,
					tview.NewTableCell("").SetSelectable(false))
				for i := 1; i < len(workHeader); i++ {
					w.Table.SetCell(rowCount, i,
						tview.NewTableCell("").SetSelectable(false))
				}
				rowCount++
			}

			date = targetDate
			// date row
			w.Table.SetCell(rowCount, 0,
				tview.
					NewTableCell(fmt.Sprintln(date, chronoWork.CreatedAt.Weekday())).
					SetAlign(tview.AlignCenter).
					SetTextColor(tcell.ColorWhite).
					SetBackgroundColor(tcell.ColorMediumPurple.TrueColor()).
					SetSelectable(false))
			for i := 1; i < len(workHeader); i++ {
				w.Table.SetCell(rowCount, i,
					tview.NewTableCell("").
						SetBackgroundColor(tcell.ColorMediumPurple.TrueColor()).
						SetSelectable(false))
			}
			rowCount++
		}
		w.configureTable(rowCount, chronoWork)
		rowCount++
	}

	if len(activeTrackingChronoWorks) > 0 {
		for i, chronoWork := range chronoWorks {
			if chronoWork.ID == activeTrackingChronoWorks[0].ID {
				w.Table.Select(i+1, 0)
				break
			}
		}
	}
	return nil
}

func (w *Work) goToTop() {
	w.Table.ScrollToBeginning().Select(1, 0)
}

func (w *Work) goToBottom() {
	w.Table.ScrollToEnd().Select(w.Table.GetRowCount()-1, 0)
}

func (w *Work) configureTable(row int, chronoWork models.ChronoWork) {
	// ID
	w.Table.SetCell(row, 0,
		tview.
			NewTableCell(fmt.Sprintf("%d", chronoWork.ID)).
			SetAlign(tview.AlignCenter).
			SetExpansion(0))
	// TotalTime
	w.Table.SetCell(row, 1,
		tview.
			NewTableCell(pkg.FormatTime(chronoWork.TotalSeconds)).
			SetAlign(tview.AlignCenter).
			SetExpansion(0))
	// Title
	w.Table.SetCell(row, 2,
		tview.
			NewTableCell(chronoWork.Title).
			SetAlign(tview.AlignCenter).
			SetExpansion(1))
	// Project
	w.Table.SetCell(row, 3,
		tview.
			NewTableCell(chronoWork.ProjectType.Name).
			SetAlign(tview.AlignCenter).
			SetExpansion(1))
	// Tags
	w.Table.SetCell(row, 4,
		tview.
			NewTableCell(chronoWork.Tag.Name).
			SetAlign(tview.AlignCenter).
			SetExpansion(1))
	// TRACKING
	if chronoWork.IsTracking {
		w.Table.SetCell(row, 5,
			tview.
				NewTableCell("Yes").
				SetAlign(tview.AlignCenter).
				SetTextColor(tcell.ColorGreen).
				SetExpansion(0))
	} else {
		w.Table.SetCell(row, 5,
			tview.
				NewTableCell("No").
				SetAlign(tview.AlignCenter).
				SetTextColor(tcell.ColorRed).
				SetExpansion(0))
	}
}
