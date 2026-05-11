// Package aggregate provides field-level aggregation over log lines,
// counting occurrences of distinct values for a given key.
package aggregate

import (
	"fmt"
	"io"
	"sort"
	"sync"
)

// Counter accumulates value frequencies for a single field key.
type Counter struct {
	mu     sync.Mutex
	key    string
	counts map[string]int
	total  int
}

// New returns a Counter that tracks distinct values for key.
func New(key string) *Counter {
	return &Counter{
		key:    key,
		counts: make(map[string]int),
	}
}

// Add records the value observed for the counter's key.
// An empty value is recorded under the label "<empty>".
func (c *Counter) Add(value string) {
	if value == "" {
		value = "<empty>"
	}
	c.mu.Lock()
	c.counts[value]++
	c.total++
	c.mu.Unlock()
}

// Total returns the number of times Add has been called.
func (c *Counter) Total() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.total
}

// Results returns a snapshot of value → count pairs sorted by count descending.
func (c *Counter) Results() []Result {
	c.mu.Lock()
	defer c.mu.Unlock()

	out := make([]Result, 0, len(c.counts))
	for v, n := range c.counts {
		out = append(out, Result{Value: v, Count: n})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Value < out[j].Value
	})
	return out
}

// WriteTo writes a human-readable summary to w.
func (c *Counter) WriteTo(w io.Writer) {
	c.mu.Lock()
	key := c.key
	total := c.total
	c.mu.Unlock()

	fmt.Fprintf(w, "field: %s  total: %d\n", key, total)
	for _, r := range c.Results() {
		pct := 0.0
		if total > 0 {
			pct = float64(r.Count) / float64(total) * 100
		}
		fmt.Fprintf(w, "  %-30s %6d  (%5.1f%%)\n", r.Value, r.Count, pct)
	}
}

// Result holds a single value and its observed count.
type Result struct {
	Value string
	Count int
}
