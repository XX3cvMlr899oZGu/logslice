package mask_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/mask"
)

var sink string

func BenchmarkApplySingleField(b *testing.B) {
	m := mask.New([]string{"password"})
	line := `{"ts":"2024-01-01T10:00:00Z","level":"info","password":"supersecret","user":"alice"}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = m.Apply(line)
	}
}

func BenchmarkApplyMultipleFields(b *testing.B) {
	m := mask.New([]string{"password", "token", "apikey", "secret"})
	line := `{"ts":"2024-01-01T10:00:00Z","password":"pw","token":"tok","apikey":"ak","secret":"s","msg":"ok"}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = m.Apply(line)
	}
}

func BenchmarkApplyNoMatch(b *testing.B) {
	m := mask.New([]string{"password"})
	line := `{"ts":"2024-01-01T10:00:00Z","level":"info","msg":"nothing sensitive here at all"}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = m.Apply(line)
	}
}
