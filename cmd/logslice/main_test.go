package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func writeTempLogFile(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "log-*.log")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	defer f.Close()
	for _, l := range lines {
		fmt.Fprintln(f, l)
	}
	return f.Name()
}

func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "logslice")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("build failed: %v", err)
	}
	return bin
}

func TestMainMissingFlags(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit, got nil")
	}
	if !strings.Contains(stderr.String(), "required") {
		t.Errorf("expected 'required' in stderr, got: %s", stderr.String())
	}
}

func TestMainExtractsRange(t *testing.T) {
	base := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	var lines []string
	for i := 0; i < 10; i++ {
		ts := base.Add(time.Duration(i) * time.Minute).Format(time.RFC3339)
		lines = append(lines, fmt.Sprintf("%s level=info msg=\"line %d\"", ts, i))
	}
	logFile := writeTempLogFile(t, lines)
	bin := buildBinary(t)

	var stdout bytes.Buffer
	cmd := exec.Command(bin,
		"-f", logFile,
		"-start", "2024-01-01T10:02:00",
		"-end", "2024-01-01T10:04:00",
	)
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := stdout.String()
	if !strings.Contains(result, "line 2") {
		t.Errorf("expected line 2 in output, got:\n%s", result)
	}
	if strings.Contains(result, "line 0") {
		t.Errorf("line 0 should not be in output, got:\n%s", result)
	}
}
