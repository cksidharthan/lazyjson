package ui

import (
	"fmt"

	"github.com/cksidharthan/lazyjson/internal/model"
	"github.com/jroimartin/gocui"
)

// Layout sets up the UI layout
func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Title view
	if v, err := g.SetView("title", 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "JSON Viewer"
		fmt.Fprintf(v, "  File: %s", model.GetFilePath())
	}

	// Filter view (moved to top)
	if v, err := g.SetView("filter", 0, 3, maxX-1, 5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Filter"
		v.Editable = true
		v.Editor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
			filterEditor(g, v, key, ch, mod)
		})
	}

	// Main view for JSON tree
	if v, err := g.SetView("tree", 0, 6, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "JSON Tree"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		RenderTree(v)

		if _, err := g.SetCurrentView("tree"); err != nil {
			return err
		}
	}

	// Help view (at bottom)
	if v, err := g.SetView("help", 0, maxY-4, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Help"
		fmt.Fprintln(v, "  Navigation │ ↑/↓: Move cursor   Enter: Expand/Collapse node   e: Expand all   c: Collapse all")
		fmt.Fprintln(v, "  Filter     │ /: Enable filter   Esc: Clear filter   Enter: Apply filter")
		fmt.Fprintln(v, "  Other      │ q: Quit")
	}

	return nil
}
