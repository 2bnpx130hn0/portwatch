package history

import (
	"fmt"
	"io"
	"os"
	"time"
)

// BaselineCmd manages baseline save/compare operations.
type BaselineCmd struct {
	HistoryPath  string
	BaselinePath string
	Format       string
	Out          io.Writer
}

// NewBaselineCmd creates a BaselineCmd with defaults.
func NewBaselineCmd(historyPath, baselinePath, format string) *BaselineCmd {
	return &BaselineCmd{
		HistoryPath:  historyPath,
		BaselinePath: baselinePath,
		Format:       format,
		Out:          os.Stdout,
	}
}

// Save writes current history as a new baseline.
func (c *BaselineCmd) Save() error {
	h, err := New(c.HistoryPath)
	if err != nil {
		return fmt.Errorf("load history: %w", err)
	}
	if err := SaveBaseline(h.Entries(), c.BaselinePath); err != nil {
		return fmt.Errorf("save baseline: %w", err)
	}
	fmt.Fprintf(c.Out, "baseline saved to %s (%d entries)\n", c.BaselinePath, len(h.Entries()))
	return nil
}

// Compare loads the baseline and diffs it against current history.
func (c *BaselineCmd) Compare() error {
	b, err := LoadBaseline(c.BaselinePath)
	if err != nil {
		return fmt.Errorf("load baseline: %w", err)
	}
	h, err := New(c.HistoryPath)
	if err != nil {
		return fmt.Errorf("load history: %w", err)
	}

	added, removed := BaselineDiff(b.Entries, h.Entries())

	w := c.Out
	if w == nil {
		w = os.Stdout
	}

	fmt.Fprintf(w, "baseline from %s\n", b.CreatedAt.Format(time.RFC3339))
	fmt.Fprintf(w, "added (%d):\n", len(added))
	for _, e := range added {
		fmt.Fprintf(w, "  + %s:%d [%s]\n", e.Protocol, e.Port, e.Action)
	}
	fmt.Fprintf(w, "removed (%d):\n", len(removed))
	for _, e := range removed {
		fmt.Fprintf(w, "  - %s:%d [%s]\n", e.Protocol, e.Port, e.Action)
	}
	return nil
}
