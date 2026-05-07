// Package tail provides functionality for following a log file in real-time,
// similar to `tail -f`, emitting new lines as they are appended.
package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// DefaultPollInterval is how often the tailer checks for new data when the
// file has not grown since the last read.
const DefaultPollInterval = 200 * time.Millisecond

// Line represents a single line emitted by the tailer.
type Line struct {
	Text   string
	Offset int64 // byte offset of the start of this line in the file
}

// Tailer follows a file, sending new lines to a channel as they appear.
type Tailer struct {
	path         string
	pollInterval time.Duration
	lines        chan Line
	errs         chan error
}

// New creates a Tailer for the given file path.
// pollInterval controls how frequently the file is polled for new content;
// pass 0 to use DefaultPollInterval.
func New(path string, pollInterval time.Duration) *Tailer {
	if pollInterval <= 0 {
		pollInterval = DefaultPollInterval
	}
	return &Tailer{
		path:         path,
		pollInterval: pollInterval,
		lines:        make(chan Line, 64),
		errs:         make(chan error, 1),
	}
}

// Lines returns the channel on which new lines are delivered.
func (t *Tailer) Lines() <-chan Line { return t.lines }

// Errs returns the channel on which a terminal error is delivered.
// At most one error will be sent; the channel is closed after that.
func (t *Tailer) Errs() <-chan error { return t.errs }

// Follow opens the file, seeks to the end, and begins emitting lines until
// ctx is cancelled. It is intended to be run in its own goroutine.
func (t *Tailer) Follow(ctx context.Context) {
	defer close(t.lines)
	defer close(t.errs)

	f, err := os.Open(t.path)
	if err != nil {
		t.errs <- err
		return
	}
	defer f.Close()

	// Seek to end so we only tail new content.
	offset, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		t.errs <- err
		return
	}

	reader := bufio.NewReader(f)
	ticker := time.NewTicker(t.pollInterval)
	defer ticker.Stop()

	for {
		// Drain all available lines before waiting for the next tick.
		for {
			text, readErr := reader.ReadString('\n')
			if len(text) > 0 {
				// Strip the trailing newline for consistency with other
				// components in the pipeline.
				if len(text) > 0 && text[len(text)-1] == '\n' {
					text = text[:len(text)-1]
				}
				select {
				case t.lines <- Line{Text: text, Offset: offset}:
				case <-ctx.Done():
					return
				}
				offset += int64(len(text)) + 1
			}
			if readErr == io.EOF {
				// No more data right now; wait for the next poll tick.
				break
			}
			if readErr != nil {
				t.errs <- readErr
				return
			}
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Reset the reader so it picks up any new bytes written since the
			// last read (the underlying *os.File position is already advanced).
			reader.Reset(f)
		}
	}
}
