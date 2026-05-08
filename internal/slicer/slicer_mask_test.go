package slicer_test

import (
	"bufio"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/mask"
	"github.com/yourorg/logslice/internal/slicer"
)

func buildMaskedLogFile(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "mask-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	base := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	lines := []string{
		`{"ts":"2024-01-01T10:00:00Z","level":"info","password":"hunter2","msg":"login"}`,
		`{"ts":"2024-01-01T10:01:00Z","level":"info","token":"abc123","msg":"auth"}`,
		`{"ts":"2024-01-01T10:02:00Z","level":"warn","msg":"no secrets here"}`,
	}
	_ = base
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return f.Name()
}

func TestSliceMasksFields(t *testing.T) {
	path := buildMaskedLogFile(t)
	out := strings.Builder{}

	m := mask.New([]string{"password", "token"})

	err := slicer.Slice(slicer.Options{
		FilePath:  path,
		From:      time.Date(2024, 1, 1, 9, 59, 0, 0, time.UTC),
		To:        time.Date(2024, 1, 1, 10, 5, 0, 0, time.UTC),
		Output:    &out,
		Masker:    m,
	})
	if err != nil {
		t.Fatalf("Slice error: %v", err)
	}

	result := out.String()
	if strings.Contains(result, "hunter2") {
		t.Errorf("password value not masked in output")
	}
	if strings.Contains(result, "abc123") {
		t.Errorf("token value not masked in output")
	}
	if !strings.Contains(result, "no secrets here") {
		t.Errorf("unrelated line missing from output")
	}
}

func TestSliceNilMaskerPassesThrough(t *testing.T) {
	path := buildMaskedLogFile(t)
	out := strings.Builder{}

	err := slicer.Slice(slicer.Options{
		FilePath: path,
		From:     time.Date(2024, 1, 1, 9, 59, 0, 0, time.UTC),
		To:       time.Date(2024, 1, 1, 10, 5, 0, 0, time.UTC),
		Output:   &out,
	})
	if err != nil {
		t.Fatalf("Slice error: %v", err)
	}

	sc := bufio.NewScanner(strings.NewReader(out.String()))
	count := 0
	for sc.Scan() {
		count++
	}
	if count != 3 {
		t.Errorf("expected 3 lines, got %d", count)
	}
}
