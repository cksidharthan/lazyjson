package model

import (
	"fmt"
)

// Node represents a node in the JSON tree
type Node struct {
	Key         string
	Value       interface{}
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

// BuildTree constructs a tree from JSON data
func BuildTree(key string, value interface{}, parent *Node, path string) *Node {
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
	case map[string]interface{}:
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
	case []interface{}:
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

// IsExpandable returns true if the node can be expanded/collapsed
func IsExpandable(node *Node) bool {
	return node.Type == "object" || node.Type == "array"
}
