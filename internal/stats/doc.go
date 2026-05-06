// Package stats collects and reports metrics produced during a logslice
// operation.
//
// Usage:
//
//	s := stats.New()
//	// ... process lines ...
//	s.RecordLine(matched, bytesRead, bytesWritten)
//	s.Finish()
//	s.WriteTo(os.Stderr)
//
// Stats is intentionally simple and not goroutine-safe; it is designed to be
// used from a single processing goroutine and inspected after the slice
// operation completes.
package stats
