// Package truncate clips log lines that exceed a configured byte length.
//
// Usage:
//
//	tr, err := truncate.New(120, truncate.WithSuffix(" [truncated]"))
//	if err != nil {
//		log.Fatal(err)
//	}
//	short := tr.Apply(longLine)
//
// A max of 0 disables truncation and every line is passed through unchanged.
// The suffix (default "...") is counted as part of the maximum length, so the
// returned string is never longer than max bytes.
package truncate
