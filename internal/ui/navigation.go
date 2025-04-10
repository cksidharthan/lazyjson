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

	// Get the line number under the cursor
	_, cy := v.Cursor()

	// Retrieve the node path associated with this line number
	nodePath, exists := lineToNodePath[cy]
	if !exists || nodePath == "" {
		// If no path found for this line (e.g., empty view), do nothing
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
	restoreAndEnsureVisibleAfterFilter(model.GetRootNode())
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

// restoreAndEnsureVisibleAfterFilter traverses the node tree after filtering.
// It expands the parents of all nodes that match the filter and restores
// the pre-filter expansion state for all nodes.
func restoreAndEnsureVisibleAfterFilter(node *model.Node) {
	if node == nil {
		return
	}

	// --- Step 1: Ensure Visibility of Matches ---
	// If this node (or one of its children) matched the filter,
	// we must ensure its parent hierarchy is expanded so it's visible.
	if node.MatchFilter {
		parent := node.Parent
		for parent != nil {
			// Check the *saved* state before expanding. If the user had
			// explicitly collapsed a parent before filtering, we respect that
			// unless that parent itself matched the filter (handled when parent is processed).
			if savedState, exists := model.GetNodeExpansionState(parent.Path); !exists || savedState {
				parent.Expanded = true
			}
			parent = parent.Parent
		}
	}

	// --- Step 2: Restore Pre-Filter Expansion State ---
	// Regardless of matching, restore the node's original expansion state
	// that was saved before the filter was applied.
	if savedState, exists := model.GetNodeExpansionState(node.Path); exists {
		if model.IsExpandable(node) {
			node.Expanded = savedState
		}
	} else if model.IsExpandable(node) {
		// Default to collapsed if no prior state exists (e.g., new file)
		node.Expanded = false
	}

	// --- Step 3: Recurse ---
	// Process children. Crucially, we only need to recurse if the current
	// node is *now* expanded (after state restoration) OR if it's the root.
	// We must always process the children of the root node.
	if node.Expanded || node.Parent == nil {
		for _, child := range node.Children {
			restoreAndEnsureVisibleAfterFilter(child)
		}
	}
}

// Quit exits the application
func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
