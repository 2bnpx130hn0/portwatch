package history

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// CorrelateCmd runs a correlation analysis over a history file.
type CorrelateCmd struct {
	Path       string
	Window     time.Duration
	MinEntries int
	Action     string
	Format     string
	Writer     io.Writer
}

// NewCorrelateCmd returns a CorrelateCmd with sensible defaults.
func NewCorrelateCmd(path string) *CorrelateCmd {
	return &CorrelateCmd{
		Path:       path,
		Window:     5 * time.Minute,
		MinEntries: 2,
		Format:     "text",
		Writer:     os.Stdout,
	}
}

// Run loads history, runs correlation, and prints results.
func (c *CorrelateCmd) Run() error {
	h, err := New(c.Path)
	if err != nil {
		return fmt.Errorf("load history: %w", err)
	}

	opts := CorrelateOptions{
		Window:     c.Window,
		MinEntries: c.MinEntries,
		Action:     c.Action,
	}

	results := Correlate(h.Entries(), opts)

	w := c.Writer
	if w == nil {
		w = os.Stdout
	}

	if c.Format == "json" {
		return c.writeJSON(w, results)
	}
	return c.writeText(w, results)
}

func (c *CorrelateCmd) writeText(w io.Writer, results []Correlation) error {
	if len(results) == 0 {
		fmt.Fprintln(w, "no correlations found")
		return nil
	}
	for _, r := range results {
		fmt.Fprintf(w, "port=%-6d protocol=%-5s entries=%d score=%.0f\n",
			r.Port, r.Protocol, len(r.Entries), r.Score)
		for _, e := range r.Entries {
			fmt.Fprintf(w, "  [%s] action=%s\n",
				e.Timestamp.Format(time.RFC3339), e.Action)
		}
	}
	return nil
}

func (c *CorrelateCmd) writeJSON(w io.Writer, results []Correlation) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
