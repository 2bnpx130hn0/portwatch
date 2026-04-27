package history

import (
	"time"
)

// BurstResult represents a detected burst of activity for a port/protocol pair.
type BurstResult struct {
	Port     int
	Protocol string
	Action   string
	Count    int
	Window   time.Duration
	First    time.Time
	Last     time.Time
}

// BurstOptions configures burst detection.
type BurstOptions struct {
	// Window is the rolling time window to evaluate.
	Window time.Duration
	// Threshold is the minimum number of events in the window to qualify as a burst.
	Threshold int
	// Action filters entries by action (empty means all).
	Action string
	// Protocol filters entries by protocol (empty means all).
	Protocol string
	// Since ignores entries before this time.
	Since time.Time
}

// DetectBursts finds port/protocol pairs that exceed the event threshold
// within the configured rolling window.
func DetectBursts(entries []Entry, opts BurstOptions) []BurstResult {
	if opts.Window <= 0 {
		opts.Window = time.Hour
	}
	if opts.Threshold <= 0 {
		opts.Threshold = 5
	}

	type key struct {
		port     int
		proto    string
		action   string
	}

	groups := make(map[key][]time.Time)

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
		k := key{port: e.Port, proto: strings.ToLower(e.Protocol), action: strings.ToLower(e.Action)}
		groups[k] = append(groups[k], e.Timestamp)
	}

	var results []BurstResult
	for k, times := range groups {
		if len(times) < 2 {
			continue
		}
		// Sort ascending.
		sortTimes(times)
		// Sliding window count.
		maxCount := 0
		var winFirst, winLast time.Time
		for i := 0; i < len(times); i++ {
			count := 1
			for j := i + 1; j < len(times); j++ {
				if times[j].Sub(times[i]) <= opts.Window {
					count++
				} else {
					break
				}
			}
			if count > maxCount {
				maxCount = count
				winFirst = times[i]
				if i+count-1 < len(times) {
					winLast = times[i+count-1]
				}
			}
		}
		if maxCount >= opts.Threshold {
			results = append(results, BurstResult{
				Port:     k.port,
				Protocol: k.proto,
				Action:   k.action,
				Count:    maxCount,
				Window:   opts.Window,
				First:    winFirst,
				Last:     winLast,
			})
		}
	}
	return results
}

func sortTimes(ts []time.Time) {
	for i := 1; i < len(ts); i++ {
		for j := i; j > 0 && ts[j].Before(ts[j-1]); j-- {
			ts[j], ts[j-1] = ts[j-1], ts[j]
		}
	}
}
