package ui

import (
	"fmt"
	"strings"

	"github.com/cksidharthan/lazyjson/internal/filter"
	"github.com/cksidharthan/lazyjson/internal/model"

	"github.com/jroimartin/gocui"
)

// SetupKeybindings configures all keyboard shortcuts
func SetupKeybindings(g *gocui.Gui) error {
	// Quit
	if err := g.SetKeybinding("", 'q', gocui.ModNone, Quit); err != nil {
		return err
	}

	// Navigation
	if err := g.SetKeybinding("tree", gocui.KeyArrowDown, gocui.ModNone, MoveDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("tree", gocui.KeyArrowUp, gocui.ModNone, MoveUp); err != nil {
		return err
	}

	// Expand/collapse
	if err := g.SetKeybinding("tree", gocui.KeyEnter, gocui.ModNone, ToggleExpand); err != nil {
		return err
	}

	// Activate filter
	if err := g.SetKeybinding("tree", '/', gocui.ModNone, ActivateFilter); err != nil {
		return err
	}

	// Clear filter with Escape
	if err := g.SetKeybinding("tree", gocui.KeyEsc, gocui.ModNone, ClearFilter); err != nil {
		return err
	}

	// Return to tree from filter
	if err := g.SetKeybinding("filter", gocui.KeyEsc, gocui.ModNone, ReturnToTree); err != nil {
		return err
	}

	// Expand all nodes
	if err := g.SetKeybinding("tree", 'e', gocui.ModNone, ExpandAll); err != nil {
		return err
	}

	// Collapse all nodes
	if err := g.SetKeybinding("tree", 'c', gocui.ModNone, CollapseAll); err != nil {
		return err
	}

	return nil
}

// filterEditor handles input in the filter field
func filterEditor(g *gocui.Gui, v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case key == gocui.KeyEnter:
		// Apply filter
		filterText := strings.TrimSpace(v.Buffer())
		filter.SetFilterText(filterText)

		if filterText != "" {
			filter.SetShowAll(false)
			rootNode := model.GetRootNode()
			filter.ApplyFilter(rootNode, filterText)
			filter.ExpandFilterMatches(rootNode)

			// Update matches view
			if matchesView, err := g.View("matches"); err == nil {
				matchesView.Clear()
				fmt.Fprintf(matchesView, " %d", filter.GetMatchCount())
			}
		} else {
			filter.SetShowAll(true)
			filter.ClearFilter(model.GetRootNode())
			// Clear matches view
			if matchesView, err := g.View("matches"); err == nil {
				matchesView.Clear()
			}
		}

		g.SetCurrentView("tree")
		treeView, _ := g.View("tree")
		RenderTree(treeView)
	default:
		gocui.DefaultEditor.Edit(v, key, ch, mod)
	}
}
