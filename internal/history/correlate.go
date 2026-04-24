package history

import (
	"sort"
	"strings"
	"time"
)

// Correlation holds a group of entries that share a common pattern
// within a short time window, suggesting related activity.
type Correlation struct {
	Port     int
	Protocol string
	Entries  []Entry
	Score    float64
}

// CorrelateOptions controls how correlation is performed.
type CorrelateOptions struct {
	// Window is the time range within which entries are considered related.
	Window time.Duration
	// MinEntries is the minimum number of entries required to form a correlation.
	MinEntries int
	// Action filters entries to a specific action (optional).
	Action string
}

// Correlate groups entries by port+protocol that appear together within
// a sliding time window, returning groups ordered by score descending.
func Correlate(entries []Entry, opts CorrelateOptions) []Correlation {
	if opts.Window <= 0 {
		opts.Window = 5 * time.Minute
	}
	if opts.MinEntries < 2 {
		opts.MinEntries = 2
	}

	filtered := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if opts.Action == "" || strings.EqualFold(e.Action, opts.Action) {
			filtered = append(filtered, e)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp.Before(filtered[j].Timestamp)
	})

	type key struct {
		port     int
		protocol string
	}

	groups := map[key][]Entry{}
	for i, anchor := range filtered {
		k := key{port: anchor.Port, protocol: strings.ToLower(anchor.Protocol)}
		if _, seen := groups[k]; seen {
			continue
		}
		var cluster []Entry
		for j := i; j < len(filtered); j++ {
			if filtered[j].Port == anchor.Port &&
				strings.EqualFold(filtered[j].Protocol, anchor.Protocol) &&
				filtered[j].Timestamp.Sub(anchor.Timestamp) <= opts.Window {
				cluster = append(cluster, filtered[j])
			}
		}
		if len(cluster) >= opts.MinEntries {
			groups[k] = cluster
		}
	}

	result := make([]Correlation, 0, len(groups))
	for k, g := range groups {
		result = append(result, Correlation{
			Port:     k.port,
			Protocol: k.protocol,
			Entries:  g,
			Score:    float64(len(g)),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Score != result[j].Score {
			return result[i].Score > result[j].Score
		}
		return result[i].Port < result[j].Port
	})
	return result
}
