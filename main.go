package main

import (
	"log"

	"github.com/awesome-gocui/gocui"
	"github.com/coderj001/lazypost/pkg/gui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManager(gui.GocuiManager{})

	// Set up key bindings
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	// Set up Tab key to switch between views
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}

	// Start the main loop
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	currentView := g.CurrentView().Name()

	views := []string{"(1)", "floatingBox"}
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
