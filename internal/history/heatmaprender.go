package history

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var dayAbbrev = [7]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

// RenderHeatmap writes the heatmap to w in the given format ("text" or "json").
// If w is nil, os.Stdout is used.
func RenderHeatmap(cells []HeatmapCell, format string, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	switch strings.ToLower(format) {
	case "json":
		renderHeatmapJSON(cells, w)
	default:
		renderHeatmapText(cells, w)
	}
}

func renderHeatmapText(cells []HeatmapCell, w io.Writer) {
	// Build a 7×24 grid indexed by [day][hour].
	grid := [7][24]int{}
	for _, c := range cells {
		grid[int(c.DayOfWeek)][c.Hour] = c.Count
	}

	// Header row: hours 0–23.
	fmt.Fprintf(w, "%-4s", "")
	for h := 0; h < 24; h++ {
		fmt.Fprintf(w, "%3d", h)
	}
	fmt.Fprintln(w)

	for d := 0; d < 7; d++ {
		fmt.Fprintf(w, "%-4s", dayAbbrev[d])
		for h := 0; h < 24; h++ {
			count := grid[d][h]
			if count == 0 {
				fmt.Fprintf(w, "%3s", ".")
			} else {
				fmt.Fprintf(w, "%3d", count)
			}
		}
		fmt.Fprintln(w)
	}
}

func renderHeatmapJSON(cells []HeatmapCell, w io.Writer) {
	type jsonCell struct {
		Day   string `json:"day"`
		Hour  int    `json:"hour"`
		Count int    `json:"count"`
	}
	out := make([]jsonCell, 0, len(cells))
	for _, c := range cells {
		out = append(out, jsonCell{
			Day:   time.Weekday(c.DayOfWeek).String(),
			Hour:  c.Hour,
			Count: c.Count,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(out)
}
