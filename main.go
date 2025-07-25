package main

import (
	"fmt"
	"log"

	"github.com/spossner/ad-reporting-merger/internal/config"
	"github.com/spossner/ad-reporting-merger/internal/filesystem"
	"github.com/spossner/ad-reporting-merger/internal/processor"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize file operations
	fileOps, err := filesystem.NewFileOperations(cfg.GetWorkDir())
	if err != nil {
		log.Fatalf("Failed to initialize file operations: %v", err)
	}

	// Change to work directory
	err = fileOps.ChangeToWorkDir()
	if err != nil {
		log.Fatalf("Failed to change to work directory: %v", err)
	}

	// Initialize processor
	proc := processor.NewProcessor(fileOps)

	// Process all groups
	results := proc.ProcessAllGroups(cfg.GetGroups())

	// Display results
	for _, result := range results {
		fmt.Printf("Processing group: %s\n", result.Group.Prefix)
		if result.Error != nil {
			fmt.Printf("Error: %v\n", result.Error)
			continue
		}
		for _, date := range result.DatesFound {
			fmt.Printf("  %s\n", date)
		}
		fmt.Printf("Merged group: %s -> %s (Duration: %v)\n",
			result.Group.Prefix, result.OutputFile, result.Duration)
	}
}
