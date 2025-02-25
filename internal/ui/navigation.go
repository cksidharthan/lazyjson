package ui

import (
	"github.com/cksidharthan/lazyjson/internal/filter"
	"github.com/cksidharthan/lazyjson/internal/model"

	"github.com/jroimartin/gocui"
)

// MoveDown handles down arrow navigation
func MoveDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

// MoveUp handles up arrow navigation
func MoveUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

// ToggleExpand expands or collapses the selected node
func ToggleExpand(g *gocui.Gui, v *gocui.View) error {
	// Get current line
	_, cy := v.Cursor()
	_, oy := v.Origin()
	lineNum := cy + oy

	// Find node at this line
	rootNode := model.GetRootNode()
	showAll := filter.IsShowingAll()

	var node *model.Node
	var findNode func(n *model.Node, depth int, isLast bool, prefix string, lineCounter *int) bool
	findNode = func(n *model.Node, depth int, isLast bool, prefix string, lineCounter *int) bool {
		// Skip if filtering is active and node doesn't match
		if !showAll && !n.MatchFilter {
			return false
		}

		if *lineCounter == lineNum {
			node = n
			return true
		}
		*lineCounter++

		if n.Expanded {
			newPrefix := prefix
			if depth > 0 {
				if isLast {
					newPrefix += "    "
				} else {
					newPrefix += "│   "
				}
			}

			visibleChildren := 0
			for _, child := range n.Children {
				if showAll || child.MatchFilter {
					visibleChildren++
				}
			}

			visibleCount := 0
			for _, child := range n.Children {
				if showAll || child.MatchFilter {
					if findNode(child, depth+1, visibleCount == visibleChildren-1, newPrefix, lineCounter) {
						return true
					}
					visibleCount++
				}
			}
		}
		return false
	}

	lineCounter := 0
	findNode(rootNode, 0, true, "", &lineCounter)

	if node != nil && model.IsExpandable(node) {
		node.Expanded = !node.Expanded
		RenderTree(v)
	}

	return nil
}

// ActivateFilter activates the filter input field
func ActivateFilter(g *gocui.Gui, v *gocui.View) error {
	filterView, err := g.View("filter")
	if err != nil {
		return err
	}

	filterView.Clear()
	g.SetCurrentView("filter")
	return nil
}

// ReturnToTree returns focus to the tree view
func ReturnToTree(g *gocui.Gui, v *gocui.View) error {
	g.SetCurrentView("tree")
	return nil
}

// ClearFilter clears the active filter
func ClearFilter(g *gocui.Gui, v *gocui.View) error {
	filter.SetShowAll(true)
	filter.ClearFilter(model.GetRootNode())

	filterView, err := g.View("filter")
	if err != nil {
		return err
	}
	filterView.Clear()
	filter.SetFilterText("")

	RenderTree(v)
	return nil
}

// Quit exits the application
func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
