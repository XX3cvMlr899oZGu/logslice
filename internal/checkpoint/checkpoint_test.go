package checkpoint_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/logslice/internal/checkpoint"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "logslice.checkpoint")
}

func TestSaveAndLoad(t *testing.T) {
	path := tempPath(t)
	want := checkpoint.State{File: "/var/log/app.log", Offset: 4096, Lines: 200}

	if err := checkpoint.Save(path, want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := checkpoint.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.File != want.File {
		t.Errorf("File: got %q, want %q", got.File, want.File)
	}
	if got.Offset != want.Offset {
		t.Errorf("Offset: got %d, want %d", got.Offset, want.Offset)
	}
	if got.Lines != want.Lines {
		t.Errorf("Lines: got %d, want %d", got.Lines, want.Lines)
	}
	if got.SavedAt.IsZero() {
		t.Error("SavedAt should be set after Save")
	}
}

func TestLoadNotFound(t *testing.T) {
	_, err := checkpoint.Load("/nonexistent/path/cp.json")
	if !errors.Is(err, checkpoint.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestRemoveDeletesFile(t *testing.T) {
	path := tempPath(t)
	if err := checkpoint.Save(path, checkpoint.State{File: "f"}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := checkpoint.Remove(path); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Error("expected file to be deleted")
	}
}

func TestRemoveIdempotent(t *testing.T) {
	if err := checkpoint.Remove("/nonexistent/cp.json"); err != nil {
		t.Errorf("Remove on missing file should not error, got %v", err)
	}
}

func TestSaveIsAtomic(t *testing.T) {
	path := tempPath(t)
	// Save twice; second save should overwrite cleanly with no .tmp residue.
	for i := 0; i < 2; i++ {
		if err := checkpoint.Save(path, checkpoint.State{File: "f", Offset: int64(i * 100)}); err != nil {
			t.Fatalf("Save %d: %v", i, err)
		}
	}
	if _, err := os.Stat(path + ".tmp"); !errors.Is(err, os.ErrNotExist) {
		t.Error("temp file should not remain after successful Save")
	}
}
