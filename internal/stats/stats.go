package stats

import "github.com/cksidharthan/lazyjson/internal/model"

// Stats holds statistics about the JSON structure
type Stats struct {
	TotalNodes  int
	MaxDepth    int
	ObjectCount int
	ArrayCount  int
	ValueCount  int
}

// CalculateStats calculates statistics for a JSON tree
func CalculateStats(node *model.Node) Stats {
	stats := Stats{}
	calculateStatsRecursive(node, 0, &stats)
	return stats
}

func calculateStatsRecursive(node *model.Node, depth int, stats *Stats) {
	if node == nil {
		return
	}

	// Update total nodes
	stats.TotalNodes++

	// Update max depth
	if depth > stats.MaxDepth {
		stats.MaxDepth = depth
	}

	// Count by type
	switch node.Type {
	case "object":
		stats.ObjectCount++
	case "array":
		stats.ArrayCount++
	case "string", "number", "boolean", "null":
		stats.ValueCount++
	}

	// Process children
	for _, child := range node.Children {
		calculateStatsRecursive(child, depth+1, stats)
	}
}
