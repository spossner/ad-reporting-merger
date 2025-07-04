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

The application uses a modular package structure with clear separation of concerns:

### Package Structure
```
internal/
├── config/          # Configuration management with embedded JSON
├── processor/       # Main processing orchestration
├── merger/          # CSV merging functionality
├── detector/        # Duplicate detection logic
└── filesystem/      # File system operations
```

### Core Components
- **Config Package**: Manages group definitions via embedded JSON configuration
- **Processor Package**: Orchestrates the complete processing pipeline
- **Merger Package**: Handles CSV file merging with chronological sorting
- **Detector Package**: Implements duplicate detection using MD5 hashing
- **Filesystem Package**: Abstracts file operations and path handling

### Processing Pipeline
1. **Configuration Loading**: Embedded JSON config loaded via `go:embed`
2. **File Discovery**: Find files matching group prefixes in work directory
3. **Duplicate Detection**: MD5 hash comparison to identify duplicate files
4. **Chronological Sorting**: Files sorted by date extracted from first CSV row
5. **Merging**: CSV files combined, skipping headers from subsequent files
6. **Cleanup**: Source files removed after successful merging

### Key Types
- `config.Group`: Defines file processing groups with prefix patterns and output filenames
- `processor.ProcessingResult`: Contains detailed results including timing and error information
- Interfaces for testability: Each component can be easily mocked for unit testing

## File Processing Behavior

- Works in configurable directory (defaults to ~/Downloads)
- Skips header rows when merging CSV files
- Prints first 10 characters of second row (date) for verification
- Handles errors gracefully and continues processing other groups
- Uses buffered I/O for efficient file operations
- **Automatically deletes source files after successful merging** - merged files are removed from work directory
- Provides detailed processing results including timing and error information

## Configuration

Groups are defined in `internal/config/groups.json` and embedded at compile time:
- **AdManager Reporting** → `raw.csv`
- **Revenue per AdUnit** → `raw-revenue.csv`

To modify groups, edit `internal/config/groups.json` and rebuild the application.

## Dependencies

The project uses only Go standard library packages:
- `bufio`, `crypto/md5`, `encoding/csv` for file processing
- `os`, `path/filepath` for file system operations
- `fmt`, `log` for output and error handling
- `sort`, `strings` for data manipulation