package main

import (
	"bufio"
	"crypto/md5"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Group struct {
	prefix string
	output string
}

var (
	groups = []Group{
		{
			prefix: "AdManager Reporting",
			output: "raw.csv",
		},
		{
			prefix: "Revenue per AdUnit",
			output: "raw-revenue.csv",
		},
	}
)

func main() {
	// Navigate to ~/Downloads
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Unable to determine user home directory")
	}
	err = os.Chdir(filepath.Join(home, "Downloads"))
	if err != nil {
		log.Fatalf("Unable to change to ~/Downloads: %v", err)
	}

	for _, group := range groups {
		fmt.Println("Processing group:", group.prefix)
		err := processPattern(group.prefix, group.output)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Merged group:", group.prefix, "->", group.output)
		deleteFiles(group.prefix)
	}
}

func processPattern(prefix, output string) error {
	files := findFiles(prefix)
	if hasDuplicateContent(files) {
		return fmt.Errorf("duplicate file content found in group: %s", prefix)
	}

	if len(files) == 0 {
		return fmt.Errorf("no files found for pattern: %s", prefix)
	}

	// Sort files by the first 10 characters (assumed date in yyyy-mm-dd format)
	sort.Slice(files, func(i, j int) bool {
		dateI := readFirstDate(files[i])
		dateJ := readFirstDate(files[j])
		return dateI < dateJ
	})

	mergeFiles(files, output)
	return nil
}

func findFiles(prefix string) []string {
	var matched []string
	entries, err := os.ReadDir(".")
	if err != nil {
		log.Fatalf("Unable to read directory: %v", err)
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), prefix) {
			matched = append(matched, e.Name())
		}
	}
	return matched
}

func hasDuplicateContent(files []string) bool {
	hashes := make(map[string]string) // contentHash -> filename
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Unable to read file: %v", err)
		}
		hash := fmt.Sprintf("%x", md5.Sum(content))
		if prev, exists := hashes[hash]; exists {
			fmt.Printf("Duplicate files: %s and %s\n", file, prev)
			return true
		}
		hashes[hash] = file
	}
	return false
}

func readFirstDate(file string) string {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	// Skip header row
	_, err = r.Read()
	if err != nil {
		return ""
	}

	row, err := r.Read()
	if err != nil || len(row) == 0 {
		return ""
	}
	return row[0] // date in the first column
}

func mergeFiles(files []string, outFile string) {
	out, err := os.Create(outFile)
	if err != nil {
		log.Fatalf("Unable to create output file: %v", err)
	}
	defer out.Close()

	writer := bufio.NewWriter(out)
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			log.Fatalf("Unable to open file: %v", err)
		}
		scanner := bufio.NewScanner(f)
		var line int
		for scanner.Scan() {
			line += 1
			if line == 1 {
				continue // skip header
			}
			row := scanner.Text()
			if line == 2 {
				fmt.Printf("%s\n", row[:10]) // print first 10 characters of the second row
			}
			writer.WriteString(row + "\n")
		}
		f.Close()
	}
	writer.Flush()
}

func deleteFiles(prefix string) {
	files := findFiles(prefix)
	for _, f := range files {
		err := os.Remove(f)
		if err != nil {
			fmt.Printf("Could not delete %s: %v\n", f, err)
		}
	}
}
