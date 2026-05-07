package tail

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

const defaultPollInterval = 250 * time.Millisecond

// Options configures Tailer behaviour.
type Options struct {
	// PollInterval is how often the file is checked for new content.
	// Defaults to 250ms when zero.
	PollInterval time.Duration
}

// Tailer follows a file, emitting new lines as they are appended.
type Tailer struct {
	path   string
	opts   Options
	lines  chan string
	stopCh chan struct{}
}

// New opens path and starts tailing it. Returns an error if the file cannot
// be opened.
func New(path string, opts Options) (*Tailer, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("tail: open %s: %w", path, err)
	}
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		f.Close()
		return nil, fmt.Errorf("tail: seek %s: %w", path, err)
	}

	if opts.PollInterval == 0 {
		opts.PollInterval = defaultPollInterval
	}

	t := &Tailer{
		path:   path,
		opts:   opts,
		lines:  make(chan string, 64),
		stopCh: make(chan struct{}),
	}

	go t.follow(f)
	return t, nil
}

// Lines returns the channel on which tailed lines are delivered.
func (t *Tailer) Lines() <-chan string {
	return t.lines
}

// Stop signals the tailer to stop and waits until the background goroutine
// has exited.
func (t *Tailer) Stop() {
	close(t.stopCh)
}

func (t *Tailer) follow(f *os.File) {
	defer close(t.lines)
	defer f.Close()

	reader := bufio.NewReader(f)
	ticker := time.NewTicker(t.opts.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-t.stopCh:
			return
		case <-ticker.C:
			if err := t.drainLines(f, reader); err != nil {
				return
			}
		}
	}
}

func (t *Tailer) drainLines(f *os.File, r *bufio.Reader) error {
	for {
		line, err := r.ReadString('\n')
		if len(line) > 0 {
			// Strip trailing newline before sending.
			if len(line) > 0 && line[len(line)-1] == '\n' {
				line = line[:len(line)-1]
			}
			select {
			case t.lines <- line:
			case <-t.stopCh:
				return fmt.Errorf("stopped")
			}
		}
		if err == io.EOF {
			// Check for truncation / rotation.
			if rotated, _ := isRotated(f); rotated {
				newF, err := os.Open(t.path)
				if err != nil {
					return err
				}
				f.Close()
				*f = *newF
				r.Reset(f)
			}
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func isRotated(f *os.File) (bool, error) {
	cur, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return false, err
	}
	info, err := f.Stat()
	if err != nil {
		return false, err
	}
	return info.Size() < cur, nil
}
