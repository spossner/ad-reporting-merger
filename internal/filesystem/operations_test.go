package filesystem

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileOperations(t *testing.T) {
	t.Run("with regular path", func(t *testing.T) {
		ops, err := NewFileOperations("/tmp")
		require.NoError(t, err)
		assert.Equal(t, "/tmp", ops.workDir)
	})

	t.Run("with home path", func(t *testing.T) {
		ops, err := NewFileOperations("~/test")
		require.NoError(t, err)
		home, _ := os.UserHomeDir()
		expected := filepath.Join(home, "test")
		assert.Equal(t, expected, ops.workDir)
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
		require.NoError(t, err, "Failed to create test file %s", file)
		f.Close()
	}

	// Test finding files
	ops, err := NewFileOperations(tmpDir)
	require.NoError(t, err)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err = ops.ChangeToWorkDir()
	require.NoError(t, err)

	found, err := ops.FindFiles("AdManager Reporting")
	require.NoError(t, err)

	assert.Len(t, found, 2)

	for _, file := range found {
		assert.True(t, strings.HasPrefix(file, "AdManager Reporting"), "File %s doesn't have expected prefix", file)
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
	require.NoError(t, err)

	// Copy files
	ops, err := NewFileOperations(destDir)
	require.NoError(t, err)

	err = ops.CopyTestFiles(sourceDir, destDir)
	require.NoError(t, err)

	// Verify file was copied
	copiedFile := filepath.Join(destDir, "test.csv")
	copiedContent, err := os.ReadFile(copiedFile)
	require.NoError(t, err)

	assert.Equal(t, content, string(copiedContent))

	// Verify original still exists
	_, err = os.Stat(testFile)
	assert.NoError(t, err, "Original file should still exist")
}