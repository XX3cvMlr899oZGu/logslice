package parser

import (
	"fmt"
	"time"
)

// CommonFormats is an ordered list of timestamp formats tried during parsing.
var CommonFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
	"02/Jan/2006:15:04:05 -0700", // Common Log Format
}

// ParseTimestamp attempts to extract and parse a timestamp from a raw log line.
// It tries each format in CommonFormats and returns the first successful match.
func ParseTimestamp(line string) (time.Time, error) {
	if len(line) == 0 {
		return time.Time{}, fmt.Errorf("empty line")
	}

	// Try progressively longer prefixes of the line to find a parseable timestamp.
	// Most timestamps appear at the start of a log line.
	maxPrefix := 40
	if len(line) < maxPrefix {
		maxPrefix = len(line)
	}

	for _, format := range CommonFormats {
		candidate := line
		if len(format) < maxPrefix {
			candidate = line[:len(format)]
			if len(candidate) > maxPrefix {
				continue
			}
		}
		t, err := time.Parse(format, candidate)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("no recognisable timestamp found in line prefix: %q", line[:maxPrefix])
}
