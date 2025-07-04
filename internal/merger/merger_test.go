package merger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMergeFiles(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "merger_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	file1 := filepath.Join(tmpDir, "file1.csv")
	file2 := filepath.Join(tmpDir, "file2.csv")
	file3 := filepath.Join(tmpDir, "file3.csv")
	output := filepath.Join(tmpDir, "output.csv")

	content1 := "Date,Value\n2025-01-01,100\n2025-01-01,150\n"
	content2 := "Date,Value\n2025-01-02,200\n2025-01-02,250\n"
	content3 := "Date,Value\n2025-01-03,300\n2025-01-03,350\n"

	err = os.WriteFile(file1, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}

	err = os.WriteFile(file2, []byte(content2), 0644)
	if err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}

	err = os.WriteFile(file3, []byte(content3), 0644)
	if err != nil {
		t.Fatalf("Failed to create file3: %v", err)
	}

	merger := NewCSVMerger()

	t.Run("merge files", func(t *testing.T) {
		err = merger.MergeFiles([]string{file1, file2, file3}, output)
		if err != nil {
			t.Fatalf("Failed to merge files: %v", err)
		}

		// Read output file
		outputContent, err := os.ReadFile(output)
		if err != nil {
			t.Fatalf("Failed to read output file: %v", err)
		}

		outputStr := string(outputContent)
		lines := strings.Split(strings.TrimSpace(outputStr), "\n")

		// Should have 6 lines (2 from each file, headers skipped)
		if len(lines) != 6 {
			t.Errorf("Expected 6 lines, got %d", len(lines))
		}

		// Check that data is in chronological order
		if !strings.HasPrefix(lines[0], "2025-01-01") {
			t.Errorf("Expected first line to start with 2025-01-01, got %s", lines[0])
		}

		if !strings.HasPrefix(lines[4], "2025-01-03") {
			t.Errorf("Expected fifth line to start with 2025-01-03, got %s", lines[4])
		}
	})

	t.Run("empty file list", func(t *testing.T) {
		err = merger.MergeFiles([]string{}, output)
		if err == nil {
			t.Error("Expected error for empty file list")
		}
	})

	t.Run("nonexistent file", func(t *testing.T) {
		err = merger.MergeFiles([]string{"nonexistent.csv"}, output)
		if err == nil {
			t.Error("Expected error for nonexistent file")
		}
	})
}

func TestReadFirstDate(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "merger_date_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test file
	testFile := filepath.Join(tmpDir, "test.csv")
	content := "Date,Value\n2025-01-01,100\n2025-01-02,200\n"
	err = os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	merger := NewCSVMerger()

	t.Run("read first date", func(t *testing.T) {
		date := merger.readFirstDate(testFile)
		if date != "2025-01-01" {
			t.Errorf("Expected first date to be 2025-01-01, got %s", date)
		}
	})

	t.Run("nonexistent file", func(t *testing.T) {
		date := merger.readFirstDate("nonexistent.csv")
		if date != "" {
			t.Errorf("Expected empty string for nonexistent file, got %s", date)
		}
	})
}