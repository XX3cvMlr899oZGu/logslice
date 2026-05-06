package filter

import "fmt"

// UnknownLevelError is returned when an unrecognised level name is provided.
type UnknownLevelError struct {
	Name string
}

func (e *UnknownLevelError) Error() string {
	return fmt.Sprintf("filter: unknown log level %q; valid values are debug, info, warn, error", e.Name)
}
