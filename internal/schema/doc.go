// Package schema provides utilities for detecting the structure of log lines
// and extracting fields from common structured formats.
//
// Supported formats:
//
//   - JSON: lines beginning with '{' are parsed as JSON objects.
//   - Logfmt: lines containing key=value or key="quoted value" pairs.
//
// Field extraction returns a flat string-to-string map regardless of the
// underlying format, making downstream processing format-agnostic.
//
// Selectors allow filtering lines by field value using == and != operators:
//
//	sel, err := schema.ParseSelector("level==error")
//	if err == nil && sel.Match(schema.Extract(line)) {
//	    // process line
//	}
package schema
