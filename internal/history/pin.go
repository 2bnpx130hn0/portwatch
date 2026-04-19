package history

import "time"

// Pin marks entries matching port+protocol as pinned via a label.
func Pin(entries []Entry, port int, protocol string) []Entry {
	out := make([]Entry, len(entries))
	copy(out, entries)
	for i, e := range out {
		if e.Port == port && equalFold(e.Protocol, protocol) {
			if out[i].Labels == nil {
				out[i].Labels = map[string]string{}
			}
			out[i].Labels["pinned"] = "true"
			out[i].Labels["pinned_at"] = time.Now().UTC().Format(time.RFC3339)
		}
	}
	return out
}

// Unpin removes the pinned label from matching entries.
func Unpin(entries []Entry, port int, protocol string) []Entry {
	out := make([]Entry, len(entries))
	copy(out, entries)
	for i, e := range out {
		if e.Port == port && equalFold(e.Protocol, protocol) {
			if out[i].Labels != nil {
				delete(out[i].Labels, "pinned")
				delete(out[i].Labels, "pinned_at")
			}
		}
	}
	return out
}

// FilterPinned returns only entries that have been pinned.
func FilterPinned(entries []Entry) []Entry {
	var out []Entry
	for _, e := range entries {
		if e.Labels != nil && e.Labels["pinned"] == "true" {
			out = append(out, e)
		}
	}
	return out
}
