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
	// Unique, dot-separated path for node identification (e.g., "root.data.items[0].name")
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
func BuildTree(key string, value any, parent *Node, currentPath string) *Node {
	node := &Node{
		Key:         key,
		Value:       value,
		Parent:      parent,
		Children:    []*Node{},
		// Expand all expandable nodes by default initially.
		// Expanded:    depth <= 3, // Old logic: only expand top 3 levels
		Path:        currentPath,
		MatchFilter: false,
	}

	// Determine node type and set initial expansion state
	switch v := value.(type) {
	case map[string]any:
		node.Type = "object"
		// Objects are expandable, expand by default
		node.Expanded = true
		for k, val := range v {
			// Construct the unique path for the child node
			childPath := currentPath
			if childPath != "" {
				childPath += "."
			}
			childPath += k
			child := BuildTree(k, val, node, childPath)
			node.Children = append(node.Children, child)
		}
	case []any:
		node.Type = "array"
		// Arrays are expandable, expand by default
		node.Expanded = true
		for i, val := range v {
			k := fmt.Sprintf("[%d]", i)
			// Construct the unique path for the array element
			childPath := currentPath + k
			child := BuildTree(k, val, node, childPath)
			node.Children = append(node.Children, child)
		}
	case string:
		node.Type = "string"
		// Non-expandable nodes are never expanded
		node.Expanded = false
	case float64, float32:
		node.Type = "number"
		node.Expanded = false
	case int, int8, int16, int32, int64:
		node.Type = "number"
		node.Expanded = false
	case bool:
		node.Type = "boolean"
		node.Expanded = false
	case nil:
		node.Type = "null"
		node.Expanded = false
	default:
		node.Type = "unknown"
		node.Expanded = false
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

// findNodeByPath recursively searches for a node by its unique path
func findNodeByPath(node *Node, targetPath string) *Node {
	if node == nil {
		return nil
	}
	if node.Path == targetPath {
		return node
	}
	for _, child := range node.Children {
		if found := findNodeByPath(child, targetPath); found != nil {
			return found
		}
	}
	return nil
}

// ToggleNodeExpansionByPath finds a node by its unique path and toggles its expansion state
// It returns true if a node was found and toggled, false otherwise.
func ToggleNodeExpansionByPath(targetPath string) bool {
	if targetPath == "" {
		return false
	}
	nodeToToggle := findNodeByPath(rootNode, targetPath)

	if nodeToToggle != nil && IsExpandable(nodeToToggle) {
		nodeToToggle.Expanded = !nodeToToggle.Expanded
		return true
	}
	return false
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
