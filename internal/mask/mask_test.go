package mask_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/mask"
)

func TestNoFieldsIsIdentity(t *testing.T) {
	m := mask.New(nil)
	line := `{"level":"info","msg":"hello"}`
	if got := m.Apply(line); got != line {
		t.Fatalf("expected identity, got %q", got)
	}
}

func TestMaskQuotedJSONValue(t *testing.T) {
	m := mask.New([]string{"password"})
	input := `{"user":"alice","password":"s3cr3t","ok":true}`
	got := m.Apply(input)
	if contains(got, "s3cr3t") {
		t.Fatalf("secret still present: %q", got)
	}
	if !contains(got, `"password"`) {
		t.Fatalf("key missing from output: %q", got)
	}
}

func TestMaskKeyValueFormat(t *testing.T) {
	m := mask.New([]string{"token"})
	input := `time=2024-01-01 level=info token=abc123 msg=ok`
	got := m.Apply(input)
	if contains(got, "abc123") {
		t.Fatalf("token value still present: %q", got)
	}
}

func TestMaskMultipleFields(t *testing.T) {
	m := mask.New([]string{"password", "token"})
	input := `{"password":"hunter2","token":"xyz","user":"bob"}`
	got := m.Apply(input)
	if contains(got, "hunter2") || contains(got, "xyz") {
		t.Fatalf("secrets still present: %q", got)
	}
	if !contains(got, `"user":"bob"`) {
		t.Fatalf("unmasked field altered: %q", got)
	}
}

func TestCustomPlaceholder(t *testing.T) {
	m := mask.New([]string{"secret"}, mask.WithPlaceholder("REDACTED"))
	input := `{"secret":"topsecret"}`
	got := m.Apply(input)
	if !contains(got, "REDACTED") {
		t.Fatalf("custom placeholder not found: %q", got)
	}
}

func TestEmptyPlaceholderFallsBackToDefault(t *testing.T) {
	m := mask.New([]string{"key"}, mask.WithPlaceholder(""))
	input := `{"key":"value"}`
	got := m.Apply(input)
	if contains(got, "value") {
		t.Fatalf("value not masked: %q", got)
	}
}

func TestMaskDoesNotAlterUnrelatedFields(t *testing.T) {
	m := mask.New([]string{"password"})
	input := `{"level":"warn","msg":"bad login","password":"oops"}`
	got := m.Apply(input)
	if !contains(got, `"level":"warn"`) {
		t.Fatalf("unrelated field altered: %q", got)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
