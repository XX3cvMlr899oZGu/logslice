package annotate_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/annotate"
)

func TestNoOptionsIsIdentity(t *testing.T) {
	a, err := annotate.New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := a.Apply("hello world")
	if got != "hello world" {
		t.Errorf("expected identity, got %q", got)
	}
}

func TestSingleFieldAppended(t *testing.T) {
	a, err := annotate.New(annotate.WithField("host", "srv-1"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := a.Apply("log line")
	want := "log line host=srv-1"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMultipleFieldsAppended(t *testing.T) {
	a, err := annotate.New(
		annotate.WithField("env", "prod"),
		annotate.WithField("region", "us-east-1"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := a.Apply("msg")
	want := "msg env=prod region=us-east-1"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSuffixAppended(t *testing.T) {
	a, err := annotate.New(annotate.WithSuffix("[ANNOTATED]"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := a.Apply("some log")
	want := "some log [ANNOTATED]"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFieldAndSuffixCombined(t *testing.T) {
	a, err := annotate.New(
		annotate.WithField("svc", "api"),
		annotate.WithSuffix("EOF"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := a.Apply("entry")
	want := "entry svc=api EOF"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEmptyKeyReturnsError(t *testing.T) {
	_, err := annotate.New(annotate.WithField("", "value"))
	if err == nil {
		t.Error("expected error for empty key, got nil")
	}
}

func TestApplyEmptyLine(t *testing.T) {
	a, err := annotate.New(annotate.WithField("k", "v"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := a.Apply("")
	want := " k=v"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
