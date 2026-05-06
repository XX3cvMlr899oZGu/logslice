package rotate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Segment represents a single rotated log file with its modification time.
type Segment struct {
	Path    string
	ModTime time.Time
}

// FindSegments scans the directory containing baseFile and returns all files
// matching the glob pattern derived from baseFile's base name (e.g. app.log,
// app.log.1, app.log.2024-01-01, …). Results are ordered oldest-first by
// modification time so callers can process them in chronological order.
func FindSegments(baseFile string) ([]Segment, error) {
	dir := filepath.Dir(baseFile)
	base := filepath.Base(baseFile)

	pattern := filepath.Join(dir, base+"*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("rotate: glob %q: %w", pattern, err)
	}

	// Also include the base file itself if it exists.
	matches = dedupe(append(matches, baseFile))

	var segments []Segment
	for _, m := range matches {
		info, err := os.Stat(m)
		if err != nil {
			// File may have been rotated away between glob and stat; skip it.
			continue
		}
		segments = append(segments, Segment{
			Path:    m,
			ModTime: info.ModTime(),
		})
	}

	sort.Slice(segments, func(i, j int) bool {
		return segments[i].ModTime.Before(segments[j].ModTime)
	})

	return segments, nil
}

// dedupe returns a slice with duplicate strings removed while preserving order.
func dedupe(ss []string) []string {
	seen := make(map[string]struct{}, len(ss))
	out := ss[:0]
	for _, s := range ss {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	return out
}
