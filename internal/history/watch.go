package history

import "time"

// WatchEvent represents a single port change event recorded during a watch cycle.
type WatchEvent struct {
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
	Note      string    `json:"note,omitempty"`
}

// WatchSummary aggregates events from a single watcher tick.
type WatchSummary struct {
	CycleAt  time.Time    `json:"cycle_at"`
	Added    []WatchEvent `json:"added"`
	Removed  []WatchEvent `json:"removed"`
	Alerted  int          `json:"alerted"`
	Allowed  int          `json:"allowed"`
}

// NewWatchSummary builds a WatchSummary from a slice of history entries recorded
// after a single scan cycle. Only entries whose timestamp is at or after cycleAt
// are included.
func NewWatchSummary(entries []Entry, cycleAt time.Time) WatchSummary {
	s := WatchSummary{CycleAt: cycleAt}
	for _, e := range entries {
		if e.Timestamp.Before(cycleAt) {
			continue
		}
		ev := WatchEvent{
			Port:      e.Port,
			Protocol:  e.Protocol,
			Action:    e.Action,
			Timestamp: e.Timestamp,
			Note:      e.Note,
		}
		switch e.Action {
		case "alert":
			s.Added = append(s.Added, ev)
			s.Alerted++
		case "allow":
			s.Added = append(s.Added, ev)
			s.Allowed++
		case "removed":
			s.Removed = append(s.Removed, ev)
		}
	}
	return s
}

// HasChanges returns true if the summary contains any added or removed events.
func (s WatchSummary) HasChanges() bool {
	return len(s.Added) > 0 || len(s.Removed) > 0
}
