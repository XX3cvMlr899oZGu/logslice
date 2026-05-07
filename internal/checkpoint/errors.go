package checkpoint

import (
	"errors"
	"fmt"
)

// ErrNotFound is returned by Load when no checkpoint file exists at the
// requested path.
var ErrNotFound = errors.New("checkpoint: no checkpoint file found")

// ErrCorrupt is returned by Load when a checkpoint file exists but cannot
// be parsed or has failed integrity validation.
var ErrCorrupt = errors.New("checkpoint: checkpoint file is corrupt")

// PathError wraps an underlying error with the file path that caused it,
// providing more context for debugging.
type PathError struct {
	Path string
	Err  error
}

func (e *PathError) Error() string {
	return fmt.Sprintf("checkpoint: error at path %q: %v", e.Path, e.Err)
}

func (e *PathError) Unwrap() error {
	return e.Err
}
