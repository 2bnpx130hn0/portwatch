package history

import (
	"encoding/json"
	"os"
	"time"
)

// Baseline represents a saved reference set of history entries.
type Baseline struct {
	CreatedAt time.Time `json:"created_at"`
	Entries   []Entry   `json:"entries"`
}

// SaveBaseline writes the current entries as a baseline to path.
func SaveBaseline(entries []Entry, path string) error {
	b := Baseline{
		CreatedAt: time.Now().UTC(),
		Entries:   entries,
	}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadBaseline reads a baseline from path.
func LoadBaseline(path string) (Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Baseline{}, err
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return Baseline{}, err
	}
	return b, nil
}

// BaselineDiff compares current entries against a baseline.
// Returns entries added (in current but not baseline) and removed (in baseline but not current).
func BaselineDiff(baseline, current []Entry) (added, removed []Entry) {
	key := func(e Entry) string {
		return e.Protocol + ":" + itoa(e.Port)
	}

	baseMap := make(map[string]struct{}, len(baseline))
	for _, e := range baseline {
		baseMap[key(e)] = struct{}{}
	}
	currMap := make(map[string]struct{}, len(current))
	for _, e := range current {
		currMap[key(e)] = struct{}{}
	}

	for _, e := range current {
		if _, ok := baseMap[key(e)]; !ok {
			added = append(added, e)
		}
	}
	for _, e := range baseline {
		if _, ok := currMap[key(e)]; !ok {
			removed = append(removed, e)
		}
	}
	return
}
