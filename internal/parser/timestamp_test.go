package parser

import (
	"testing"
	"time"
)

func TestParseTimestamp(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantErr bool
		wantUTC string // expected UTC time in RFC3339 (empty == skip check)
	}{
		{
			name:    "RFC3339Nano",
			line:    "2024-03-15T12:34:56.789012345Z INFO server started",
			wantUTC: "2024-03-15T12:34:56Z",
		},
		{
			name:    "RFC3339 with offset",
			line:    "2024-03-15T12:34:56+02:00 WARN disk usage high",
			wantUTC: "2024-03-15T10:34:56Z",
		},
		{
			name:    "space-separated datetime",
			line:    "2024-03-15 12:34:56 ERROR connection refused",
			wantUTC: "2024-03-15T12:34:56Z",
		},
		{
			name:    "common log format",
			line:    "10/Oct/2000:13:55:36 -0700 GET /index.html HTTP/1.0",
			wantUTC: "2000-10-10T20:55:36Z",
		},
		{
			name:    "empty line",
			line:    "",
			wantErr: true,
		},
		{
			name:    "no timestamp",
			line:    "this line has no timestamp at all",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseTimestamp(tc.line)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error but got time %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantUTC != "" {
				want, _ := time.Parse(time.RFC3339, tc.wantUTC)
				if !got.UTC().Truncate(time.Second).Equal(want.UTC()) {
					t.Errorf("got %v, want %v", got.UTC(), want.UTC())
				}
			}
		})
	}
}
