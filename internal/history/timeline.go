package history

import (
	"sort"
	"time"
)

// TimelineEntry represents aggregated activity within a time bucket.
type TimelineEntry struct {
	Bucket   time.Time
	Total    int
	ByAction map[string]int
}

// TimelineOptions controls how the timeline is bucketed.
type TimelineOptions struct {
	BucketSize time.Duration // e.g. time.Hour, 24*time.Hour
	Since      time.Time
	Until      time.Time
}

// Timeline groups history entries into time buckets for trend analysis.
func Timeline(entries []Entry, opts TimelineOptions) []TimelineEntry {
	if opts.BucketSize <= 0 {
		opts.BucketSize = time.Hour
	}

	buckets := map[time.Time]*TimelineEntry{}

	for _, e := range entries {
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if !opts.Until.IsZero() && e.Timestamp.After(opts.Until) {
			continue
		}
		bucket := e.Timestamp.Truncate(opts.BucketSize)
		if _, ok := buckets[bucket]; !ok {
			buckets[bucket] = &TimelineEntry{
				Bucket:   bucket,
				ByAction: map[string]int{},
			}
		}
		buckets[bucket].Total++
		buckets[bucket].ByAction[e.Action]++
	}

	result := make([]TimelineEntry, 0, len(buckets))
	for _, te := range buckets {
		result = append(result, *te)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Bucket.Before(result[j].Bucket)
	})
	return result
}
