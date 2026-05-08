package truncate

import (
	"strings"
	"testing"
)

func TestZeroMaxPassesThrough(t *testing.T) {
	tr, err := New(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	long := strings.Repeat("a", 500)
	if got := tr.Apply(long); got != long {
		t.Errorf("expected unchanged line, got len %d", len(got))
	}
}

func TestNegativeMaxReturnsError(t *testing.T) {
	_, err := New(-1)
	if err == nil {
		t.Fatal("expected error for negative max")
	}
}

func TestShortLineUnchanged(t *testing.T) {
	tr, _ := New(80)
	line := "short line"
	if got := tr.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestExactLengthUnchanged(t *testing.T) {
	tr, _ := New(10)
	line := "1234567890"
	if got := tr.Apply(line); got != line {
		t.Errorf("expected unchanged line at exact max")
	}
}

func TestLongLineTruncated(t *testing.T) {
	tr, _ := New(20)
	line := strings.Repeat("x", 40)
	got := tr.Apply(line)
	if len(got) != 20 {
		t.Errorf("expected len 20, got %d", len(got))
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("expected default suffix, got %q", got)
	}
}

func TestCustomSuffix(t *testing.T) {
	tr, _ := New(15, WithSuffix(" [cut]"))
	line := strings.Repeat("a", 30)
	got := tr.Apply(line)
	if len(got) != 15 {
		t.Errorf("expected len 15, got %d", len(got))
	}
	if !strings.HasSuffix(got, " [cut]") {
		t.Errorf("expected custom suffix, got %q", got)
	}
}

func TestSuffixLongerThanMax(t *testing.T) {
	// suffix "..." is 3 bytes; max 2 means keep=0, suffix clipped to 2 bytes
	tr, _ := New(2)
	line := "hello world"
	got := tr.Apply(line)
	if len(got) > 2 {
		t.Errorf("result must not exceed max, got len %d: %q", len(got), got)
	}
}

func TestApplyEmptyLine(t *testing.T) {
	tr, _ := New(10)
	if got := tr.Apply(""); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}
