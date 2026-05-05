// Package reader provides efficient log file reading primitives for logslice.
//
// It exposes two main capabilities:
//
//  1. Sequential line reading via LineReader, which wraps a file or any
//     io.Reader and iterates lines with minimal allocations.
//
//  2. Binary search via SeekToTime, which locates the byte offset of the
//     first log line whose timestamp is >= a given target time. This allows
//     O(log N) seek into large, time-sorted log files without reading every
//     line from the start.
//
// Typical usage:
//
//	f, _ := os.Open("app.log")
//	offset, _ := reader.SeekToTime(f, startTime)
//	f.Seek(offset, io.SeekStart)
//	lr := reader.NewLineReaderFromFile(f)
//	for lr.Next() {
//		// process lr.Line()
//	}
package reader
