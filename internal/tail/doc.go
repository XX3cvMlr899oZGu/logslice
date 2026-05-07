// Package tail provides functionality for following a log file in real-time,
// similar to `tail -f`. It watches for new lines appended to a file and emits
// them through a channel, respecting optional filter and output configurations.
//
// The Tailer handles file rotation by detecting when a file is truncated or
// replaced, automatically reopening it to resume tailing from the beginning.
package tail
