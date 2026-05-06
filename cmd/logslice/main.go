package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/logslice/internal/slicer"
)

const timeLayout = "2006-01-02T15:04:05"

func main() {
	var (
		filePath  = flag.String("f", "", "Path to the log file (required)")
		startStr  = flag.String("start", "", "Start timestamp, format: 2006-01-02T15:04:05 (required)")
		endStr    = flag.String("end", "", "End timestamp, format: 2006-01-02T15:04:05 (required)")
		outPath   = flag.String("o", "", "Output file path (default: stdout)")
	)
	flag.Parse()

	if *filePath == "" || *startStr == "" || *endStr == "" {
		fmt.Fprintln(os.Stderr, "error: -f, -start, and -end are required")
		flag.Usage()
		os.Exit(1)
	}

	start, err := time.Parse(timeLayout, *startStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: invalid -start value %q: %v\n", *startStr, err)
		os.Exit(1)
	}

	end, err := time.Parse(timeLayout, *endStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: invalid -end value %q: %v\n", *endStr, err)
		os.Exit(1)
	}

	var out *os.File
	if *outPath == "" {
		out = os.Stdout
	} else {
		out, err = os.Create(*outPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: cannot create output file %q: %v\n", *outPath, err)
			os.Exit(1)
		}
		defer out.Close()
	}

	opts := slicer.Options{
		FilePath: *filePath,
		Start:    start,
		End:      end,
		Writer:   out,
	}

	if err := slicer.Slice(opts); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
