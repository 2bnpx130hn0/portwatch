package history

import (
	"sort"
	"time"
)

// HotspotEntry represents a port/protocol pair ranked by activity frequency.
type HotspotEntry struct {
	Port     int
	Protocol string
	Count    int
	Actions  map[string]int
	LastSeen time.Time
}

// HotspotOptions controls filtering for Hotspot detection.
type HotspotOptions struct {
	Since    time.Time
	Action   string
	Protocol string
	TopN     int
}

// Hotspot identifies the most frequently seen port/protocol pairs within the
// given options window, returning up to TopN results ordered by count desc.
func Hotspot(entries []Entry, opts HotspotOptions) []HotspotEntry {
	type key struct {
		port  int
		proto string
	}

	agg := make(map[key]*HotspotEntry)

	for _, e := range entries {
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if opts.Action != "" && !strings.EqualFold(e.Action, opts.Action) {
			continue
		}
		if opts.Protocol != "" && !strings.EqualFold(e.Protocol, opts.Protocol) {
			continue
		}

		k := key{port: e.Port, proto: normProto(e.Protocol)}
		if agg[k] == nil {
			agg[k] = &HotspotEntry{
				Port:     e.Port,
				Protocol: normProto(e.Protocol),
				Actions:  make(map[string]int),
			}
		}
		h := agg[k]
		h.Count++
		h.Actions[strings.ToLower(e.Action)]++
		if e.Timestamp.After(h.LastSeen) {
			h.LastSeen = e.Timestamp
		}
	}

	result := make([]HotspotEntry, 0, len(agg))
	for _, v := range agg {
		result = append(result, *v)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		return result[i].Port < result[j].Port
	})

	n := opts.TopN
	if n <= 0 || n > len(result) {
		n = len(result)
	}
	return result[:n]
}
