package history

import "time"

// Flag marks entries matching port/protocol as flagged for review.
func Flag(entries []Entry, port int, protocol string) []Entry {
	updated := make([]Entry, len(entries))
	for i, e := range entries {
		if e.Port == port && equalFold(e.Protocol, protocol) {
			if e.Labels == nil {
				e.Labels = map[string]string{}
			}
			e.Labels["flagged"] = "true"
			e.Labels["flagged_at"] = time.Now().UTC().Format(time.RFC3339)
		}
		updated[i] = e
	}
	return updated
}

// Unflag removes the flagged label from matching entries.
func Unflag(entries []Entry, port int, protocol string) []Entry {
	updated := make([]Entry, len(entries))
	for i, e := range entries {
		if e.Port == port && equalFold(e.Protocol, protocol) {
			if e.Labels != nil {
				delete(e.Labels, "flagged")
				delete(e.Labels, "flagged_at")
			}
		}
		updated[i] = e
	}
	return updated
}

// FilterFlagged returns only entries that have been flagged.
func FilterFlagged(entries []Entry) []Entry {
	var out []Entry
	for _, e := range entries {
		if e.Labels != nil && e.Labels["flagged"] == "true" {
			out = append(out, e)
		}
	}
	return out
}
