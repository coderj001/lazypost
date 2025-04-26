package main

import (
	"log"

	"github.com/jroimartin/gocui"
	"github.com/coderj001/lazypost/pkg/gui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManager(gui.GocuiManager{})

	//KeyBinding
	gui.KeyBindSetup(g)

	// Start the main loop
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
