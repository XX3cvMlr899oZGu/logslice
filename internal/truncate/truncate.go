// Package truncate provides line-length truncation for log output.
// Lines exceeding the configured maximum byte length are clipped and
// an optional suffix (e.g. "...") is appended to signal truncation.
package truncate

import "errors"

// ErrNegativeMax is returned when a negative max length is supplied.
var ErrNegativeMax = errors.New("truncate: max length must be >= 0")

// Truncator clips lines that exceed a maximum byte length.
type Truncator struct {
	max    int
	suffix []byte
}

// Option configures a Truncator.
type Option func(*Truncator)

// WithSuffix sets the bytes appended to a truncated line.
// The suffix length is counted against max.
func WithSuffix(s string) Option {
	return func(t *Truncator) {
		t.suffix = []byte(s)
	}
}

// New creates a Truncator that clips lines to at most max bytes.
// If max is 0 lines are returned unchanged.
func New(max int, opts ...Option) (*Truncator, error) {
	if max < 0 {
		return nil, ErrNegativeMax
	}
	t := &Truncator{max: max, suffix: []byte("...")}
	for _, o := range opts {
		o(t)
	}
	return t, nil
}

// Apply returns line truncated to the configured maximum.
// If max is 0 the original line is returned unchanged.
func (t *Truncator) Apply(line string) string {
	if t.max == 0 || len(line) <= t.max {
		return line
	}
	suf := t.suffix
	keep := t.max - len(suf)
	if keep < 0 {
		keep = 0
		suf = suf[:t.max]
	}
	buf := make([]byte, 0, t.max)
	buf = append(buf, line[:keep]...)
	buf = append(buf, suf...)
	return string(buf)
}
