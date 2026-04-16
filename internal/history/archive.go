package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Archiver moves old entries out of the active history into dated archive files.
type Archiver struct {
	history *History
	archiveDir string
}

// NewArchiver creates an Archiver that writes archives to archiveDir.
func NewArchiver(h *History, archiveDir string) *Archiver {
	return &Archiver{history: h, archiveDir: archiveDir}
}

// Archive moves entries older than maxAge into a dated archive file and
// removes them from the active history. Returns the number of entries archived.
func (a *Archiver) Archive(maxAge time.Duration) (int, error) {
	if err := os.MkdirAll(a.archiveDir, 0o755); err != nil {
		return 0, fmt.Errorf("archive: mkdir: %w", err)
	}

	cutoff := time.Now().Add(-maxAge)
	var keep, old []Entry

	for _, e := range a.history.entries {
		if e.Timestamp.Before(cutoff) {
			old = append(old, e)
		} else {
			keep = append(keep, e)
		}
	}

	if len(old) == 0 {
		return 0, nil
	}

	filename := fmt.Sprintf("archive-%s.json", time.Now().UTC().Format("20060102-150405"))
	path := filepath.Join(a.archiveDir, filename)

	f, err := os.Create(path)
	if err != nil {
		return 0, fmt.Errorf("archive: create file: %w", err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(old); err != nil {
		return 0, fmt.Errorf("archive: encode: %w", err)
	}

	a.history.entries = keep
	if err := a.history.persist(); err != nil {
		return 0, fmt.Errorf("archive: persist: %w", err)
	}

	return len(old), nil
}
