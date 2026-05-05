// Package slicer combines the reader and parser packages to extract
// log lines within a specific time range from a structured log file.
package slicer

import (
	"fmt"
	"io"
	"time"

	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/reader"
)

// Options configures the behaviour of a Slice operation.
type Options struct {
	// From is the inclusive start of the time range.
	From time.Time
	// To is the inclusive end of the time range.
	To time.Time
}

// Result holds metadata about a completed slice operation.
type Result struct {
	// LinesWritten is the number of log lines copied to the writer.
	LinesWritten int
}

// Slice opens the log file at path, seeks to the first line whose timestamp
// is >= opts.From using a binary search, then streams lines to w until a
// line whose timestamp exceeds opts.To is encountered.
func Slice(path string, opts Options, w io.Writer) (Result, error) {
	if opts.From.IsZero() || opts.To.IsZero() {
		return Result{}, fmt.Errorf("slicer: From and To must both be set")
	}
	if opts.To.Before(opts.From) {
		return Result{}, fmt.Errorf("slicer: To must not be before From")
	}

	offset, err := reader.SeekToTime(path, opts.From)
	if err != nil {
		return Result{}, fmt.Errorf("slicer: seek failed: %w", err)
	}

	lr, err := reader.NewLineReader(path)
	if err != nil {
		return Result{}, fmt.Errorf("slicer: open failed: %w", err)
	}
	defer lr.Close()

	if err := lr.Seek(offset); err != nil {
		return Result{}, fmt.Errorf("slicer: seek to offset %d failed: %w", offset, err)
	}

	var result Result
	for {
		line, err := lr.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return result, fmt.Errorf("slicer: read error: %w", err)
		}

		ts, parseErr := parser.ParseTimestamp(line)
		if parseErr != nil {
			// Skip lines we cannot parse a timestamp from.
			continue
		}

		if ts.After(opts.To) {
			break
		}

		if _, writeErr := fmt.Fprintln(w, line); writeErr != nil {
			return result, fmt.Errorf("slicer: write error: %w", writeErr)
		}
		result.LinesWritten++
	}

	return result, nil
}
