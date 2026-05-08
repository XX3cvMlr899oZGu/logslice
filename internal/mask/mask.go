// Package mask provides field-level masking for structured log lines.
// It replaces sensitive field values with a fixed placeholder string,
// supporting both JSON-style quoted values and bare word tokens.
package mask

import (
	"regexp"
	"strings"
)

const defaultPlaceholder = "***"

// Masker replaces the values of named fields in a log line.
type Masker struct {
	patterns    []*regexp.Regexp
	placeholder string
}

// Option configures a Masker.
type Option func(*Masker)

// WithPlaceholder overrides the default replacement string.
func WithPlaceholder(p string) Option {
	return func(m *Masker) {
		if p != "" {
			m.placeholder = p
		}
	}
}

// New creates a Masker that will redact the given field names.
// Field names are matched case-sensitively.
func New(fields []string, opts ...Option) *Masker {
	m := &Masker{placeholder: defaultPlaceholder}
	for _, o := range opts {
		o(m)
	}
	for _, f := range fields {
		escaped := regexp.QuoteMeta(f)
		// Matches: "field":"value" or "field":bare or field=value or field="value"
		pat := regexp.MustCompile(
			`(?i)("` + escaped + `"\s*:\s*|` + escaped + `\s*=\s*)("[^"]*"|\S+)`,
		)
		m.patterns = append(m.patterns, pat)
	}
	return m
}

// Apply returns a copy of line with all configured field values masked.
func (m *Masker) Apply(line string) string {
	for _, pat := range m.patterns {
		line = pat.ReplaceAllStringFunc(line, func(match string) string {
			// Preserve the key portion; replace only the value.
			loc := pat.FindStringSubmatchIndex(match)
			if len(loc) < 6 {
				return match
			}
			key := match[loc[2]:loc[3]]
			val := match[loc[4]:loc[5]]
			masked := m.placeholder
			if strings.HasPrefix(val, `"`) {
				masked = `"` + m.placeholder + `"`
			}
			return key + masked
		})
	}
	return line
}
