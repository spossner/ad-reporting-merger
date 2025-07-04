package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spossner/ad-reporting-merger/internal/config"
	"github.com/spossner/ad-reporting-merger/internal/filesystem"
	"github.com/spossner/ad-reporting-merger/internal/processor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)

	// Get absolute path to testdata
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	testDataSource := filepath.Join(wd, "testdata", "source")

	// Copy test files to temp directory
	err = fileOps.CopyTestFiles(testDataSource, tmpDir)
	require.NoError(t, err)

	// Save current directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Change to test directory
	err = fileOps.ChangeToWorkDir()
	require.NoError(t, err)

	// Load config
	cfg, err := config.LoadConfig()
	require.NoError(t, err)

	// Create processor
	proc := processor.NewProcessor(fileOps)

	// Process all groups
	results := proc.ProcessAllGroups(cfg.GetGroups())

	// Verify results
	assert.Len(t, results, 2)

	for _, result := range results {
		assert.NoError(t, result.Error, "Expected no error for group %s", result.Group.Prefix)
		assert.Equal(t, 3, result.FilesFound, "Expected 3 files found for group %s", result.Group.Prefix)
		assert.Equal(t, 3, result.FilesMerged, "Expected 3 files merged for group %s", result.Group.Prefix)

		// Check output file exists
		outputPath := filepath.Join(tmpDir, result.OutputFile)
		_, err := os.Stat(outputPath)
		assert.NoError(t, err, "Output file %s should exist", result.OutputFile)

		// Verify output content
		content, err := os.ReadFile(outputPath)
		if !assert.NoError(t, err, "Failed to read output file %s", result.OutputFile) {
			continue
		}

		contentStr := string(content)
		lines := strings.Split(strings.TrimSpace(contentStr), "\n")

		// Should have 9 lines (3 files × 3 lines each, minus headers)
		assert.Len(t, lines, 9, "Expected 9 lines in output file %s", result.OutputFile)

		// Check chronological order
		assert.True(t, strings.HasPrefix(lines[0], "2025-01-01"), "Expected first line to start with 2025-01-01 in %s", result.OutputFile)
		assert.True(t, strings.HasPrefix(lines[8], "2025-01-03"), "Expected last line to start with 2025-01-03 in %s", result.OutputFile)
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

	assert.Equal(t, 0, sourceFileCount, "Expected all source files to be deleted")
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
	require.NoError(t, err)

	// Get absolute path to testdata
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	testDataSource := filepath.Join(wd, "testdata", "source")

	// Copy test files to temp directory (preserving originals)
	err = fileOps.CopyTestFiles(testDataSource, tmpDir)
	require.NoError(t, err)

	// Verify original test files still exist
	originalFiles, err := os.ReadDir(testDataSource)
	require.NoError(t, err)
	assert.NotEmpty(t, originalFiles, "Original test files should still exist")

	// Verify copied files exist
	copiedFiles, err := os.ReadDir(tmpDir)
	require.NoError(t, err)
	assert.Len(t, copiedFiles, len(originalFiles), "Expected same number of copied files as originals")

	t.Logf("Successfully copied %d test files to %s", len(copiedFiles), tmpDir)
	t.Logf("Original test files preserved in %s", testDataSource)
}