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

// ToggleExpand toggles the expansion state of the current node using its path
func ToggleExpand(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	// Get the line number under the cursor (relative to the view origin)
	_, cy := v.Cursor()
	// Get the view's origin (top-left visible line)
	_, oy := v.Origin()
	// Calculate the absolute line number in the buffer
	absoluteLine := oy + cy

	// Retrieve the node path associated with this absolute line number
	nodePath, exists := lineToNodePath[absoluteLine]
	if !exists || nodePath == "" {
		// If no path found for this line (e.g., empty view or error), do nothing
		return nil
	}

	// Toggle the node using its path via the model function
	toggled := model.ToggleNodeExpansionByPath(nodePath)

	// Re-render the tree only if a node was actually toggled
	if toggled {
		// Preserve cursor position and origin
		curX, curY := v.Cursor()
		origX, origY := v.Origin()

		RenderTree(v)

		// Restore cursor position and origin
		v.SetCursor(curX, curY)
		v.SetOrigin(origX, origY)
	}

	return nil
}

// ActivateFilter activates the filter input field
func ActivateFilter(g *gocui.Gui, v *gocui.View) error {
	filterView, err := g.View("filter")
	if err != nil {
		return err
	}

	// Clear the filter view before activating it
	filterView.Clear()
	filterView.SetCursor(0, 0)

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

// expandParentsUpwards ensures all parent nodes up to the root are expanded.
func expandParentsUpwards(node *model.Node) {
	parent := node.Parent
	for parent != nil {
		if model.IsExpandable(parent) {
			parent.Expanded = true
		}
		parent = parent.Parent
	}
}

// It expands nodes that match the filter, their entire parent hierarchy,
// and their entire subtree. Nodes not part of a matching hierarchy are collapsed.
// Returns true if the current node or any of its descendants matched the filter.
func ExpandMatchingHierarchy(node *model.Node) bool {
	if node == nil {
		return false
	}

	// Step 1: Recursively check if any children match
	isAnyChildMatching := false
	for _, child := range node.Children {
		// We need the result of the recursive call for the current node's logic
		if ExpandMatchingHierarchy(child) {
			isAnyChildMatching = true
		}
	}

	// Step 2: Determine if the current node is part of the matching hierarchy
	inMatchingHierarchy := node.MatchFilter || isAnyChildMatching

	// Step 3: Set expansion state based on hierarchy status
	if model.IsExpandable(node) {
		node.Expanded = inMatchingHierarchy // Expand if in hierarchy, collapse otherwise
	} else {
		node.Expanded = false // Non-expandable nodes are always collapsed
	}

	// Step 4: If in hierarchy, expand parents
	if inMatchingHierarchy {
		expandParentsUpwards(node) // Ensure visibility from root
	}

	// Step 5: Return whether this node is in the matching hierarchy
	return inMatchingHierarchy
}
