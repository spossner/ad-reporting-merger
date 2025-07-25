package merger

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
)

type CSVMerger struct{}

func NewCSVMerger() *CSVMerger {
	return &CSVMerger{}
}

func (m *CSVMerger) MergeFiles(files []string, output string) ([]string, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files to merge")
	}

	// Sort files by the first date in each file
	sort.Slice(files, func(i, j int) bool {
		dateI := m.readFirstDate(files[i])
		dateJ := m.readFirstDate(files[j])
		return dateI < dateJ
	})

	out, err := os.Create(output)
	if err != nil {
		return nil, fmt.Errorf("unable to create output file: %w", err)
	}
	defer out.Close()

	writer := bufio.NewWriter(out)
	defer writer.Flush()

	dates := make([]string, 0, len(files))
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("unable to open file %s: %w", file, err)
		}

		scanner := bufio.NewScanner(f)
		var line int
		for scanner.Scan() {
			line++
			if line == 1 {
				continue // skip header
			}
			row := scanner.Text()
			if line == 2 {
				dates = append(dates, row[:10]) // track first 10 characters of the second row as date
			}
			writer.WriteString(row + "\n")
		}
		f.Close()

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading file %s: %w", file, err)
		}
	}

	return dates, nil
}

func (m *CSVMerger) readFirstDate(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
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
