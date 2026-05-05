package reader

import (
	"io"
	"os"
	"time"

	"github.com/user/logslice/internal/parser"
)

// SeekToTime performs a binary search over a log file to find the byte offset
// of the first line whose timestamp is >= target. The file must be sorted by
// time. Returns the offset, or 0 if no suitable line is found.
func SeekToTime(f *os.File, target time.Time) (int64, error) {
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	size := fi.Size()
	if size == 0 {
		return 0, nil
	}

	lo, hi := int64(0), size
	for lo < hi {
		mid := (lo + hi) / 2
		lineStart, line, err := readLineAt(f, mid, size)
		if err != nil {
			return 0, err
		}
		if line == "" {
			break
		}
		ts, err := parser.ParseTimestamp(line)
		if err != nil || ts.Before(target) {
			lo = lineStart + int64(len(line)) + 1
		} else {
			hi = lineStart
		}
	}
	return lo, nil
}

// readLineAt seeks to offset mid, then scans forward to find the next complete
// line. Returns the line's start offset and its text.
func readLineAt(f *os.File, mid, size int64) (int64, string, error) {
	if _, err := f.Seek(mid, io.SeekStart); err != nil {
		return 0, "", err
	}
	buf := make([]byte, 4096)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return 0, "", err
	}
	buf = buf[:n]

	start := 0
	if mid > 0 {
		// skip partial line
		for start < len(buf) && buf[start] != '\n' {
			start++
		}
		start++ // skip the newline itself
	}
	if start >= len(buf) {
		return mid + int64(start), "", nil
	}
	end := start
	for end < len(buf) && buf[end] != '\n' {
		end++
	}
	lineStart := mid + int64(start)
	return lineStart, string(buf[start:end]), nil
}
