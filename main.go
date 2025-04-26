package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func createUrlEndpointView(g *gocui.Gui, maxX int) error {
	v, err := g.SetView("UrlEndpoint", 0, 0, maxX-1, 2)
	if err != nil {
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

		// Make sure this view is initially focused
		if _, err := g.SetCurrentView("UrlEndpoint"); err != nil {
			return err
		}
	}
	return nil
}

func createParamsView(g *gocui.Gui, maxX, maxY int) error {
	v, err := g.SetView("Params", 0, 3, maxX-1, maxY/2)
	if err != nil {
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
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// UrlEndpoint view -  top section
	if err := createUrlEndpointView(g, maxX); err != nil {
		return err
	}

	// Parameters view -  middle-left section
	if err := createParamsView(g, maxX, maxY); err != nil {
		return err
	}

	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}

	// Set up Tab key to switch between views
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	currentView := g.CurrentView().Name()
	// views := []string{"UrlEndpoint", "Params", "Results"}
	views := []string{"UrlEndpoint", "Params"}

	currentIdx := -1
	for i, name := range views {
		if name == currentView {
			currentIdx = i
			break
		}
	}

	nextIdx := (currentIdx + 1) % len(views)
	g.SetCurrentView(views[nextIdx])
	return nil
}

// quit handles Ctrl+C to exit the application
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	//KeyBinding
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	// Start the main loop
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
