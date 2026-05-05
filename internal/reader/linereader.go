// Package reader provides utilities for reading log files efficiently,
// including line-by-line iteration and binary search for time-range slicing.
package reader

import (
	"bufio"
	"io"
	"os"
)

// LineReader wraps a file and provides sequential line reading.
type LineReader struct {
	file    *os.File
	scanner *bufio.Scanner
	line    string
	lineNum int
}

// NewLineReader opens the file at path and returns a LineReader.
func NewLineReader(path string) (*LineReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	return &LineReader{file: f, scanner: scanner}, nil
}

// Next advances to the next line. Returns false when done or on error.
func (r *LineReader) Next() bool {
	if r.scanner.Scan() {
		r.line = r.scanner.Text()
		r.lineNum++
		return true
	}
	return false
}

// Line returns the current line text.
func (r *LineReader) Line() string {
	return r.line
}

// LineNum returns the current 1-based line number.
func (r *LineReader) LineNum() int {
	return r.lineNum
}

// Err returns any scanner error (not io.EOF).
func (r *LineReader) Err() error {
	return r.scanner.Err()
}

// Close releases the underlying file.
func (r *LineReader) Close() error {
	return r.file.Close()
}

// ReadAllLines reads all lines from an io.Reader into a string slice.
func ReadAllLines(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
