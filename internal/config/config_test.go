package config

import (
	"flag"
	"testing"
	"time"
)

func newFS() *flag.FlagSet {
	return flag.NewFlagSet("test", flag.ContinueOnError)
}

func TestParseValidOptions(t *testing.T) {
	opts, err := Parse(newFS(), []string{
		"--file", "/var/log/app.log",
		"--from", "2024-01-01T00:00:00",
		"--to", "2024-01-01T01:00:00",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.FilePath != "/var/log/app.log" {
		t.Errorf("expected file path, got %q", opts.FilePath)
	}
	expectedFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if !opts.From.Equal(expectedFrom) {
		t.Errorf("unexpected From: %v", opts.From)
	}
}

func TestParseMissingFile(t *testing.T) {
	_, err := Parse(newFS(), []string{
		"--from", "2024-01-01T00:00:00",
		"--to", "2024-01-01T01:00:00",
	})
	if err != ErrMissingFile {
		t.Errorf("expected ErrMissingFile, got %v", err)
	}
}

func TestParseMissingFrom(t *testing.T) {
	_, err := Parse(newFS(), []string{
		"--file", "app.log",
		"--to", "2024-01-01T01:00:00",
	})
	if err != ErrMissingFrom {
		t.Errorf("expected ErrMissingFrom, got %v", err)
	}
}

func TestParseMissingTo(t *testing.T) {
	_, err := Parse(newFS(), []string{
		"--file", "app.log",
		"--from", "2024-01-01T00:00:00",
	})
	if err != ErrMissingTo {
		t.Errorf("expected ErrMissingTo, got %v", err)
	}
}

func TestParseInvalidRange(t *testing.T) {
	_, err := Parse(newFS(), []string{
		"--file", "app.log",
		"--from", "2024-01-01T02:00:00",
		"--to", "2024-01-01T01:00:00",
	})
	if err != ErrInvalidRange {
		t.Errorf("expected ErrInvalidRange, got %v", err)
	}
}

func TestParseOptionalFlags(t *testing.T) {
	opts, err := Parse(newFS(), []string{
		"--file", "app.log",
		"--from", "2024-01-01T00:00:00",
		"--to", "2024-01-01T01:00:00",
		"--level", "WARN",
		"--keyword", "timeout",
		"--output", "/tmp/out.log",
		"--stats",
		"--progress",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Level != "WARN" {
		t.Errorf("expected level WARN, got %q", opts.Level)
	}
	if opts.Keyword != "timeout" {
		t.Errorf("expected keyword 'timeout', got %q", opts.Keyword)
	}
	if !opts.Stats {
		t.Error("expected Stats=true")
	}
	if !opts.Progress {
		t.Error("expected Progress=true")
	}
}
