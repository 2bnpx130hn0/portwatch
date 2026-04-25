package history

import (
	"math"
	"strings"
	"time"
)

// OutlierResult represents a port/protocol pair that exhibits outlier behaviour
// relative to the overall activity distribution.
type OutlierResult struct {
	Port     int
	Protocol string
	Count    int
	ZScore   float64
	Action   string
}

// OutlierOptions controls how DetectOutliers selects candidates.
type OutlierOptions struct {
	// MinZScore is the minimum absolute z-score to be considered an outlier.
	// Defaults to 2.0 if zero.
	MinZScore float64
	// Since restricts entries to those after this time (zero means no filter).
	Since time.Time
	// Action filters entries to a specific action (empty means all).
	Action string
}

// DetectOutliers identifies port/protocol pairs whose event counts deviate
// significantly from the mean count across all observed pairs.
func DetectOutliers(entries []Entry, opts OutlierOptions) []OutlierResult {
	if opts.MinZScore == 0 {
		opts.MinZScore = 2.0
	}

	// Filter entries.
	var filtered []Entry
	for _, e := range entries {
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if opts.Action != "" && !strings.EqualFold(e.Action, opts.Action) {
			continue
		}
		filtered = append(filtered, e)
	}

	if len(filtered) == 0 {
		return nil
	}

	// Accumulate counts per (port, protocol, action) key.
	type key struct {
		port  int
		proto string
		act   string
	}
	counts := make(map[key]int)
	for _, e := range filtered {
		k := key{port: e.Port, proto: strings.ToLower(e.Protocol), act: strings.ToLower(e.Action)}
		counts[k]++
	}

	// Compute mean and std-dev of counts.
	vals := make([]float64, 0, len(counts))
	for _, c := range counts {
		vals = append(vals, float64(c))
	}
	m, sd := meanStddev(vals)
	if sd == 0 {
		return nil
	}

	// Collect outliers.
	var results []OutlierResult
	for k, c := range counts {
		z := math.Abs((float64(c) - m) / sd)
		if z >= opts.MinZScore {
			results = append(results, OutlierResult{
				Port:     k.port,
				Protocol: k.proto,
				Count:    c,
				ZScore:   math.Round(z*100) / 100,
				Action:   k.act,
			})
		}
	}

	// Sort descending by z-score.
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].ZScore > results[j-1].ZScore; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}
	return results
}
