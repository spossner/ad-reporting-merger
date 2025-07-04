package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spossner/ad-reporting-merger/internal/config"
	"github.com/spossner/ad-reporting-merger/internal/filesystem"
	"github.com/spossner/ad-reporting-merger/internal/processor"
)

func TestIntegration(t *testing.T) {
	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Setup file operations
	fileOps, err := filesystem.NewFileOperations(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create file operations: %v", err)
	}

	// Get absolute path to testdata
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	testDataSource := filepath.Join(wd, "testdata", "source")

	// Copy test files to temp directory
	err = fileOps.CopyTestFiles(testDataSource, tmpDir)
	if err != nil {
		t.Fatalf("Failed to copy test files: %v", err)
	}

	// Save current directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Change to test directory
	err = fileOps.ChangeToWorkDir()
	if err != nil {
		t.Fatalf("Failed to change to work dir: %v", err)
	}

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Create processor
	proc := processor.NewProcessor(fileOps)

	// Process all groups
	results := proc.ProcessAllGroups(cfg.GetGroups())

	// Verify results
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	for _, result := range results {
		if result.Error != nil {
			t.Errorf("Expected no error for group %s, got %v", result.Group.Prefix, result.Error)
		}

		if result.FilesFound != 3 {
			t.Errorf("Expected 3 files found for group %s, got %d", result.Group.Prefix, result.FilesFound)
		}

		if result.FilesMerged != 3 {
			t.Errorf("Expected 3 files merged for group %s, got %d", result.Group.Prefix, result.FilesMerged)
		}

		// Check output file exists
		outputPath := filepath.Join(tmpDir, result.OutputFile)
		_, err := os.Stat(outputPath)
		if err != nil {
			t.Errorf("Output file %s should exist: %v", result.OutputFile, err)
		}

		// Verify output content
		content, err := os.ReadFile(outputPath)
		if err != nil {
			t.Errorf("Failed to read output file %s: %v", result.OutputFile, err)
			continue
		}

		contentStr := string(content)
		lines := strings.Split(strings.TrimSpace(contentStr), "\n")

		// Should have 9 lines (3 files × 3 lines each, minus headers)
		if len(lines) != 9 {
			t.Errorf("Expected 9 lines in output file %s, got %d", result.OutputFile, len(lines))
		}

		// Check chronological order
		if !strings.HasPrefix(lines[0], "2025-01-01") {
			t.Errorf("Expected first line to start with 2025-01-01 in %s", result.OutputFile)
		}

		if !strings.HasPrefix(lines[8], "2025-01-03") {
			t.Errorf("Expected last line to start with 2025-01-03 in %s", result.OutputFile)
		}
	}

	// Verify source files were deleted
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read temp dir: %v", err)
	}

	sourceFileCount := 0
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "AdManager Reporting") || 
		   strings.HasPrefix(entry.Name(), "Revenue per AdUnit") {
			sourceFileCount++
		}
	}

	if sourceFileCount != 0 {
		t.Errorf("Expected all source files to be deleted, but found %d", sourceFileCount)
	}
}

func TestIntegrationWithTestSetup(t *testing.T) {
	// This test demonstrates how to run the application with test data
	// without deleting the original test files

	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "integration_setup_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Setup file operations
	fileOps, err := filesystem.NewFileOperations(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create file operations: %v", err)
	}

	// Get absolute path to testdata
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	testDataSource := filepath.Join(wd, "testdata", "source")

	// Copy test files to temp directory (preserving originals)
	err = fileOps.CopyTestFiles(testDataSource, tmpDir)
	if err != nil {
		t.Fatalf("Failed to copy test files: %v", err)
	}

	// Verify original test files still exist
	originalFiles, err := os.ReadDir(testDataSource)
	if err != nil {
		t.Fatalf("Failed to read test data source: %v", err)
	}

	if len(originalFiles) == 0 {
		t.Error("Original test files should still exist")
	}

	// Verify copied files exist
	copiedFiles, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read temp dir: %v", err)
	}

	if len(copiedFiles) != len(originalFiles) {
		t.Errorf("Expected %d copied files, got %d", len(originalFiles), len(copiedFiles))
	}

	t.Logf("Successfully copied %d test files to %s", len(copiedFiles), tmpDir)
	t.Logf("Original test files preserved in %s", testDataSource)
}