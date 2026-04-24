package history

import (
	"time"
)

// HeatmapCell represents activity intensity for a given hour-of-day and day-of-week.
type HeatmapCell struct {
	DayOfWeek time.Weekday `json:"day_of_week"`
	Hour      int         `json:"hour"`
	Count     int         `json:"count"`
}

// HeatmapOptions controls how the heatmap is generated.
type HeatmapOptions struct {
	Since  time.Time
	Action string
	Proto  string
}

// Heatmap builds a 7×24 activity grid from history entries.
// Each cell holds the number of events for that (weekday, hour) pair.
func Heatmap(entries []Entry, opts HeatmapOptions) []HeatmapCell {
	counts := make(map[[2]int]int)

	for _, e := range entries {
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if opts.Action != "" && !strings.EqualFold(e.Action, opts.Action) {
			continue
		}
		if opts.Proto != "" && !strings.EqualFold(e.Protocol, opts.Proto) {
			continue
		}
		key := [2]int{int(e.Timestamp.Weekday()), e.Timestamp.Hour()}
		counts[key]++
	}

	var cells []HeatmapCell
	for day := 0; day < 7; day++ {
		for hour := 0; hour < 24; hour++ {
			key := [2]int{day, hour}
			cells = append(cells, HeatmapCell{
				DayOfWeek: time.Weekday(day),
				Hour:      hour,
				Count:     counts[key],
			})
		}
	}
	return cells
}
