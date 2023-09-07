package widgets

import (
	"fmt"
	"log"
	"sort"
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
	dateSeparateRow = 3
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
						if cw.ID != chronoWork.ID || !pkg.IsToday(cw.CreatedAt) {
							if err := cw.StopTrackingChronoWorks(db.DB); err != nil {
								log.Println(err)
								break
							}
						}
					}
				}
				if pkg.IsToday(chronoWork.CreatedAt) {
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
				} else {
					// chronowork copy
					newChronoWork, err := models.CreateChronoWork(db.DB, chronoWork.Title, chronoWork.ProjectTypeID, chronoWork.TagID)
					if err != nil {
						log.Println(err)
						break
					}
					newChronoWork.StartTrackingChronoWork(db.DB)
					timer.SetStartTimer(newChronoWork.StartTime)
					timer.SetCalculateSeconds(tui)
					timer.SetTimerText(newChronoWork)
					w.ReStoreTable(pkg.RelativeStartTime(), pkg.TodayEndTime())
					w.Table.Select(1, 0)
				}
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
		tableCell := tview.NewTableCell(header).
			SetTextColor(tcell.ColorWhite).
			SetBackgroundColor(tcell.ColorPurple).
			SetSelectable(false)
		if header != "ID" {
			tableCell.SetExpansion(1)
		}
		if header == "TotalTime" || header == "TRACKING" {
			tableCell.SetAlign(tview.AlignCenter)
		} else {
			tableCell.SetAlign(tview.AlignLeft)
		}

		w.Table.SetCell(0, i, tableCell)
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

	if len(chronoWorks) == 0 {
		return nil
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
		sort.Slice(chronoWorks, func(i, j int) bool {
			return chronoWorks[i].CreatedAt.After(chronoWorks[j].CreatedAt)
		})
	}

	groupedChronoWorks := map[string][]models.ChronoWork{}
	for _, work := range chronoWorks {
		dateStr := work.CreatedAt.Format("2006/01/02")
		groupedChronoWorks[dateStr] = append(groupedChronoWorks[dateStr], work)
	}

	today := time.Now()
	rowCount := 1
	for dateStr, chronoWorks := range groupedChronoWorks {
		totalSecondsByDay := 0
		date, err := time.Parse("2006/01/02", dateStr)
		if err != nil {
			log.Println(err)
			return err
		}
		if date.Year() != today.Year() || date.Month() != today.Month() || date.Day() != today.Day() {
			for i := 0; i < dateSeparateRow; i++ {
				w.insertBlankRow(rowCount)
				rowCount++
			}
			w.insertDateRow(rowCount, date.Format("2006/01/02"), date.Weekday())
			rowCount++
		}
		for _, chronoWork := range chronoWorks {
			w.configureTable(rowCount, chronoWork)
			rowCount++
			totalSecondsByDay += chronoWork.TotalSeconds
		}
		w.insertTotalSecondsByDayRow(rowCount, totalSecondsByDay, len(chronoWorks))
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

func (w *Work) insertBlankRow(rowCount int) {
	for i := 0; i < len(workHeader); i++ {
		w.Table.SetCell(rowCount, i, tview.NewTableCell("").SetSelectable(false))
	}
}

func (w *Work) insertTotalSecondsByDayRow(rowCount int, totalSecondsByDay int, count int) {
	w.Table.SetCell(rowCount, 0,
		tview.NewTableCell("Total").
			SetAlign(tview.AlignLeft).
			SetTextColor(tcell.ColorWhite).
			SetBackgroundColor(tcell.ColorRebeccaPurple).
			SetSelectable(false))
	w.Table.SetCell(rowCount, 1,
		tview.NewTableCell(pkg.FormatTime(totalSecondsByDay)).
			SetAlign(tview.AlignCenter).
			SetTextColor(tcell.ColorWhite).
			SetBackgroundColor(tcell.ColorRebeccaPurple).
			SetSelectable(false))
	w.Table.SetCell(rowCount, 2,
		tview.NewTableCell(fmt.Sprintf("count:%d", count)).
			SetAlign(tview.AlignLeft).
			SetTextColor(tcell.ColorWhite).
			SetBackgroundColor(tcell.ColorRebeccaPurple).
			SetSelectable(false))

	for i := 3; i < len(workHeader); i++ {
		w.Table.SetCell(rowCount, i,
			tview.NewTableCell("").
				SetBackgroundColor(tcell.ColorRebeccaPurple).
				SetSelectable(false))
	}
}

func (w *Work) insertDateRow(rowCount int, date string, weekday time.Weekday) {
	w.Table.SetCell(rowCount, 0,
		tview.NewTableCell(fmt.Sprintf("%s %s", date, weekday)).
			SetAlign(tview.AlignCenter).
			SetTextColor(tcell.ColorWhite).
			SetBackgroundColor(tcell.ColorMediumPurple.TrueColor()).
			SetSelectable(false))

	for i := 1; i < len(workHeader); i++ {
		if workHeader[i] == "TRACKING" {
			w.Table.SetCell(rowCount, i,
				tview.NewTableCell("Copy to Today").
					SetAlign(tview.AlignCenter).
					SetBackgroundColor(tcell.ColorMediumPurple.TrueColor()).
					SetSelectable(false))
		} else {
			w.Table.SetCell(rowCount, i,
				tview.NewTableCell("").
					SetBackgroundColor(tcell.ColorMediumPurple.TrueColor()).
					SetSelectable(false))
		}
	}
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
			SetAlign(tview.AlignLeft).
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
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
	// Project
	w.Table.SetCell(row, 3,
		tview.
			NewTableCell(chronoWork.ProjectType.Name).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
	// Tags
	w.Table.SetCell(row, 4,
		tview.
			NewTableCell(chronoWork.Tag.Name).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
	// TRACKING
	trackingCell := tview.NewTableCell("").SetAlign(tview.AlignCenter).SetExpansion(0)
	setText := "Yes"
	setColor := tcell.ColorGreen
	if pkg.IsToday(chronoWork.CreatedAt) {
		if !chronoWork.IsTracking {
			setText = "No"
			setColor = tcell.ColorRed
		}
	} else {
		setText = "Copy"
		if !chronoWork.IsTracking {
			setColor = tcell.ColorRed
		}
	}
	trackingCell.SetText(setText).SetTextColor(setColor)
	w.Table.SetCell(row, 5, trackingCell)
}
