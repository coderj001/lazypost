package ui

import (
	"github.com/jroimartin/gocui"
	"testing"
)

func TestUILayoutSetup_Success(t *testing.T) {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		t.Fatalf("gocui.NewGui failed: %v", err)
	}
	defer g.Close()

	l := NewUILayout()
	err = l.Setup(g)
	if err != nil {
		t.Fatalf("UILayout.Setup failed: %v", err)
	}

	names := []string{"UrlEndpoint", "Params", "Headers", "ResponseBody"}
	for _, name := range names {
		view := l.View(name)
		if view == nil {
			t.Errorf("view %q not created", name)
		}
	}
}

// This test demonstrates intent; gocui does not easily allow duplicate view names in the same gui instance
// so we document the expected lack rather than inject a false test.
func TestUILayoutSetup_DuplicateView(t *testing.T) {
	t.Skip("gocui refuses duplicate SetView calls; test not feasible without deeper patching/mocking")
}
