package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry represents a single recorded port event.
type Entry struct {
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}

// History manages a persistent log of port events.
type History struct {
	mu      sync.Mutex
	path    string
	entries []Entry
}

// New loads history from path, or starts empty if the file does not exist.
func New(path string) (*History, error) {
	h := &History{path: path}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return h, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &h.entries); err != nil {
		return nil, err
	}
	return h, nil
}

// Record appends an entry, setting Timestamp if zero, and persists to disk.
func (h *History) Record(e Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	h.entries = append(h.entries, e)
	return h.persist()
}

// All returns a copy of all entries.
func (h *History) All() []Entry {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]Entry, len(h.entries))
	copy(out, h.entries)
	return out
}

// persist writes entries to disk; caller must hold h.mu.
func (h *History) persist() error {
	data, err := json.MarshalIndent(h.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0o644)
}
