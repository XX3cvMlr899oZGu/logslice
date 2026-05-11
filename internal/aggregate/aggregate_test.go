package aggregate_test

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"github.com/yourorg/logslice/internal/aggregate"
)

func TestNewCounterStartsEmpty(t *testing.T) {
	c := aggregate.New("level")
	if c.Total() != 0 {
		t.Fatalf("expected 0, got %d", c.Total())
	}
	if len(c.Results()) != 0 {
		t.Fatal("expected empty results")
	}
}

func TestAddIncrementsTotal(t *testing.T) {
	c := aggregate.New("level")
	c.Add("info")
	c.Add("info")
	c.Add("error")
	if c.Total() != 3 {
		t.Fatalf("expected 3, got %d", c.Total())
	}
}

func TestResultsSortedByCountDesc(t *testing.T) {
	c := aggregate.New("level")
	for i := 0; i < 5; i++ {
		c.Add("info")
	}
	for i := 0; i < 2; i++ {
		c.Add("error")
	}
	c.Add("warn")

	res := c.Results()
	if len(res) != 3 {
		t.Fatalf("expected 3 results, got %d", len(res))
	}
	if res[0].Value != "info" || res[0].Count != 5 {
		t.Errorf("unexpected first result: %+v", res[0])
	}
	if res[1].Value != "error" || res[1].Count != 2 {
		t.Errorf("unexpected second result: %+v", res[1])
	}
}

func TestEmptyValueRecordedAsPlaceholder(t *testing.T) {
	c := aggregate.New("level")
	c.Add("")
	res := c.Results()
	if len(res) != 1 || res[0].Value != "<empty>" {
		t.Errorf("expected <empty>, got %+v", res)
	}
}

func TestWriteToContainsKeyAndPercent(t *testing.T) {
	c := aggregate.New("status")
	c.Add("200")
	c.Add("200")
	c.Add("500")

	var buf bytes.Buffer
	c.WriteTo(&buf)
	out := buf.String()

	if !strings.Contains(out, "status") {
		t.Error("output missing field name")
	}
	if !strings.Contains(out, "66.7") {
		t.Errorf("output missing percentage: %s", out)
	}
}

func TestAddIsConcurrencySafe(t *testing.T) {
	c := aggregate.New("level")
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Add("info")
		}()
	}
	wg.Wait()
	if c.Total() != 100 {
		t.Fatalf("expected 100, got %d", c.Total())
	}
}
