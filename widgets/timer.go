package widgets

import (
	"context"
	"time"

	"github.com/niiharamegumu/ChronoWork/models"
	"github.com/niiharamegumu/ChronoWork/pkg"
	"github.com/niiharamegumu/ChronoWork/service"
	"github.com/rivo/tview"
)

type Timer struct {
	Wrapper     *tview.Grid
	Time        *tview.TextView
	Title       *tview.TextView
	ProjectName *tview.TextView
	TagName     *tview.TextView
	StartTime   time.Time
	cancelCtx   context.Context
	cancelFunc  context.CancelFunc
}

func NewTimer() *Timer {
	time := tview.NewTextView().
		SetLabel("Timer : ").
		SetText("00:00:00")
	title := tview.NewTextView().
		SetLabel("TItle : ")
	projectName := tview.NewTextView().
		SetLabel("Project Name : ")
	tagName := tview.NewTextView().
		SetLabel("Tag Name : ")
	timer := &Timer{
		Wrapper: tview.NewGrid().
			SetRows(0, 2, 2, 2).
			SetColumns(0).
			AddItem(time, 0, 0, 1, 1, 0, 0, false).
			AddItem(title, 1, 0, 1, 1, 0, 0, false).
			AddItem(projectName, 2, 0, 1, 1, 0, 0, false).
			AddItem(tagName, 3, 0, 1, 1, 0, 0, false),
		Time:        time,
		Title:       title,
		ProjectName: projectName,
		TagName:     tagName,
	}
	return timer
}

func (t *Timer) SetTimerText(c models.ChronoWork) {
	t.Title.SetText(c.Title)
	t.ProjectName.SetText(c.ProjectType.Name)
	t.TagName.SetText(c.Tag.Name)
}

func (t *Timer) SetStartTimer(startTime time.Time) {
	t.StartTime = startTime
}

func (t *Timer) SetCalculateSeconds(tui *service.TUI) {
	t.cancelCtx, t.cancelFunc = context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-t.cancelCtx.Done():
				return
			default:
				seconds := int(time.Now().Sub(t.StartTime).Seconds())
				tui.App.QueueUpdateDraw(func() {
					t.Time.SetText(pkg.FormatTime(seconds))
				})
				time.Sleep(time.Second)
			}
		}
	}()
}

func (t *Timer) StopCalculateSeconds() {
	t.cancelFunc()
}

func (t *Timer) ResetSetText() {
	t.Time.SetText("00:00:00")
	t.Title.SetText("")
	t.ProjectName.SetText("")
	t.TagName.SetText("")
}
