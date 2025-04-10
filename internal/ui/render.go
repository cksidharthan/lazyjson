package ui

import (
	"github.com/cksidharthan/lazyjson/internal/model"
	"github.com/jroimartin/gocui"
)

// lineToNodePath maps the visible line number in the tree view
// to the unique Path string of the model.Node displayed on that line.
// This is rebuilt every time the tree is rendered and used by ToggleExpand
// to identify which node needs to be toggled based on cursor position.
var lineToNodePath map[int]string

// RenderTree renders the JSON tree in the view
func RenderTree(v *gocui.View) {
	if v == nil {
		return
	}

	// Reset the map and clear the view
	lineToNodePath = make(map[int]string)
	v.Clear()
	rootNode := model.GetRootNode()
	currentLine := 0

	// Recursive function to traverse and render the tree
	var renderNode func(node *model.Node, depth int, isLast bool, prefix string)
	renderNode = func(node *model.Node, depth int, isLast bool, prefix string) {
		// Store the path for this line number *before* rendering the line content
		lineToNodePath[currentLine] = node.Path

		// Render the node's line with proper indentation
		renderNodeLine(v, node, depth, isLast, prefix)

		// Increment line counter *after* rendering the line
		currentLine++

		// If node is expanded, render its children
		if node.Expanded {
			// Calculate the new prefix for children
			newPrefix := calculateChildPrefix(prefix, depth, isLast)

			// Count visible children for determining which is the last one
			visibleChildren := len(node.Children)

			// Render each child
			for i, child := range node.Children {
				isLastChild := (i == visibleChildren-1)
				renderNode(child, depth+1, isLastChild, newPrefix)
			}
		}
	}

	// Start rendering from the root node
	renderNode(rootNode, 0, true, "")

	// Ensure cursor stays within bounds after re-render
	_, y := v.Cursor()
	if y >= currentLine && currentLine > 0 {
		v.SetCursor(0, currentLine-1)
	}
}

// Helper function to render a single node line
func renderNodeLine(v *gocui.View, node *model.Node, depth int, isLast bool, prefix string) {
	// Add indentation and tree connectors for non-root nodes
	if depth > 0 {
		v.Write([]byte(prefix))
		if isLast {
			ColorizeTreeSymbol(v, "└── ")
		} else {
			ColorizeTreeSymbol(v, "├── ")
		}
	}

	// Render the node content
	if node.Key == "root" {
		// Special handling for root node
		ColorizeRoot(v, "JSON Root")
		if len(node.Children) > 0 {
			v.Write([]byte(" "))
			if node.Expanded {
				ColorizeExpand(v, "[-]")
			} else {
				ColorizeExpand(v, "[+]")
			}
		}
	} else {
		// Regular node with key and value
		ColorizeKey(v, node.Key)
		v.Write([]byte(": "))
		ColorizeValue(v, model.GetNodeValueString(node), model.GetNodeValueType(node))

		// Add expand/collapse indicator if node has children
		if model.IsExpandable(node) {
			v.Write([]byte(" "))
			if node.Expanded {
				ColorizeExpand(v, "[-]")
			} else {
				ColorizeExpand(v, "[+]")
			}
		}
	}
	
	// End the line
	v.Write([]byte("\n"))
}

// Helper function to calculate the prefix for child nodes
func calculateChildPrefix(parentPrefix string, parentDepth int, parentIsLast bool) string {
	if parentDepth == 0 {
		return ""
	}
	
	if parentIsLast {
		return parentPrefix + "    " // Space after last item
	} else {
		return parentPrefix + "│   " // Vertical line for non-last items
	}
}

// Helper function to count visible children
func countVisibleChildren(node *model.Node) int {
	return len(node.Children)
}
