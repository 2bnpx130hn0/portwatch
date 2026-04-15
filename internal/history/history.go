package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry represents a single recorded port change event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Proto     string    `json:"proto"`
	Port      int       `json:"port"`
	Action    string    `json:"action"`
	Message   string    `json:"message"`
}

// History manages an append-only log of port change events.
type History struct {
	mu      sync.Mutex
	path    string
	entries []Entry
}

// New loads existing history from path (if any) and returns a History.
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

// Record appends an entry and persists the log to disk.
func (h *History) Record(e Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	h.entries = append(h.entries, e)
	return h.save()
}

// Entries returns a copy of all recorded entries.
func (h *History) Entries() []Entry {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]Entry, len(h.entries))
	copy(out, h.entries)
	return out
}

// save writes the current entries to disk (caller must hold mu).
func (h *History) save() error {
	data, err := json.MarshalIndent(h.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0o644)
}
