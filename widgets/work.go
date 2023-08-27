package widgets

import (
	"fmt"
	"log"
	"strconv"

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
			SetBorders(true).
			SetSelectable(true, false).
			SetFixed(1, 1).
			SetBordersColor(tview.Styles.BorderColor),
	}
	return work
}

func (w *Work) GenerateInitWork(tui *service.TUI) (*Work, error) {
	w.setHeader()
	if err := w.setBody(); err != nil {
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
				tui.SetFocus("mainForm")
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
					tui.SetFocus("mainForm")
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
					tui.SetFocus("mainForm")
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
								if err := w.ReStoreTable(); err != nil {
									log.Println(err)
								}
							}
						}
						tui.DeleteModal()
						tui.SetFocus("mainContent")
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
				w.ReStoreTable()
			}
		}
		return event
	})
}

func (w *Work) ReStoreTable() error {
	w.Table.Clear()
	w.setHeader()
	if err := w.setBody(); err != nil {
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
					SetTextColor(tcell.ColorPurple).
					SetSelectable(false).
					SetExpansion(1),
			)
	}
}

func (w *Work) setBody() error {
	var chronoWork models.ChronoWork
	var chronoWorks []models.ChronoWork
	var err error

	chronoWorks, err = chronoWork.FindInRangeByTime(db.DB, pkg.TodayStartTime(), pkg.TodayEndTime())
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

	for i, chronoWork := range chronoWorks {
		w.configureTable(i, chronoWork)
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
	w.Table.SetCell(row+1, 0,
		tview.
			NewTableCell(fmt.Sprintf("%d", chronoWork.ID)).
			SetAlign(tview.AlignCenter).
			SetExpansion(1))
	// TotalTime
	w.Table.SetCell(row+1, 1,
		tview.
			NewTableCell(pkg.FormatTime(chronoWork.TotalSeconds)).
			SetAlign(tview.AlignCenter).
			SetExpansion(1))
	// Title
	w.Table.SetCell(row+1, 2,
		tview.
			NewTableCell(chronoWork.Title).
			SetAlign(tview.AlignCenter).
			SetExpansion(1))
	// Project
	w.Table.SetCell(row+1, 3,
		tview.
			NewTableCell(chronoWork.ProjectType.Name).
			SetAlign(tview.AlignCenter).
			SetExpansion(1))
	// Tags
	w.Table.SetCell(row+1, 4,
		tview.
			NewTableCell(chronoWork.Tag.Name).
			SetAlign(tview.AlignCenter).
			SetExpansion(1))
	// TRACKING
	if chronoWork.IsTracking {
		w.Table.SetCell(row+1, 5,
			tview.
				NewTableCell("Yes").
				SetAlign(tview.AlignCenter).
				SetTextColor(tcell.ColorGreen).
				SetExpansion(1))
	} else {
		w.Table.SetCell(row+1, 5,
			tview.
				NewTableCell("No").
				SetAlign(tview.AlignCenter).
				SetTextColor(tcell.ColorRed).
				SetExpansion(1))
	}
}
