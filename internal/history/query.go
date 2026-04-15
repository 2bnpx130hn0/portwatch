package history

import (
	"time"
)

// Filter holds optional criteria for querying history entries.
type Filter struct {
	Protocol string
	Port     int
	Action   string
	Since    time.Time
	Limit    int
}

// Query returns history entries matching the given filter.
// Zero-value fields in the filter are treated as "no constraint".
func (h *History) Query(f Filter) []Entry {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var results []Entry
	for _, e := range h.entries {
		if f.Protocol != "" && e.Protocol != f.Protocol {
			continue
		}
		if f.Port != 0 && e.Port != f.Port {
			continue
		}
		if f.Action != "" && e.Action != f.Action {
			continue
		}
		if !f.Since.IsZero() && e.Timestamp.Before(f.Since) {
			continue
		}
		results = append(results, e)
		if f.Limit > 0 && len(results) >= f.Limit {
			break
		}
	}
	return results
}

// Latest returns the n most recent history entries.
// If n <= 0 all entries are returned.
func (h *History) Latest(n int) []Entry {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if n <= 0 || n >= len(h.entries) {
		copy := make([]Entry, len(h.entries))
		copy2 := append(copy[:0], h.entries...)
		return copy2
	}
	start := len(h.entries) - n
	return append([]Entry{}, h.entries[start:]...)
}
