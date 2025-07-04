package detector

import (
	"crypto/md5"
	"fmt"
	"os"
)

type DuplicateDetector struct{}

func NewDuplicateDetector() *DuplicateDetector {
	return &DuplicateDetector{}
}

func (d *DuplicateDetector) HasDuplicates(files []string) (bool, error) {
	hashes := make(map[string]string) // contentHash -> filename
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return false, fmt.Errorf("unable to read file %s: %w", file, err)
		}
		hash := fmt.Sprintf("%x", md5.Sum(content))
		if prev, exists := hashes[hash]; exists {
			fmt.Printf("Duplicate files: %s and %s\n", file, prev)
			return true, nil
		}
		hashes[hash] = file
	}
	return false, nil
}