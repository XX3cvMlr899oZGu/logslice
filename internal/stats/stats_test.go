package stats_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/stats"
)

func TestRecordLineCountsCorrectly(t *testing.T) {
	s := stats.New()

	s.RecordLine(true, 100, 100)
	s.RecordLine(true, 80, 80)
	s.RecordLine(false, 60, 0)

	if s.TotalLines != 3 {
		t.Errorf("TotalLines: want 3, got %d", s.TotalLines)
	}
	if s.MatchedLines != 2 {
		t.Errorf("MatchedLines: want 2, got %d", s.MatchedLines)
	}
	if s.SkippedLines != 1 {
		t.Errorf("SkippedLines: want 1, got %d", s.SkippedLines)
	}
	if s.BytesRead != 240 {
		t.Errorf("BytesRead: want 240, got %d", s.BytesRead)
	}
	if s.BytesWritten != 180 {
		t.Errorf("BytesWritten: want 180, got %d", s.BytesWritten)
	}
}

func TestFinishSetsDuration(t *testing.T) {
	s := stats.New()
	s.Finish()

	if s.Duration <= 0 {
		t.Error("expected Duration > 0 after Finish")
	}
}

func TestWriteToProducesOutput(t *testing.T) {
	s := stats.New()
	s.RecordLine(true, 50, 50)
	s.RecordLine(false, 30, 0)
	s.Finish()

	var buf bytes.Buffer
	n, err := s.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo error: %v", err)
	}
	if n == 0 {
		t.Error("expected non-zero bytes written")
	}

	out := buf.String()
	for _, want := range []string{"total=2", "matched=1", "skipped=1", "read=80", "written=50"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q: %s", want, out)
		}
	}
}

func TestWriteToEmptyStats(t *testing.T) {
	s := stats.New()
	s.Finish()

	var buf bytes.Buffer
	_, err := s.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo error: %v", err)
	}
	if !strings.Contains(buf.String(), "total=0") {
		t.Errorf("expected total=0 in output: %s", buf.String())
	}
}
