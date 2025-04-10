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
	// expansionStates stores the desired expansion state (true=expanded)
	// keyed by the unique node Path. Used to preserve state across filtering.
	expansionStates map[string]bool = make(map[string]bool)
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

// StoreNodeExpansionState stores the expansion state of a node by its path
func StoreNodeExpansionState(path string, expanded bool) {
	if path == "" {
		return
	}
	expansionStates[path] = expanded
}

// GetNodeExpansionState retrieves the stored expansion state of a node by its path
func GetNodeExpansionState(path string) (bool, bool) {
	if path == "" {
		return false, false
	}
	// Returns the stored state and whether a state was actually stored for this path.
	expanded, exists := expansionStates[path]
	return expanded, exists
}

// ClearExpansionStates clears all stored expansion states, typically called
// before saving new states or when clearing filters.
func ClearExpansionStates() {
	expansionStates = make(map[string]bool)
}

// BuildTree constructs a tree from JSON data
func BuildTree(key string, value any, parent *Node, currentPath string) *Node {
	// Calculate current depth
	depth := 0
	p := parent
	for p != nil {
		depth++
		p = p.Parent
	}

	node := &Node{
		Key:      key,
		Value:    value,
		Parent:   parent,
		Children: []*Node{},
		// Auto-expand nodes up to depth 3
		Expanded:    depth <= 3,
		// Unique, dot-separated path for node identification (e.g., "root.data.items[0].name")
		Path:        currentPath,
		MatchFilter: false,
	}

	switch v := value.(type) {
	case map[string]any:
		node.Type = "object"
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
		for i, val := range v {
			k := fmt.Sprintf("[%d]", i)
			// Construct the unique path for the array element
			childPath := currentPath + k
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

// saveExpansionStateRecursive is a helper to recursively save states
func saveExpansionStateRecursive(node *Node) {
	if node == nil {
		return
	}
	// Store the state regardless of expandability, as it might become expandable later
	StoreNodeExpansionState(node.Path, node.Expanded)
	for _, child := range node.Children {
		saveExpansionStateRecursive(child)
	}
}

// SaveAllExpansionStates iterates through the entire tree and saves
// the current expansion state of each node to the expansionStates map.
func SaveAllExpansionStates() {
	ClearExpansionStates() // Clear previous states before saving new ones
	saveExpansionStateRecursive(rootNode)
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
		// Optionally update stored state if needed immediately, though maybe better after render?
		// StoreNodeExpansionState(nodeToToggle.Path, nodeToToggle.Expanded)
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
	StoreNodeExpansionState(node.Path, node.Expanded)
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
		StoreNodeExpansionState(node.Path, node.Expanded)
	}

	for _, child := range node.Children {
		CollapseAll(child)
	}
}
