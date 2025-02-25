package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cksidharthan/lazyjson/internal/model"
	"github.com/cksidharthan/lazyjson/internal/parser"
	"github.com/cksidharthan/lazyjson/internal/ui"

	"github.com/jroimartin/gocui"
)

func Start() {
	if len(os.Args) < 2 {
		fmt.Println("Error: No JSON file specified")
		fmt.Println("Usage: lazyjson <json-file>")
		fmt.Println("\nExample: lazyjson data.json")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' does not exist\n", filename)
		os.Exit(1)
	}

	// Parse JSON file
	data, err := parser.LoadJSON(filename)
	if err != nil {
		log.Fatalf("Error loading JSON: %v", err)
	}

	// Get absolute path
	absPath, err := filepath.Abs(filename)
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}
	model.SetFilePath(absPath)

	// Build tree model
	rootNode := model.BuildTree("root", data, nil, "")
	rootNode.Expanded = true
	model.SetRootNode(rootNode)

	// Initialize GUI
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalf("Failed to create GUI: %v", err)
	}
	defer g.Close()

	// Setup UI components
	g.SetManagerFunc(ui.Layout)
	if err := ui.SetupKeybindings(g); err != nil {
		log.Fatalf("Failed to set keybindings: %v", err)
	}

	g.Cursor = true
	g.Mouse = true

	// Start main loop
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalf("Main loop error: %v", err)
	}
}
