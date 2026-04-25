package history

import (
	"sort"
	"time"
)

// Session represents a contiguous run of activity for a port/protocol pair.
type Session struct {
	Port     int       `json:"port"`
	Protocol string    `json:"protocol"`
	Action   string    `json:"action"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Duration time.Duration `json:"duration"`
	Count    int       `json:"count"`
}

// SessionOptions controls how sessions are built.
type SessionOptions struct {
	Gap      time.Duration // max gap between events to be considered same session
	Action   string
	Protocol string
	Since    time.Time
}

// BuildSessions groups history entries into sessions based on temporal proximity.
func BuildSessions(entries []Entry, opts SessionOptions) []Session {
	if opts.Gap <= 0 {
		opts.Gap = 5 * time.Minute
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

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp.Before(filtered[j].Timestamp)
	})

	type key struct {
		port  int
		proto string
	}

	open := map[key]*Session{}
	var sessions []Session

	for _, e := range filtered {
		k := key{e.Port, strings.ToLower(e.Protocol)}
		s, ok := open[k]
		if !ok || e.Timestamp.Sub(s.End) > opts.Gap {
			if ok {
				sessions = append(sessions, *s)
			}
			new := &Session{
				Port:     e.Port,
				Protocol: strings.ToLower(e.Protocol),
				Action:   e.Action,
				Start:    e.Timestamp,
				End:      e.Timestamp,
				Count:    1,
			}
			open[k] = new
		} else {
			s.End = e.Timestamp
			s.Count++
		}
	}

	for _, s := range open {
		s.Duration = s.End.Sub(s.Start)
		sessions = append(sessions, *s)
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Start.Before(sessions[j].Start)
	})

	return sessions
}
