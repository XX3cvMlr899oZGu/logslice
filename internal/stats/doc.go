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
// The Stats type tracks the following metrics:
//   - Total lines seen and lines matched by the slice filter
//   - Bytes read from the input source
//   - Bytes written to the output destination
//   - Elapsed wall-clock time between New() and Finish()
//
// Stats is intentionally simple and not goroutine-safe; it is designed to be
// used from a single processing goroutine and inspected after the slice
// operation completes.
package stats
