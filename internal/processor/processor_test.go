package processor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spossner/ad-reporting-merger/internal/config"
	"github.com/spossner/ad-reporting-merger/internal/filesystem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessGroup(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "processor_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	file1 := filepath.Join(tmpDir, "AdManager Reporting_2025-01-01.csv")
	file2 := filepath.Join(tmpDir, "AdManager Reporting_2025-01-02.csv")
	
	content1 := "Date,Value\n2025-01-01,100\n2025-01-01,150\n"
	content2 := "Date,Value\n2025-01-02,200\n2025-01-02,250\n"

	err = os.WriteFile(file1, []byte(content1), 0644)
	require.NoError(t, err)

	err = os.WriteFile(file2, []byte(content2), 0644)
	require.NoError(t, err)

	// Setup processor
	fileOps, err := filesystem.NewFileOperations(tmpDir)
	require.NoError(t, err)

	processor := NewProcessor(fileOps)

	// Save current directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Change to test directory
	err = fileOps.ChangeToWorkDir()
	require.NoError(t, err)

	group := config.Group{
		Prefix: "AdManager Reporting",
		Output: "test-output.csv",
	}

	t.Run("successful processing", func(t *testing.T) {
		result := processor.ProcessGroup(group)
		
		assert.NoError(t, result.Error)
		assert.Equal(t, 2, result.FilesFound)
		assert.Equal(t, 2, result.FilesMerged)

		// Check output file exists
		outputPath := filepath.Join(tmpDir, "test-output.csv")
		_, err := os.Stat(outputPath)
		assert.NoError(t, err, "Output file should exist")

		// Check source files were deleted
		_, err = os.Stat(file1)
		assert.Error(t, err, "Source file1 should have been deleted")

		_, err = os.Stat(file2)
		assert.Error(t, err, "Source file2 should have been deleted")
	})
}

func TestProcessGroupNoDuplicates(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "processor_dup_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create duplicate files
	file1 := filepath.Join(tmpDir, "AdManager Reporting_2025-01-01.csv")
	file2 := filepath.Join(tmpDir, "AdManager Reporting_2025-01-02.csv")
	
	content := "Date,Value\n2025-01-01,100\n" // Same content

	err = os.WriteFile(file1, []byte(content), 0644)
	require.NoError(t, err)

	err = os.WriteFile(file2, []byte(content), 0644)
	require.NoError(t, err)

	// Setup processor
	fileOps, err := filesystem.NewFileOperations(tmpDir)
	require.NoError(t, err)

	processor := NewProcessor(fileOps)

	// Save current directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Change to test directory
	err = fileOps.ChangeToWorkDir()
	require.NoError(t, err)

	group := config.Group{
		Prefix: "AdManager Reporting",
		Output: "test-output.csv",
	}

	t.Run("duplicate detection", func(t *testing.T) {
		result := processor.ProcessGroup(group)
		
		assert.Error(t, result.Error, "Expected error for duplicate files")
		assert.Equal(t, 2, result.FilesFound)
		assert.Equal(t, 0, result.FilesMerged, "Expected 0 files merged due to duplicates")
	})
}

func TestProcessAllGroups(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "processor_all_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Setup processor
	fileOps, err := filesystem.NewFileOperations(tmpDir)
	require.NoError(t, err)

	processor := NewProcessor(fileOps)

	// Save current directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Change to test directory
	err = fileOps.ChangeToWorkDir()
	require.NoError(t, err)

	groups := []config.Group{
		{Prefix: "AdManager Reporting", Output: "test1.csv"},
		{Prefix: "Revenue per AdUnit", Output: "test2.csv"},
	}

	t.Run("process all groups", func(t *testing.T) {
		results := processor.ProcessAllGroups(groups)
		
		assert.Len(t, results, 2)

		// Both should have errors since no files exist
		for i, result := range results {
			assert.Error(t, result.Error, "Expected error for group %d (no files)", i)
		}
	})
}