package ui

import (
	"fmt"
	"strings"

	"github.com/cksidharthan/lazyjson/internal/model"
	"github.com/cksidharthan/lazyjson/internal/stats"
	"github.com/jroimartin/gocui"
)

// Layout sets up the UI layout
func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Title view with stats
	if v, err := g.SetView("title", 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "JSON Viewer"
		v.FgColor = gocui.ColorYellow

		// Get JSON stats
		rootNode := model.GetRootNode()
		jsonStats := stats.CalculateStats(rootNode)
		statsStr := fmt.Sprintf("\033[36mNodes:\033[0m \033[33m%d\033[0m \033[36mDepth:\033[0m \033[33m%d\033[0m \033[36mObjects:\033[0m \033[33m%d\033[0m \033[36mArrays:\033[0m \033[33m%d\033[0m \033[36mValues:\033[0m \033[33m%d\033[0m",
			jsonStats.TotalNodes,
			jsonStats.MaxDepth,
			jsonStats.ObjectCount,
			jsonStats.ArrayCount,
			jsonStats.ValueCount)

		// Left side content
		leftContent := fmt.Sprintf("  \033[36mFile:\033[0m \033[33m%s\033[0m", model.GetFilePath())

		// Strip ANSI codes for length calculation
		cleanLeftContent := strings.ReplaceAll(leftContent, "\033[36m", "")
		cleanLeftContent = strings.ReplaceAll(cleanLeftContent, "\033[33m", "")
		cleanLeftContent = strings.ReplaceAll(cleanLeftContent, "\033[0m", "")

		cleanStatsStr := strings.ReplaceAll(statsStr, "\033[36m", "")
		cleanStatsStr = strings.ReplaceAll(cleanStatsStr, "\033[33m", "")
		cleanStatsStr = strings.ReplaceAll(cleanStatsStr, "\033[0m", "")

		// Calculate padding for right alignment
		padding := maxX - len(cleanLeftContent) - len(cleanStatsStr) - 3
		if padding < 1 {
			padding = 1
		}

		// Write content with right-aligned stats
		fmt.Fprintf(v, "%s%s%s",
			leftContent,
			strings.Repeat(" ", padding),
			statsStr)
	}

	// Filter view
	if v, err := g.SetView("filter", 0, 3, maxX-1, 5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Filter"
		v.Editable = true
		v.FgColor = gocui.ColorGreen
		v.Editor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
			filterEditor(g, v, key, ch, mod)
		})
	}

	// Main view for JSON tree
	if v, err := g.SetView("tree", 0, 6, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "JSON Tree"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		RenderTree(v)

		if _, err := g.SetCurrentView("tree"); err != nil {
			return err
		}
	}

	// Help view
	if v, err := g.SetView("help", 0, maxY-4, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Help"
		v.FgColor = gocui.ColorCyan

		// Navigation help with custom colors
		fmt.Fprintf(v, "  \033[33mNavigation\033[0m │ \033[36m↑/↓\033[0m: Move cursor   \033[36mEnter\033[0m: Expand/Collapse   \033[36me\033[0m: Expand all   \033[36mc\033[0m: Collapse all\n")
		fmt.Fprintf(v, "  \033[33mFilter\033[0m     │ \033[36m/\033[0m: Enable filter   \033[36mEsc\033[0m: Clear filter   \033[36mEnter\033[0m: Apply filter\n")
		fmt.Fprintf(v, "  \033[33mOther\033[0m      │ \033[36mq\033[0m: Quit")
	}

	return nil
}
