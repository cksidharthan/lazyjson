package model

// GetKey returns the node's key
func (n *Node) GetKey() string {
	return n.Key
}

// GetValue returns the node's value
func (n *Node) GetValue() interface{} {
	return n.Value
}

// GetType returns the node's type
func (n *Node) GetType() string {
	return n.Type
}

// GetChildren returns the node's children as a slice of filter.Node
func (n *Node) GetChildren() []interface{} {
	children := make([]interface{}, len(n.Children))
	for i, child := range n.Children {
		children[i] = child
	}
	return children
}

// GetParent returns the node's parent as a filter.Node
func (n *Node) GetParent() interface{} {
	return n.Parent
}

// SetMatchFilter sets whether the node matches the filter
func (n *Node) SetMatchFilter(match bool) {
	n.MatchFilter = match
}

// IsMatchFilter returns whether the node matches the filter
func (n *Node) IsMatchFilter() bool {
	return n.MatchFilter
}

// SetExpanded sets whether the node is expanded
func (n *Node) SetExpanded(expanded bool) {
	n.Expanded = expanded
}
