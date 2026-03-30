package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

// UILayout defines the public contract for TUI layout initialization.
type UILayout interface {
	// Setup must build/register all views and bindings; returns error if it fails.
	Setup(gui *gocui.Gui) error
	// (Optional) View returns a named *gocui.View.
	View(name string) *gocui.View
}

type appLayout struct {
	views map[string]*gocui.View
}

// NewUILayout returns a UILayout object.
func NewUILayout() UILayout {
	return &appLayout{views: make(map[string]*gocui.View)}
}

// Setup implements UILayout; builds and registers all main views.
func (a *appLayout) Setup(gui *gocui.Gui) error {
	maxX, maxY := gui.Size()

	// UrlEndpoint view
	v, err := gui.SetView("UrlEndpoint", 0, 0, maxX-1, 2)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == nil {
		v.Title = "(1) GET"
		v.Wrap = true
		v.Autoscroll = false
		v.SelFgColor = gocui.ColorWhite
		v.BgColor = gocui.ColorCyan
		fmt.Fprintln(v, "https://httpbin.org/get")
	}
	a.views["UrlEndpoint"] = v

	// Params view
	width2 := maxX / 3
	pv, err := gui.SetView("Params", 0, 3, width2, maxY)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == nil {
		pv.Title = "(2) Params"
		pv.Wrap = true
		pv.Autoscroll = true
		pv.Editable = true
		pv.Editor = gocui.DefaultEditor
		pv.SelFgColor = gocui.ColorWhite
		fmt.Fprintln(pv, "param1=value1\nparam2=value2")
	}
	a.views["Params"] = pv

	// Headers view
	height1 := maxY / 2
	hv, err := gui.SetView("Headers", 0, height1, 2*width2-1, maxY-1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == nil {
		hv.Title = "(3) Headers"
		hv.Wrap = true
		hv.Autoscroll = true
		hv.Editable = true
		hv.Editor = gocui.DefaultEditor
		hv.SelFgColor = gocui.ColorWhite
		fmt.Fprintln(hv, "Content-Type=application/json\nAccept=application/json")
	}
	a.views["Headers"] = hv

	// ResponseBody view
	rv, err := gui.SetView("ResponseBody", width2, 3, maxX-1, maxY-1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == nil {
		rv.Title = "(4)"
		rv.Wrap = true
		rv.Autoscroll = true
		rv.Editable = false
		rv.Editor = gocui.DefaultEditor
		rv.SelFgColor = gocui.ColorWhite
		fmt.Fprintln(rv, "Response will appear here after sending request")
	}
	a.views["ResponseBody"] = rv

	return nil
}

// View implements UILayout's View accessor.
func (a *appLayout) View(name string) *gocui.View {
	return a.views[name]
}
