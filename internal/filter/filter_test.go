package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
)

func TestNewUnknownLevel(t *testing.T) {
	_, err := filter.New(filter.Options{MinLevel: "trace"})
	if err == nil {
		t.Fatal("expected error for unknown level, got nil")
	}
}

func TestAcceptKeyword(t *testing.T) {
	f, err := filter.New(filter.Options{Keyword: "database"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Accept("2024-01-01T00:00:00Z INFO database connection established") {
		t.Error("expected line with keyword to be accepted")
	}
	if f.Accept("2024-01-01T00:00:00Z INFO user logged in") {
		t.Error("expected line without keyword to be rejected")
	}
}

func TestAcceptMinLevel(t *testing.T) {
	f, err := filter.New(filter.Options{MinLevel: "warn"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cases := []struct {
		line   string
		want   bool
	}{
		{"2024-01-01T00:00:00Z ERROR something failed", true},
		{"2024-01-01T00:00:00Z WARN disk space low", true},
		{"2024-01-01T00:00:00Z INFO server started", false},
		{"2024-01-01T00:00:00Z DEBUG request received", false},
	}
	for _, tc := range cases {
		got := f.Accept(tc.line)
		if got != tc.want {
			t.Errorf("Accept(%q) = %v, want %v", tc.line, got, tc.want)
		}
	}
}

func TestAcceptNoFilters(t *testing.T) {
	f, err := filter.New(filter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := []string{
		"2024-01-01T00:00:00Z DEBUG anything",
		"2024-01-01T00:00:00Z INFO something",
		"plain text with no level",
	}
	for _, l := range lines {
		if !f.Accept(l) {
			t.Errorf("expected all lines accepted with no filters, got false for %q", l)
		}
	}
}

func TestAcceptCombinedFilters(t *testing.T) {
	f, err := filter.New(filter.Options{MinLevel: "error", Keyword: "timeout"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Accept("2024-01-01T00:00:00Z ERROR connection timeout") {
		t.Error("expected ERROR+keyword line to be accepted")
	}
	if f.Accept("2024-01-01T00:00:00Z ERROR disk full") {
		t.Error("expected ERROR line without keyword to be rejected")
	}
	if f.Accept("2024-01-01T00:00:00Z WARN timeout approaching") {
		t.Error("expected WARN+keyword line to be rejected (below min level)")
	}
}

func TestAcceptKeywordCaseInsensitive(t *testing.T) {
	f, err := filter.New(filter.Options{Keyword: "timeout"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cases := []struct {
		line string
		want bool
	}{
		{"2024-01-01T00:00:00Z ERROR connection TIMEOUT", true},
		{"2024-01-01T00:00:00Z ERROR connection Timeout", true},
		{"2024-01-01T00:00:00Z ERROR connection timeout", true},
		{"2024-01-01T00:00:00Z ERROR disk full", false},
	}
	for _, tc := range cases {
		got := f.Accept(tc.line)
		if got != tc.want {
			t.Errorf("Accept(%q) = %v, want %v", tc.line, got, tc.want)
		}
	}
}
