package ui

import (
	"fmt"

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
		fmt.Fprintln(v, "  Use ↑/↓ to navigate, Enter to expand/collapse, / to filter, q to quit")
	}

	// Filter view
	if v, err := g.SetView("filter", 0, maxY-7, maxX-1, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Filter"
		v.Editable = true
		v.Editor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
			filterEditor(g, v, key, ch, mod)
		})
	}

	// Help view
	if v, err := g.SetView("help", 0, maxY-4, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Help"
		fmt.Fprintln(v, "  ↑/↓: Navigate   Enter: Expand/Collapse   /: Filter   Escape: Clear Filter   q: Quit")
	}

	// Main view for JSON tree
	if v, err := g.SetView("tree", 0, 3, maxX-1, maxY-8); err != nil {
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

	return nil
}
