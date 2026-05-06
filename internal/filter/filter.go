// Package filter provides log line filtering capabilities based on
// log level and keyword matching for use with logslice.
package filter

import "strings"

// Level represents a log severity level.
type Level int

const (
	LevelAll Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
)

// levelNames maps string representations to Level constants.
var levelNames = map[string]Level{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
}

// Options holds configuration for a Filter.
type Options struct {
	// MinLevel filters out lines below this severity. Empty means accept all.
	MinLevel string
	// Keyword filters lines that do not contain this substring. Empty means accept all.
	Keyword string
}

// Filter decides whether log lines should be included in output.
type Filter struct {
	minLevel Level
	keyword  string
}

// New creates a Filter from the given Options.
// Returns an error if MinLevel is unrecognised.
func New(opts Options) (*Filter, error) {
	f := &Filter{keyword: opts.Keyword}
	if opts.MinLevel == "" {
		f.minLevel = LevelAll
		return f, nil
	}
	lvl, ok := levelNames[strings.ToLower(opts.MinLevel)]
	if !ok {
		return nil, &UnknownLevelError{Name: opts.MinLevel}
	}
	f.minLevel = lvl
	return f, nil
}

// Accept returns true if the given log line passes all active filters.
func (f *Filter) Accept(line string) bool {
	if f.keyword != "" && !strings.Contains(line, f.keyword) {
		return false
	}
	if f.minLevel == LevelAll {
		return true
	}
	return lineLevel(line) >= f.minLevel
}

// lineLevel attempts to detect the severity level embedded in a log line.
func lineLevel(line string) Level {
	upper := strings.ToUpper(line)
	switch {
	case strings.Contains(upper, "ERROR"):
		return LevelError
	case strings.Contains(upper, "WARN"):
		return LevelWarn
	case strings.Contains(upper, "INFO"):
		return LevelInfo
	case strings.Contains(upper, "DEBUG"):
		return LevelDebug
	default:
		return LevelAll
	}
}
