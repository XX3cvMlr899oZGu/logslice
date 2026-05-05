// Package parser provides utilities for extracting and interpreting structured
// information from log lines, with a focus on timestamp detection.
//
// Timestamp parsing
//
// ParseTimestamp attempts to identify a timestamp at the beginning of a log
// line by trying a curated list of common formats (CommonFormats). The formats
// are ordered from most specific to least specific to avoid ambiguous matches.
//
// Supported formats include:
//   - RFC 3339 / ISO 8601 (with and without sub-second precision)
//   - Space-separated datetime ("2006-01-02 15:04:05")
//   - Slash-separated datetime ("2006/01/02 15:04:05")
//   - Apache/Nginx Common Log Format ("02/Jan/2006:15:04:05 -0700")
//
// Example usage:
//
//	t, err := parser.ParseTimestamp("2024-03-15T08:00:00Z application ready")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(t)
package parser
