package filter

import (
	"fmt"
	"strings"
)

// These imports may cause circular imports, so we use interface{} instead of model.Node
// In a real application, we would define interfaces to solve this problem

var (
	filterText string
	showAll    bool = true
	matchCount int  = 0
)

// SetFilterText sets the current filter text
func SetFilterText(text string) {
	filterText = text
}

// GetFilterText returns the current filter text
func GetFilterText() string {
	return filterText
}

// SetShowAll sets whether to show all nodes or only filtered nodes
func SetShowAll(value bool) {
	showAll = value
}

// IsShowingAll returns whether all nodes are being shown
func IsShowingAll() bool {
	return showAll
}

// GetMatchCount returns the current match count
func GetMatchCount() int {
	return matchCount
}

// ResetMatchCount resets the match count to zero
func ResetMatchCount() {
	matchCount = 0
}

// The Node interface to avoid circular imports
type Node interface {
	GetKey() string
	GetValue() any
	GetType() string
	GetChildren() []Node
	GetParent() Node
	SetMatchFilter(match bool)
	IsMatchFilter() bool
	SetExpanded(expanded bool)
}

// ApplyFilter applies a filter to the tree and returns whether the node matches
func ApplyFilter(node any, filter string) bool {
	// Reset match count when starting from root
	if n := node.(Node); n.GetParent() == nil {
		ResetMatchCount()
	}

	n := node.(Node)

	// Convert to lowercase for case-insensitive matching
	filterLower := strings.ToLower(filter)

	// Check if current node matches filter
	keyMatch := strings.Contains(strings.ToLower(n.GetKey()), filterLower)

	valueMatch := false
	nodeType := n.GetType()
	if nodeType == "string" {
		if strVal, ok := n.GetValue().(string); ok {
			valueMatch = strings.Contains(strings.ToLower(strVal), filterLower)
		}
	} else if nodeType == "number" || nodeType == "boolean" {
		valueMatch = strings.Contains(strings.ToLower(fmt.Sprintf("%v", n.GetValue())), filterLower)
	}

	// Mark this node as matching if it does
	matches := keyMatch || valueMatch
	if matches {
		matchCount++
	}
	n.SetMatchFilter(matches)

	// Apply filter recursively to children
	childMatch := false
	for _, child := range n.GetChildren() {
		if ApplyFilter(child, filter) {
			childMatch = true
		}
	}

	// Node matches if it or any of its children match
	matches = matches || childMatch
	n.SetMatchFilter(matches)

	return matches
}

// ClearFilter clears the filter from all nodes
func ClearFilter(node any) {
	n := node.(Node)
	n.SetMatchFilter(false)

	for _, child := range n.GetChildren() {
		ClearFilter(child)
	}
}

// ExpandFilterMatches expands nodes that match the filter
func ExpandFilterMatches(node any) {
	n := node.(Node)

	if n.IsMatchFilter() && (n.GetType() == "object" || n.GetType() == "array") {
		n.SetExpanded(true)

		// If this node matches, expand all its parents too
		parent := n.GetParent()
		for parent != nil {
			if pNode, ok := parent.(Node); ok {
				pNode.SetExpanded(true)
				parent = pNode.GetParent()
			} else {
				break
			}
		}
	}

	for _, child := range n.GetChildren() {
		ExpandFilterMatches(child)
	}
}
