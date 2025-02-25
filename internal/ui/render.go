package ui

import (
	"github.com/cksidharthan/lazyjson/internal/filter"
	"github.com/cksidharthan/lazyjson/internal/model"
	"github.com/jroimartin/gocui"
)

// RenderTree renders the JSON tree in the view
func RenderTree(v *gocui.View) {
	v.Clear()
	rootNode := model.GetRootNode()
	showAll := filter.IsShowingAll()

	var traverse func(node *model.Node, depth int, isLast bool, prefix string)
	traverse = func(node *model.Node, depth int, isLast bool, prefix string) {
		// Skip if filtering is active and node doesn't match
		if !showAll && !node.MatchFilter {
			return
		}

		// Create the line prefix with appropriate indentation
		if depth > 0 {
			v.Write([]byte(prefix))
			if isLast {
				ColorizeTreeSymbol(v, "└── ")
			} else {
				ColorizeTreeSymbol(v, "├── ")
			}
		}

		// Create the display string
		if node.Key == "root" {
			ColorizeRoot(v, "JSON Root")
			// Always show expand/collapse for root if it has children
			if len(node.Children) > 0 {
				v.Write([]byte(" "))
				if node.Expanded {
					ColorizeExpand(v, "[-]")
				} else {
					ColorizeExpand(v, "[+]")
				}
			}
		} else {
			ColorizeKey(v, node.Key)
			v.Write([]byte(": "))
			ColorizeValue(v, model.GetNodeValueString(node), model.GetNodeValueType(node))

			// Add expand/collapse indicator for non-root nodes
			if model.IsExpandable(node) {
				v.Write([]byte(" "))
				if node.Expanded {
					ColorizeExpand(v, "[-]")
				} else {
					ColorizeExpand(v, "[+]")
				}
			}
		}
		v.Write([]byte("\n"))

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
				if showAll || child.MatchFilter {
					traverse(child, depth+1, visibleCount == visibleChildren-1, newPrefix)
					visibleCount++
				}
			}
		}
	}

	traverse(rootNode, 0, true, "")
}
