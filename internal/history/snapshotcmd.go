package history

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// SnapshotCmd provides high-level operations for snapshot management.
type SnapshotCmd struct {
	Dir string
}

// NewSnapshotCmd creates a SnapshotCmd rooted at dir.
func NewSnapshotCmd(dir string) *SnapshotCmd {
	return &SnapshotCmd{Dir: dir}
}

// Save captures entries as a named snapshot file.
func (s *SnapshotCmd) Save(name string, entries []Entry) (string, error) {
	if err := os.MkdirAll(s.Dir, 0o755); err != nil {
		return "", err
	}
	if name == "" {
		name = time.Now().UTC().Format("20060102-150405")
	}
	path := filepath.Join(s.Dir, name+".json")
	snap := TakeSnapshot(entries)
	if err := SaveSnapshot(path, snap); err != nil {
		return "", err
	}
	return path, nil
}

// Compare loads two named snapshots and renders their diff.
func (s *SnapshotCmd) Compare(nameA, nameB, format string, w io.Writer) error {
	pathA := filepath.Join(s.Dir, nameA+".json")
	pathB := filepath.Join(s.Dir, nameB+".json")
	a, err := LoadSnapshot(pathA)
	if err != nil {
		return fmt.Errorf("load %s: %w", nameA, err)
	}
	b, err := LoadSnapshot(pathB)
	if err != nil {
		return fmt.Errorf("load %s: %w", nameB, err)
	}
	RenderSnapshot(a, b, format, w)
	return nil
}

// List returns the names of all snapshots in the directory.
func (s *SnapshotCmd) List() ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(s.Dir, "*.json"))
	if err != nil {
		return nil, err
	}
	names := make([]string, len(matches))
	for i, m := range matches {
		base := filepath.Base(m)
		names[i] = base[:len(base)-5]
	}
	return names, nil
}
