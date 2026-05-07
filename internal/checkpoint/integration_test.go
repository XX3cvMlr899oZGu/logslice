package checkpoint_test

import (
	"errors"
	"testing"

	"github.com/example/logslice/internal/checkpoint"
)

// TestResumeWorkflow simulates the full resume lifecycle:
// save progress mid-run → load on restart → remove on completion.
func TestResumeWorkflow(t *testing.T) {
	path := tempPath(t)

	// First run: no checkpoint exists yet.
	_, err := checkpoint.Load(path)
	if !errors.Is(err, checkpoint.ErrNotFound) {
		t.Fatalf("expected ErrNotFound on first run, got %v", err)
	}

	// Mid-run: persist progress after processing 500 lines.
	mid := checkpoint.State{File: "/logs/service.log", Offset: 8192, Lines: 500}
	if err := checkpoint.Save(path, mid); err != nil {
		t.Fatalf("mid-run Save: %v", err)
	}

	// Restart: load the saved state.
	resume, err := checkpoint.Load(path)
	if err != nil {
		t.Fatalf("resume Load: %v", err)
	}
	if resume.Offset != mid.Offset {
		t.Errorf("resume Offset: got %d, want %d", resume.Offset, mid.Offset)
	}
	if resume.Lines != mid.Lines {
		t.Errorf("resume Lines: got %d, want %d", resume.Lines, mid.Lines)
	}

	// Completion: remove checkpoint.
	if err := checkpoint.Remove(path); err != nil {
		t.Fatalf("Remove: %v", err)
	}

	// Verify it is gone.
	_, err = checkpoint.Load(path)
	if !errors.Is(err, checkpoint.ErrNotFound) {
		t.Errorf("expected ErrNotFound after Remove, got %v", err)
	}
}
