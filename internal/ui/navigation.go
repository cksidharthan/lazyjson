package ui

import (
	"github.com/cksidharthan/lazyjson/internal/filter"
	"github.com/cksidharthan/lazyjson/internal/model"
	"github.com/jroimartin/gocui"
	"strings"
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

// ToggleExpand toggles the expansion state of the current node
func ToggleExpand(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, cy := v.Cursor()
		line, err := v.Line(cy)
		if err != nil {
			return err
		}

		// Clean up the line by removing tree symbols and whitespace
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "├── ")
		line = strings.TrimPrefix(line, "└── ")
		line = strings.TrimSuffix(line, " [-]")
		line = strings.TrimSuffix(line, " [+]")

		// Toggle node expansion
		model.ToggleNodeExpansion(line)
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

	// Don't clear the filter view, just set cursor to end of content
	if filterView.Buffer() != "" {
		lines := len(filterView.BufferLines())
		lastLine := lines - 1
		lastLineContent := filterView.BufferLines()[lastLine]
		filterView.SetCursor(len(lastLineContent), lastLine)
	}
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
	filterView, err := g.View("filter")
	if err != nil {
		return err
	}

	// Clear the filter view
	filterView.Clear()
	filterView.SetCursor(0, 0)

	// Reset filter state
	filter.SetFilterText("")
	filter.SetShowAll(true)
	filter.ClearFilter(model.GetRootNode())

	// Update tree view
	treeView, err := g.View("tree")
	if err != nil {
		return err
	}
	RenderTree(treeView)

	return nil
}

// ExpandAll expands all nodes in the tree
func ExpandAll(g *gocui.Gui, v *gocui.View) error {
	rootNode := model.GetRootNode()
	model.ExpandAll(rootNode)
	RenderTree(v)
	return nil
}

// CollapseAll collapses all nodes in the tree except the root node
func CollapseAll(g *gocui.Gui, v *gocui.View) error {
	rootNode := model.GetRootNode()
	// Keep root node expanded but collapse all children
	rootNode.Expanded = true
	for _, child := range rootNode.Children {
		model.CollapseAll(child)
	}
	RenderTree(v)
	return nil
}

// Quit exits the application
func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
