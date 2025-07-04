package detector

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHasDuplicates(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "detector_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	file1 := filepath.Join(tmpDir, "file1.csv")
	file2 := filepath.Join(tmpDir, "file2.csv")
	file3 := filepath.Join(tmpDir, "file3.csv")

	content1 := "Date,Value\n2025-01-01,100\n"
	content2 := "Date,Value\n2025-01-02,200\n"
	content3 := "Date,Value\n2025-01-01,100\n" // Same as file1

	err = os.WriteFile(file1, []byte(content1), 0644)
	require.NoError(t, err)

	err = os.WriteFile(file2, []byte(content2), 0644)
	require.NoError(t, err)

	err = os.WriteFile(file3, []byte(content3), 0644)
	require.NoError(t, err)

	detector := NewDuplicateDetector()

	t.Run("no duplicates", func(t *testing.T) {
		hasDuplicates, err := detector.HasDuplicates([]string{file1, file2})
		require.NoError(t, err)
		assert.False(t, hasDuplicates, "Expected no duplicates")
	})

	t.Run("with duplicates", func(t *testing.T) {
		hasDuplicates, err := detector.HasDuplicates([]string{file1, file2, file3})
		require.NoError(t, err)
		assert.True(t, hasDuplicates, "Expected duplicates")
	})

	t.Run("single file", func(t *testing.T) {
		hasDuplicates, err := detector.HasDuplicates([]string{file1})
		require.NoError(t, err)
		assert.False(t, hasDuplicates, "Expected no duplicates for single file")
	})

	t.Run("nonexistent file", func(t *testing.T) {
		_, err := detector.HasDuplicates([]string{"nonexistent.csv"})
		assert.Error(t, err, "Expected error for nonexistent file")
	})
}