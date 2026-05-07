// Package checkpoint provides resumable slice state persistence.
// It records the last successfully processed byte offset so that
// a interrupted run can continue from where it left off.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// State holds the persisted progress of a slice operation.
type State struct {
	File   string    `json:"file"`
	Offset int64     `json:"offset"`
	Lines  int64     `json:"lines"`
	SavedAt time.Time `json:"saved_at"`
}

// Save writes state to the given path as JSON, overwriting any
// existing file atomically via a temp-file rename.
func Save(path string, s State) error {
	s.SavedAt = time.Now().UTC()
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// Load reads a previously saved State from path.
// Returns ErrNotFound if the file does not exist.
func Load(path string) (State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return State{}, ErrNotFound
		}
		return State{}, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return State{}, err
	}
	return s, nil
}

// Remove deletes the checkpoint file at path.
// Returns nil if the file does not exist.
func Remove(path string) error {
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
