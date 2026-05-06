// Package rotate detects and enumerates rotated log file segments within a
// directory. Many log management systems (logrotate, newsyslog, etc.) produce
// a sequence of files such as:
//
//	app.log        – the current, active log file
//	app.log.1      – the previous rotation
//	app.log.2      – the rotation before that
//
// rotate.FindSegments scans a directory for all files matching a given prefix
// and returns them as an ordered []Segment slice, oldest first, so that
// callers can iterate from the beginning of time toward the present.
//
// This package is intentionally stateless and performs no I/O beyond
// directory enumeration; it does not open, parse, or modify any log file.
package rotate
