# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go utility for merging ad reporting CSV files from the Downloads directory. It processes files with specific prefixes and consolidates them into unified CSV outputs while handling duplicate detection and chronological sorting.

## Common Commands

### Build and Run
```bash
go build -o ad-reporting-merger
./ad-reporting-merger
```

### Development
```bash
go run main.go          # Run directly without building
go fmt ./...            # Format all Go files
go vet ./...            # Run static analysis
go mod tidy             # Clean up dependencies
```

### Testing
```bash
go test ./...           # Run all tests (if any exist)
```

## Code Architecture

The application is structured as a single-file Go program with these key components:

### Core Data Structure
- `Group` struct defines file processing groups with prefix patterns and output filenames
- Two predefined groups: "AdManager Reporting" → `raw.csv` and "Revenue per AdUnit" → `raw-revenue.csv`

### Processing Pipeline
1. **File Discovery**: `findFiles()` searches ~/Downloads for files matching group prefixes
2. **Duplicate Detection**: `hasDuplicateContent()` uses MD5 hashing to identify duplicate files
3. **Chronological Sorting**: Files are sorted by date extracted from the first data row
4. **Merging**: `mergeFiles()` combines CSV files, skipping headers from subsequent files
5. **Cleanup**: `deleteFiles()` removes source files after successful merging

### Key Functions
- `processPattern()`: Main processing logic for each group
- `readFirstDate()`: Extracts date from first CSV row for sorting
- `deleteFiles()`: Cleanup function that removes source files after successful merging

## File Processing Behavior

- Works exclusively in ~/Downloads directory
- Skips header rows when merging CSV files
- Prints first 10 characters of second row (date) for verification
- Handles errors gracefully and continues processing other groups
- Uses buffered I/O for efficient file operations
- **Automatically deletes source files after successful merging** - merged files are removed from Downloads

## Dependencies

The project uses only Go standard library packages:
- `bufio`, `crypto/md5`, `encoding/csv` for file processing
- `os`, `path/filepath` for file system operations
- `fmt`, `log` for output and error handling
- `sort`, `strings` for data manipulation