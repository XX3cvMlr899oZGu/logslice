// Package annotate provides line annotation support for logslice.
// It prepends a configurable prefix or injects key=value fields into
// each log line before it is written to the output.
package annotate

import (
	"fmt"
	"strings"
)

// Annotator appends static key=value pairs to every log line.
type Annotator struct {
	fields []field
	suffix string
}

type field struct {
	key   string
	value string
}

// Option configures an Annotator.
type Option func(*Annotator)

// WithField adds a static key=value annotation to every line.
func WithField(key, value string) Option {
	return func(a *Annotator) {
		a.fields = append(a.fields, field{key: key, value: value})
	}
}

// WithSuffix sets a raw string that is appended after all key=value fields.
func WithSuffix(s string) Option {
	return func(a *Annotator) {
		a.suffix = s
	}
}

// New constructs an Annotator with the supplied options.
// Returns an error if any key is empty.
func New(opts ...Option) (*Annotator, error) {
	a := &Annotator{}
	for _, o := range opts {
		o(a)
	}
	for _, f := range a.fields {
		if strings.TrimSpace(f.key) == "" {
			return nil, fmt.Errorf("annotate: field key must not be empty")
		}
	}
	return a, nil
}

// Apply appends the configured annotations to line and returns the result.
// If no fields or suffix are configured the original line is returned unchanged.
func (a *Annotator) Apply(line string) string {
	if len(a.fields) == 0 && a.suffix == "" {
		return line
	}
	var sb strings.Builder
	sb.WriteString(line)
	for _, f := range a.fields {
		sb.WriteByte(' ')
		sb.WriteString(f.key)
		sb.WriteByte('=')
		sb.WriteString(f.value)
	}
	if a.suffix != "" {
		sb.WriteByte(' ')
		sb.WriteString(a.suffix)
	}
	return sb.String()
}
