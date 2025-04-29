package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

// Views
const (
	UrlEndpointView = "UrlEndpoint"
	ParamsView      = "ParamsView"
	// TODO: need to be implemented
	ResponseBodyView = "ResponseBody"
	HeadersView      = "Headers"
)

var (
	AllViews = []string{
		UrlEndpointView,
		ParamsView,
		ResponseBodyView,
	}
	active = 0
)

func createUrlEndpointView(g *gocui.Gui, maxX int) error {
	v, err := g.SetView(UrlEndpointView, 0, 0, maxX-1, 2)
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
		if _, err := g.SetCurrentView(UrlEndpointView); err != nil {
			return err
		}
	}
	return nil
}

func createParamsView(g *gocui.Gui, maxX, maxY int) error {
	width2 := maxX / 3
	v, err := g.SetView(ParamsView, 0, 3, width2-1, maxY-1)
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

func createResponseBodyView(g *gocui.Gui, maxX, maxY int) error {
	width2 := maxX / 3
	v, err := g.SetView(ResponseBodyView, width2, 3, maxX-1, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "(3)"
		v.Wrap = true
		v.Autoscroll = true
		v.Editable = true
		v.Editor = gocui.DefaultEditor
		v.SelFgColor = gocui.ColorWhite
		fmt.Fprintln(v, "ResponseBody")
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

	// ResponseBody view - right cover
	if err := createResponseBodyView(g, maxX, maxY); err != nil {
		return err
	}

	return nil
}

// Add keybindings here
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

	currentIdx := -1
	for i, name := range AllViews {
		if name == currentView {
			currentIdx = i
			break
		}
	}

	nextIdx := (currentIdx + 1) % len(AllViews)
	g.SetCurrentView(AllViews[nextIdx])
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
	
	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

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
