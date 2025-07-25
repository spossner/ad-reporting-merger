package processor

import (
	"fmt"
	"time"

	"github.com/spossner/ad-reporting-merger/internal/config"
	"github.com/spossner/ad-reporting-merger/internal/detector"
	"github.com/spossner/ad-reporting-merger/internal/filesystem"
	"github.com/spossner/ad-reporting-merger/internal/merger"
)

type ProcessingResult struct {
	Group       config.Group
	FilesFound  int
	FilesMerged int
	DatesFound  []string
	OutputFile  string
	Duration    time.Duration
	Error       error
}

type Processor struct {
	fileOps  *filesystem.FileOperations
	detector *detector.DuplicateDetector
	merger   *merger.CSVMerger
}

func NewProcessor(fileOps *filesystem.FileOperations) *Processor {
	return &Processor{
		fileOps:  fileOps,
		detector: detector.NewDuplicateDetector(),
		merger:   merger.NewCSVMerger(),
	}
}

func (p *Processor) ProcessGroup(group config.Group) *ProcessingResult {
	start := time.Now()
	result := &ProcessingResult{
		Group:      group,
		OutputFile: group.Output,
	}

	files, err := p.fileOps.FindFiles(group.Prefix)
	if err != nil {
		result.Error = fmt.Errorf("failed to find files: %w", err)
		result.Duration = time.Since(start)
		return result
	}

	result.FilesFound = len(files)

	if len(files) == 0 {
		result.Error = fmt.Errorf("no files found for pattern: %s", group.Prefix)
		result.Duration = time.Since(start)
		return result
	}

	hasDuplicates, err := p.detector.HasDuplicates(files)
	if err != nil {
		result.Error = fmt.Errorf("failed to check duplicates: %w", err)
		result.Duration = time.Since(start)
		return result
	}

	if hasDuplicates {
		result.Error = fmt.Errorf("duplicate file content found in group: %s", group.Prefix)
		result.Duration = time.Since(start)
		return result
	}

	dates, err := p.merger.MergeFiles(files, group.Output)
	if err != nil {
		result.Error = fmt.Errorf("failed to merge files: %w", err)
		result.Duration = time.Since(start)
		return result
	}

	result.FilesMerged = len(files)
	result.DatesFound = dates

	// Clean up source files
	err = p.fileOps.DeleteFiles(files)
	if err != nil {
		result.Error = fmt.Errorf("failed to delete source files: %w", err)
		result.Duration = time.Since(start)
		return result
	}

	result.Duration = time.Since(start)
	return result
}

func (p *Processor) ProcessAllGroups(groups []config.Group) []*ProcessingResult {
	results := make([]*ProcessingResult, len(groups))
	for i, group := range groups {
		results[i] = p.ProcessGroup(group)
	}
	return results
}
