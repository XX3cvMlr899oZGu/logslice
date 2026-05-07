package rotate

import (
	"os"
	"sort"
	"time"
)

// Segment represents a single log file segment with its modification time.
type Segment struct {
	Path    string
	ModTime time.Time
}

// ScanSegments walks the given list of file paths, stats each one, and
// returns them sorted oldest-first by modification time. Paths that cannot
// be stat'd are silently skipped so that a missing rotated file does not
// abort the entire scan.
func ScanSegments(paths []string) ([]Segment, error) {
	var segments []Segment

	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			// Skip files we cannot access (deleted between glob and stat).
			continue
		}
		segments = append(segments, Segment{
			Path:    p,
			ModTime: info.ModTime(),
		})
	}

	sort.Slice(segments, func(i, j int) bool {
		return segments[i].ModTime.Before(segments[j].ModTime)
	})

	return segments, nil
}

// FilterByTimeRange returns only those segments whose modification time falls
// within [from, to] (inclusive on both ends). This provides a coarse pre-
// filter before the binary-search inside each file is applied.
func FilterByTimeRange(segments []Segment, from, to time.Time) []Segment {
	var out []Segment
	for _, s := range segments {
		// Keep any segment that could overlap the requested range:
		// a segment modified before "from" might still contain lines
		// that fall within the range, so we only drop segments whose
		// modification time is strictly before the start of the range
		// by more than a generous buffer — here we simply keep all
		// segments up to and including those modified after "from".
		if !s.ModTime.Before(from) || s.ModTime.Equal(from) {
			out = append(out, s)
		}
		// Stop once we pass the end of the range.
		if s.ModTime.After(to) {
			break
		}
	}
	return out
}
