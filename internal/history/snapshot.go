package history

import (
	"encoding/json"
	"os"
	"time"
)

// Snapshot captures a point-in-time view of history entries.
type Snapshot struct {
	TakenAt time.Time `json:"taken_at"`
	Entries []Entry   `json:"entries"`
}

// TakeSnapshot creates a Snapshot from the given entries.
func TakeSnapshot(entries []Entry) Snapshot {
	copied := make([]Entry, len(entries))
	copy(copied, entries)
	return Snapshot{
		TakenAt: time.Now().UTC(),
		Entries: copied,
	}
}

// SaveSnapshot writes a Snapshot to a JSON file at path.
func SaveSnapshot(path string, snap Snapshot) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(snap)
}

// LoadSnapshot reads a Snapshot from a JSON file at path.
func LoadSnapshot(path string) (Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return Snapshot{}, err
	}
	defer f.Close()
	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}

// SnapshotDiff returns entries present in snap but not in baseline, and vice versa.
func SnapshotDiff(baseline, snap Snapshot) (added []Entry, removed []Entry) {
	index := func(entries []Entry) map[string]struct{} {
		m := make(map[string]struct{}, len(entries))
		for _, e := range entries {
			m[e.Protocol+":"+itoa(e.Port)] = struct{}{}
		}
		return m
	}
	baseIdx := index(baseline.Entries)
	snapIdx := index(snap.Entries)
	for _, e := range snap.Entries {
		if _, ok := baseIdx[e.Protocol+":"+itoa(e.Port)]; !ok {
			added = append(added, e)
		}
	}
	for _, e := range baseline.Entries {
		if _, ok := snapIdx[e.Protocol+":"+itoa(e.Port)]; !ok {
			removed = append(removed, e)
		}
	}
	return
}
