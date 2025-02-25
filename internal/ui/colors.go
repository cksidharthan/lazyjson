package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

// ANSI color codes
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Gray      = "\033[90m"
	Red       = "\033[91m"
	Green     = "\033[92m"
	Yellow    = "\033[93m"
	Blue      = "\033[94m"
	Magenta   = "\033[95m"
	Cyan      = "\033[96m"
	White     = "\033[97m"
)

// ColorizeKey returns a colored key string
func ColorizeKey(v *gocui.View, key string) {
	fmt.Fprintf(v, "\033[36m%s\033[0m", key)
}

// ColorizeValue colors and writes a value based on its type
func ColorizeValue(v *gocui.View, value, valueType string) {
	switch valueType {
	case "string":
		fmt.Fprintf(v, "\033[32m%s\033[0m", value) // Green
	case "number":
		fmt.Fprintf(v, "\033[33m%s\033[0m", value) // Yellow
	case "boolean":
		fmt.Fprintf(v, "\033[35m%s\033[0m", value) // Magenta
	case "null":
		fmt.Fprintf(v, "\033[90m%s\033[0m", value) // Gray
	default:
		fmt.Fprintf(v, "\033[37m%s\033[0m", value) // White
	}
}

// ColorizeTreeSymbol colors and writes a tree symbol
func ColorizeTreeSymbol(v *gocui.View, symbol string) {
	fmt.Fprintf(v, "\033[90m%s\033[0m", symbol) // Gray
}

// ColorizeExpand colors and writes an expand/collapse indicator
func ColorizeExpand(v *gocui.View, symbol string) {
	fmt.Fprintf(v, "\033[34m%s\033[0m", symbol) // Blue
}

// ColorizeRoot colors and writes the root node text
func ColorizeRoot(v *gocui.View, text string) {
	fmt.Fprintf(v, "\033[1;37m%s\033[0m", text) // Bold White
}
