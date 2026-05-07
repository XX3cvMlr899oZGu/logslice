// Package dedupe implements a sliding-window deduplication filter for log
// lines. It is useful when log files contain repeated identical entries
// (e.g. health-check noise) that should be suppressed in sliced output.
//
// Usage:
//
//	df := dedupe.New(64) // remember last 64 distinct lines
//	for _, line := range lines {
//		if !df.IsDuplicate(line) {
//			fmt.Println(line)
//		}
//	}
package dedupe
