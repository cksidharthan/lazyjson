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
		v.FgColor = gocui.ColorYellow
		fmt.Fprintf(v, "  \033[36mFile:\033[0m \033[33m%s\033[0m", model.GetFilePath())
	}

	// Filter view
	if v, err := g.SetView("filter", 0, 3, maxX-1, 5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Filter"
		v.Editable = true
		v.FgColor = gocui.ColorGreen
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

	// Help view
	if v, err := g.SetView("help", 0, maxY-4, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Help"
		v.FgColor = gocui.ColorCyan

		// Navigation help with custom colors
		fmt.Fprintf(v, "  \033[33mNavigation\033[0m │ \033[36m↑/↓\033[0m: Move cursor   \033[36mEnter\033[0m: Expand/Collapse   \033[36me\033[0m: Expand all   \033[36mc\033[0m: Collapse all\n")
		fmt.Fprintf(v, "  \033[33mFilter\033[0m     │ \033[36m/\033[0m: Enable filter   \033[36mEsc\033[0m: Clear filter   \033[36mEnter\033[0m: Apply filter\n")
		fmt.Fprintf(v, "  \033[33mOther\033[0m      │ \033[36mq\033[0m: Quit")
	}

	return nil
}
