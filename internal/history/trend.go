package history

import (
	"sort"
	"time"
)

// TrendPoint represents event count in a time window.
type TrendPoint struct {
	Window time.Time
	Count  int
}

// Trend options for computing trends.
type TrendOptions struct {
	BucketSize time.Duration
	Action     string
	Protocol   string
	Since      time.Time
}

// Trend computes a count-over-time trend from history entries.
func Trend(entries []Entry, opts TrendOptions) []TrendPoint {
	if opts.BucketSize <= 0 {
		opts.BucketSize = time.Hour
	}

	buckets := map[time.Time]int{}

	for _, e := range entries {
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if opts.Action != "" && !equalFold(e.Action, opts.Action) {
			continue
		}
		if opts.Protocol != "" && !equalFold(e.Protocol, opts.Protocol) {
			continue
		}
		window := e.Timestamp.Truncate(opts.BucketSize)
		buckets[window]++
	}

	points := make([]TrendPoint, 0, len(buckets))
	for w, c := range buckets {
		points = append(points, TrendPoint{Window: w, Count: c})
	}
	sort.Slice(points, func(i, j int) bool {
		return points[i].Window.Before(points[j].Window)
	})
	return points
}
