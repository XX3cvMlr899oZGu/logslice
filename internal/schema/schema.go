// Package schema provides log line field extraction and schema detection
// for structured log formats (JSON, logfmt, and plain key=value).
package schema

import (
	"encoding/json"
	"strings"
)

// Format represents a detected log line format.
type Format int

const (
	FormatUnknown Format = iota
	FormatJSON
	FormatLogfmt
)

// Fields holds the key-value pairs extracted from a log line.
type Fields map[string]string

// Detect returns the Format of the given log line.
func Detect(line string) Format {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) == 0 {
		return FormatUnknown
	}
	if trimmed[0] == '{' {
		return FormatJSON
	}
	if strings.Contains(trimmed, "=") {
		return FormatLogfmt
	}
	return FormatUnknown
}

// Extract parses key-value fields from a log line using the detected format.
// Returns an empty Fields map if the line cannot be parsed.
func Extract(line string) Fields {
	switch Detect(line) {
	case FormatJSON:
		return extractJSON(line)
	case FormatLogfmt:
		return extractLogfmt(line)
	default:
		return Fields{}
	}
}

func extractJSON(line string) Fields {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return Fields{}
	}
	out := make(Fields, len(raw))
	for k, v := range raw {
		switch sv := v.(type) {
		case string:
			out[k] = sv
		default:
			b, _ := json.Marshal(v)
			out[k] = string(b)
		}
	}
	return out
}

func extractLogfmt(line string) Fields {
	out := Fields{}
	for _, pair := range strings.Fields(line) {
		idx := strings.IndexByte(pair, '=')
		if idx < 0 {
			continue
		}
		k := pair[:idx]
		v := strings.Trim(pair[idx+1:], `"`)
		out[k] = v
	}
	return out
}
