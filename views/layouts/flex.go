package layouts

import "github.com/rivo/tview"

func FlexRow() *tview.Flex {
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)

	return flex
}
