// Package rotate provides support for detecting and handling rotated or
// rolled log files, allowing logslice to transparently read across file
// boundaries when a log directory contains sequentially numbered or
// date-stamped log segments.
package rotate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Segment represents a single log file segment on disk.
type Segment struct {
	Path  string
	Index int // relative order among siblings; 0 = most recent
}

// FindSegments scans dir for files whose base names start with prefix and
// returns them ordered from oldest to newest (ascending index).
// The primary file (exactly equal to prefix) is treated as index 0 and
// placed last so callers process it after the rotated copies.
func FindSegments(dir, prefix string) ([]Segment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("rotate: read dir %q: %w", dir, err)
	}

	var segments []Segment
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if name == prefix {
			segments = append(segments, Segment{Path: filepath.Join(dir, name), Index: 0})
			continue
		}
		if strings.HasPrefix(name, prefix+".") || strings.HasPrefix(name, prefix+"-") {
			segments = append(segments, Segment{Path: filepath.Join(dir, name), Index: len(segments) + 1})
		}
	}

	if len(segments) == 0 {
		return nil, fmt.Errorf("rotate: no segments found for prefix %q in %q", prefix, dir)
	}

	// Sort so rotated files (index > 0) come first, ordered by path name
	// (lexicographic), then the live file (index 0) last.
	sort.Slice(segments, func(i, j int) bool {
		if segments[i].Index == 0 {
			return false
		}
		if segments[j].Index == 0 {
			return true
		}
		return segments[i].Path < segments[j].Path
	})

	// Re-assign contiguous indices after sort.
	for i := range segments {
		segments[i].Index = i
	}

	return segments, nil
}
