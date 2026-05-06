// Package filter provides post-extraction filtering of log lines for logslice.
//
// After SeekToTime has located the relevant byte range in a log file, the
// filter package can be used to further narrow output by log severity level
// and/or keyword substring matching.
//
// Usage:
//
//	f, err := filter.New(filter.Options{
//		MinLevel: "warn",
//		Keyword:  "database",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	if f.Accept(line) {
//		fmt.Println(line)
//	}
package filter
