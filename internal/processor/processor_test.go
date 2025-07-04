package processor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spossner/ad-reporting-merger/internal/config"
	"github.com/spossner/ad-reporting-merger/internal/filesystem"
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
	if err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}

	err = os.WriteFile(file2, []byte(content2), 0644)
	if err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}

	// Setup processor
	fileOps, err := filesystem.NewFileOperations(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create file operations: %v", err)
	}

	processor := NewProcessor(fileOps)

	// Save current directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Change to test directory
	err = fileOps.ChangeToWorkDir()
	if err != nil {
		t.Fatalf("Failed to change to work dir: %v", err)
	}

	group := config.Group{
		Prefix: "AdManager Reporting",
		Output: "test-output.csv",
	}

	t.Run("successful processing", func(t *testing.T) {
		result := processor.ProcessGroup(group)
		
		if result.Error != nil {
			t.Errorf("Expected no error, got %v", result.Error)
		}

		if result.FilesFound != 2 {
			t.Errorf("Expected 2 files found, got %d", result.FilesFound)
		}

		if result.FilesMerged != 2 {
			t.Errorf("Expected 2 files merged, got %d", result.FilesMerged)
		}

		// Check output file exists
		outputPath := filepath.Join(tmpDir, "test-output.csv")
		_, err := os.Stat(outputPath)
		if err != nil {
			t.Errorf("Output file should exist: %v", err)
		}

		// Check source files were deleted
		_, err = os.Stat(file1)
		if err == nil {
			t.Error("Source file1 should have been deleted")
		}

		_, err = os.Stat(file2)
		if err == nil {
			t.Error("Source file2 should have been deleted")
		}
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
	if err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}

	err = os.WriteFile(file2, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}

	// Setup processor
	fileOps, err := filesystem.NewFileOperations(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create file operations: %v", err)
	}

	processor := NewProcessor(fileOps)

	// Save current directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Change to test directory
	err = fileOps.ChangeToWorkDir()
	if err != nil {
		t.Fatalf("Failed to change to work dir: %v", err)
	}

	group := config.Group{
		Prefix: "AdManager Reporting",
		Output: "test-output.csv",
	}

	t.Run("duplicate detection", func(t *testing.T) {
		result := processor.ProcessGroup(group)
		
		if result.Error == nil {
			t.Error("Expected error for duplicate files")
		}

		if result.FilesFound != 2 {
			t.Errorf("Expected 2 files found, got %d", result.FilesFound)
		}

		if result.FilesMerged != 0 {
			t.Errorf("Expected 0 files merged due to duplicates, got %d", result.FilesMerged)
		}
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
	if err != nil {
		t.Fatalf("Failed to create file operations: %v", err)
	}

	processor := NewProcessor(fileOps)

	// Save current directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Change to test directory
	err = fileOps.ChangeToWorkDir()
	if err != nil {
		t.Fatalf("Failed to change to work dir: %v", err)
	}

	groups := []config.Group{
		{Prefix: "AdManager Reporting", Output: "test1.csv"},
		{Prefix: "Revenue per AdUnit", Output: "test2.csv"},
	}

	t.Run("process all groups", func(t *testing.T) {
		results := processor.ProcessAllGroups(groups)
		
		if len(results) != 2 {
			t.Errorf("Expected 2 results, got %d", len(results))
		}

		// Both should have errors since no files exist
		for i, result := range results {
			if result.Error == nil {
				t.Errorf("Expected error for group %d (no files)", i)
			}
		}
	})
}