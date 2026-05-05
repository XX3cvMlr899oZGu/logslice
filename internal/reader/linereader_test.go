package reader

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestLineReaderSequential(t *testing.T) {
	content := "line one\nline two\nline three\n"
	path := writeTempFile(t, content)

	r, err := NewLineReader(path)
	if err != nil {
		t.Fatalf("NewLineReader error: %v", err)
	}
	defer r.Close()

	expected := []string{"line one", "line two", "line three"}
	var got []string
	for r.Next() {
		got = append(got, r.Line())
	}
	if r.Err() != nil {
		t.Fatalf("scanner error: %v", r.Err())
	}
	if len(got) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(got))
	}
	for i, line := range expected {
		if got[i] != line {
			t.Errorf("line %d: expected %q, got %q", i+1, line, got[i])
		}
	}
}

func TestLineReaderLineNum(t *testing.T) {
	path := writeTempFile(t, "a\nb\nc\n")
	r, err := NewLineReader(path)
	if err != nil {
		t.Fatalf("NewLineReader error: %v", err)
	}
	defer r.Close()

	count := 0
	for r.Next() {
		count++
		if r.LineNum() != count {
			t.Errorf("expected LineNum %d, got %d", count, r.LineNum())
		}
	}
}

func TestReadAllLines(t *testing.T) {
	input := "alpha\nbeta\ngamma"
	lines, err := ReadAllLines(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ReadAllLines error: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[1] != "beta" {
		t.Errorf("expected 'beta', got %q", lines[1])
	}
}

func TestLineReaderMissingFile(t *testing.T) {
	_, err := NewLineReader("/nonexistent/path/file.log")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
