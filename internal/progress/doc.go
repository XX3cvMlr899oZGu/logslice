// Package progress provides lightweight progress reporting for logslice
// operations. It is designed to give operators visibility into how far
// through a large log file the slicer has advanced.
//
// Basic usage:
//
//	rep := progress.New(os.Stderr, fileSize, 2*time.Second)
//	rep.Start()
//	defer rep.Stop()
//
//	// inside your processing loop:
//	rep.Add(int64(len(line)))
//
// When total is 0 the reporter omits the percentage and only reports
// the raw byte count, which is useful when the file size is unavailable
// (e.g. when reading from stdin or a named pipe).
package progress
