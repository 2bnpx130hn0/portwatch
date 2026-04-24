package history

import (
	"sort"
	"strings"
	"time"
)

// PatternResult holds a detected recurring pattern for a port/protocol pair.
type PatternResult struct {
	Port     int
	Protocol string
	Action   string
	Count    int
	FirstSeen time.Time
	LastSeen  time.Time
	AvgIntervalSecs float64
}

// PatternOptions controls pattern detection behaviour.
type PatternOptions struct {
	MinOccurrences int
	Since          time.Time
	Action         string // optional filter
	Protocol       string // optional filter
}

// DetectPatterns identifies port/protocol/action combinations that recur at
// least MinOccurrences times within the entries, returning results ordered by
// count descending.
func DetectPatterns(entries []Entry, opts PatternOptions) []PatternResult {
	if opts.MinOccurrences <= 0 {
		opts.MinOccurrences = 2
	}

	type key struct {
		port     int
		protocol string
		action   string
	}

	type bucket struct {
		times []time.Time
	}

	buckets := make(map[key]*bucket)

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
		k := key{port: e.Port, protocol: strings.ToLower(e.Protocol), action: strings.ToLower(e.Action)}
		if buckets[k] == nil {
			buckets[k] = &bucket{}
		}
		buckets[k].times = append(buckets[k].times, e.Timestamp)
	}

	var results []PatternResult
	for k, b := range buckets {
		if len(b.times) < opts.MinOccurrences {
			continue
		}
		sort.Slice(b.times, func(i, j int) bool { return b.times[i].Before(b.times[j]) })
		avg := 0.0
		if len(b.times) > 1 {
			total := b.times[len(b.times)-1].Sub(b.times[0]).Seconds()
			avg = total / float64(len(b.times)-1)
		}
		results = append(results, PatternResult{
			Port:            k.port,
			Protocol:        k.protocol,
			Action:          k.action,
			Count:           len(b.times),
			FirstSeen:       b.times[0],
			LastSeen:        b.times[len(b.times)-1],
			AvgIntervalSecs: avg,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Count != results[j].Count {
			return results[i].Count > results[j].Count
		}
		return results[i].Port < results[j].Port
	})
	return results
}
