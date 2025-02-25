package model

import (
	"fmt"
)

// Node represents a node in the JSON tree
type Node struct {
	Key         string
	Value       any
	Parent      *Node
	Children    []*Node
	Expanded    bool
	Path        string
	Type        string
	MatchFilter bool
}

var (
	rootNode *Node
	currNode *Node
	filePath string
)

// SetRootNode sets the root node for the model
func SetRootNode(node *Node) {
	rootNode = node
	currNode = rootNode
}

// GetRootNode returns the root node
func GetRootNode() *Node {
	return rootNode
}

// SetFilePath sets the file path for the current JSON file
func SetFilePath(path string) {
	filePath = path
}

// GetFilePath returns the current JSON file path
func GetFilePath() string {
	return filePath
}

// BuildTree constructs a tree from JSON data
func BuildTree(key string, value any, parent *Node, path string) *Node {
	node := &Node{
		Key:         key,
		Value:       value,
		Parent:      parent,
		Children:    []*Node{},
		Expanded:    false,
		Path:        path,
		MatchFilter: false,
	}

	switch v := value.(type) {
	case map[string]any:
		node.Type = "object"
		for k, val := range v {
			childPath := path
			if childPath != "" {
				childPath += "."
			}
			childPath += k
			child := BuildTree(k, val, node, childPath)
			node.Children = append(node.Children, child)
		}
	case []any:
		node.Type = "array"
		for i, val := range v {
			k := fmt.Sprintf("[%d]", i)
			childPath := path + k
			child := BuildTree(k, val, node, childPath)
			node.Children = append(node.Children, child)
		}
	case string:
		node.Type = "string"
	case float64:
		node.Type = "number"
	case bool:
		node.Type = "boolean"
	case nil:
		node.Type = "null"
	}

	return node
}

// GetNodeValueString returns a string representation of the node's value
func GetNodeValueString(node *Node) string {
	switch node.Type {
	case "object":
		if node.Expanded {
			return fmt.Sprintf("{%d items}", len(node.Children))
		}
		return "{...}"
	case "array":
		if node.Expanded {
			return fmt.Sprintf("[%d items]", len(node.Children))
		}
		return "[...]"
	case "string":
		return fmt.Sprintf("\"%v\"", node.Value)
	case "number", "boolean":
		return fmt.Sprintf("%v", node.Value)
	case "null":
		return "null"
	default:
		return "unknown"
	}
}

// GetNodeValueType returns the type of the node's value as a string
func GetNodeValueType(node *Node) string {
	if node == nil || node.Value == nil {
		return "null"
	}

	switch v := node.Value.(type) {
	case string:
		return "string"
	case float64, int, int64, float32:
		return "number"
	case bool:
		return "boolean"
	case map[string]interface{}:
		return "object"
	case []interface{}:
		return "array"
	default:
		return fmt.Sprintf("%T", v)
	}
}

// IsExpandable returns true if the node has children and can be expanded
func IsExpandable(node *Node) bool {
	return node != nil && len(node.Children) > 0
}

// ToggleNodeExpansion finds and toggles the expansion state of a node based on its line representation
func ToggleNodeExpansion(line string) {
	rootNode := GetRootNode()
	showAll := true // We want to find the node even if it's filtered out

	var findAndToggle func(n *Node, depth int, isLast bool, prefix string, lineCounter *int) bool
	findAndToggle = func(n *Node, depth int, isLast bool, prefix string, lineCounter *int) bool {
		// Skip if node is nil
		if n == nil {
			return false
		}

		// Build the line representation for this node
		var nodeLine string
		if n.Key == "root" {
			nodeLine = "JSON Root"
		} else {
			nodeLine = n.Key + ": " + GetNodeValueString(n)
		}

		// Check if this is the line we're looking for
		if line == nodeLine && IsExpandable(n) {
			n.Expanded = !n.Expanded
			return true
		}

		if n.Expanded {
			for _, child := range n.Children {
				if showAll || child.MatchFilter {
					if findAndToggle(child, depth+1, false, prefix, lineCounter) {
						return true
					}
				}
			}
		}
		return false
	}

	findAndToggle(rootNode, 0, true, "", nil)
}

// ExpandAll expands all nodes in the tree
func ExpandAll(node *Node) {
	if node == nil {
		return
	}

	node.Expanded = true
	for _, child := range node.Children {
		ExpandAll(child)
	}
}

// CollapseAll collapses all nodes in the tree except the root node
func CollapseAll(node *Node) {
	if node == nil {
		return
	}

	// Don't collapse the root node
	if node.Key != "root" {
		node.Expanded = false
	}

	for _, child := range node.Children {
		CollapseAll(child)
	}
}
