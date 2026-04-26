package history

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// FlowCmd runs the port-flow analysis command.
type FlowCmd struct {
	HistoryPath string
	Protocol    string
	Action      string
	Since       time.Duration
	Window      time.Duration
	MinCount    int
	Format      string
	TopN        int
}

// Run loads history, computes flow edges, and renders results.
func (c *FlowCmd) Run() error {
	h, err := New(c.HistoryPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("load history: %w", err)
	}

	entries := h.Entries()
	if len(entries) == 0 {
		fmt.Fprintln(os.Stdout, "no history entries found")
		return nil
	}

	opts := FlowOptions{
		Protocol: strings.ToLower(c.Protocol),
		Action:   strings.ToLower(c.Action),
		Window:   c.Window,
		MinCount: c.MinCount,
	}
	if c.Since > 0 {
		opts.Since = time.Now().Add(-c.Since)
	}

	edges := BuildFlow(entries, opts)

	if c.TopN > 0 && len(edges) > c.TopN {
		edges = edges[:c.TopN]
	}

	return RenderFlow(edges, c.Format, os.Stdout)
}
