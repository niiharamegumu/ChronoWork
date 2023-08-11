package widgets

import "github.com/rivo/tview"

func SimpleText(text string) *tview.TextView {
	textView := tview.NewTextView().
		SetText(text).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	return textView
}
