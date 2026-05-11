package schema_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/schema"
)

func TestDetectJSON(t *testing.T) {
	if got := schema.Detect(`{"level":"info","msg":"ok"}`); got != schema.FormatJSON {
		t.Fatalf("expected FormatJSON, got %v", got)
	}
}

func TestDetectLogfmt(t *testing.T) {
	if got := schema.Detect(`level=info msg=hello`); got != schema.FormatLogfmt {
		t.Fatalf("expected FormatLogfmt, got %v", got)
	}
}

func TestDetectUnknown(t *testing.T) {
	if got := schema.Detect(`plain text log line`); got != schema.FormatUnknown {
		t.Fatalf("expected FormatUnknown, got %v", got)
	}
}

func TestDetectEmpty(t *testing.T) {
	if got := schema.Detect(``); got != schema.FormatUnknown {
		t.Fatalf("expected FormatUnknown for empty, got %v", got)
	}
}

func TestExtractJSON(t *testing.T) {
	fields := schema.Extract(`{"level":"error","msg":"boom","code":42}`)
	if fields["level"] != "error" {
		t.Errorf("level: got %q", fields["level"])
	}
	if fields["msg"] != "boom" {
		t.Errorf("msg: got %q", fields["msg"])
	}
	if fields["code"] != "42" {
		t.Errorf("code: got %q", fields["code"])
	}
}

func TestExtractLogfmt(t *testing.T) {
	fields := schema.Extract(`level=warn msg="disk full" host=srv1`)
	if fields["level"] != "warn" {
		t.Errorf("level: got %q", fields["level"])
	}
	if fields["msg"] != "disk full" {
		t.Errorf("msg: got %q", fields["msg"])
	}
	if fields["host"] != "srv1" {
		t.Errorf("host: got %q", fields["host"])
	}
}

func TestExtractUnknownReturnsEmpty(t *testing.T) {
	fields := schema.Extract(`this is just a plain log line`)
	if len(fields) != 0 {
		t.Errorf("expected empty fields, got %v", fields)
	}
}

func TestExtractInvalidJSONReturnsEmpty(t *testing.T) {
	fields := schema.Extract(`{not valid json`)
	if len(fields) != 0 {
		t.Errorf("expected empty fields for invalid JSON, got %v", fields)
	}
}

func TestExtractJSONNestedValueSkipped(t *testing.T) {
	// Nested objects are not representable as a flat string value;
	// Extract should omit keys whose values are objects or arrays.
	fields := schema.Extract(`{"level":"info","meta":{"host":"srv1"}}`)
	if fields["level"] != "info" {
		t.Errorf("level: got %q", fields["level"])
	}
	if _, ok := fields["meta"]; ok {
		t.Errorf("expected nested object key 'meta' to be omitted, but it was present")
	}
}
