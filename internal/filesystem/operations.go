package filesystem

import (
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