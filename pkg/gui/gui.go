package gui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type GocuiManager struct{}

func (GocuiManager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("URL", 0, 0, maxX-1, maxY/2-1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "(1)"
		v.Wrap = true
		v.Autoscroll = false
		v.SelBgColor = gocui.ColorCyan
		v.SelFgColor = gocui.ColorWhite
		fmt.Fprintln(v, "https://httpbin.org/get")
		v.BgColor = gocui.ColorGreen
		g.SetCurrentView("URL")
	}
	return nil
}
