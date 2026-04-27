package history

import (
	"sort"
	"time"
)

// RankOptions controls how entries are ranked.
type RankOptions struct {
	Action   string
	Protocol string
	Since    time.Time
	TopN     int
}

// RankResult holds a ranked port/protocol pair with its computed score.
type RankResult struct {
	Port     int
	Protocol string
	Count    int
	Score    float64
}

// Rank scores entries by frequency and recency, returning the top results.
func Rank(entries []Entry, opts RankOptions) []RankResult {
	type key struct {
		port  int
		proto string
	}

	type accumulator struct {
		count  int
		recent time.Time
	}

	now := time.Now()
	acc := make(map[key]*accumulator)

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
		k := key{port: e.Port, proto: strings.ToLower(e.Protocol)}
		if acc[k] == nil {
			acc[k] = &accumulator{}
		}
		acc[k].count++
		if e.Timestamp.After(acc[k].recent) {
			acc[k].recent = e.Timestamp
		}
	}

	results := make([]RankResult, 0, len(acc))
	for k, a := range acc {
		age := now.Sub(a.recent).Hours()
		recencyBoost := 1.0 / (1.0 + age/24.0)
		score := float64(a.count) * (1.0 + recencyBoost)
		results = append(results, RankResult{
			Port:     k.port,
			Protocol: k.proto,
			Count:    a.count,
			Score:    score,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if opts.TopN > 0 && len(results) > opts.TopN {
		results = results[:opts.TopN]
	}
	return results
}
