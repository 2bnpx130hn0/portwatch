package history

import "strings"

// Label sets a named label (key=value) on matching entries.
func Label(entries []Entry, port int, protocol, key, value string) []Entry {
	out := make([]Entry, len(entries))
	for i, e := range entries {
		if e.Port == port && strings.EqualFold(e.Protocol, protocol) {
			if e.Labels == nil {
				e.Labels = map[string]string{}
			} else {
				copy := map[string]string{}
				for k, v := range e.Labels {
					copy[k] = v
				}
				e.Labels = copy
			}
			e.Labels[key] = value
		}
		out[i] = e
	}
	return out
}

// RemoveLabel removes a label key from matching entries.
func RemoveLabel(entries []Entry, port int, protocol, key string) []Entry {
	out := make([]Entry, len(entries))
	for i, e := range entries {
		if e.Port == port && strings.EqualFold(e.Protocol, protocol) {
			if e.Labels != nil {
				copy := map[string]string{}
				for k, v := range e.Labels {
					if k != key {
						copy[k] = v
					}
				}
				e.Labels = copy
			}
		}
		out[i] = e
	}
	return out
}

// FilterByLabel returns entries that have the given key=value label.
func FilterByLabel(entries []Entry, key, value string) []Entry {
	var out []Entry
	for _, e := range entries {
		if v, ok := e.Labels[key]; ok && v == value {
			out = append(out, e)
		}
	}
	return out
}
