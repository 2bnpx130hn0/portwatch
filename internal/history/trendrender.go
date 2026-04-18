package history

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// RenderTrend writes trend points in the given format.
func RenderTrend(points []TrendPoint, format string, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}
	if format == "" {
		format = "text"
	}
	switch strings.ToLower(format) {
	case "json":
		return renderTrendJSON(points, w)
	default:
		return renderTrendText(points, w)
	}
}

func renderTrendText(points []TrendPoint, w io.Writer) error {
	if len(points) == 0 {
		_, err := fmt.Fprintln(w, "no trend data")
		return err
	}
	for _, p := range points {
		_, err := fmt.Fprintf(w, "%s  %d\n", p.Window.Format("2006-01-02 15:04"), p.Count)
		if err != nil {
			return err
		}
	}
	return nil
}

func renderTrendJSON(points []TrendPoint, w io.Writer) error {
	type jsonPoint struct {
		Window string `json:"window"`
		Count  int    `json:"count"`
	}
	out := make([]jsonPoint, len(points))
	for i, p := range points {
		out[i] = jsonPoint{Window: p.Window.Format("2006-01-02T15:04:05Z07:00"), Count: p.Count}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
