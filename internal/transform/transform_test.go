package transform_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/transform"
)

func TestNoOptionsIsIdentity(t *testing.T) {
	tr := transform.New()
	input := `2024-01-01T00:00:00Z level=info msg="hello world"`
	if got := tr.Apply(input); got != input {
		t.Errorf("expected identity, got %q", got)
	}
}

func TestStripANSIRemovesEscapes(t *testing.T) {
	tr := transform.New(transform.WithStripANSI())
	input := "\x1b[32mINFO\x1b[0m message"
	want := "INFO message"
	if got := tr.Apply(input); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestStripANSILeavesPlainTextAlone(t *testing.T) {
	tr := transform.New(transform.WithStripANSI())
	input := "plain log line with no escapes"
	if got := tr.Apply(input); got != input {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestRedactKeyQuotedValue(t *testing.T) {
	tr := transform.New(transform.WithRedactKeys("password"))
	input := `level=info msg="login" password="s3cr3t"`
	want := `level=info msg="login" password=[REDACTED]`
	if got := tr.Apply(input); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRedactKeyUnquotedValue(t *testing.T) {
	tr := transform.New(transform.WithRedactKeys("token"))
	input := `level=warn token=abc123 user=alice`
	want := `level=warn token=[REDACTED] user=alice`
	if got := tr.Apply(input); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRedactMultipleKeys(t *testing.T) {
	tr := transform.New(transform.WithRedactKeys("password", "token"))
	input := `password="hunter2" token=xyz user=bob`
	want := `password=[REDACTED] token=[REDACTED] user=bob`
	if got := tr.Apply(input); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRedactKeyNotPresent(t *testing.T) {
	tr := transform.New(transform.WithRedactKeys("secret"))
	input := `level=info msg="nothing sensitive here"`
	if got := tr.Apply(input); got != input {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestCombinedStripAndRedact(t *testing.T) {
	tr := transform.New(
		transform.WithStripANSI(),
		transform.WithRedactKeys("api_key"),
	)
	input := "\x1b[31mERROR\x1b[0m api_key=\"topsecret\" msg=\"oops\""
	want := "ERROR api_key=[REDACTED] msg=\"oops\""
	if got := tr.Apply(input); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
