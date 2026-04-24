package history

import (
	"fmt"
	"io"
	"os"
	"time"
)

// PatternCmd detects and renders recurring port/protocol/action patterns from
// a history file.
type PatternCmd struct {
	HistoryPath    string
	MinOccurrences int
	Action         string
	Protocol       string
	SinceHours     int
	Format         string
	Writer         io.Writer
}

// NewPatternCmd returns a PatternCmd with sensible defaults.
func NewPatternCmd(historyPath string) *PatternCmd {
	return &PatternCmd{
		HistoryPath:    historyPath,
		MinOccurrences: 2,
		Format:         "text",
		Writer:         os.Stdout,
	}
}

// Run loads history, detects patterns, and writes results to the writer.
func (c *PatternCmd) Run() error {
	h, err := New(c.HistoryPath)
	if err != nil {
		return fmt.Errorf("pattern: load history: %w", err)
	}

	opts := PatternOptions{
		MinOccurrences: c.MinOccurrences,
		Action:         c.Action,
		Protocol:       c.Protocol,
	}
	if c.SinceHours > 0 {
		opts.Since = time.Now().Add(-time.Duration(c.SinceHours) * time.Hour)
	}

	results := DetectPatterns(h.Entries(), opts)

	w := c.Writer
	if w == nil {
		w = os.Stdout
	}

	if c.Format == "json" {
		return renderPatternJSON(w, results)
	}
	return renderPatternText(w, results)
}

func renderPatternText(w io.Writer, results []PatternResult) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "no recurring patterns detected")
		return err
	}
	fmt.Fprintf(w, "%-6s %-8s %-8s %6s %10s\n", "PORT", "PROTO", "ACTION", "COUNT", "AVG_INTERVAL")
	for _, r := range results {
		interval := "-"
		if r.AvgIntervalSecs > 0 {
			d := time.Duration(r.AvgIntervalSecs) * time.Second
			interval = d.Round(time.Second).String()
		}
		fmt.Fprintf(w, "%-6d %-8s %-8s %6d %10s\n",
			r.Port, r.Protocol, r.Action, r.Count, interval)
	}
	return nil
}

func renderPatternJSON(w io.Writer, results []PatternResult) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "[]")
		return err
	}
	fmt.Fprintln(w, "[")
	for i, r := range results {
		comma := ","
		if i == len(results)-1 {
			comma = ""
		}
		fmt.Fprintf(w, `  {"port":%d,"protocol":%q,"action":%q,"count":%d,"avg_interval_secs":%.2f}%s`+"\n",
			r.Port, r.Protocol, r.Action, r.Count, r.AvgIntervalSecs, comma)
	}
	fmt.Fprintln(w, "]")
	return nil
}
