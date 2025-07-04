package filesystem

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileOperations struct {
	workDir string
}

func NewFileOperations(workDir string) (*FileOperations, error) {
	expandedDir, err := expandPath(workDir)
	if err != nil {
		return nil, err
	}
	return &FileOperations{workDir: expandedDir}, nil
}

func (f *FileOperations) ChangeToWorkDir() error {
	return os.Chdir(f.workDir)
}

func (f *FileOperations) FindFiles(prefix string) ([]string, error) {
	var matched []string
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), prefix) {
			matched = append(matched, e.Name())
		}
	}
	return matched, nil
}

func (f *FileOperations) DeleteFiles(files []string) error {
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			return err
		}
	}
	return nil
}

// CopyTestFiles copies test files from source to destination, preserving originals
func (f *FileOperations) CopyTestFiles(sourceDir, destDir string) error {
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		sourcePath := filepath.Join(sourceDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		err := f.copyFile(sourcePath, destPath)
		if err != nil {
			return fmt.Errorf("failed to copy file %s: %w", entry.Name(), err)
		}
	}

	return nil
}

func (f *FileOperations) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}