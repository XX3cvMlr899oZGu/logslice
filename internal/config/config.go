// Package config parses and validates CLI options for logslice.
package config

import (
	"errors"
	"flag"
	"time"
)

// Options holds all runtime configuration derived from CLI flags.
type Options struct {
	FilePath  string
	From      time.Time
	To        time.Time
	Level     string
	Keyword   string
	Output    string
	Stats     bool
	Progress  bool
}

// ErrMissingFile is returned when no input file is specified.
var ErrMissingFile = errors.New("config: --file is required")

// ErrMissingFrom is returned when --from is not provided.
var ErrMissingFrom = errors.New("config: --from is required")

// ErrMissingTo is returned when --to is not provided.
var ErrMissingTo = errors.New("config: --to is required")

// ErrInvalidRange is returned when --from is not before --to.
var ErrInvalidRange = errors.New("config: --from must be before --to")

const tsLayout = "2006-01-02T15:04:05"

// Parse reads flags from the provided FlagSet and args, returning validated Options.
func Parse(fs *flag.FlagSet, args []string) (*Options, error) {
	file := fs.String("file", "", "path to log file")
	from := fs.String("from", "", "start timestamp (2006-01-02T15:04:05)")
	to := fs.String("to", "", "end timestamp (2006-01-02T15:04:05)")
	level := fs.String("level", "", "minimum log level filter (DEBUG|INFO|WARN|ERROR)")
	keyword := fs.String("keyword", "", "keyword substring filter")
	output := fs.String("output", "", "output file path (default: stdout)")
	stats := fs.Bool("stats", false, "print statistics after slicing")
	progress := fs.Bool("progress", false, "show progress during slicing")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if *file == "" {
		return nil, ErrMissingFile
	}
	if *from == "" {
		return nil, ErrMissingFrom
	}
	if *to == "" {
		return nil, ErrMissingTo
	}

	parsedFrom, err := time.Parse(tsLayout, *from)
	if err != nil {
		return nil, errors.New("config: invalid --from: " + err.Error())
	}
	parsedTo, err := time.Parse(tsLayout, *to)
	if err != nil {
		return nil, errors.New("config: invalid --to: " + err.Error())
	}
	if !parsedFrom.Before(parsedTo) {
		return nil, ErrInvalidRange
	}

	return &Options{
		FilePath: *file,
		From:     parsedFrom,
		To:       parsedTo,
		Level:    *level,
		Keyword:  *keyword,
		Output:   *output,
		Stats:    *stats,
		Progress: *progress,
	}, nil
}
