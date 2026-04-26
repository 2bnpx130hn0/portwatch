package history

import (
	"sort"
	"strings"
	"time"
)

// FlowEdge represents a directional transition between two ports observed within a time window.
type FlowEdge struct {
	FromPort int
	ToPort   int
	Protocol string
	Count    int
	FirstSeen time.Time
	LastSeen  time.Time
}

// FlowOptions controls how port flow transitions are computed.
type FlowOptions struct {
	Protocol string
	Action   string
	Since    time.Time
	Window   time.Duration // max gap between events to form an edge
	MinCount int
}

// BuildFlow detects sequential port transitions within a time window,
// returning directed edges sorted by count descending.
func BuildFlow(entries []Entry, opts FlowOptions) []FlowEdge {
	if opts.Window <= 0 {
		opts.Window = 5 * time.Minute
	}

	filtered := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if opts.Protocol != "" && !strings.EqualFold(e.Protocol, opts.Protocol) {
			continue
		}
		if opts.Action != "" && !strings.EqualFold(e.Action, opts.Action) {
			continue
		}
		filtered = append(filtered, e)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp.Before(filtered[j].Timestamp)
	})

	type edgeKey struct {
		from, to int
		proto    string
	}

	type edgeVal struct {
		count     int
		firstSeen time.Time
		lastSeen  time.Time
	}

	edges := map[edgeKey]*edgeVal{}

	for i := 1; i < len(filtered); i++ {
		prev := filtered[i-1]
		curr := filtered[i]
		if curr.Timestamp.Sub(prev.Timestamp) > opts.Window {
			continue
		}
		if prev.Port == curr.Port {
			continue
		}
		proto := strings.ToLower(curr.Protocol)
		k := edgeKey{from: prev.Port, to: curr.Port, proto: proto}
		if v, ok := edges[k]; ok {
			v.count++
			if curr.Timestamp.After(v.lastSeen) {
				v.lastSeen = curr.Timestamp
			}
		} else {
			edges[k] = &edgeVal{count: 1, firstSeen: prev.Timestamp, lastSeen: curr.Timestamp}
		}
	}

	result := make([]FlowEdge, 0, len(edges))
	for k, v := range edges {
		if opts.MinCount > 0 && v.count < opts.MinCount {
			continue
		}
		result = append(result, FlowEdge{
			FromPort:  k.from,
			ToPort:    k.to,
			Protocol:  k.proto,
			Count:     v.count,
			FirstSeen: v.firstSeen,
			LastSeen:  v.lastSeen,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		return result[i].FromPort < result[j].FromPort
	})
	return result
}
