package history

import (
	"sort"
	"time"
)

// MergeOptions controls how two entry slices are merged.
type MergeOptions struct {
	DeduplicateWindow time.Duration // entries within this window with same port/protocol/action are considered duplicates
}

// Merge combines two slices of Entry, optionally deduplicating entries that
// fall within a time window and share the same port, protocol, and action.
// The result is sorted ascending by timestamp.
func Merge(a, b []Entry, opts MergeOptions) []Entry {
	combined := make([]Entry, 0, len(a)+len(b))
	combined = append(combined, a...)
	combined = append(combined, b...)

	sort.Slice(combined, func(i, j int) bool {
		return combined[i].Timestamp.Before(combined[j].Timestamp)
	})

	if opts.DeduplicateWindow <= 0 {
		return combined
	}

	return deduplicate(combined, opts.DeduplicateWindow)
}

func deduplicate(entries []Entry, window time.Duration) []Entry {
	if len(entries) == 0 {
		return entries
	}

	result := []Entry{entries[0]}

	for i := 1; i < len(entries); i++ {
		e := entries[i]
		last := result[len(result)-1]

		if e.Port == last.Port &&
			stringsEqualFold(e.Protocol, last.Protocol) &&
			stringsEqualFold(e.Action, last.Action) &&
			e.Timestamp.Sub(last.Timestamp) <= window {
			continue
		}

		result = append(result, e)
	}

	return result
}

func stringsEqualFold(a, b string) bool {
	return equalFold(a, b)
}
