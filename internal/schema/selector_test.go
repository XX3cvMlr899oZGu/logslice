package schema_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/schema"
)

func TestParseSelectorEquals(t *testing.T) {
	s, err := schema.ParseSelector("level==error")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.String() != "level==error" {
		t.Errorf("got %q", s.String())
	}
}

func TestParseSelectorNotEquals(t *testing.T) {
	s, err := schema.ParseSelector("level!=debug")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.String() != "level!=debug" {
		t.Errorf("got %q", s.String())
	}
}

func TestParseSelectorInvalid(t *testing.T) {
	_, err := schema.ParseSelector("levelinfo")
	if err == nil {
		t.Fatal("expected error for invalid selector")
	}
}

func TestParseSelectorEmptyField(t *testing.T) {
	_, err := schema.ParseSelector("==value")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestSelectorMatchEquals(t *testing.T) {
	s, _ := schema.ParseSelector("level==error")
	fields := schema.Fields{"level": "error", "msg": "boom"}
	if !s.Match(fields) {
		t.Error("expected match")
	}
}

func TestSelectorNoMatchEquals(t *testing.T) {
	s, _ := schema.ParseSelector("level==error")
	fields := schema.Fields{"level": "info"}
	if s.Match(fields) {
		t.Error("expected no match")
	}
}

func TestSelectorMatchNotEquals(t *testing.T) {
	s, _ := schema.ParseSelector("level!=debug")
	fields := schema.Fields{"level": "info"}
	if !s.Match(fields) {
		t.Error("expected match")
	}
}

func TestSelectorMissingFieldNotEquals(t *testing.T) {
	s, _ := schema.ParseSelector("env!=prod")
	fields := schema.Fields{"level": "info"}
	if !s.Match(fields) {
		t.Error("missing field should satisfy !=")
	}
}

func TestSelectorMissingFieldEquals(t *testing.T) {
	s, _ := schema.ParseSelector("env==prod")
	fields := schema.Fields{"level": "info"}
	if s.Match(fields) {
		t.Error("missing field should not satisfy ==")
	}
}
