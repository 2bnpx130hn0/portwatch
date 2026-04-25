package history

import (
	"math"
	"sort"
	"time"
)

// VelocityEntry represents the rate of change for a port/protocol pair.
type VelocityEntry struct {
	Port     int
	Protocol string
	Action   string
	Rate     float64 // events per hour
	Delta    float64 // change in rate compared to previous window
	Count    int
}

// VelocityOptions controls how velocity is computed.
type VelocityOptions struct {
	WindowSize time.Duration // size of each window (default: 1h)
	Lookback   time.Duration // total lookback period (default: 24h)
	MinEvents  int           // minimum events to include (default: 2)
	Action     string        // optional action filter
	Protocol   string        // optional protocol filter
}

// Velocity computes the rate of change of port events over time windows.
func Velocity(entries []Entry, opts VelocityOptions) []VelocityEntry {
	if opts.WindowSize == 0 {
		opts.WindowSize = time.Hour
	}
	if opts.Lookback == 0 {
		opts.Lookback = 24 * time.Hour
	}
	if opts.MinEvents == 0 {
		opts.MinEvents = 2
	}

	now := time.Now()
	cutoff := now.Add(-opts.Lookback)

	type key struct {
		port     int
		proto    string
		action   string
	}

	type windowCounts struct {
		current  int
		previous int
	}

	currentStart := now.Add(-opts.WindowSize)
	previousStart := currentStart.Add(-opts.WindowSize)

	buckets := make(map[key]*windowCounts)

	for _, e := range entries {
		if e.Timestamp.Before(cutoff) {
			continue
		}
		if opts.Action != "" && !strings.EqualFold(e.Action, opts.Action) {
			continue
		}
		if opts.Protocol != "" && !strings.EqualFold(e.Protocol, opts.Protocol) {
			continue
		}
		k := key{port: e.Port, proto: strings.ToLower(e.Protocol), action: strings.ToLower(e.Action)}
		if buckets[k] == nil {
			buckets[k] = &windowCounts{}
		}
		if !e.Timestamp.Before(currentStart) {
			buckets[k].current++
		} else if !e.Timestamp.Before(previousStart) {
			buckets[k].previous++
		}
	}

	var results []VelocityEntry
	hours := opts.WindowSize.Hours()
	if hours == 0 {
		hours = 1
	}

	for k, w := range buckets {
		total := w.current + w.previous
		if total < opts.MinEvents {
			continue
		}
		rate := float64(w.current) / hours
		prevRate := float64(w.previous) / hours
		delta := rate - prevRate
		results = append(results, VelocityEntry{
			Port:     k.port,
			Protocol: k.proto,
			Action:   k.action,
			Rate:     math.Round(rate*100) / 100,
			Delta:    math.Round(delta*100) / 100,
			Count:    total,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if math.Abs(results[i].Delta) != math.Abs(results[j].Delta) {
			return math.Abs(results[i].Delta) > math.Abs(results[j].Delta)
		}
		return results[i].Port < results[j].Port
	})

	return results
}
