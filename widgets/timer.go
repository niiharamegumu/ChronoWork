package widgets

import (
	"context"
	"time"

	"github.com/niiharamegumu/ChronoWork/pkg"
	"github.com/niiharamegumu/ChronoWork/service"
	"github.com/rivo/tview"
)

type Timer struct {
	Timer      *tview.TextView
	StartTime  time.Time
	cancelCtx  context.Context
	cancelFunc context.CancelFunc
}

func NewTimer() *Timer {
	timer := &Timer{
		Timer: tview.NewTextView().
			SetMaxLines(1).
			SetTextAlign(tview.AlignCenter).
			SetText("00s"),
	}
	return timer
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
					t.Timer.SetText(pkg.FormatTime(seconds))
				})
				time.Sleep(time.Second)
			}
		}
	}()
}

func (t *Timer) StopCalculateSeconds() {
	t.cancelFunc()
}
