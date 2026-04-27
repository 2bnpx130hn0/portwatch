package history

import (
	"fmt"
	"io"
	"os"
	"time"
)

// WindowCmd runs the sliding window aggregation command.
type WindowCmd struct {
	HistoryFile string
	Size        time.Duration
	Step        time.Duration
	Action      string
	Protocol    string
	Since       time.Time
	Format      string
	Out         io.Writer
}

// Run loads history, applies the window aggregation, and renders results.
func (c *WindowCmd) Run() error {
	if c.Out == nil {
		c.Out = os.Stdout
	}
	if c.HistoryFile == "" {
		return fmt.Errorf("history file path is required")
	}
	h, err := New(c.HistoryFile)
	if err != nil {
		return fmt.Errorf("loading history: %w", err)
	}
	entries := h.Entries()
	if len(entries) == 0 {
		fmt.Fprintln(c.Out, "no history entries found")
		return nil
	}
	size := c.Size
	if size <= 0 {
		size = time.Hour
	}
	buckets := Window(entries, WindowOptions{
		Size:     size,
		Step:     c.Step,
		Action:   c.Action,
		Protocol: c.Protocol,
		Since:    c.Since,
	})
	RenderWindow(buckets, c.Format, c.Out)
	return nil
}

// NewWindowCmd constructs a WindowCmd with defaults.
func NewWindowCmd(historyFile string) *WindowCmd {
	return &WindowCmd{
		HistoryFile: historyFile,
		Size:        time.Hour,
		Format:      "text",
	}
}
