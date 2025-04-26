package gui

import (
	"log"

	"github.com/jroimartin/gocui"
)


func KeyBindSetup(g *gocui.Gui) {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	// Set up Tab key to switch between views
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}
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
