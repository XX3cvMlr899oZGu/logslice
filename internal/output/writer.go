// Package output provides utilities for writing sliced log segments
// to various destinations such as stdout or files.
package output

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Options configures the output destination for sliced log lines.
type Options struct {
	// OutputPath is the file path to write output to.
	// If empty, output is written to Stdout.
	OutputPath string

	// Stdout is the writer used when OutputPath is empty.
	// Defaults to os.Stdout if nil.
	Stdout io.Writer
}

// Writer wraps a buffered writer and tracks lines written.
type Writer struct {
	w          *bufio.Writer
	closer     io.Closer
	LinesWritten int
}

// New creates a new Writer based on the provided Options.
// If Options.OutputPath is set, output is written to that file.
// Otherwise output goes to Options.Stdout (or os.Stdout by default).
func New(opts Options) (*Writer, error) {
	var w io.Writer
	var closer io.Closer

	if opts.OutputPath != "" {
		f, err := os.Create(opts.OutputPath)
		if err != nil {
			return nil, fmt.Errorf("output: create file %q: %w", opts.OutputPath, err)
		}
		w = f
		closer = f
	} else {
		if opts.Stdout != nil {
			w = opts.Stdout
		} else {
			w = os.Stdout
		}
		closer = io.NopCloser(w)
	}

	return &Writer{
		w:      bufio.NewWriter(w),
		closer: closer,
	}, nil
}

// WriteLine writes a single line followed by a newline character.
func (wr *Writer) WriteLine(line string) error {
	_, err := fmt.Fprintln(wr.w, line)
	if err != nil {
		return fmt.Errorf("output: write line: %w", err)
	}
	wr.LinesWritten++
	return nil
}

// Close flushes buffered data and closes any underlying file.
func (wr *Writer) Close() error {
	if err := wr.w.Flush(); err != nil {
		return fmt.Errorf("output: flush: %w", err)
	}
	return wr.closer.Close()
}
