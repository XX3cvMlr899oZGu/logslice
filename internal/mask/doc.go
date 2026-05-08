// Package mask implements field-level value masking for log lines.
//
// It is designed to remove or replace sensitive data (passwords, tokens,
// API keys, etc.) before log lines are written to an output destination.
//
// Both JSON-encoded log lines and key=value formatted lines are supported.
//
// Example usage:
//
//	m := mask.New([]string{"password", "token"})
//	safe := m.Apply(rawLine)
//
// By default masked values are replaced with "***". Use WithPlaceholder
// to supply a custom replacement string.
package mask
