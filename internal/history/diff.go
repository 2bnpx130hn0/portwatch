package history

import "time"

// DiffEntry represents a change between two history snapshots.
type DiffEntry struct {
	Port     int
	Protocol string
	Action   string
	Added    bool
	Removed  bool
	At       time.Time
}

// Diff compares two slices of Entry and returns what was added or removed
// between the older (before) and newer (after) sets.
func Diff(before, after []Entry) []DiffEntry {
	key := func(e Entry) string {
		return e.Protocol + ":" + itoa(e.Port)
	}

	beforeMap := make(map[string]Entry, len(before))
	for _, e := range before {
		beforeMap[key(e)] = e
	}

	afterMap := make(map[string]Entry, len(after))
	for _, e := range after {
		afterMap[key(e)] = e
	}

	var diffs []DiffEntry

	for k, e := range afterMap {
		if _, exists := beforeMap[k]; !exists {
			diffs = append(diffs, DiffEntry{
				Port:     e.Port,
				Protocol: e.Protocol,
				Action:   e.Action,
				Added:    true,
				At:       e.Timestamp,
			})
		}
	}

	for k, e := range beforeMap {
		if _, exists := afterMap[k]; !exists {
			diffs = append(diffs, DiffEntry{
				Port:     e.Port,
				Protocol: e.Protocol,
				Action:   e.Action,
				Removed:  true,
				At:       e.Timestamp,
			})
		}
	}

	return diffs
}
