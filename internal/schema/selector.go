package schema

import (
	"fmt"
	"strings"
)

// Selector evaluates a simple field==value or field!=value predicate
// against extracted log line fields.
type Selector struct {
	field string
	op    string
	value string
}

// ParseSelector parses an expression of the form "field==value" or "field!=value".
func ParseSelector(expr string) (*Selector, error) {
	for _, op := range []string{"!=", "=="} {
		if idx := strings.Index(expr, op); idx >= 0 {
			field := strings.TrimSpace(expr[:idx])
			value := strings.TrimSpace(expr[idx+len(op):])
			if field == "" {
				return nil, fmt.Errorf("schema: empty field name in selector %q", expr)
			}
			return &Selector{field: field, op: op, value: value}, nil
		}
	}
	return nil, fmt.Errorf("schema: selector %q must contain == or !=", expr)
}

// Match reports whether the given Fields satisfy the selector.
func (s *Selector) Match(f Fields) bool {
	actual, ok := f[s.field]
	switch s.op {
	case "==":
		return ok && actual == s.value
	case "!=":
		return !ok || actual != s.value
	}
	return false
}

// String returns the canonical string representation of the selector.
func (s *Selector) String() string {
	return s.field + s.op + s.value
}
