package ui

import (
	"fmt"
	"github.com/cksidharthan/lazyjson/internal/filter"
	"github.com/cksidharthan/lazyjson/internal/model"

	"github.com/jroimartin/gocui"
)

// RenderTree renders the JSON tree in the view
func RenderTree(v *gocui.View) {
	v.Clear()
	lines := []string{}

	rootNode := model.GetRootNode()
	showAll := filter.IsShowingAll()

	var traverse func(node *model.Node, depth int, isLast bool, prefix string) bool
	traverse = func(node *model.Node, depth int, isLast bool, prefix string) bool {
		// Skip if filtering is active and node doesn't match
		if !showAll && !node.MatchFilter {
			return false
		}

		// Create the line prefix with appropriate indentation
		indent := prefix
		if depth > 0 {
			if isLast {
				indent += "└── "
			} else {
				indent += "├── "
			}
		}

		// Create the display string
		var display string
		if node.Key == "root" {
			display = "JSON Root"
		} else {
			keyDisplay := node.Key
			valueDisplay := model.GetNodeValueString(node)
			display = fmt.Sprintf("%s: %s", keyDisplay, valueDisplay)
		}

		// Add expand/collapse indicator
		if model.IsExpandable(node) {
			if node.Expanded {
				display += " [-]"
			} else {
				display += " [+]"
			}
		}

		// Highlight matched nodes when filtering
		if !showAll && node.MatchFilter {
			display = "* " + display
		}

		lines = append(lines, indent+display)

		// If node is expanded, traverse its children
		if node.Expanded {
			newPrefix := prefix
			if depth > 0 {
				if isLast {
					newPrefix += "    "
				} else {
					newPrefix += "│   "
				}
			}

			visibleChildren := 0
			for _, child := range node.Children {
				if showAll || child.MatchFilter {
					visibleChildren++
				}
			}

			visibleCount := 0
			for _, child := range node.Children {
				if traverse(child, depth+1, visibleCount == visibleChildren-1, newPrefix) {
					visibleCount++
				}
			}
		}

		return true
	}

	traverse(rootNode, 0, true, "")

	// Display the lines with proper scrolling
	for _, line := range lines {
		fmt.Fprintln(v, line)
	}

	// Show filter status
	filterText := filter.GetFilterText()
	if !showAll {
		v.Title = fmt.Sprintf("JSON Tree (Filtered by: %s)", filterText)
	} else {
		v.Title = "JSON Tree"
	}
}
