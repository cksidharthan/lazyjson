package filter

import (
	"fmt"
	"strings"
	
	"github.com/cksidharthan/lazyjson/internal/interfaces"
)

// These imports may cause circular imports, so we use interface{} instead of model.Node
// In a real application, we would define interfaces to solve this problem

var (
	filterText string
	showAll    bool = true
	matchCount int
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

// Node interface represents a node in the JSON tree
// This is defined here to avoid circular imports with the model package
type Node = interfaces.Node

// ApplyFilter applies a filter to the tree and returns whether the node matches
func ApplyFilter(node Node, filter string) bool {
	// Reset match count when starting from root
	if node.GetParent() == nil {
		ResetMatchCount()
	}

	// Empty filter matches everything
	if filter == "" {
		node.SetMatchFilter(true)
		return true
	}

	// Convert to lowercase for case-insensitive matching
	filterLower := strings.ToLower(filter)

	// Check if current node matches filter
	keyMatch := strings.Contains(strings.ToLower(node.GetKey()), filterLower)

	// Check value match based on node type
	valueMatch := false
	nodeType := node.GetType()
	switch nodeType {
	case "string":
		if strVal, ok := node.GetValue().(string); ok {
			valueMatch = strings.Contains(strings.ToLower(strVal), filterLower)
		}
	case "number", "boolean":
		valueMatch = strings.Contains(strings.ToLower(fmt.Sprintf("%v", node.GetValue())), filterLower)
	}

	// Mark this node as matching if it matches directly
	directMatch := keyMatch || valueMatch
	if directMatch {
		matchCount++
	}

	// Apply filter recursively to children
	childMatch := false
	for _, child := range node.GetChildren() {
		if ApplyFilter(child, filter) {
			childMatch = true
		}
	}

	// Node matches if it or any of its children match
	matches := directMatch || childMatch
	node.SetMatchFilter(matches)

	return matches
}

// ClearFilter clears the filter from all nodes
func ClearFilter(node Node) {
	if node == nil {
		return
	}

	node.SetMatchFilter(false)

	for _, child := range node.GetChildren() {
		ClearFilter(child)
	}
}
