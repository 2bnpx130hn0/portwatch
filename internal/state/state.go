package state

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Snapshot holds the port state at a point in time.
type Snapshot struct {
	Timestamp time.Time        `json:"timestamp"`
	Ports     map[string][]int `json:"ports"` // protocol -> list of ports
}

// Store persists and retrieves port snapshots.
type Store struct {
	mu       sync.RWMutex
	filePath string
	current  *Snapshot
}

// New creates a new Store backed by the given file path.
func New(filePath string) *Store {
	return &Store{filePath: filePath}
}

// Load reads the last snapshot from disk. Returns nil if none exists.
func (s *Store) Load() (*Snapshot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, err
	}
	s.current = &snap
	return &snap, nil
}

// Save writes the given snapshot to disk and updates the in-memory state.
func (s *Store) Save(snap *Snapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(s.filePath, data, 0o600); err != nil {
		return err
	}
	s.current = snap
	return nil
}

// Current returns the in-memory snapshot without hitting disk.
func (s *Store) Current() *Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.current
}

// Diff computes ports added and removed between two snapshots.
func Diff(prev, next *Snapshot) (added, removed map[string][]int) {
	added = make(map[string][]int)
	removed = make(map[string][]int)

	for proto, ports := range next.Ports {
		prevPorts := portSet(prev.Ports[proto])
		for _, p := range ports {
			if !prevPorts[p] {
				added[proto] = append(added[proto], p)
			}
		}
	}
	for proto, ports := range prev.Ports {
		nextPorts := portSet(next.Ports[proto])
		for _, p := range ports {
			if !nextPorts[p] {
				removed[proto] = append(removed[proto], p)
			}
		}
	}
	return added, removed
}

func portSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
