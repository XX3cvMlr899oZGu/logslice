package output

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteLineToStdout(t *testing.T) {
	var buf bytes.Buffer
	w, err := New(Options{Stdout: &buf})
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	lines := []string{"line one", "line two", "line three"}
	for _, l := range lines {
		if err := w.WriteLine(l); err != nil {
			t.Fatalf("WriteLine(%q) error: %v", l, err)
		}
	}

	if err := w.Close(); err != nil {
		t.Fatalf("Close() error: %v", err)
	}

	if w.LinesWritten != len(lines) {
		t.Errorf("LinesWritten = %d, want %d", w.LinesWritten, len(lines))
	}

	got := buf.String()
	for _, l := range lines {
		if !strings.Contains(got, l) {
			t.Errorf("output missing line %q", l)
		}
	}
}

func TestWriteLineToFile(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "out.log")

	w, err := New(Options{OutputPath: outPath})
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	expected := []string{"alpha", "beta", "gamma"}
	for _, l := range expected {
		if err := w.WriteLine(l); err != nil {
			t.Fatalf("WriteLine(%q) error: %v", l, err)
		}
	}

	if err := w.Close(); err != nil {
		t.Fatalf("Close() error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}

	got := string(data)
	for _, l := range expected {
		if !strings.Contains(got, l) {
			t.Errorf("file missing line %q", l)
		}
	}

	if w.LinesWritten != len(expected) {
		t.Errorf("LinesWritten = %d, want %d", w.LinesWritten, len(expected))
	}
}

func TestNewInvalidPath(t *testing.T) {
	_, err := New(Options{OutputPath: "/nonexistent/dir/out.log"})
	if err == nil {
		t.Error("expected error for invalid output path, got nil")
	}
}
