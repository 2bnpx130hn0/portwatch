package history

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// VelocityCmd runs the velocity analysis command.
type VelocityCmd struct {
	HistoryPath string
	Window      time.Duration
	Lookback    time.Duration
	MinEvents   int
	Action      string
	Protocol    string
	Format      string
	Writer      io.Writer
}

// Run loads history and prints velocity results.
func (c *VelocityCmd) Run() error {
	w := c.Writer
	if w == nil {
		w = os.Stdout
	}

	h, err := New(c.HistoryPath)
	if err != nil {
		return fmt.Errorf("loading history: %w", err)
	}

	opts := VelocityOptions{
		WindowSize: c.Window,
		Lookback:   c.Lookback,
		MinEvents:  c.MinEvents,
		Action:     c.Action,
		Protocol:   c.Protocol,
	}

	results := Velocity(h.Entries(), opts)

	switch c.Format {
	case "json":
		return renderVelocityJSON(w, results)
	default:
		return renderVelocityText(w, results)
	}
}

func renderVelocityText(w io.Writer, results []VelocityEntry) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "no velocity data")
		return err
	}
	_, err := fmt.Fprintf(w, "%-8s %-8s %-8s %10s %10s %6s\n",
		"PORT", "PROTO", "ACTION", "RATE/HR", "DELTA/HR", "COUNT")
	if err != nil {
		return err
	}
	for _, r := range results {
		_, err = fmt.Fprintf(w, "%-8d %-8s %-8s %10.2f %10.2f %6d\n",
			r.Port, r.Protocol, r.Action, r.Rate, r.Delta, r.Count)
		if err != nil {
			return err
		}
	}
	return nil
}

func renderVelocityJSON(w io.Writer, results []VelocityEntry) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
