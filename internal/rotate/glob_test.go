package rotate

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// touchFile creates a file at path and sets its modification time to modTime.
func touchFile(t *testing.T, path string, modTime time.Time) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create %s: %v", path, err)
	}
	f.Close()
	if err := os.Chtimes(path, modTime, modTime); err != nil {
		t.Fatalf("chtimes %s: %v", path, err)
	}
}

func TestFindSegmentsOrderedOldestFirst(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "app.log")

	now := time.Now().Truncate(time.Second)
	touchFile(t, base, now)
	touchFile(t, base+".2", now.Add(-2*time.Hour))
	touchFile(t, base+".1", now.Add(-1*time.Hour))

	segs, err := FindSegments(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(segs) != 3 {
		t.Fatalf("expected 3 segments, got %d", len(segs))
	}
	for i := 1; i < len(segs); i++ {
		if segs[i].ModTime.Before(segs[i-1].ModTime) {
			t.Errorf("segment %d (%s) is older than segment %d (%s)",
				i, segs[i].Path, i-1, segs[i-1].Path)
		}
	}
}

func TestFindSegmentsSingleFile(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "app.log")
	touchFile(t, base, time.Now())

	segs, err := FindSegments(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(segs) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(segs))
	}
	if segs[0].Path != base {
		t.Errorf("expected path %s, got %s", base, segs[0].Path)
	}
}

func TestFindSegmentsNoneFound(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "missing.log")

	segs, err := FindSegments(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(segs) != 0 {
		t.Fatalf("expected 0 segments, got %d", len(segs))
	}
}

func TestFindSegmentsIgnoresUnrelated(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "app.log")
	touchFile(t, base, time.Now())
	touchFile(t, filepath.Join(dir, "other.log"), time.Now())

	segs, err := FindSegments(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, s := range segs {
		if filepath.Base(s.Path) == "other.log" {
			t.Errorf("unrelated file other.log should not be included")
		}
	}
}
