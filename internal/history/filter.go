package history

import "time"

// Filter holds criteria for querying history entries.
type Filter struct {
	Protocol string
	Port     int
	Action   string
	Since    time.Time
	Limit    int
}

// Apply returns entries matching all non-zero filter fields.
func (f Filter) Apply(entries []Entry) []Entry {
	var out []Entry
	for _, e := range entries {
		if f.Protocol != "" && !strings.EqualFold(e.Protocol, f.Protocol) {
			continue
		}
		if f.Port != 0 && e.Port != f.Port {
			continue
		}
		if f.Action != "" && !strings.EqualFold(e.Action, f.Action) {
			continue
		}
		if !f.Since.IsZero() && e.Timestamp.Before(f.Since) {
			continue
		}
		out = append(out, e)
		if f.Limit > 0 && len(out) >= f.Limit {
			break
		}
	}
	return out
}
