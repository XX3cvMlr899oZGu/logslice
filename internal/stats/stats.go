// Package stats provides collection and reporting of log slicing statistics.
package stats

import (
	"fmt"
	"io"
	"time"
)

// Stats holds metrics gathered during a slice operation.
type Stats struct {
	TotalLines   int
	MatchedLines int
	SkippedLines int
	BytesRead    int64
	BytesWritten int64
	Duration     time.Duration
	start        time.Time
}

// New creates a new Stats instance and starts the internal timer.
func New() *Stats {
	return &Stats{start: time.Now()}
}

// RecordLine records a processed line, marking whether it was matched or skipped.
func (s *Stats) RecordLine(matched bool, bytesRead int, bytesWritten int) {
	s.TotalLines++
	s.BytesRead += int64(bytesRead)
	if matched {
		s.MatchedLines++
		s.BytesWritten += int64(bytesWritten)
	} else {
		s.SkippedLines++
	}
}

// Finish stops the timer and records the total duration.
func (s *Stats) Finish() {
	s.Duration = time.Since(s.start)
}

// WriteTo writes a human-readable summary of the stats to w.
func (s *Stats) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w,
		"lines: total=%d matched=%d skipped=%d | bytes: read=%d written=%d | duration=%s\n",
		s.TotalLines,
		s.MatchedLines,
		s.SkippedLines,
		s.BytesRead,
		s.BytesWritten,
		s.Duration.Round(time.Millisecond),
	)
	return int64(n), err
}
