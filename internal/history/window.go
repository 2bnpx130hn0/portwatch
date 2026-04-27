package history

import (
	"time"
)

// WindowOptions configures a sliding window aggregation.
type WindowOptions struct {
	Size     time.Duration
	Step     time.Duration
	Action   string
	Protocol string
	Since    time.Time
}

// WindowBucket holds aggregated counts for a time window.
type WindowBucket struct {
	Start  time.Time
	End    time.Time
	Count  int
	Ports  map[int]int
	Actions map[string]int
}

// Window performs a sliding window aggregation over history entries.
func Window(entries []Entry, opts WindowOptions) []WindowBucket {
	if len(entries) == 0 || opts.Size <= 0 {
		return nil
	}
	step := opts.Step
	if step <= 0 {
		step = opts.Size
	}

	filtered := make([]Entry, 0, len(entries))
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
		filtered = append(filtered, e)
	}
	if len(filtered) == 0 {
		return nil
	}

	var earliest, latest time.Time
	for i, e := range filtered {
		if i == 0 || e.Timestamp.Before(earliest) {
			earliest = e.Timestamp
		}
		if i == 0 || e.Timestamp.After(latest) {
			latest = e.Timestamp
		}
	}

	var buckets []WindowBucket
	for start := earliest; !start.After(latest); start = start.Add(step) {
		end := start.Add(opts.Size)
		b := WindowBucket{
			Start:   start,
			End:     end,
			Ports:   make(map[int]int),
			Actions: make(map[string]int),
		}
		for _, e := range filtered {
			if !e.Timestamp.Before(start) && e.Timestamp.Before(end) {
				b.Count++
				b.Ports[e.Port]++
				b.Actions[e.Action]++
			}
		}
		buckets = append(buckets, b)
	}
	return buckets
}
