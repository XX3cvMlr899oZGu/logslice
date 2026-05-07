// Package checkpoint implements lightweight resumable-run support for
// logslice. When slicing very large log files the process may be
// interrupted (SIGINT, power loss, etc.). Checkpoint persists the last
// successfully written byte offset together with the source filename so
// that the next invocation can seek directly to that position and
// continue without re-processing already-emitted lines.
//
// Typical usage:
//
//	state, err := checkpoint.Load(cpPath)
//	if errors.Is(err, checkpoint.ErrNotFound) {
//	    state = checkpoint.State{File: logPath}
//	}
//	// ... after each batch ...
//	checkpoint.Save(cpPath, state)
//	// ... on success ...
//	checkpoint.Remove(cpPath)
package checkpoint
