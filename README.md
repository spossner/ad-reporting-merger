# Ad Reporting Merger

A Go utility for merging ad reporting CSV files from the Downloads directory. Processes files with specific prefixes and consolidates them into unified CSV outputs with duplicate detection and chronological sorting.

## Usage

### Build and Run
```bash
go build -o ad-reporting-merger
./ad-reporting-merger
```

### Direct Execution
```bash
go run main.go
```

### Development Commands
```bash
go fmt ./...            # Format all Go files
go vet ./...            # Run static analysis
go mod tidy             # Clean up dependencies
go test ./...           # Run all tests
go test -v ./...        # Run tests with verbose output
go test ./... -cover    # Run tests with coverage report
```

## Configuration

The tool uses embedded JSON configuration in `internal/config/groups.json`. To modify settings, edit this file and rebuild the application.

### Default Configuration
- **Work Directory**: `~/Downloads`
- **Processing Groups**:
  - Files with "AdManager Reporting" prefix → merged into `raw.csv`
  - Files with "Revenue per AdUnit" prefix → merged into `raw-revenue.csv`

### Configuration Structure
```json
{
  "work_dir": "~/Downloads",
  "groups": [
    {
      "prefix": "AdManager Reporting",
      "output": "raw.csv"
    },
    {
      "prefix": "Revenue per AdUnit", 
      "output": "raw-revenue.csv"
    }
  ]
}
```

## Features

- **Duplicate Detection**: Uses MD5 hashing to identify and skip duplicate files
- **Chronological Sorting**: Orders files by date extracted from CSV content
- **Header Management**: Automatically skips headers when merging files
- **Automatic Cleanup**: Removes source files after successful merging
- **Error Handling**: Continues processing other groups if errors occur
- **Performance Monitoring**: Provides detailed timing and processing statistics