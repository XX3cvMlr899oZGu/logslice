package checkpoint

import "errors"

// ErrNotFound is returned by Load when no checkpoint file exists at the
// requested path.
var ErrNotFound = errors.New("checkpoint: no checkpoint file found")
