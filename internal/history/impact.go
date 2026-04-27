package history

import (
	"sort"
	"strings"
	"time"
)

// ImpactEntry summarizes the impact of a port/protocol pair.
type ImpactEntry struct {
	Port     int
	Protocol string
	Total    int
	Alerts   int
	Warnings int
	Allowed  int
	FirstSeen time.Time
	LastSeen  time.Time
	Score    float64
}

// ImpactOptions controls how Impact analysis is performed.
type ImpactOptions struct {
	Since    time.Time
	Action   string
	TopN     int
}

// Impact ranks port/protocol pairs by their overall event impact.
// Score = alerts*3 + warnings*1.5 + allowed*0.5, boosted by recency.
func Impact(entries []Entry, opts ImpactOptions) []ImpactEntry {
	type key struct {
		port  int
		proto string
	}

	type bucket struct {
		total, alerts, warnings, allowed int
		first, last                      time.Time
	}

	m := make(map[key]*bucket)
	now := time.Now()

	for _, e := range entries {
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if opts.Action != "" && !strings.EqualFold(e.Action, opts.Action) {
			continue
		}
		k := key{port: e.Port, proto: strings.ToLower(e.Protocol)}
		b, ok := m[k]
		if !ok {
			b = &bucket{first: e.Timestamp, last: e.Timestamp}
			m[k] = b
		}
		b.total++
		switch strings.ToLower(e.Action) {
		case "alert":
			b.alerts++
		case "warn":
			b.warnings++
		default:
			b.allowed++
		}
		if e.Timestamp.Before(b.first) {
			b.first = e.Timestamp
		}
		if e.Timestamp.After(b.last) {
			b.last = e.Timestamp
		}
	}

	result := make([]ImpactEntry, 0, len(m))
	for k, b := range m {
		base := float64(b.alerts)*3.0 + float64(b.warnings)*1.5 + float64(b.allowed)*0.5
		age := now.Sub(b.last).Hours()
		recency := 1.0
		if age < 1 {
			recency = 2.0
		} else if age < 24 {
			recency = 1.5
		}
		result = append(result, ImpactEntry{
			Port:      k.port,
			Protocol:  k.proto,
			Total:     b.total,
			Alerts:    b.alerts,
			Warnings:  b.warnings,
			Allowed:   b.allowed,
			FirstSeen: b.first,
			LastSeen:  b.last,
			Score:     base * recency,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})

	if opts.TopN > 0 && len(result) > opts.TopN {
		result = result[:opts.TopN]
	}
	return result
}
