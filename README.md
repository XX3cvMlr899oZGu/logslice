# logslice

Fast log file slicer that extracts time-range segments from large structured log files.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

## Usage

Extract log entries between two timestamps:

```bash
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" --file app.log
```

Pipe output to a new file:

```bash
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" --file app.log > slice.log
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--file` | Path to the log file | required |
| `--from` | Start timestamp (RFC3339) | required |
| `--to` | End timestamp (RFC3339) | required |
| `--format` | Timestamp format in logs | `RFC3339` |
| `--field` | Timestamp field name (JSON logs) | `time` |

### Supported Log Formats

- JSON structured logs
- Common log format (CLF)
- Custom formats via `--format` flag

## Building from Source

```bash
git clone https://github.com/yourusername/logslice.git
cd logslice
go build ./...
```

## License

MIT © 2024