package interfaces

// Node represents a node in the JSON tree
type Node interface {
	// GetKey returns the node's key
	GetKey() string
	
	// GetValue returns the node's value
	GetValue() any
	
	// GetType returns the node's type
	GetType() string
	
	// GetChildren returns the node's children
	GetChildren() []Node
	
	// GetParent returns the node's parent
	GetParent() Node
	
	// SetMatchFilter sets whether the node matches the filter
	SetMatchFilter(match bool)
	
	// IsMatchFilter returns whether the node matches the filter
	IsMatchFilter() bool
	
	// SetExpanded sets whether the node is expanded
	SetExpanded(expanded bool)
	
	// IsExpanded returns whether the node is expanded
	IsExpanded() bool
}
