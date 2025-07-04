package filesystem

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewFileOperations(t *testing.T) {
	t.Run("with regular path", func(t *testing.T) {
		ops, err := NewFileOperations("/tmp")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if ops.workDir != "/tmp" {
			t.Errorf("Expected workDir to be /tmp, got %s", ops.workDir)
		}
	})

	t.Run("with home path", func(t *testing.T) {
		ops, err := NewFileOperations("~/test")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		home, _ := os.UserHomeDir()
		expected := filepath.Join(home, "test")
		if ops.workDir != expected {
			t.Errorf("Expected workDir to be %s, got %s", expected, ops.workDir)
		}
	})
}

func TestFindFiles(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "filesystem_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	testFiles := []string{
		"AdManager Reporting_2025-01-01.csv",
		"AdManager Reporting_2025-01-02.csv",
		"Revenue per AdUnit_2025-01-01.csv",
		"other_file.txt",
	}

	for _, file := range testFiles {
		f, err := os.Create(filepath.Join(tmpDir, file))
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
		f.Close()
	}

	// Test finding files
	ops, err := NewFileOperations(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create FileOperations: %v", err)
	}

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err = ops.ChangeToWorkDir()
	if err != nil {
		t.Fatalf("Failed to change to work dir: %v", err)
	}

	found, err := ops.FindFiles("AdManager Reporting")
	if err != nil {
		t.Fatalf("Failed to find files: %v", err)
	}

	if len(found) != 2 {
		t.Errorf("Expected 2 files, got %d", len(found))
	}

	for _, file := range found {
		if !strings.HasPrefix(file, "AdManager Reporting") {
			t.Errorf("File %s doesn't have expected prefix", file)
		}
	}
}

func TestCopyTestFiles(t *testing.T) {
	// Create temporary directories
	sourceDir, err := os.MkdirTemp("", "source_test")
	if err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	defer os.RemoveAll(sourceDir)

	destDir, err := os.MkdirTemp("", "dest_test")
	if err != nil {
		t.Fatalf("Failed to create dest dir: %v", err)
	}
	defer os.RemoveAll(destDir)

	// Create test file in source
	testFile := filepath.Join(sourceDir, "test.csv")
	content := "Date,Value\n2025-01-01,100\n"
	err = os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Copy files
	ops, err := NewFileOperations(destDir)
	if err != nil {
		t.Fatalf("Failed to create FileOperations: %v", err)
	}

	err = ops.CopyTestFiles(sourceDir, destDir)
	if err != nil {
		t.Fatalf("Failed to copy test files: %v", err)
	}

	// Verify file was copied
	copiedFile := filepath.Join(destDir, "test.csv")
	copiedContent, err := os.ReadFile(copiedFile)
	if err != nil {
		t.Fatalf("Failed to read copied file: %v", err)
	}

	if string(copiedContent) != content {
		t.Errorf("Copied content doesn't match original")
	}

	// Verify original still exists
	_, err = os.Stat(testFile)
	if err != nil {
		t.Errorf("Original file was deleted: %v", err)
	}
}