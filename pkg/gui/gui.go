package gui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type GocuiManager struct{}

func (GocuiManager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	currentView := g.CurrentView()
	var currentViewName string
	if currentView != nil {
		currentViewName = currentView.Name()
	}

	// UrlEndpoint view -  top section
	if v, err := g.SetView("UrlEndpoint", 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "(1) GET"
		v.Wrap = true
		v.Autoscroll = false
		v.Editable = true
		v.Editor = gocui.DefaultEditor
		v.SelFgColor = gocui.ColorWhite
		v.BgColor = gocui.ColorCyan
		fmt.Fprintln(v, "https://httpbin.org/get")

		// make sure call once
		if _, err := g.SetCurrentView("UrlEndpoint"); err != nil {
			return err
		}
	}

	// Parameters view -  middle-left section
	if v, err := g.SetView("Params", 0, 3, maxX-1, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "(2)"
		v.Wrap = true
		v.Autoscroll = true
		v.Editable = true
		v.Editor = gocui.DefaultEditor
		v.SelFgColor = gocui.ColorWhite
		fmt.Fprintln(v, "param1=value1\nparam2=value2")

	}

	applyViewColors(g, currentViewName)
	return nil
}

func applyViewColors(g *gocui.Gui, currentViewName string) {
	viewNames := []string{"UrlEndpoint", "Params"}

	for _, name := range viewNames {
		view, err := g.View(name)
		if err != nil {
			continue
		}

		if name == currentViewName {
			view.SelFgColor = gocui.ColorGreen
		} else {
			view.SelFgColor = gocui.ColorWhite
		}
	}
}
