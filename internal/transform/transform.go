// Package transform provides line transformation utilities for log output,
// such as stripping ANSI colour codes or redacting sensitive fields.
package transform

import (
	"regexp"
	"strings"
)

// ansiEscape matches ANSI terminal escape sequences.
var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// Transformer applies a chain of transformations to a log line.
type Transformer struct {
	stripANSI bool
	redactKeys []string
}

// Option configures a Transformer.
type Option func(*Transformer)

// WithStripANSI instructs the Transformer to remove ANSI escape sequences.
func WithStripANSI() Option {
	return func(t *Transformer) { t.stripANSI = true }
}

// WithRedactKeys instructs the Transformer to replace the values of the given
// JSON-style keys with "[REDACTED]".
func WithRedactKeys(keys ...string) Option {
	return func(t *Transformer) {
		t.redactKeys = append(t.redactKeys, keys...)
	}
}

// New returns a Transformer configured with the supplied options.
func New(opts ...Option) *Transformer {
	t := &Transformer{}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Apply runs all configured transformations on line and returns the result.
func (t *Transformer) Apply(line string) string {
	if t.stripANSI {
		line = ansiEscape.ReplaceAllString(line, "")
	}
	for _, key := range t.redactKeys {
		line = redactKey(line, key)
	}
	return line
}

// redactKey replaces the value of a key in a loosely structured log line.
// It handles both quoted-string values and unquoted word values.
func redactKey(line, key string) string {
	// Match: key="value" or key=value (no spaces)
	pattern := regexp.MustCompile(
		`(?i)(` + regexp.QuoteMeta(key) + `=)("[^"]*"|\S+)`,
	)
	return pattern.ReplaceAllStringFunc(line, func(match string) string {
		idx := strings.Index(match, "=")
		if idx < 0 {
			return match
		}
		return match[:idx+1] + "[REDACTED]"
	})
}
