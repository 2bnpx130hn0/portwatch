package history

import (
	"time"
)

// DedupeOptions controls deduplication behaviour.
type DedupeOptions struct {
	// Window is the time window within which identical events are considered
	// duplicates. Zero means exact timestamp match only.
	Window time.Duration
	// Fields controls which fields must match for an entry to be a duplicate.
	// Supported values: "port", "protocol", "action". Defaults to all three.
	Fields []string
}

// dedupeKey represents the fields used to identify a duplicate.
type dedupeKey struct {
	Port     int
	Protocol string
	Action   string
}

// Dedupe removes duplicate entries from a slice of history entries.
// Entries are considered duplicates when they share the same port, protocol,
// and action within the configured time window. The first occurrence is kept.
func Dedupe(entries []Entry, opts DedupeOptions) []Entry {
	if len(entries) == 0 {
		return entries
	}

	fields := opts.Fields
	if len(fields) == 0 {
		fields = []string{"port", "protocol", "action"}
	}

	usePort := containsField(fields, "port")
	useProtocol := containsField(fields, "protocol")
	useAction := containsField(fields, "action")

	// seen maps a dedupeKey to the timestamp of the first occurrence.
	seen := make(map[dedupeKey]time.Time)
	out := make([]Entry, 0, len(entries))

	for _, e := range entries {
		key := dedupeKey{}
		if usePort {
			key.Port = e.Port
		}
		if useProtocol {
			key.Protocol = normProto(e.Protocol)
		}
		if useAction {
			key.Action = e.Action
		}

		first, exists := seen[key]
		if exists {
			if opts.Window == 0 {
				// exact match: skip if timestamps are equal
				if e.Timestamp.Equal(first) {
					continue
				}
			} else {
				diff := e.Timestamp.Sub(first)
				if diff < 0 {
					diff = -diff
				}
				if diff <= opts.Window {
					continue
				}
			}
		}

		seen[key] = e.Timestamp
		out = append(out, e)
	}

	return out
}

func containsField(fields []string, target string) bool {
	for _, f := range fields {
		if f == target {
			return true
		}
	}
	return false
}
