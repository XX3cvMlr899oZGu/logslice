package rotate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logslice/internal/rotate"
)

func makeFiles(t *testing.T, dir string, names []string) {
	t.Helper()
	for _, n := range names {
		p := filepath.Join(dir, n)
		if err := os.WriteFile(p, []byte("line\n"), 0o644); err != nil {
			t.Fatalf("create %s: %v", p, err)
		}
	}
}

func TestFindSegmentsOrderedOldestFirst(t *testing.T) {
	dir := t.TempDir()
	makeFiles(t, dir, []string{"app.log", "app.log.1", "app.log.2"})

	segs, err := rotate.FindSegments(dir, "app.log")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(segs) != 3 {
		t.Fatalf("expected 3 segments, got %d", len(segs))
	}
	// Live file must be last.
	last := segs[len(segs)-1]
	if filepath.Base(last.Path) != "app.log" {
		t.Errorf("expected live file last, got %q", last.Path)
	}
	// Rotated files must precede live file.
	if filepath.Base(segs[0].Path) != "app.log.1" && filepath.Base(segs[0].Path) != "app.log.2" {
		t.Errorf("unexpected first segment: %q", segs[0].Path)
	}
}

func TestFindSegmentsSingleFile(t *testing.T) {
	dir := t.TempDir()
	makeFiles(t, dir, []string{"service.log"})

	segs, err := rotate.FindSegments(dir, "service.log")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(segs) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(segs))
	}
	if segs[0].Index != 0 {
		t.Errorf("expected index 0, got %d", segs[0].Index)
	}
}

func TestFindSegmentsNoneFound(t *testing.T) {
	dir := t.TempDir()
	_, err := rotate.FindSegments(dir, "missing.log")
	if err == nil {
		t.Fatal("expected error for missing prefix, got nil")
	}
}

func TestFindSegmentsIgnoresUnrelated(t *testing.T) {
	dir := t.TempDir()
	makeFiles(t, dir, []string{"app.log", "other.log", "app.log.1"})

	segs, err := rotate.FindSegments(dir, "app.log")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, s := range segs {
		base := filepath.Base(s.Path)
		if base == "other.log" {
			t.Errorf("unrelated file included: %q", s.Path)
		}
	}
}
