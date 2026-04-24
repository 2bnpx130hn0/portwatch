package history

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// ForecastCmd runs a forecast over a history file and renders results.
type ForecastCmd struct {
	Path       string
	Protocol   string
	Action     string
	Steps      int
	BucketSize time.Duration
	Since      time.Duration
	Format     string
	Out        io.Writer
}

func NewForecastCmd(path string) *ForecastCmd {
	return &ForecastCmd{
		Path:       path,
		Steps:      3,
		BucketSize: time.Hour,
		Format:     "text",
		Out:        os.Stdout,
	}
}

func (c *ForecastCmd) Run() error {
	h, err := New(c.Path)
	if err != nil {
		return fmt.Errorf("forecast: load history: %w", err)
	}

	opts := ForecastOptions{
		Protocol:   c.Protocol,
		Action:     c.Action,
		Steps:      c.Steps,
		BucketSize: c.BucketSize,
	}
	if c.Since > 0 {
		opts.Since = time.Now().Add(-c.Since)
	}

	results := Forecast(h.Entries(), opts)

	switch c.Format {
	case "json":
		return c.renderJSON(results)
	default:
		return c.renderText(results)
	}
}

func (c *ForecastCmd) renderText(results []ForecastResult) error {
	w := c.Out
	if w == nil {
		w = os.Stdout
	}
	if len(results) == 0 {
		fmt.Fprintln(w, "no forecast data available")
		return nil
	}
	fmt.Fprintf(w, "%-25s  %-8s  %-8s  %s\n", "BUCKET", "PROTOCOL", "ACTION", "PREDICTED")
	for _, r := range results {
		proto := r.Protocol
		if proto == "" {
			proto = "any"
		}
		action := r.Action
		if action == "" {
			action = "any"
		}
		fmt.Fprintf(w, "%-25s  %-8s  %-8s  %.2f\n",
			r.Bucket.Format(time.RFC3339), proto, action, r.Predicted)
	}
	return nil
}

func (c *ForecastCmd) renderJSON(results []ForecastResult) error {
	w := c.Out
	if w == nil {
		w = os.Stdout
	}
	return json.NewEncoder(w).Encode(results)
}
